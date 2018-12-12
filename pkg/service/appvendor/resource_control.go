// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package appvendor

import (
	"context"
	"time"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/pi"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
)

type Resource struct {
}

func (resource *Resource) DescribeVendorVerifyInfos(ctx context.Context, req *pb.DescribeVendorVerifyInfosRequest) ([]*models.AppVendor, uint32, error) {
	var vendors []*models.AppVendor
	offset := pbutil.GetOffsetFromRequest(req)
	limit := pbutil.GetLimitFromRequest(req)

	var vendorColumns = db.GetColumnsFromStruct(&models.AppVendor{})
	query := pi.Global().DB(ctx).
		Select(vendorColumns...).
		From(constants.TableVendor).
		Offset(offset).
		Limit(limit).
		Where(manager.BuildFilterConditions(req, constants.TableVendor))

	query = manager.AddQueryOrderDir(query, req, "submit_time")
	_, err := query.Load(&vendors)
	count, err := query.Count()

	if err != nil {
		logger.Error(ctx, "Failed to Describe VendorVerifyInfos: [%+v]", req)
		return nil, 0, err
	}

	return vendors, count, err
}

func (resource *Resource) GetVendorVerifyInfo(ctx context.Context, userID string) (*models.AppVendor, error) {
	vendor := &models.AppVendor{}
	vendorColumns := db.GetColumnsFromStruct(&models.AppVendor{})
	err := pi.Global().DB(ctx).
		Select(vendorColumns...).
		From(constants.TableVendor).
		Where(db.Eq(constants.ColumnUserId, userID)).
		LoadOne(&vendor)

	if err != nil {
		logger.Error(ctx, "Failed to Get VendorVerifyInfo:userID= [%s]", userID)
		return nil, err
	}
	return vendor, err
}

func (resource *Resource) PassVendorVerifyInfo(ctx context.Context, userID string) (string, error) {
	_, err := pi.Global().DB(ctx).
		Update(constants.TableVendor).
		Set("status", "passed").
		Set("status_time", time.Now()).
		Where(db.Eq(constants.ColumnUserId, userID)).
		Exec()
	if err != nil {
		logger.Error(ctx, "Failed to Pass VendorVerifyInfo：userID= [%s]", userID)
		return "", err
	}
	return userID, err
}

func (resource *Resource) RejectVendorVerifyInfo(ctx context.Context, userID string, rejectmsg string) (string, error) {
	_, err := pi.Global().DB(ctx).
		Update(constants.TableVendor).
		Set("status", "rejected").
		Set("status_time", time.Now()).
		Set("reject_message", rejectmsg).
		Where(db.Eq(constants.ColumnUserId, userID)).
		Exec()
	if err != nil {
		logger.Error(ctx, "Failed to Reject VendorVerifyInfo:userID= [%s]", userID)
		return "", err
	}
	return userID, err
}

func (resource *Resource) UpdateVendorVerifyInfo(ctx context.Context, userID string, attributes map[string]interface{}) (string, error) {
	var err error
	if len(attributes) != 0 {
		_, err = pi.Global().DB(ctx).
			Update(constants.TableVendor).
			SetMap(attributes).
			Where(db.Eq("user_id", userID)).
			Exec()
		if err != nil {
			logger.Error(ctx, "Failed to update appvendor:userID= [%s]", userID)
			return "", err
		}
	}
	return userID, nil
}

func (resource *Resource) CreateVendorVerifyInfo(ctx context.Context, vendor models.AppVendor) (string, error) {
	_, err := pi.Global().DB(ctx).
		InsertInto(constants.TableVendor).
		Record(vendor).
		Exec()
	if err != nil {
		logger.Error(ctx, "Failed to Create appvendor [%+v]", vendor)
		return "", err
	}
	return vendor.UserId, nil
}
