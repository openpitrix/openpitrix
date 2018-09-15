// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package iam

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/gerr"
	"openpitrix.io/openpitrix/pkg/logger"
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
		Select(models.UserColumns...).From(models.UserTableName).
		Offset(offset).Limit(limit).
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
		Select(models.GroupColumns...).From(models.GroupTableName).
		Offset(offset).Limit(limit).
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
	var attributes = manager.BuildUpdateAttributes(req,
		"username", "email", "role", "status", "description",
	)
	if len(attributes) == 0 {
		return &pbiam.ModifyUserResponse{}, nil
	}

	_, err := pi.Global().DB(ctx).
		Update(models.UserTableName).SetMap(attributes).
		Where(db.Eq(models.ColumnUserId, req.UserId.GetValue())).
		Exec()
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDeleteResourcesFailed)
	}

	reply := &pbiam.ModifyUserResponse{
		UserId: req.UserId,
	}

	return reply, nil
}

func (p *Server) DeleteUsers(ctx context.Context, req *pbiam.DeleteUsersRequest) (*pbiam.DeleteUsersResponse, error) {
	_, err := pi.Global().DB(ctx).
		DeleteFrom(models.UserTableName).
		Where(db.Eq(models.ColumnUserId, req.UserId)).
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
	email := req.GetEmail().GetValue()

	// check email exists
	query := pi.Global().DB(ctx).
		Select(models.UserColumns...).
		From(models.UserTableName).Limit(1).
		Where(db.Eq(models.ColumnEmail, email))
	if count, err := query.Count(); err == nil && count > 0 {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal,
			fmt.Errorf("email(%q) exists", email),
			gerr.ErrorCreateResourcesFailed,
		)
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.GetPassword().GetValue()), bcrypt.DefaultCost)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
	}

	var newUser = models.NewUser(
		getUsernameFromEmail(req.GetEmail().GetValue()),
		string(hashedPass),
		req.GetEmail().GetValue(),
		req.GetRole().GetValue(),
		req.GetDescription().GetValue(),
	)

	_, err = pi.Global().DB(ctx).
		InsertInto(models.UserTableName).
		Columns(models.UserColumns...).
		Record(newUser).
		Exec()
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
	}

	reply := &pbiam.CreateUserResponse{
		UserId: pbutil.ToProtoString(newUser.UserId),
	}

	return reply, nil
}

func (p *Server) CreatePasswordReset(ctx context.Context, req *pbiam.CreatePasswordResetRequest) (*pbiam.CreatePasswordResetResponse, error) {
	var user_id = req.GetUserId().GetValue()
	var user_info models.User

	query := pi.Global().DB(ctx).
		Select(models.UserColumns...).
		From(models.UserTableName).
		Where(db.Eq(models.ColumnUserId, user_id))
	err := query.LoadOne(&user_info)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal,
			fmt.Errorf("user(%q) not fount", user_id),
			gerr.ErrorCreateResourcesFailed,
		)
	}

	if user_info.Password != req.GetPassword().GetValue() {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal,
			fmt.Errorf("ivalid password"),
			gerr.ErrorCreateResourcesFailed,
		)
	}

	var newUserPasswordReset = models.NewUserPasswordReset(user_id)

	_, err = pi.Global().DB(ctx).
		InsertInto(models.UserPasswordResetTableName).
		Columns(models.UserPasswordResetColumns...).
		Record(newUserPasswordReset).
		Exec()
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
	}

	reply := &pbiam.CreatePasswordResetResponse{
		UserId:  pbutil.ToProtoString(user_id),
		ResetId: pbutil.ToProtoString(newUserPasswordReset.ResetId),
	}

	return reply, nil
}

func (p *Server) ChangePassword(ctx context.Context, req *pbiam.ChangePasswordRequest) (*pbiam.ChangePasswordResponse, error) {
	reset_id := req.GetResetId().GetValue()
	new_password := req.GetNewPassword().GetValue()

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(new_password), bcrypt.DefaultCost)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorUpgradeResourceFailed)
	}

	var reset_info models.UserPasswordReset

	query := pi.Global().DB(ctx).
		Select(models.UserColumns...).
		From(models.UserPasswordResetTableName).Limit(1).
		Where(db.Eq(models.ColumnResetId, reset_id))
	err = query.LoadOne(&reset_info)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDeleteResourcesFailed)
	}

	_, err = pi.Global().DB(ctx).
		Update(models.UserTableName).Set(models.ColumnPassword, string(hashedPass)).
		Where(db.Eq(models.ColumnUserId, reset_info.UserId)).
		Exec()
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDeleteResourcesFailed)
	}

	reply := &pbiam.ChangePasswordResponse{
		UserId: pbutil.ToProtoString(reset_info.UserId),
	}

	return reply, nil
}

