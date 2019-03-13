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

func ReqToVendorVerifyInfo(ctx context.Context, req *pb.SubmitVendorVerifyInfoRequest) *VendorVerifyInfo {
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
	Vendor.OwnerPath = ctxutil.GetSender(ctx).GetOwnerPath()
	t := time.Now()
	Vendor.SubmitTime = &t
	Vendor.StatusTime = time.Now()
	return &Vendor
}

func VendorVerifyInfoSetToPbSet(inVendors []*VendorVerifyInfo) []*pb.VendorVerifyInfo {
	var pbVendors []*pb.VendorVerifyInfo
	for _, inVendor := range inVendors {
		pbVendors = append(pbVendors, VendorVerifyInfoToPb(inVendor))
	}
	return pbVendors
}

func VendorVerifyInfoToPb(inVendor *VendorVerifyInfo) *pb.VendorVerifyInfo {
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
	pbVendor.Owner = pbutil.ToProtoString(inVendor.Owner)
	pbVendor.OwnerPath = inVendor.OwnerPath.ToProtoString()
	return &pbVendor
}

type VendorStatistics struct {
	UserId            string
	CompanyName       string
	ActiveAppCount    uint32
	ClusterCountMonth uint32
	ClusterCountTotal uint32
}

func VendorStatisticsToPb(vendorStatistics *VendorStatistics) *pb.VendorStatistics {
	pbVendor := pb.VendorStatistics{}
	pbVendor.UserId = pbutil.ToProtoString(vendorStatistics.UserId)
	pbVendor.CompanyName = pbutil.ToProtoString(vendorStatistics.CompanyName)
	pbVendor.ActiveAppCount = vendorStatistics.ActiveAppCount
	pbVendor.ClusterCountMonth = vendorStatistics.ClusterCountMonth
	pbVendor.ClusterCountTotal = vendorStatistics.ClusterCountTotal
	return &pbVendor
}

func VendorStatisticsSetToPbSet(vendorStatisticsSet []*VendorStatistics) []*pb.VendorStatistics {
	var pbVendorStatisticsSet []*pb.VendorStatistics
	for _, vendorStatistics := range vendorStatisticsSet {
		pbVendorStatisticsSet = append(pbVendorStatisticsSet, VendorStatisticsToPb(vendorStatistics))
	}
	return pbVendorStatisticsSet
}
