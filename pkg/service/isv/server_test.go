// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package isv

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
		UserId:            "testuserID",
		CompanyName:       pbutil.ToProtoString("CompanyName23"),
		CompanyWebsite:    pbutil.ToProtoString("www.baidu.com"),
		CompanyProfile:    pbutil.ToProtoString("CompanyProfile"),
		AuthorizerName:    pbutil.ToProtoString("AuthorizerName2"),
		AuthorizerEmail:   pbutil.ToProtoString("testemail@163.com"),
		AuthorizerPhone:   pbutil.ToProtoString("15827656666"),
		BankName:          pbutil.ToProtoString("BankName"),
		BankAccountName:   pbutil.ToProtoString("BankAccountName"),
		BankAccountNumber: pbutil.ToProtoString("6226820011200783"),
	}
	resp, _ := s.SubmitVendorVerifyInfo(ctx, req)
	logger.Info(nil, "Test Passed,TestSubmitVendorVerifyInfo %s", resp.GetUserId())

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
	//statuses = append(statuses, "submitted")
	//statuses = append(statuses, "passed")
	//statuses = append(statuses, "rejected")

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
	logger.Info(nil, "Test Passed,TestDescribeVendorVerifyInfos TotalCount=%d", resp.GetTotalCount())
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
		UserId: "testuserID",
	}
	resp, _ := s.PassVendorVerifyInfo(ctx, req)
	logger.Info(nil, "Test Passed,TestPassVendorVerifyInfo %s ", resp.GetUserId())
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
		UserId:        "testuserID",
		RejectMessage: pbutil.ToProtoString("RejectMsg test"),
	}
	resp, _ := s.RejectVendorVerifyInfo(ctx, req)
	logger.Info(nil, "Test Passed,TestRejectVendorVerifyInfo %s", resp.GetUserId())
}
