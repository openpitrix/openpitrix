// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package account

import (
	"context"
	"strings"

	pbim "kubesphere.io/im/pkg/pb"

	pbam "openpitrix.io/iam/pkg/pb"
	"openpitrix.io/openpitrix/pkg/client/iam/im"
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
	"openpitrix.io/openpitrix/pkg/util/stringutil"
)

var (
	_           pb.AccountManagerServer = (*Server)(nil)
	imClient, _                         = im.NewClient()
)

const OwnerKey = "owner"
const OwnerPathKey = "owner_path"

func getRoleUserIds(ctx context.Context, roleIds, userIds []string) ([]string, error) {
	var roleUserIds []string
	res, err := amClient.DescribeRolesWithUser(ctx, &pbam.DescribeRolesRequest{
		RoleId: roleIds,
	})
	if err != nil {
		return nil, err
	}
	for _, role := range res.RoleSet {
		roleUserIds = append(roleUserIds, role.UserIdSet...)
	}

	var retUserIds []string
	if len(userIds) == 0 {
		retUserIds = roleUserIds
	} else {
		for _, userId := range roleUserIds {
			if stringutil.StringIn(userId, userIds) {
				retUserIds = append(retUserIds, userId)
			}
		}
	}

	return retUserIds, nil
}

func getUser(ctx context.Context, userId string) (*pbim.User, error) {
	getUserResponse, err := imClient.GetUser(ctx, &pbim.GetUserRequest{
		UserId: userId,
	})
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}
	return getUserResponse.User, err
}

func createAndJoinRootGroup(ctx context.Context, userId string) error {
	rep, err := imClient.CreateGroup(ctx, &pbim.CreateGroupRequest{
		ParentGroupId: "",
		GroupName:     "root",
	})
	if err != nil {
		return err
	}

	groupId := rep.GroupId
	_, err = imClient.ModifyGroup(ctx, &pbim.ModifyGroupRequest{
		GroupId: groupId,
		Extra: map[string]string{
			OwnerKey:     userId,
			OwnerPathKey: groupId + ":" + userId,
		},
	})
	if err != nil {
		return err
	}

	_, err = imClient.JoinGroup(ctx, &pbim.JoinGroupRequest{
		GroupId: []string{groupId},
		UserId:  []string{userId},
	})

	if err != nil {
		return err
	}
	return nil
}

func getSystemUserId(ctx context.Context) (string, error) {
	var userId string
	getRoleWithUserResponse, err := amClient.GetRoleWithUser(ctx, &pbam.GetRoleRequest{
		RoleId: constants.RoleGlobalAdmin,
	})
	if err != nil {
		return userId, err
	}
	userIds := getRoleWithUserResponse.Role.UserIdSet
	if len(userIds) == 0 {
		logger.Error(ctx, "There is no global admin user")
		return userId, gerr.New(ctx, gerr.Internal, gerr.ErrorInternalError)
	}
	listUsersResponse, err := imClient.ListUsers(ctx, &pbim.ListUsersRequest{
		UserId: userIds,
		Status: []string{constants.StatusActive},
	})
	if err != nil {
		return userId, err
	}
	if len(listUsersResponse.UserSet) == 0 {
		logger.Error(ctx, "There is no active global admin user")
		return userId, gerr.New(ctx, gerr.Internal, gerr.ErrorInternalError)
	}
	userId = listUsersResponse.UserSet[0].UserId
	return userId, nil
}

func getRootGroupId(ctx context.Context, userId string) (string, error) {
	var rootGroupId string

	if userId == constants.UserSystem {
		var err error
		userId, err = getSystemUserId(ctx)
		if err != nil {
			return rootGroupId, err
		}
	}

	getUserWithGroupRes, err := imClient.GetUserWithGroup(ctx, &pbim.GetUserRequest{
		UserId: userId,
	})
	if err != nil {
		logger.Error(ctx, "Get user [%s] failed: %+v", userId, err)
		return rootGroupId, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}

	// root group is same
	if len(getUserWithGroupRes.User.GroupSet) == 0 || len(getUserWithGroupRes.User.GroupSet[0].GroupPath) == 0 {
		logger.Error(ctx, "Failed to get root group for user [%s]", userId)
		return rootGroupId, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}
	groupPath := getUserWithGroupRes.User.GroupSet[0].GroupPath
	rootGroupId = strings.Split(groupPath, ".")[0]
	return rootGroupId, nil
}

