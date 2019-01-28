// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package account

import (
	"context"
	"fmt"
	"strings"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
	"openpitrix.io/openpitrix/pkg/util/stringutil"
)

type Client struct {
	pb.AccountManagerClient
}

func NewClient() (*Client, error) {
	conn, err := manager.NewClient(constants.AccountServiceHost, constants.AccountServicePort)
	if err != nil {
		return nil, err
	}
	return &Client{
		AccountManagerClient: pb.NewAccountManagerClient(conn),
	}, nil
}

func (c *Client) GetUsers(ctx context.Context, userIds []string) ([]*pb.User, error) {
	var internalUsers []*pb.User
	var noInternalUserIds []string
	for _, userId := range userIds {
		if stringutil.StringIn(userId, constants.InternalUsers) {
			internalUsers = append(internalUsers, &pb.User{
				UserId: pbutil.ToProtoString(userId),
				Role:   pbutil.ToProtoString(constants.RoleGlobalAdmin),
			})
		} else {
			noInternalUserIds = append(noInternalUserIds, userId)
		}
	}

	if len(noInternalUserIds) == 0 {
		return internalUsers, nil
	}

	response, err := c.DescribeUsers(ctx, &pb.DescribeUsersRequest{
		UserId: noInternalUserIds,
	})
	if err != nil {
		logger.Error(ctx, "Describe users %s failed: %+v", noInternalUserIds, err)
		return nil, err
	}
	if len(response.UserSet) != len(noInternalUserIds) {
		logger.Error(ctx, "Describe users %s with return count [%d]", userIds, len(response.UserSet)+len(internalUsers))
		return nil, fmt.Errorf("describe users %s with return count [%d]", userIds, len(response.UserSet)+len(internalUsers))
	}
	response.UserSet = append(response.UserSet, internalUsers...)
	return response.UserSet, nil
}

func (c *Client) GetUserGroupPath(ctx context.Context, userId string) (string, error) {
	var userGroupPath string

	var userIds []string
	userIds = append(userIds, userId)

	response, err := c.DescribeGroups(ctx, &pb.DescribeGroupsRequest{
		UserId: userIds,
	})
	if err != nil {
		logger.Error(ctx, "Describe groups %s failed: %+v", userIds, err)
		return "", err
	}

	respGroupSet := response.GroupSet

	//If one uer under different Group, get the highest Group Path.
	if len(respGroupSet) > 1 {
		minLevel := len(strings.Split(respGroupSet[0].GroupPath.GetValue(), "."))
		for _, group := range response.GroupSet {
			if len(strings.Split(group.GroupPath.GetValue(), ".")) < minLevel {
				minLevel = len(strings.Split(group.GroupPath.GetValue(), "."))
				userGroupPath = group.GroupPath.GetValue()
			}
		}

	} else if len(respGroupSet) == 1 {
		userGroupPath = response.GroupSet[0].GetGroupPath().GetValue()
	} else {
		return "", nil
	}

	return userGroupPath, nil

}

func GetUsers(ctx context.Context, userIds []string) ([]*pb.User, error) {
	client, err := NewClient()
	if err != nil {
		logger.Error(ctx, "Failed to create im client: %+v", err)
		return nil, err
	}
	response, err := client.GetUsers(ctx, userIds)
	if err != nil {
		return nil, err
	}
	return response, err
}

func GetIsvFromUsers(ctx context.Context, userIds []string) ([]*pb.User, error) {
	client, err := NewClient()
	if err != nil {
		logger.Error(ctx, "Failed to create im client: %+v", err)
		return nil, err
	}

	var owners []string
	for _, userId := range userIds {
		response, err := client.GetUserGroupOwner(ctx, &pb.GetUserGroupOwnerRequest{
			UserId: userId,
		})
		if err != nil {
			return nil, err
		}
		owners = append(owners, response.Owner)
	}

	response, err := client.GetUsers(ctx, owners)
	if err != nil {
		return nil, err
	}
	return response, err
}
