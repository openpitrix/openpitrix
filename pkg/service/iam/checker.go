// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package iam

import (
	"context"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/pb"
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
			StringChosen("role", constants.SupportRoles).
			Exec()
	case *pb.DeleteUsersRequest:
		return manager.NewChecker(ctx, r).
			Required("user_id").
			Exec()
	case *pb.CreateUserRequest:
		return manager.NewChecker(ctx, r).
			Required("email", "password").
			StringChosen("role", constants.SupportRoles).
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
			Required("name").
			Exec()
	case *pb.ModifyGroupRequest:
		return manager.NewChecker(ctx, r).
			Required("group_id").
			Exec()
	case *pb.DeleteGroupsRequest:
		return manager.NewChecker(ctx, r).
			Required("group_id").
			Exec()
	case *pb.JoinGroupRequest:
		return manager.NewChecker(ctx, r).
			Required("group_id", "user_id").
			Exec()
	case *pb.LeaveGroupRequest:
		return manager.NewChecker(ctx, r).
			Required("group_id", "user_id").
			Exec()
	case *pb.CreateClientRequest:
		return manager.NewChecker(ctx, r).
			Required("user_id").
			Exec()
	case *pb.AuthRequest:
		return manager.NewChecker(ctx, r).
			Required("grant_type", "client_id", "client_secret").
			StringChosen("grant_type", constants.GrantTypeAuths).
			Exec()
	case *pb.TokenRequest:
		return manager.NewChecker(ctx, r).
			Required("grant_type", "client_id", "client_secret", "refresh_token").
			StringChosen("grant_type", constants.GrantTypeTokens).
			Exec()
	}

	logger.Warn(ctx, "checker unknown type: %T", req)
	return nil
}
