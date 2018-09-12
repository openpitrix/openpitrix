// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package iam

import (
	"context"
	"fmt"

	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/gerr"
	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb/iam"
	"openpitrix.io/openpitrix/pkg/pi"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
)

var (
	_ pbiam.AccountManagerServer = (*Server)(nil)
)

func (p *Server) DescribeUsers(ctx context.Context, req *pbiam.DescribeUsersRequest) (*pbiam.DescribeUsersResponse, error) {
	var (
		offset = pbutil.GetOffsetFromRequest(req)
		limit  = pbutil.GetLimitFromRequest(req)
	)

	var query = pi.Global().DB(ctx).
		Select(models.UserColumns...).
		From(models.UserTableName).
		Offset(offset).
		Limit(limit).
		Where(manager.BuildFilterConditions(req, models.UserTableName))

	var users []*models.User
	_, err := query.Load(&users)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}

	count, err := query.Count()
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}

	reply := &pbiam.DescribeUsersResponse{
		UserSet:    models.UsersToPbs(users),
		TotalCount: count,
	}
	return reply, nil
}

func (p *Server) DescribeGroups(ctx context.Context, req *pbiam.DescribeGroupsRequest) (*pbiam.DescribeGroupsResponse, error) {
	var (
		offset = pbutil.GetOffsetFromRequest(req)
		limit  = pbutil.GetLimitFromRequest(req)
	)

	var query = pi.Global().DB(ctx).
		Select(models.GroupColumns...).
		From(models.GroupTableName).
		Offset(offset).
		Limit(limit).
		Where(manager.BuildFilterConditions(req, models.GroupTableName))

	var groups []*models.Group
	_, err := query.Load(&groups)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}

	count, err := query.Count()
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}

	reply := &pbiam.DescribeGroupsResponse{
		GroupSet:   models.GroupsToPbs(groups),
		TotalCount: count,
	}
	return reply, nil
}

func (p *Server) ModifyUser(ctx context.Context, req *pbiam.ModifyUserRequest) (*pbiam.ModifyUserResponse, error) {
	if req.UserId == nil {
		return &pbiam.ModifyUserResponse{}, nil
	}

	var m = make(map[string]interface{})
	if req.Email != nil {
		m["email"] = req.Email.GetValue()
	}
	if req.Username != nil {
		m["name"] = req.Username.GetValue()
	}
	if req.Role != nil {
		m["role"] = req.Role.GetValue()
	}

	_, err := pi.Global().DB(ctx).
		Update(models.UserTableName).
		SetMap(m).
		Where(db.Eq(models.ColumnRuntimeId, []string{})).
		Exec()
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDeleteResourcesFailed)
	}

	return &pbiam.ModifyUserResponse{}, nil
}

func (p *Server) DeleteUsers(ctx context.Context, req *pbiam.DeleteUsersRequest) (*pbiam.DeleteUsersResponse, error) {
	_, err := pi.Global().DB(ctx).
		DeleteFrom(models.UserTableName).
		Where(db.Eq(models.ColumnName, req.UserId)).
		Exec()
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDeleteResourcesFailed)
	}

	reply := &pbiam.DeleteUsersResponse{
		UserId: req.UserId,
	}

	return reply, nil
}

func (p *Server) CreateUser(ctx context.Context, req *pbiam.CreateUserRequest) (*pbiam.CreateUserResponse, error) {
	var newUser = models.NewUser(
		getUsernameFromEmail(req.GetEmail().GetValue()),
		req.GetPassword().GetValue(),
		req.GetEmail().GetValue(),
		req.GetRole().GetValue(),
		"",
	)

	_, err := pi.Global().DB(ctx).
		InsertInto(models.UserTableName).
		Columns(models.UserColumns...).
		Record(newUser).
		Exec()
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
	}

	return &pbiam.CreateUserResponse{}, nil
}

func (p *Server) CreatePasswordReset(ctx context.Context, req *pbiam.CreatePasswordResetRequest) (*pbiam.CreatePasswordResetResponse, error) {

	/*

		var userId = req.GetUserId().GetValue()
		var userInfo models.User

		query := pi.Global().DB(ctx).
			Select(models.UserColumns...).
			From(models.UserTableName).Limit(1).
			Where(db.Eq(models.ColumnName, userId))
		err := query.LoadOne(&userInfo)
		if err != nil {
			return nil, fmt.Errorf("user(%q) not fount", userId)
		}

		if userInfo.Password != req.GetPassword().GetValue() {
			return nil, fmt.Errorf("user(%q) password failed", req.UserId)
		}

		tokStr, err := MakeJwtToken(p.TokenConfig.Secret, func(opt *JwtToken) {
			opt.UserId = userId
			opt.TokenType = TokenType_ResetPassword
			opt.ExpiresAt = int64(time.Second * time.Duration(p.DurationSeconds))
		})
		if err != nil {
			return nil, err
		}

		reply := &pbiam.CreatePasswordResetResponse{
			UserId:  pbutil.ToProtoString(userId),
			ResetId: pbutil.ToProtoString(tokStr),
		}

		return reply, nil
	*/

	return nil, fmt.Errorf("TODO")
}