func getUserPortal(ctx context.Context, userId string) (string, error) {
	if userId == constants.UserSystem {
		return constants.PortalGlobalAdmin, nil
	}
	var portal string
	response, err := amClient.DescribeRoles(ctx, &pbam.DescribeRolesRequest{
		UserId: []string{userId},
	})
	if err != nil {
		return portal, err
	}
	if len(response.RoleSet) == 0 {
		logger.Error(ctx, "Failed to get role for user [%s]", userId)
		return portal, gerr.New(ctx, gerr.Internal, gerr.ErrorInternalError)
	}
	portal = response.RoleSet[0].Portal
	return portal, nil
}

func (p *Server) DescribeUsers(ctx context.Context, req *pb.DescribeUsersRequest) (*pb.DescribeUsersResponse, error) {
	s := ctxutil.GetSender(ctx)
	var (
		offset = pbutil.GetOffsetFromRequest(req)
		limit  = pbutil.GetLimitFromRequest(req)
	)

	senderPortal, err := getUserPortal(ctx, s.UserId)
	if err != nil {
		return nil, err
	}

	var rootGroupIds []string
	rootGroupId, err := getRootGroupId(ctx, s.UserId)
	if err != nil {
		return nil, err
	}

	if len(req.GetRootGroupId()) > 0 {
		err = CheckRootGroupIds(ctx, req.GetRootGroupId(), rootGroupId)
		if err != nil {
			return nil, err
		}
		rootGroupIds = append(rootGroupIds, req.GetRootGroupId()...)
	} else if senderPortal != constants.PortalGlobalAdmin {
		rootGroupIds = append(rootGroupIds, rootGroupId)
	}

	if len(req.GetRoleId()) > 0 {
		userIds, err := getRoleUserIds(ctx, req.GetRoleId(), req.GetUserId())
		if err != nil {
			return nil, err
		}
		req.UserId = userIds
		if len(req.UserId) == 0 {
			return &pb.DescribeUsersResponse{
				UserSet:    []*pb.User{},
				TotalCount: 0,
			}, nil
		}
	}

	res, err := imClient.ListUsers(ctx, &pbim.ListUsersRequest{
		SortKey:     req.GetSortKey().GetValue(),
		Reverse:     req.GetReverse().GetValue(),
		Offset:      uint32(offset),
		Limit:       uint32(limit),
		SearchWord:  []string{req.SearchWord.GetValue()},
		RootGroupId: rootGroupIds,
		GroupId:     req.GetGroupId(),
		UserId:      req.GetUserId(),
		Username:    req.GetUsername(),
		Email:       req.GetEmail(),
		PhoneNumber: req.GetPhoneNumber(),
		Status:      req.GetStatus(),
	})
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}

	reply := &pb.DescribeUsersResponse{
		UserSet:    models.ToPbUsers(res.UserSet),
		TotalCount: res.GetTotal(),
	}
	return reply, nil
}

