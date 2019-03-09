// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"time"

	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/sender"
	"openpitrix.io/openpitrix/pkg/util/idutil"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
)

const UncategorizedId = "ctg-uncategorized"

func NewCategoryId() string {
	return idutil.GetUuid("ctg-")
}

type Category struct {
	CategoryId  string
	Name        string
	Description string
	Locale      string
	OwnerPath   sender.OwnerPath
	Owner       string
	Icon        string
	CreateTime  time.Time
	UpdateTime  *time.Time
}

var CategoryColumns = db.GetColumnsFromStruct(&Category{})

func NewCategory(name, locale, description string, ownerPath sender.OwnerPath) *Category {
	if locale == "" {
		locale = "{}"
	}
	return &Category{
		CategoryId:  NewCategoryId(),
		Name:        name,
		Locale:      locale,
		Description: description,
		Owner:       ownerPath.Owner(),
		OwnerPath:   ownerPath,
		CreateTime:  time.Now(),
	}
}

func CategoryToPb(category *Category) *pb.Category {
	pbCategory := pb.Category{}
	pbCategory.CategoryId = pbutil.ToProtoString(category.CategoryId)
	pbCategory.Name = pbutil.ToProtoString(category.Name)
	pbCategory.Locale = pbutil.ToProtoString(category.Locale)
	pbCategory.OwnerPath = category.OwnerPath.ToProtoString()
	pbCategory.Owner = pbutil.ToProtoString(category.Owner)
	pbCategory.Description = pbutil.ToProtoString(category.Description)
	pbCategory.CreateTime = pbutil.ToProtoTimestamp(category.CreateTime)
	pbCategory.Icon = pbutil.ToProtoString(category.Icon)
	if category.UpdateTime != nil {
		pbCategory.UpdateTime = pbutil.ToProtoTimestamp(*category.UpdateTime)
	}
	return &pbCategory
}

func CategoriesToPbs(categories []*Category) (pbCategories []*pb.Category) {
	for _, app := range categories {
		pbCategories = append(pbCategories, CategoryToPb(app))
	}
	return
}
