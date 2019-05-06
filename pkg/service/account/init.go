// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package account

import (
	"context"
	"os"

	pbim "kubesphere.io/im/pkg/pb"

	pbam "openpitrix.io/iam/pkg/pb"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/pi"
	"openpitrix.io/openpitrix/pkg/sender"
	"openpitrix.io/openpitrix/pkg/util/ctxutil"
)

func initIAMClient() {
	clientId := os.Getenv("IAM_INIT_CLIENT_ID")
	clientSecret := os.Getenv("IAM_INIT_CLIENT_SECRET")
	const userId = constants.UserSystem

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
	var userId string

	user, isUserExist, isGroupExist := validateUserAndGroupExist(ctx, email)
	if !isUserExist {
		// create user
		createUserResponse, err := imClient.CreateUser(ctx, &pbim.CreateUserRequest{
			Email:    email,
			Username: getUsernameFromEmail(email),
			Password: password,
		})
		if err != nil {
			panic(err)
		} else {
			logger.Info(ctx, "Create new user with email [%s] done", email)
		}
		userId = createUserResponse.UserId
	} else {
		logger.Info(ctx, "User with email [%s] already exist", email)
		userId = user.UserId
	}

	isEmailPasswordMatched := validateUserPassword(ctx, userId, password)
	if !isEmailPasswordMatched {
		_, err := imClient.ModifyPassword(ctx, &pbim.ModifyPasswordRequest{
			UserId:   userId,
			Password: password,
		})
		if err != nil {
			panic(err)
		} else {
			logger.Info(ctx, "Init IAM admin account [%s] done, update [%s] password", email, userId)
		}
	} else {
		logger.Info(ctx, "User [%s] with email [%s] password no need to update", userId, email)
	}

	ctx = ctxutil.ContextWithSender(ctx, sender.GetSystemSender())
	isRoleBoundUser := validateRoleUser(ctx, constants.RoleGlobalAdmin, userId)
	if !isRoleBoundUser {
		_, err := amClient.BindUserRole(ctx, &pbam.BindUserRoleRequest{
			UserId: []string{userId},
			RoleId: []string{constants.RoleGlobalAdmin},
		})
		if err != nil {
			panic(err)
		} else {
			logger.Info(ctx, "Bind user [%s] global admin role done", userId)
		}
	} else {
		logger.Info(ctx, "User [%s] already bound with global admin role", userId)
	}

	if !isGroupExist {
		err := createAndJoinRootGroup(ctx, userId)
		if err != nil {
			panic(err)
		} else {
			logger.Info(ctx, "Create and join root group for user [%s] done", userId)
		}
	}
}
