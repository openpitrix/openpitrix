// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package am

import (
	"context"

	"openpitrix.io/iam/pkg/pb/am"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/manager"
)

type Client struct {
	pbam.AccessManagerClient
}

func NewClient() (*Client, error) {
	conn, err := manager.NewClient(constants.AMServiceHost, constants.AMServicePort)
	if err != nil {
		return nil, err
	}

	return &Client{
		AccessManagerClient: pbam.NewAccessManagerClient(conn),
	}, nil
}

func GetRoleUsers(ctx context.Context, roles []string) ([]*pbam.UserWithRole, error) {
	client, err := NewClient()
	if err != nil {
		logger.Error(ctx, "Failed to create am client: %+v", err)
		return nil, err
	}
	response, err := client.DescribeUsersWithRole(ctx, &pbam.DescribeUsersWithRoleRequest{
		RoleId: roles,
	})
	if err != nil {
		return nil, err
	}
	return response.User, err
}
