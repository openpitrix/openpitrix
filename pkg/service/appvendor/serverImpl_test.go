package appvendor

import (
	"context"
	"testing"
	"time"

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
		UserId:             "testuserID33",
		CompanyName:        pbutil.ToProtoString("CompanyName33"),
		CompanyWebsite:     pbutil.ToProtoString("CompanyWebsite"),
		CompanyProfile:     pbutil.ToProtoString("CompanyProfile"),
		AuthorizerName:     pbutil.ToProtoString("AuthorizerName"),
		AuthorizerEmail:    pbutil.ToProtoString("AuthorizerEmail"),
		AuthorizerPhone:    pbutil.ToProtoString("AuthorizerPhone"),
		AuthorizerPosition: pbutil.ToProtoString("AuthorizerPosition"),
		BankName:           pbutil.ToProtoString("BankName"),
		BankAccountName:    pbutil.ToProtoString("BankAccountName"),
		BankAccountNumber:  pbutil.ToProtoString("BankAccountNumber"),
	}
	s.SubmitVendorVerifyInfo(ctx, req)
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
	userids = append(userids, "testuserID2")
	userids = append(userids, "testuserID3")

	var statuses []string
	statuses = append(statuses, "pending")

	var req = &vendor.DescribeVendorVerifyInfosRequest{
		SearchWord: pbutil.ToProtoString("CompanyName"),
		SortKey:    pbutil.ToProtoString("user_id"),
		Reverse:    pbutil.ToProtoBool(false),
		Limit:      10,
		Offset:     0,
		UserId:     userids,
		Status:     statuses,
	}
	s.DescribeVendorVerifyInfos(ctx, req)
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
		UserId: "testuserID33",
	}
	s.GetVendorVerifyInfo(ctx, req)
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
	s.PassVendorVerifyInfo(ctx, req)
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
		UserId:    "testuserID",
		RejectMsg: pbutil.ToProtoString("RejectMsg test"),
	}
	s.RejectVendorVerifyInfo(ctx, req)
}
