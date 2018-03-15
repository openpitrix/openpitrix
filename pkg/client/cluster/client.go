// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package cluster

import (
	"context"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/pb"
)

func NewClusterManagerClient(ctx context.Context) (pb.ClusterManagerClient, error) {
	conn, err := manager.NewClient(ctx, constants.ClusterManagerHost, constants.ClusterManagerPort)
	if err != nil {
		return nil, err
	}
	return pb.NewClusterManagerClient(conn), err
}

func ModifyCluster(request *pb.ModifyClusterRequest) error {
	ctx := context.Background()
	client, err := NewClusterManagerClient(ctx)
	if err != nil {
		return err
	}
	_, err = client.ModifyCluster(ctx, request)
	if err != nil {
		return err
	}
	return nil
}

func ModifyClusterNode(request *pb.ModifyClusterNodeRequest) error {
	ctx := context.Background()
	client, err := NewClusterManagerClient(ctx)
	if err != nil {
		return err
	}
	_, err = client.ModifyClusterNode(ctx, request)
	if err != nil {
		return err
	}
	return nil
}
