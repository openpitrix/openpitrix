// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package app

import (
	"bytes"
	"context"
	"strings"
	"time"

	repoClient "openpitrix.io/openpitrix/pkg/client/repo"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/gerr"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/pi"
	"openpitrix.io/openpitrix/pkg/service/category/categoryutil"
	"openpitrix.io/openpitrix/pkg/util/gziputil"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
	"openpitrix.io/openpitrix/pkg/util/senderutil"
)

func (p *Server) getAppVersion(ctx context.Context, versionId string) (*models.AppVersion, error) {
	version := &models.AppVersion{}
	err := pi.Global().DB(ctx).
		Select(models.AppVersionColumns...).
		From(constants.TableAppVersion).
		Where(db.Eq(constants.ColumnVersionId, versionId)).
		LoadOne(&version)
	if err != nil {
		return nil, err
	}
	return version, nil
}

func (p *Server) getAppVersions(ctx context.Context, versionIds []string) ([]*models.AppVersion, error) {
	var versions []*models.AppVersion
	_, err := pi.Global().DB(ctx).
		Select(models.AppVersionColumns...).
		From(constants.TableAppVersion).
		Where(db.Eq(constants.ColumnVersionId, versionIds)).
		Load(&versions)
	if err != nil {
		return nil, err
	}
	return versions, nil
}

func (p *Server) SyncRepo(ctx context.Context, req *pb.SyncRepoRequest) (*pb.SyncRepoResponse, error) {
	var res = &pb.SyncRepoResponse{}
	var repoId = req.GetRepoId()
	if len(repoId) == 0 {
		return res, nil
	}
	var failed = func(reason string) (*pb.SyncRepoResponse, error) {
		res.Failed = true
		res.Result = reason
		return res, nil
	}

	repoManagerClient, err := repoClient.NewRepoManagerClient()
	if err != nil {
		return failed("internal error")
	}
	describeRepoReq := pb.DescribeReposRequest{
		RepoId: []string{repoId},
	}
	describeRepoRes, err := repoManagerClient.DescribeRepos(ctx, &describeRepoReq)
	if err != nil {
		logger.Error(ctx, "Failed to describe repo [%s], %+v", repoId, err)
		return failed("internal error")
	}
	if describeRepoRes.TotalCount == 0 {
		logger.Error(ctx, "Failed to describe repo [%s], repo not exists", repoId)
		return failed("internal error")
	}
	repo := describeRepoRes.RepoSet[0]
	err = newRepoProxy(repo).SyncRepo(ctx)
	if err != nil {
		return failed(err.Error())
	}
	return res, nil
}
func (p *Server) DescribeApps(ctx context.Context, req *pb.DescribeAppsRequest) (*pb.DescribeAppsResponse, error) {
	var apps []*models.App
	offset := pbutil.GetOffsetFromRequest(req)
	limit := pbutil.GetLimitFromRequest(req)
	categoryIds := req.GetCategoryId()

	query := pi.Global().DB(ctx).
		Select(models.AppColumns...).
		From(constants.TableApp).
		Offset(offset).
		Limit(limit).
		Where(manager.BuildFilterConditions(req, constants.TableApp))
	if len(categoryIds) > 0 {
		subqueryStmt := pi.Global().DB(ctx).
			Select(constants.ColumnResouceId).
			From(constants.TableCategoryResource).
			Where(db.Eq(constants.ColumnStatus, constants.StatusEnabled)).
			Where(db.Eq(constants.ColumnCategoryId, categoryIds))
		query = query.Where(db.Eq(constants.ColumnAppId, []*db.SelectQuery{subqueryStmt}))
	}
	// TODO: validate sort_key
	query = manager.AddQueryOrderDir(query, req, constants.ColumnCreateTime)
	_, err := query.Load(&apps)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}
	count, err := query.Count()
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}

	appSet, err := formatAppSet(ctx, apps)
	if err != nil {
		return nil, err
	}

	res := &pb.DescribeAppsResponse{
		AppSet:     appSet,
		TotalCount: count,
	}
	return res, nil
}

