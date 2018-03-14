// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package app

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/utils"
	"openpitrix.io/openpitrix/pkg/utils/sender"
)

func (p *Server) getApp(appId string) (*models.App, error) {
	app := &models.App{}
	err := p.Db.
		Select(models.AppColumns...).
		From(models.AppTableName).
		Where(db.Eq("app_id", appId)).
		LoadOne(&app)
	if err != nil {
		return nil, err
	}
	return app, nil
}

func (p *Server) DescribeApps(ctx context.Context, req *pb.DescribeAppsRequest) (*pb.DescribeAppsResponse, error) {
	s := sender.GetSenderFromContext(ctx)
	logger.Infof("Got sender: %+v", s)
	logger.Debugf("Got req: %+v", req)
	var apps []*models.App
	offset := utils.GetOffsetFromRequest(req)
	limit := utils.GetLimitFromRequest(req)

	query := p.Db.
		Select(models.AppColumns...).
		From(models.AppTableName).
		Offset(offset).
		Limit(limit).
		Where(manager.BuildFilterConditions(req, models.AppTableName))
	_, err := query.Load(&apps)
	if err != nil {
		// TODO: err_code should be implementation
		return nil, status.Errorf(codes.Internal, "DescribeApps: %+v", err)
	}
	count, err := query.Count()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "DescribeApps: %+v", err)
	}

	res := &pb.DescribeAppsResponse{
		AppSet:     models.AppsToPbs(apps),
		TotalCount: count,
	}
	return res, nil
}

func (p *Server) CreateApp(ctx context.Context, req *pb.CreateAppRequest) (*pb.CreateAppResponse, error) {
	// TODO: validate CreateAppRequest
	s := sender.GetSenderFromContext(ctx)
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
		"maintainers", "sources", "readme")
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

func (p *Server) DeleteApp(ctx context.Context, req *pb.DeleteAppRequest) (*pb.DeleteAppResponse, error) {
	// TODO: check resource permission
	appId := req.GetAppId().GetValue()
	_, err := p.getApp(appId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get app [%s]", appId)
	}

	_, err = p.Db.
		Update(models.AppTableName).
		Set("status", constants.StatusDeleted).
		Where(db.Eq("app_id", appId)).
		Exec()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "DeleteApp: %+v", err)
	}

	app, err := p.getApp(appId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get app [%s]", appId)
	}

	return &pb.DeleteAppResponse{
		App: models.AppToPb(app),
	}, nil
}
