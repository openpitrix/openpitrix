// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package app

import (
	"context"

	accountclient "openpitrix.io/openpitrix/pkg/client/account"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/gerr"
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

func (c *Client) DescribeAppsWithAppVendorUserId(ctx context.Context, appVendorUserId string, limit uint32, offset uint32) ([]*pb.App, int32, error) {
	account, err := accountclient.NewClient()
	if err != nil {
		return nil, 0, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}

	groupPath, _ := account.GetUserGroupPath(ctx, appVendorUserId)
	var groupPaths []string
	groupPaths = append(groupPaths, groupPath)

	req := pb.DescribeAppsRequest{
		OwnerPath: groupPaths,
		Limit:     limit,
		Offset:    offset,
	}

	response, err := c.DescribeApps(ctx, &req)
	if err != nil {
		logger.Error(ctx, "Describe apps failed: %+v", err)
		return nil, 0, err
	}
	return response.AppSet, int32(response.TotalCount), nil
}
