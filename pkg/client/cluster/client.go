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
	"openpitrix.io/openpitrix/pkg/util/pbutil"
)

type Client struct {
	pb.ClusterManagerClient
}

func NewClient(ctx context.Context) (*Client, error) {
	conn, err := manager.NewClient(ctx, constants.ClusterManagerHost, constants.ClusterManagerPort)
	if err != nil {
		return nil, err
	}
	return &Client{
		ClusterManagerClient: pb.NewClusterManagerClient(conn),
	}, nil
}

func (c *Client) GetClusterNodes(ctx context.Context, nodeIds []string) ([]*pb.ClusterNode, error) {
	response, err := c.DescribeClusterNodes(ctx, &pb.DescribeClusterNodesRequest{
		NodeId: nodeIds,
	})
	if err != nil {
		logger.Error("Describe cluster nodes %s failed: %+v", nodeIds, err)
		return nil, err
	}
	if len(response.ClusterNodeSet) != len(nodeIds) {
		logger.Error("Describe cluster nodes %s with return count [%d]", nodeIds, len(response.ClusterNodeSet))
		return nil, fmt.Errorf("describe cluster nodes %s with return count [%d]", nodeIds, len(response.ClusterNodeSet))
	}
	return response.ClusterNodeSet, nil
}

func (c *Client) GetClusters(ctx context.Context, clusterIds []string) ([]*pb.Cluster, error) {
	response, err := c.DescribeClusters(ctx, &pb.DescribeClustersRequest{
		ClusterId: clusterIds,
	})
	if err != nil {
		logger.Error("Describe clusters %s failed: %+v", clusterIds, err)
		return nil, err
	}
	if len(response.ClusterSet) != len(clusterIds) {
		logger.Error("Describe clusters %s with return count [%d]", clusterIds, len(response.ClusterSet))
		return nil, fmt.Errorf("describe clusters %s with return count [%d]", clusterIds, len(response.ClusterSet))
	}
	return response.ClusterSet, nil
}

func (c *Client) GetClusterWrappers(ctx context.Context, clusterIds []string) ([]*models.ClusterWrapper, error) {
	pbClusterSet, err := c.GetClusters(ctx, clusterIds)
	if err != nil {
		return nil, err
	}
	var clusterWrappers []*models.ClusterWrapper
	for _, pbCluster := range pbClusterSet {
		clusterWrappers = append(clusterWrappers, models.PbToClusterWrapper(pbCluster))
	}
	return clusterWrappers, nil
}

func (c *Client) ModifyClusterTransitionStatus(ctx context.Context, clusterId string, transitionStatus string) error {
	_, err := c.ModifyCluster(ctx, &pb.ModifyClusterRequest{
		Cluster: &pb.Cluster{
			ClusterId:        pbutil.ToProtoString(clusterId),
			TransitionStatus: pbutil.ToProtoString(transitionStatus),
		},
	})
	return err
}

func (c *Client) ModifyClusterStatus(ctx context.Context, clusterId string, status string) error {
	_, err := c.ModifyCluster(ctx, &pb.ModifyClusterRequest{
		Cluster: &pb.Cluster{
			ClusterId: pbutil.ToProtoString(clusterId),
			Status:    pbutil.ToProtoString(status),
		},
	})
	return err
}
