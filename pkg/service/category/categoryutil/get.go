// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package categoryutil

import (
	"openpitrix.io/openpitrix/pkg/client"
	categoryclient "openpitrix.io/openpitrix/pkg/client/category"
	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
	"openpitrix.io/openpitrix/pkg/util/stringutil"
)

func GetResourcesCategories(d *db.Database, resourceIds []string) (map[string][]*pb.ResourceCategory, error) {
	var rcmap = make(map[string][]*pb.ResourceCategory)
	var categoryResources []*models.CategoryResource
	_, err := d.Select(models.CategoryResourceColumns...).
		From(models.CategoryResourceTableName).
		Where(db.Eq(models.ColumnResouceId, resourceIds)).
		Load(&categoryResources)
	if err != nil {
		return rcmap, err
	}
	var categoryIds []string
	for _, r := range categoryResources {
		categoryIds = append(categoryIds, r.CategoryId)
		var categorySet []*pb.ResourceCategory
		if cs, ok := rcmap[r.ResourceId]; ok {
			categorySet = cs
		}
		categorySet = append(categorySet, &pb.ResourceCategory{
			CategoryId: pbutil.ToProtoString(r.CategoryId),
			Status:     pbutil.ToProtoString(r.Status),
			CreateTime: pbutil.ToProtoTimestamp(r.CreateTime),
			StatusTime: pbutil.ToProtoTimestamp(r.StatusTime),
		})
		rcmap[r.ResourceId] = categorySet
	}
	ctx := client.GetSystemUserContext()
	c, err := categoryclient.NewCategoryManagerClient()
	if err != nil {
		return rcmap, err
	}
	descParams := pb.DescribeCategoriesRequest{
		CategoryId: stringutil.Unique(categoryIds),
	}
	resp, err := c.DescribeCategories(ctx, &descParams)
	if err != nil {
		return rcmap, err
	}
	categories := resp.CategorySet
	categoryMap := make(map[string]*pb.Category)
	for _, category := range categories {
		categoryMap[category.GetCategoryId().GetValue()] = category
	}
	for _, categorySet := range rcmap {
		for _, rCategory := range categorySet {
			category := categoryMap[rCategory.GetCategoryId().GetValue()]
			rCategory.Name = category.GetName()
			rCategory.Locale = category.GetLocale()
		}
	}
	return rcmap, nil
}
