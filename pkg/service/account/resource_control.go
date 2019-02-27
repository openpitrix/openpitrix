// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package account

import (
	"context"

	pbim "kubesphere.io/im/pkg/pb"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/gerr"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pi"
	"openpitrix.io/openpitrix/pkg/util/stringutil"
)

func getClient(ctx context.Context, clientId, clientSecret string) (*models.UserClient, error) {
	var userClient = models.UserClient{}
	err := pi.Global().DB(ctx).
		Select(models.UserClientColumns...).
		From(constants.TableUserClient).
		Where(db.Eq(constants.ColumnClientId, clientId)).
		Where(db.Eq(constants.ColumnClientSecret, clientSecret)).
		LoadOne(&userClient)
	if err != nil {
		return nil, gerr.New(ctx, gerr.PermissionDenied, gerr.ErrorAuthFailure)
	}
	return &userClient, nil
}

type clientIface interface {
	GetClientId() string
	GetClientSecret() string
}

func validateClientCredentials(ctx context.Context, client clientIface) (*pbim.User, *models.UserClient, error) {
	userClient, err := getClient(ctx, client.GetClientId(), client.GetClientSecret())
	if err != nil {
		return nil, nil, err
	}
	if userClient.Status != constants.StatusActive {
		logger.Error(ctx, "User client [%+v] is not active", userClient)
		return nil, nil, gerr.New(ctx, gerr.PermissionDenied, gerr.ErrorAuthFailure)
	}
	userId := userClient.UserId

	if stringutil.StringIn(userId, constants.InternalUsers) {
		return &pbim.User{
			UserId: userId,
			Extra:  map[string]string{"role": constants.RoleGlobalAdmin},
		}, userClient, nil
	}

	// check the credential's user
	user, err := getUser(ctx, userId)
	if err != nil {
		return nil, nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorResourceNotFound, userId)
	}
	if user.Status != constants.StatusActive {
		logger.Error(ctx, "User [%s] is not active user", userId)
		return nil, nil, gerr.New(ctx, gerr.PermissionDenied, gerr.ErrorAuthFailure)
	}
	if user.Extra["role"] != constants.RoleGlobalAdmin {
		logger.Error(ctx, "User [%s] is not global admin", userId)
		return nil, nil, gerr.New(ctx, gerr.PermissionDenied, gerr.ErrorAuthFailure)
	}
	return user, userClient, nil
}

func getTokenByRefreshToken(ctx context.Context, refreshToken string) (*models.Token, error) {
	var token = models.Token{}
	err := pi.Global().DB(ctx).
		Select(models.TokenColumns...).
		From(constants.TableToken).
		Where(db.Eq(constants.ColumnRefreshToken, refreshToken)).
		Where(db.Eq(constants.ColumnStatus, constants.StatusActive)).
		LoadOne(&token)
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func getLastToken(ctx context.Context, clientId, userId, scope string) (*models.Token, error) {
	var token = models.Token{}
	err := pi.Global().DB(ctx).
		Select(models.TokenColumns...).
		From(constants.TableToken).
		Where(db.Eq(constants.ColumnUserId, userId)).
		Where(db.Eq(constants.ColumnClientId, clientId)).
		Where(db.Eq(constants.ColumnScope, scope)).
		Where(db.Eq(constants.ColumnStatus, constants.StatusActive)).
		OrderDir(constants.ColumnCreateTime, false).
		LoadOne(&token)
	if err != nil {
		if err == db.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &token, nil
}

func newToken(ctx context.Context, clientId, scope, userId string) (*models.Token, error) {
	var err error
	var token *models.Token
	var i = 0
	for {
		i++
		if i == 10 {
			return nil, err
		}
		token = models.NewToken(clientId, userId, scope)
		_, err = pi.Global().DB(ctx).InsertInto(constants.TableToken).Record(token).Exec()
		if err != nil {
			continue
		}
		return token, nil
	}
}
