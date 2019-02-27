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

var client = testutil.GetClient(testutil.GetClientConfig())
var userId = constants.UserSystem

const Service = "openpitrix-market-manager"

func deleteMarketUser(t *testing.T, marketId string) {
	userLeaveMarketParams := market_manager.NewUserLeaveMarketParams()
	userLeaveMarketParams.SetBody(
		&models.OpenpitrixUserLeaveMarketRequest{
			MarketID: []string{marketId},
			UserID:   []string{userId},
		})
	leaveResp, err := client.MarketManager.UserLeaveMarket(userLeaveMarketParams, nil)
	testutil.NoError(t, err, []string{Service})
	t.Log(leaveResp)
}

// temporary comment
func testMarket(t *testing.T) {

	testMarketName := "test_market"
	testVisibility1 := constants.VisibilityPublic
	testVisibility2 := constants.VisibilityPrivate
	testStatus := constants.StatusEnabled
	testDescription := "test_create_market"

	// create market
	createMarketParams := market_manager.NewCreateMarketParams()
	createMarketParams.SetBody(
		&models.OpenpitrixCreateMarketRequest{
			Name:        testMarketName,
			Visibility:  testVisibility1,
			Description: testDescription,
		})
	createMarketResp, err := client.MarketManager.CreateMarket(createMarketParams, nil)

	testutil.NoError(t, err, []string{Service})
	t.Log(createMarketResp)

	marketId := createMarketResp.Payload.MarketID

	deleteMarketUser(t, marketId)
	// modify market
	modifyMarketParams := market_manager.NewModifyMarketParams()
	modifyMarketParams.SetBody(
		&models.OpenpitrixModifyMarketRequest{
			MarketID:    marketId,
			Name:        testMarketName,
			Status:      testStatus,
			Visibility:  testVisibility2,
			Description: testDescription,
		})
	modifyMarketResp, err := client.MarketManager.ModifyMarket(modifyMarketParams, nil)

	testutil.NoError(t, err, []string{Service})

	t.Log(modifyMarketResp)

	// describe market
	describeMarketsParams := market_manager.NewDescribeMarketsParams()
	describeMarketsParams.SetMarketID([]string{marketId})
	describeMarketsResp, err := client.MarketManager.DescribeMarkets(describeMarketsParams, nil)
	testutil.NoError(t, err, []string{Service})

	markets := describeMarketsResp.Payload.MarketSet
	require.Equal(t, 1, len(markets))

	market := markets[0]

	require.Equal(t, testMarketName, market.Name)
	require.Equal(t, testVisibility2, market.Visibility)
	require.Equal(t, testStatus, market.Status)
	require.Equal(t, testDescription, market.Description)
	t.Log(describeMarketsResp)

	// delete market
	deleteMarketParams := market_manager.NewDeleteMarketsParams()
	deleteMarketParams.WithBody(&models.OpenpitrixDeleteMarketsRequest{
		MarketID: []string{marketId},
	})
	deleteMarketResp, err := client.MarketManager.DeleteMarkets(deleteMarketParams, nil)
	testutil.NoError(t, err, []string{Service})
	t.Log(deleteMarketResp)

	// delete market_user

}

// temporary comment
func testMarketUser(t *testing.T) {
	testMarketName := "test_marketUser"

	// create a market for market_user
	testVisibility := constants.VisibilityPublic

	createMarektParams := market_manager.NewCreateMarketParams()
	createMarektParams.SetBody(
		&models.OpenpitrixCreateMarketRequest{
			Name:        testMarketName,
			Visibility:  testVisibility,
			Description: "test_market3",
		})

	createResp, err := client.MarketManager.CreateMarket(createMarektParams, nil)
	testutil.NoError(t, err, []string{Service})
	marketId := createResp.Payload.MarketID

	// describe marker_user
	describMarketUserParams := market_manager.NewDescribeMarketUsersParams()
	describMarketUserParams.SetMarketID([]string{marketId})
	describMarketUserParams.SetUserID([]string{userId})

	describMarketUserResp, err := client.MarketManager.DescribeMarketUsers(describMarketUserParams, nil)
	testutil.NoError(t, err, []string{Service})

	marketUsers := describMarketUserResp.Payload.MarketUserSet
	marketUser := marketUsers[0]

	require.Equal(t, marketId, marketUser.MarketID)
	require.Equal(t, userId, marketUser.UserID)
	//require.Equal(t, userId, marketUser.Owner)
	t.Log(describMarketUserResp)

	// user leave market
	userLeaveMarketParams := market_manager.NewUserLeaveMarketParams()
	userLeaveMarketParams.SetBody(
		&models.OpenpitrixUserLeaveMarketRequest{
			MarketID: []string{marketId},
			UserID:   []string{userId},
		})
	leaveResp, err := client.MarketManager.UserLeaveMarket(userLeaveMarketParams, nil)
	testutil.NoError(t, err, []string{Service})
	t.Log(leaveResp)

	// user join market
	userJoinMarketParams := market_manager.NewUserJoinMarketParams()
	userJoinMarketParams.SetBody(
		&models.OpenpitrixUserJoinMarketRequest{
			MarketID: []string{marketId},
			UserID:   []string{userId},
		})
	joinResp, err := client.MarketManager.UserJoinMarket(userJoinMarketParams, nil)
	testutil.NoError(t, err, []string{Service})
	t.Log(joinResp)

	// delete old date
	deleteMarketUser(t, marketId)

	deleteParams := market_manager.NewDeleteMarketsParams()
	deleteParams.SetBody(&models.OpenpitrixDeleteMarketsRequest{
		MarketID: []string{marketId},
	})
	_, err = client.MarketManager.DeleteMarkets(deleteParams, nil)
	testutil.NoError(t, err, []string{Service})

}
