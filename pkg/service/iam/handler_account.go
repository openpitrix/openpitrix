// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package iam

import (
	"context"
	"time"

	"openpitrix.io/iam/pkg/pb/im"
	"openpitrix.io/openpitrix/pkg/client/iam2"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/gerr"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/pi"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
)

var (
	_         pb.AccountManagerServer = (*Server)(nil)
	client, _                         = clientiam2.NewClient()
)

func formatUsers(pbimUsers []*pbim.User) []*pb.User {
	var users []*pb.User
	for _, u := range pbimUsers {
		var role = ""
		if r, ok := u.Extra["role"]; ok {
			role = r
		}
		users = append(users, &pb.User{
			UserId:      pbutil.ToProtoString(u.Uid),
			Username:    pbutil.ToProtoString(u.Name),
			Email:       pbutil.ToProtoString(u.Email),
			Description: pbutil.ToProtoString(u.Description),
			Status:      pbutil.ToProtoString(u.Status),
			CreateTime:  u.CreateTime,
			UpdateTime:  u.UpdateTime,
			StatusTime:  u.StatusTime,
			Role:        pbutil.ToProtoString(role),
		})
	}
	return users
}

func (p *Server) DescribeUsers(ctx context.Context, req *pb.DescribeUsersRequest) (*pb.DescribeUsersResponse, error) {
	var (
		offset = pbutil.GetOffsetFromRequest(req)
		limit  = pbutil.GetLimitFromRequest(req)
	)

	res, err := client.ListUsers(ctx, &pbim.ListUsersRequest{
		Limit:      int32(limit),
		Offset:     int32(offset),
		SortKey:    req.GetSortKey().GetValue(),
		Reverse:    req.GetReverse().GetValue(),
		SearchWord: req.GetSearchWord().GetValue(),

		Gid:    req.GetGroupId(),
		Uid:    req.GetUserId(),
		Status: req.GetStatus(),
	})
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}

	reply := &pb.DescribeUsersResponse{
		UserSet:    formatUsers(res.User),
		TotalCount: uint32(res.GetTotal()),
	}
	return reply, nil
}

func (p *Server) DescribeGroups(ctx context.Context, req *pb.DescribeGroupsRequest) (*pb.DescribeGroupsResponse, error) {
	var (
		offset = pbutil.GetOffsetFromRequest(req)
		limit  = pbutil.GetLimitFromRequest(req)
	)
	res, err := client.ListGroups(ctx, &pbim.ListGroupsRequest{
		Limit:      int32(limit),
		Offset:     int32(offset),
		SortKey:    req.GetSortKey().GetValue(),
		Reverse:    req.GetReverse().GetValue(),
		SearchWord: req.GetSearchWord().GetValue(),

		Gid:    req.GetGroupId(),
		Uid:    req.GetUserId(),
		Status: req.GetStatus(),
	})
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}
	var groups []*pb.Group
	for _, u := range res.Group {
		groups = append(groups, &pb.Group{
			GroupId:     pbutil.ToProtoString(u.Gid),
			Name:        pbutil.ToProtoString(u.Name),
			Description: pbutil.ToProtoString(u.Description),
			Status:      pbutil.ToProtoString(u.Status),
			CreateTime:  u.CreateTime,
			UpdateTime:  u.UpdateTime,
			StatusTime:  u.StatusTime,
		})
	}

	reply := &pb.DescribeGroupsResponse{
		GroupSet:   groups,
		TotalCount: uint32(res.GetTotal()),
	}
	return reply, nil
}

func getUser(ctx context.Context, userId string) (*pbim.User, error) {
	user, err := client.GetUser(ctx, &pbim.UserId{
		Uid: userId,
	})
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}
	return user, err
}

