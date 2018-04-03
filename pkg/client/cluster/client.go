// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package cluster

import (
	"context"
	"fmt"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/logger"
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

func GetClusterNodes(ctx context.Context, client pb.ClusterManagerClient, nodeIds []string) ([]*pb.ClusterNode, error) {
	response, err := client.DescribeClusterNodes(ctx, &pb.DescribeClusterNodesRequest{
		NodeId: nodeIds,
	})
	if err != nil {
		logger.Errorf("Describe cluster nodes %s failed: %+v", nodeIds, err)
		return nil, err
	}
	if len(response.ClusterNodeSet) != len(nodeIds) {
		logger.Errorf("Describe cluster nodes %s with return count [%d]", nodeIds, len(response.ClusterNodeSet))
		return nil, fmt.Errorf("describe cluster nodes %s with return count [%d]", nodeIds, len(response.ClusterNodeSet))
	}
	return response.ClusterNodeSet, nil
}
