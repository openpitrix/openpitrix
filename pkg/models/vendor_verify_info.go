// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"context"
	"time"

	"openpitrix.io/openpitrix/pkg/logger"

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

func (vendor *AppVendor) ParseInputParams2Obj(req *pb.SubmitVendorVerifyInfoRequest) *AppVendor {
	Vendor := AppVendor{}

	//Vendor.UserId = req.GetUserId()
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
	Vendor.UpdateTime = time.Now()
	Vendor.StatusTime = time.Now()
	return &Vendor
}

func (vendor *AppVendor) ParseVendorSet2pbSet(ctx context.Context, invendors []*AppVendor) []*pb.VendorVerifyInfo {
	var pbVendors []*pb.VendorVerifyInfo
	for _, invendor := range invendors {
		var pbVendor *pb.VendorVerifyInfo
		pbVendor = vendor.ParseVendor2pb(ctx, invendor)
		pbVendors = append(pbVendors, pbVendor)
	}
	return pbVendors
}

func (vendor *AppVendor) ParseVendor2pb(ctx context.Context, invendor *AppVendor) *pb.VendorVerifyInfo {
	pbVendor := pb.VendorVerifyInfo{}
	logger.Info(nil, "test=[%+v]", invendor.UserId)
	pbVendor.UserId = pbutil.ToProtoString(invendor.UserId)
	pbVendor.CompanyName = pbutil.ToProtoString(invendor.CompanyName)
	pbVendor.CompanyWebsite = pbutil.ToProtoString(invendor.CompanyWebsite)
	pbVendor.CompanyProfile = pbutil.ToProtoString(invendor.CompanyProfile)
	pbVendor.AuthorizerName = pbutil.ToProtoString(invendor.AuthorizerName)
	pbVendor.AuthorizerEmail = pbutil.ToProtoString(invendor.AuthorizerEmail)
	pbVendor.AuthorizerPhone = pbutil.ToProtoString(invendor.AuthorizerPhone)
	pbVendor.BankName = pbutil.ToProtoString(invendor.BankName)
	pbVendor.BankAccountName = pbutil.ToProtoString(invendor.BankAccountName)
	pbVendor.BankAccountNumber = pbutil.ToProtoString(invendor.BankAccountNumber)
	pbVendor.Status = pbutil.ToProtoString(invendor.Status)
	pbVendor.RejectMessage = pbutil.ToProtoString(invendor.RejectMessage)
	if invendor.SubmitTime != nil {
		pbVendor.SubmitTime = pbutil.ToProtoTimestamp(*invendor.SubmitTime)
	}
	pbVendor.UpdateTime = pbutil.ToProtoTimestamp(invendor.UpdateTime)
	pbVendor.StatusTime = pbutil.ToProtoTimestamp(invendor.StatusTime)
	return &pbVendor

}
