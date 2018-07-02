// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package category

import (
	"context"
	"time"

	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/gerr"
	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
	"openpitrix.io/openpitrix/pkg/util/senderutil"
)

func (p *Server) DescribeCategories(ctx context.Context, req *pb.DescribeCategoriesRequest) (*pb.DescribeCategoriesResponse, error) {
	var categories []*models.Category
	offset := pbutil.GetOffsetFromRequest(req)
	limit := pbutil.GetLimitFromRequest(req)

	query := p.Db.
		Select(models.CategoryColumns...).
		From(models.CategoryTableName).
		Offset(offset).
		Limit(limit).
		Where(manager.BuildFilterConditions(req, models.CategoryTableName))
	// TODO: validate sort_key
	query = manager.AddQueryOrderDir(query, req, models.ColumnCreateTime)
	_, err := query.Load(&categories)
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}
	count, err := query.Count()
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}

	res := &pb.DescribeCategoriesResponse{
		CategorySet: models.CategoriesToPbs(categories),
		TotalCount:  count,
	}
	return res, nil
}

func (p *Server) CreateCategory(ctx context.Context, req *pb.CreateCategoryRequest) (*pb.CreateCategoryResponse, error) {
	s := senderutil.GetSenderFromContext(ctx)
	category := models.NewCategory(
		req.GetName().GetValue(),
		req.GetLocale().GetValue(),
		req.GetDescription().GetValue(),
		s.UserId)

	_, err := p.Db.
		InsertInto(models.CategoryTableName).
		Columns(models.CategoryColumns...).
		Record(category).
		Exec()
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
	}

	res := &pb.CreateCategoryResponse{
		CategoryId: pbutil.ToProtoString(category.CategoryId),
	}
	return res, nil
}

func (p *Server) ModifyCategory(ctx context.Context, req *pb.ModifyCategoryRequest) (*pb.ModifyCategoryResponse, error) {
	// TODO: check resource permission
	categoryId := req.GetCategoryId().GetValue()
	_, err := p.getCategory(categoryId)
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}

	attributes := manager.BuildUpdateAttributes(req,
		models.ColumnName, models.ColumnLocale, models.ColumnDescription)
	attributes[models.ColumnUpdateTime] = time.Now()
	_, err = p.Db.
		Update(models.CategoryTableName).
		SetMap(attributes).
		Where(db.Eq(models.ColumnCategoryId, categoryId)).
		Exec()
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorModifyResourcesFailed)
	}
	res := &pb.ModifyCategoryResponse{
		CategoryId: req.GetCategoryId(),
	}
	return res, nil
}

func (p *Server) DeleteCategories(ctx context.Context, req *pb.DeleteCategoriesRequest) (*pb.DeleteCategoriesResponse, error) {
	categoryIds := req.GetCategoryId()

	_, err := p.Db.
		DeleteFrom(models.CategoryTableName).
		Where(db.Eq(models.ColumnCategoryId, categoryIds)).
		Exec()
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorDeleteResourcesFailed)
	}

	return &pb.DeleteCategoriesResponse{
		CategoryId: categoryIds,
	}, nil
}
