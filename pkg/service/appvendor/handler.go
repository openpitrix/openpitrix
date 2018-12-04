package appvendor

import (
	"context"

	"openpitrix.io/openpitrix/pkg/gerr"
	"openpitrix.io/openpitrix/pkg/logger"
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
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}
	var vendor models.Vendor //need use a vendor object to call function
	vendorPbSet, err := vendor.ParseVendorSet2pbSet(ctx, vendors)
	if err != nil {
		return nil, err
	}
	res := &pb.DescribeVendorVerifyInfosResponse{
		VendorVerifyInfoSet: vendorPbSet,
		TotalCount:          count,
	}
	return res, nil
}

func (h *Handler) GetVendorVerifyInfo(ctx context.Context, req *pb.GetVendorVerifyInfoRequest) (*pb.VendorVerifyInfo, error) {
	userID := req.GetUserId()
	resource := Resource{}
	vendor, err := resource.GetVendorVerifyInfo(ctx, userID)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.NotFound, err, gerr.ErrorResourceNotFound, userID)
	}
	logger.Debug(ctx, "Got vendor verify info: [%+v]", vendor)
	vendor2pb := vendor.ParseVendor2pb(ctx, vendor)
	return vendor2pb, nil
}

func (h *Handler) SubmitVendorVerifyInfo(ctx context.Context, req *pb.SubmitVendorVerifyInfoRequest) (*pb.SubmitVendorVerifyInfoResponse, error) {
	vendor := &models.Vendor{}
	vendor = vendor.ParseInputParams2Obj(req)

	resource := Resource{}
	userID, err := resource.SubmitVendorVerifyInfo(ctx, *vendor)

	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
	}

	res := &pb.SubmitVendorVerifyInfoResponse{
		UserId: pbutil.ToProtoString(userID),
	}
	return res, err
}

func (h *Handler) PassVendorVerifyInfo(ctx context.Context, req *pb.PassVendorVerifyInfoRequest) (*pb.PassVendorVerifyInfoResponse, error) {
	userID := req.GetUserId()
	resource := Resource{}
	userID, err := resource.PassVendorVerifyInfo(ctx, userID)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
	}
	res := &pb.PassVendorVerifyInfoResponse{
		UserId: pbutil.ToProtoString(userID),
	}
	return res, err
}

func (h *Handler) RejectVendorVerifyInfo(ctx context.Context, req *pb.RejectVendorVerifyInfoRequest) (*pb.RejectVendorVerifyInfoResponse, error) {
	userID := req.GetUserId()
	rejectmsg := req.GetRejectMsg().GetValue()
	resource := Resource{}
	userID, err := resource.RejectVendorVerifyInfo(ctx, userID, rejectmsg)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
	}
	res := &pb.RejectVendorVerifyInfoResponse{
		UserId: pbutil.ToProtoString(userID),
	}
	return res, err
}
