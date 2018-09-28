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
	Cancel:  {constants.StatusSubmitted, constants.StatusPassed},
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

func checkAppVersionHandlePermission(ctx context.Context, action Action, version *models.AppVersion) error {
	allowedStatus, ok := VersionFiniteStatusMachine[action]
	if !ok {
		return gerr.New(ctx, gerr.Internal, gerr.ErrorInternalError)
	}
	versionStatus := version.Status
	if !stringutil.StringIn(versionStatus, allowedStatus) {
		return gerr.New(ctx, gerr.FailedPrecondition,
			gerr.ErrorAppVersionIncorrectStatus, version.VersionId, versionStatus)
	}
	return nil
}

func deleteApp(ctx context.Context, appId string) error {
	_, err := pi.Global().DB(ctx).
		DeleteFrom(constants.TableApp).
		Where(db.Eq(constants.ColumnAppId, appId)).
		Exec()
	if err != nil {
		return gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDeleteResourcesFailed)
	}
	return err
}

func deleteVersion(ctx context.Context, versionId string) error {
	_, err := pi.Global().DB(ctx).
		DeleteFrom(constants.TableAppVersion).
		Where(db.Eq(constants.ColumnVersionId, versionId)).
		Exec()
	if err != nil {
		return gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDeleteResourcesFailed)
	}
	return err
}

func insertVersion(ctx context.Context, version *models.AppVersion) error {
	_, err := pi.Global().DB(ctx).
		InsertInto(constants.TableAppVersion).
		Record(version).
		Exec()
	if err != nil {
		logger.Error(ctx, "Failed to insert version [%+v]", version)
		return gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
	}
	return nil
}

func updateVersion(ctx context.Context, versionId string, attributes map[string]interface{}) error {
	attributes[constants.ColumnUpdateTime] = time.Now()
	_, err := pi.Global().DB(ctx).
		Update(constants.TableAppVersion).
		SetMap(attributes).
		Where(db.Eq(constants.ColumnVersionId, versionId)).
		Exec()
	if err != nil {
		return gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorModifyResourceFailed, versionId)
	}
	return nil
}

func updateApp(ctx context.Context, appId string, attributes map[string]interface{}) error {
	attributes[constants.ColumnUpdateTime] = time.Now()
	_, err := pi.Global().DB(ctx).
		Update(constants.TableApp).
		SetMap(attributes).
		Where(db.Eq(constants.ColumnAppId, appId)).
		Exec()
	if err != nil {
		return gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorModifyResourcesFailed)
	}
	return nil
}

func updateVersionStatus(ctx context.Context, version *models.AppVersion, status string, attributes ...map[string]interface{}) error {
	var attr = make(map[string]interface{})
	for _, a := range attributes {
		for k, v := range a {
			attr[k] = v
		}
	}
	attr[constants.ColumnStatus] = status
	attr[constants.ColumnStatusTime] = time.Now()
	err := updateVersion(ctx, version.VersionId, attr)
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
	attributes := make(map[string]interface{})
	activeVersion, err := getLatestAppVersion(ctx, appId, constants.StatusActive)
	if err != nil {
		return err
	}
	if activeVersion != nil {
		if activeVersion.Description != app.Description {
			attributes[constants.ColumnDescription] = activeVersion.Description
		}
		if activeVersion.Home != app.Home {
			attributes[constants.ColumnHome] = activeVersion.Home
		}
		if activeVersion.Icon != app.Icon {
			attributes[constants.ColumnIcon] = activeVersion.Icon
		}
		if activeVersion.Screenshots != app.Screenshots {
			attributes[constants.ColumnScreenshots] = activeVersion.Screenshots
		}
		if activeVersion.Maintainers != app.Maintainers {
			attributes[constants.ColumnMaintainers] = activeVersion.Maintainers
		}
		if activeVersion.Keywords != app.Keywords {
			attributes[constants.ColumnKeywords] = activeVersion.Keywords
		}
		if activeVersion.Sources != app.Sources {
			attributes[constants.ColumnSources] = activeVersion.Sources
		}
		if activeVersion.Readme != app.Readme {
			attributes[constants.ColumnReadme] = activeVersion.Readme
		}
		if constants.StatusActive != app.Status {
			attributes[constants.ColumnStatus] = constants.StatusActive
			attributes[constants.ColumnStatusTime] = time.Now()
		}
	} else {
		statusCountMap, err := groupAppVersionStatus(ctx, appId)
		if err != nil {
			return err
		}
		status := computeAppStatus(statusCountMap)
		if status != app.Status {
			attributes[constants.ColumnStatus] = status
			attributes[constants.ColumnStatusTime] = time.Now()
		}
	}

	if len(attributes) == 0 {
		return nil
	}
	attributes[constants.ColumnUpdateTime] = time.Now()

	_, err = pi.Global().DB(ctx).
		Update(constants.TableApp).
		SetMap(attributes).
		Where(db.Eq(constants.ColumnAppId, appId)).
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
		From(constants.TableApp).
		Where(db.Eq(constants.ColumnAppId, appId)).
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
		From(constants.TableAppVersion).
		Where(db.Eq(constants.ColumnAppId, appId))
	if len(status) > 0 {
		stmt.Where(db.Eq(constants.ColumnStatus, status))
	}
	err := stmt.
		OrderDir(constants.ColumnSequence, false).
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

func getAppCategories(ctx context.Context, appId string) ([]*pb.ResourceCategory, error) {
	rcmap, err := categoryutil.GetResourcesCategories(ctx, pi.Global().DB(ctx), []string{appId})
	if err != nil {
		return nil, err
	}
	return rcmap[appId], nil
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
		From(constants.TableAppVersion).
		Where(db.Eq(constants.ColumnAppId, appId)).
		Where(db.Neq(constants.ColumnStatus, constants.StatusDeleted)).
		LoadOne(&sequence)
	if err != nil {
		return sequence, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourceFailed, appId)
	}
	return sequence, nil
}

