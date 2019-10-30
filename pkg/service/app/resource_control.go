// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package app

import (
	"context"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	"openpitrix.io/openpitrix/pkg/client/isv"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/gerr"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/pi"
	"openpitrix.io/openpitrix/pkg/service/category/categoryutil"
	"openpitrix.io/openpitrix/pkg/util/ctxutil"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
	"openpitrix.io/openpitrix/pkg/util/stringutil"
)

type Action int32

const (
	Submit Action = iota
	Review
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
	Modify:  {constants.StatusDraft, constants.StatusRejected, constants.StatusActive},
	Submit:  {constants.StatusDraft, constants.StatusRejected},
	Review:  {constants.StatusSubmitted, constants.StatusInReview},
	Cancel:  {constants.StatusSubmitted, constants.StatusInReview, constants.StatusPassed},
	Release: {constants.StatusPassed},
	Delete: {
		constants.StatusSuspended,
		constants.StatusDraft,
		constants.StatusPassed,
		constants.StatusRejected,
	},
	Pass:    {constants.StatusInReview},
	Reject:  {constants.StatusInReview},
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

func checkModifyAppPermission(ctx context.Context, app *models.App) error {
	var count uint32
	var err error

	if app.Status == constants.StatusDeleted {
		return gerr.New(ctx, gerr.FailedPrecondition, gerr.ErrorResourceAlreadyDeleted, app.AppId)
	}

	count, err = pi.Global().DB(ctx).
		Select("").
		From(constants.TableAppVersion).
		Where(db.Eq(constants.ColumnAppId, app.AppId)).
		Where(db.Eq(constants.ColumnStatus, constants.StatusInReview)).
		Where(db.Eq(constants.ColumnActive, app.Active)).
		Count()
	if err != nil {
		return gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}
	if count > 0 {
		return gerr.New(ctx, gerr.FailedPrecondition, gerr.ErrorAppVersionInReview)
	}
	return nil
}

func getAppVersion(ctx context.Context, versionId string) (*models.AppVersion, error) {
	version := &models.AppVersion{}
	err := pi.Global().DB(ctx).
		Select(models.AppVersionColumns...).
		From(constants.TableAppVersion).
		Where(db.Eq(constants.ColumnVersionId, versionId)).
		Where(db.Eq(constants.ColumnActive, false)).
		LoadOne(&version)
	if err != nil {
		return nil, err
	}
	return version, nil
}

func checkAppName(ctx context.Context, name string) error {
	count, err := pi.Global().DB(ctx).
		Select(models.AppColumns...).
		From(constants.TableApp).
		Where(db.Eq(constants.ColumnName, name)).
		Where(db.Neq(constants.ColumnStatus, constants.StatusDeleted)).
		Where(db.Eq(constants.ColumnActive, false)).
		Count()
	if err != nil {
		return gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}
	if count > 0 {
		return gerr.New(ctx, gerr.PermissionDenied, gerr.ErrorAppNameExists, name)
	}

	return nil
}

func checkAppVersionName(ctx context.Context, appId, name string) error {
	count, err := pi.Global().DB(ctx).
		Select(models.AppVersionColumns...).
		From(constants.TableAppVersion).
		Where(db.Eq(constants.ColumnName, name)).
		Where(db.Eq(constants.ColumnAppId, appId)).
		Where(db.Neq(constants.ColumnStatus, constants.StatusDeleted)).
		Count()
	if err != nil {
		return gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}
	if count > 0 {
		return gerr.New(ctx, gerr.PermissionDenied, gerr.ErrorAppVersionExists, appId, name)
	}

	return nil
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
	err = addAppVersionAudit(ctx, version, constants.StatusDraft, constants.OperatorTypeDeveloper, "")
	if err != nil {
		return err
	}
	return nil
}

func updateVersion(ctx context.Context, versionId string, attributes map[string]interface{}) error {
	attributes[constants.ColumnUpdateTime] = time.Now()
	_, err := pi.Global().DB(ctx).
		Update(constants.TableAppVersion).
		SetMap(attributes).
		Where(db.Eq(constants.ColumnVersionId, versionId)).
		Where(db.Eq(constants.ColumnActive, false)).
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
		Where(db.Eq(constants.ColumnActive, false)).
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

func syncActiveVersion(ctx context.Context, version *models.AppVersion) error {
	var existActiveVersion *models.AppVersion
	_, err := pi.Global().DB(ctx).
		Select(models.AppVersionColumns...).
		From(constants.TableAppVersion).
		Where(db.Eq(constants.ColumnVersionId, version.VersionId)).
		Where(db.Eq(constants.ColumnActive, true)).
		Load(&existActiveVersion)
	if err != nil {
		return err
	}
	if existActiveVersion == nil {
		_, err := pi.Global().DB(ctx).
			InsertInto(constants.TableAppVersion).
			Record(version).
			Exec()
		return err
	}
	var updateAttr = make(map[string]interface{})
	if existActiveVersion.PackageName != version.PackageName {
		updateAttr[constants.ColumnPackageName] = version.PackageName
	}
	if existActiveVersion.Name != version.Name {
		updateAttr[constants.ColumnName] = version.Name
	}
	if existActiveVersion.Description != version.Description {
		updateAttr[constants.ColumnDescription] = version.Description
	}
	if existActiveVersion.Status != version.Status {
		updateAttr[constants.ColumnStatus] = version.Status
	}
	if len(updateAttr) == 0 {
		return nil
	}
	_, err = pi.Global().DB(ctx).
		Update(constants.TableAppVersion).SetMap(updateAttr).
		Where(db.Eq(constants.ColumnVersionId, version.VersionId)).
		Where(db.Eq(constants.ColumnActive, true)).
		Exec()
	return err
}

func syncActiveApp(ctx context.Context, app *models.App) error {
	var existActiveApp *models.App
	_, err := pi.Global().DB(ctx).
		Select(models.AppColumns...).
		From(constants.TableApp).
		Where(db.Eq(constants.ColumnAppId, app.AppId)).
		Where(db.Eq(constants.ColumnActive, true)).
		Load(&existActiveApp)
	if err != nil {
		return err
	}
	if existActiveApp == nil {
		_, err := pi.Global().DB(ctx).
			InsertInto(constants.TableApp).
			Record(app).
			Exec()
		return err
	}
	var updateAttr = make(map[string]interface{})
	if existActiveApp.RepoId != app.RepoId {
		updateAttr[constants.ColumnRepoId] = app.RepoId
	}
	if existActiveApp.Name != app.Name {
		updateAttr[constants.ColumnName] = app.Name
	}
	if existActiveApp.Description != app.Description {
		updateAttr[constants.ColumnDescription] = app.Description
	}
	if existActiveApp.Status != app.Status {
		updateAttr[constants.ColumnStatus] = app.Status
	}
	if existActiveApp.Home != app.Home {
		updateAttr[constants.ColumnHome] = app.Home
	}
	if existActiveApp.Icon != app.Icon {
		updateAttr[constants.ColumnIcon] = app.Icon
	}
	if existActiveApp.Screenshots != app.Screenshots {
		updateAttr[constants.ColumnScreenshots] = app.Screenshots
	}
	if existActiveApp.Maintainers != app.Maintainers {
		updateAttr[constants.ColumnMaintainers] = app.Maintainers
	}
	if existActiveApp.Keywords != app.Keywords {
		updateAttr[constants.ColumnKeywords] = app.Keywords
	}
	if existActiveApp.Sources != app.Sources {
		updateAttr[constants.ColumnSources] = app.Sources
	}
	if existActiveApp.Readme != app.Readme {
		updateAttr[constants.ColumnReadme] = app.Readme
	}
	if existActiveApp.ChartName != app.ChartName {
		updateAttr[constants.ColumnChartName] = app.ChartName
	}
	if len(updateAttr) == 0 {
		return nil
	}
	_, err = pi.Global().DB(ctx).
		Update(constants.TableApp).SetMap(updateAttr).
		Where(db.Eq(constants.ColumnAppId, app.AppId)).
		Where(db.Eq(constants.ColumnActive, true)).
		Exec()
	return err
}

// sync app version and app infos to active app version and active app
// it will active or deactive app_version/app
func syncAppVersion(ctx context.Context, versionId string) error {
	version, err := getAppVersion(ctx, versionId)
	if err != nil {
		return err
	}
	app, err := getApp(ctx, version.AppId)
	if err != nil {
		return err
	}
	version.Active = true
	app.Active = true

	err = syncActiveVersion(ctx, version)
	if err != nil {
		return err
	}
	return syncActiveApp(ctx, app)
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
		if app.Status != constants.StatusActive {
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
		Where(db.Eq(constants.ColumnActive, false)).
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
		Where(db.Eq(constants.ColumnActive, false)).
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
		Where(db.Eq(constants.ColumnActive, false)).
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

func getVendorMap(ctx context.Context, userIds []string) (map[string]*pb.VendorVerifyInfo, error) {
	userIds = stringutil.Unique(userIds)
	var isvs []string
	for _, uid := range userIds {
		if stringutil.StringIn(uid, constants.InternalUsers) {
			continue
		}
		isvs = append(isvs, uid)
	}

	vendorVerifyInfoSet, err := isv.GetVendorInfos(ctx, isvs)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}
	var vendorMap = make(map[string]*pb.VendorVerifyInfo)
	for _, v := range vendorVerifyInfoSet {
		vendorMap[v.Owner.GetValue()] = v
	}
	return vendorMap, nil
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

func formatAppSet(ctx context.Context, apps []*models.App, active bool) ([]*pb.App, error) {
	var appIds []string
	var userIds []string
	for _, app := range apps {
		appIds = append(appIds, app.AppId)
		userIds = append(userIds, app.Isv)
	}
	var pbApps []*pb.App
	appsVersionTypes, err := getAppsVersionTypes(ctx, appIds, active)
	if err != nil {
		return pbApps, err
	}
	rcMap, err := categoryutil.GetResourcesCategories(ctx, pi.Global().DB(ctx), appIds)
	if err != nil {
		return pbApps, err
	}
	vendorMap, err := getVendorMap(ctx, userIds)
	if err != nil {
		return pbApps, err
	}
	for _, app := range apps {

		var pbApp *pb.App
		pbApp, err := formatApp(ctx, app)
		if err != nil {
			return pbApps, err
		}
		if appVersionType, ok := appsVersionTypes[app.AppId]; ok {
			pbApp.AppVersionTypes = pbutil.ToProtoString(appVersionType)
		}
		if categorySet, ok := rcMap[app.AppId]; ok {
			pbApp.CategorySet = categorySet
		}
		if vendor, ok := vendorMap[app.Isv]; ok {
			pbApp.CompanyJoinTime = vendor.StatusTime
			pbApp.CompanyName = vendor.CompanyName
			pbApp.CompanyProfile = vendor.CompanyProfile
			pbApp.CompanyWebsite = vendor.CompanyWebsite
		}

		pbApps = append(pbApps, pbApp)
	}
	return pbApps, nil
}

func loadAppNameMap(ctx context.Context, appIds []string) (map[string]string, error) {
	var appNameMap = make(map[string]string)
	if len(appIds) > 0 {
		_, err := pi.Global().DB(ctx).
			Select(constants.ColumnAppId, constants.ColumnName).
			From(constants.TableApp).
			Where(db.Eq(constants.ColumnAppId, appIds)).
			Load(&appNameMap)
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
		}
	}
	return appNameMap, nil
}

func loadVersionMap(ctx context.Context, versionIds []string) (map[string]*models.AppVersion, error) {
	var versions []*models.AppVersion
	versionMaps := make(map[string]*models.AppVersion)
	_, err := pi.Global().DB(ctx).
		Select(models.AppVersionColumns...).
		From(constants.TableAppVersion).
		Where(db.Eq(constants.ColumnVersionId, versionIds)).
		Load(&versions)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}

	for _, version := range versions {
		versionMaps[version.VersionId] = version
	}

	return versionMaps, nil
}

func formatAppVersionReviewSet(ctx context.Context, appVersionReviews []*models.AppVersionReview) ([]*pb.AppVersionReview, error) {
	var (
		appIds     []string
		versionIds []string
	)
	for _, a := range appVersionReviews {
		appIds = append(appIds, a.AppId)
		versionIds = append(versionIds, a.VersionId)
	}
	appNameMap, err := loadAppNameMap(ctx, appIds)
	if err != nil {
		return nil, err
	}
	versionMap, err := loadVersionMap(ctx, versionIds)
	if err != nil {
		return nil, err
	}
	pbAvrs := models.AppVersionReviewsToPbs(appVersionReviews)
	for _, a := range pbAvrs {
		appId := a.GetAppId().GetValue()
		versionId := a.GetVersionId().GetValue()
		if appName, ok := appNameMap[appId]; ok {
			a.AppName = pbutil.ToProtoString(appName)
		}
		if version, ok := versionMap[versionId]; ok {
			a.VersionName = pbutil.ToProtoString(version.Name)
			a.VersionType = pbutil.ToProtoString(version.Type)
		}
	}

	return pbAvrs, nil
}

func formatAppVersionAuditSet(ctx context.Context, appVersionAudits []*models.AppVersionAudit) ([]*pb.AppVersionAudit, error) {
	var (
		appIds     []string
		versionIds []string
	)
	for _, a := range appVersionAudits {
		appIds = append(appIds, a.AppId)
		versionIds = append(versionIds, a.VersionId)
	}
	appNameMap, err := loadAppNameMap(ctx, appIds)
	if err != nil {
		return nil, err
	}
	versionMap, err := loadVersionMap(ctx, versionIds)
	if err != nil {
		return nil, err
	}
	pbAvas := models.AppVersionAuditsToPbs(appVersionAudits)
	for _, a := range pbAvas {
		appId := a.GetAppId().GetValue()
		versionId := a.GetVersionId().GetValue()
		if appName, ok := appNameMap[appId]; ok {
			a.AppName = pbutil.ToProtoString(appName)
		}
		if version, ok := versionMap[versionId]; ok {
			a.VersionName = pbutil.ToProtoString(version.Name)
			a.VersionType = pbutil.ToProtoString(version.Type)
		}
	}

	return pbAvas, nil
}

func getAppsVersionTypes(ctx context.Context, appIds []string, active bool) (map[string]string, error) {
	var appsVersionTypes = make(map[string]string)

	if len(appIds) == 0 {
		return appsVersionTypes, nil
	}

	_, err := pi.Global().DB(ctx).
		Select(constants.ColumnAppId, "GROUP_CONCAT(DISTINCT type ORDER BY type SEPARATOR ',')").
		From(constants.TableAppVersion).
		GroupBy(constants.ColumnAppId).
		Where(db.Eq(constants.ColumnAppId, appIds)).
		Where(db.Eq(constants.ColumnActive, active)).
		Load(&appsVersionTypes)

	return appsVersionTypes, err
}

func resortAppVersions(ctx context.Context, appId string) error {
	queryFunc := func(active bool) (versions models.AppVersions, err error) {
		_, err = pi.Global().DB(ctx).
			Select(
				constants.ColumnVersionId, constants.ColumnName,
				constants.ColumnSequence, constants.ColumnCreateTime,
				constants.ColumnActive,
			).
			From(constants.TableAppVersion).
			Where(db.Eq(constants.ColumnActive, active)).
			Where(db.Eq(constants.ColumnAppId, appId)).
			Where(db.Neq(constants.ColumnStatus, constants.StatusDeleted)).
			Load(&versions)
		return
	}

	sortFunc := func(versions models.AppVersions) error {
		sort.Sort(versions)
		for i, version := range versions {
			if version.Sequence != uint32(i) {
				_, err := pi.Global().DB(ctx).
					Update(constants.TableAppVersion).
					Set(constants.ColumnSequence, i).
					Set(constants.ColumnUpdateTime, time.Now()).
					Where(db.Eq(constants.ColumnVersionId, version.VersionId)).
					Where(db.Eq(constants.ColumnActive, version.Active)).
					Exec()
				if err != nil {
					return err
				}
			}
		}
		return nil
	}

	var versions models.AppVersions
	var err error
	versions, err = queryFunc(false)
	if err != nil {
		return err
	}
	err = sortFunc(versions)
	if err != nil {
		return err
	}
	versions, err = queryFunc(true)
	if err != nil {
		return err
	}
	err = sortFunc(versions)
	if err != nil {
		return err
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
		Where(db.Eq(constants.ColumnActive, false)).
		Where(db.Neq(constants.ColumnVersionId, ignoredVersionIds)).
		Exec()
	return err
}

func clearRepoAppVersions(ctx context.Context, repoId string, ignoredVersionIds []string) error {
	_, err := pi.Global().DB(ctx).
		Update(constants.TableAppVersion).
		Set(constants.ColumnStatus, constants.StatusDeleted).
		Set(constants.ColumnStatusTime, time.Now()).
		Set(constants.ColumnUpdateTime, time.Now()).
		Where(db.Eq(constants.ColumnAppId,
			[]*db.SelectQuery{pi.Global().DB(ctx).
				Select(constants.ColumnAppId).
				From(constants.TableApp).
				Where(db.Eq(constants.ColumnRepoId, repoId))})).
		Where(db.Eq(constants.ColumnActive, false)).
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
		Where(db.Eq(constants.ColumnActive, false)).
		Where(db.Neq(constants.ColumnAppId, ignoredAppIds)).
		Exec()
	return err
}

var versionStatusPriority = map[string]int32{
	constants.StatusActive:    7,
	constants.StatusRejected:  6,
	constants.StatusPassed:    5,
	constants.StatusSubmitted: 4,
	constants.StatusSuspended: 3,
	constants.StatusDraft:     2,
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
	c, ok := versionStatusPriority[currentStatus]
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
		Where(db.Eq(constants.ColumnAppId, appId)).
		Where(db.Eq(constants.ColumnActive, false)).
		GroupBy("status").
		Load(&statusCountMap)
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
	{constants.StatusDeleted, constants.StatusDraft},
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

func addAppVersionAudit(ctx context.Context, version *models.AppVersion, status, operatorType, message string) error {
	s := ctxutil.GetSender(ctx)
	versionAudit := models.NewAppVersionAudit(version.VersionId, version.AppId, status, s.UserId, operatorType, s.GetOwnerPath())
	versionAudit.Message = message
	if version.Status == constants.StatusSubmitted {
		versionAudit.ReviewId = version.ReviewId
	}

	_, err := pi.Global().DB(ctx).
		InsertInto(constants.TableAppVersionAudit).
		Record(versionAudit).
		Exec()
	if err != nil {
		logger.Error(ctx, "Failed to insert version audit [%+v]", versionAudit)
		return gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
	}
	return nil
}

func matchPackageFailedError(err error, res *pb.ValidatePackageResponse) {
	var errStr = err.Error()
	var matchedError = ""
	var errorDetails = make(map[string]string)
	switch {
	// Helm errors
	case strings.HasPrefix(errStr, "no files in chart archive"),
		strings.HasPrefix(errStr, "no files in app archive"):

		matchedError = "no files in package"

	case strings.HasPrefix(errStr, "chart yaml not in base directory"),
		strings.HasPrefix(errStr, "chart metadata (Chart.yaml) missing"):

		errorDetails["Chart.yaml"] = "not found"

	case strings.HasPrefix(errStr, "invalid chart (Chart.yaml): name must not be empty"):

		errorDetails["Chart.yaml"] = "package name must not be empty"

	case strings.HasPrefix(errStr, "values.toml is illegal"):

		errorDetails["values.toml"] = errStr

	case strings.HasPrefix(errStr, "error reading"):

		matched := regexp.MustCompile("error reading (.+): (.+)").FindStringSubmatch(errStr)
		if len(matched) > 0 {
			errorDetails[matched[1]] = matched[2]
		}

		// Devkit errors
	case strings.HasPrefix(errStr, "[package.json] not in base directory"):

		errorDetails["package.json"] = "not found"

	case strings.HasPrefix(errStr, "missing file ["):

		matched := regexp.MustCompile("missing file \\[(.+)]").FindStringSubmatch(errStr)
		if len(matched) > 0 {
			errorDetails[matched[1]] = "not found"
		}

	case strings.HasPrefix(errStr, "failed to parse"),
		strings.HasPrefix(errStr, "failed to render"),
		strings.HasPrefix(errStr, "failed to load"),
		strings.HasPrefix(errStr, "failed to decode"):

		matched := regexp.MustCompile("failed to (.+) (.+): (.+)").FindStringSubmatch(errStr)
		if len(matched) > 0 {
			errorDetails[matched[2]] = fmt.Sprintf("%s failed, %s", matched[1], matched[3])
		}

	default:
		matchedError = errStr
	}
	if len(errorDetails) > 0 {
		res.ErrorDetails = errorDetails
	}
	if len(matchedError) > 0 {
		res.Error = pbutil.ToProtoString(matchedError)
	}
}