func (p *Server) CreateApp(ctx context.Context, req *pb.CreateAppRequest) (*pb.CreateAppResponse, error) {
	// TODO: validate CreateAppRequest
	// TODO: check categories

	s := senderutil.GetSenderFromContext(ctx)
	newApp := models.NewApp(
		req.GetName().GetValue(),
		req.GetRepoId().GetValue(),
		req.GetDescription().GetValue(),
		s.UserId,
		req.GetChartName().GetValue())

	newApp.Home = req.GetHome().GetValue()
	newApp.Icon = req.GetIcon().GetValue()
	newApp.Screenshots = req.GetScreenshots().GetValue()
	newApp.Sources = req.GetSources().GetValue()
	newApp.Readme = req.GetReadme().GetValue()
	newApp.Keywords = req.GetKeywords().GetValue()

	if req.GetStatus() != nil {
		newApp.Status = req.GetStatus().GetValue()
	} else {
		newApp.Status = pi.Global().GlobalConfig().GetAppDefaultStatus()
	}

	_, err := pi.Global().DB(ctx).
		InsertInto(constants.TableApp).
		Record(newApp).
		Exec()
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
	}

	err = categoryutil.SyncResourceCategories(
		ctx,
		pi.Global().DB(ctx),
		newApp.AppId,
		categoryutil.DecodeCategoryIds(req.GetCategoryId().GetValue()),
	)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
	}

	res := &pb.CreateAppResponse{
		AppId: pbutil.ToProtoString(newApp.AppId),
	}

	if req.GetPackage() != nil {
		version := models.NewAppVersion(newApp.AppId, newApp.Name, newApp.Description, newApp.Owner, "")
		err = insertVersion(ctx, version)
		if err != nil {
			return nil, err
		}
		err = newVersionProxy(version).AddPackageFile(
			ctx, req.GetPackage().GetValue(), true, false)
		if err != nil {
			deleteApp(ctx, newApp.AppId)
			deleteVersion(ctx, version.VersionId)
			return nil, err
		}

		res.VersionId = pbutil.ToProtoString(version.VersionId)
	}

	return res, nil
}

func (p *Server) ModifyApp(ctx context.Context, req *pb.ModifyAppRequest) (*pb.ModifyAppResponse, error) {
	err := models.CheckAppPermission(ctx, req.GetAppId().GetValue())
	if err != nil {
		return nil, err
	}
	appId := req.GetAppId().GetValue()
	app, err := getApp(ctx, appId)
	if err != nil {
		return nil, err
	}
	if app.Status == constants.StatusDeleted {
		return nil, gerr.NewWithDetail(ctx, gerr.FailedPrecondition, err, gerr.ErrorResourceAlreadyDeleted, appId)
	}

	attributes := manager.BuildUpdateAttributes(req,
		"name", "repo_id", "owner", "chart_name",
		"description", "home", "icon", "screenshots",
		"maintainers", "sources", "readme", "keywords")

	if req.GetStatus() != nil {
		attributes[constants.ColumnStatus] = req.GetStatus().GetValue()
	}
	err = updateApp(ctx, appId, attributes)
	if err != nil {
		return nil, err
	}

	if req.GetCategoryId() != nil {
		err = categoryutil.SyncResourceCategories(
			ctx,
			pi.Global().DB(ctx),
			appId,
			categoryutil.DecodeCategoryIds(req.GetCategoryId().GetValue()),
		)
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorModifyResourcesFailed)
		}
	}

	res := &pb.ModifyAppResponse{
		AppId: req.GetAppId(),
	}
	return res, nil
}

func (p *Server) DeleteApps(ctx context.Context, req *pb.DeleteAppsRequest) (*pb.DeleteAppsResponse, error) {
	err := models.CheckAppPermission(ctx, req.GetAppId()...)
	if err != nil {
		return nil, err
	}
	appIds := req.GetAppId()

	_, err = pi.Global().DB(ctx).
		Update(constants.TableApp).
		Set(constants.ColumnStatus, constants.StatusDeleted).
		Where(db.Eq(constants.ColumnAppId, appIds)).
		Exec()
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDeleteResourcesFailed)
	}

	return &pb.DeleteAppsResponse{
		AppId: appIds,
	}, nil
}

