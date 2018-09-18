// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package iam

import (
	"context"

	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/pb/iam"
)

func (p *Server) Checker(ctx context.Context, req interface{}) error {
	switch r := req.(type) {
	case *pbiam.DescribeUsersRequest:
		return manager.NewChecker(ctx, r).
			Required().
			Exec()
	case *pbiam.ModifyUserRequest:
		return manager.NewChecker(ctx, r).
			Required("user_id").
			Exec()
	case *pbiam.DeleteUsersRequest:
		return manager.NewChecker(ctx, r).
			Required("user_id").
			Exec()
	case *pbiam.CreateUserRequest:
		return manager.NewChecker(ctx, r).
			Required("email", "password").
			Exec()
	case *pbiam.CreatePasswordResetRequest:
		return manager.NewChecker(ctx, r).
			Required("user_id", "password").
			Exec()
	case *pbiam.ChangePasswordRequest:
		return manager.NewChecker(ctx, r).
			Required("new_password", "reset_id").
			Exec()
	case *pbiam.GetPasswordResetRequest:
		return manager.NewChecker(ctx, r).
			Required("reset_id").
			Exec()
	case *pbiam.ValidateUserPasswordRequest:
		return manager.NewChecker(ctx, r).
			Required("email", "password").
			Exec()
	case *pbiam.DescribeGroupsRequest:
		return manager.NewChecker(ctx, r).
			Required().
			Exec()
	case *pbiam.CreateGroupRequest:
		return manager.NewChecker(ctx, r).
			Required("name").
			Exec()
	case *pbiam.ModifyGroupRequest:
		return manager.NewChecker(ctx, r).
			Required("group_id").
			Exec()
	case *pbiam.DeleteGroupsRequest:
		return manager.NewChecker(ctx, r).
			Required("group_id").
			Exec()
	case *pbiam.JoinGroupRequest:
		return manager.NewChecker(ctx, r).
			Required("group_id", "user_id").
			Exec()
	case *pbiam.LeaveGroupRequest:
		return manager.NewChecker(ctx, r).
			Required("group_id", "user_id").
			Exec()
	}

	logger.Warn(ctx, "checker unknown type: %T", req)
	return nil
}
