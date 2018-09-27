// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package iam

import (
	"context"
	"fmt"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/pb"
)

type Client struct {
	pb.AccountManagerClient
}

func NewClient() (*Client, error) {
	conn, err := manager.NewClient(constants.IAMServiceHost, constants.IAMServicePort)
	if err != nil {
		return nil, err
	}
	return &Client{
		AccountManagerClient: pb.NewAccountManagerClient(conn),
	}, nil
}

func (c *Client) GetUsers(ctx context.Context, userIds []string) ([]*pb.User, error) {
	response, err := c.DescribeUsers(ctx, &pb.DescribeUsersRequest{
		UserId: userIds,
	})
	if err != nil {
		logger.Error(ctx, "Describe users %s failed: %+v", userIds, err)
		return nil, err
	}
	if len(response.UserSet) != len(userIds) {
		logger.Error(ctx, "Describe users %s with return count [%d]", userIds, len(response.UserSet))
		return nil, fmt.Errorf("describe users %s with return count [%d]", userIds, len(response.UserSet))
	}
	return response.UserSet, nil
}
