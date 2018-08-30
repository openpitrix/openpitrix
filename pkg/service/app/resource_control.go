// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package app

import (
	"bytes"
	"context"
	"fmt"
	"time"

	repoclient "openpitrix.io/openpitrix/pkg/client/repo"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/devkit"
	"openpitrix.io/openpitrix/pkg/devkit/app"
	"openpitrix.io/openpitrix/pkg/gerr"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/pi"
	"openpitrix.io/openpitrix/pkg/repoiface"
	"openpitrix.io/openpitrix/pkg/repoiface/indexer"
	"openpitrix.io/openpitrix/pkg/service/category/categoryutil"
	"openpitrix.io/openpitrix/pkg/util/stringutil"
)

type Action int32

const (
	Submit Action = iota
	Cancel
	Release
	Delete
	Pass
	Reject
	Suspend
	Recover
	Modify
)

// Action => []version.status
var VersionFiniteStatusMachine = map[Action][]string{
	// TODO: only admin can modify active app version
	Modify:  {constants.StatusDraft, constants.StatusRejected, constants.StatusActive},
	Submit:  {constants.StatusDraft, constants.StatusRejected},
	Cancel:  {constants.StatusSubmitted},
	Release: {constants.StatusPassed},
	Delete: {
		constants.StatusSuspended,
		constants.StatusDraft,
		constants.StatusPassed,
		constants.StatusRejected,
	},
	Pass:    {constants.StatusSubmitted, constants.StatusRejected},
	Reject:  {constants.StatusSubmitted},
	Suspend: {constants.StatusActive},
	Recover: {constants.StatusSuspended},
}

func checkAppVersionHandlePermission(
	ctx context.Context, action Action, versionId string,
) (*models.AppVersion, error) {
	// TODO: check admin/developer permission
	//sender := senderutil.GetSenderFromContext(ctx)
	version := models.AppVersion{}
	err := pi.Global().DB(ctx).
		Select(models.AppVersionColumns...).
		From(models.AppVersionTableName).
		Where(db.Eq(models.ColumnVersionId, versionId)).
		LoadOne(&version)
	if err != nil {
		if err == db.ErrNotFound {
			return nil, gerr.New(ctx, gerr.NotFound, gerr.ErrorResourceNotFound, versionId)
		}
		return nil, gerr.New(ctx, gerr.Internal, gerr.ErrorInternalError)
	}

	allowedStatus, ok := VersionFiniteStatusMachine[action]
	if !ok {
		return nil, gerr.New(ctx, gerr.Internal, gerr.ErrorInternalError)
	}
	versionStatus := version.Status
	if !stringutil.StringIn(versionStatus, allowedStatus) {
		return nil, gerr.New(ctx, gerr.FailedPrecondition,
			gerr.ErrorAppVersionIncorrectStatus, version.VersionId, versionStatus)
	}
	return &version, nil
}

func insertVersion(ctx context.Context, version *models.AppVersion) error {
	_, err := pi.Global().DB(ctx).
		InsertInto(models.AppVersionTableName).
		Columns(models.AppVersionColumns...).
		Record(version).
		Exec()
	if err != nil {
		logger.Error(ctx, "Failed to insert version [%+v]", version)
		return gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
	}
	return nil
}

func updateVersion(ctx context.Context, versionId string, attributes map[string]interface{}) error {
	attributes[models.ColumnUpdateTime] = time.Now()
	_, err := pi.Global().DB(ctx).
		Update(models.AppVersionTableName).
		SetMap(attributes).
		Where(db.Eq(models.ColumnVersionId, versionId)).
		Exec()
	if err != nil {
		return gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorModifyResourceFailed, versionId)
	}
	return nil
}

func updateApp(ctx context.Context, appId string, attributes map[string]interface{}) error {
	attributes[models.ColumnUpdateTime] = time.Now()
	_, err := pi.Global().DB(ctx).
		Update(models.AppTableName).
		SetMap(attributes).
		Where(db.Eq(models.ColumnAppId, appId)).
		Exec()
	if err != nil {
		return gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorModifyResourcesFailed)
	}
	return nil
}

func updateVersionStatus(ctx context.Context, version *models.AppVersion, status string) error {
	err := updateVersion(ctx, version.VersionId, map[string]interface{}{
		models.ColumnStatus:     status,
		models.ColumnStatusTime: time.Now(),
	})
	if err != nil {
		return err
	}
	err = syncAppStatus(ctx, version.AppId)
	if err != nil {
		return err
	}
	return nil
}

