// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package app

import (
	"context"
	"io/ioutil"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
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
		Where(db.Eq("version_id", versionId)).
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
		Where(db.Eq("version_id", versionIds)).
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

	query := p.Db.
		Select(models.AppColumns...).
		From(models.AppTableName).
		Offset(offset).
		Limit(limit).
		Where(manager.BuildFilterConditions(req, models.AppTableName))
	// TODO: validate sort_key
	query = manager.AddQueryOrderDir(query, req, models.ColumnCreateTime)
	// TODO: add category_id join query
	_, err := query.Load(&apps)
	if err != nil {
		// TODO: err_code should be implementation
		return nil, status.Errorf(codes.Internal, "DescribeApps: %+v", err)
	}
	count, err := query.Count()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "DescribeApps: %+v", err)
	}

	appSet, err := p.formatAppSet(apps)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "DescribeApps: %+v", err)
	}

	res := &pb.DescribeAppsResponse{
		AppSet:     appSet,
		TotalCount: count,
	}
	return res, nil
}

func (p *Server) CreateApp(ctx context.Context, req *pb.CreateAppRequest) (*pb.CreateAppResponse, error) {
	// TODO: validate CreateAppRequest
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
		return nil, status.Errorf(codes.Internal, "CreateApp: %+v", err)
	}

	res := &pb.CreateAppResponse{
		App: models.AppToPb(newApp),
	}
	return res, nil
}

func (p *Server) ModifyApp(ctx context.Context, req *pb.ModifyAppRequest) (*pb.ModifyAppResponse, error) {
	// TODO: check resource permission
	appId := req.GetAppId().GetValue()
	app, err := p.getApp(appId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get app [%s]", appId)
	}

	attributes := manager.BuildUpdateAttributes(req,
		"name", "repo_id", "owner", "chart_name",
		"description", "home", "icon", "screenshots",
		"maintainers", "sources", "readme", "keywords")
	attributes["update_time"] = time.Now()
	_, err = p.Db.
		Update(models.AppTableName).
		SetMap(attributes).
		Where(db.Eq("app_id", appId)).
		Exec()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "ModifyApp: %+v", err)
	}
	app, err = p.getApp(appId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get app [%s]", appId)
	}

	res := &pb.ModifyAppResponse{
		App: models.AppToPb(app),
	}
	return res, nil
}

func (p *Server) DeleteApps(ctx context.Context, req *pb.DeleteAppsRequest) (*pb.DeleteAppsResponse, error) {
	// TODO: check resource permission
	err := manager.CheckParamsRequired(req, "app_id")
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	appIds := req.GetAppId()

	_, err = p.Db.
		Update(models.AppTableName).
		Set("status", constants.StatusDeleted).
		Where(db.Eq("app_id", appIds)).
		Exec()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "DeleteApps: %+v", err)
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
		return nil, status.Errorf(codes.Internal, "CreateAppVersion: %+v", err)
	}

	res := &pb.CreateAppVersionResponse{
		AppVersion: models.AppVersionToPb(newAppVersion),
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
		// TODO: err_code should be implementation
		return nil, status.Errorf(codes.Internal, "DescribeAppVersions: %+v", err)
	}
	count, err := query.Count()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "DescribeAppVersions: %+v", err)
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
		return nil, status.Errorf(codes.Internal, "Failed to get app version [%s]", versionId)
	}

	attributes := manager.BuildUpdateAttributes(req, "name", "description", "package_name", "sequence")
	_, err = p.Db.
		Update(models.AppVersionTableName).
		SetMap(attributes).
		Where(db.Eq("version_id", versionId)).
		Exec()
	attributes["update_time"] = time.Now()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "ModifyAppVersion: %+v", err)
	}
	version, err = p.getAppVersion(versionId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get app version [%s]", versionId)
	}

	res := &pb.ModifyAppVersionResponse{
		AppVersion: models.AppVersionToPb(version),
	}
	return res, nil

}

func (p *Server) DeleteAppVersions(ctx context.Context, req *pb.DeleteAppVersionsRequest) (*pb.DeleteAppVersionsResponse, error) {
	// TODO: check resource permission
	err := manager.CheckParamsRequired(req, "version_id")
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	versionIds := req.GetVersionId()

	_, err = p.Db.
		Update(models.AppVersionTableName).
		Set("status", constants.StatusDeleted).
		Where(db.Eq("version_id", versionIds)).
		Exec()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "DeleteAppVersions: %+v", err)
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
		return nil, status.Errorf(codes.Internal, "Failed to get app version [%s]", versionId)
	}
	logger.Debug("Got app version: [%+v]", version)
	packageUrl := version.PackageName
	resp, err := httputil.HttpGet(packageUrl)
	if err != nil {
		logger.Error("Failed to http get [%s], error: %+v", packageUrl, err)
		return nil, status.Errorf(codes.Internal, "Failed to get app version [%s]", versionId)
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error("Failed to read http response [%s], error: %+v", packageUrl, err)
		return nil, status.Errorf(codes.Internal, "Failed to get app version [%s]", versionId)
	}
	return &pb.GetAppVersionPackageResponse{
		Package: content,
	}, nil
}

