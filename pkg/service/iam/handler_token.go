// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package iam

import (
	"context"
	"time"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/gerr"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/pi"
	"openpitrix.io/openpitrix/pkg/util/senderutil"
)

var (
	_ pb.TokenManagerServer = (*Server)(nil)
)

func (p *Server) CreateClient(ctx context.Context, req *pb.CreateClientRequest) (*pb.CreateClientResponse, error) {
	sender := senderutil.GetSenderFromContext(ctx)
	if !sender.IsGlobalAdmin() {
		return nil, gerr.New(ctx, gerr.PermissionDenied, gerr.ErrorPermissionDenied)
	}
	userId := req.UserId
	client := models.NewUserClient(userId)
	_, err := pi.Global().DB(ctx).InsertInto(constants.TableUserClient).Record(client).Exec()
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
	}
	return &pb.CreateClientResponse{
		UserId:       client.UserId,
		ClientId:     client.ClientId,
		ClientSecret: client.ClientSecret,
	}, nil
}

func (p *Server) Token(ctx context.Context, req *pb.TokenRequest) (*pb.TokenResponse, error) {
	if req.GrantType == constants.GrantTypePassword {
		if req.Username == "" {
			return nil, gerr.New(ctx, gerr.InvalidArgument, gerr.ErrorParameterShouldNotBeEmpty, "username")
		}
		if req.Password == "" {
			return nil, gerr.New(ctx, gerr.InvalidArgument, gerr.ErrorParameterShouldNotBeEmpty, "password")
		}
	}
	if req.GrantType == constants.GrantTypeRefreshToken {
		if req.RefreshToken == "" {
			return nil, gerr.New(ctx, gerr.InvalidArgument, gerr.ErrorParameterShouldNotBeEmpty, "refresh_token")
		}
	}
	// validate client credentials
	user, userClient, err := validateClientCredentials(ctx, req)
	if err != nil {
		return nil, err
	}
	// if grant_type is password, switch user
	if req.GrantType == constants.GrantTypePassword {
		user, err = getUser(ctx, map[string]interface{}{constants.ColumnEmail: req.Username})
		if err != nil {
			return nil, gerr.New(ctx, gerr.PermissionDenied, gerr.ErrorEmailPasswordNotMatched)
		}
		if user.Status != constants.StatusActive {
			return nil, gerr.New(ctx, gerr.PermissionDenied, gerr.ErrorPermissionDenied)
		}
		if !validateUserPassword(user, req.Password) {
			return nil, gerr.New(ctx, gerr.PermissionDenied, gerr.ErrorEmailPasswordNotMatched)
		}
	}
	var token *models.Token
	if req.GrantType == constants.GrantTypeRefreshToken {
		token, err = getTokenByRefreshToken(ctx, req.RefreshToken)
		if err != nil {
			return nil, gerr.New(ctx, gerr.PermissionDenied, gerr.ErrorAuthFailure)
		}
		if token.CreateTime.Add(p.RefreshTokenExpireTime).Unix() <= time.Now().Unix() {
			return nil, gerr.New(ctx, gerr.PermissionDenied, gerr.ErrorRefreshTokenExpired)
		}
	} else {
		// reuse exist token
		token, err = getLastToken(ctx, userClient.ClientId, user.UserId, req.Scope)
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
		}
		// token not exists or expired
		if token == nil || token.CreateTime.Add(p.RefreshTokenExpireTime).Unix() <= time.Now().Unix() {
			// generate access token
			token, err = newToken(ctx, userClient.ClientId, req.Scope, user)
			if err != nil {
				return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
			}
		}
	}

	accessToken, err := senderutil.Generate(
		p.IAMConfig.SecretKey, p.IAMConfig.ExpireTime, user.UserId, user.Role,
	)
	if err != nil {
		return nil, gerr.New(ctx, gerr.Internal, gerr.ErrorInternalError)
	}

	return &pb.TokenResponse{
		TokenType:    senderutil.TokenType,
		ExpiresIn:    int32(p.ExpireTime.Seconds()),
		AccessToken:  accessToken,
		RefreshToken: token.RefreshToken,
	}, nil
}
