// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package iam

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/gerr"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/pi"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
)

var (
	_ pb.AccountManagerServer = (*Server)(nil)
)

func (p *Server) DescribeUsers(ctx context.Context, req *pb.DescribeUsersRequest) (*pb.DescribeUsersResponse, error) {
	var (
		offset = pbutil.GetOffsetFromRequest(req)
		limit  = pbutil.GetLimitFromRequest(req)
	)

	var query = pi.Global().DB(ctx).
		Select(models.UserColumns...).From(constants.TableUser).
		Offset(offset).Limit(limit).
		Where(manager.BuildFilterConditions(req, constants.TableUser))

	query = manager.AddQueryOrderDir(query, req, constants.ColumnCreateTime)

	var users []*models.User
	_, err := query.Load(&users)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}

	count, err := query.Count()
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}

	reply := &pb.DescribeUsersResponse{
		UserSet:    models.UsersToPbs(users),
		TotalCount: count,
	}
	return reply, nil
}

func (p *Server) DescribeGroups(ctx context.Context, req *pb.DescribeGroupsRequest) (*pb.DescribeGroupsResponse, error) {
	// TODO: add filter condition
	var (
		offset = pbutil.GetOffsetFromRequest(req)
		limit  = pbutil.GetLimitFromRequest(req)
	)

	var query = pi.Global().DB(ctx).
		Select(models.GroupColumns...).
		From(constants.TableGroup).
		Offset(offset).Limit(limit).
		Where(manager.BuildFilterConditions(req, constants.TableGroup))

	query = manager.AddQueryOrderDir(query, req, constants.ColumnCreateTime)

	var groups []*models.Group
	_, err := query.Load(&groups)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}

	count, err := query.Count()
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}

	reply := &pb.DescribeGroupsResponse{
		GroupSet:   models.GroupsToPbs(groups),
		TotalCount: count,
	}
	return reply, nil
}

func (p *Server) ModifyUser(ctx context.Context, req *pb.ModifyUserRequest) (*pb.ModifyUserResponse, error) {
	password := req.GetPassword().GetValue()
	if password != "" {
		hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
		}
		req.Password = pbutil.ToProtoString(string(hashedPass))
	}

	var attributes = manager.BuildUpdateAttributes(req,
		"username", "email", "role", "status", "description", "password",
	)
	if len(attributes) == 0 {
		return &pb.ModifyUserResponse{
			UserId: req.UserId,
		}, nil
	}

	_, err := pi.Global().DB(ctx).
		Update(constants.TableUser).
		SetMap(attributes).
		Where(db.Eq(constants.ColumnUserId, req.UserId.GetValue())).
		Exec()
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDeleteResourcesFailed)
	}

	reply := &pb.ModifyUserResponse{
		UserId: req.UserId,
	}

	return reply, nil
}

func (p *Server) DeleteUsers(ctx context.Context, req *pb.DeleteUsersRequest) (*pb.DeleteUsersResponse, error) {
	_, err := pi.Global().DB(ctx).
		Update(constants.TableUser).
		Set(constants.ColumnStatus, constants.StatusDeleted).
		Where(db.Eq(constants.ColumnUserId, req.UserId)).
		Exec()
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDeleteResourcesFailed)
	}

	reply := &pb.DeleteUsersResponse{
		UserId: req.UserId,
	}

	return reply, nil
}

func (p *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	email := req.GetEmail().GetValue()

	// check email exists
	query := pi.Global().DB(ctx).
		Select(models.UserColumns...).
		From(constants.TableUser).Limit(1).
		Where(db.Eq(constants.ColumnEmail, email)).
		Where(db.Neq(constants.ColumnStatus, constants.StatusDeleted))
	if count, err := query.Count(); err == nil && count > 0 {
		return nil, gerr.New(ctx, gerr.FailedPrecondition, gerr.ErrorEmailExists, email)
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
		InsertInto(constants.TableUser).
		Record(newUser).
		Exec()
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
	}

	reply := &pb.CreateUserResponse{
		UserId: pbutil.ToProtoString(newUser.UserId),
	}

	return reply, nil
}

func (p *Server) CreatePasswordReset(ctx context.Context, req *pb.CreatePasswordResetRequest) (*pb.CreatePasswordResetResponse, error) {
	var userId = req.GetUserId().GetValue()
	var userInfo models.User

	query := pi.Global().DB(ctx).
		Select(models.UserColumns...).
		From(constants.TableUser).
		Where(db.Eq(constants.ColumnUserId, userId))
	err := query.LoadOne(&userInfo)
	if err != nil {
		if err == db.ErrNotFound {
			return nil, gerr.New(ctx, gerr.InvalidArgument, gerr.ErrorResourceNotFound, userId)
		}
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}

	err = bcrypt.CompareHashAndPassword([]byte(userInfo.Password), []byte(req.GetPassword().GetValue()))
	if err != nil {
		return nil, gerr.New(ctx, gerr.InvalidArgument, gerr.ErrorPasswordIncorrect)
	}

	var newUserPasswordReset = models.NewUserPasswordReset(userId)

	_, err = pi.Global().DB(ctx).
		InsertInto(constants.TableUserPasswordReset).
		Columns(models.UserPasswordResetColumns...).
		Record(newUserPasswordReset).
		Exec()
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}

	reply := &pb.CreatePasswordResetResponse{
		UserId:  pbutil.ToProtoString(userId),
		ResetId: pbutil.ToProtoString(newUserPasswordReset.ResetId),
	}

	return reply, nil
}