func (p *Server) DescribeUsersDetail(ctx context.Context, req *pb.DescribeUsersRequest) (*pb.DescribeUsersDetailResponse, error) {
	s := ctxutil.GetSender(ctx)
	var (
		offset = pbutil.GetOffsetFromRequest(req)
		limit  = pbutil.GetLimitFromRequest(req)
	)

	senderPortal, err := getUserPortal(ctx, s.UserId)
	if err != nil {
		return nil, err
	}

	var rootGroupIds []string
	rootGroupId, err := getRootGroupId(ctx, s.UserId)
	if err != nil {
		return nil, err
	}

	if len(req.GetRootGroupId()) > 0 {
		err = CheckRootGroupIds(ctx, req.GetRootGroupId(), rootGroupId)
		if err != nil {
			return nil, err
		}
		rootGroupIds = append(rootGroupIds, req.GetRootGroupId()...)
	} else if senderPortal != constants.PortalGlobalAdmin {
		rootGroupIds = append(rootGroupIds, rootGroupId)
	}

	if len(req.GetRoleId()) > 0 {
		userIds, err := getRoleUserIds(ctx, req.GetRoleId(), req.GetUserId())
		if err != nil {
			return nil, err
		}
		req.UserId = userIds
		if len(req.UserId) == 0 {
			return &pb.DescribeUsersDetailResponse{
				UserDetailSet: []*pb.UserDetail{},
				TotalCount:    0,
			}, nil
		}
	}

	res, err := imClient.ListUsersWithGroup(ctx, &pbim.ListUsersRequest{
		SortKey:     req.GetSortKey().GetValue(),
		Reverse:     req.GetReverse().GetValue(),
		Offset:      uint32(offset),
		Limit:       uint32(limit),
		SearchWord:  []string{req.SearchWord.GetValue()},
		RootGroupId: rootGroupIds,
		GroupId:     req.GetGroupId(),
		UserId:      req.GetUserId(),
		Status:      req.GetStatus(),
		Username:    req.GetUsername(),
		Email:       req.GetEmail(),
		PhoneNumber: req.GetPhoneNumber(),
	})
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}

	var userDetails []*pb.UserDetail
	for _, userWithGroup := range res.UserSet {
		res, err := amClient.DescribeRoles(ctx, &pbam.DescribeRolesRequest{
			UserId: []string{userWithGroup.User.UserId},
		})
		if err != nil {
			return nil, err
		}

		userDetails = append(userDetails, &pb.UserDetail{
			User:     models.ToPbUser(userWithGroup.User),
			GroupSet: models.ToPbGroups(userWithGroup.GroupSet),
			RoleSet:  models.ToPbRoles(res.RoleSet),
		})
	}

	reply := &pb.DescribeUsersDetailResponse{
		UserDetailSet: userDetails,
		TotalCount:    res.GetTotal(),
	}
	return reply, nil
}

func (p *Server) DescribeGroups(ctx context.Context, req *pb.DescribeGroupsRequest) (*pb.DescribeGroupsResponse, error) {
	s := ctxutil.GetSender(ctx)
	var (
		offset = pbutil.GetOffsetFromRequest(req)
		limit  = pbutil.GetLimitFromRequest(req)
	)

	var rootGroupIds []string
	rootGroupId, err := getRootGroupId(ctx, s.UserId)
	if err != nil {
		return nil, err
	}

	if len(req.GetRootGroupId()) > 0 {
		err = CheckRootGroupIds(ctx, req.GetRootGroupId(), rootGroupId)
		if err != nil {
			return nil, err
		}
		rootGroupIds = append(rootGroupIds, req.GetRootGroupId()...)
	} else {
		rootGroupIds = append(rootGroupIds, rootGroupId)
	}

	res, err := imClient.ListGroups(ctx, &pbim.ListGroupsRequest{
		Limit:         uint32(limit),
		Offset:        uint32(offset),
		SortKey:       req.GetSortKey().GetValue(),
		Reverse:       req.GetReverse().GetValue(),
		SearchWord:    []string{req.SearchWord.GetValue()},
		RootGroupId:   rootGroupIds,
		ParentGroupId: req.GetParentGroupId(),
		GroupId:       req.GetGroupId(),
		GroupPath:     req.GetGroupPath(),
		Status:        req.GetStatus(),
		GroupName:     req.GetGroupName(),
	})
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}

	return &pb.DescribeGroupsResponse{
		GroupSet:   models.ToPbGroups(res.GroupSet),
		TotalCount: res.GetTotal(),
	}, nil
}

func (p *Server) DescribeGroupsDetail(ctx context.Context, req *pb.DescribeGroupsRequest) (*pb.DescribeGroupsDetailResponse, error) {
	s := ctxutil.GetSender(ctx)
	var (
		offset = pbutil.GetOffsetFromRequest(req)
		limit  = pbutil.GetLimitFromRequest(req)
	)

	var rootGroupIds []string
	rootGroupId, err := getRootGroupId(ctx, s.UserId)
	if err != nil {
		return nil, err
	}

	if len(req.GetRootGroupId()) > 0 {
		err = CheckRootGroupIds(ctx, req.GetRootGroupId(), rootGroupId)
		if err != nil {
			return nil, err
		}
		rootGroupIds = append(rootGroupIds, req.GetRootGroupId()...)
	} else {
		rootGroupIds = append(rootGroupIds, rootGroupId)
	}

	err = CheckRootGroupIds(ctx, req.GetRootGroupId(), rootGroupId)
	if err != nil {
		return nil, err
	}
	rootGroupIds = append(rootGroupIds, req.GetRootGroupId()...)

	res, err := imClient.ListGroupsWithUser(ctx, &pbim.ListGroupsRequest{
		Limit:         uint32(limit),
		Offset:        uint32(offset),
		SortKey:       req.GetSortKey().GetValue(),
		Reverse:       req.GetReverse().GetValue(),
		SearchWord:    []string{req.SearchWord.GetValue()},
		RootGroupId:   rootGroupIds,
		ParentGroupId: req.GetParentGroupId(),
		GroupId:       req.GetGroupId(),
		GroupPath:     req.GetGroupPath(),
		Status:        req.GetStatus(),
		GroupName:     req.GetGroupName(),
	})
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}
	var groupDetails []*pb.GroupDetail
	for _, groupWithUser := range res.GroupSet {
		groupDetails = append(groupDetails, &pb.GroupDetail{
			Group:   models.ToPbGroup(groupWithUser.Group),
			UserSet: models.ToPbUsers(groupWithUser.UserSet),
		})
	}

	return &pb.DescribeGroupsDetailResponse{
		GroupDetailSet: groupDetails,
		TotalCount:     res.GetTotal(),
	}, nil
}

