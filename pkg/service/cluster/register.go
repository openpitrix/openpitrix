// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package cluster

import (
	"context"
	"time"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pi"
)

func RegisterClusterNode(ctx context.Context, clusterNode *models.ClusterNode) error {
	clusterNode.NodeId = models.NewClusterNodeId()
	clusterNode.CreateTime = time.Now()
	clusterNode.StatusTime = time.Now()
	_, err := pi.Global().DB(ctx).
		InsertInto(constants.TableClusterNode).
		Record(clusterNode).
		Exec()
	if err != nil {
		logger.Error(ctx, "Failed to insert table [%s] with cluster id [%s]: %+v",
			constants.TableClusterNode, clusterNode.ClusterId, err)
		return err
	}
	return nil
}

func RegisterClusterRole(ctx context.Context, clusterRole *models.ClusterRole) error {
	_, err := pi.Global().DB(ctx).
		InsertInto(constants.TableClusterRole).
		Record(clusterRole).
		Exec()
	if err != nil {
		logger.Error(ctx, "Failed to insert table [%s] with cluster id [%s]: %+v",
			constants.TableClusterRole, clusterRole.ClusterId, err)
		return err
	}
	return nil
}

func RegisterClusterWrapper(ctx context.Context, clusterWrapper *models.ClusterWrapper) error {
	clusterId := clusterWrapper.Cluster.ClusterId
	owner := clusterWrapper.Cluster.Owner
	ownerPath := clusterWrapper.Cluster.OwnerPath
	// register cluster
	if clusterWrapper.Cluster != nil {
		now := time.Now()
		clusterWrapper.Cluster.CreateTime = now
		clusterWrapper.Cluster.StatusTime = now
		if clusterWrapper.Cluster.UpgradeTime == nil {
			clusterWrapper.Cluster.UpgradeTime = &now
		}
		_, err := pi.Global().DB(ctx).
			InsertInto(constants.TableCluster).
			Record(clusterWrapper.Cluster).
			Exec()
		if err != nil {
			logger.Error(ctx, "Failed to insert table [%s] with cluster id [%s]: %+v",
				constants.TableCluster, clusterWrapper.Cluster.ClusterId, err)
			return err
		}
	}

	// register cluster node
	newClusterNodes := make(map[string]*models.ClusterNodeWithKeyPairs)
	for _, clusterNodeWithKeyPairs := range clusterWrapper.ClusterNodesWithKeyPairs {
		clusterNodeWithKeyPairs.ClusterId = clusterId
		clusterNodeWithKeyPairs.Owner = owner
		clusterNodeWithKeyPairs.OwnerPath = ownerPath
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
			InsertInto(constants.TableClusterCommon).
			Record(clusterCommon).
			Exec()
		if err != nil {
			logger.Error(ctx, "Failed to insert table [%s] with cluster id [%s]: %+v",
				constants.TableClusterCommon, clusterWrapper.Cluster.ClusterId, err)
			return err
		}
	}

	// register cluster link
	for _, clusterLink := range clusterWrapper.ClusterLinks {
		clusterLink.ClusterId = clusterId
		clusterLink.Owner = owner
		clusterLink.OwnerPath = ownerPath
		_, err := pi.Global().DB(ctx).
			InsertInto(constants.TableClusterLink).
			Record(clusterLink).
			Exec()
		if err != nil {
			logger.Error(ctx, "Failed to insert table [%s] with cluster id [%s]: %+v",
				constants.TableClusterLink, clusterWrapper.Cluster.ClusterId, err)
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
				InsertInto(constants.TableClusterLoadbalancer).
				Record(clusterLoadbalancer).
				Exec()
			if err != nil {
				logger.Error(ctx, "Failed to insert table [%s] with cluster id [%s]: %+v",
					constants.TableClusterLoadbalancer, clusterWrapper.Cluster.ClusterId, err)
				return err
			}
		}
	}

	return nil
}
