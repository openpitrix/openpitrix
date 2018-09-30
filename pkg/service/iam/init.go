// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package iam

import (
	"context"
	"os"

	"golang.org/x/crypto/bcrypt"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/models"
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
	user, err := getUser(ctx, map[string]interface{}{
		constants.ColumnEmail: email,
	})
	if err != nil {
		if err == db.ErrNotFound {
			// create new user
			hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if err != nil {
				panic(err)
			}

			var newUser = models.NewUser(
				getUsernameFromEmail(email),
				string(hashedPass),
				email,
				constants.RoleGlobalAdmin,
				"",
			)

			_, err = pi.Global().DB(ctx).
				InsertInto(constants.TableUser).
				Record(newUser).
				Exec()
			if err != nil {
				panic(err)
			}
			logger.Info(nil, "Init IAM admin account [%s] done", email)
		} else {
			panic(err)
		}
		return
	}
	// check password and update
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err == nil {
		logger.Info(nil, "Init IAM admin account [%s] done, no need update [%s] password",
			email, user.UserId)
		return
	}
	// password not matched, update it
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	_, err = pi.Global().DB(ctx).
		Update(constants.TableUser).Set(constants.ColumnPassword, string(hashedPass)).
		Where(db.Eq(constants.ColumnUserId, user.UserId)).
		Exec()
	if err != nil {
		panic(err)
	}
	logger.Info(nil, "Init IAM admin account [%s] done, update [%s] password",
		email, user.UserId)
}
