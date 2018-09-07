// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package cluster

import (
	"context"
	"time"

	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pi"
)

func RegisterClusterNode(ctx context.Context, clusterNode *models.ClusterNode) error {
	clusterNode.NodeId = models.NewClusterNodeId()
	clusterNode.CreateTime = time.Now()
	clusterNode.StatusTime = time.Now()
	_, err := pi.Global().DB(ctx).
		InsertInto(models.ClusterNodeTableName).
		Columns(models.ClusterNodeColumns...).
		Record(clusterNode).
		Exec()
	if err != nil {
		logger.Error(ctx, "Failed to insert table [%s] with cluster id [%s]: %+v",
			models.ClusterNodeTableName, clusterNode.ClusterId, err)
		return err
	}
	return nil
}

func RegisterClusterRole(ctx context.Context, clusterRole *models.ClusterRole) error {
	_, err := pi.Global().DB(ctx).
		InsertInto(models.ClusterRoleTableName).
		Columns(models.ClusterRoleColumns...).
		Record(clusterRole).
		Exec()
	if err != nil {
		logger.Error(ctx, "Failed to insert table [%s] with cluster id [%s]: %+v",
			models.ClusterRoleTableName, clusterRole.ClusterId, err)
		return err
	}
	return nil
}

func RegisterClusterWrapper(ctx context.Context, clusterWrapper *models.ClusterWrapper) error {
	clusterId := clusterWrapper.Cluster.ClusterId
	owner := clusterWrapper.Cluster.Owner
	// register cluster
	if clusterWrapper.Cluster != nil {
		now := time.Now()
		clusterWrapper.Cluster.CreateTime = now
		clusterWrapper.Cluster.StatusTime = now
		if clusterWrapper.Cluster.UpgradeTime == nil {
			clusterWrapper.Cluster.UpgradeTime = &now
		}
		_, err := pi.Global().DB(ctx).
			InsertInto(models.ClusterTableName).
			Columns(models.ClusterColumns...).
			Record(clusterWrapper.Cluster).
			Exec()
		if err != nil {
			logger.Error(ctx, "Failed to insert table [%s] with cluster id [%s]: %+v",
				models.ClusterTableName, clusterWrapper.Cluster.ClusterId, err)
			return err
		}
	}

	// register cluster node
	newClusterNodes := make(map[string]*models.ClusterNodeWithKeyPairs)
	for _, clusterNodeWithKeyPairs := range clusterWrapper.ClusterNodesWithKeyPairs {
		clusterNodeWithKeyPairs.ClusterId = clusterId
		clusterNodeWithKeyPairs.Owner = owner
		err := RegisterClusterNode(ctx, clusterNodeWithKeyPairs.ClusterNode)
		if err != nil {
			return err
		}
		newClusterNodes[clusterNodeWithKeyPairs.NodeId] = clusterNodeWithKeyPairs
	}

	clusterWrapper.ClusterNodesWithKeyPairs = newClusterNodes

	// register cluster common
	for _, clusterCommon := range clusterWrapper.ClusterCommons {
		clusterCommon.ClusterId = clusterId
		_, err := pi.Global().DB(ctx).
			InsertInto(models.ClusterCommonTableName).
			Columns(models.ClusterCommonColumns...).
			Record(clusterCommon).
			Exec()
		if err != nil {
			logger.Error(ctx, "Failed to insert table [%s] with cluster id [%s]: %+v",
				models.ClusterCommonTableName, clusterWrapper.Cluster.ClusterId, err)
			return err
		}
	}

	// register cluster link
	for _, clusterLink := range clusterWrapper.ClusterLinks {
		clusterLink.ClusterId = clusterId
		clusterLink.Owner = owner
		_, err := pi.Global().DB(ctx).
			InsertInto(models.ClusterLinkTableName).
			Columns(models.ClusterLinkColumns...).
			Record(clusterLink).
			Exec()
		if err != nil {
			logger.Error(ctx, "Failed to insert table [%s] with cluster id [%s]: %+v",
				models.ClusterLinkTableName, clusterWrapper.Cluster.ClusterId, err)
			return err
		}
	}

	// register cluster role
	for _, clusterRole := range clusterWrapper.ClusterRoles {
		clusterRole.ClusterId = clusterId
		err := RegisterClusterRole(ctx, clusterRole)
		if err != nil {
			return err
		}
	}

	// register cluster loadbalancer
	for _, clusterLoadbalancers := range clusterWrapper.ClusterLoadbalancers {
		for _, clusterLoadbalancer := range clusterLoadbalancers {
			clusterLoadbalancer.ClusterId = clusterId
			_, err := pi.Global().DB(ctx).
				InsertInto(models.ClusterLoadbalancerTableName).
				Columns(models.ClusterLoadbalancerColumns...).
				Record(clusterLoadbalancer).
				Exec()
			if err != nil {
				logger.Error(ctx, "Failed to insert table [%s] with cluster id [%s]: %+v",
					models.ClusterLoadbalancerTableName, clusterWrapper.Cluster.ClusterId, err)
				return err
			}
		}
	}

	return nil
}