func (p *Server) CreateAppVersion(ctx context.Context, req *pb.CreateAppVersionRequest) (*pb.CreateAppVersionResponse, error) {
	// TODO: validate CreateAppVersionRequest
	s := senderutil.GetSenderFromContext(ctx)
	newAppVersion := models.NewAppVersion(
		req.GetAppId().GetValue(),
		req.GetName().GetValue(),
		req.GetDescription().GetValue(),
		s.UserId,
		req.GetPackageName().GetValue())

	if req.Sequence != nil {
		newAppVersion.Sequence = req.Sequence.GetValue()
	} else {
		seq, err := getBigestSequence(ctx, newAppVersion.AppId)
		if err != nil {
			return nil, err
		}
		newAppVersion.Sequence = seq + 1
	}

	newAppVersion.Home = req.GetHome().GetValue()
	newAppVersion.Icon = req.GetIcon().GetValue()
	newAppVersion.Screenshots = req.GetScreenshots().GetValue()
	newAppVersion.Maintainers = req.GetMaintainers().GetValue()
	newAppVersion.Keywords = req.GetKeywords().GetValue()
	newAppVersion.Sources = req.GetSources().GetValue()
	newAppVersion.Readme = req.GetReadme().GetValue()

	if req.GetStatus() != nil {
		newAppVersion.Status = req.GetStatus().GetValue()
	} else {
		newAppVersion.Status = pi.Global().GlobalConfig().GetAppDefaultStatus()
	}

	err := insertVersion(ctx, newAppVersion)
	if err != nil {
		return nil, err
	}

	if req.GetPackage() != nil {
		var syncAppName = false
		if newAppVersion.Sequence == 1 {
			syncAppName = true
		}
		err = newVersionProxy(newAppVersion).AddPackageFile(
			ctx, req.GetPackage().GetValue(), syncAppName, false)
		if err != nil {
			deleteVersion(ctx, newAppVersion.VersionId)
			return nil, err
		}
	}

	res := &pb.CreateAppVersionResponse{
		VersionId: pbutil.ToProtoString(newAppVersion.VersionId),
	}
	return res, nil

}

func (p *Server) DescribeAppVersions(ctx context.Context, req *pb.DescribeAppVersionsRequest) (*pb.DescribeAppVersionsResponse, error) {
	var versions []*models.AppVersion
	offset := pbutil.GetOffsetFromRequest(req)
	limit := pbutil.GetLimitFromRequest(req)

	query := pi.Global().DB(ctx).
		Select(models.AppVersionColumns...).
		From(constants.TableAppVersion).
		Offset(offset).
		Limit(limit).
		Where(manager.BuildFilterConditions(req, constants.TableAppVersion))
	query = manager.AddQueryOrderDir(query, req, constants.ColumnSequence)
	_, err := query.Load(&versions)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}
	count, err := query.Count()
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}

	res := &pb.DescribeAppVersionsResponse{
		AppVersionSet: models.AppVersionsToPbs(versions),
		TotalCount:    count,
	}
	return res, nil

}

func (p *Server) ModifyAppVersion(ctx context.Context, req *pb.ModifyAppVersionRequest) (*pb.ModifyAppVersionResponse, error) {
	err := models.CheckAppVersionPermission(ctx, req.GetVersionId().GetValue())
	if err != nil {
		return nil, err
	}
	versionId := req.GetVersionId().GetValue()
	version, err := checkAppVersionHandlePermission(ctx, Modify, versionId)
	if err != nil {
		return nil, err
	}

	attributes := manager.BuildUpdateAttributes(req,
		"name", "description", "package_name",
		"sequence", "home", "icon", "screenshots",
		"maintainers", "keywords", "sources", "readme")

	if version.Status != constants.StatusActive {
		attributes[constants.ColumnStatus] = constants.StatusDraft
	}
	if req.GetStatus() != nil {
		attributes[constants.ColumnStatus] = req.GetStatus().GetValue()
	}

	err = updateVersion(ctx, versionId, attributes)
	if err != nil {
		return nil, err
	}

	if req.GetPackage() != nil {
		err = newVersionProxy(version).AddPackageFile(
			ctx, req.GetPackage().GetValue(), false, true)
		if err != nil {
			return nil, err
		}
	}

	res := &pb.ModifyAppVersionResponse{
		VersionId: req.GetVersionId(),
	}
	return res, nil

}