func (p *Server) ModifyUser(ctx context.Context, req *pb.ModifyUserRequest) (*pb.ModifyUserResponse, error) {
	userId := req.GetUserId().GetValue()
	_, err := CheckUsersPermission(ctx, []string{userId})
	if err != nil {
		return nil, err
	}

	password := req.GetPassword().GetValue()
	if password != "" {
		_, err := imClient.ModifyPassword(ctx, &pbim.ModifyPasswordRequest{
			UserId:   userId,
			Password: password,
		})
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
		}
	}

	_, err = imClient.ModifyUser(ctx, &pbim.ModifyUserRequest{
		UserId:      userId,
		Username:    req.GetUsername().GetValue(),
		Email:       req.GetEmail().GetValue(),
		PhoneNumber: req.GetPhoneNumber().GetValue(),
		Description: req.GetDescription().GetValue(),
	})
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}

	return &pb.ModifyUserResponse{
		UserId: req.UserId,
	}, nil
}

func (p *Server) DeleteUsers(ctx context.Context, req *pb.DeleteUsersRequest) (*pb.DeleteUsersResponse, error) {
	return nil, gerr.New(ctx, gerr.PermissionDenied, gerr.ErrorPermissionDenied)
	//userIds := req.GetUserId()
	//_, err := CheckUsersPermission(ctx, userIds)
	//if err != nil {
	//	return nil, err
	//}
	//_, err = amClient.UnbindUserRole(ctx, &pbam.UnbindUserRoleRequest{
	//	UserId: userIds,
	//})
	//if err != nil {
	//	return nil, gerr.NewWithDetail(ctx, gerr.PermissionDenied, err, gerr.ErrorCannotDeleteUsers)
	//}
	//_, err = imClient.DeleteUsers(ctx, &pbim.DeleteUsersRequest{
	//	UserId: userIds,
	//})
	//if err != nil {
	//	return nil, gerr.NewWithDetail(ctx, gerr.PermissionDenied, err, gerr.ErrorCannotDeleteUsers)
	//}
	//
	//return &pb.DeleteUsersResponse{
	//	UserId: userIds,
	//}, nil
}

