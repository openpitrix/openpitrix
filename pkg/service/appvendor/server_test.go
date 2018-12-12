// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package appvendor

import (
	"context"
	"testing"
	"time"

	"openpitrix.io/openpitrix/pkg/logger"

	vendor "openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
)

func TestNewServer(t *testing.T) {
	if !*tTestingEnvEnabled {
		t.Skip("testing env disabled")
	}
	InitGlobelSetting()
}

func TestSubmitVendorVerifyInfo(t *testing.T) {
	if !*tTestingEnvEnabled {
		t.Skip("testing env disabled")
	}
	InitGlobelSetting()
	s, _ := NewServer()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var req = &vendor.SubmitVendorVerifyInfoRequest{
		UserId:            pbutil.ToProtoString("testuserID"),
		CompanyName:       pbutil.ToProtoString("CompanyName1"),
		CompanyWebsite:    pbutil.ToProtoString("CompanyWebsite1"),
		CompanyProfile:    pbutil.ToProtoString("CompanyProfile"),
		AuthorizerName:    pbutil.ToProtoString("AuthorizerName"),
		AuthorizerEmail:   pbutil.ToProtoString("AuthorizerEmail"),
		AuthorizerPhone:   pbutil.ToProtoString("AuthorizerPhone"),
		BankName:          pbutil.ToProtoString("BankName"),
		BankAccountName:   pbutil.ToProtoString("BankAccountName"),
		BankAccountNumber: pbutil.ToProtoString("BankAccountNumber"),
	}
	resp, _ := s.SubmitVendorVerifyInfo(ctx, req)
	logger.Info(nil, "TestSubmitVendorVerifyInfo %s", resp.GetUserId())

}

func TestDescribeVendorVerifyInfos(t *testing.T) {
	if !*tTestingEnvEnabled {
		t.Skip("testing env disabled")
	}
	InitGlobelSetting()
	s, _ := NewServer()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var userids []string
	userids = append(userids, "testuserID")
	userids = append(userids, "testuserID1")
	userids = append(userids, "testuserID2")

	var statuses []string
	statuses = append(statuses, "new")
	statuses = append(statuses, "submitted")
	statuses = append(statuses, "passed")
	statuses = append(statuses, "rejected")

	var req = &vendor.DescribeVendorVerifyInfosRequest{
		SearchWord: pbutil.ToProtoString("AuthorizerName"),
		SortKey:    pbutil.ToProtoString("user_id"),
		Reverse:    pbutil.ToProtoBool(false),
		Limit:      10,
		Offset:     0,
		UserId:     userids,
		Status:     statuses,
	}
	resp, _ := s.DescribeVendorVerifyInfos(ctx, req)
	logger.Info(nil, "TestDescribeVendorVerifyInfos TotalCount=%d", resp.GetTotalCount())
}

func TestGetVendorVerifyInfo(t *testing.T) {
	if !*tTestingEnvEnabled {
		t.Skip("testing env disabled")
	}
	InitGlobelSetting()
	s, _ := NewServer()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	var req = &vendor.GetVendorVerifyInfoRequest{
		UserId: pbutil.ToProtoString("testuserID"),
	}
	resp, _ := s.GetVendorVerifyInfo(ctx, req)
	logger.Info(nil, "TestGetVendorVerifyInfo %s", resp.GetUserId())

}

func TestPassVendorVerifyInfo(t *testing.T) {
	if !*tTestingEnvEnabled {
		t.Skip("testing env disabled")
	}
	InitGlobelSetting()
	s, _ := NewServer()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	var req = &vendor.PassVendorVerifyInfoRequest{
		UserId: pbutil.ToProtoString("testuserID"),
	}
	resp, _ := s.PassVendorVerifyInfo(ctx, req)
	logger.Info(nil, "TestPassVendorVerifyInfo %s ", resp.GetUserId())
}

func TestRejectVendorVerifyInfo(t *testing.T) {
	if !*tTestingEnvEnabled {
		t.Skip("testing env disabled")
	}
	InitGlobelSetting()
	s, _ := NewServer()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	var req = &vendor.RejectVendorVerifyInfoRequest{
		UserId:        pbutil.ToProtoString("testuserID"),
		RejectMessage: pbutil.ToProtoString("RejectMsg test"),
	}
	resp, _ := s.RejectVendorVerifyInfo(ctx, req)
	logger.Info(nil, "TestRejectVendorVerifyInfo %s", resp.GetUserId())
}
