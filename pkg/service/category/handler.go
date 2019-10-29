// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package category

import (
	"context"
	"time"

	attachmentclient "openpitrix.io/openpitrix/pkg/client/attachment"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/gerr"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/pi"
	"openpitrix.io/openpitrix/pkg/util/ctxutil"
	"openpitrix.io/openpitrix/pkg/util/imageutil"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
	"openpitrix.io/openpitrix/pkg/util/stringutil"
)

func (p *Server) DescribeCategories(ctx context.Context, req *pb.DescribeCategoriesRequest) (*pb.DescribeCategoriesResponse, error) {
	var categories []*models.Category
	offset := pbutil.GetOffsetFromRequest(req)
	limit := pbutil.GetLimitFromRequest(req)

	displayColumns := manager.GetDisplayColumns(req.GetDisplayColumns(), models.CategoryColumns)
	query := pi.Global().DB(ctx).
		Select(displayColumns...).
		From(constants.TableCategory).
		Offset(offset).
		Limit(limit).
		Where(manager.BuildPermissionFilter(ctx)).
		Where(manager.BuildFilterConditions(req, constants.TableCategory))
	// TODO: validate sort_key
	query = manager.AddQueryOrderDir(query, req, constants.ColumnCreateTime)
	if len(displayColumns) > 0 {
		_, err := query.Load(&categories)
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
		}
	}
	count, err := query.Count()
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}

	res := &pb.DescribeCategoriesResponse{
		CategorySet: models.CategoriesToPbs(categories),
		TotalCount:  count,
	}
	return res, nil
}

func (p *Server) CreateCategory(ctx context.Context, req *pb.CreateCategoryRequest) (*pb.CreateCategoryResponse, error) {
	s := ctxutil.GetSender(ctx)
	category := models.NewCategory(
		req.GetName().GetValue(),
		req.GetLocale().GetValue(),
		req.GetDescription().GetValue(),
		s.GetOwnerPath())

	var iconAttachmentId string
	if req.GetIcon() != nil {
		// upload icon attachment
		attachmentManagerClient, err := attachmentclient.NewAttachmentManagerClient()
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
		}
		icon := req.GetIcon().GetValue()
		content, err := imageutil.Thumbnail(ctx, icon)
		if err != nil {
			logger.Error(ctx, "Make thumbnail failed: %+v", err)
			return nil, gerr.NewWithDetail(ctx, gerr.InvalidArgument, err, gerr.ErrorImageDecodeFailed)
		}
		createAttachmentRes, err := attachmentManagerClient.CreateAttachment(ctx, &pb.CreateAttachmentRequest{
			AttachmentContent: content,
		})
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
		}
		iconAttachmentId = createAttachmentRes.AttachmentId
	}
	category.Icon = iconAttachmentId

	_, err := pi.Global().DB(ctx).
		InsertInto(constants.TableCategory).
		Record(category).
		Exec()
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
	}

	res := &pb.CreateCategoryResponse{
		CategoryId: pbutil.ToProtoString(category.CategoryId),
	}
	return res, nil
}

func (p *Server) ModifyCategory(ctx context.Context, req *pb.ModifyCategoryRequest) (*pb.ModifyCategoryResponse, error) {
	// TODO: check resource permission
	categoryId := req.GetCategoryId().GetValue()
	category, err := p.getCategory(ctx, categoryId)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}

	attributes := manager.BuildUpdateAttributes(req,
		constants.ColumnName, constants.ColumnLocale, constants.ColumnDescription)
	attributes[constants.ColumnUpdateTime] = time.Now()

	if req.GetIcon() != nil {
		// upload or replace icon attachment
		attachmentManagerClient, err := attachmentclient.NewAttachmentManagerClient()
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
		}
		content, err := imageutil.Thumbnail(ctx, req.GetIcon().GetValue())
		if err != nil {
			logger.Error(ctx, "Make thumbnail failed: %+v", err)
			return nil, gerr.NewWithDetail(ctx, gerr.InvalidArgument, err, gerr.ErrorImageDecodeFailed)
		}
		if category.Icon == "" {
			createAttachmentRes, err := attachmentManagerClient.CreateAttachment(ctx, &pb.CreateAttachmentRequest{
				AttachmentContent: content,
			})
			if err != nil {
				return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
			}
			attributes[constants.ColumnIcon] = createAttachmentRes.AttachmentId
		} else {
			_, err := attachmentManagerClient.ReplaceAttachment(ctx, &pb.ReplaceAttachmentRequest{
				AttachmentId:      category.Icon,
				AttachmentContent: content,
			})
			if err != nil {
				return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
			}
		}
	}

	_, err = pi.Global().DB(ctx).
		Update(constants.TableCategory).
		SetMap(attributes).
		Where(db.Eq(constants.ColumnCategoryId, categoryId)).
		Exec()
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorModifyResourcesFailed)
	}
	res := &pb.ModifyCategoryResponse{
		CategoryId: req.GetCategoryId(),
	}
	return res, nil
}

func (p *Server) DeleteCategories(ctx context.Context, req *pb.DeleteCategoriesRequest) (*pb.DeleteCategoriesResponse, error) {
	var err error
	var count uint32
	var categoryIds []string

	categoryIds = req.GetCategoryId()
	if stringutil.StringIn(models.UncategorizedId, categoryIds) {
		return nil, gerr.New(ctx, gerr.PermissionDenied, gerr.ErrorCannotDeleteDefaultCategory)
	}
	count, err = countRelations(ctx, categoryIds)
	if err != nil {
		return nil, gerr.New(ctx, gerr.Internal, gerr.ErrorDeleteResourcesFailed)
	}
	if !req.Force.GetValue() && count > 0 {
		return nil, gerr.New(ctx, gerr.FailedPrecondition, gerr.ErrorDeleteResourcesFailed)
	}
	if count > 0 {
		err = deleteRelations(ctx, categoryIds)
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDeleteResourcesFailed)
		}
	}
	err = deleteCateogries(ctx, categoryIds)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDeleteResourcesFailed)
	}
	return &pb.DeleteCategoriesResponse{
		CategoryId: categoryIds,
	}, nil
}
