// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// +build integration

package isv

import (
	"testing"

	"openpitrix.io/openpitrix/pkg/logger"
	isv_manager "openpitrix.io/openpitrix/test/client/isv_manager"
	"openpitrix.io/openpitrix/test/models"
	"openpitrix.io/openpitrix/test/testutil"
)

const Service = "hyperpitrix"

var clientConfig = testutil.GetClientConfig()

// temporary comment
func testIsv(t *testing.T) {
	client := testutil.GetClient(clientConfig)
	testUserID := "appvendor_test_userID"

	// SubmitVendorVerifyInfo
	logger.Info(nil, "Test1 SubmitVendorVerifyInfo")
	submitParams := isv_manager.NewSubmitVendorVerifyInfoParams()
	submitParams.SetBody(
		&models.OpenpitrixSubmitVendorVerifyInfoRequest{
			AuthorizerEmail:   "testemail@163.com",
			AuthorizerName:    "AuthorizerName",
			AuthorizerPhone:   "15827656666",
			BankAccountName:   "BankAccountName",
			BankAccountNumber: "6226820011200783",
			BankName:          "BankName",
			CompanyName:       "CompanyName",
			CompanyProfile:    "CompanyProfile",
			CompanyWebsite:    "www.baidu.com",
			UserID:            testUserID,
		})
	submitResp, err := client.IsvManager.SubmitVendorVerifyInfo(submitParams, nil)
	testutil.NoError(t, err, []string{Service})
	userId := submitResp.Payload.UserID
	if userId != testUserID {
		t.Fatalf("failed to SubmitVendorVerifyInfo:UserID= [%s]", testUserID)
	}

	// DescribeVendorVerifyInfos
	logger.Info(nil, "%s", "Test2 DescribeVendorVerifyInfos**************************")
	describeParams := isv_manager.NewDescribeVendorVerifyInfosParams()

	var userids []string
	userids = append(userids, testUserID)

	describeParams.SetUserID(userids)

	logger.Info(nil, "test describeParams=[%+v]", describeParams)
	describeResp, err := client.IsvManager.DescribeVendorVerifyInfos(describeParams, nil)
	testutil.NoError(t, err, []string{Service})
	AppVendors := describeResp.Payload.VendorVerifyInfoSet
	logger.Info(nil, "test describeParams result,AppVendors=[%+v]", AppVendors)
	respUserId := AppVendors[0].UserID
	logger.Info(nil, "test describeParams result,respUserId=[%+s]", respUserId)
	if userId != testUserID {
		t.Fatalf("failed to SubmitVendorVerifyInfo:UserID= [%s]", testUserID)
	}

	// PassVendorVerifyInfo
	passParams := isv_manager.NewPassVendorVerifyInfoParams()
	passParams.SetBody(
		&models.OpenpitrixPassVendorVerifyInfoRequest{
			UserID: testUserID,
		})
	passResp, err := client.IsvManager.PassVendorVerifyInfo(passParams, nil)
	testutil.NoError(t, err, []string{Service, "openpitrix-account-service", "openpitrix-im-service", "openpitrix-am-service"})
	t.Log(passResp)

	describeParams1 := isv_manager.NewDescribeVendorVerifyInfosParams()
	var userids1 []string
	userids1 = append(userids1, testUserID)
	var statuses1 []string
	statuses1 = append(statuses1, "passed")

	describeParams1.SetUserID(userids1)
	describeParams1.SetStatus(statuses1)

	describeResp1, err := client.IsvManager.DescribeVendorVerifyInfos(describeParams1, nil)
	testutil.NoError(t, err, []string{Service})
	AppVendors1 := describeResp1.Payload.VendorVerifyInfoSet
	respAppvendor1 := AppVendors1[0]
	logger.Info(nil, "test describeParams result,respAppvendor=[%+v]", respAppvendor1)

	if respAppvendor1.Status != "passed" {
		t.Fatalf("failed to SubmitVendorVerifyInfo:UserID= [%s]", testUserID)
	}

	// RejectVendorVerifyInfo
	logger.Info(nil, "%s", "Test5 GetVendorVerifyInfo**************************")
	rejectParams := isv_manager.NewRejectVendorVerifyInfoParams()

	rejectParams.SetBody(&models.OpenpitrixRejectVendorVerifyInfoRequest{
		UserID:        testUserID,
		RejectMessage: "RejectMsg test",
	})
	rejectResp, err := client.IsvManager.RejectVendorVerifyInfo(rejectParams, nil)
	testutil.NoError(t, err, []string{Service})
	t.Log(rejectResp)

	describeParams2 := isv_manager.NewDescribeVendorVerifyInfosParams()
	var userids2 []string
	userids2 = append(userids2, testUserID)
	var statuses2 []string
	statuses2 = append(statuses2, "rejected")

	describeParams2.SetUserID(userids2)
	describeParams2.SetStatus(statuses2)

	describeResp2, err := client.IsvManager.DescribeVendorVerifyInfos(describeParams2, nil)
	testutil.NoError(t, err, []string{Service})
	AppVendors2 := describeResp2.Payload.VendorVerifyInfoSet
	respAppvendor2 := AppVendors2[0]
	logger.Info(nil, "test describeParams result,respAppvendor=[%+v]", respAppvendor2)

	if respAppvendor2.Status != "rejected" {
		t.Fatalf("failed to SubmitVendorVerifyInfo:UserID= [%s]", testUserID)
	}

	// DescribeAppVendorStatistics
	logger.Info(nil, "%s", "Test6 DescribeAppVendorStatistics**************************")
	describeStaParams := isv_manager.NewDescribeAppVendorStatisticsParams()

	var userids3 []string
	userids3 = append(userids3, testUserID)
	var statuses3 []string
	statuses3 = append(statuses3, "rejected")

	describeStaParams.SetUserID(userids3)
	describeStaParams.SetStatus(statuses3)

	describeStaResp, err := client.IsvManager.DescribeAppVendorStatistics(describeStaParams, nil)
	testutil.NoError(t, err, []string{Service})
	vendorStatisticsVendors := describeStaResp.Payload.VendorVerifyStatisticsSet
	if len(vendorStatisticsVendors) == 0 {
		t.Fatalf("failed to DescribeAppVendorStatistics with params [%+v]", describeStaParams)
	} else {
		t.Logf("success to DescribeAppVendorStatistics:vendorStatisticsVendors cnt= [%d]", len(vendorStatisticsVendors))
	}

	t.Log("test AppVendor finish, all test is ok")

}
