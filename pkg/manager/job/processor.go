// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package job

import (
	clientutil "openpitrix.io/openpitrix/pkg/client"
	clusterclient "openpitrix.io/openpitrix/pkg/client/cluster"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/plugins"
)

type Processor struct {
	Job *models.Job
}

func NewProcessor(job *models.Job) *Processor {
	return &Processor{
		Job: job,
	}
}

// Pre process when job is start
func (j *Processor) Pre() error {
	var err error
	ctx := clientutil.GetSystemUserContext()
	client, err := clusterclient.NewClusterManagerClient(ctx)
	if err != nil {
		logger.Errorf("Executing job [%s] pre processor failed: %+v", j.Job.JobId, err)
		return err
	}
	switch j.Job.JobAction {
	case constants.ActionCreateCluster:
		err = clusterclient.ModifyClusterTransitionStatus(ctx, client, j.Job.ClusterId, constants.StatusCreating)
	case constants.ActionUpgradeCluster:
		err = clusterclient.ModifyClusterTransitionStatus(ctx, client, j.Job.ClusterId, constants.StatusUpgrading)
	case constants.ActionRollbackCluster:
		err = clusterclient.ModifyClusterTransitionStatus(ctx, client, j.Job.ClusterId, constants.StatusRollbacking)
	case constants.ActionResizeCluster:
		err = clusterclient.ModifyClusterTransitionStatus(ctx, client, j.Job.ClusterId, constants.StatusResizing)
	case constants.ActionAddClusterNodes:
		err = clusterclient.ModifyClusterTransitionStatus(ctx, client, j.Job.ClusterId, constants.StatusScaling)
	case constants.ActionDeleteClusterNodes:
		err = clusterclient.ModifyClusterTransitionStatus(ctx, client, j.Job.ClusterId, constants.StatusScaling)
	case constants.ActionStopClusters:
		err = clusterclient.ModifyClusterTransitionStatus(ctx, client, j.Job.ClusterId, constants.StatusStopping)
	case constants.ActionStartClusters:
		err = clusterclient.ModifyClusterTransitionStatus(ctx, client, j.Job.ClusterId, constants.StatusStarting)
	case constants.ActionDeleteClusters:
		err = clusterclient.ModifyClusterTransitionStatus(ctx, client, j.Job.ClusterId, constants.StatusDeleting)
	case constants.ActionRecoverClusters:
		err = clusterclient.ModifyClusterTransitionStatus(ctx, client, j.Job.ClusterId, constants.StatusRecovering)
	case constants.ActionCeaseClusters:
		err = clusterclient.ModifyClusterTransitionStatus(ctx, client, j.Job.ClusterId, constants.StatusCeasing)
	case constants.ActionUpdateClusterEnv:
		err = clusterclient.ModifyClusterTransitionStatus(ctx, client, j.Job.ClusterId, constants.StatusUpdating)
	default:
		logger.Errorf("Unknown job action [%s]", j.Job.JobAction)
	}

	if err != nil {
		logger.Panicf("Executing job [%s] pre processor failed: %+v", j.Job.JobId, err)
	}
	return err
}

