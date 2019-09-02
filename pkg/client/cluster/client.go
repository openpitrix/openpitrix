// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package cluster

import (
	"context"
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes/wrappers"

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

func NewClient() (*Client, error) {
	conn, err := manager.NewClient(constants.ClusterManagerHost, constants.ClusterManagerPort)
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
		logger.Error(ctx, "Describe cluster nodes %s failed: %+v", nodeIds, err)
		return nil, err
	}
	if len(response.ClusterNodeSet) != len(nodeIds) {
		logger.Error(ctx, "Describe cluster nodes %s with return count [%d]", nodeIds, len(response.ClusterNodeSet))
		return nil, fmt.Errorf("describe cluster nodes %s with return count [%d]", nodeIds, len(response.ClusterNodeSet))
	}
	return response.ClusterNodeSet, nil
}

func (c *Client) GetClusters(ctx context.Context, clusterIds []string) ([]*pb.Cluster, error) {
	response, err := c.DescribeClusters(ctx, &pb.DescribeClustersRequest{
		ClusterId: clusterIds,
	})
	if err != nil {
		logger.Error(ctx, "Describe clusters %s failed: %+v", clusterIds, err)
		return nil, err
	}
	if int(response.TotalCount) != len(clusterIds) {
		debugResponse, err := c.DescribeDebugClusters(ctx, &pb.DescribeClustersRequest{
			ClusterId: clusterIds,
		})
		if err != nil {
			logger.Error(ctx, "Describe debug clusters %s failed: %+v", clusterIds, err)
			return nil, err
		}

		totalCount := response.TotalCount + debugResponse.TotalCount

		if int(totalCount) != len(clusterIds) {
			logger.Error(ctx, "Describe clusters %s with return count [%d]", clusterIds, totalCount)
			return nil, fmt.Errorf("describe clusters %s with return count [%d]", clusterIds, totalCount)
		}

		response.ClusterSet = append(response.ClusterSet, debugResponse.ClusterSet...)
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
			StatusTime:       pbutil.ToProtoTimestamp(time.Now()),
		},
	})
	return err
}

func (c *Client) ModifyClusterStatus(ctx context.Context, clusterId string, status string) error {
	_, err := c.ModifyCluster(ctx, &pb.ModifyClusterRequest{
		Cluster: &pb.Cluster{
			ClusterId:  pbutil.ToProtoString(clusterId),
			Status:     pbutil.ToProtoString(status),
			StatusTime: pbutil.ToProtoTimestamp(time.Now()),
		},
	})
	return err
}

func (c *Client) ModifyClusterNodeTransitionStatus(ctx context.Context, nodeId string, transitionStatus string) error {
	_, err := c.ModifyClusterNode(ctx, &pb.ModifyClusterNodeRequest{
		ClusterNode: &pb.ClusterNode{
			NodeId:           pbutil.ToProtoString(nodeId),
			TransitionStatus: pbutil.ToProtoString(transitionStatus),
		},
	})
	return err
}

func (c *Client) ModifyClusterNodeStatus(ctx context.Context, nodeId string, status string) error {
	_, err := c.ModifyClusterNode(ctx, &pb.ModifyClusterNodeRequest{
		ClusterNode: &pb.ClusterNode{
			NodeId: pbutil.ToProtoString(nodeId),
			Status: pbutil.ToProtoString(status),
		},
	})
	return err
}

func (c *Client) DescribeClustersWithFrontgateId(ctx context.Context, frontgateId string, status []string, debug bool) ([]*pb.Cluster, error) {
	var request *pb.DescribeClustersRequest
	var response *pb.DescribeClustersResponse
	var err error
	if status == nil {
		request = &pb.DescribeClustersRequest{
			FrontgateId: []string{frontgateId},
		}
	} else {
		request = &pb.DescribeClustersRequest{
			FrontgateId: []string{frontgateId},
			Status:      status,
		}
	}
	if debug {
		response, err = c.DescribeDebugClusters(ctx, request)
		if err != nil {
			logger.Error(ctx, "Describe debug clusters with frontgate [%s] failed: %+v", frontgateId, err)
			return nil, err
		}
	} else {
		response, err = c.DescribeClusters(ctx, request)
		if err != nil {
			logger.Error(ctx, "Describe clusters with frontgate [%s] failed: %+v", frontgateId, err)
			return nil, err
		}
	}

	return response.ClusterSet, nil
}

func (c *Client) DeleteAndCeaseClusters(ctx context.Context, clusterIds []string, force bool) error {
	var err error

	deleteReq := &pb.DeleteClustersRequest{
		ClusterId: clusterIds,
		Force: &wrappers.BoolValue{
			Value: force,
		},
	}
	_, err = c.DeleteClusters(ctx, deleteReq)
	if err != nil {
		return err
	}
	ceaseReq := &pb.CeaseClustersRequest{
		ClusterId: clusterIds,
		Force: &wrappers.BoolValue{
			Value: force,
		},
	}
	_, err = c.CeaseClusters(ctx, ceaseReq)
	if err != nil {
		return err
	}
	return nil
}
