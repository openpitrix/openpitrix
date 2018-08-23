// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package app

import (
	"context"

	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/pi"
	"openpitrix.io/openpitrix/pkg/service/category/categoryutil"
)

func (p *Server) getApp(ctx context.Context, appId string) (*models.App, error) {
	app := &models.App{}
	err := pi.Global().DB(ctx).
		Select(models.AppColumns...).
		From(models.AppTableName).
		Where(db.Eq(models.ColumnAppId, appId)).
		LoadOne(&app)
	if err != nil {
		return nil, err
	}
	return app, nil
}

func (p *Server) getApps(ctx context.Context, appIds []string) ([]*models.App, error) {
	var apps []*models.App
	_, err := pi.Global().DB(ctx).
		Select(models.AppColumns...).
		From(models.AppTableName).
		Where(db.Eq(models.ColumnAppId, appIds)).
		Load(&apps)
	if err != nil {
		return nil, err
	}
	return apps, nil
}

func (p *Server) getLatestAppVersion(ctx context.Context, appId string) (*models.AppVersion, error) {
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

func (p *Server) formatApp(ctx context.Context, app *models.App) (*pb.App, error) {
	pbApp := models.AppToPb(app)

	latestAppVersion, err := p.getLatestAppVersion(ctx, app.AppId)
	if err != nil {
		return nil, err
	}
	pbApp.LatestAppVersion = models.AppVersionToPb(latestAppVersion)

	return pbApp, nil
}

func (p *Server) formatAppSet(ctx context.Context, apps []*models.App) ([]*pb.App, error) {
	var pbApps []*pb.App
	var appIds []string
	for _, app := range apps {
		var pbApp *pb.App
		pbApp, err := p.formatApp(ctx, app)
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