func (p *Server) GetPasswordReset(ctx context.Context, req *pbiam.GetPasswordResetRequest) (*pbiam.GetPasswordResetResponse, error) {
	var reset_id = req.GetResetId().GetValue()
	var reset_info models.UserPasswordReset

	query := pi.Global().DB(ctx).
		Select(models.UserColumns...).
		From(models.UserPasswordResetTableName).Limit(1).
		Where(db.Eq(models.ColumnResetId, reset_id))
	err := query.LoadOne(&reset_info)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
	}

	reply := &pbiam.GetPasswordResetResponse{
		ResetId: pbutil.ToProtoString(reset_id),
		UserId:  pbutil.ToProtoString(reset_info.ResetId),
	}

	return reply, nil
}
func (p *Server) ValidateUserPassword(ctx context.Context, req *pbiam.ValidateUserPasswordRequest) (*pbiam.ValidateUserPasswordResponse, error) {
	var email = req.GetEmail()
	var userInfo models.User

	query := pi.Global().DB(ctx).
		Select(models.UserColumns...).
		From(models.UserTableName).Limit(1).
		Where(db.Eq(models.ColumnEmail, email))
	err := query.LoadOne(&userInfo)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal,
			fmt.Errorf("user(%q) not fount", email),
			gerr.ErrorCreateResourcesFailed,
		)
	}

	var validated = true
	err = bcrypt.CompareHashAndPassword([]byte(userInfo.Password), []byte(req.GetPassword()))
	if err != nil {
		validated = false
	}

	reply := &pbiam.ValidateUserPasswordResponse{
		Validated: validated,
	}

	return reply, nil
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

	reply := &pbiam.CreateGroupResponse{
		GroupId: pbutil.ToProtoString(newGroup.GroupId),
	}

	return reply, nil
}

func (p *Server) ModifyGroup(ctx context.Context, req *pbiam.ModifyGroupRequest) (*pbiam.ModifyGroupResponse, error) {
	group_id := req.GetGroupId().GetValue()

	var attributes = manager.BuildUpdateAttributes(req,
		"name", "status", "description",
	)

	if len(attributes) == 0 {
		return &pbiam.ModifyGroupResponse{}, nil
	}

	_, err := pi.Global().DB(ctx).
		Update(models.GroupTableName).SetMap(attributes).
		Where(db.Eq(models.ColumnGroupId, group_id)).
		Exec()
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDeleteResourcesFailed)
	}

	reply := &pbiam.ModifyGroupResponse{
		GroupId: pbutil.ToProtoString(group_id),
	}

	return reply, nil
}
func (p *Server) DeleteGroups(ctx context.Context, req *pbiam.DeleteGroupsRequest) (*pbiam.DeleteGroupsResponse, error) {
	_, err := pi.Global().DB(ctx).
		DeleteFrom(models.GroupTableName).
		Where(db.Eq(models.ColumnGroupId, req.GroupId)).
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
	var reply = &pbiam.JoinGroupResponse{}
	var lastErr error

	for _, gid := range req.GroupId {
		for _, uid := range req.UserId {
			query := pi.Global().DB(ctx).
				Select().From(models.GroupMemberTableName).Limit(1).
				Where(db.Eq(models.ColumnGroupId, gid)).
				Where(db.Eq(models.ColumnUserId, uid))
			if count, err := query.Count(); err != nil || count > 0 {
				logger.Warn(ctx, "gid(%v)/uid(%v): %v", gid, uid, err)

				if err != nil {
					lastErr = err
				}
				continue
			}

			_, err := pi.Global().DB(ctx).
				InsertInto(models.GroupMemberTableName).
				Record(models.NewGroupMember(gid, uid)).
				Exec()
			if err != nil {
				logger.Warn(ctx, "gid(%v)/uid(%v): %v", gid, uid, err)

				if err != nil {
					lastErr = err
				}
				continue
			}

			reply.GroupId = append(reply.GroupId, gid)
			reply.UserId = append(reply.UserId, uid)
		}
	}
	if lastErr != nil {
		return reply, gerr.NewWithDetail(ctx, gerr.Internal, lastErr, gerr.ErrorUpgradeResourceFailed)
	}

	return reply, nil
}

func (p *Server) LeaveGroup(ctx context.Context, req *pbiam.LeaveGroupRequest) (*pbiam.LeaveGroupResponse, error) {
	_, err := pi.Global().DB(ctx).
		DeleteFrom(models.GroupMemberTableName).
		Where(db.And(
			db.Eq(models.ColumnGroupId, req.GroupId),
			db.Eq(models.ColumnUserId, req.UserId),
		)).Exec()
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorUpgradeResourceFailed)
	}

	var reply = &pbiam.LeaveGroupResponse{
		GroupId: req.GroupId,
		UserId:  req.UserId,
	}

	return reply, nil
}