func (p *Server) ChangePassword(ctx context.Context, req *pbiam.ChangePasswordRequest) (*pbiam.ChangePasswordResponse, error) {

	return nil, fmt.Errorf("TODO")

	/*
		token, err := ValidateJwtToken(req.GetResetId().GetValue(), p.TokenConfig.Secret)
		if err != nil {
			return nil, err
		}

		if token.TokenType != TokenType_ResetPassword {
			return nil, fmt.Errorf("invalid token type")
		}

		_, err = pi.Global().DB(ctx).
			Update(models.UserTableName).
			Set("password", req.GetNewPassword().GetValue()).
			Where(db.Eq("id", []string{token.UserId})).
			Exec()
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDeleteResourcesFailed)
		}

		reply := &pbiam.ChangePasswordResponse{
			UserId: pbutil.ToProtoString(token.UserId),
		}

		return reply, nil
	*/
}

func (p *Server) GetPasswordReset(context.Context, *pbiam.GetPasswordResetRequest) (*pbiam.GetPasswordResetResponse, error) {
	// reset id is token
	// parse user id from jwt token
	return nil, fmt.Errorf("TODO")
}
func (p *Server) ValidateUserPassword(context.Context, *pbiam.ValidateUserPasswordRequest) (*pbiam.ValidateUserPasswordResponse, error) {
	// email => name
	return nil, fmt.Errorf("TODO")
}

func (p *Server) CreateGroup(ctx context.Context, req *pbiam.CreateGroupRequest) (*pbiam.CreateGroupResponse, error) {
	var newGroup = models.NewGroup(
		req.GetName().GetValue(),
		req.GetDescription().GetValue(),
	)

	_, err := pi.Global().DB(ctx).
		InsertInto(models.GroupTableName).
		Columns(models.GroupColumns...).
		Record(newGroup).
		Exec()
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
	}

	return &pbiam.CreateGroupResponse{}, nil
}

func (p *Server) ModifyGroup(ctx context.Context, req *pbiam.ModifyGroupRequest) (*pbiam.ModifyGroupResponse, error) {
	if req.GroupId == nil {
		return &pbiam.ModifyGroupResponse{}, nil
	}

	var m = make(map[string]interface{})
	if req.Name != nil {
		m["name"] = req.Name.GetValue()
	}
	if req.Description != nil {
		m["description"] = req.Description.GetValue()
	}

	_, err := pi.Global().DB(ctx).
		Update(models.UserTableName).
		SetMap(m).
		Where(db.Eq(models.ColumnRuntimeId, []string{})).
		Exec()
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDeleteResourcesFailed)
	}

	return &pbiam.ModifyGroupResponse{}, nil
}
func (p *Server) DeleteGroups(ctx context.Context, req *pbiam.DeleteGroupsRequest) (*pbiam.DeleteGroupsResponse, error) {
	_, err := pi.Global().DB(ctx).
		DeleteFrom(models.GroupTableName).
		Where(db.Eq(models.ColumnName, req.GroupId)).
		Exec()
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDeleteResourcesFailed)
	}

	reply := &pbiam.DeleteGroupsResponse{
		GroupId: req.GroupId,
	}

	return reply, nil
}

func (p *Server) JoinGroup(ctx context.Context, req *pbiam.JoinGroupRequest) (*pbiam.JoinGroupResponse, error) {
	if len(req.GroupId) != 1 || len(req.UserId) != 1 {
		return nil, fmt.Errorf("TODO")
	}

	var newGroupMember = models.NewGroupMember(
		req.GroupId[0],
		req.UserId[0],
	)

	_, err := pi.Global().DB(ctx).
		InsertInto(models.GroupMemberTableName).
		Record(newGroupMember).
		Exec()
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
	}

	return &pbiam.JoinGroupResponse{}, nil
}

func (p *Server) LeaveGroup(ctx context.Context, req *pbiam.LeaveGroupRequest) (*pbiam.LeaveGroupResponse, error) {
	if len(req.GroupId) != 1 || len(req.UserId) != 1 {
		return nil, fmt.Errorf("TODO")
	}

	_, err := pi.Global().DB(ctx).
		DeleteFrom(models.GroupMemberTableName).
		Where(db.Eq("group_id", req.GroupId)).
		Where(db.Eq("user_id", req.UserId)).
		Exec()
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
	}

	return &pbiam.LeaveGroupResponse{}, nil
}
