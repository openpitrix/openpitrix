// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"context"
	"time"

	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
)

type AppVendor struct {
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
	SubmitTime        *time.Time
	StatusTime        time.Time
	UpdateTime        time.Time
}

func (vendor *AppVendor) ParseReq2Vendor(req *pb.SubmitVendorVerifyInfoRequest) *AppVendor {
	Vendor := AppVendor{}
	Vendor.UserId = req.GetUserId().GetValue()
	Vendor.CompanyName = req.GetCompanyName().GetValue()
	Vendor.CompanyWebsite = req.GetCompanyWebsite().GetValue()
	Vendor.CompanyProfile = req.GetCompanyProfile().GetValue()
	Vendor.AuthorizerName = req.GetAuthorizerName().GetValue()
	Vendor.AuthorizerEmail = req.GetAuthorizerEmail().GetValue()
	Vendor.AuthorizerPhone = req.GetAuthorizerPhone().GetValue()
	Vendor.BankName = req.GetBankName().GetValue()
	Vendor.BankAccountName = req.GetBankAccountName().GetValue()
	Vendor.BankAccountNumber = req.GetBankAccountNumber().GetValue()
	Vendor.Status = "submitted"
	t := time.Now()
	Vendor.SubmitTime = &t
	Vendor.StatusTime = time.Now()
	return &Vendor
}

func (vendor *AppVendor) ParseVendorSet2pbSet(ctx context.Context, inVendors []*AppVendor) []*pb.VendorVerifyInfo {
	var pbVendors []*pb.VendorVerifyInfo
	for _, invendor := range inVendors {
		var pbVendor *pb.VendorVerifyInfo
		pbVendor = vendor.ParseVendor2Pb(ctx, invendor)
		pbVendors = append(pbVendors, pbVendor)
	}
	return pbVendors
}

func (vendor *AppVendor) ParseVendor2Pb(ctx context.Context, inVendor *AppVendor) *pb.VendorVerifyInfo {
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
	if inVendor.SubmitTime != nil {
		pbVendor.SubmitTime = pbutil.ToProtoTimestamp(*inVendor.SubmitTime)
	}
	pbVendor.StatusTime = pbutil.ToProtoTimestamp(inVendor.StatusTime)
	return &pbVendor

}
