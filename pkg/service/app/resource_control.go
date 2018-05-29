// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package app

import (
	"fmt"

	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
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

func (p *Server) getCategory(categoryId string) (*models.Category, error) {
	categories, err := p.getCategories([]string{categoryId})
	if err != nil {
		return nil, err
	}
	if len(categories) == 0 {
		logger.Error("Failed to get category [%s]", categoryId)
		return nil, fmt.Errorf("failed to get category [%s]", categoryId)
	}
	return categories[0], nil
}

func (p *Server) getCategories(categoryIds []string) ([]*models.Category, error) {
	var categories []*models.Category
	_, err := p.Db.
		Select(models.CategoryColumns...).
		From(models.CategoryTableName).
		Where(db.Eq("category_id", categoryIds)).
		Load(&categories)
	if err != nil {
		return nil, err
	}
	return categories, nil
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

func (p *Server) getAppCategorySet(appId string) ([]*pb.AppCategory, error) {
	var categoryResouces []*models.CategoryResource
	_, err := p.Db.
		Select(models.CategoryResourceColumns...).
		From(models.CategoryResourceTableName).
		Where(db.Eq(models.ColumnResouceId, appId)).
		Load(&categoryResouces)
	if err != nil {
		return nil, err
	}
	var categoryResourceMap = make(map[string]*models.CategoryResource)
	var categoryIds []string
	for _, categoryResouce := range categoryResouces {
		categoryResourceMap[categoryResouce.CategoryId] = categoryResouce
		categoryIds = append(categoryIds, categoryResouce.CategoryId)
	}
	categories, err := p.getCategories(categoryIds)
	if err != nil {
		return nil, err
	}
	var pbCategories []*pb.AppCategory
	for _, category := range categories {
		categoryResouce := categoryResourceMap[category.CategoryId]
		pbCategories = append(
			pbCategories,
			&pb.AppCategory{
				CategoryId: pbutil.ToProtoString(category.CategoryId),
				Name:       pbutil.ToProtoString(category.Name),
				Locale:     pbutil.ToProtoString(category.Locale),
				Status:     pbutil.ToProtoString(categoryResouce.Status),
				CreateTime: pbutil.ToProtoTimestamp(categoryResouce.CreateTime),
				StatusTime: pbutil.ToProtoTimestamp(categoryResouce.StatusTime),
			},
		)
	}

	return pbCategories, nil
}

func (p *Server) formatApp(app *models.App) (*pb.App, error) {
	pbApp := models.AppToPb(app)

	latestAppVersion, err := p.getLatestAppVersion(app.AppId)
	if err != nil {
		return nil, err
	}
	pbApp.LatestAppVersion = models.AppVersionToPb(latestAppVersion)

	appCategorySet, err := p.getAppCategorySet(app.AppId)
	if err != nil {
		return nil, err
	}
	pbApp.AppCategorySet = appCategorySet

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
