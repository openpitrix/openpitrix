// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package isv

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
	"openpitrix.io/openpitrix/pkg/util/ctxutil"
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
		Where(manager.BuildPermissionFilter(ctx)).
		Where(manager.BuildFilterConditions(req, constants.TableVendorVerifyInfo))

	query = manager.AddQueryOrderDir(query, req, constants.ColumnSubmitTime)
	if len(displayColumns) > 0 {
		_, err := query.Load(&vendors)
		if err != nil {
			logger.Error(ctx, "Failed to describe vendor verify info: %+v", err)
			return nil, 0, err
		}
	}

	count, err := query.Count()
	if err != nil {
		logger.Error(ctx, "Failed to describe vendor verify info count: %+v", err)
		return nil, 0, err
	}

	return vendors, count, err
}

func PassVendorVerifyInfo(ctx context.Context, appVendorUserId string) error {
	sender := ctxutil.GetSender(ctx)
	approver := sender.UserId

	_, err := pi.Global().DB(ctx).
		Update(constants.TableVendorVerifyInfo).
		Set(constants.ColumnStatus, constants.StatusPassed).
		Set(constants.ColumnApprover, approver).
		Set(constants.ColumnStatusTime, time.Now()).
		Where(db.Eq(constants.ColumnUserId, appVendorUserId)).
		Exec()
	if err != nil {
		logger.Error(ctx, "Failed to pass vendor [%s] verify info: %+v", appVendorUserId, err)
		return err
	}
	return err
}

func RejectVendorVerifyInfo(ctx context.Context, appVendorUserId string, rejectMsg string, approver string) error {
	_, err := pi.Global().DB(ctx).
		Update(constants.TableVendorVerifyInfo).
		Set(constants.ColumnStatus, constants.StatusRejected).
		Set(constants.ColumnStatusTime, time.Now()).
		Set(constants.ColumnRejectMessage, rejectMsg).
		Set(constants.ColumnApprover, approver).
		Where(db.Eq(constants.ColumnUserId, appVendorUserId)).
		Exec()
	if err != nil {
		logger.Error(ctx, "Failed to reject vendor [%s] verify info: %+v", appVendorUserId, err)
		return err
	}
	return err
}

func UpdateVendorVerifyInfo(ctx context.Context, req *pb.SubmitVendorVerifyInfoRequest) error {
	appVendorUserId := req.UserId
	attributes := manager.BuildUpdateAttributes(req, constants.ColumnCompanyName, constants.ColumnCompanyWebsite, constants.ColumnCompanyProfile,
		constants.ColumnAuthorizerName, constants.ColumnAuthorizerEmail, constants.ColumnAuthorizerPhone, constants.ColumnBankName, constants.ColumnBankAccountName,
		constants.ColumnBankAccountNumber)
	attributes[constants.ColumnStatus] = constants.StatusSubmitted
	attributes[constants.ColumnSubmitTime] = time.Now()
	attributes[constants.ColumnStatusTime] = time.Now()

	logger.Debug(ctx, "Update vendor verify info attributes: [%+v]", attributes)

	var err error
	if len(attributes) != 0 {
		_, err = pi.Global().DB(ctx).
			Update(constants.TableVendorVerifyInfo).
			SetMap(attributes).
			Where(db.Eq(constants.ColumnUserId, appVendorUserId)).
			Exec()
		if err != nil {
			logger.Error(ctx, "Failed to update vendor [%s] verify info: %+v", appVendorUserId, err)
			return err
		}
	}
	return nil
}

func CreateVendorVerifyInfo(ctx context.Context, vendor *models.VendorVerifyInfo) error {
	_, err := pi.Global().DB(ctx).
		InsertInto(constants.TableVendorVerifyInfo).
		Record(vendor).
		Exec()
	if err != nil {
		logger.Error(ctx, "Failed to create vendor verify info: %+v", err)
		return err
	}
	return nil
}

func GetVendorVerifyInfoCountByCompanyName(ctx context.Context, companyName string) (uint32, error) {
	count, err := pi.Global().DB(ctx).
		Select(constants.ColumnCompanyName).
		From(constants.TableVendorVerifyInfo).
		Where(db.Eq(constants.ColumnCompanyName, companyName)).
		Count()
	if err != nil {
		logger.Error(ctx, "Failed to get company name count: %+v", err)
		return 0, err
	}
	return count, nil
}

func GetVendorVerifyInfo(ctx context.Context, appVendorUserId string) (*models.VendorVerifyInfo, error) {
	vendor := &models.VendorVerifyInfo{}
	vendorColumns := db.GetColumnsFromStruct(&models.VendorVerifyInfo{})
	err := pi.Global().DB(ctx).
		Select(vendorColumns...).
		From(constants.TableVendorVerifyInfo).
		Where(db.Eq(constants.ColumnUserId, appVendorUserId)).
		LoadOne(&vendor)
	if err != nil {
		logger.Error(ctx, "Failed to get vendor [%s] verify info: %+v", appVendorUserId, err)
		return nil, err
	}
	return vendor, err
}
