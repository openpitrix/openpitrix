// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package account

import (
	"context"
	"time"

	pbam "openpitrix.io/iam/pkg/pb/am"
	pbim "openpitrix.io/iam/pkg/pb/im"
	"openpitrix.io/openpitrix/pkg/client/im"
	nfclient "openpitrix.io/openpitrix/pkg/client/notification"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/gerr"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/pi"
	"openpitrix.io/openpitrix/pkg/util/ctxutil"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
)

var (
	_           pb.AccountManagerServer = (*Server)(nil)
	imClient, _                         = im.NewClient()
)

const OwnerKey = "owner"

func formatUsers(pbamUsers []*pbam.UserWithRole) []*pb.User {
	var users []*pb.User
	for _, u := range pbamUsers {
		var rs []*pb.Role
		for _, r := range u.Role {
			rs = append(rs, pbRole(r))
		}
		users = append(users, &pb.User{
			UserId:      pbutil.ToProtoString(u.UserId),
			Username:    pbutil.ToProtoString(u.Username),
			Email:       pbutil.ToProtoString(u.Email),
			PhoneNumber: pbutil.ToProtoString(u.PhoneNumber),
			Description: pbutil.ToProtoString(u.Description),
			Status:      pbutil.ToProtoString(u.Status),
			CreateTime:  u.CreateTime,
			UpdateTime:  u.UpdateTime,
			StatusTime:  u.StatusTime,
			Role:        rs,
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

	res, err := amClient.DescribeUsersWithRole(ctx, &pbam.DescribeUsersWithRoleRequest{
		Limit:      int32(limit),
		Offset:     int32(offset),
		SortKey:    req.GetSortKey().GetValue(),
		Reverse:    req.GetReverse().GetValue(),
		SearchWord: req.GetSearchWord().GetValue(),

		GroupId: req.GetGroupId(),
		UserId:  req.GetUserId(),
		Status:  req.GetStatus(),
		RoleId:  req.GetRole(),
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
	res, err := imClient.ListGroups(ctx, &pbim.ListGroupsRequest{
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
	user, err := imClient.GetUser(ctx, &pbim.UserId{
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
		_, err = imClient.ModifyPassword(ctx, &pbim.Password{
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
		user.Username = req.GetUsername().GetValue()
	}
	user.UpdateTime = pbutil.ToProtoTimestamp(time.Now())

	user, err = imClient.ModifyUser(ctx, user)
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
		user, err := imClient.GetUser(ctx, &pbim.UserId{
			UserId: uid,
		})
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
		}
		user.Status = constants.StatusDeleted
		user.StatusTime = pbutil.ToProtoTimestamp(time.Now())

		user, err = imClient.ModifyUser(ctx, user)
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
		}
	}

	return &pb.DeleteUsersResponse{
		UserId: uids,
	}, nil
}

func (p *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	s := ctxutil.GetSender(ctx)
	email := req.GetEmail().GetValue()

	res, err := imClient.ListUsers(ctx, &pbim.ListUsersRequest{
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

	role := req.GetRole().GetValue()
	user, err := imClient.CreateUser(ctx, &pbim.User{
		Email:       req.GetEmail().GetValue(),
		PhoneNumber: req.GetPhoneNumber().GetValue(),
		Username:    getUsernameFromEmail(req.GetEmail().GetValue()),
		Password:    req.GetPassword().GetValue(),
		Description: req.GetDescription().GetValue(),
		Status:      constants.StatusActive,
	})
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}

	_, err = amClient.BindUserRole(ctx, &pbam.BindUserRoleRequest{
		UserId: []string{user.UserId},
		RoleId: []string{role},
	})
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}

	var emailNotifications []*models.EmailNotification

	if role == constants.RoleIsv {
		emailNotifications = append(emailNotifications, &models.EmailNotification{
			Title:       constants.AdminInviteIsvNotifyTitle.GetDefaultMessage(),
			Content:     constants.AdminInviteIsvNotifyContent.GetDefaultMessage(user.Username, user.Email, req.GetPassword().GetValue()),
			Owner:       s.UserId,
			ContentType: constants.NfContentTypeInvite,
			Addresses:   []string{user.Email},
		})
	} else {
		emailNotifications = append(emailNotifications, &models.EmailNotification{
			Title:       constants.AdminInviteUserNotifyTitle.GetDefaultMessage(),
			Content:     constants.AdminInviteUserNotifyContent.GetDefaultMessage(user.Username, user.Email, req.GetPassword().GetValue()),
			Owner:       s.UserId,
			ContentType: constants.NfContentTypeInvite,
			Addresses:   []string{user.Email},
		})
	}

	nfclient.SendEmailNotification(ctx, emailNotifications)

	reply := &pb.CreateUserResponse{
		UserId: pbutil.ToProtoString(user.UserId),
	}

	return reply, nil
}

func (p *Server) IsvCreateUser(ctx context.Context, req *pb.IsvCreateUserRequest) (*pb.IsvCreateUserResponse, error) {
	s := ctxutil.GetSender(ctx)
	email := req.GetEmail().GetValue()

	res, err := imClient.ListUsers(ctx, &pbim.ListUsersRequest{
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

	user, err := imClient.CreateUser(ctx, &pbim.User{
		Email:       req.GetEmail().GetValue(),
		PhoneNumber: req.GetPhoneNumber().GetValue(),
		Username:    getUsernameFromEmail(req.GetEmail().GetValue()),
		Password:    req.GetPassword().GetValue(),
		Description: req.GetDescription().GetValue(),
		Status:      constants.StatusActive,
	})
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}
	_, err = amClient.BindUserRole(ctx, &pbam.BindUserRoleRequest{
		UserId: []string{user.UserId},
		RoleId: []string{constants.RoleUser},
	})
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}

	var emailNotifications []*models.EmailNotification
	listUsersResponse, err := imClient.ListUsers(ctx, &pbim.ListUsersRequest{
		UserId: []string{s.UserId},
	})
	if err != nil || len(listUsersResponse.User) != 1 {
		logger.Error(ctx, "Failed to describe users [%s]: %+v", s.UserId, err)
	} else {
		emailNotifications = append(emailNotifications, &models.EmailNotification{
			Title:       constants.IsvInviteMemberNotifyTitle.GetDefaultMessage(listUsersResponse.User[0].Username),
			Content:     constants.IsvInviteMemberNotifyContent.GetDefaultMessage(user.Username, user.Email, req.GetPassword().GetValue()),
			Owner:       s.UserId,
			ContentType: constants.NfContentTypeInvite,
			Addresses:   []string{user.Email},
		})
	}
	nfclient.SendEmailNotification(ctx, emailNotifications)

	reply := &pb.IsvCreateUserResponse{
		UserId: pbutil.ToProtoString(user.UserId),
	}

	return reply, nil
}

func (p *Server) CreatePasswordReset(ctx context.Context, req *pb.CreatePasswordResetRequest) (*pb.CreatePasswordResetResponse, error) {
	uid := req.GetUserId().GetValue()
	password := req.GetPassword().GetValue()

	_, err := imClient.GetUser(ctx, &pbim.UserId{
		UserId: uid,
	})
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}
	b, err := imClient.ComparePassword(ctx, &pbim.Password{
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

	_, err = imClient.ModifyPassword(ctx, &pbim.Password{
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
	res, err := imClient.ListUsers(ctx, &pbim.ListUsersRequest{
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
	b, err := imClient.ComparePassword(ctx, &pbim.Password{
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
	group, err := imClient.CreateGroup(ctx, &pbim.Group{
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

	group, err := imClient.GetGroup(ctx, &pbim.GroupId{
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

	group, err = imClient.ModifyGroup(ctx, group)
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
		group, err := imClient.GetGroup(ctx, &pbim.GroupId{
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

		group, err = imClient.ModifyGroup(ctx, group)
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
		}
	}

	return &pb.DeleteGroupsResponse{
		GroupId: gids,
	}, nil
}

func (p *Server) JoinGroup(ctx context.Context, req *pb.JoinGroupRequest) (*pb.JoinGroupResponse, error) {
	_, err := imClient.JoinGroup(ctx, &pbim.JoinGroupRequest{
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
	_, err := imClient.LeaveGroup(ctx, &pbim.LeaveGroupRequest{
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
	res, err := imClient.ListGroups(ctx, &pbim.ListGroupsRequest{
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