func (p *Server) ModifyUser(ctx context.Context, req *pb.ModifyUserRequest) (*pb.ModifyUserResponse, error) {
	uid := req.GetUserId().GetValue()

	user, err := getUser(ctx, uid)
	if err != nil {
		return nil, err
	}

	password := req.GetPassword().GetValue()
	if password != "" {
		_, err = client.ModifyPassword(ctx, &pbim.Password{
			Uid:      req.GetUserId().GetValue(),
			Password: password,
		})
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
		}
	}
	if req.GetDescription() != nil {
		user.Description = req.GetDescription().GetValue()
	}
	if req.GetEmail() != nil {
		user.Email = req.GetEmail().GetValue()
	}
	if req.GetUsername() != nil {
		user.Name = req.GetUsername().GetValue()
	}
	user.UpdateTime = pbutil.ToProtoTimestamp(time.Now())

	user, err = client.ModifyUser(ctx, user)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}

	reply := &pb.ModifyUserResponse{
		UserId: req.UserId,
	}

	return reply, nil
}

func (p *Server) DeleteUsers(ctx context.Context, req *pb.DeleteUsersRequest) (*pb.DeleteUsersResponse, error) {
	uids := req.GetUserId()

	for _, uid := range uids {
		user, err := client.GetUser(ctx, &pbim.UserId{
			Uid: uid,
		})
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
		}
		user.Status = constants.StatusDeleted
		user.StatusTime = pbutil.ToProtoTimestamp(time.Now())

		user, err = client.ModifyUser(ctx, user)
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
		}
	}

	return &pb.DeleteUsersResponse{
		UserId: uids,
	}, nil
}

func (p *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	email := req.GetEmail().GetValue()

	res, err := client.ListUsers(ctx, &pbim.ListUsersRequest{
		Limit:  1,
		Status: []string{constants.StatusActive},
		Email:  []string{email},
	})
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}
	if res.Total > 0 {
		return nil, gerr.New(ctx, gerr.FailedPrecondition, gerr.ErrorEmailExists, email)
	}

	user, err := client.CreateUser(ctx, &pbim.User{
		Email:       req.GetEmail().GetValue(),
		Name:        getUsernameFromEmail(req.GetEmail().GetValue()),
		Password:    req.GetPassword().GetValue(),
		Description: req.GetDescription().GetValue(),
		Status:      constants.StatusActive,
		Extra:       map[string]string{"role": req.GetRole().GetValue()},
	})
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}

	reply := &pb.CreateUserResponse{
		UserId: pbutil.ToProtoString(user.Uid),
	}

	return reply, nil
}

func (p *Server) CreatePasswordReset(ctx context.Context, req *pb.CreatePasswordResetRequest) (*pb.CreatePasswordResetResponse, error) {
	uid := req.GetUserId().GetValue()
	password := req.GetPassword().GetValue()

	_, err := client.GetUser(ctx, &pbim.UserId{
		Uid: uid,
	})
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}
	b, err := client.ComparePassword(ctx, &pbim.Password{
		Uid:      uid,
		Password: password,
	})
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}
	if !b.Value {
		return nil, gerr.New(ctx, gerr.PermissionDenied, gerr.ErrorPasswordIncorrect)
	}

	var newUserPasswordReset = models.NewUserPasswordReset(uid)

	_, err = pi.Global().DB(ctx).
		InsertInto(constants.TableUserPasswordReset).
		Columns(models.UserPasswordResetColumns...).
		Record(newUserPasswordReset).
		Exec()
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}

	reply := &pb.CreatePasswordResetResponse{
		UserId:  pbutil.ToProtoString(uid),
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
		From(constants.TableUserPasswordReset).
		Limit(1).
		Where(db.Eq(constants.ColumnStatus, constants.StatusActive)).
		Where(db.Eq(constants.ColumnResetId, resetId))
	err := query.LoadOne(&resetInfo)
	if err != nil {
		if err == db.ErrNotFound {
			return nil, gerr.New(ctx, gerr.InvalidArgument, gerr.ErrorResourceNotFound, resetId)
		}
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}

	_, err = client.ModifyPassword(ctx, &pbim.Password{
		Uid:      resetInfo.UserId,
		Password: newPassword,
	})
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}

	_, err = pi.Global().DB(ctx).
		Update(constants.TableUserPasswordReset).
		Set(constants.ColumnStatus, constants.StatusUsed).
		Where(db.Eq(constants.ColumnResetId, resetId)).
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
		Select(models.UserPasswordResetColumns...).
		From(constants.TableUserPasswordReset).
		Limit(1).
		Where(db.Eq(constants.ColumnStatus, constants.StatusActive)).
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

