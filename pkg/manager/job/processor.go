// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package job

import (
	clusterClient "openpitrix.io/openpitrix/pkg/client/cluster"
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

// Post process when job is done
func (j *Processor) Post() {
	var err error
	switch j.Job.JobAction {
	case constants.ActionCreateCluster:
		err = clusterClient.ModifyCluster(&pb.ModifyClusterRequest{
			ClusterId:        utils.ToProtoString(j.Job.ClusterId),
			Status:           utils.ToProtoString(constants.StatusActive),
			TransitionStatus: utils.ToProtoString(""),
		})
	case constants.ActionUpgradeCluster:
		err = clusterClient.ModifyCluster(&pb.ModifyClusterRequest{
			ClusterId:        utils.ToProtoString(j.Job.ClusterId),
			Status:           utils.ToProtoString(constants.StatusActive),
			TransitionStatus: utils.ToProtoString(""),
		})
	case constants.ActionRollbackCluster:
		err = clusterClient.ModifyCluster(&pb.ModifyClusterRequest{
			ClusterId:        utils.ToProtoString(j.Job.ClusterId),
			Status:           utils.ToProtoString(constants.StatusActive),
			TransitionStatus: utils.ToProtoString(""),
		})
	case constants.ActionResizeCluster:
		err = clusterClient.ModifyCluster(&pb.ModifyClusterRequest{
			ClusterId:        utils.ToProtoString(j.Job.ClusterId),
			Status:           utils.ToProtoString(constants.StatusActive),
			TransitionStatus: utils.ToProtoString(""),
		})
	case constants.ActionAddClusterNodes:
		err = clusterClient.ModifyCluster(&pb.ModifyClusterRequest{
			ClusterId:        utils.ToProtoString(j.Job.ClusterId),
			Status:           utils.ToProtoString(constants.StatusActive),
			TransitionStatus: utils.ToProtoString(""),
		})
	case constants.ActionDeleteClusterNodes:
		err = clusterClient.ModifyCluster(&pb.ModifyClusterRequest{
			ClusterId:        utils.ToProtoString(j.Job.ClusterId),
			Status:           utils.ToProtoString(constants.StatusActive),
			TransitionStatus: utils.ToProtoString(""),
		})
	case constants.ActionStopClusters:
		err = clusterClient.ModifyCluster(&pb.ModifyClusterRequest{
			ClusterId:        utils.ToProtoString(j.Job.ClusterId),
			Status:           utils.ToProtoString(constants.StatusStopped),
			TransitionStatus: utils.ToProtoString(""),
		})
	case constants.ActionStartClusters:
		err = clusterClient.ModifyCluster(&pb.ModifyClusterRequest{
			ClusterId:        utils.ToProtoString(j.Job.ClusterId),
			Status:           utils.ToProtoString(constants.StatusActive),
			TransitionStatus: utils.ToProtoString(""),
		})
	case constants.ActionDeleteClusters:
		err = clusterClient.ModifyCluster(&pb.ModifyClusterRequest{
			ClusterId:        utils.ToProtoString(j.Job.ClusterId),
			Status:           utils.ToProtoString(constants.StatusDeleted),
			TransitionStatus: utils.ToProtoString(""),
		})
	case constants.ActionRecoverClusters:
		err = clusterClient.ModifyCluster(&pb.ModifyClusterRequest{
			ClusterId:        utils.ToProtoString(j.Job.ClusterId),
			Status:           utils.ToProtoString(constants.StatusActive),
			TransitionStatus: utils.ToProtoString(""),
		})
	case constants.ActionCeaseClusters:
		err = clusterClient.ModifyCluster(&pb.ModifyClusterRequest{
			ClusterId:        utils.ToProtoString(j.Job.ClusterId),
			Status:           utils.ToProtoString(constants.StatusCeased),
			TransitionStatus: utils.ToProtoString(""),
		})
	case constants.ActionUpdateClusterEnv:
		err = clusterClient.ModifyCluster(&pb.ModifyClusterRequest{
			ClusterId:        utils.ToProtoString(j.Job.ClusterId),
			Status:           utils.ToProtoString(constants.StatusActive),
			TransitionStatus: utils.ToProtoString(""),
		})
	}

	if err != nil {
		logger.Errorf("Executing job [%s] post processor failed: %+v", j.Job.JobId, err)
	}
}
