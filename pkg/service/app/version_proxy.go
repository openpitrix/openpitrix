// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package app

import (
	"bytes"
	"context"

	clientutil "openpitrix.io/openpitrix/pkg/client"
	repoclient "openpitrix.io/openpitrix/pkg/client/repo"
	"openpitrix.io/openpitrix/pkg/gerr"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/repoiface"
	"openpitrix.io/openpitrix/pkg/util/gziputil"
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

func (vp *versionProxy) GetPackageFile(ctx context.Context, includeFiles ...string) (map[string][]byte, error) {
	rreader, err := vp.GetRepoReader(ctx)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourceFailed, vp.version.VersionId)
	}
	content, err := rreader.ReadFile(ctx, vp.version.PackageName)
	if err != nil {
		logger.Error(ctx, "Failed to read [%s] package, error: %+v", vp.version.VersionId, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourceFailed, vp.version.VersionId)
	}

	archiveFiles, err := gziputil.LoadArchive(bytes.NewReader(content), includeFiles...)
	if err != nil {
		logger.Error(ctx, "Failed to load [%s] package, error: %+v", vp.version.VersionId, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourceFailed, vp.version.VersionId)
	}
	return archiveFiles, err
}

func (vp *versionProxy) GetRepoReader(ctx context.Context) (*repoiface.Reader, error) {
	repo, err := vp.GetRepo(ctx)
	if err != nil {
		return nil, err
	}
	return repoiface.NewReader(ctx, repo)
}
