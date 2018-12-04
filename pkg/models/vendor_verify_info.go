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

type Vendor struct {
	UserId             string
	CompanyName        string
	CompanyWebsite     string
	CompanyProfile     string
	AuthorizerName     string
	AuthorizerEmail    string
	AuthorizerPosition string
	AuthorizerPhone    string
	BankName           string
	BankAccountName    string
	BankAccountNumber  string
	Status             string
	RejectMessage      string
	SubmitTime         time.Time
	StatusTime         time.Time
	UpdateTime         time.Time
}

func (vendor *Vendor) ParseInputParams2Obj(in *pb.SubmitVendorVerifyInfoRequest) *Vendor {
	Vendor := Vendor{}
	Vendor.UserId = in.GetUserId()
	Vendor.CompanyName = in.GetCompanyName().GetValue()
	Vendor.CompanyWebsite = in.GetCompanyWebsite().GetValue()
	Vendor.CompanyProfile = in.GetCompanyProfile().GetValue()
	Vendor.AuthorizerName = in.GetAuthorizerName().GetValue()
	Vendor.AuthorizerEmail = in.GetAuthorizerEmail().GetValue()
	Vendor.AuthorizerPosition = in.GetAuthorizerPosition().GetValue()
	Vendor.AuthorizerPhone = in.GetAuthorizerPhone().GetValue()
	Vendor.BankName = in.GetBankName().GetValue()
	Vendor.BankAccountName = in.GetBankAccountName().GetValue()
	Vendor.BankAccountNumber = in.GetBankAccountNumber().GetValue()
	Vendor.Status = "pending"
	Vendor.SubmitTime = time.Now()
	Vendor.UpdateTime = time.Now()
	Vendor.StatusTime = time.Now()
	return &Vendor
}

func (vendor *Vendor) ParseVendorSet2pbSet(ctx context.Context, invendors []*Vendor) ([]*pb.VendorVerifyInfo, error) {
	var pbVendors []*pb.VendorVerifyInfo
	for _, invendor := range invendors {
		var pbVendor *pb.VendorVerifyInfo
		pbVendor = vendor.ParseVendor2pb(ctx, invendor)
		pbVendors = append(pbVendors, pbVendor)
	}
	return pbVendors, nil
}

func (vendor *Vendor) ParseVendor2pb(ctx context.Context, invendor *Vendor) *pb.VendorVerifyInfo {
	logger.Info(nil, "test="+invendor.UserId)
	pbVendor := pb.VendorVerifyInfo{}
	pbVendor.UserId = pbutil.ToProtoString(invendor.UserId)
	pbVendor.CompanyName = pbutil.ToProtoString(invendor.CompanyName)
	pbVendor.CompanyWebsite = pbutil.ToProtoString(invendor.CompanyWebsite)
	pbVendor.CompanyProfile = pbutil.ToProtoString(invendor.CompanyProfile)
	pbVendor.AuthorizerName = pbutil.ToProtoString(invendor.AuthorizerName)
	pbVendor.AuthorizerEmail = pbutil.ToProtoString(invendor.AuthorizerEmail)
	pbVendor.AuthorizerPosition = pbutil.ToProtoString(invendor.AuthorizerPosition)
	pbVendor.AuthorizerPhone = pbutil.ToProtoString(invendor.AuthorizerPhone)
	pbVendor.BankName = pbutil.ToProtoString(invendor.BankName)
	pbVendor.BankAccountName = pbutil.ToProtoString(invendor.BankAccountName)
	pbVendor.BankAccountNumber = pbutil.ToProtoString(invendor.BankAccountNumber)
	pbVendor.Status = pbutil.ToProtoString(invendor.Status)
	pbVendor.RejectMessage = pbutil.ToProtoString(invendor.RejectMessage)
	pbVendor.SubmitTime = pbutil.ToProtoTimestamp(invendor.SubmitTime)
	pbVendor.UpdateTime = pbutil.ToProtoTimestamp(invendor.UpdateTime)
	pbVendor.StatusTime = pbutil.ToProtoTimestamp(invendor.StatusTime)
	return &pbVendor

}
