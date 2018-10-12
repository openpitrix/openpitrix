// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package app

import (
	"context"

	clientutil "openpitrix.io/openpitrix/pkg/client"
	repoclient "openpitrix.io/openpitrix/pkg/client/repo"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/gerr"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/repoiface"
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
	systemCtx := clientutil.SetSystemUserToContext(ctx)
	describeResp, err := rc.DescribeRepos(systemCtx, &pb.DescribeReposRequest{
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

func (vp *versionProxy) AddPackageFile(ctx context.Context, newPackage []byte, syncAppName bool, replacePackage bool) error {
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

	pkgVersion, err := rreader.LoadPackage(ctx, newPackage)
	if err != nil {
		logger.Error(ctx, "Failed to load package, error: %+v", err)
		return gerr.NewWithDetail(ctx, gerr.InvalidArgument, err, gerr.ErrorLoadPackageFailed, err.Error())
	}

	appVersions, err := rreader.GetAppVersions(ctx)
	if err != nil {
		return gerr.NewWithDetail(ctx, gerr.PermissionDenied, err, gerr.ErrorCannotAccessRepo)
	}

	appName := pkgVersion.GetName()
	version := pkgVersion.GetVersion()
	if app.Name != appName ||
		app.ChartName != appName {
		if !syncAppName {
			// cannot change app name
			return gerr.New(ctx, gerr.InvalidArgument, gerr.ErrorCannotChangeAppName)
		} else {
			// check app name in index.yaml
			if _, ok := appVersions[appName]; ok {
				return gerr.New(ctx, gerr.InvalidArgument, gerr.ErrorAppNameExists, appName)
			}

			err = updateApp(ctx, appId, map[string]interface{}{
				constants.ColumnName:      appName,
				constants.ColumnChartName: appName,
			})
			if err != nil {
				return err
			}
		}
	}

	if !replacePackage || version != vp.version.GetSemver() {
		if versionMap, ok := appVersions[appName]; ok {
			if appVersion, ok := versionMap[version]; ok {
				return gerr.New(ctx, gerr.InvalidArgument, gerr.ErrorAppVersionExists, version, appVersion)
			}
		}
	}

	err = updateVersion(ctx, versionId, map[string]interface{}{
		constants.ColumnName:        pkgVersion.GetVersionName(),
		constants.ColumnDescription: pkgVersion.GetDescription(),
		constants.ColumnPackageName: pkgVersion.GetPackageName(),
		constants.ColumnHome:        pkgVersion.GetHome(),
		constants.ColumnIcon:        pkgVersion.GetIcon(),
		constants.ColumnScreenshots: pkgVersion.GetScreenshots(),
		constants.ColumnSources:     pkgVersion.GetSources(),
		constants.ColumnKeywords:    pkgVersion.GetKeywords(),
	})
	if err != nil {
		return err
	}

	err = resortAppVersions(ctx, appId)
	if err != nil {
		return err
	}

	if replacePackage {
		err = rreader.DeletePackage(ctx, appName, vp.version.GetSemver())
		if err != nil {
			logger.Error(ctx, "Failed to delete [%s] old package, error: %+v", vp.version.VersionId, err)
			return gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDeleteResourcesFailed)
		}
	}

	err = rreader.CheckWrite(ctx)
	if err != nil {
		return gerr.NewWithDetail(ctx, gerr.PermissionDenied, err, gerr.ErrorCannotWriteRepo, vp.repo.RepoId)
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

func (vp *versionProxy) DeletePackageFile(ctx context.Context) error {
	rreader, err := vp.GetRepoReader(ctx)
	if err != nil {
		return err
	}
	app, err := vp.GetApp(ctx)
	if err != nil {
		return err
	}
	err = rreader.DeletePackage(ctx, app.ChartName, vp.version.GetSemver())
	if err != nil {
		logger.Error(ctx, "Failed to delete version [%s] from repo, error: %+v", vp.version.VersionId, err)
		return gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDeleteResourceFailed, vp.version.VersionId)
	}
	return nil
}