func syncAppStatus(ctx context.Context, appId string) error {
	app, err := getApp(ctx, appId)
	if err != nil {
		return err
	}
	if app.Status == constants.StatusSuspended ||
		app.Status == constants.StatusDeleted {
		return nil
	}
	var status string
	countActive, err := pi.Global().DB(ctx).
		Select(models.AppVersionColumns...).
		From(models.AppVersionTableName).
		Where(db.Eq(models.ColumnAppId, appId)).
		Where(db.Eq(models.ColumnStatus, constants.StatusActive)).
		Count()
	if err != nil {
		return gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}
	if countActive > 0 {
		status = constants.StatusActive
	} else {
		status = constants.StatusDraft
	}

	_, err = pi.Global().DB(ctx).
		Update(models.AppTableName).
		Set(models.ColumnStatus, status).
		Set(models.ColumnStatusTime, time.Now()).
		Where(db.Eq(models.ColumnAppId, appId)).
		Exec()
	if err != nil {
		return gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorModifyResourceFailed, appId)
	}

	return nil
}

func getApp(ctx context.Context, appId string) (*models.App, error) {
	app := &models.App{}
	err := pi.Global().DB(ctx).
		Select(models.AppColumns...).
		From(models.AppTableName).
		Where(db.Eq(models.ColumnAppId, appId)).
		LoadOne(&app)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourceFailed, appId)
	}
	return app, nil
}

func getLatestAppVersion(ctx context.Context, appId string) (*models.AppVersion, error) {
	appVersion := &models.AppVersion{}
	err := pi.Global().DB(ctx).
		Select(models.AppVersionColumns...).
		From(models.AppVersionTableName).
		Where(db.Eq(models.ColumnAppId, appId)).
		OrderDir(models.ColumnSequence, false).
		LoadOne(&appVersion)
	if err != nil {
		if err == db.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}
	return appVersion, nil
}

func formatApp(ctx context.Context, app *models.App) (*pb.App, error) {
	pbApp := models.AppToPb(app)

	latestAppVersion, err := getLatestAppVersion(ctx, app.AppId)
	if err != nil {
		return nil, err
	}
	pbApp.LatestAppVersion = models.AppVersionToPb(latestAppVersion)

	return pbApp, nil
}

func formatAppSet(ctx context.Context, apps []*models.App) ([]*pb.App, error) {
	var pbApps []*pb.App
	var appIds []string
	for _, app := range apps {
		var pbApp *pb.App
		pbApp, err := formatApp(ctx, app)
		if err != nil {
			return pbApps, err
		}
		appIds = append(appIds, app.AppId)
		pbApps = append(pbApps, pbApp)
	}
	rcmap, err := categoryutil.GetResourcesCategories(ctx, pi.Global().DB(ctx), appIds)
	if err != nil {
		return pbApps, err
	}
	for _, pbApp := range pbApps {
		if categorySet, ok := rcmap[pbApp.GetAppId().GetValue()]; ok {
			pbApp.CategorySet = categorySet
		}
	}
	return pbApps, nil
}

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
	riface, err := vp.GetRepoInterface(ctx)
	if err != nil {
		return nil, err
	}
	return riface.ReadFile(ctx, vp.version.PackageName)
}

func (vp *versionProxy) GetRepoInterface(ctx context.Context) (repoiface.RepoInterface, error) {
	repo, err := vp.GetRepo(ctx)
	if err != nil {
		return nil, err
	}
	return repoiface.New(ctx, repo.Type.GetValue(), repo.Url.GetValue(), repo.Credential.GetValue())
}

func (vp *versionProxy) ModifyPackageFile(ctx context.Context, newPackage []byte) error {
	riface, err := vp.GetRepoInterface(ctx)
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

	appVersion := &app.Version{}
	appVersion.Metadata = pkg.Metadata
	pkgVersion := indexer.AppVersionWrapper{Version: appVersion}
	appVersionName := pkgVersion.GetVersion()
	if pkgVersion.GetAppVersion() != "" {
		appVersionName += fmt.Sprintf(" [%s]", pkgVersion.GetAppVersion())
	}
	packageName := fmt.Sprintf("%s-%s.tgz", appVersion.Name, appVersion.Version)

	err = updateApp(ctx, appId, map[string]interface{}{
		models.ColumnName:        pkgVersion.GetName(),
		models.ColumnDescription: pkgVersion.GetDescription(),
		models.ColumnChartName:   pkgVersion.GetName(),
		models.ColumnHome:        pkgVersion.GetHome(),
		models.ColumnIcon:        pkgVersion.GetIcon(),
		models.ColumnScreenshots: pkgVersion.GetScreenshots(),
		models.ColumnSources:     pkgVersion.GetSources(),
		models.ColumnKeywords:    pkgVersion.GetKeywords(),
	})
	if err != nil {
		return err
	}

	err = updateVersion(ctx, versionId, map[string]interface{}{
		models.ColumnName:        appVersionName,
		models.ColumnDescription: pkgVersion.GetDescription(),
		models.ColumnPackageName: packageName,
	})
	if err != nil {
		return err
	}

	err = riface.WriteFile(ctx, packageName, newPackage)
	if err != nil {
		logger.Error(ctx, "Failed to write [%s] package, error: %+v", vp.version.VersionId, err)
		return gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorCreateResourceFailed, vp.version.VersionId)
	}
	return nil
}