func (p *Server) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error) {
	resetId := req.GetResetId().GetValue()
	newPassword := req.GetNewPassword().GetValue()

	var resetInfo models.UserPasswordReset

	query := pi.Global().DB(ctx).
		Select(models.UserPasswordResetColumns...).
		From(constants.TableUserPasswordReset).Limit(1).
		Where(db.Eq(constants.ColumnResetId, resetId))
	err := query.LoadOne(&resetInfo)
	if err != nil {
		if err == db.ErrNotFound {
			return nil, gerr.New(ctx, gerr.InvalidArgument, gerr.ErrorResourceNotFound, resetId)
		}
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}

	_, err = pi.Global().DB(ctx).
		Update(constants.TableUser).
		Set(constants.ColumnPassword, string(hashedPass)).
		Where(db.Eq(constants.ColumnUserId, resetInfo.UserId)).
		Exec()
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}

	reply := &pb.ChangePasswordResponse{
		UserId: pbutil.ToProtoString(resetInfo.UserId),
	}

	return reply, nil
}

func (p *Server) GetPasswordReset(ctx context.Context, req *pb.GetPasswordResetRequest) (*pb.GetPasswordResetResponse, error) {
	var resetId = req.GetResetId()
	var resetInfo models.UserPasswordReset

	query := pi.Global().DB(ctx).
		Select(models.UserColumns...).
		From(constants.TableUserPasswordReset).Limit(1).
		Where(db.Eq(constants.ColumnResetId, resetId))
	err := query.LoadOne(&resetInfo)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
	}

	reply := &pb.GetPasswordResetResponse{
		ResetId: resetId,
		UserId:  resetInfo.ResetId,
	}

	return reply, nil
}
func (p *Server) ValidateUserPassword(ctx context.Context, req *pb.ValidateUserPasswordRequest) (*pb.ValidateUserPasswordResponse, error) {
	var email = req.GetEmail()
	var userInfo models.User

	query := pi.Global().DB(ctx).
		Select(models.UserColumns...).
		From(constants.TableUser).Limit(1).
		Where(db.Eq(constants.ColumnEmail, email))
	err := query.LoadOne(&userInfo)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal,
			fmt.Errorf("user(%q) not found", email),
			gerr.ErrorCreateResourcesFailed,
		)
	}

	var validated = true
	err = bcrypt.CompareHashAndPassword([]byte(userInfo.Password), []byte(req.GetPassword()))
	if err != nil {
		validated = false
	}

	reply := &pb.ValidateUserPasswordResponse{
		Validated: validated,
	}

	return reply, nil
}

func (p *Server) CreateGroup(ctx context.Context, req *pb.CreateGroupRequest) (*pb.CreateGroupResponse, error) {
	var newGroup = models.NewGroup(
		req.GetName().GetValue(),
		req.GetDescription().GetValue(),
	)

	_, err := pi.Global().DB(ctx).
		InsertInto(constants.TableGroup).
		Columns(models.GroupColumns...).
		Record(newGroup).
		Exec()
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
	}

	reply := &pb.CreateGroupResponse{
		GroupId: pbutil.ToProtoString(newGroup.GroupId),
	}

	return reply, nil
}

func (p *Server) ModifyGroup(ctx context.Context, req *pb.ModifyGroupRequest) (*pb.ModifyGroupResponse, error) {
	groupId := req.GetGroupId().GetValue()

	var attributes = manager.BuildUpdateAttributes(req,
		"name", "status", "description",
	)

	if len(attributes) == 0 {
		return &pb.ModifyGroupResponse{}, nil
	}

	_, err := pi.Global().DB(ctx).
		Update(constants.TableGroup).SetMap(attributes).
		Where(db.Eq(constants.ColumnGroupId, groupId)).
		Exec()
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDeleteResourcesFailed)
	}

	reply := &pb.ModifyGroupResponse{
		GroupId: pbutil.ToProtoString(groupId),
	}

	return reply, nil
}
func (p *Server) DeleteGroups(ctx context.Context, req *pb.DeleteGroupsRequest) (*pb.DeleteGroupsResponse, error) {
	_, err := pi.Global().DB(ctx).
		DeleteFrom(constants.TableGroup).
		Where(db.Eq(constants.ColumnGroupId, req.GroupId)).
		Exec()
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDeleteResourcesFailed)
	}

	reply := &pb.DeleteGroupsResponse{
		GroupId: req.GroupId,
	}

	return reply, nil
}

func (p *Server) JoinGroup(ctx context.Context, req *pb.JoinGroupRequest) (*pb.JoinGroupResponse, error) {
	var reply = &pb.JoinGroupResponse{}
	var lastErr error

	for _, gid := range req.GroupId {
		for _, uid := range req.UserId {
			query := pi.Global().DB(ctx).
				Select().From(constants.TableGroupMember).Limit(1).
				Where(db.Eq(constants.ColumnGroupId, gid)).
				Where(db.Eq(constants.ColumnUserId, uid))
			if count, err := query.Count(); err != nil || count > 0 {
				logger.Warn(ctx, "gid(%v)/uid(%v): %v", gid, uid, err)

				if err != nil {
					lastErr = err
				}
				continue
			}

			_, err := pi.Global().DB(ctx).
				InsertInto(constants.TableGroupMember).
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

func (p *Server) LeaveGroup(ctx context.Context, req *pb.LeaveGroupRequest) (*pb.LeaveGroupResponse, error) {
	_, err := pi.Global().DB(ctx).
		DeleteFrom(constants.TableGroupMember).
		Where(db.And(
			db.Eq(constants.ColumnGroupId, req.GroupId),
			db.Eq(constants.ColumnUserId, req.UserId),
		)).Exec()
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorUpgradeResourceFailed)
	}

	var reply = &pb.LeaveGroupResponse{
		GroupId: req.GroupId,
		UserId:  req.UserId,
	}

	return reply, nil
}