func validateUserPassword(ctx context.Context, email, password string) (*pbim.User, bool, error) {
	res, err := client.ListUsers(ctx, &pbim.ListUsersRequest{
		Limit:  1,
		Status: []string{constants.StatusActive},
		Email:  []string{email},
	})
	if err != nil {
		return nil, false, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}
	if res.Total == 0 {
		return nil, false, gerr.New(ctx, gerr.FailedPrecondition, gerr.ErrorEmailPasswordNotMatched)
	}

	user := res.User[0]
	b, err := client.ComparePassword(ctx, &pbim.Password{
		Uid:      user.Uid,
		Password: password,
	})
	if err != nil {
		return nil, false, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}
	return user, b.Value, nil
}

func (p *Server) ValidateUserPassword(ctx context.Context, req *pb.ValidateUserPasswordRequest) (*pb.ValidateUserPasswordResponse, error) {
	_, b, err := validateUserPassword(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, err
	}
	return &pb.ValidateUserPasswordResponse{
		Validated: b,
	}, nil
}

func (p *Server) CreateGroup(ctx context.Context, req *pb.CreateGroupRequest) (*pb.CreateGroupResponse, error) {
	group, err := client.CreateGroup(ctx, &pbim.Group{
		Name:        req.GetName().GetValue(),
		Description: req.GetDescription().GetValue(),
		Status:      constants.StatusActive,
	})
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}

	reply := &pb.CreateGroupResponse{
		GroupId: pbutil.ToProtoString(group.Gid),
	}

	return reply, nil
}

func (p *Server) ModifyGroup(ctx context.Context, req *pb.ModifyGroupRequest) (*pb.ModifyGroupResponse, error) {
	gid := req.GetGroupId().GetValue()

	group, err := client.GetGroup(ctx, &pbim.GroupId{
		Gid: gid,
	})
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}

	if req.GetDescription() != nil {
		group.Description = req.GetDescription().GetValue()
	}
	if req.GetName() != nil {
		group.Name = req.GetName().GetValue()
	}
	group.UpdateTime = pbutil.ToProtoTimestamp(time.Now())

	group, err = client.ModifyGroup(ctx, group)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}

	reply := &pb.ModifyGroupResponse{
		GroupId: req.GetGroupId(),
	}

	return reply, nil
}
func (p *Server) DeleteGroups(ctx context.Context, req *pb.DeleteGroupsRequest) (*pb.DeleteGroupsResponse, error) {
	gids := req.GetGroupId()

	for _, gid := range gids {
		group, err := client.GetGroup(ctx, &pbim.GroupId{
			Gid: gid,
		})
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
		}
		group.Status = constants.StatusDeleted
		group.StatusTime = pbutil.ToProtoTimestamp(time.Now())

		group, err = client.ModifyGroup(ctx, group)
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
		}
	}

	return &pb.DeleteGroupsResponse{
		GroupId: gids,
	}, nil
}

func (p *Server) JoinGroup(ctx context.Context, req *pb.JoinGroupRequest) (*pb.JoinGroupResponse, error) {
	_, err := client.JoinGroup(ctx, &pbim.JoinGroupRequest{
		Gid: req.GetGroupId(),
		Uid: req.GetUserId(),
	})
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}

	return &pb.JoinGroupResponse{
		GroupId: req.GetGroupId(),
		UserId:  req.GetUserId(),
	}, nil
}

func (p *Server) LeaveGroup(ctx context.Context, req *pb.LeaveGroupRequest) (*pb.LeaveGroupResponse, error) {
	_, err := client.LeaveGroup(ctx, &pbim.LeaveGroupRequest{
		Gid: req.GetGroupId(),
		Uid: req.GetUserId(),
	})
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}

	return &pb.LeaveGroupResponse{
		GroupId: req.GetGroupId(),
		UserId:  req.GetUserId(),
	}, nil
}
