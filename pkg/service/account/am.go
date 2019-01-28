// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package account

import (
	"context"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/util/ctxutil"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
)

func (p *Server) Checker(ctx context.Context, req interface{}) error {
	switch r := req.(type) {
	case *pb.DescribeUsersRequest:
		return manager.NewChecker(ctx, r).
			Required().
			Exec()
	case *pb.ModifyUserRequest:
		return manager.NewChecker(ctx, r).
			Required("user_id").
			StringChosen("role", constants.AllRoles).
			Exec()
	case *pb.DeleteUsersRequest:
		return manager.NewChecker(ctx, r).
			Role(constants.AllAdminRoles).
			Required("user_id").
			Exec()
	case *pb.CreateUserRequest:
		return manager.NewChecker(ctx, r).
			Role(constants.AllAdminRoles).
			Required("email", "password").
			StringChosen("role", constants.AllRoles).
			Exec()
	case *pb.IsvCreateUserRequest:
		return manager.NewChecker(ctx, r).
			Role(constants.AllIsvRoles).
			Required("email", "password").
			StringChosen("role", constants.AllDeveloperRoles).
			Exec()
	case *pb.CreatePasswordResetRequest:
		return manager.NewChecker(ctx, r).
			Required("user_id", "password").
			Exec()
	case *pb.ChangePasswordRequest:
		return manager.NewChecker(ctx, r).
			Required("new_password", "reset_id").
			Exec()
	case *pb.GetPasswordResetRequest:
		return manager.NewChecker(ctx, r).
			Role(constants.AllAdminRoles).
			Required("reset_id").
			Exec()
	case *pb.ValidateUserPasswordRequest:
		return manager.NewChecker(ctx, r).
			Required("email", "password").
			Exec()
	case *pb.DescribeGroupsRequest:
		return manager.NewChecker(ctx, r).
			Required().
			Exec()
	case *pb.CreateGroupRequest:
		return manager.NewChecker(ctx, r).
			Role(constants.AllAdminRoles).
			Required("name").
			Exec()
	case *pb.ModifyGroupRequest:
		return manager.NewChecker(ctx, r).
			Role(constants.AllAdminRoles).
			Required("group_id").
			Exec()
	case *pb.DeleteGroupsRequest:
		return manager.NewChecker(ctx, r).
			Role(constants.AllAdminRoles).
			Required("group_id").
			Exec()
	case *pb.JoinGroupRequest:
		return manager.NewChecker(ctx, r).
			Role(constants.AllAdminRoles).
			Required("group_id", "user_id").
			Exec()
	case *pb.LeaveGroupRequest:
		return manager.NewChecker(ctx, r).
			Role(constants.AllAdminRoles).
			Required("group_id", "user_id").
			Exec()
	case *pb.CreateClientRequest:
		return manager.NewChecker(ctx, r).
			Role(constants.AllAdminRoles).
			Required("user_id").
			Exec()
	case *pb.TokenRequest:
		return manager.NewChecker(ctx, r).
			Required("grant_type", "client_id", "client_secret").
			StringChosen("grant_type", constants.GrantTypeTokens).
			Exec()
	}

	logger.Warn(ctx, "checker unknown type: %T", req)
	return nil
}

func (p *Server) Builder(ctx context.Context, req interface{}) interface{} {
	sender := ctxutil.GetSender(ctx)
	switch r := req.(type) {
	case *pb.CreatePasswordResetRequest:
		if !sender.IsGlobalAdmin() {
			r.UserId = pbutil.ToProtoString(sender.UserId)
		}
		return r
	case *pb.ModifyUserRequest:
		if !sender.IsGlobalAdmin() {
			r.UserId = pbutil.ToProtoString(sender.UserId)
			r.Role = nil
			r.Email = nil
			r.Password = nil
		}
		return r
	}
	return req
}
