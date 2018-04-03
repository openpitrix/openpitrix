// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package job

import (
	"context"

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
func (j *Processor) Pre() error {
	var err error
	ctx := context.Background()
	client, err := clusterclient.NewClusterManagerClient(ctx)
	if err != nil {
		logger.Errorf("Executing job [%s] pre processor failed: %+v", j.Job.JobId, err)
		return err
	}
	switch j.Job.JobAction {
	case constants.ActionCreateCluster:
		_, err = client.ModifyCluster(ctx, &pb.ModifyClusterRequest{
			ClusterId:        utils.ToProtoString(j.Job.ClusterId),
			TransitionStatus: utils.ToProtoString(constants.StatusCreating),
		})
	case constants.ActionUpgradeCluster:
		_, err = client.ModifyCluster(ctx, &pb.ModifyClusterRequest{
			ClusterId:        utils.ToProtoString(j.Job.ClusterId),
			TransitionStatus: utils.ToProtoString(constants.StatusUpgrading),
		})
	case constants.ActionRollbackCluster:
		_, err = client.ModifyCluster(ctx, &pb.ModifyClusterRequest{
			ClusterId:        utils.ToProtoString(j.Job.ClusterId),
			TransitionStatus: utils.ToProtoString(constants.StatusRollbacking),
		})
	case constants.ActionResizeCluster:
		_, err = client.ModifyCluster(ctx, &pb.ModifyClusterRequest{
			ClusterId:        utils.ToProtoString(j.Job.ClusterId),
			TransitionStatus: utils.ToProtoString(constants.StatusResizing),
		})
	case constants.ActionAddClusterNodes:
		_, err = client.ModifyCluster(ctx, &pb.ModifyClusterRequest{
			ClusterId:        utils.ToProtoString(j.Job.ClusterId),
			TransitionStatus: utils.ToProtoString(constants.StatusScaling),
		})
	case constants.ActionDeleteClusterNodes:
		_, err = client.ModifyCluster(ctx, &pb.ModifyClusterRequest{
			ClusterId:        utils.ToProtoString(j.Job.ClusterId),
			TransitionStatus: utils.ToProtoString(constants.StatusScaling),
		})
	case constants.ActionStopClusters:
		_, err = client.ModifyCluster(ctx, &pb.ModifyClusterRequest{
			ClusterId:        utils.ToProtoString(j.Job.ClusterId),
			TransitionStatus: utils.ToProtoString(constants.StatusStopping),
		})
	case constants.ActionStartClusters:
		_, err = client.ModifyCluster(ctx, &pb.ModifyClusterRequest{
			ClusterId:        utils.ToProtoString(j.Job.ClusterId),
			TransitionStatus: utils.ToProtoString(constants.StatusStarting),
		})
	case constants.ActionDeleteClusters:
		_, err = client.ModifyCluster(ctx, &pb.ModifyClusterRequest{
			ClusterId:        utils.ToProtoString(j.Job.ClusterId),
			TransitionStatus: utils.ToProtoString(constants.StatusDeleting),
		})
	case constants.ActionRecoverClusters:
		_, err = client.ModifyCluster(ctx, &pb.ModifyClusterRequest{
			ClusterId:        utils.ToProtoString(j.Job.ClusterId),
			TransitionStatus: utils.ToProtoString(constants.StatusRecovering),
		})
	case constants.ActionCeaseClusters:
		_, err = client.ModifyCluster(ctx, &pb.ModifyClusterRequest{
			ClusterId:        utils.ToProtoString(j.Job.ClusterId),
			TransitionStatus: utils.ToProtoString(constants.StatusCeasing),
		})
	case constants.ActionUpdateClusterEnv:
		_, err = client.ModifyCluster(ctx, &pb.ModifyClusterRequest{
			ClusterId:        utils.ToProtoString(j.Job.ClusterId),
			TransitionStatus: utils.ToProtoString(constants.StatusUpdating),
		})
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
	ctx := context.Background()
	client, err := clusterclient.NewClusterManagerClient(ctx)
	if err != nil {
		logger.Errorf("Executing job [%s] post processor failed: %+v", j.Job.JobId, err)
		return err
	}
	switch j.Job.JobAction {
	case constants.ActionCreateCluster:
		_, err = client.ModifyCluster(ctx, &pb.ModifyClusterRequest{
			ClusterId: utils.ToProtoString(j.Job.ClusterId),
			Status:    utils.ToProtoString(constants.StatusActive),
		})
	case constants.ActionUpgradeCluster:
		_, err = client.ModifyCluster(ctx, &pb.ModifyClusterRequest{
			ClusterId: utils.ToProtoString(j.Job.ClusterId),
			Status:    utils.ToProtoString(constants.StatusActive),
		})
	case constants.ActionRollbackCluster:
		_, err = client.ModifyCluster(ctx, &pb.ModifyClusterRequest{
			ClusterId: utils.ToProtoString(j.Job.ClusterId),
			Status:    utils.ToProtoString(constants.StatusActive),
		})
	case constants.ActionResizeCluster:
		_, err = client.ModifyCluster(ctx, &pb.ModifyClusterRequest{
			ClusterId: utils.ToProtoString(j.Job.ClusterId),
			Status:    utils.ToProtoString(constants.StatusActive),
		})
	case constants.ActionAddClusterNodes:
		_, err = client.ModifyCluster(ctx, &pb.ModifyClusterRequest{
			ClusterId: utils.ToProtoString(j.Job.ClusterId),
			Status:    utils.ToProtoString(constants.StatusActive),
		})
	case constants.ActionDeleteClusterNodes:
		_, err = client.ModifyCluster(ctx, &pb.ModifyClusterRequest{
			ClusterId: utils.ToProtoString(j.Job.ClusterId),
			Status:    utils.ToProtoString(constants.StatusActive),
		})
	case constants.ActionStopClusters:
		_, err = client.ModifyCluster(ctx, &pb.ModifyClusterRequest{
			ClusterId: utils.ToProtoString(j.Job.ClusterId),
			Status:    utils.ToProtoString(constants.StatusStopped),
		})
	case constants.ActionStartClusters:
		_, err = client.ModifyCluster(ctx, &pb.ModifyClusterRequest{
			ClusterId: utils.ToProtoString(j.Job.ClusterId),
			Status:    utils.ToProtoString(constants.StatusActive),
		})
	case constants.ActionDeleteClusters:
		_, err = client.ModifyCluster(ctx, &pb.ModifyClusterRequest{
			ClusterId: utils.ToProtoString(j.Job.ClusterId),
			Status:    utils.ToProtoString(constants.StatusDeleted),
		})
	case constants.ActionRecoverClusters:
		_, err = client.ModifyCluster(ctx, &pb.ModifyClusterRequest{
			ClusterId: utils.ToProtoString(j.Job.ClusterId),
			Status:    utils.ToProtoString(constants.StatusActive),
		})
	case constants.ActionCeaseClusters:
		_, err = client.ModifyCluster(ctx, &pb.ModifyClusterRequest{
			ClusterId: utils.ToProtoString(j.Job.ClusterId),
			Status:    utils.ToProtoString(constants.StatusCeased),
		})
	case constants.ActionUpdateClusterEnv:
		_, err = client.ModifyCluster(ctx, &pb.ModifyClusterRequest{
			ClusterId: utils.ToProtoString(j.Job.ClusterId),
			Status:    utils.ToProtoString(constants.StatusActive),
		})
	default:
		logger.Errorf("Unknown job action [%s]", j.Job.JobAction)
	}

	if err != nil {
		logger.Errorf("Executing job [%s] post processor failed: %+v", j.Job.JobId, err)
	}
	return err
}

func (j *Processor) Final() {
	ctx := context.Background()
	client, err := clusterclient.NewClusterManagerClient(ctx)
	if err != nil {
		logger.Errorf("Executing job [%s] final processor failed: %+v", j.Job.JobId, err)
		return
	}
	_, err = client.ModifyCluster(ctx, &pb.ModifyClusterRequest{
		ClusterId:        utils.ToProtoString(j.Job.ClusterId),
		TransitionStatus: utils.ToProtoString(""),
	})
	if err != nil {
		logger.Errorf("Executing job [%s] final processor failed: %+v", j.Job.JobId, err)
	}
}
