// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package job

import (
	clusterclient "openpitrix.io/openpitrix/pkg/client/cluster"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/utils"
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
func (j *Processor) Pre() {
	var err error
	switch j.Job.JobAction {
	case constants.ActionCreateCluster:
		err = clusterclient.ModifyCluster(&pb.ModifyClusterRequest{
			ClusterId:        utils.ToProtoString(j.Job.ClusterId),
			TransitionStatus: utils.ToProtoString(constants.StatusCreating),
		})
	case constants.ActionUpgradeCluster:
		err = clusterclient.ModifyCluster(&pb.ModifyClusterRequest{
			ClusterId:        utils.ToProtoString(j.Job.ClusterId),
			TransitionStatus: utils.ToProtoString(constants.StatusUpgrading),
		})
	case constants.ActionRollbackCluster:
		err = clusterclient.ModifyCluster(&pb.ModifyClusterRequest{
			ClusterId:        utils.ToProtoString(j.Job.ClusterId),
			TransitionStatus: utils.ToProtoString(constants.StatusRollbacking),
		})
	case constants.ActionResizeCluster:
		err = clusterclient.ModifyCluster(&pb.ModifyClusterRequest{
			ClusterId:        utils.ToProtoString(j.Job.ClusterId),
			TransitionStatus: utils.ToProtoString(constants.StatusResizing),
		})
	case constants.ActionAddClusterNodes:
		err = clusterclient.ModifyCluster(&pb.ModifyClusterRequest{
			ClusterId:        utils.ToProtoString(j.Job.ClusterId),
			TransitionStatus: utils.ToProtoString(constants.StatusScaling),
		})
	case constants.ActionDeleteClusterNodes:
		err = clusterclient.ModifyCluster(&pb.ModifyClusterRequest{
			ClusterId:        utils.ToProtoString(j.Job.ClusterId),
			TransitionStatus: utils.ToProtoString(constants.StatusScaling),
		})
	case constants.ActionStopClusters:
		err = clusterclient.ModifyCluster(&pb.ModifyClusterRequest{
			ClusterId:        utils.ToProtoString(j.Job.ClusterId),
			TransitionStatus: utils.ToProtoString(constants.StatusStopping),
		})
	case constants.ActionStartClusters:
		err = clusterclient.ModifyCluster(&pb.ModifyClusterRequest{
			ClusterId:        utils.ToProtoString(j.Job.ClusterId),
			TransitionStatus: utils.ToProtoString(constants.StatusStarting),
		})
	case constants.ActionDeleteClusters:
		err = clusterclient.ModifyCluster(&pb.ModifyClusterRequest{
			ClusterId:        utils.ToProtoString(j.Job.ClusterId),
			TransitionStatus: utils.ToProtoString(constants.StatusDeleting),
		})
	case constants.ActionRecoverClusters:
		err = clusterclient.ModifyCluster(&pb.ModifyClusterRequest{
			ClusterId:        utils.ToProtoString(j.Job.ClusterId),
			TransitionStatus: utils.ToProtoString(constants.StatusRecovering),
		})
	case constants.ActionCeaseClusters:
		err = clusterclient.ModifyCluster(&pb.ModifyClusterRequest{
			ClusterId:        utils.ToProtoString(j.Job.ClusterId),
			TransitionStatus: utils.ToProtoString(constants.StatusCeasing),
		})
	case constants.ActionUpdateClusterEnv:
		err = clusterclient.ModifyCluster(&pb.ModifyClusterRequest{
			ClusterId:        utils.ToProtoString(j.Job.ClusterId),
			TransitionStatus: utils.ToProtoString(constants.StatusUpdating),
		})
	default:
		logger.Errorf("Unknown job action [%s]", j.Job.JobAction)
	}

	if err != nil {
		logger.Panicf("Executing job [%s] pre processor failed: %+v", j.Job.JobId, err)
	}
}

// Post process when job is done
func (j *Processor) Post() {
	var err error
	switch j.Job.JobAction {
	case constants.ActionCreateCluster:
		err = clusterclient.ModifyCluster(&pb.ModifyClusterRequest{
			ClusterId:        utils.ToProtoString(j.Job.ClusterId),
			Status:           utils.ToProtoString(constants.StatusActive),
			TransitionStatus: utils.ToProtoString(""),
		})
	case constants.ActionUpgradeCluster:
		err = clusterclient.ModifyCluster(&pb.ModifyClusterRequest{
			ClusterId:        utils.ToProtoString(j.Job.ClusterId),
			Status:           utils.ToProtoString(constants.StatusActive),
			TransitionStatus: utils.ToProtoString(""),
		})
	case constants.ActionRollbackCluster:
		err = clusterclient.ModifyCluster(&pb.ModifyClusterRequest{
			ClusterId:        utils.ToProtoString(j.Job.ClusterId),
			Status:           utils.ToProtoString(constants.StatusActive),
			TransitionStatus: utils.ToProtoString(""),
		})
	case constants.ActionResizeCluster:
		err = clusterclient.ModifyCluster(&pb.ModifyClusterRequest{
			ClusterId:        utils.ToProtoString(j.Job.ClusterId),
			Status:           utils.ToProtoString(constants.StatusActive),
			TransitionStatus: utils.ToProtoString(""),
		})
	case constants.ActionAddClusterNodes:
		err = clusterclient.ModifyCluster(&pb.ModifyClusterRequest{
			ClusterId:        utils.ToProtoString(j.Job.ClusterId),
			Status:           utils.ToProtoString(constants.StatusActive),
			TransitionStatus: utils.ToProtoString(""),
		})
	case constants.ActionDeleteClusterNodes:
		err = clusterclient.ModifyCluster(&pb.ModifyClusterRequest{
			ClusterId:        utils.ToProtoString(j.Job.ClusterId),
			Status:           utils.ToProtoString(constants.StatusActive),
			TransitionStatus: utils.ToProtoString(""),
		})
	case constants.ActionStopClusters:
		err = clusterclient.ModifyCluster(&pb.ModifyClusterRequest{
			ClusterId:        utils.ToProtoString(j.Job.ClusterId),
			Status:           utils.ToProtoString(constants.StatusStopped),
			TransitionStatus: utils.ToProtoString(""),
		})
	case constants.ActionStartClusters:
		err = clusterclient.ModifyCluster(&pb.ModifyClusterRequest{
			ClusterId:        utils.ToProtoString(j.Job.ClusterId),
			Status:           utils.ToProtoString(constants.StatusActive),
			TransitionStatus: utils.ToProtoString(""),
		})
	case constants.ActionDeleteClusters:
		err = clusterclient.ModifyCluster(&pb.ModifyClusterRequest{
			ClusterId:        utils.ToProtoString(j.Job.ClusterId),
			Status:           utils.ToProtoString(constants.StatusDeleted),
			TransitionStatus: utils.ToProtoString(""),
		})
	case constants.ActionRecoverClusters:
		err = clusterclient.ModifyCluster(&pb.ModifyClusterRequest{
			ClusterId:        utils.ToProtoString(j.Job.ClusterId),
			Status:           utils.ToProtoString(constants.StatusActive),
			TransitionStatus: utils.ToProtoString(""),
		})
	case constants.ActionCeaseClusters:
		err = clusterclient.ModifyCluster(&pb.ModifyClusterRequest{
			ClusterId:        utils.ToProtoString(j.Job.ClusterId),
			Status:           utils.ToProtoString(constants.StatusCeased),
			TransitionStatus: utils.ToProtoString(""),
		})
	case constants.ActionUpdateClusterEnv:
		err = clusterclient.ModifyCluster(&pb.ModifyClusterRequest{
			ClusterId:        utils.ToProtoString(j.Job.ClusterId),
			Status:           utils.ToProtoString(constants.StatusActive),
			TransitionStatus: utils.ToProtoString(""),
		})
	default:
		logger.Errorf("Unknown job action [%s]", j.Job.JobAction)
	}

	if err != nil {
		logger.Panicf("Executing job [%s] post processor failed: %+v", j.Job.JobId, err)
	}
}
