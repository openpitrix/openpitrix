package appvendor

import (
	"context"
	"time"

	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/gerr"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
)

type Handler struct {
}

func (h *Handler) DescribeVendorVerifyInfos(ctx context.Context, req *pb.DescribeVendorVerifyInfosRequest) (*pb.DescribeVendorVerifyInfosResponse, error) {
	resource := Resource{}
	vendors, count, err := resource.DescribeVendorVerifyInfos(ctx, req)
	if err != nil {
		logger.Error(ctx, "Failed to Describe VendorVerifyInfos [%+v]", req)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}

	var vendor models.AppVendor //need use a appvendor object to call function
	vendorPbSet := vendor.ParseVendorSet2pbSet(ctx, vendors)

	res := &pb.DescribeVendorVerifyInfosResponse{
		VendorVerifyInfoSet: vendorPbSet,
		TotalCount:          count,
	}
	return res, nil
}

func (h *Handler) GetVendorVerifyInfo(ctx context.Context, req *pb.GetVendorVerifyInfoRequest) (*pb.VendorVerifyInfo, error) {
	userID := req.GetUserId().GetValue()
	resource := Resource{}
	vendor, err := resource.GetVendorVerifyInfo(ctx, userID)

	if vendor == nil && err != nil {
		vendor = &models.AppVendor{}
		vendor.UserId = userID
		vendor.Status = "new"
		vendor.UpdateTime = time.Now()
		vendor.StatusTime = time.Now()
		userID, err = resource.CreateVendorVerifyInfo(ctx, *vendor)
		if err != nil {
			logger.Error(ctx, "appvendor does not exit,create new appvendor failed: [%+v]", vendor)
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorCreateResourceFailed)
		}
		logger.Debug(ctx, "appvendor does not exit,create new appvendor verify info: [%+v]", vendor)
	}

	if err != nil {
		logger.Error(ctx, "get VendorVerifyInfo failed: [%+v]", userID)
		return nil, gerr.NewWithDetail(ctx, gerr.NotFound, err, gerr.ErrorResourceNotFound, userID)
	}

	logger.Info(ctx, "Got VendorVerifyInfo: [%+v]", vendor)
	vendor2pb := vendor.ParseVendor2Pb(ctx, vendor)

	return vendor2pb, nil
}

func (h *Handler) PassVendorVerifyInfo(ctx context.Context, req *pb.PassVendorVerifyInfoRequest) (*pb.PassVendorVerifyInfoResponse, error) {
	userID := req.GetUserId().GetValue()
	resource := Resource{}
	userID, err := resource.PassVendorVerifyInfo(ctx, userID)
	if err != nil {
		logger.Error(ctx, "Pass VendorVerifyInfo failed: [%+v]", userID)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorUpdateResourceFailed)
	}
	res := &pb.PassVendorVerifyInfoResponse{
		UserId: pbutil.ToProtoString(userID),
	}
	return res, err
}

func (h *Handler) RejectVendorVerifyInfo(ctx context.Context, req *pb.RejectVendorVerifyInfoRequest) (*pb.RejectVendorVerifyInfoResponse, error) {
	userID := req.GetUserId().GetValue()
	rejectmsg := req.GetRejectMessage().GetValue()
	resource := Resource{}
	userID, err := resource.RejectVendorVerifyInfo(ctx, userID, rejectmsg)
	if err != nil {
		logger.Error(ctx, "Reject VendorVerifyInfo failed: [%+v]", userID)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorUpdateResourceFailed)
	}
	res := &pb.RejectVendorVerifyInfoResponse{
		UserId: pbutil.ToProtoString(userID),
	}
	return res, err
}

func (h *Handler) SubmitVendorVerifyInfo(ctx context.Context, req *pb.SubmitVendorVerifyInfoRequest) (*pb.SubmitVendorVerifyInfoResponse, error) {
	var modifyType string
	var err error
	var userID string
	userID = req.UserId.GetValue()

	resource := Resource{}
	info, err := resource.GetVendorVerifyInfo(ctx, userID)

	if info == nil && err != nil {
		modifyType = "new"
	} else if info != nil && err == nil {
		modifyType = "update"
	} else {
		logger.Error(ctx, "Submit VendorVerifyInfo failed: [%+v]", userID)
		return nil, err
	}

	if modifyType == "new" {
		vendor := &models.AppVendor{}
		vendor = vendor.ParseReq2Vendor(req)
		userID, err = resource.CreateVendorVerifyInfo(ctx, *vendor)
		if err != nil {
			logger.Error(ctx, "Submit VendorVerifyInfo failed: [%+v]", vendor)
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
		}
	}

	if modifyType == "update" {
		var appVendorColumns = db.GetColumnsFromStruct(&models.AppVendor{})
		attributes := manager.BuildUpdateAttributes(req, appVendorColumns...)
		delete(attributes, "user_id")
		logger.Debug(ctx, "UpdateAppVendor got attributes: [%+v]", attributes)

		userID, err = resource.UpdateVendorVerifyInfo(ctx, userID, attributes)
		if err != nil {
			logger.Error(ctx, "Submit VendorVerifyInfo failed: [%+v]", userID)
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorUpdateResourceFailed)
		}
	}
	res := &pb.SubmitVendorVerifyInfoResponse{
		UserId: pbutil.ToProtoString(userID),
	}
	return res, nil
}
