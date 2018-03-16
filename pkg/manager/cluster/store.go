package cluster

import (
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pi"
)

type Store struct {
	*pi.Pi
}

func (s *Store) RegisterCluster(cluster *models.Cluster) error {
	// TODO: Need to insert into multiple tables
	_, err := s.Db.
		InsertInto(models.ClusterTableName).
		Columns(models.ClusterColumns...).
		Record(cluster).
		Exec()
	if err != nil {
		logger.Errorf("Failed to insert cluster [%s]: %+v", cluster.ClusterId, err)
		return err
	}
	return nil
}
