// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package app

import (
	"context"
	"io/ioutil"
	"strings"
	"time"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/gerr"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/service/category/categoryutil"
	"openpitrix.io/openpitrix/pkg/util/gziputil"
	"openpitrix.io/openpitrix/pkg/util/httputil"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
	"openpitrix.io/openpitrix/pkg/util/senderutil"
)

func (p *Server) getAppVersion(versionId string) (*models.AppVersion, error) {
	version := &models.AppVersion{}
	err := p.Db.
		Select(models.AppVersionColumns...).
		From(models.AppVersionTableName).
		Where(db.Eq(models.ColumnVersionId, versionId)).
		LoadOne(&version)
	if err != nil {
		return nil, err
	}
	return version, nil
}

func (p *Server) getAppVersions(versionIds []string) ([]*models.AppVersion, error) {
	var versions []*models.AppVersion
	_, err := p.Db.
		Select(models.AppVersionColumns...).
		From(models.AppVersionTableName).
		Where(db.Eq(models.ColumnVersionId, versionIds)).
		Load(&versions)
	if err != nil {
		return nil, err
	}
	return versions, nil
}

func (p *Server) DescribeApps(ctx context.Context, req *pb.DescribeAppsRequest) (*pb.DescribeAppsResponse, error) {
	var apps []*models.App
	offset := pbutil.GetOffsetFromRequest(req)
	limit := pbutil.GetLimitFromRequest(req)
	categoryIds := req.GetCategoryId()

	query := p.Db.
		Select(models.AppColumns...).
		From(models.AppTableName).
		Offset(offset).
		Limit(limit).
		Where(manager.BuildFilterConditions(req, models.AppTableName))
	if len(categoryIds) > 0 {
		subqueryStmt := p.Db.
			Select(models.ColumnResouceId).
			From(models.CategoryResourceTableName).
			Where(db.Eq(models.ColumnCategoryId, categoryIds))
		query = query.Where(db.Eq(models.ColumnAppId, []*db.SelectQuery{subqueryStmt}))
	}
	// TODO: validate sort_key
	query = manager.AddQueryOrderDir(query, req, models.ColumnCreateTime)
	_, err := query.Load(&apps)
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}
	count, err := query.Count()
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}

	appSet, err := p.formatAppSet(apps)
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
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

	_, err := p.Db.
		InsertInto(models.AppTableName).
		Columns(models.AppColumns...).
		Record(newApp).
		Exec()
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
	}

	err = categoryutil.SyncResourceCategories(p.Db, newApp.AppId, strings.Split(req.GetCategoryId().GetValue(), ","))
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
	}

	res := &pb.CreateAppResponse{
		AppId: pbutil.ToProtoString(newApp.AppId),
	}
	return res, nil
}

func (p *Server) ModifyApp(ctx context.Context, req *pb.ModifyAppRequest) (*pb.ModifyAppResponse, error) {
	// TODO: check resource permission
	appId := req.GetAppId().GetValue()
	app, err := p.getApp(appId)
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}
	if app.Status == constants.StatusDeleted {
		return nil, gerr.NewWithDetail(gerr.FailedPrecondition, err, gerr.ErrorResourceAlreadyDeleted, appId)
	}

	attributes := manager.BuildUpdateAttributes(req,
		"name", "repo_id", "owner", "chart_name",
		"description", "home", "icon", "screenshots",
		"maintainers", "sources", "readme", "keywords")
	attributes["update_time"] = time.Now()
	_, err = p.Db.
		Update(models.AppTableName).
		SetMap(attributes).
		Where(db.Eq(models.ColumnAppId, appId)).
		Exec()
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorModifyResourcesFailed)
	}

	err = categoryutil.SyncResourceCategories(p.Db, appId, strings.Split(req.GetCategoryId().GetValue(), ","))
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorModifyResourcesFailed)
	}

	res := &pb.ModifyAppResponse{
		AppId: req.GetAppId(),
	}
	return res, nil
}

func (p *Server) DeleteApps(ctx context.Context, req *pb.DeleteAppsRequest) (*pb.DeleteAppsResponse, error) {
	// TODO: check resource permission
	appIds := req.GetAppId()

	_, err := p.Db.
		Update(models.AppTableName).
		Set(models.ColumnStatus, constants.StatusDeleted).
		Where(db.Eq(models.ColumnAppId, appIds)).
		Exec()
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorDeleteResourcesFailed)
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
	}

	_, err := p.Db.
		InsertInto(models.AppVersionTableName).
		Columns(models.AppVersionColumns...).
		Record(newAppVersion).
		Exec()
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
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

	query := p.Db.
		Select(models.AppVersionColumns...).
		From(models.AppVersionTableName).
		Offset(offset).
		Limit(limit).
		Where(manager.BuildFilterConditions(req, models.AppVersionTableName))
	query = manager.AddQueryOrderDir(query, req, models.ColumnSequence)
	_, err := query.Load(&versions)
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}
	count, err := query.Count()
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}

	res := &pb.DescribeAppVersionsResponse{
		AppVersionSet: models.AppVersionsToPbs(versions),
		TotalCount:    count,
	}
	return res, nil

}

