// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package iam

import (
	"context"
	"fmt"

	"openpitrix.io/openpitrix/pkg/pb/iam"
)

var (
	_ pbiam.AccountManagerServer = (*Server)(nil)
)

func (p *Server) DescribeUsers(context.Context, *pbiam.DescribeUsersRequest) (*pbiam.DescribeUsersResponse, error) {
	return nil, fmt.Errorf("TODO")
}
func (p *Server) ModifyUser(context.Context, *pbiam.ModifyUserRequest) (*pbiam.ModifyUserResponse, error) {
	return nil, fmt.Errorf("TODO")
}
func (p *Server) DeleteUsers(context.Context, *pbiam.DeleteUsersRequest) (*pbiam.DeleteUsersResponse, error) {
	return nil, fmt.Errorf("TODO")
}
func (p *Server) InviteUsers(context.Context, *pbiam.InviteUsersRequest) (*pbiam.InviteUsersResponse, error) {
	return nil, fmt.Errorf("TODO")
}
func (p *Server) ChangePassword(context.Context, *pbiam.ChangePasswordRequest) (*pbiam.ChangePasswordResponse, error) {
	return nil, fmt.Errorf("TODO")
}
func (p *Server) CreatePasswordReset(context.Context, *pbiam.CreatePasswordResetRequest) (*pbiam.CreatePasswordResetResponse, error) {
	return nil, fmt.Errorf("TODO")
}
func (p *Server) CreateUser(context.Context, *pbiam.CreateUserRequest) (*pbiam.CreateUserResponse, error) {
	return nil, fmt.Errorf("TODO")
}
func (p *Server) GetPasswordReset(context.Context, *pbiam.GetPasswordResetRequest) (*pbiam.GetPasswordResetResponse, error) {
	return nil, fmt.Errorf("TODO")
}
func (p *Server) ValidateUserPassword(context.Context, *pbiam.ValidateUserPasswordRequest) (*pbiam.ValidateUserPasswordResponse, error) {
	return nil, fmt.Errorf("TODO")
}
func (p *Server) CreateGroup(context.Context, *pbiam.CreateGroupRequest) (*pbiam.CreateGroupResponse, error) {
	return nil, fmt.Errorf("TODO")
}
func (p *Server) DescribeGroups(context.Context, *pbiam.DescribeGroupsRequest) (*pbiam.DescribeUsersResponse, error) {
	return nil, fmt.Errorf("TODO")
}
func (p *Server) ModifyGroup(context.Context, *pbiam.ModifyGroupRequest) (*pbiam.ModifyGroupResponse, error) {
	return nil, fmt.Errorf("TODO")
}
func (p *Server) DeleteGroups(context.Context, *pbiam.DeleteGroupsRequest) (*pbiam.DeleteGroupsResponse, error) {
	return nil, fmt.Errorf("TODO")
}
func (p *Server) JoinGroup(context.Context, *pbiam.JoinGroupRequest) (*pbiam.JoinGroupResponse, error) {
	return nil, fmt.Errorf("TODO")
}
func (p *Server) LeaveGroup(context.Context, *pbiam.LeaveGroupRequest) (*pbiam.LeaveGroupResponse, error) {
	return nil, fmt.Errorf("TODO")
}
