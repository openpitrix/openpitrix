// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package account

import (
	"context"
	"os"

	pbim "openpitrix.io/iam/pkg/pb/im"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/pi"
)

func initIAMClient() {
	clientId := os.Getenv("IAM_INIT_CLIENT_ID")
	clientSecret := os.Getenv("IAM_INIT_CLIENT_SECRET")
	const userId = "system"

	if clientId == "" || clientSecret == "" {
		return
	}
	_, err := pi.Global().DB(context.Background()).InsertBySql(
		`insert into user_client (client_id, user_id, client_secret, status, description)
values (?, ?, ?, 'active', '')
on duplicate key update user_id = ?, client_secret = ?, status = 'active';`,
		clientId, userId, clientSecret,
		userId, clientSecret,
	).Exec()
	if err != nil {
		logger.Error(nil, "Init default IAM client [%s] [%s] failed", clientId, clientSecret)
		panic(err)
	}
	logger.Info(nil, "Init IAM client [%s] done", clientId)
}

func initIAMAccount() {
	email := os.Getenv("IAM_INIT_ACCOUNT_EMAIL")
	password := os.Getenv("IAM_INIT_ACCOUNT_PASSWORD")

	if email == "" || password == "" {
		return
	}
	ctx := context.Background()
	user, b, err := validateUserPassword(ctx, email, password)
	if err != nil {
		logger.Info(ctx, "Validate user password failed, create new user")
		// create user
		_, err = imClient.CreateUser(ctx, &pbim.User{
			Email:    email,
			Username: getUsernameFromEmail(email),
			Password: password,
			Status:   constants.StatusActive,
			Extra:    map[string]string{"role": constants.RoleGlobalAdmin},
		})
		if err != nil {
			logger.Info(ctx, "Create new user failed, error: %+v", err)
		}
		return
	}
	if b {
		// password matched
		return
	}
	userId := user.UserId
	_, err = imClient.ModifyPassword(ctx, &pbim.Password{
		UserId:   userId,
		Password: password,
	})
	if err != nil {
		panic(err)
	}
	logger.Info(ctx, "Init IAM admin account [%s] done, update [%s] password", email, userId)
}
