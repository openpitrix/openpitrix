// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package access

import (
	"context"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/pb"
)

type Client struct {
	pb.AccessManagerClient
}

func NewClient() (*Client, error) {
	conn, err := manager.NewClient(constants.AccountServiceHost, constants.AccountServicePort)
	if err != nil {
		return nil, err
	}
	return &Client{
		AccessManagerClient: pb.NewAccessManagerClient(conn),
	}, nil
}

func GetActionBundleRoles(ctx context.Context, actionBundleIds []string, statuses []string) ([]*pb.Role, error) {
	client, err := NewClient()
	if err != nil {
		logger.Error(ctx, "Failed to create account manager client: %+v", err)
		return nil, err
	}

	response, err := client.DescribeRoles(ctx, &pb.DescribeRolesRequest{
		ActionBundleId: actionBundleIds,
		Status:         statuses,
	})
	if err != nil {
		logger.Error(ctx, "Describe users failed: %+v", err)
		return nil, err
	}

	return response.RoleSet, nil
}
