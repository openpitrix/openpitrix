// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package category

import (
	"fmt"

	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/models"
)

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
