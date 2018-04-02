// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package runtimeenv

import (
	"context"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/pb"
)

func NewRuntimeEnvManagerClient(ctx context.Context) (pb.RuntimeEnvManagerClient, error) {
	conn, err := manager.NewClient(ctx, constants.RuntimeManagerHost, constants.RuntimeEnvManagerPort)
	if err != nil {
		return nil, err
	}
	return pb.NewRuntimeEnvManagerClient(conn), err
}

func DescribeRuntimeEnvs(request *pb.DescribeRuntimeEnvsRequest) (*pb.DescribeRuntimeEnvsResponse, error) {
	ctx := context.Background()
	client, err := NewRuntimeEnvManagerClient(ctx)
	if err != nil {
		return nil, err
	}
	response, err := client.DescribeRuntimeEnvs(ctx, request)
	if err != nil {
		return nil, err
	}
	return response, err
}
