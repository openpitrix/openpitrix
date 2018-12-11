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
		UserId:      pbutil.ToProtoString("testuserID"),
		CompanyName: pbutil.ToProtoString("CompanyName-1"),
		//CompanyWebsite:    pbutil.ToProtoString("CompanyWebsite1"),
		//CompanyProfile:    pbutil.ToProtoString("CompanyProfile"),
		AuthorizerName:    pbutil.ToProtoString("AuthorizerName"),
		AuthorizerEmail:   pbutil.ToProtoString("AuthorizerEmail"),
		AuthorizerPhone:   pbutil.ToProtoString("AuthorizerPhone"),
		BankName:          pbutil.ToProtoString("BankName"),
		BankAccountName:   pbutil.ToProtoString("BankAccountName"),
		BankAccountNumber: pbutil.ToProtoString("BankAccountNumber"),
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
	userids = append(userids, "testuserID1")
	userids = append(userids, "testuserID2")

	var statuses []string
	statuses = append(statuses, "passed")

	var req = &vendor.DescribeVendorVerifyInfosRequest{
		SearchWord: pbutil.ToProtoString("AuthorizerName"),
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
		UserId: pbutil.ToProtoString("testuserID"),
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
		UserId: pbutil.ToProtoString("testuserID"),
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
		UserId:        pbutil.ToProtoString("testuserID"),
		RejectMessage: pbutil.ToProtoString("RejectMsg test"),
	}
	s.RejectVendorVerifyInfo(ctx, req)
}
