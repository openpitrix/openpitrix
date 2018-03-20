package cluster

import (
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pi"
)

type Register struct {
	*pi.Pi
}

func (r *Register) RegisterClusterWrapper(clusterId, runtimeEnvId, frontgateId, owner string, clusterWrapper *models.ClusterWrapper) error {
	// register cluster
	if clusterWrapper.Cluster != nil {
		clusterWrapper.Cluster.ClusterId = clusterId
		clusterWrapper.Cluster.RuntimeEnvId = runtimeEnvId
		clusterWrapper.Cluster.FrontgateId = frontgateId
		clusterWrapper.Cluster.Owner = owner
		_, err := r.Db.
			InsertInto(models.ClusterTableName).
			Columns(models.ClusterColumns...).
			Record(clusterWrapper.Cluster).
			Exec()
		if err != nil {
			logger.Errorf("Failed to insert table [%s] with cluster id [%s]: %+v",
				models.ClusterTableName, clusterWrapper.Cluster.ClusterId, err)
			return err
		}
	}

	// register cluster node
	for _, clusterNode := range clusterWrapper.ClusterNodes {
		clusterNode.ClusterId = clusterId
		clusterNode.NodeId = models.NewClusterNodeId()
		clusterNode.Owner = owner
		_, err := r.Db.
			InsertInto(models.ClusterNodeTableName).
			Columns(models.ClusterNodeColumns...).
			Record(clusterNode).
			Exec()
		if err != nil {
			logger.Errorf("Failed to insert table [%s] with cluster id [%s]: %+v",
				models.ClusterNodeTableName, clusterWrapper.Cluster.ClusterId, err)
			return err
		}
	}

	// register cluster common
	for _, clusterCommon := range clusterWrapper.ClusterCommons {
		clusterCommon.ClusterId = clusterId
		_, err := r.Db.
			InsertInto(models.ClusterCommonTableName).
			Columns(models.ClusterCommonColumns...).
			Record(clusterCommon).
			Exec()
		if err != nil {
			logger.Errorf("Failed to insert table [%s] with cluster id [%s]: %+v",
				models.ClusterCommonTableName, clusterWrapper.Cluster.ClusterId, err)
			return err
		}
	}

	// register cluster link
	for _, clusterLink := range clusterWrapper.ClusterLinks {
		clusterLink.ClusterId = clusterId
		clusterLink.Owner = owner
		_, err := r.Db.
			InsertInto(models.ClusterLinkTableName).
			Columns(models.ClusterLinkColumns...).
			Record(clusterLink).
			Exec()
		if err != nil {
			logger.Errorf("Failed to insert table [%s] with cluster id [%s]: %+v",
				models.ClusterLinkTableName, clusterWrapper.Cluster.ClusterId, err)
			return err
		}
	}

	// register cluster role
	for _, clusterRole := range clusterWrapper.ClusterRoles {
		clusterRole.ClusterId = clusterId
		_, err := r.Db.
			InsertInto(models.ClusterRoleTableName).
			Columns(models.ClusterRoleColumns...).
			Record(clusterRole).
			Exec()
		if err != nil {
			logger.Errorf("Failed to insert table [%s] with cluster id [%s]: %+v",
				models.ClusterRoleTableName, clusterWrapper.Cluster.ClusterId, err)
			return err
		}
	}

	// register cluster loadbalancer
	for _, clusterLoadbalancer := range clusterWrapper.ClusterLoadbalancers {
		clusterLoadbalancer.ClusterId = clusterId
		_, err := r.Db.
			InsertInto(models.ClusterLoadbalancerTableName).
			Columns(models.ClusterLoadbalancerColumns...).
			Record(clusterLoadbalancer).
			Exec()
		if err != nil {
			logger.Errorf("Failed to insert table [%s] with cluster id [%s]: %+v",
				models.ClusterLoadbalancerTableName, clusterWrapper.Cluster.ClusterId, err)
			return err
		}
	}

	return nil
}
