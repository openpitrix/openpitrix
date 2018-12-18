// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package appvendor

import (
	"context"
	"time"

	"openpitrix.io/openpitrix/pkg/manager"

	"openpitrix.io/openpitrix/pkg/gerr"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
)

func (s *Server) DescribeVendorVerifyInfos(ctx context.Context, req *pb.DescribeVendorVerifyInfosRequest) (*pb.DescribeVendorVerifyInfosResponse, error) {
	status := req.GetStatus()
	_, err := VerifyStatus(ctx, status...)
	if err != nil {
		return nil, err
	}

	vendors, count, err := DescribeVendorVerifyInfos(ctx, req)
	if err != nil {
		logger.Error(ctx, "Failed to describe vendorVerifyInfos, error: %+v", err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}

	var vendor models.AppVendor //need use a appvendor object to call function
	vendorPbSet := vendor.ParseVendorSet2PbSet(ctx, vendors)

	res := &pb.DescribeVendorVerifyInfosResponse{
		VendorVerifyInfoSet: vendorPbSet,
		TotalCount:          count,
	}
	return res, nil
}

func (s *Server) GetVendorVerifyInfo(ctx context.Context, req *pb.GetVendorVerifyInfoRequest) (*pb.VendorVerifyInfo, error) {
	userID := req.GetUserId().GetValue()
	vendor, err := GetVendorVerifyInfo(ctx, userID)

	if vendor == nil && err != nil {
		vendor = &models.AppVendor{}
		vendor.UserId = userID
		vendor.Status = models.StatusNew
		vendor.StatusTime = time.Now()
		userID, err = CreateVendorVerifyInfo(ctx, *vendor)
		if err != nil {
			logger.Error(ctx, "vendorVerifyInfo does not exit,create new vendorVerifyInfo failed [%s], %+v", userID, err)
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorCreateResourceFailed)
		}
		logger.Debug(ctx, "vendorVerifyInfo does not exit,create new vendorVerifyInfo verify info,[%+v]", vendor)
	}

	if err != nil {
		logger.Error(ctx, "Failed to get vendorVerifyInfo [%s], %+v", userID, err)
		return nil, gerr.NewWithDetail(ctx, gerr.NotFound, err, gerr.ErrorResourceNotFound, userID)
	}

	logger.Debug(ctx, "Got VendorVerifyInfo: [%+v]", vendor)
	vendor2pb := vendor.ParseVendor2Pb(ctx, vendor)

	return vendor2pb, nil
}

func (s *Server) SubmitVendorVerifyInfo(ctx context.Context, req *pb.SubmitVendorVerifyInfoRequest) (*pb.SubmitVendorVerifyInfoResponse, error) {
	err := s.validateSubmitParams(ctx, req)
	if err != nil {
		return nil, err
	}

	var userID string
	ifExist, err := s.checkVendorVerifyInfoIfExit(ctx, req.UserId)
	if err != nil {
		logger.Error(ctx, "Failed to get vendorVerifyInfo [%s], %+v", req.UserId, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
	}

	if ifExist {
		attributes := manager.BuildUpdateAttributes(req, models.ColumnCompanyName, models.ColumnCompanyWebsite, models.ColumnCompanyProfile,
			models.ColumnAuthorizerName, models.ColumnAuthorizerEmail, models.ColumnAuthorizerPhone, models.ColumnBankName, models.ColumnBankAccountName,
			models.ColumnBankAccountNumber)
		attributes[models.ColumnStatus] = models.StatusSubmitted
		attributes[models.ColumnSubmitTime] = time.Now()
		attributes[models.ColumnStatusTime] = time.Now()
		logger.Debug(ctx, "SubmitVendorVerifyInfo got attributes: [%+v]", attributes)
		userID, err = UpdateVendorVerifyInfo(ctx, req.UserId, attributes)
		if err != nil {
			logger.Error(ctx, "Failed to submit vendorVerifyInfo [%s], %+v", userID, err)
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorUpdateResourceFailed)
		}
	} else {
		vendor := &models.AppVendor{}
		vendor = vendor.ParseReq2Vendor(req)
		userID, err = CreateVendorVerifyInfo(ctx, *vendor)
		if err != nil {
			logger.Error(ctx, "Failed to submit vendorVerifyInfo [%+v], %+v", vendor, err)
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
		}
		logger.Debug(ctx, "vendorVerifyInfo does not exit,create new vendorVerifyInfo verify info,[%+v]", vendor)

	}
	res := &pb.SubmitVendorVerifyInfoResponse{
		UserId: pbutil.ToProtoString(userID),
	}
	return res, nil
}

func (s *Server) PassVendorVerifyInfo(ctx context.Context, req *pb.PassVendorVerifyInfoRequest) (*pb.PassVendorVerifyInfoResponse, error) {
	userID := req.GetUserId()
	userID, err := PassVendorVerifyInfo(ctx, userID)
	if err != nil {
		logger.Error(ctx, "Failed to pass vendorVerifyInfo [%s], %+v", userID, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorUpdateResourceFailed)
	}
	res := &pb.PassVendorVerifyInfoResponse{
		UserId: pbutil.ToProtoString(userID),
	}
	return res, err
}

func (s *Server) RejectVendorVerifyInfo(ctx context.Context, req *pb.RejectVendorVerifyInfoRequest) (*pb.RejectVendorVerifyInfoResponse, error) {
	userID := req.GetUserId()
	rejectmsg := req.GetRejectMessage().GetValue()
	userID, err := RejectVendorVerifyInfo(ctx, userID, rejectmsg)
	if err != nil {
		logger.Error(ctx, "Failed to reject vendorVerifyInfo [%s], %+v", userID, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorUpdateResourceFailed)
	}
	res := &pb.RejectVendorVerifyInfoResponse{
		UserId: pbutil.ToProtoString(userID),
	}
	return res, err
}

func (s *Server) checkVendorVerifyInfoIfExit(ctx context.Context, userID string) (bool, error) {
	info, err := GetVendorVerifyInfo(ctx, userID)
	if info == nil && err != nil {
		return false, nil
	} else if info != nil && err == nil {
		return true, nil
	} else {
		return false, err
	}
}

func (s *Server) validateSubmitParams(ctx context.Context, req *pb.SubmitVendorVerifyInfoRequest) error {
	isUrlFmt, err := VerifyUrl(ctx, req.CompanyWebsite.GetValue())
	if !isUrlFmt {
		logger.Error(ctx, "Failed to validateSubmitParams [%s].", req.CompanyWebsite.GetValue())
		return err
	}

	isEmailFmt, err := VerifyEmailFmt(ctx, req.AuthorizerEmail.GetValue())
	if !isEmailFmt {
		logger.Error(ctx, "Failed to validateSubmitParams [%s].", req.AuthorizerEmail.GetValue())
		return err
	}

	isPhoneFmt, err := VerifyPhoneFmt(ctx, req.AuthorizerPhone.GetValue())
	if !isPhoneFmt {
		logger.Error(ctx, "Failed to validateSubmitParams [%s].", req.AuthorizerPhone.GetValue())
		return err
	}

	isBankAccountNumberFmt, err := VerifyBankAccountNumberFmt(ctx, req.BankAccountNumber.GetValue())
	if !isBankAccountNumberFmt {
		logger.Error(ctx, "Failed to validateSubmitParams [%s].", req.BankAccountNumber.GetValue())
		return gerr.NewWithDetail(ctx, gerr.Internal, nil, gerr.ErrorValidateFailed)
	}
	return nil
}