func (p *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	s := ctxutil.GetSender(ctx)
	email := req.GetEmail().GetValue()
	roleId := req.GetRoleId().GetValue()
	phoneNumber := req.GetPhoneNumber().GetValue()
	username := getUsernameFromEmail(email)
	password := req.GetPassword().GetValue()
	description := req.GetDescription().GetValue()

	res, err := imClient.ListUsers(ctx, &pbim.ListUsersRequest{
		Limit:  1,
		Status: []string{constants.StatusActive},
		Email:  []string{email},
	})
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}
	if res.Total > 0 {
		return nil, gerr.New(ctx, gerr.PermissionDenied, gerr.ErrorEmailExists, email)
	}

	createUserResponse, err := imClient.CreateUser(ctx, &pbim.CreateUserRequest{
		Email:       email,
		PhoneNumber: phoneNumber,
		Username:    username,
		Password:    password,
		Description: description,
	})
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}

	userId := createUserResponse.UserId

	// get sender portal

	senderPortal, err := getUserPortal(ctx, s.UserId)
	if err != nil {
		return nil, err
	}

	// get user portal
	getRoleResponse, err := amClient.GetRole(ctx, &pbam.GetRoleRequest{
		RoleId: roleId,
	})
	if err != nil {
		return nil, err
	}
	userPortal := getRoleResponse.Role.Portal

	// create group and join group
	if senderPortal != userPortal {
		err := createAndJoinRootGroup(ctx, userId)
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
		}
	} else {
		rootGroupId, err := getRootGroupId(ctx, s.UserId)
		if err != nil {
			return nil, err
		}
		_, err = imClient.JoinGroup(ctx, &pbim.JoinGroupRequest{
			GroupId: []string{rootGroupId},
			UserId:  []string{userId},
		})
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
		}
	}

	_, err = amClient.BindUserRole(ctx, &pbam.BindUserRoleRequest{
		UserId: []string{userId},
		RoleId: []string{roleId},
	})
	if err != nil {
		logger.Error(ctx, "Failed to bind user [%s] with role [%s]: %+v", userId, roleId, err)
		_, deleteErr := imClient.DeleteUsers(ctx, &pbim.DeleteUsersRequest{
			UserId: []string{createUserResponse.UserId},
		})
		if deleteErr != nil {
			logger.Error(ctx, "Failed to delete user [%s]: %+v", userId, err)
		}
		return nil, err
	}

	//get the username of current login user.
	resp, err := imClient.GetUser(ctx, &pbim.GetUserRequest{UserId: s.UserId})
	if err != nil {
		logger.Error(ctx, "Failed to get user [%s], %+v", s.UserId, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}
	senderUserName := resp.GetUser().GetUsername()
	platformName := pi.Global().GlobalConfig().BasicCfg.PlatformName
	platformUrl := pi.Global().GlobalConfig().BasicCfg.PlatformUrl

	if !stringutil.StringIn(s.UserId, constants.InternalUsers) {
		var emailNotifications []*models.EmailNotification
		if roleId == constants.RoleIsv {
			emailNotifications = append(emailNotifications, &models.EmailNotification{
				Title:       constants.AdminInviteIsvNotifyTitle.GetDefaultMessage(senderUserName),
				Content:     constants.AdminInviteIsvNotifyContent.GetDefaultMessage(platformName, username, senderUserName, platformName, platformUrl, platformUrl, platformUrl, email, password),
				Owner:       s.UserId,
				ContentType: constants.NfContentTypeInvite,
				Addresses:   []string{email},
			})
		} else {
			emailNotifications = append(emailNotifications, &models.EmailNotification{
				Title:       constants.AdminInviteUserNotifyTitle.GetDefaultMessage(senderUserName),
				Content:     constants.AdminInviteUserNotifyContent.GetDefaultMessage(platformName, username, senderUserName, platformName, platformUrl, platformUrl, platformUrl, email, password),
				Owner:       s.UserId,
				ContentType: constants.NfContentTypeInvite,
				Addresses:   []string{email},
			})
		}

		nfclient.SendEmailNotification(ctx, emailNotifications)
	}

	reply := &pb.CreateUserResponse{
		UserId: pbutil.ToProtoString(createUserResponse.UserId),
	}

	return reply, nil
}