// Post process when job is done
func (j *Processor) Post() error {
	var err error
	ctx := clientutil.GetSystemUserContext()
	client, err := clusterclient.NewClusterManagerClient(ctx)
	if err != nil {
		logger.Errorf("Executing job [%s] post processor failed: %+v", j.Job.JobId, err)
		return err
	}
	switch j.Job.JobAction {
	case constants.ActionCreateCluster:
		providerInterface, err := plugins.GetProviderPlugin(j.Job.Provider)
		if err != nil {
			logger.Errorf("No such provider [%s]. ", j.Job.Provider)
			return err
		}
		err = providerInterface.UpdateClusterStatus(j.Job)
		if err != nil {
			logger.Errorf("Executing job [%s] post processor failed: %+v", j.Job.JobId, err)
			return err
		}

		err = clusterclient.ModifyClusterStatus(ctx, client, j.Job.ClusterId, constants.StatusActive)
	case constants.ActionUpgradeCluster:
		providerInterface, err := plugins.GetProviderPlugin(j.Job.Provider)
		if err != nil {
			logger.Errorf("No such provider [%s]. ", j.Job.Provider)
			return err
		}
		err = providerInterface.UpdateClusterStatus(j.Job)
		if err != nil {
			logger.Errorf("Executing job [%s] post processor failed: %+v", j.Job.JobId, err)
			return err
		}

		err = clusterclient.ModifyClusterStatus(ctx, client, j.Job.ClusterId, constants.StatusActive)
	case constants.ActionRollbackCluster:
		providerInterface, err := plugins.GetProviderPlugin(j.Job.Provider)
		if err != nil {
			logger.Errorf("No such provider [%s]. ", j.Job.Provider)
			return err
		}
		err = providerInterface.UpdateClusterStatus(j.Job)
		if err != nil {
			logger.Errorf("Executing job [%s] post processor failed: %+v", j.Job.JobId, err)
			return err
		}

		err = clusterclient.ModifyClusterStatus(ctx, client, j.Job.ClusterId, constants.StatusActive)
	case constants.ActionResizeCluster:
		err = clusterclient.ModifyClusterStatus(ctx, client, j.Job.ClusterId, constants.StatusActive)
	case constants.ActionAddClusterNodes:
		// delete node record from db when pre check is failed
		if j.Job.Status == constants.StatusFailed {
			clusterWrappers, err := clusterclient.GetClusterWrappers(ctx, client, []string{j.Job.ClusterId})
			if err != nil {
				logger.Errorf("No such cluster [%s], %+v ", j.Job.ClusterId, err)
				return err
			}
			clusterWrapper := clusterWrappers[0]
			var deleteNodeIds []string
			for _, clusterNode := range clusterWrapper.ClusterNodes {
				if clusterNode.Status == constants.StatusPending && clusterNode.TransitionStatus == "" {
					deleteNodeIds = append(deleteNodeIds, clusterNode.NodeId)
				}
			}
			if len(deleteNodeIds) > 0 {
				client.DeleteTableClusterNodes(ctx, &pb.DeleteTableClusterNodesRequest{
					NodeId: deleteNodeIds,
				})
			}
		}
		err = clusterclient.ModifyClusterStatus(ctx, client, j.Job.ClusterId, constants.StatusActive)
	case constants.ActionDeleteClusterNodes:
		err = clusterclient.ModifyClusterStatus(ctx, client, j.Job.ClusterId, constants.StatusActive)
	case constants.ActionStopClusters:
		err = clusterclient.ModifyClusterStatus(ctx, client, j.Job.ClusterId, constants.StatusStopped)
	case constants.ActionStartClusters:
		err = clusterclient.ModifyClusterStatus(ctx, client, j.Job.ClusterId, constants.StatusActive)
	case constants.ActionDeleteClusters:
		err = clusterclient.ModifyClusterStatus(ctx, client, j.Job.ClusterId, constants.StatusDeleted)
	case constants.ActionRecoverClusters:
		err = clusterclient.ModifyClusterStatus(ctx, client, j.Job.ClusterId, constants.StatusActive)
	case constants.ActionCeaseClusters:
		err = clusterclient.ModifyClusterStatus(ctx, client, j.Job.ClusterId, constants.StatusCeased)
	case constants.ActionUpdateClusterEnv:
		err = clusterclient.ModifyClusterStatus(ctx, client, j.Job.ClusterId, constants.StatusActive)
	default:
		logger.Errorf("Unknown job action [%s]", j.Job.JobAction)
	}

	if err != nil {
		logger.Errorf("Executing job [%s] post processor failed: %+v", j.Job.JobId, err)
	}
	return err
}

func (j *Processor) Final() {
	ctx := clientutil.GetSystemUserContext()
	client, err := clusterclient.NewClusterManagerClient(ctx)
	if err != nil {
		logger.Errorf("Executing job [%s] final processor failed: %+v", j.Job.JobId, err)
		return
	}
	// TODO: modify cluster status to `active or deleted`
	err = clusterclient.ModifyClusterTransitionStatus(ctx, client, j.Job.ClusterId, "")
	if err != nil {
		logger.Errorf("Executing job [%s] final processor failed: %+v", j.Job.JobId, err)
	}
}