func resortAppVersions(ctx context.Context, appId string) error {
	var versions models.AppVersions
	_, err := pi.Global().DB(ctx).
		Select(constants.ColumnVersionId, constants.ColumnName, constants.ColumnSequence, constants.ColumnCreateTime).
		From(constants.TableAppVersion).
		Where(db.Eq(constants.ColumnAppId, appId)).
		Where(db.Neq(constants.ColumnStatus, constants.StatusDeleted)).
		Load(&versions)
	if err != nil {
		return gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourceFailed, appId)
	}
	sort.Sort(versions)
	for i, version := range versions {
		if version.Sequence != uint32(i) {
			err = updateVersion(ctx, version.VersionId, map[string]interface{}{
				constants.ColumnSequence: i,
			})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func clearAppVersions(ctx context.Context, appId string, ignoredVersionIds []string) error {
	_, err := pi.Global().DB(ctx).
		Update(constants.TableAppVersion).
		Set(constants.ColumnStatus, constants.StatusDeleted).
		Set(constants.ColumnStatusTime, time.Now()).
		Set(constants.ColumnUpdateTime, time.Now()).
		Where(db.Eq(constants.ColumnAppId, appId)).
		Where(db.Neq(constants.ColumnVersionId, ignoredVersionIds)).
		Exec()
	return err
}

func clearApps(ctx context.Context, repoId string, ignoredAppIds []string) error {
	_, err := pi.Global().DB(ctx).
		Update(constants.TableApp).
		Set(constants.ColumnStatus, constants.StatusDeleted).
		Set(constants.ColumnStatusTime, time.Now()).
		Set(constants.ColumnUpdateTime, time.Now()).
		Where(db.Eq(constants.ColumnRepoId, repoId)).
		Where(db.Neq(constants.ColumnAppId, ignoredAppIds)).
		Exec()
	return err
}

var versionStatusPriority = map[string]int32{
	constants.StatusActive:    7,
	constants.StatusRejected:  6,
	constants.StatusPassed:    5,
	constants.StatusSubmitted: 4,
	constants.StatusDraft:     3,
	constants.StatusSuspended: 2,
	constants.StatusDeleted:   1,
}

func getAppDefaultStatus(repo *pb.Repo) string {
	var defaultStatus = repo.GetAppDefaultStatus().GetValue()
	if len(defaultStatus) == 0 {
		return pi.Global().GlobalConfig().GetAppDefaultStatus()
	}
	return defaultStatus
}

func getAppVersionStatus(defaultStatus, currentStatus string) string {
	d, ok := versionStatusPriority[defaultStatus]
	if !ok {
		return currentStatus
	}
	c, ok := versionStatusPriority[defaultStatus]
	if !ok {
		return defaultStatus
	}
	if c >= d {
		return currentStatus
	}
	return defaultStatus
}

func groupAppVersionStatus(ctx context.Context, appId string) (map[string]int32, error) {
	var statusCountMap = make(map[string]int32)
	_, err := pi.Global().DB(ctx).
		Select("status", "count(version_id)").
		From(constants.TableAppVersion).
		Where(db.Eq(constants.ColumnAppId, appId)).GroupBy("status").Load(&statusCountMap)
	return statusCountMap, err
}

var versionStatusToAppStatus = [][]string{
	// from => to
	{constants.StatusActive, constants.StatusActive},
	{constants.StatusRejected, constants.StatusDraft},
	{constants.StatusPassed, constants.StatusDraft},
	{constants.StatusSubmitted, constants.StatusDraft},
	{constants.StatusDraft, constants.StatusDraft},
	{constants.StatusSuspended, constants.StatusSuspended},
	{constants.StatusDeleted, constants.StatusDeleted},
}

// compute status from exist app version status
func computeAppStatus(statusCountMap map[string]int32) string {
	for _, vs := range versionStatusToAppStatus {
		if c, ok := statusCountMap[vs[0]]; ok && c > 0 {
			return vs[1]
		}
	}
	if len(statusCountMap) == 0 {
		return constants.StatusDeleted
	}
	return constants.StatusDraft
}
