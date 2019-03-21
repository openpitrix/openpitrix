// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file

package isv

import (
	"context"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/pb"
)

type Client struct {
	pb.IsvManagerClient
}

func NewClient() (*Client, error) {
	conn, err := manager.NewClient(constants.IsvManagerHost, constants.IsvManagerPort)
	if err != nil {
		return nil, err
	}
	return &Client{
		IsvManagerClient: pb.NewIsvManagerClient(conn),
	}, nil
}

func GetVendorInfos(ctx context.Context, userIds []string) ([]*pb.VendorVerifyInfo, error) {
	client, err := NewClient()
	if err != nil {
		logger.Error(ctx, "Failed to create app vendor client: %+v", err)
		return nil, err
	}
	response, err := client.DescribeVendorVerifyInfos(ctx, &pb.DescribeVendorVerifyInfosRequest{
		UserId: userIds,
	})
	if err != nil {
		return nil, err
	}
	return response.VendorVerifyInfoSet, err
}
