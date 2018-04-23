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
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/utils"
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

func GetClusters(ctx context.Context, client pb.ClusterManagerClient, clusterIds []string) ([]*pb.Cluster, error) {
	response, err := client.DescribeClusters(ctx, &pb.DescribeClustersRequest{
		ClusterId: clusterIds,
	})
	if err != nil {
		logger.Errorf("Describe clusters %s failed: %+v", clusterIds, err)
		return nil, err
	}
	if len(response.ClusterSet) != len(clusterIds) {
		logger.Errorf("Describe clusters %s with return count [%d]", clusterIds, len(response.ClusterSet))
		return nil, fmt.Errorf("describe clusters %s with return count [%d]", clusterIds, len(response.ClusterSet))
	}
	return response.ClusterSet, nil
}

func GetClusterWrappers(ctx context.Context, client pb.ClusterManagerClient, clusterIds []string) ([]*models.ClusterWrapper, error) {
	pbClusterSet, err := GetClusters(ctx, client, clusterIds)
	if err != nil {
		return nil, err
	}
	var clusterWrappers []*models.ClusterWrapper
	for _, pbCluster := range pbClusterSet {
		clusterWrappers = append(clusterWrappers, models.PbToClusterWrapper(pbCluster))
	}
	return clusterWrappers, nil
}

func ModifyClusterTransitionStatus(ctx context.Context, client pb.ClusterManagerClient, clusterId string, transitionStatus string) error {
	_, err := client.ModifyCluster(ctx, &pb.ModifyClusterRequest{
		Cluster: &pb.Cluster{
			ClusterId:        utils.ToProtoString(clusterId),
			TransitionStatus: utils.ToProtoString(transitionStatus),
		},
	})
	return err
}
