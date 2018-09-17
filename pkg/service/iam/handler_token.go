// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

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
	"openpitrix.io/openpitrix/pkg/util/jwtutil"
)

var (
	_ pbiam.TokenManagerServer = (*Server)(nil)
)

func (p *Server) CreateClient(ctx context.Context, req *pbiam.CreateClientRequest) (*pbiam.CreateClientResponse, error) {
	// unsupport ClientSecret
	err := fmt.Errorf("Unimplemented")
	logger.Warn(nil, "%+v", err)
	return nil, err

	var (
		user_id       = req.GetUserId()
		client_id     = models.NewUserClientId()
		client_secret = "12345678"
		description   = ""
	)

	_, err = pi.Global().DB(ctx).
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
	// unsupport ClientSecret
	if req.ClientId != "" && req.ClientSecret != "" {
		err := fmt.Errorf("Unimplemented")
		logger.Warn(nil, "%+v", err)
		return nil, err
	}

	// user & password
	if req.Username == "" && req.Password == "" {
		err := fmt.Errorf("invalid user(%q) or password(%q)", req.Username, req.Password)
		logger.Warn(nil, "%+v", err)
		return nil, err
	}

	// query user info from mysql
	query := pi.Global().DB(ctx).
		Select(models.UserColumns...).
		From(models.UserTableName).Limit(1).
		Where(db.Eq(models.ColumnName, req.Username))

	var userInfo models.User
	_, err := query.Load(&userInfo)
	if err != nil {
		return nil, fmt.Errorf("user(%q) not fount", req.Username)
	}

	tokStr, err := jwtutil.MakeToken(userInfo.Id,
		p.TokenConfig.Secret, time.Second*time.Duration(p.DurationSeconds),
		nil,
	)
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return nil, err
	}

	reply := &pbiam.AuthResponse{
		TokenType:    "", // TODO
		ExpiresIn:    "", // TODO
		AccessToken:  "", // TODO
		RefreshToken: "", // TODO
		IdToken:      tokStr,
	}

	return reply, nil
}

func (p *Server) Token(ctx context.Context, req *pbiam.TokenRequest) (*pbiam.TokenResponse, error) {
	err := fmt.Errorf("Unimplemented")
	logger.Warn(nil, "%+v", err)
	return nil, err
}
