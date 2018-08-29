// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package app

import (
	"context"
	"fmt"
	"time"

	repoclient "openpitrix.io/openpitrix/pkg/client/repo"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/gerr"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/pi"
	"openpitrix.io/openpitrix/pkg/repoiface"
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

func updateVersionStatus(ctx context.Context, version *models.AppVersion, status string) error {
	_, err := pi.Global().DB(ctx).
		Update(models.AppVersionTableName).
		Set(models.ColumnStatus, status).
		Set(models.ColumnStatusTime, time.Now()).
		Where(db.Eq(models.ColumnVersionId, version.VersionId)).
		Exec()
	if err != nil {
		return gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorModifyResourceFailed, version.VersionId)
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

func getRepo(ctx context.Context, repoId string) (*pb.Repo, error) {
	rc, err := repoclient.NewRepoManagerClient()
	if err != nil {
		return nil, err
	}
	describeResp, err := rc.DescribeRepos(ctx, &pb.DescribeReposRequest{
		RepoId: []string{repoId},
	})
	if err != nil {
		return nil, err
	}
	if len(describeResp.RepoSet) == 0 {
		return nil, fmt.Errorf("repo [%s] not exists", repoId)
	}

	return describeResp.RepoSet[0], nil
}

func getPackageFile(ctx context.Context, version *models.AppVersion) ([]byte, error) {
	riface, err := getRepoInterface(ctx, version)
	if err != nil {
		return nil, err
	}
	return riface.ReadFile(ctx, version.PackageName)
}

func getRepoInterface(ctx context.Context, version *models.AppVersion) (repoiface.RepoInterface, error) {
	app, err := getApp(ctx, version.AppId)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourceFailed, version.AppId)
	}
	repo, err := getRepo(ctx, app.RepoId)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourceFailed, app.RepoId)
	}
	return repoiface.New(ctx, repo.Type.GetValue(), repo.Url.GetValue(), repo.Credential.GetValue())
}

func modifyPackageFile(ctx context.Context, version *models.AppVersion, newPackage []byte) error {
	riface, err := getRepoInterface(ctx, version)
	if err != nil {
		return err
	}
	return riface.WriteFile(ctx, version.PackageName, newPackage)
}