func (p *Server) ModifyAppVersion(ctx context.Context, req *pb.ModifyAppVersionRequest) (*pb.ModifyAppVersionResponse, error) {
	// TODO: check resource permission
	versionId := req.GetVersionId().GetValue()
	version, err := p.getAppVersion(versionId)
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}
	if version.Status == constants.StatusDeleted {
		return nil, gerr.NewWithDetail(gerr.FailedPrecondition, err, gerr.ErrorResourceAlreadyDeleted, versionId)
	}

	attributes := manager.BuildUpdateAttributes(req, "name", "description", "package_name", "sequence")
	_, err = p.Db.
		Update(models.AppVersionTableName).
		SetMap(attributes).
		Where(db.Eq(models.ColumnVersionId, versionId)).
		Exec()
	attributes["update_time"] = time.Now()
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorModifyResourcesFailed)
	}

	res := &pb.ModifyAppVersionResponse{
		VersionId: req.GetVersionId(),
	}
	return res, nil

}

func (p *Server) DeleteAppVersions(ctx context.Context, req *pb.DeleteAppVersionsRequest) (*pb.DeleteAppVersionsResponse, error) {
	// TODO: check resource permission
	versionIds := req.GetVersionId()

	_, err := p.Db.
		Update(models.AppVersionTableName).
		Set(models.ColumnStatus, constants.StatusDeleted).
		Where(db.Eq(models.ColumnVersionId, versionIds)).
		Exec()
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorDeleteResourceFailed, strings.Join(versionIds, ","))
	}

	return &pb.DeleteAppVersionsResponse{
		VersionId: versionIds,
	}, nil
}

func (p *Server) GetAppVersionPackage(ctx context.Context, req *pb.GetAppVersionPackageRequest) (*pb.GetAppVersionPackageResponse, error) {
	// TODO: check resource permission
	versionId := req.GetVersionId().GetValue()
	version, err := p.getAppVersion(versionId)
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.NotFound, err, gerr.ErrorResourceNotFound, versionId)
	}
	logger.Debug("Got app version: [%+v]", version)
	packageUrl := version.PackageName
	resp, err := httputil.HttpGet(packageUrl)
	if err != nil {
		logger.Error("Failed to http get [%s], error: %+v", packageUrl, err)
		return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorDescribeResourceFailed, versionId)
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error("Failed to read http response [%s], error: %+v", packageUrl, err)
		return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorDescribeResourceFailed, versionId)
	}
	return &pb.GetAppVersionPackageResponse{
		Package:   content,
		VersionId: req.GetVersionId(),
	}, nil
}

func (p *Server) GetAppVersionPackageFiles(ctx context.Context, req *pb.GetAppVersionPackageFilesRequest) (*pb.GetAppVersionPackageFilesResponse, error) {
	// TODO: check resource permission
	versionId := req.GetVersionId().GetValue()
	includeFiles := req.Files
	version, err := p.getAppVersion(versionId)
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.NotFound, err, gerr.ErrorResourceNotFound, versionId)
	}
	packageUrl := version.PackageName
	resp, err := httputil.HttpGet(packageUrl)
	if err != nil {
		logger.Error("Failed to http get [%s], error: %+v", packageUrl, err)
		return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorDescribeResourceFailed, versionId)
	}
	archiveFiles, err := gziputil.LoadArchive(resp.Body, includeFiles...)
	if err != nil {
		logger.Error("Failed to load package [%s] archive, error: %+v", packageUrl, err)
		return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorDescribeResourceFailed, versionId)
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
	appCount, err := p.Db.Select(models.ColumnAppId).From(models.AppTableName).Count()
	if err != nil {
		logger.Error("Failed to get app count, error: %+v", err)
		return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}
	res.AppCount = appCount

	err = p.Db.Select("COUNT(DISTINCT repo_id)").From(models.AppTableName).LoadOne(&res.RepoCount)
	if err != nil {
		logger.Error("Failed to get repo count, error: %+v", err)
		return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}

	time2week := time.Now().Add(-14 * 24 * time.Hour)
	var as []*appStatistic
	_, err = p.Db.
		Select("DATE_FORMAT(create_time, '%Y-%m-%d')", "COUNT(app_id)").
		From(models.AppTableName).
		GroupBy("DATE_FORMAT(create_time, '%Y-%m-%d')").
		Where(db.Gte(models.ColumnCreateTime, time2week)).
		Limit(14).Load(&as)

	if err != nil {
		logger.Error("Failed to get app statistics, error: %+v", err)
		return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}
	for _, a := range as {
		res.LastTwoWeekCreated[a.Date] = a.Count
	}

	var rs []*repoStatistic
	_, err = p.Db.
		Select("repo_id", "COUNT(app_id)").
		From(models.AppTableName).
		GroupBy(models.ColumnRepoId).
		OrderDir("COUNT(app_id)", false).
		Limit(10).Load(&rs)

	if err != nil {
		logger.Error("Failed to get repo statistics, error: %+v", err)
		return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}
	for _, a := range rs {
		res.TopTenRepos[a.RepoId] = a.Count
	}

	return res, nil
}
