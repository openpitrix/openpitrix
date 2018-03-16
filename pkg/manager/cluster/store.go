package cluster

import (
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pi"
)

type Store struct {
	*pi.Pi
}

func (s *Store) RegisterClusterWrapper(clusterWrapper *models.ClusterWrapper) error {
	if clusterWrapper.Cluster != nil {
		_, err := s.Db.
			InsertInto(models.ClusterTableName).
			Columns(models.ClusterColumns...).
			Record(clusterWrapper.Cluster).
			Exec()
		if err != nil {
			logger.Errorf("Failed to insert cluster table with cluster id [%s]: %+v",
				clusterWrapper.Cluster.ClusterId, err)
			return err
		}
	}

	for _, clusterNode := range clusterWrapper.ClusterNodes {
		_, err := s.Db.
			InsertInto(models.ClusterNodeTableName).
			Columns(models.ClusterNodeColumns...).
			Record(clusterNode).
			Exec()
		if err != nil {
			logger.Errorf("Failed to insert cluster_node table with cluster id [%s]: %+v",
				clusterWrapper.Cluster.ClusterId, err)
			return err
		}
	}

	for _, clusterCommon := range clusterWrapper.ClusterCommons {
		_, err := s.Db.
			InsertInto(models.ClusterCommonTableName).
			Columns(models.ClusterCommonColumns...).
			Record(clusterCommon).
			Exec()
		if err != nil {
			logger.Errorf("Failed to insert cluster_common table with cluster id [%s]: %+v",
				clusterWrapper.Cluster.ClusterId, err)
			return err
		}
	}

	for _, clusterLink := range clusterWrapper.ClusterLinks {
		_, err := s.Db.
			InsertInto(models.ClusterLinkTableName).
			Columns(models.ClusterLinkColumns...).
			Record(clusterLink).
			Exec()
		if err != nil {
			logger.Errorf("Failed to insert cluster_link table with cluster id [%s]: %+v",
				clusterWrapper.Cluster.ClusterId, err)
			return err
		}
	}

	for _, clusterRole := range clusterWrapper.ClusterRoles {
		_, err := s.Db.
			InsertInto(models.ClusterRoleTableName).
			Columns(models.ClusterRoleColumns...).
			Record(clusterRole).
			Exec()
		if err != nil {
			logger.Errorf("Failed to insert cluster_role table with cluster id [%s]: %+v",
				clusterWrapper.Cluster.ClusterId, err)
			return err
		}
	}

	return nil
}
