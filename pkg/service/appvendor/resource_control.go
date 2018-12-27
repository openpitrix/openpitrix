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

func DescribeVendorVerifyInfos(ctx context.Context, req *pb.DescribeVendorVerifyInfosRequest) ([]*models.VendorVerifyInfo, uint32, error) {
	var vendors []*models.VendorVerifyInfo
	offset := pbutil.GetOffsetFromRequest(req)
	limit := pbutil.GetLimitFromRequest(req)

	var vendorColumns = db.GetColumnsFromStruct(&models.VendorVerifyInfo{})

	displayColumns := manager.GetDisplayColumns(req.GetDisplayColumns(), vendorColumns)
	query := pi.Global().DB(ctx).
		Select(displayColumns...).
		From(constants.TableVendorVerifyInfo).
		Offset(offset).
		Limit(limit).
		Where(manager.BuildFilterConditions(req, constants.TableVendorVerifyInfo))

	query = manager.AddQueryOrderDir(query, req, "submit_time")
	if len(displayColumns) > 0 {
		_, err := query.Load(&vendors)
		if err != nil {
			logger.Error(ctx, "Failed to describe vendorVerifyInfos [%v], error: %+v.", req, err)
			return nil, 0, err
		}
	}

	count, err := query.Count()
	if err != nil {
		logger.Error(ctx, "Failed to describe vendorVerifyInfos count [%v], error: %+v.", req, err)
		return nil, 0, err
	}

	return vendors, count, err
}

func GetVendorVerifyInfo(ctx context.Context, userID string) (*models.VendorVerifyInfo, error) {
	vendor := &models.VendorVerifyInfo{}
	vendorColumns := db.GetColumnsFromStruct(&models.VendorVerifyInfo{})
	err := pi.Global().DB(ctx).
		Select(vendorColumns...).
		From(constants.TableVendorVerifyInfo).
		Where(db.Eq(constants.ColumnUserId, userID)).
		LoadOne(&vendor)
	if err != nil {
		//logger.Error(ctx, "Failed to get vendorVerifyInfo [%s], vendorVerifyInfo not exists.", userID)
		return nil, err
	}
	return vendor, err
}

func PassVendorVerifyInfo(ctx context.Context, userID string) (string, error) {
	_, err := pi.Global().DB(ctx).
		Update(constants.TableVendorVerifyInfo).
		Set(constants.ColumnStatus, constants.StatusPassed).
		Set(constants.ColumnStatusTime, time.Now()).
		Where(db.Eq(constants.ColumnUserId, userID)).
		Exec()
	if err != nil {
		logger.Error(ctx, "Failed to pass vendorVerifyInfo [%s].", userID)
		return "", err
	}
	return userID, err
}

func RejectVendorVerifyInfo(ctx context.Context, userID string, rejectmsg string) (string, error) {
	_, err := pi.Global().DB(ctx).
		Update(constants.TableVendorVerifyInfo).
		Set(constants.ColumnStatus, "rejected").
		Set(constants.ColumnStatusTime, time.Now()).
		Set(constants.ColumnRejectMessage, rejectmsg).
		Where(db.Eq(constants.ColumnUserId, userID)).
		Exec()
	if err != nil {
		logger.Error(ctx, "Failed to reject vendorVerifyInfo [%s].", userID)
		return "", err
	}
	return userID, err
}

func UpdateVendorVerifyInfo(ctx context.Context, userID string, attributes map[string]interface{}) (string, error) {
	var err error
	if len(attributes) != 0 {
		_, err = pi.Global().DB(ctx).
			Update(constants.TableVendorVerifyInfo).
			SetMap(attributes).
			Where(db.Eq(constants.ColumnUserId, userID)).
			Exec()
		if err != nil {
			logger.Error(ctx, "Failed to update vendorVerifyInfo [%s].", userID)
			return "", err
		}
	}
	return userID, nil
}

func CreateVendorVerifyInfo(ctx context.Context, vendor models.VendorVerifyInfo) (string, error) {
	_, err := pi.Global().DB(ctx).
		InsertInto(constants.TableVendorVerifyInfo).
		Record(vendor).
		Exec()
	if err != nil {
		logger.Error(ctx, "Failed to create vendorVerifyInfo [%+v].", vendor)
		return "", err
	}
	return vendor.UserId, nil
}