func (p *Server) GetAppVersionPackageFiles(ctx context.Context, req *pb.GetAppVersionPackageFilesRequest) (*pb.GetAppVersionPackageFilesResponse, error) {
	// TODO: check resource permission
	versionId := req.GetVersionId().GetValue()
	includeFiles := req.Files
	version, err := p.getAppVersion(versionId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get app version [%s]", versionId)
	}
	packageUrl := version.PackageName
	resp, err := httputil.HttpGet(packageUrl)
	if err != nil {
		logger.Error("Failed to http get [%s], error: %+v", packageUrl, err)
		return nil, status.Errorf(codes.Internal, "Failed to http get [%s]", versionId)
	}
	archiveFiles, err := gziputil.LoadArchive(resp.Body, includeFiles...)
	if err != nil {
		logger.Error("Failed to load package [%s] archive, error: %+v", packageUrl, err)
		return nil, status.Errorf(codes.Internal, "Failed to load package [%s] archiv", versionId)
	}
	return &pb.GetAppVersionPackageFilesResponse{
		Files: archiveFiles,
	}, nil
}

func (p *Server) DescribeCategories(ctx context.Context, req *pb.DescribeCategoriesRequest) (*pb.DescribeCategoriesResponse, error) {
	var categories []*models.Category
	offset := pbutil.GetOffsetFromRequest(req)
	limit := pbutil.GetLimitFromRequest(req)

	query := p.Db.
		Select(models.CategoryColumns...).
		From(models.CategoryTableName).
		Offset(offset).
		Limit(limit).
		Where(manager.BuildFilterConditions(req, models.CategoryTableName))
	// TODO: validate sort_key
	query = manager.AddQueryOrderDir(query, req, models.ColumnCreateTime)
	_, err := query.Load(&categories)
	if err != nil {
		// TODO: err_code should be implementation
		return nil, status.Errorf(codes.Internal, "DescribeCategories: %+v", err)
	}
	count, err := query.Count()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "DescribeCategories: %+v", err)
	}

	res := &pb.DescribeCategoriesResponse{
		CategorySet: models.CategoriesToPbs(categories),
		TotalCount:  count,
	}
	return res, nil
}

func (p *Server) CreateCategory(ctx context.Context, req *pb.CreateCategoryRequest) (*pb.CreateCategoryResponse, error) {
	s := senderutil.GetSenderFromContext(ctx)
	category := models.NewCategory(req.GetName().GetValue(), req.GetLocale().GetValue(), s.UserId)

	_, err := p.Db.
		InsertInto(models.CategoryTableName).
		Columns(models.CategoryColumns...).
		Record(category).
		Exec()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "CreateCategory: %+v", err)
	}

	res := &pb.CreateCategoryResponse{
		Category: models.CategoryToPb(category),
	}
	return res, nil
}

func (p *Server) ModifyCategory(ctx context.Context, req *pb.ModifyCategoryRequest) (*pb.ModifyCategoryResponse, error) {
	// TODO: check resource permission
	categoryId := req.GetCategoryId().GetValue()
	category, err := p.getCategory(categoryId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get category [%s]", categoryId)
	}

	attributes := manager.BuildUpdateAttributes(req, "name", "locale")
	attributes["update_time"] = time.Now()
	_, err = p.Db.
		Update(models.CategoryTableName).
		SetMap(attributes).
		Where(db.Eq("category_id", categoryId)).
		Exec()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "ModifyCategory: %+v", err)
	}
	category, err = p.getCategory(categoryId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get category [%s]", categoryId)
	}

	res := &pb.ModifyCategoryResponse{
		Category: models.CategoryToPb(category),
	}
	return res, nil
}

func (p *Server) DeleteCategories(ctx context.Context, req *pb.DeleteCategoriesRequest) (*pb.DeleteCategoriesResponse, error) {
	err := manager.CheckParamsRequired(req, "category_id")
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	categoryIds := req.GetCategoryId()

	_, err = p.Db.
		DeleteFrom(models.CategoryTableName).
		Where(db.Eq("category_id", categoryIds)).
		Exec()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "DeleteCategories: %+v", err)
	}

	return &pb.DeleteCategoriesResponse{
		CategoryId: categoryIds,
	}, nil
}
