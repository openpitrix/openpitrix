// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// +build integration

package account

import (
	"testing"

	"github.com/stretchr/testify/require"

	"openpitrix.io/openpitrix/test/client/account_manager"
	"openpitrix.io/openpitrix/test/models"
	"openpitrix.io/openpitrix/test/testutil"
)

var clientConfig = testutil.GetClientConfig()

func TestLogin(t *testing.T) {
	client := testutil.GetClient(clientConfig)
	validateParams := account_manager.NewValidateUserPasswordParams()
	validateParams.SetBody(&models.OpenpitrixValidateUserPasswordRequest{
		Email:    "admin@op.com",
		Password: "passw0rd",
	})
	validateResp, err := client.AccountManager.ValidateUserPassword(validateParams, nil)
	require.NoError(t, err)
	require.True(t, validateResp.Payload.Validated)
}

func TestDescribeUsers(t *testing.T) {
	client := testutil.GetClient(clientConfig)
	describeParams := account_manager.NewDescribeUsersParams()
	describeResp, err := client.AccountManager.DescribeUsers(describeParams, nil)
	require.NoError(t, err)
	require.True(t, describeResp.Payload.TotalCount > 0)
}