func (p *Server) IsvCreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	s := ctxutil.GetSender(ctx)
	email := req.GetEmail().GetValue()
	roleId := req.GetRoleId().GetValue()
	phoneNumber := req.GetPhoneNumber().GetValue()
	username := getUsernameFromEmail(email)
	password := req.GetPassword().GetValue()
	description := req.GetDescription().GetValue()

	res, err := imClient.ListUsers(ctx, &pbim.ListUsersRequest{
		Limit:  1,
		Status: []string{constants.StatusActive},
		Email:  []string{email},
	})
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}
	if res.Total > 0 {
		return nil, gerr.New(ctx, gerr.PermissionDenied, gerr.ErrorEmailExists, email)
	}

	// isv can not create user with role isv
	if roleId == constants.RoleIsv {
		return nil, gerr.New(ctx, gerr.PermissionDenied, gerr.ErrorCannotCreateUserWithRole, roleId)
	}

	createUserResponse, err := imClient.CreateUser(ctx, &pbim.CreateUserRequest{
		Email:       email,
		PhoneNumber: phoneNumber,
		Username:    username,
		Password:    password,
		Description: description,
	})
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}

	userId := createUserResponse.UserId

	rootGroupId, err := getRootGroupId(ctx, s.UserId)
	if err != nil {
		return nil, err
	}
	_, err = imClient.JoinGroup(ctx, &pbim.JoinGroupRequest{
		GroupId: []string{rootGroupId},
		UserId:  []string{userId},
	})
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}

	_, err = amClient.BindUserRole(ctx, &pbam.BindUserRoleRequest{
		UserId: []string{createUserResponse.UserId},
		RoleId: []string{roleId},
	})
	if err != nil {
		logger.Error(ctx, "Failed to bind user [%s] with role [%s]: %+v", createUserResponse.UserId, roleId, err)
		_, deleteErr := imClient.DeleteUsers(ctx, &pbim.DeleteUsersRequest{
			UserId: []string{createUserResponse.UserId},
		})
		if deleteErr != nil {
			logger.Error(ctx, "Failed to delete user [%s]: %+v", createUserResponse.UserId, err)
		}
		return nil, err
	}

	if !stringutil.StringIn(s.UserId, constants.InternalUsers) {
		var emailNotifications []*models.EmailNotification
		getUserResponse, err := imClient.GetUser(ctx, &pbim.GetUserRequest{
			UserId: s.UserId,
		})
		if err != nil {
			logger.Error(ctx, "Failed to get user [%s]: %+v", s.UserId, err)
		} else {
			senderUserName := getUserResponse.User.Username
			platformName := pi.Global().GlobalConfig().BasicCfg.PlatformName
			platformUrl := pi.Global().GlobalConfig().BasicCfg.PlatformUrl

			emailNotifications = append(emailNotifications, &models.EmailNotification{
				Title:       constants.IsvInviteMemberNotifyTitle.GetDefaultMessage(senderUserName, platformName),
				Content:     constants.IsvInviteMemberNotifyContent.GetDefaultMessage(platformName, username, senderUserName, platformName, platformUrl, platformUrl, platformUrl, email, password),
				Owner:       s.UserId,
				ContentType: constants.NfContentTypeInvite,
				Addresses:   []string{email},
			})
		}
		nfclient.SendEmailNotification(ctx, emailNotifications)
	}

	reply := &pb.CreateUserResponse{
		UserId: pbutil.ToProtoString(createUserResponse.UserId),
	}

	return reply, nil
}

