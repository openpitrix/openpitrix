// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package categoryutil

import (
	"context"
	"strings"
	"time"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/util/stringutil"
)

func DecodeCategoryIds(s string) []string {
	if len(s) == 0 {
		return []string{models.UncategorizedId}
	}
	return strings.Split(s, ",")
}

func SyncResourceCategories(ctx context.Context, d *db.Conn, appId string, categoryIds []string) error {
	if len(categoryIds) == 0 {
		categoryIds = append(categoryIds, models.UncategorizedId)
	}
	var existCategoryIds []string
	_, err := d.
		Select(models.ColumnCategoryId).
		From(models.CategoryResourceTableName).
		Where(db.Eq(models.ColumnResouceId, appId)).
		Load(&existCategoryIds)
	if err != nil {
		logger.Error(ctx, "Failed to load resource [%s] categories", appId)
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
		updateStmt := d.
			Update(models.CategoryResourceTableName).
			Set(models.ColumnStatus, constants.StatusDisabled).
			Set(models.ColumnStatusTime, time.Now()).
			Where(db.Eq(models.ColumnResouceId, appId)).
			Where(db.Eq(models.ColumnCategoryId, disableIds))
		_, err = updateStmt.Exec()
		if err != nil {
			logger.Error(ctx, "Failed to set resource [%s] categories [%s] to disabled", appId, disableIds)
			return err
		}
	}
	if len(enableIds) > 0 {
		updateStmt := d.
			Update(models.CategoryResourceTableName).
			Set(models.ColumnStatus, constants.StatusEnabled).
			Set(models.ColumnStatusTime, time.Now()).
			Where(db.Eq(models.ColumnResouceId, appId)).
			Where(db.Eq(models.ColumnCategoryId, enableIds))
		_, err = updateStmt.Exec()
		if err != nil {
			logger.Error(ctx, "Failed to set resource [%s] categories [%s] to enabled", appId, enableIds)
			return err
		}
	}
	if len(createIds) > 0 {
		insertStmt := d.
			InsertInto(models.CategoryResourceTableName).
			Columns(models.CategoryResourceColumns...)
		for _, categoryId := range createIds {
			insertStmt = insertStmt.Record(
				models.NewCategoryResource(categoryId, appId, constants.StatusEnabled),
			)
		}
		_, err = insertStmt.Exec()
		if err != nil {
			logger.Error(ctx, "Failed to create resource [%s] categories [%s]", appId, createIds)
			return err
		}
	}

	return nil
}
