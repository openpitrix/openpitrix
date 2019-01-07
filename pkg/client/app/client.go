// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package app

import (
	"context"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/pb"
)

type Client struct {
	pb.AppManagerClient
}

func NewAppManagerClient() (*Client, error) {
	conn, err := manager.NewClient(constants.AppManagerHost, constants.AppManagerPort)
	if err != nil {
		return nil, err
	}
	return &Client{
		AppManagerClient: pb.NewAppManagerClient(conn),
	}, nil
}

func (c *Client) DescribeActiveAppsWithOwner(ctx context.Context, userId string, limit uint32, offset uint32) ([]*pb.App, int32, error) {
	var appUserIds []string
	appUserIds = append(appUserIds, userId)
	request := &pb.DescribeAppsRequest{
		Owner:  appUserIds,
		Limit:  limit,
		Offset: offset,
	}

	response, err := c.DescribeActiveApps(ctx, request)

	if err != nil {
		logger.Error(ctx, "Describe ActiveApps with Owners [%s] failed: %+v", userId, err)
		return nil, 0, err
	}
	return response.AppSet, int32(response.TotalCount), nil
}
