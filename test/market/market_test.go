// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// +build integration

package market

import (
	"testing"

	"github.com/stretchr/testify/require"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/test/client/market_manager"
	"openpitrix.io/openpitrix/test/models"
	"openpitrix.io/openpitrix/test/testutil"
)

var clientConfig = testutil.GetClientConfig()

func TestMarket(t *testing.T) {
	client := testutil.GetClient(clientConfig)
	testMarketName := "test_market"
	testVisibility1 := constants.VisibilityPublic
	testVisibility2 := constants.VisibilityPrivate
	testStatus := constants.StatusDisabled
	testDescription := "test_create_market"

	// delete old market

	// create market
	createMarketParams := market_manager.NewCreateMarketParams()
	createMarketParams.SetBody(
		&models.OpenpitrixCreateMarketRequest{
			Name:        testMarketName,
			Visibility:  testVisibility1,
			Description: testDescription,
		})
	createMarketResp, err := client.MarketManager.CreateMarket(createMarketParams, nil)

	require.NoError(t, err)
	t.Log(createMarketResp)

	marketId1 := createMarketResp.Payload.MarketID
	// modify market
	testMarketName2 := "test_market2"
	modifyMarketParams := market_manager.NewModifyMarketParams()
	modifyMarketParams.SetBody(
		&models.OpenpitrixModifyMarketRequest{
			MarketID:    marketId1,
			Name:        testMarketName2,
			Status:      testStatus,
			Visibility:  testVisibility2,
			Description: testDescription,
		})
	modifyMarketResp, err := client.MarketManager.ModifyMarket(modifyMarketParams, nil)

	require.NoError(t, err)

	t.Log(modifyMarketResp)

	// describe market
	describeMarketsParams := market_manager.NewDescribeMarketsParams()
	describeMarketsParams.SetMarketID([]string{marketId1})
	describeMarketsResp, err := client.MarketManager.DescribeMarkets(describeMarketsParams, nil)
	require.NoError(t, err)

	markets := describeMarketsResp.Payload.MarketSet
	require.Equal(t, 1, len(markets))

	market := markets[0]

	require.Equal(t, testMarketName2, market.Name)
	require.Equal(t, testVisibility2, market.Visibility)
	require.Equal(t, testStatus, market.Status)
	require.Equal(t, testDescription, market.Description)
	t.Log(describeMarketsResp)

	// delete market
	deleteMarketParams := market_manager.NewDeleteMarketsParams()
	deleteMarketParams.WithBody(&models.OpenpitrixDeleteMarketsRequest{
		MarketID: []string{marketId1},
	})
	deleteMarketResp, err := client.MarketManager.DeleteMarkets(deleteMarketParams, nil)
	require.NoError(t, err)
	t.Log(deleteMarketResp)
}

func TestMarketUser(t *testing.T) {
	// create a market for market_user
	client := testutil.GetClient(clientConfig)
	testMarketName := "test_marketUser"
	testVisibility := constants.VisibilityPublic

	createMarektParams := market_manager.NewCreateMarketParams()
	createMarektParams.SetBody(
		&models.OpenpitrixCreateMarketRequest{
			Name:        testMarketName,
			Visibility:  testVisibility,
			Description: "test_market3",
		})

	createResp, err := client.MarketManager.CreateMarket(createMarektParams, nil)
	require.NoError(t, err)
	marketId := createResp.Payload.MarketID

	userId := "system"
	// describe marker_user
	describMarketUserParams := market_manager.NewDescribeMarketUsersParams()
	describMarketUserParams.SetMarketID([]string{marketId})
	describMarketUserParams.SetUserID([]string{userId})

	describMarketUserResp, err := client.MarketManager.DescribeMarketUsers(describMarketUserParams, nil)
	require.NoError(t, err)

	marketUsers := describMarketUserResp.Payload.MarketUserSet
	marketUser := marketUsers[0]

	require.Equal(t, marketId, marketUser.MarketID)
	require.Equal(t, userId, marketUser.UserID)
	require.Equal(t, userId, marketUser.Owner)
	t.Log(describMarketUserResp)

	// user leave market
	userLeaveMarketParams := market_manager.NewUserLeaveMarketParams()
	userLeaveMarketParams.SetBody(
		&models.OpenpitrixUserLeaveMarketRequest{
			MarketID: []string{marketId},
			UserID:   []string{userId},
		})
	leaveResp, err := client.MarketManager.UserLeaveMarket(userLeaveMarketParams, nil)
	require.NoError(t, err)
	t.Log(leaveResp)

	// user join market
	userJoinMarketParams := market_manager.NewUserJoinMarketParams()
	userJoinMarketParams.SetBody(
		&models.OpenpitrixUserJoinMarketRequest{
			MarketID: []string{"mkt-AkEG1pQVGZPL"},
			UserID:   []string{userId},
		})
	joinResp, err := client.MarketManager.UserJoinMarket(userJoinMarketParams, nil)
	require.NoError(t, err)
	t.Log(joinResp)
}
