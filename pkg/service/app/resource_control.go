// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package app

import (
	"time"

	"openpitrix.io/openpitrix/pkg/client"
	categoryclient "openpitrix.io/openpitrix/pkg/client/category"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
	"openpitrix.io/openpitrix/pkg/util/stringutil"
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

func (p *Server) syncAppCategories(appId string, categoryIds []string) error {
	var existCategoryIds []string
	_, err := p.Db.
		Select(models.ColumnCategoryId).
		From(models.CategoryResourceTableName).
		Where(db.Eq(models.ColumnResouceId, appId)).
		Load(&existCategoryIds)
	if err != nil {
		logger.Error("Failed to load app [%s] categories", appId)
		return err
	}
	disableIds := stringutil.Diff(existCategoryIds, categoryIds)
	createIds := stringutil.Diff(categoryIds, existCategoryIds)
	var enableIds []string
	for _, id := range existCategoryIds {
		if stringutil.StringIn(id, categoryIds) {
			enableIds = append(enableIds, id)
		}
	}
	if len(disableIds) > 0 {
		updateStmt := p.Db.
			Update(models.CategoryResourceTableName).
			Set(models.ColumnStatus, constants.StatusDisabled).
			Set(models.ColumnStatusTime, time.Now()).
			Where(db.Eq(models.ColumnResouceId, appId)).
			Where(db.Eq(models.ColumnCategoryId, disableIds))
		_, err = updateStmt.Exec()
		if err != nil {
			logger.Error("Failed to set app [%s] categories [%s] to disabled", appId, disableIds)
			return err
		}
	}
	if len(enableIds) > 0 {
		updateStmt := p.Db.
			Update(models.CategoryResourceTableName).
			Set(models.ColumnStatus, constants.StatusEnabled).
			Set(models.ColumnStatusTime, time.Now()).
			Where(db.Eq(models.ColumnResouceId, appId)).
			Where(db.Eq(models.ColumnCategoryId, enableIds))
		_, err = updateStmt.Exec()
		if err != nil {
			logger.Error("Failed to set app [%s] categories [%s] to enabled", appId, enableIds)
			return err
		}
	}
	if len(createIds) > 0 {
		insertStmt := p.Db.
			InsertInto(models.CategoryResourceTableName).
			Columns(models.CategoryResourceColumns...)
		for _, categoryId := range createIds {
			insertStmt = insertStmt.Record(
				models.NewCategoryResource(categoryId, appId, constants.StatusEnabled),
			)
		}
		_, err = insertStmt.Exec()
		if err != nil {
			logger.Error("Failed to create app [%s] categories [%s]", appId, createIds)
			return err
		}
	}

	return nil
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
	var appIds []string
	var pbAppMap = make(map[string]*pb.App)
	for _, app := range apps {
		var pbApp *pb.App
		pbApp, err = p.formatApp(app)
		if err != nil {
			return
		}
		appIds = append(appIds, app.AppId)
		pbAppMap[app.AppId] = pbApp
	}

	var categoryResources []*models.CategoryResource
	_, err = p.Db.Select(models.CategoryResourceColumns...).
		From(models.CategoryResourceTableName).
		Where(db.Eq(models.ColumnResouceId, appIds)).
		Load(&categoryResources)
	if err != nil {
		return
	}
	var categoryIds []string
	for _, r := range categoryResources {
		categoryIds = append(categoryIds, r.CategoryId)
		pbApp := pbAppMap[r.ResourceId]
		pbApp.AppCategorySet = append(pbApp.AppCategorySet, &pb.AppCategory{
			CategoryId: pbutil.ToProtoString(r.CategoryId),
			Status:     pbutil.ToProtoString(r.Status),
			CreateTime: pbutil.ToProtoTimestamp(r.CreateTime),
			StatusTime: pbutil.ToProtoTimestamp(r.StatusTime),
		})
	}
	ctx := client.GetSystemUserContext()
	c, err := categoryclient.NewCategoryManagerClient(ctx)
	if err != nil {
		return pbApps, err
	}
	descParams := pb.DescribeCategoriesRequest{
		CategoryId: stringutil.Unique(categoryIds),
	}
	resp, err := c.DescribeCategories(ctx, &descParams)
	if err != nil {
		return pbApps, err
	}
	categories := resp.CategorySet
	categoryMap := make(map[string]*pb.Category)
	for _, category := range categories {
		categoryMap[category.GetCategoryId().GetValue()] = category
	}
	for _, pbApp := range pbAppMap {
		for _, appCategory := range pbApp.AppCategorySet {
			category := categoryMap[appCategory.GetCategoryId().GetValue()]
			appCategory.Name = category.GetName()
			appCategory.Locale = category.GetLocale()
		}
		pbApps = append(pbApps, pbApp)
	}
	return
}
