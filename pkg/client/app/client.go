// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package app

import (
	"context"
	"fmt"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
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

type DescribeAppsApi struct{}

func (p *DescribeAppsApi) SetRequest(ctx context.Context, req interface{}, limit uint32, offset uint32) error {
	switch r := req.(type) {
	case *pb.DescribeAppsRequest:
		r.Limit = limit
		r.Offset = offset
	default:
		return fmt.Errorf("invalid req")
	}
	return nil
}

func (p *DescribeAppsApi) Describe(ctx context.Context, req interface{}, advancedParams ...string) (pbutil.DescribeResponse, error) {
	switch r := req.(type) {
	case *pb.DescribeAppsRequest:
		appClient, err := NewAppManagerClient()
		if err != nil {
			return nil, err
		}
		if len(advancedParams) > 0 && advancedParams[0] == "active" {
			return appClient.DescribeActiveApps(ctx, r)
		} else {
			return appClient.DescribeApps(ctx, r)
		}
	default:
		return nil, fmt.Errorf("invalid req")
	}
}
