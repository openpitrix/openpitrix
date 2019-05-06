// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package account

import (
	"context"
	"fmt"
	"strings"

	pbim "kubesphere.io/im/pkg/pb"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/gerr"
	"openpitrix.io/openpitrix/pkg/sender"
	"openpitrix.io/openpitrix/pkg/util/ctxutil"
)

func CheckGroupsPermission(ctx context.Context, groupIds []string) ([]*pbim.Group, error) {
	s := ctxutil.GetSender(ctx)

	listGroupResponse, err := imClient.ListGroups(ctx, &pbim.ListGroupsRequest{
		GroupId: groupIds,
		Status:  []string{constants.StatusActive},
	})
	if err != nil {
		return nil, err
	}

	if len(groupIds) != len(listGroupResponse.GroupSet) {
		return nil, gerr.New(ctx, gerr.PermissionDenied, gerr.ErrorGroupNotFound, strings.Join(groupIds, ","))
	}

	for _, group := range listGroupResponse.GroupSet {
		ownerPath := sender.OwnerPath(group.Extra[OwnerPathKey])
		if !ownerPath.CheckPermission(s) {
			return nil, gerr.New(ctx, gerr.PermissionDenied, gerr.ErrorGroupAccessDenied, group.GroupId)
		}
	}
	return listGroupResponse.GroupSet, nil
}

func CheckUsersPermission(ctx context.Context, userIds []string) ([]*pbim.UserWithGroup, error) {
	s := ctxutil.GetSender(ctx)

	listUsersWithGroupResponse, err := imClient.ListUsersWithGroup(ctx, &pbim.ListUsersRequest{
		UserId: userIds,
		Status: []string{constants.StatusActive},
	})
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}

	if len(userIds) != len(listUsersWithGroupResponse.UserSet) {
		return nil, gerr.New(ctx, gerr.PermissionDenied, gerr.ErrorUserNotFound, strings.Join(userIds, ","))
	}

	for _, userWithGroup := range listUsersWithGroupResponse.UserSet {
		for _, group := range userWithGroup.GroupSet {
			ownerPath := sender.OwnerPath(group.Extra[OwnerPathKey])
			if !ownerPath.CheckPermission(s) {
				return nil, gerr.New(ctx, gerr.PermissionDenied, gerr.ErrorUserAccessDenied, userWithGroup.User.UserId)
			}
		}
	}
	return listUsersWithGroupResponse.UserSet, nil
}

func CheckRootGroupIds(ctx context.Context, checkGroupIds []string, rootGroupId string) error {
	s := ctxutil.GetSender(ctx)
	if s.UserId == constants.UserSystem {
		return nil
	}
	listGroupsResponse, err := imClient.ListGroups(ctx, &pbim.ListGroupsRequest{
		GroupId: checkGroupIds,
	})
	if err != nil {
		return err
	}
	for _, group := range listGroupsResponse.GroupSet {
		if !strings.HasPrefix(group.GroupPath, rootGroupId) {
			err = fmt.Errorf("invalid root group id [%s]", group.GroupId)
			return gerr.NewWithDetail(ctx, gerr.InvalidArgument, err, gerr.ErrorValidateFailed)
		}
	}
	return nil
}
