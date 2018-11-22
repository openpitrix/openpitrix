// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package cluster

import (
	"context"
	"fmt"
	"strings"
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

func (f *Frontgate) getFrontgateFromDb(ctx context.Context, vpcId, userId string) ([]*models.Cluster, error) {
	var frontgates []*models.Cluster
	statuses := []string{constants.StatusActive, constants.StatusPending, constants.StatusStopped}
	_, err := pi.Global().DB(ctx).
		Select(models.ClusterColumns...).
		From(constants.TableCluster).
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

func (f *Frontgate) activate(ctx context.Context, frontgate *models.Cluster) error {
	if frontgate.TransitionStatus != "" {
		logger.Warn(ctx, "Frontgate cluster [%s] is in [%s] transition status, please try laster",
			frontgate.ClusterId, frontgate.TransitionStatus)
		err := fmt.Errorf("frontgate service is [%s], please try later", frontgate.TransitionStatus)
		return gerr.NewWithDetail(ctx, gerr.PermissionDenied, err, gerr.ErrorResourceTransitionStatus, "frontgate", constants.StatusUpdating)
	}

	if frontgate.Status == constants.StatusActive {
		return nil
	} else if frontgate.Status == constants.StatusStopped {
		err := f.StartCluster(ctx, frontgate)
		if err != nil {
			return gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorStartResourceFailed, frontgate.ClusterId)
		}
		return nil
	} else {
		err := fmt.Errorf("frontgate cluster [%s] is in wrong status [%s]", frontgate.ClusterId, frontgate.Status)
		return gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorResourceTransitionStatus, frontgate.ClusterId, frontgate.Status)
	}
}

func (f *Frontgate) GetFrontgate(ctx context.Context, frontgateId string) (*models.Cluster, error) {
	var frontgate *models.Cluster
	err := pi.Global().DB(ctx).
		Select(models.ClusterColumns...).
		From(constants.TableCluster).
		Where(db.Eq("cluster_id", frontgateId)).
		LoadOne(&frontgate)
	if err != nil {
		return nil, err
	}
	return frontgate, nil
}

func (f *Frontgate) ActivateFrontgate(ctx context.Context, frontgateId string) error {
	frontgate, err := f.GetFrontgate(ctx, frontgateId)
	if err != nil {
		return err
	}

	return f.activate(ctx, frontgate)
}

func (f *Frontgate) GetActiveFrontgate(ctx context.Context, clusterWrapper *models.ClusterWrapper) (*models.Cluster, error) {
	var frontgate *models.Cluster
	vpcId := clusterWrapper.Cluster.VpcId
	owner := clusterWrapper.Cluster.Owner
	err := pi.Global().Etcd(ctx).DlockWithTimeout(constants.ClusterPrefix+vpcId, 600*time.Second, func() error {
		// Check vpc status
		providerInterface, err := plugins.GetProviderPlugin(ctx, f.Runtime.Provider)
		if err != nil {
			return gerr.NewWithDetail(ctx, gerr.NotFound, err, gerr.ErrorProviderNotFound, f.Runtime.Provider)
		}
		vpc, err := providerInterface.DescribeVpc(ctx, f.Runtime.RuntimeId, vpcId)
		if err != nil {
			return gerr.NewWithDetail(ctx, gerr.PermissionDenied, err, gerr.ErrorResourceNotFound, vpcId)
		}
		if vpc == nil {
			err = fmt.Errorf("describe vpc [%s] failed", vpcId)
			return gerr.NewWithDetail(ctx, gerr.PermissionDenied, err, gerr.ErrorDescribeResourceFailed, vpcId)
		}
		if vpc.Status != constants.StatusActive && vpc.Status != constants.StatusAvailable && vpc.Status != strings.Title(constants.StatusAvailable) {
			err = fmt.Errorf("vpc [%s] is not active or available", vpcId)
			return gerr.NewWithDetail(ctx, gerr.PermissionDenied, err, gerr.ErrorResourceNotInStatus, vpcId, constants.StatusActive)
		}
		if vpc.TransitionStatus != "" {
			err = fmt.Errorf("vpc [%s] is now updating", vpcId)
			return gerr.NewWithDetail(ctx, gerr.PermissionDenied, err, gerr.ErrorResourceTransitionStatus, vpcId, constants.StatusUpdating)
		}

		frontgates, err := f.getFrontgateFromDb(ctx, vpcId, owner)
		if err != nil {
			return gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourceFailed, vpcId)
		}
		if len(frontgates) == 0 {
			frontgateId, err := f.CreateCluster(ctx, clusterWrapper)
			frontgate = &models.Cluster{ClusterId: frontgateId}
			if err != nil {
				return gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorCreateResourceFailed, frontgateId)
			}
			return nil
		} else if len(frontgates) == 1 {
			frontgate = frontgates[0]
			err = f.activate(ctx, frontgate)
			return err
		} else {
			logger.Critical(ctx, "More than one non-ceased frontgate cluster in the vpc [%s] for user [%s]", vpcId, owner)
			err = fmt.Errorf("more than one non-ceased frontgate cluster")
			return gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
		}
	})
	if err != nil {
		return nil, err
	}

	return frontgate, nil
}
