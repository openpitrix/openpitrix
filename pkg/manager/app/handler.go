// Copyright 2017 The OpenPitrix Authors. All rights reserved.
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
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/utils"
	"openpitrix.io/openpitrix/pkg/utils/sender"
)

func (p *Server) getApp(appId string) (*models.App, error) {
	app := &models.App{}
	err := p.db.
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

	query := p.db.
		Select(models.AppColumns...).
		From(models.AppTableName).
		Offset(offset).
		Limit(limit)
	// TODO: filter condition
	if len(req.GetName()) > 0 {
		query = query.Where(db.Eq("name", req.GetName()))
	}
	if len(req.GetRepoId()) > 0 {
		query = query.Where(db.Eq("repo_id", req.GetRepoId()))
	}
	if len(req.GetAppId()) > 0 {
		query = query.Where(db.Eq("app_id", req.GetAppId()))
	}
	if len(req.GetStatus()) > 0 {
		query = query.Where(db.Eq("status", req.GetStatus()))
	}

	count, err := query.Load(&apps)
	if err != nil {
		// TODO: err_code should be implementation
		return nil, status.Errorf(codes.Internal, "DescribeApps: %+v", err)
	}

	res := &pb.DescribeAppsResponse{
		AppSet:     models.AppsToPbs(apps),
		TotalCount: uint32(count),
	}
	return res, nil
}

func (p *Server) CreateApp(ctx context.Context, req *pb.CreateAppRequest) (*pb.CreateAppResponse, error) {
	// TODO: validate CreateAppRequest
	s := sender.GetSenderFromContext(ctx)
	newApp := models.NewApp(req.GetName(), req.GetRepoId(), req.GetDescription(), s.UserId)

	_, err := p.db.
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
	appId := req.GetAppId()
	app, err := p.getApp(appId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get app [%s]", appId)
	}

	// TODO: use reflect parse attributes
	attributes := make(map[string]interface{})
	if len(req.Name) > 0 {
		attributes["name"] = req.Name
	}
	if len(req.RepoId) > 0 {
		attributes["repo_id"] = req.RepoId
	}
	if len(req.Owner) > 0 {
		attributes["owner"] = req.Owner
	}
	if len(req.ChartName) > 0 {
		attributes["chart_name"] = req.ChartName
	}
	if len(req.Description) > 0 {
		attributes["description"] = req.Description
	}
	if len(req.Home) > 0 {
		attributes["home"] = req.Home
	}
	if len(req.Icon) > 0 {
		attributes["icon"] = req.Icon
	}
	if len(req.Screenshots) > 0 {
		attributes["screenshots"] = req.Screenshots
	}
	if len(req.Maintainers) > 0 {
		attributes["maintainers"] = req.Maintainers
	}
	if len(req.Sources) > 0 {
		attributes["sources"] = req.Sources
	}
	if len(req.Readme) > 0 {
		attributes["readme"] = req.Readme
	}

	_, err = p.db.
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
	appId := req.GetAppId()
	_, err := p.getApp(appId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get app [%s]", appId)
	}

	_, err = p.db.
		Update(models.AppTableName).
		Set("status", constants.StatusDeleted).
		Where(db.Eq("app_id", appId)).
		Exec()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "DeleteApp: %+v", err)
	}

	appId = req.GetAppId()
	app, err := p.getApp(appId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get app [%s]", appId)
	}

	return &pb.DeleteAppResponse{
		App: models.AppToPb(app),
	}, nil
}