func (p *Server) DeleteAppVersions(ctx context.Context, req *pb.DeleteAppVersionsRequest) (*pb.DeleteAppVersionsResponse, error) {
	err := models.CheckAppVersionPermission(ctx, req.GetVersionId()...)
	if err != nil {
		return nil, err
	}
	versionIds := req.GetVersionId()
	for _, versionId := range versionIds {
		_, err := checkAppVersionHandlePermission(ctx, Delete, versionId)
		if err != nil {
			return nil, err
		}
	}

	_, err = pi.Global().DB(ctx).
		Update(constants.TableAppVersion).
		Set(constants.ColumnStatus, constants.StatusDeleted).
		Where(db.Eq(constants.ColumnVersionId, versionIds)).
		Exec()
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDeleteResourceFailed, strings.Join(versionIds, ","))
	}

	return &pb.DeleteAppVersionsResponse{
		VersionId: versionIds,
	}, nil
}

func (p *Server) GetAppVersionPackage(ctx context.Context, req *pb.GetAppVersionPackageRequest) (*pb.GetAppVersionPackageResponse, error) {
	versionId := req.GetVersionId().GetValue()
	version, err := p.getAppVersion(ctx, versionId)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.NotFound, err, gerr.ErrorResourceNotFound, versionId)
	}
	logger.Debug(ctx, "Got app version: [%+v]", version)

	content, err := newVersionProxy(version).GetPackageFile(ctx)
	if err != nil {
		return nil, err
	}

	return &pb.GetAppVersionPackageResponse{
		Package:   content,
		AppId:     pbutil.ToProtoString(version.AppId),
		VersionId: req.GetVersionId(),
	}, nil
}

func (p *Server) GetAppVersionPackageFiles(ctx context.Context, req *pb.GetAppVersionPackageFilesRequest) (*pb.GetAppVersionPackageFilesResponse, error) {
	versionId := req.GetVersionId().GetValue()
	includeFiles := req.Files
	version, err := p.getAppVersion(ctx, versionId)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.NotFound, err, gerr.ErrorResourceNotFound, versionId)
	}

	content, err := newVersionProxy(version).GetPackageFile(ctx)
	if err != nil {
		logger.Error(ctx, "Failed to load [%s] package, error: %+v", versionId, err)
		return nil, err
	}

	archiveFiles, err := gziputil.LoadArchive(bytes.NewReader(content), includeFiles...)
	if err != nil {
		logger.Error(ctx, "Failed to load [%s] package, error: %+v", versionId, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourceFailed, versionId)
	}
	return &pb.GetAppVersionPackageFilesResponse{
		Files:     archiveFiles,
		VersionId: req.GetVersionId(),
	}, nil
}

type appStatistic struct {
	Date  string `db:"DATE_FORMAT(create_time, '%Y-%m-%d')"`
	Count uint32 `db:"COUNT(app_id)"`
}
type repoStatistic struct {
	RepoId string `db:"repo_id"`
	Count  uint32 `db:"COUNT(app_id)"`
}