func (p *Server) CreatePasswordReset(ctx context.Context, req *pb.CreatePasswordResetRequest) (*pb.CreatePasswordResetResponse, error) {
	userId := req.GetUserId().GetValue()

	_, err := CheckUsersPermission(ctx, []string{userId})
	if err != nil {
		return nil, err
	}

	password := req.GetPassword().GetValue()
	_, err = imClient.GetUser(ctx, &pbim.GetUserRequest{
		UserId: userId,
	})
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}

	b, err := imClient.ComparePassword(ctx, &pbim.ComparePasswordRequest{
		UserId:   userId,
		Password: password,
	})
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}
	if !b.Ok {
		return nil, gerr.New(ctx, gerr.PermissionDenied, gerr.ErrorPasswordIncorrect)
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

	_, err = imClient.ModifyPassword(ctx, &pbim.ModifyPasswordRequest{
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

func validateUserAndGroupExist(ctx context.Context, email string) (user *pbim.User, isUserExist bool, isGroupExist bool) {
	isUserExist = false
	isGroupExist = false
	res, err := imClient.ListUsersWithGroup(ctx, &pbim.ListUsersRequest{
		Limit:  1,
		Status: []string{constants.StatusActive},
		Email:  []string{email},
	})
	if err != nil {
		logger.Error(ctx, "List users with group failed: %+v", err)
		return
	} else if res.Total == 0 {
		return
	}

	isUserExist = true
	user = res.UserSet[0].User

	if len(res.UserSet[0].GroupSet) > 0 {
		isGroupExist = true
	}
	return
}

func validateUserPassword(ctx context.Context, userId, password string) bool {
	b, err := imClient.ComparePassword(ctx, &pbim.ComparePasswordRequest{
		UserId:   userId,
		Password: password,
	})
	if err != nil {
		logger.Error(ctx, "Compare password failed: %+v", err)
		return false
	}
	return b.Ok
}

func (p *Server) ValidateUserPassword(ctx context.Context, req *pb.ValidateUserPasswordRequest) (*pb.ValidateUserPasswordResponse, error) {
	user, isUserExist, _ := validateUserAndGroupExist(ctx, req.GetEmail())
	if !isUserExist {
		return nil, gerr.New(ctx, gerr.NotFound, gerr.ErrorEmailNotExists, req.GetEmail())
	}

	ok := validateUserPassword(ctx, user.UserId, req.GetPassword())
	if !ok {
		return nil, gerr.New(ctx, gerr.PermissionDenied, gerr.ErrorEmailPasswordNotMatched)
	}
	return &pb.ValidateUserPasswordResponse{
		Validated: true,
	}, nil
}

func (p *Server) CreateGroup(ctx context.Context, req *pb.CreateGroupRequest) (*pb.CreateGroupResponse, error) {
	s := ctxutil.GetSender(ctx)
	createGroupResponse, err := imClient.CreateGroup(ctx, &pbim.CreateGroupRequest{
		ParentGroupId: req.GetParentGroupId().GetValue(),
		GroupName:     req.GetName().GetValue(),
		Description:   req.GetDescription().GetValue(),
		Extra: map[string]string{
			OwnerKey:     s.UserId,
			OwnerPathKey: string(s.OwnerPath),
		},
	})
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}

	reply := &pb.CreateGroupResponse{
		GroupId: pbutil.ToProtoString(createGroupResponse.GroupId),
	}

	return reply, nil
}

func (p *Server) ModifyGroup(ctx context.Context, req *pb.ModifyGroupRequest) (*pb.ModifyGroupResponse, error) {
	groupId := req.GetGroupId().GetValue()
	_, err := CheckGroupsPermission(ctx, []string{groupId})
	if err != nil {
		return nil, err
	}

	_, err = imClient.ModifyGroup(ctx, &pbim.ModifyGroupRequest{
		GroupId:       groupId,
		ParentGroupId: req.GetParentGroupId().GetValue(),
		GroupName:     req.GetName().GetValue(),
		Description:   req.GetDescription().GetValue(),
	})
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}

	reply := &pb.ModifyGroupResponse{
		GroupId: req.GetGroupId(),
	}

	return reply, nil
}

func (p *Server) DeleteGroups(ctx context.Context, req *pb.DeleteGroupsRequest) (*pb.DeleteGroupsResponse, error) {
	groupIds := req.GetGroupId()

	_, err := CheckGroupsPermission(ctx, groupIds)
	if err != nil {
		return nil, err
	}

	_, err = imClient.DeleteGroups(ctx, &pbim.DeleteGroupsRequest{
		GroupId: groupIds,
	})
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.PermissionDenied, err, gerr.ErrorCannotDeleteGroups)
	}

	return &pb.DeleteGroupsResponse{
		GroupId: groupIds,
	}, nil
}

func (p *Server) JoinGroup(ctx context.Context, req *pb.JoinGroupRequest) (*pb.JoinGroupResponse, error) {
	groupIds := req.GetGroupId()
	userIds := req.GetUserId()
	_, err := CheckGroupsPermission(ctx, groupIds)
	if err != nil {
		return nil, err
	}
	userWithGroups, err := CheckUsersPermission(ctx, userIds)
	if err != nil {
		return nil, err
	}

	var oldGroupIds []string
	for _, userWithGroup := range userWithGroups {
		for _, group := range userWithGroup.GroupSet {
			if !stringutil.StringIn(group.GroupId, oldGroupIds) {
				oldGroupIds = append(oldGroupIds, group.GroupId)
			}
		}
	}

	_, err = imClient.LeaveGroup(ctx, &pbim.LeaveGroupRequest{
		GroupId: oldGroupIds,
		UserId:  userIds,
	})
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.PermissionDenied, err, gerr.ErrorCannotJoinGroup)
	}

	_, err = imClient.JoinGroup(ctx, &pbim.JoinGroupRequest{
		GroupId: groupIds,
		UserId:  userIds,
	})
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.PermissionDenied, err, gerr.ErrorCannotJoinGroup)
	}

	return &pb.JoinGroupResponse{
		GroupId: groupIds,
		UserId:  userIds,
	}, nil
}

func (p *Server) LeaveGroup(ctx context.Context, req *pb.LeaveGroupRequest) (*pb.LeaveGroupResponse, error) {
	return nil, gerr.New(ctx, gerr.PermissionDenied, gerr.ErrorPermissionDenied)
}
