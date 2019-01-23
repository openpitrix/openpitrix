// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// +build integration

package appvendor

import (
	"testing"

	"github.com/stretchr/testify/require"

	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/test/client/app_vendor_manager"
	"openpitrix.io/openpitrix/test/models"
	"openpitrix.io/openpitrix/test/testutil"
)

var clientConfig = testutil.GetClientConfig()

func TestAppVendor(t *testing.T) {
	client := testutil.GetClient(clientConfig)
	testUserID := "appvendor_test_userID"

	// SubmitVendorVerifyInfo
	logger.Info(nil, "Test1 SubmitVendorVerifyInfo")
	submitParams := app_vendor_manager.NewSubmitVendorVerifyInfoParams()
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
	submitResp, err := client.AppVendorManager.SubmitVendorVerifyInfo(submitParams, nil)
	require.NoError(t, err)
	userId := submitResp.Payload.UserID
	if userId != testUserID {
		t.Fatalf("failed to SubmitVendorVerifyInfo:UserID= [%s]", testUserID)
	}

	// DescribeVendorVerifyInfos
	logger.Info(nil, "%s", "Test2 DescribeVendorVerifyInfos**************************")
	describeParams := app_vendor_manager.NewDescribeVendorVerifyInfosParams()

	var userids []string
	userids = append(userids, testUserID)

	describeParams.SetUserID(userids)

	logger.Info(nil, "test describeParams=[%+v]", describeParams)
	describeResp, err := client.AppVendorManager.DescribeVendorVerifyInfos(describeParams, nil)
	require.NoError(t, err)
	AppVendors := describeResp.Payload.VendorVerifyInfoSet
	logger.Info(nil, "test describeParams result,AppVendors=[%+v]", AppVendors)
	respUserId := AppVendors[0].UserID
	logger.Info(nil, "test describeParams result,respUserId=[%+s]", respUserId)
	if userId != testUserID {
		t.Fatalf("failed to SubmitVendorVerifyInfo:UserID= [%s]", testUserID)
	}

	// PassVendorVerifyInfo
	passParams := app_vendor_manager.NewPassVendorVerifyInfoParams()
	passParams.SetBody(
		&models.OpenpitrixPassVendorVerifyInfoRequest{
			UserID: testUserID,
		})
	passResp, err := client.AppVendorManager.PassVendorVerifyInfo(passParams, nil)
	require.NoError(t, err)
	t.Log(passResp)

	describeParams1 := app_vendor_manager.NewDescribeVendorVerifyInfosParams()
	var userids1 []string
	userids1 = append(userids1, testUserID)
	var statuses1 []string
	statuses1 = append(statuses1, "passed")

	describeParams1.SetUserID(userids1)
	describeParams1.SetStatus(statuses1)

	describeResp1, err := client.AppVendorManager.DescribeVendorVerifyInfos(describeParams1, nil)
	require.NoError(t, err)
	AppVendors1 := describeResp1.Payload.VendorVerifyInfoSet
	respAppvendor1 := AppVendors1[0]
	logger.Info(nil, "test describeParams result,respAppvendor=[%+v]", respAppvendor1)

	if respAppvendor1.Status != "passed" {
		t.Fatalf("failed to SubmitVendorVerifyInfo:UserID= [%s]", testUserID)
	}

	// RejectVendorVerifyInfo
	logger.Info(nil, "%s", "Test5 GetVendorVerifyInfo**************************")
	rejectParams := app_vendor_manager.NewRejectVendorVerifyInfoParams()

	rejectParams.SetBody(&models.OpenpitrixRejectVendorVerifyInfoRequest{
		UserID:        testUserID,
		RejectMessage: "RejectMsg test",
	})
	rejectResp, err := client.AppVendorManager.RejectVendorVerifyInfo(rejectParams, nil)
	require.NoError(t, err)
	t.Log(rejectResp)

	describeParams2 := app_vendor_manager.NewDescribeVendorVerifyInfosParams()
	var userids2 []string
	userids2 = append(userids2, testUserID)
	var statuses2 []string
	statuses2 = append(statuses2, "rejected")

	describeParams2.SetUserID(userids2)
	describeParams2.SetStatus(statuses2)

	describeResp2, err := client.AppVendorManager.DescribeVendorVerifyInfos(describeParams2, nil)
	require.NoError(t, err)
	AppVendors2 := describeResp2.Payload.VendorVerifyInfoSet
	respAppvendor2 := AppVendors2[0]
	logger.Info(nil, "test describeParams result,respAppvendor=[%+v]", respAppvendor2)

	if respAppvendor2.Status != "rejected" {
		t.Fatalf("failed to SubmitVendorVerifyInfo:UserID= [%s]", testUserID)
	}

	// DescribeAppVendorStatistics
	logger.Info(nil, "%s", "Test6 DescribeAppVendorStatistics**************************")
	describeStaParams := app_vendor_manager.NewDescribeAppVendorStatisticsParams()

	var userids3 []string
	userids3 = append(userids3, testUserID)
	var statuses3 []string
	statuses3 = append(statuses3, "rejected")

	describeStaParams.SetUserID(userids3)
	describeStaParams.SetStatus(statuses3)

	describeStaResp, err := client.AppVendorManager.DescribeAppVendorStatistics(describeStaParams, nil)
	require.NoError(t, err)
	vendorStatisticsVendors := describeStaResp.Payload.VendorVerifyStatisticsSet
	if len(vendorStatisticsVendors) == 0 {
		t.Fatalf("failed to DescribeAppVendorStatistics with params [%+v]", describeStaParams)
	} else {
		t.Logf("success to DescribeAppVendorStatistics:vendorStatisticsVendors cnt= [%d]", len(vendorStatisticsVendors))
	}

	t.Log("test AppVendor finish, all test is ok")

}
