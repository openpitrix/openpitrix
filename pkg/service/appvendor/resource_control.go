package appvendor

import (
	"context"
	"time"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/pi"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
)

type Resource struct {
}

func (resource *Resource) DescribeVendorVerifyInfos(ctx context.Context, req *pb.DescribeVendorVerifyInfosRequest) ([]*models.Vendor, uint32, error) {
	var vendors []*models.Vendor
	offset := pbutil.GetOffsetFromRequest(req)
	limit := pbutil.GetLimitFromRequest(req)

	var vendorColumns = db.GetColumnsFromStruct(&models.Vendor{})
	query := pi.Global().DB(ctx).
		Select(vendorColumns...).
		From(constants.TableVendor).
		Offset(offset).
		Limit(limit).
		Where(manager.BuildFilterConditions(req, constants.TableVendor))

	query = manager.AddQueryOrderDir(query, req, constants.ColumnCreateTime)
	_, err := query.Load(&vendors)
	count, err := query.Count()
	return vendors, count, err
}

func (resource *Resource) GetVendorVerifyInfo(ctx context.Context, userID string) (*models.Vendor, error) {
	vendor := &models.Vendor{}
	vendorColumns := db.GetColumnsFromStruct(&models.Vendor{})
	err := pi.Global().DB(ctx).
		Select(vendorColumns...).
		From(constants.TableVendor).
		Where(db.Eq(constants.ColumnUserId, userID)).
		LoadOne(&vendor)
	if err != nil {
		return nil, err
	}
	return vendor, nil
}

func (resource *Resource) SubmitVendorVerifyInfo(ctx context.Context, vendor models.Vendor) (string, error) {
	var modifyType string
	var err error
	var userID string

	info, err := resource.GetVendorVerifyInfo(ctx, vendor.UserId)
	if info == nil {
		modifyType = "new"
	} else {
		modifyType = "update"
	}

	if modifyType == "new" {
		userID, err = resource.CreateVendorVerifyInfo(ctx, vendor)
		return userID, err
	}
	if modifyType == "update" {
		userID, err = resource.UpdateVendorVerifyInfo(ctx, vendor)
		return userID, err
	}
	return "", err
}

func (resource *Resource) UpdateVendorVerifyInfo(ctx context.Context, vendor models.Vendor) (string, error) {
	_, err := pi.Global().DB(ctx).
		Update(constants.TableVendor).
		Set("company_name", vendor.CompanyName).
		Set("company_website", vendor.CompanyWebsite).
		Set("company_profile", vendor.CompanyProfile).
		Set("authorizer_name", vendor.AuthorizerName).
		Set("authorizer_email", vendor.AuthorizerEmail).
		Set("authorizer_phone", vendor.AuthorizerPhone).
		Set("authorizer_position", vendor.AuthorizerPosition).
		Set("bank_name", vendor.BankName).
		Set("bank_account_name", vendor.BankAccountName).
		Set("bank_account_number", vendor.BankAccountNumber).
		Set("status", vendor.Status).
		Set("submit_time", vendor.SubmitTime).
		Set("status_time", vendor.StatusTime).
		Set("update_time", vendor.UpdateTime).
		Where(db.Eq(constants.ColumnUserId, vendor.UserId)).
		Exec()
	if err != nil {
		return "", err
	}
	return vendor.UserId, err
}

func (resource *Resource) CreateVendorVerifyInfo(ctx context.Context, vendor models.Vendor) (string, error) {
	_, err := pi.Global().DB(ctx).
		InsertInto(constants.TableVendor).
		Record(vendor).
		Exec()
	if err != nil {
		return "", err
	}
	return vendor.UserId, err
}

func (resource *Resource) PassVendorVerifyInfo(ctx context.Context, userID string) (string, error) {
	_, err := pi.Global().DB(ctx).
		Update(constants.TableVendor).
		Set("status", "passed").
		Set("status_time", time.Now()).
		Set("update_time", time.Now()).
		Where(db.Eq(constants.ColumnUserId, userID)).
		Exec()
	if err != nil {
		return "", err
	}
	return userID, err
}

func (resource *Resource) RejectVendorVerifyInfo(ctx context.Context, userID string, rejectmsg string) (string, error) {
	_, err := pi.Global().DB(ctx).
		Update(constants.TableVendor).
		Set("status", "rejected").
		Set("status_time", time.Now()).
		Set("status", "new").
		Set("update_time", time.Now()).
		Set("reject_message", rejectmsg).
		Where(db.Eq(constants.ColumnUserId, userID)).
		Exec()
	if err != nil {
		return "", err
	}
	return userID, err
}
