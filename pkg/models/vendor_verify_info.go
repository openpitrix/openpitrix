// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"context"
	"time"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/sender"
	"openpitrix.io/openpitrix/pkg/util/ctxutil"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
)

type VendorVerifyInfo struct {
	UserId            string
	CompanyName       string
	CompanyWebsite    string
	CompanyProfile    string
	AuthorizerName    string
	AuthorizerEmail   string
	AuthorizerPhone   string
	BankName          string
	BankAccountName   string
	BankAccountNumber string
	Status            string
	RejectMessage     string
	Approver          string
	Owner             string
	OwnerPath         sender.OwnerPath
	SubmitTime        *time.Time
	StatusTime        time.Time
}

func (vendor *VendorVerifyInfo) ParseReq2Vendor(ctx context.Context, req *pb.SubmitVendorVerifyInfoRequest) *VendorVerifyInfo {
	Vendor := VendorVerifyInfo{}
	Vendor.UserId = req.GetUserId()
	Vendor.CompanyName = req.GetCompanyName().GetValue()
	Vendor.CompanyWebsite = req.GetCompanyWebsite().GetValue()
	Vendor.CompanyProfile = req.GetCompanyProfile().GetValue()
	Vendor.AuthorizerName = req.GetAuthorizerName().GetValue()
	Vendor.AuthorizerEmail = req.GetAuthorizerEmail().GetValue()
	Vendor.AuthorizerPhone = req.GetAuthorizerPhone().GetValue()
	Vendor.BankName = req.GetBankName().GetValue()
	Vendor.BankAccountName = req.GetBankAccountName().GetValue()
	Vendor.BankAccountNumber = req.GetBankAccountNumber().GetValue()
	Vendor.Status = constants.StatusSubmitted
	Vendor.Owner = ctxutil.GetSender(ctx).UserId
	Vendor.OwnerPath = ctxutil.GetSender(ctx).OwnerPath
	t := time.Now()
	Vendor.SubmitTime = &t
	Vendor.StatusTime = time.Now()
	return &Vendor
}

func (vendor *VendorVerifyInfo) ParseVendorSet2PbSet(ctx context.Context, inVendors []*VendorVerifyInfo) []*pb.VendorVerifyInfo {
	var pbVendors []*pb.VendorVerifyInfo
	for _, inVendor := range inVendors {
		var pbVendor *pb.VendorVerifyInfo
		pbVendor = vendor.ParseVendor2Pb(ctx, inVendor)
		pbVendors = append(pbVendors, pbVendor)
	}
	return pbVendors
}

func (vendor *VendorVerifyInfo) ParseVendor2Pb(ctx context.Context, inVendor *VendorVerifyInfo) *pb.VendorVerifyInfo {
	pbVendor := pb.VendorVerifyInfo{}
	pbVendor.UserId = pbutil.ToProtoString(inVendor.UserId)
	pbVendor.CompanyName = pbutil.ToProtoString(inVendor.CompanyName)
	pbVendor.CompanyWebsite = pbutil.ToProtoString(inVendor.CompanyWebsite)
	pbVendor.CompanyProfile = pbutil.ToProtoString(inVendor.CompanyProfile)
	pbVendor.AuthorizerName = pbutil.ToProtoString(inVendor.AuthorizerName)
	pbVendor.AuthorizerEmail = pbutil.ToProtoString(inVendor.AuthorizerEmail)
	pbVendor.AuthorizerPhone = pbutil.ToProtoString(inVendor.AuthorizerPhone)
	pbVendor.BankName = pbutil.ToProtoString(inVendor.BankName)
	pbVendor.BankAccountName = pbutil.ToProtoString(inVendor.BankAccountName)
	pbVendor.BankAccountNumber = pbutil.ToProtoString(inVendor.BankAccountNumber)
	pbVendor.Status = pbutil.ToProtoString(inVendor.Status)
	pbVendor.RejectMessage = pbutil.ToProtoString(inVendor.RejectMessage)
	pbVendor.Approver = pbutil.ToProtoString(inVendor.Approver)
	if inVendor.SubmitTime != nil {
		pbVendor.SubmitTime = pbutil.ToProtoTimestamp(*inVendor.SubmitTime)
	}
	pbVendor.StatusTime = pbutil.ToProtoTimestamp(inVendor.StatusTime)
	return &pbVendor
}

type VendorStatistics struct {
	UserId            string
	CompanyName       string
	ActiveAppCount    int32
	ClusterCountMonth int32
	ClusterCountTotal int32
}

func (vendor *VendorStatistics) ParseVendorStatistics2Pb(ctx context.Context, inVendor *VendorStatistics) *pb.VendorStatistics {
	pbVendor := pb.VendorStatistics{}
	pbVendor.UserId = pbutil.ToProtoString(inVendor.UserId)
	pbVendor.CompanyName = pbutil.ToProtoString(inVendor.CompanyName)
	pbVendor.ActiveAppCount = pbutil.ToProtoInt32(inVendor.ActiveAppCount)
	pbVendor.ClusterCountMonth = pbutil.ToProtoInt32(inVendor.ClusterCountMonth)
	pbVendor.ClusterCountTotal = pbutil.ToProtoInt32(inVendor.ClusterCountTotal)
	return &pbVendor
}

func (vendor *VendorStatistics) ParseVendorStatisticsSet2PbSet(ctx context.Context, inVendors []*VendorStatistics) []*pb.VendorStatistics {
	var pbVendorStatisticses []*pb.VendorStatistics
	for _, inVendor := range inVendors {
		var pbVendor *pb.VendorStatistics
		pbVendor = vendor.ParseVendorStatistics2Pb(ctx, inVendor)
		pbVendorStatisticses = append(pbVendorStatisticses, pbVendor)
	}
	return pbVendorStatisticses
}
