// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package app

import (
	"bytes"
	"context"

	repoclient "openpitrix.io/openpitrix/pkg/client/repo"
	"openpitrix.io/openpitrix/pkg/devkit"
	"openpitrix.io/openpitrix/pkg/devkit/opapp"
	"openpitrix.io/openpitrix/pkg/gerr"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/repoiface"
	"openpitrix.io/openpitrix/pkg/repoiface/wrapper"
	"openpitrix.io/openpitrix/pkg/util/stringutil"
)

type versionProxy struct {
	version *models.AppVersion
	app     *models.App
	repo    *pb.Repo
}

func newVersionProxy(version *models.AppVersion) *versionProxy {
	vp := &versionProxy{}
	vp.version = version
	return vp
}

func (vp *versionProxy) GetVersion(ctx context.Context) (*models.AppVersion, error) {
	return vp.version, nil
}

func (vp *versionProxy) GetApp(ctx context.Context) (*models.App, error) {
	if vp.app != nil {
		return vp.app, nil
	}
	app, err := getApp(ctx, vp.version.AppId)
	if err != nil {
		return nil, err
	}
	vp.app = app
	return app, nil
}

func (vp *versionProxy) GetRepo(ctx context.Context) (*pb.Repo, error) {
	if vp.repo != nil {
		return vp.repo, nil
	}
	app, err := vp.GetApp(ctx)
	if err != nil {
		return nil, err
	}
	repoId := app.RepoId
	rc, err := repoclient.NewRepoManagerClient()
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourceFailed, app.RepoId)
	}
	describeResp, err := rc.DescribeRepos(ctx, &pb.DescribeReposRequest{
		RepoId: []string{repoId},
	})
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourceFailed, app.RepoId)
	}
	if len(describeResp.RepoSet) == 0 {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourceFailed, app.RepoId)
	}
	vp.repo = describeResp.RepoSet[0]
	return vp.repo, nil
}

func (vp *versionProxy) GetPackageFile(ctx context.Context) ([]byte, error) {
	rreader, err := vp.GetRepoReader(ctx)
	if err != nil {
		return nil, err
	}
	return rreader.ReadFile(ctx, vp.version.PackageName)
}

func (vp *versionProxy) GetRepoReader(ctx context.Context) (*repoiface.Reader, error) {
	repo, err := vp.GetRepo(ctx)
	if err != nil {
		return nil, err
	}
	return repoiface.NewReader(ctx, repo)
}

func (vp *versionProxy) ModifyPackageFile(ctx context.Context, newPackage []byte, syncAppName bool) error {
	rreader, err := vp.GetRepoReader(ctx)
	if err != nil {
		return err
	}
	app, err := vp.GetApp(ctx)
	if err != nil {
		return err
	}
	appId := vp.version.AppId
	versionId := vp.version.VersionId

	pkg, err := devkit.LoadArchive(bytes.NewReader(newPackage))
	if err != nil {
		logger.Error(ctx, "Failed to load package, error: %+v", err)
		return gerr.NewWithDetail(ctx, gerr.InvalidArgument, err, gerr.ErrorLoadPackageFailed, err.Error())
	}

	pkgVersion := wrapper.OpVersionWrapper{OpVersion: &opapp.OpVersion{Metadata: pkg.Metadata}}
	appName := pkgVersion.GetName()
	if app.Name != appName ||
		app.ChartName != appName {
		if !syncAppName {
			// cannot change app name
			return gerr.New(ctx, gerr.InvalidArgument, gerr.ErrorCannotChangeAppName)
		} else {
			// check app name in index.yaml
			packagesName, err := rreader.GetPackagesName(ctx)
			if err != nil {
				return gerr.NewWithDetail(ctx, gerr.PermissionDenied, err, gerr.ErrorCannotAccessRepo)
			}
			if stringutil.StringIn(appName, packagesName) {
				return gerr.New(ctx, gerr.InvalidArgument, gerr.ErrorAppNameExists)
			}

			err = updateApp(ctx, appId, map[string]interface{}{
				models.ColumnName:      appName,
				models.ColumnChartName: appName,
			})
			if err != nil {
				return err
			}
		}
	}

	err = updateVersion(ctx, versionId, map[string]interface{}{
		models.ColumnName:        pkgVersion.GetVersionName(),
		models.ColumnDescription: pkgVersion.GetDescription(),
		models.ColumnPackageName: pkgVersion.GetPackageName(),
		models.ColumnHome:        pkgVersion.GetHome(),
		models.ColumnIcon:        pkgVersion.GetIcon(),
		models.ColumnScreenshots: pkgVersion.GetScreenshots(),
		models.ColumnSources:     pkgVersion.GetSources(),
		models.ColumnKeywords:    pkgVersion.GetKeywords(),
	})
	if err != nil {
		return err
	}

	err = resortAppVersions(ctx, appId)
	if err != nil {
		return err
	}

	err = rreader.WriteFile(ctx, pkgVersion.GetPackageName(), newPackage)
	if err != nil {
		logger.Error(ctx, "Failed to write [%s] package, error: %+v", vp.version.VersionId, err)
		return gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorCreateResourceFailed, vp.version.VersionId)
	}
	err = rreader.AddPackage(ctx, newPackage)
	if err != nil {
		logger.Error(ctx, "Failed to add version [%s] into index.yaml, error: %+v", vp.version.VersionId, err)
		return gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorCreateResourceFailed, vp.version.VersionId)
	}
	return nil
}
