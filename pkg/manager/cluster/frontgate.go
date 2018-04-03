// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package cluster

import (
	"fmt"
	"time"

	runtimeenvclient "openpitrix.io/openpitrix/pkg/client/runtimeenv"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pi"
)

type Frontgate struct {
	*pi.Pi
	Runtime *runtimeenvclient.Runtime
}

func (f *Frontgate) getFrontgateFromDb(vpcId, userId string) ([]*models.Cluster, error) {
	var frontgates []*models.Cluster
	statuses := []string{constants.StatusActive, constants.StatusPending,
		constants.StatusDeleted, constants.StatusStopped}
	_, err := f.Db.
		Select(models.ClusterColumns...).
		From(models.ClusterTableName).
		Where(db.Eq("vpc_id", vpcId)).
		Where(db.Eq("owner", userId)).
		Where(db.Eq("cluster_type", constants.FrontgateClusterType)).
		Where(db.Eq("status", statuses)).
		Load(&frontgates)
	if err != nil {
		return nil, err
	}
	return frontgates, nil
}

func (f *Frontgate) activate(frontgate *models.Cluster) error {
	if frontgate.TransitionStatus != "" {
		logger.Warnf("Frontgate cluster [%s] is in [%s] transition status, please try laster",
			frontgate.ClusterId, frontgate.TransitionStatus)
		return fmt.Errorf("Frontgate service is [%s], please try later. ", frontgate.TransitionStatus)
	}

	if frontgate.Status == constants.StatusActive {
		return nil
	} else if frontgate.Status == constants.StatusStopped {
		return f.StartCluster(frontgate)
	} else if frontgate.Status == constants.StatusDeleted {
		return f.RecoverCluster(frontgate)
	} else {
		logger.Panicf("Frontgate cluster [%s] is in wrong status [%s]", frontgate.ClusterId, frontgate.Status)
		return fmt.Errorf("Frontgate cluster is in wrong status [%s]. ", frontgate.Status)
	}
}

func (f *Frontgate) GetFrontgate(frontgateId string) (*models.Cluster, error) {
	var frontgate *models.Cluster
	err := f.Db.
		Select(models.ClusterColumns...).
		From(models.ClusterTableName).
		Where(db.Eq("cluster_id", frontgateId)).
		LoadOne(&frontgate)
	if err != nil {
		return nil, err
	}
	return frontgate, nil
}

func (f *Frontgate) ActivateFrontgate(frontgateId string) error {
	frontgate, err := f.GetFrontgate(frontgateId)
	if err != nil {
		return err
	}

	return f.activate(frontgate)
}

func (f *Frontgate) GetActiveFrontgate(vpcId, userId string, register *Register) (*models.Cluster, error) {
	var frontgate *models.Cluster
	err := f.Etcd.DlockWithTimeout(constants.ClusterPrefix+vpcId, 600*time.Second, func() error {
		// Check vpc status
		vpc, err := f.Runtime.ProviderInterface.DescribeVpc(vpcId)
		if err != nil {
			return err
		}
		if vpc == nil {
			return fmt.Errorf("Describe vpc [%s] failed. ", vpcId)
		}
		if vpc.Status != constants.StatusActive {
			logger.Warnf("Vpc [%s] is not active. ", vpcId)
			return fmt.Errorf("Vpc [%s] is not active. ", vpcId)
		}
		if vpc.TransitionStatus != "" {
			logger.Warnf("Vpc [%s] is now updating. ", vpcId)
			return fmt.Errorf("Vpc [%s] is now updating. ", vpcId)
		}

		frontgates, err := f.getFrontgateFromDb(vpcId, userId)
		if err != nil {
			return err
		}
		if len(frontgates) == 0 {
			frontgateId, err := f.CreateCluster(register)
			frontgate = &models.Cluster{ClusterId: frontgateId}
			return err
		} else if len(frontgates) == 1 {
			frontgate = frontgates[0]
			err = f.activate(frontgate)
			return err
		} else {
			logger.Panicf("More than one frontgate non-ceased cluster in the vpc [%s] for user [%s]", vpcId, userId)
			return fmt.Errorf("More than one frontgate non-ceased cluster. ")
		}
	})
	if err != nil {
		return nil, err
	}

	return frontgate, nil
}
