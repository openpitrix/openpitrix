// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// +build ignore

package iam

import (
	"context"
	"fmt"
	"time"

	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/gerr"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb/iam"
	"openpitrix.io/openpitrix/pkg/pi"
)

var (
	_ pbiam.TokenManagerServer = (*Server)(nil)
)

func (p *Server) CreateClient(ctx context.Context, req *pbiam.CreateClientRequest) (*pbiam.CreateClientResponse, error) {
	var (
		user_id       = req.GetUserId()
		client_id     = models.NewUserClientId()
		client_secret = p.TokenConfig.Secret
		description   = ""
	)

	_, err := pi.Global().DB(ctx).
		InsertInto(models.UserClientTableName).
		Columns(models.UserClientColumns...).
		Record(models.NewUserClient(user_id, client_id, client_secret, description)).
		Exec()
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
	}

	reply := &pbiam.CreateClientResponse{
		UserId:       user_id,
		ClientId:     client_id,
		ClientSecret: client_secret,
	}
	return reply, nil
}

func (p *Server) Auth(ctx context.Context, req *pbiam.AuthRequest) (*pbiam.AuthResponse, error) {
	if req.Username == "" || req.Password == "" || req.ClientId != "" || req.ClientSecret != "" {
		err := fmt.Errorf("invalid req(%q)", req)
		logger.Warn(nil, "%+v", err)
		return nil, err
	}

	var (
		userInfo       models.User
		userClientInfo models.UserClient
	)

	// query user info from mysql
	{
		query := pi.Global().DB(ctx).
			Select(models.UserColumns...).
			From(models.UserTableName).Limit(1).
			Where(db.Eq(models.ColumnName, req.Username))
		err := query.LoadOne(&userInfo)
		if err != nil {
			return nil, fmt.Errorf("user(%q) not fount", req.Username)
		}
	}

	// query client info from mysql
	{
		query := pi.Global().DB(ctx).
			Select(models.UserClientColumns...).
			From(models.UserClientTableName).Limit(1).
			Where(db.Eq(models.ColumnName, req.ClientId))
		err := query.LoadOne(&userClientInfo)
		if err != nil {
			return nil, fmt.Errorf("user_client(%q) not fount", req.ClientId)
		}
	}

	tokStr, err := MakeJwtToken(p.TokenConfig.Secret, func(opt *JwtToken) {
		opt.UserId = userInfo.UserId
		opt.ClientId = req.ClientId
		opt.TokenType = TokenType_ID
		opt.ExpiresAt = int64(time.Second * time.Duration(p.DurationSeconds))
	})
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return nil, err
	}

	reply := &pbiam.AuthResponse{
		IdToken: tokStr,
	}

	return reply, nil
}

func (p *Server) Token(ctx context.Context, req *pbiam.TokenRequest) (*pbiam.TokenResponse, error) {
	err := fmt.Errorf("TODO")
	logger.Warn(nil, "%+v", err)
	return nil, err
}
