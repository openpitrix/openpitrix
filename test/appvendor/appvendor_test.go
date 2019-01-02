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
	/*===============================================================================================*/
	logger.Info(nil, "%s", "/*=================================================*/")
	logger.Info(nil, "%s", "Test1 SubmitVendorVerifyInfo**************************")
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

	submitParams.SetUserID(testUserID)

	_, err := client.AppVendorManager.SubmitVendorVerifyInfo(submitParams, nil)
	require.NoError(t, err)

	submitResp, err := client.AppVendorManager.SubmitVendorVerifyInfo(submitParams, nil)
	require.NoError(t, err)
	userId := submitResp.Payload.UserID
	if userId != testUserID {
		t.Fatalf("failed to SubmitVendorVerifyInfo:UserID= [%s]", testUserID)
	}

	// DescribeVendorVerifyInfos
	/*===============================================================================================*/
	logger.Info(nil, "%s", "/*=================================================*/")
	logger.Info(nil, "%s", "Test2 DescribeVendorVerifyInfos**************************")
	describeParams := app_vendor_manager.NewDescribeVendorVerifyInfosParams()
	logger.Info(nil, "test describeParams=[%+v]", describeParams)
	describeParams.WithUserID([]string{testUserID})
	describeResp, err := client.AppVendorManager.DescribeVendorVerifyInfos(describeParams, nil)
	require.NoError(t, err)
	AppVendors := describeResp.Payload.VendorVerifyInfoSet
	if len(AppVendors) == 0 {
		t.Fatalf("failed to DescribeVendorVerifyInfos with params [%+v]", describeParams)
	}
	if AppVendors[0].UserID != testUserID || AppVendors[0].CompanyName != "CompanyName" {
		t.Fatalf("failed to SubmitVendorVerifyInfo with params [%+v]", submitParams)
	}

	// GetVendorVerifyInfo
	/*===============================================================================================*/
	logger.Info(nil, "%s", "/*=================================================*/")
	logger.Info(nil, "%s", "Test3 GetVendorVerifyInfo**************************")
	getParams := app_vendor_manager.NewGetVendorVerifyInfoParams()
	getParams.SetUserID(&testUserID)
	getResp, err := client.AppVendorManager.GetVendorVerifyInfo(getParams, nil)

	require.NoError(t, err)
	t.Log(getResp)
	if getResp.Payload.UserID != testUserID {
		t.Fatalf("failed to GetVendorVerifyInfo:UserID= [%s]", testUserID)
	}

	// PassVendorVerifyInfo
	/*===============================================================================================*/
	logger.Info(nil, "%s", "/*=================================================*/")
	logger.Info(nil, "%s", "Test4 GetVendorVerifyInfo**************************")
	passParams := app_vendor_manager.NewPassVendorVerifyInfoParams()
	passParams.SetBody(
		&models.OpenpitrixPassVendorVerifyInfoRequest{
			UserID: testUserID,
		})
	passResp, err := client.AppVendorManager.PassVendorVerifyInfo(passParams, nil)
	require.NoError(t, err)
	t.Log(passResp)
	getParams1 := app_vendor_manager.NewGetVendorVerifyInfoParams()
	getParams1.SetUserID(&testUserID)
	getResp1, err := client.AppVendorManager.GetVendorVerifyInfo(getParams1, nil)
	require.NoError(t, err)
	t.Log(getResp1)
	if getResp1.Payload.UserID == testUserID && getResp1.Payload.Status == "passed" {
		t.Logf("success to PassVendorVerifyInfo:UserID= [%s]", testUserID)
	} else {
		t.Fatalf("failed to PassVendorVerifyInfo:UserID= [%s]", testUserID)
	}

	// RejectVendorVerifyInfo
	/*===============================================================================================*/
	logger.Info(nil, "%s", "/*=================================================*/")
	logger.Info(nil, "%s", "Test5 GetVendorVerifyInfo**************************")
	rejectParams := app_vendor_manager.NewRejectVendorVerifyInfoParams()

	rejectParams.SetBody(&models.OpenpitrixRejectVendorVerifyInfoRequest{
		UserID:        testUserID,
		RejectMessage: "RejectMsg test",
	})
	rejectResp, err := client.AppVendorManager.RejectVendorVerifyInfo(rejectParams, nil)
	require.NoError(t, err)
	t.Log(rejectResp)
	getParams2 := app_vendor_manager.NewGetVendorVerifyInfoParams()
	getParams2.SetUserID(&testUserID)
	getResp2, err := client.AppVendorManager.GetVendorVerifyInfo(getParams2, nil)
	require.NoError(t, err)
	t.Log(getResp2)
	if getResp2.Payload.UserID == testUserID && getResp2.Payload.Status == "rejected" {
		t.Logf("success to RejectVendorVerifyInfo:UserID= [%s]", testUserID)
	} else {
		t.Fatalf("failed to RejectVendorVerifyInfo:UserID= [%s]", testUserID)
	}

	// DescribeAppVendorStatistics
	/*===============================================================================================*/

	logger.Info(nil, "%s", "/*=================================================*/")
	logger.Info(nil, "%s", "Test6 DescribeAppVendorStatistics**************************")
	describeStaParams := app_vendor_manager.NewDescribeAppVendorStatisticsParams()
	logger.Info(nil, "test describeParams=[%+v]", describeParams)
	describeParams.WithUserID([]string{testUserID})

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
