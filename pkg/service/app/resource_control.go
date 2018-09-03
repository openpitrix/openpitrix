// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package app

import (
	"context"
	"sort"
	"time"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/gerr"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/pi"
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
	attributes := make(map[string]interface{})
	activeVersion, err := getLatestAppVersion(ctx, appId, constants.StatusActive)
	if err != nil {
		return err
	}
	var status string
	if activeVersion != nil {
		status = constants.StatusActive

		if activeVersion.Description != app.Description {
			attributes[models.ColumnDescription] = activeVersion.Description
		}
		if activeVersion.Home != app.Home {
			attributes[models.ColumnHome] = activeVersion.Home
		}
		if activeVersion.Icon != app.Icon {
			attributes[models.ColumnIcon] = activeVersion.Icon
		}
		if activeVersion.Screenshots != app.Screenshots {
			attributes[models.ColumnScreenshots] = activeVersion.Screenshots
		}
		if activeVersion.Maintainers != app.Maintainers {
			attributes[models.ColumnMaintainers] = activeVersion.Maintainers
		}
		if activeVersion.Keywords != app.Keywords {
			attributes[models.ColumnKeywords] = activeVersion.Keywords
		}
		if activeVersion.Sources != app.Sources {
			attributes[models.ColumnSources] = activeVersion.Sources
		}
		if activeVersion.Readme != app.Readme {
			attributes[models.ColumnReadme] = activeVersion.Readme
		}
	} else {
		status = constants.StatusDraft
	}
	if status != app.Status {
		attributes[models.ColumnStatus] = status
		attributes[models.ColumnStatusTime] = time.Now()
	}
	if len(attributes) == 0 {
		return nil
	}
	attributes[models.ColumnUpdateTime] = time.Now()

	_, err = pi.Global().DB(ctx).
		Update(models.AppTableName).
		SetMap(attributes).
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

func getLatestAppVersion(ctx context.Context, appId string, status ...string) (*models.AppVersion, error) {
	appVersion := &models.AppVersion{}
	stmt := pi.Global().DB(ctx).
		Select(models.AppVersionColumns...).
		From(models.AppVersionTableName).
		Where(db.Eq(models.ColumnAppId, appId))
	if len(status) > 0 {
		stmt.Where(db.Eq(models.ColumnStatus, status))
	}
	err := stmt.
		OrderDir(models.ColumnSequence, false).
		LoadOne(&appVersion)
	if err != nil {
		if err == db.ErrNotFound {
			return nil, nil
		}
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
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

func getBigestSequence(ctx context.Context, appId string) (uint32, error) {
	var sequence uint32
	err := pi.Global().DB(ctx).
		Select("coalesce(max(sequence), 0)").
		From(models.AppVersionTableName).
		Where(db.Eq(models.ColumnAppId, appId)).
		LoadOne(&sequence)
	if err != nil {
		return sequence, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourceFailed, appId)
	}
	return sequence, nil
}

func resortAppVersions(ctx context.Context, appId string) error {
	var versions models.AppVersions
	_, err := pi.Global().DB(ctx).
		Select(models.ColumnVersionId, models.ColumnName, models.ColumnSequence, models.ColumnCreateTime).
		From(models.AppVersionTableName).
		Where(db.Eq(models.ColumnAppId, appId)).
		Load(&versions)
	if err != nil {
		return gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourceFailed, appId)
	}
	sort.Sort(versions)
	for i, version := range versions {
		if version.Sequence != uint32(i) {
			err = updateVersion(ctx, version.VersionId, map[string]interface{}{
				models.ColumnSequence: i,
			})
			if err != nil {
				return err
			}
		}
	}
	return nil
}
