// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package app

import (
	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
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

func (p *Server) getApps(appIds []string) ([]*models.App, error) {
	var apps []*models.App
	_, err := p.Db.
		Select(models.AppColumns...).
		From(models.AppTableName).
		Where(db.Eq("app_id", appIds)).
		Load(&apps)
	if err != nil {
		return nil, err
	}
	return apps, nil
}

func (p *Server) getLatestAppVersion(appId string) (*models.AppVersion, error) {
	appVersion := &models.AppVersion{}
	err := p.Db.
		Select(models.AppVersionColumns...).
		From(models.AppVersionTableName).
		Where(db.Eq("app_id", appId)).
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

func (p *Server) formatApp(app *models.App) (*pb.App, error) {
	pbApp := models.AppToPb(app)

	latestAppVersion, err := p.getLatestAppVersion(app.AppId)
	if err != nil {
		return nil, err
	}

	pbApp.LatestAppVersion = models.AppVersionToPb(latestAppVersion)

	return pbApp, nil
}

func (p *Server) formatAppSet(apps []*models.App) (pbApps []*pb.App, err error) {
	for _, app := range apps {
		var pbApp *pb.App
		pbApp, err = p.formatApp(app)
		if err != nil {
			return
		}
		pbApps = append(pbApps, pbApp)
	}
	return
}
