// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package account

import (
	"context"
	"time"

	pbim "openpitrix.io/iam/pkg/pb/im"
	clientiam2 "openpitrix.io/openpitrix/pkg/client/iam2"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/gerr"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/pi"
	"openpitrix.io/openpitrix/pkg/util/ctxutil"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
)

var (
	_         pb.AccountManagerServer = (*Server)(nil)
	client, _                         = clientiam2.NewClient()
)

const OwnerKey = "owner"

func formatUsers(pbimUsers []*pbim.User) []*pb.User {
	var users []*pb.User
	for _, u := range pbimUsers {
		var role = ""
		if r, ok := u.Extra["role"]; ok {
			role = r
		}
		users = append(users, &pb.User{
			UserId:      pbutil.ToProtoString(u.UserId),
			Username:    pbutil.ToProtoString(u.UserName),
			Email:       pbutil.ToProtoString(u.Email),
			PhoneNumber: pbutil.ToProtoString(u.PhoneNumber),
			Description: pbutil.ToProtoString(u.Description),
			Status:      pbutil.ToProtoString(u.Status),
			CreateTime:  u.CreateTime,
			UpdateTime:  u.UpdateTime,
			StatusTime:  u.StatusTime,
			Role:        pbutil.ToProtoString(role),
			GroupId:     u.GroupId,
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

		GroupId: req.GetGroupId(),
		UserId:  req.GetUserId(),
		Status:  req.GetStatus(),
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

		GroupId: req.GetGroupId(),
		UserId:  req.GetUserId(),
		Status:  req.GetStatus(),
	})
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}
	var groups []*pb.Group
	for _, u := range res.Group {
		groups = append(groups, &pb.Group{
			ParentGroupId: pbutil.ToProtoString(u.ParentGroupId),
			GroupId:       pbutil.ToProtoString(u.GroupId),
			GroupPath:     pbutil.ToProtoString(u.GroupPath),
			Name:          pbutil.ToProtoString(u.GroupName),
			Description:   pbutil.ToProtoString(u.Description),
			Status:        pbutil.ToProtoString(u.Status),
			CreateTime:    u.CreateTime,
			UpdateTime:    u.UpdateTime,
			StatusTime:    u.StatusTime,
			UserId:        u.UserId,
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
		UserId: userId,
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
			UserId:   req.GetUserId().GetValue(),
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
		user.UserName = req.GetUsername().GetValue()
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
			UserId: uid,
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
		PhoneNumber: req.GetPhoneNumber().GetValue(),
		UserName:    getUsernameFromEmail(req.GetEmail().GetValue()),
		Password:    req.GetPassword().GetValue(),
		Description: req.GetDescription().GetValue(),
		Status:      constants.StatusActive,
		Extra:       map[string]string{"role": req.GetRole().GetValue()},
	})
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}

	reply := &pb.CreateUserResponse{
		UserId: pbutil.ToProtoString(user.UserId),
	}

	return reply, nil
}

func (p *Server) IsvCreateUser(ctx context.Context, req *pb.IsvCreateUserRequest) (*pb.IsvCreateUserResponse, error) {
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
		PhoneNumber: req.GetPhoneNumber().GetValue(),
		UserName:    getUsernameFromEmail(req.GetEmail().GetValue()),
		Password:    req.GetPassword().GetValue(),
		Description: req.GetDescription().GetValue(),
		Status:      constants.StatusActive,
		Extra:       map[string]string{"role": req.GetRole().GetValue()},
	})
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}

	reply := &pb.IsvCreateUserResponse{
		UserId: pbutil.ToProtoString(user.UserId),
	}

	return reply, nil
}

func (p *Server) CreatePasswordReset(ctx context.Context, req *pb.CreatePasswordResetRequest) (*pb.CreatePasswordResetResponse, error) {
	uid := req.GetUserId().GetValue()
	password := req.GetPassword().GetValue()

	_, err := client.GetUser(ctx, &pbim.UserId{
		UserId: uid,
	})
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}
	b, err := client.ComparePassword(ctx, &pbim.Password{
		UserId:   uid,
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
		UserId:   resetInfo.UserId,
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
		UserId:   user.UserId,
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
	s := ctxutil.GetSender(ctx)
	group, err := client.CreateGroup(ctx, &pbim.Group{
		ParentGroupId: req.GetParentGroupId().GetValue(),
		GroupName:     req.GetName().GetValue(),
		Description:   req.GetDescription().GetValue(),
		Status:        constants.StatusActive,
		Extra: map[string]string{
			OwnerKey: s.UserId,
		},
	})
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}

	reply := &pb.CreateGroupResponse{
		GroupId: pbutil.ToProtoString(group.GroupId),
	}

	return reply, nil
}

func (p *Server) ModifyGroup(ctx context.Context, req *pb.ModifyGroupRequest) (*pb.ModifyGroupResponse, error) {
	gid := req.GetGroupId().GetValue()

	group, err := client.GetGroup(ctx, &pbim.GroupId{
		GroupId: gid,
	})
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}

	if req.GetDescription() != nil {
		group.Description = req.GetDescription().GetValue()
	}
	if req.GetName() != nil {
		group.GroupName = req.GetName().GetValue()
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
			GroupId: gid,
		})
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
		}
		if len(group.UserId) > 0 {
			return nil, gerr.NewWithDetail(ctx, gerr.FailedPrecondition, err, gerr.ErrorGroupHadMembers)
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
		GroupId: req.GetGroupId(),
		UserId:  req.GetUserId(),
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
		GroupId: req.GetGroupId(),
		UserId:  req.GetUserId(),
	})
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}

	return &pb.LeaveGroupResponse{
		GroupId: req.GetGroupId(),
		UserId:  req.GetUserId(),
	}, nil
}

func (p *Server) GetUserGroupOwner(ctx context.Context, req *pb.GetUserGroupOwnerRequest) (*pb.GetUserGroupOwnerResponse, error) {
	res, err := client.ListGroups(ctx, &pbim.ListGroupsRequest{
		Limit:  1,
		UserId: []string{req.GetUserId()},
		Status: []string{constants.StatusActive},
	})
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}
	reply := &pb.GetUserGroupOwnerResponse{
		UserId: req.GetUserId(),
	}
	if res.Total == 0 {
		return reply, nil
	}
	owner, ok := res.Group[0].Extra[OwnerKey]
	if !ok {
		return reply, nil
	}

	reply.Owner = owner

	return reply, nil
}
