// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package category

import (
	"context"
	"fmt"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pi"
)

func (p *Server) getCategory(ctx context.Context, categoryId string) (*models.Category, error) {
	categories, err := p.getCategories(ctx, []string{categoryId})
	if err != nil {
		return nil, err
	}
	if len(categories) == 0 {
		logger.Error(ctx, "Failed to get category [%s]", categoryId)
		return nil, fmt.Errorf("failed to get category [%s]", categoryId)
	}
	return categories[0], nil
}

func (p *Server) getCategories(ctx context.Context, categoryIds []string) ([]*models.Category, error) {
	var categories []*models.Category
	_, err := pi.Global().DB(ctx).
		Select(models.CategoryColumns...).
		From(constants.TableCategory).
		Where(db.Eq(constants.ColumnCategoryId, categoryIds)).
		Load(&categories)
	if err != nil {
		return nil, err
	}
	return categories, nil
}
