// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package cluster

import (
	"fmt"
	"time"

	runtimeclient "openpitrix.io/openpitrix/pkg/client/runtime"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/gerr"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pi"
	"openpitrix.io/openpitrix/pkg/plugins"
)

type Frontgate struct {
	Runtime *runtimeclient.Runtime
}

func (f *Frontgate) getFrontgateFromDb(vpcId, userId string) ([]*models.Cluster, error) {
	var frontgates []*models.Cluster
	statuses := []string{constants.StatusActive, constants.StatusPending, constants.StatusStopped}
	_, err := pi.Global().Db.
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
		logger.Warn("Frontgate cluster [%s] is in [%s] transition status, please try laster",
			frontgate.ClusterId, frontgate.TransitionStatus)
		err := fmt.Errorf("frontgate service is [%s], please try later", frontgate.TransitionStatus)
		return gerr.NewWithDetail(gerr.PermissionDenied, err, gerr.ErrorResourceTransitionStatus, "frontgate", constants.StatusUpdating)
	}

	if frontgate.Status == constants.StatusActive {
		return nil
	} else if frontgate.Status == constants.StatusStopped {
		err := f.StartCluster(frontgate)
		if err != nil {
			return gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorStartResourceFailed, frontgate.ClusterId)
		}
		return nil
	} else {
		err := fmt.Errorf("frontgate cluster [%s] is in wrong status [%s]", frontgate.ClusterId, frontgate.Status)
		return gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorResourceTransitionStatus, frontgate.ClusterId, frontgate.Status)
	}
}

func (f *Frontgate) GetFrontgate(frontgateId string) (*models.Cluster, error) {
	var frontgate *models.Cluster
	err := pi.Global().Db.
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

func (f *Frontgate) GetActiveFrontgate(clusterWrapper *models.ClusterWrapper) (*models.Cluster, error) {
	var frontgate *models.Cluster
	vpcId := clusterWrapper.Cluster.VpcId
	owner := clusterWrapper.Cluster.Owner
	err := pi.Global().Etcd.DlockWithTimeout(constants.ClusterPrefix+vpcId, 600*time.Second, func() error {
		// Check vpc status
		providerInterface, err := plugins.GetProviderPlugin(f.Runtime.Provider, nil)
		if err != nil {
			return gerr.NewWithDetail(gerr.NotFound, err, gerr.ErrorProviderNotFound, f.Runtime.Provider)
		}
		vpc, err := providerInterface.DescribeVpc(f.Runtime.RuntimeId, vpcId)
		if err != nil {
			return gerr.NewWithDetail(gerr.PermissionDenied, err, gerr.ErrorResourceNotFound, vpcId)
		}
		if vpc == nil {
			err = fmt.Errorf("describe vpc [%s] failed", vpcId)
			return gerr.NewWithDetail(gerr.PermissionDenied, err, gerr.ErrorDescribeResourceFailed, vpcId)
		}
		if vpc.Status != constants.StatusActive {
			err = fmt.Errorf("vpc [%s] is not active", vpcId)
			return gerr.NewWithDetail(gerr.PermissionDenied, err, gerr.ErrorResourceNotInStatus, vpcId, constants.StatusActive)
		}
		if vpc.TransitionStatus != "" {
			err = fmt.Errorf("vpc [%s] is now updating", vpcId)
			return gerr.NewWithDetail(gerr.PermissionDenied, err, gerr.ErrorResourceTransitionStatus, vpcId, constants.StatusUpdating)
		}

		frontgates, err := f.getFrontgateFromDb(vpcId, owner)
		if err != nil {
			return gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorDescribeResourceFailed, vpcId)
		}
		if len(frontgates) == 0 {
			frontgateId, err := f.CreateCluster(clusterWrapper)
			frontgate = &models.Cluster{ClusterId: frontgateId}
			if err != nil {
				return gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorCreateResourceFailed, frontgateId)
			}
			return nil
		} else if len(frontgates) == 1 {
			frontgate = frontgates[0]
			err = f.activate(frontgate)
			return err
		} else {
			logger.Critical("More than one non-ceased frontgate cluster in the vpc [%s] for user [%s]", vpcId, owner)
			err = fmt.Errorf("more than one non-ceased frontgate cluster")
			return gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorInternalError)
		}
	})
	if err != nil {
		return nil, err
	}

	return frontgate, nil
}