func (p *Server) GetAppStatistics(ctx context.Context, req *pb.GetAppStatisticsRequest) (*pb.GetAppStatisticsResponse, error) {
	res := &pb.GetAppStatisticsResponse{
		LastTwoWeekCreated: make(map[string]uint32),
		TopTenRepos:        make(map[string]uint32),
	}
	appCount, err := pi.Global().DB(ctx).
		Select(constants.ColumnAppId).
		From(constants.TableApp).
		Where(db.Neq(constants.ColumnStatus, constants.StatusDeleted)).
		Count()
	if err != nil {
		logger.Error(ctx, "Failed to get app count, error: %+v", err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}
	res.AppCount = appCount

	err = pi.Global().DB(ctx).
		Select("COUNT(DISTINCT repo_id)").
		From(constants.TableApp).
		Where(db.Neq(constants.ColumnStatus, constants.StatusDeleted)).
		LoadOne(&res.RepoCount)
	if err != nil {
		logger.Error(ctx, "Failed to get repo count, error: %+v", err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}

	time2week := time.Now().Add(-14 * 24 * time.Hour)
	var as []*appStatistic
	_, err = pi.Global().DB(ctx).
		Select("DATE_FORMAT(create_time, '%Y-%m-%d')", "COUNT(app_id)").
		From(constants.TableApp).
		GroupBy("DATE_FORMAT(create_time, '%Y-%m-%d')").
		Where(db.Gte(constants.ColumnCreateTime, time2week)).
		Limit(14).Load(&as)

	if err != nil {
		logger.Error(ctx, "Failed to get app statistics, error: %+v", err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}
	for _, a := range as {
		res.LastTwoWeekCreated[a.Date] = a.Count
	}

	var rs []*repoStatistic
	_, err = pi.Global().DB(ctx).
		Select("repo_id", "COUNT(app_id)").
		From(constants.TableApp).
		Where(db.Neq(constants.ColumnStatus, constants.StatusDeleted)).
		GroupBy(constants.ColumnRepoId).
		OrderDir("COUNT(app_id)", false).
		Limit(10).Load(&rs)

	if err != nil {
		logger.Error(ctx, "Failed to get repo statistics, error: %+v", err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}
	for _, a := range rs {
		res.TopTenRepos[a.RepoId] = a.Count
	}

	return res, nil
}

func (p *Server) SubmitAppVersion(ctx context.Context, req *pb.SubmitAppVersionRequest) (*pb.SubmitAppVersionResponse, error) {
	err := models.CheckAppVersionPermission(ctx, req.GetVersionId().GetValue())
	if err != nil {
		return nil, err
	}
	version, err := checkAppVersionHandlePermission(ctx, Submit, req.GetVersionId().GetValue())
	if err != nil {
		return nil, err
	}
	err = updateVersionStatus(ctx, version, constants.StatusSubmitted)
	if err != nil {
		return nil, err
	}
	res := pb.SubmitAppVersionResponse{
		VersionId: pbutil.ToProtoString(version.VersionId),
	}
	return &res, nil
}

func (p *Server) CancelAppVersion(ctx context.Context, req *pb.CancelAppVersionRequest) (*pb.CancelAppVersionResponse, error) {
	err := models.CheckAppVersionPermission(ctx, req.GetVersionId().GetValue())
	if err != nil {
		return nil, err
	}
	version, err := checkAppVersionHandlePermission(ctx, Cancel, req.GetVersionId().GetValue())
	if err != nil {
		return nil, err
	}
	err = updateVersionStatus(ctx, version, constants.StatusDraft)
	if err != nil {
		return nil, err
	}
	res := pb.CancelAppVersionResponse{
		VersionId: pbutil.ToProtoString(version.VersionId),
	}
	return &res, nil
}

func (p *Server) ReleaseAppVersion(ctx context.Context, req *pb.ReleaseAppVersionRequest) (*pb.ReleaseAppVersionResponse, error) {
	err := models.CheckAppVersionPermission(ctx, req.GetVersionId().GetValue())
	if err != nil {
		return nil, err
	}
	version, err := checkAppVersionHandlePermission(ctx, Release, req.GetVersionId().GetValue())
	if err != nil {
		return nil, err
	}
	err = updateVersionStatus(ctx, version, constants.StatusActive)
	if err != nil {
		return nil, err
	}
	res := pb.ReleaseAppVersionResponse{
		VersionId: pbutil.ToProtoString(version.VersionId),
	}
	return &res, nil
}

func (p *Server) DeleteAppVersion(ctx context.Context, req *pb.DeleteAppVersionRequest) (*pb.DeleteAppVersionResponse, error) {
	err := models.CheckAppVersionPermission(ctx, req.GetVersionId().GetValue())
	if err != nil {
		return nil, err
	}
	version, err := checkAppVersionHandlePermission(ctx, Delete, req.GetVersionId().GetValue())
	if err != nil {
		return nil, err
	}
	err = newVersionProxy(version).DeletePackageFile(ctx)
	if err != nil {
		logger.Error(ctx, "Failed to delete [%s] package, error: %+v", req.GetVersionId().GetValue(), err)
		return nil, err
	}
	err = updateVersionStatus(ctx, version, constants.StatusDeleted)
	if err != nil {
		return nil, err
	}
	res := pb.DeleteAppVersionResponse{
		VersionId: pbutil.ToProtoString(version.VersionId),
	}
	return &res, nil
}

func (p *Server) PassAppVersion(ctx context.Context, req *pb.PassAppVersionRequest) (*pb.PassAppVersionResponse, error) {
	err := models.CheckAppVersionPermission(ctx, req.GetVersionId().GetValue())
	if err != nil {
		return nil, err
	}
	version, err := checkAppVersionHandlePermission(ctx, Pass, req.GetVersionId().GetValue())
	if err != nil {
		return nil, err
	}
	err = updateVersionStatus(ctx, version, constants.StatusPassed)
	if err != nil {
		return nil, err
	}
	res := pb.PassAppVersionResponse{
		VersionId: pbutil.ToProtoString(version.VersionId),
	}
	return &res, nil
}

func (p *Server) RejectAppVersion(ctx context.Context, req *pb.RejectAppVersionRequest) (*pb.RejectAppVersionResponse, error) {
	err := models.CheckAppVersionPermission(ctx, req.GetVersionId().GetValue())
	if err != nil {
		return nil, err
	}
	version, err := checkAppVersionHandlePermission(ctx, Reject, req.GetVersionId().GetValue())
	if err != nil {
		return nil, err
	}
	err = updateVersionStatus(ctx, version, constants.StatusRejected, map[string]interface{}{
		constants.ColumnMessage: req.GetMessage().GetValue(),
	})
	if err != nil {
		return nil, err
	}
	res := pb.RejectAppVersionResponse{
		VersionId: pbutil.ToProtoString(version.VersionId),
	}
	return &res, nil
}

func (p *Server) SuspendAppVersion(ctx context.Context, req *pb.SuspendAppVersionRequest) (*pb.SuspendAppVersionResponse, error) {
	err := models.CheckAppVersionPermission(ctx, req.GetVersionId().GetValue())
	if err != nil {
		return nil, err
	}
	version, err := checkAppVersionHandlePermission(ctx, Suspend, req.GetVersionId().GetValue())
	if err != nil {
		return nil, err
	}
	err = updateVersionStatus(ctx, version, constants.StatusSuspended)
	if err != nil {
		return nil, err
	}
	res := pb.SuspendAppVersionResponse{
		VersionId: pbutil.ToProtoString(version.VersionId),
	}
	return &res, nil
}

func (p *Server) RecoverAppVersion(ctx context.Context, req *pb.RecoverAppVersionRequest) (*pb.RecoverAppVersionResponse, error) {
	err := models.CheckAppVersionPermission(ctx, req.GetVersionId().GetValue())
	if err != nil {
		return nil, err
	}
	version, err := checkAppVersionHandlePermission(ctx, Recover, req.GetVersionId().GetValue())
	if err != nil {
		return nil, err
	}
	err = updateVersionStatus(ctx, version, constants.StatusActive)
	if err != nil {
		return nil, err
	}
	res := pb.RecoverAppVersionResponse{
		VersionId: pbutil.ToProtoString(version.VersionId),
	}
	return &res, nil
}
