// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package job

import (
	"openpitrix.io/openpitrix/pkg/client"
	clusterclient "openpitrix.io/openpitrix/pkg/client/cluster"
	jobclient "openpitrix.io/openpitrix/pkg/client/job"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/plugins"
	"openpitrix.io/openpitrix/pkg/util/jsonutil"
)

type Processor struct {
	Job     *models.Job
	JLogger *logger.Logger
}

func NewProcessor(job *models.Job, jLogger *logger.Logger) *Processor {
	if jLogger == nil {
		jLogger = logger.NewLogger()
	}
	return &Processor{
		Job:     job,
		JLogger: jLogger,
	}
}

// Pre process when job is start
func (j *Processor) Pre() error {
	var err error
	ctx := client.GetSystemUserContext()
	clusterClient, err := clusterclient.NewClient()
	if err != nil {
		j.JLogger.Error("Executing job pre processor failed: %+v", err)
		return err
	}
	switch j.Job.JobAction {
	case constants.ActionCreateCluster:
		err = clusterClient.ModifyClusterTransitionStatus(ctx, j.Job.ClusterId, constants.StatusCreating)
	case constants.ActionUpgradeCluster:
		err = clusterClient.ModifyClusterTransitionStatus(ctx, j.Job.ClusterId, constants.StatusUpgrading)
	case constants.ActionRollbackCluster:
		err = clusterClient.ModifyClusterTransitionStatus(ctx, j.Job.ClusterId, constants.StatusRollbacking)
	case constants.ActionResizeCluster:
		err = clusterClient.ModifyClusterTransitionStatus(ctx, j.Job.ClusterId, constants.StatusResizing)
	case constants.ActionAddClusterNodes:
		err = clusterClient.ModifyClusterTransitionStatus(ctx, j.Job.ClusterId, constants.StatusScaling)
	case constants.ActionDeleteClusterNodes:
		err = clusterClient.ModifyClusterTransitionStatus(ctx, j.Job.ClusterId, constants.StatusScaling)
	case constants.ActionStopClusters:
		err = clusterClient.ModifyClusterTransitionStatus(ctx, j.Job.ClusterId, constants.StatusStopping)
	case constants.ActionStartClusters:
		err = clusterClient.ModifyClusterTransitionStatus(ctx, j.Job.ClusterId, constants.StatusStarting)
	case constants.ActionDeleteClusters:
		err = clusterClient.ModifyClusterTransitionStatus(ctx, j.Job.ClusterId, constants.StatusDeleting)
	case constants.ActionRecoverClusters:
		err = clusterClient.ModifyClusterTransitionStatus(ctx, j.Job.ClusterId, constants.StatusRecovering)
	case constants.ActionCeaseClusters:
		err = clusterClient.ModifyClusterTransitionStatus(ctx, j.Job.ClusterId, constants.StatusCeasing)
	case constants.ActionUpdateClusterEnv:
		err = clusterClient.ModifyClusterTransitionStatus(ctx, j.Job.ClusterId, constants.StatusUpdating)
	default:
		j.JLogger.Error("Unknown job action [%s]", j.Job.JobAction)
	}

	if err != nil {
		j.JLogger.Critical("Executing job pre processor failed: %+v", err)
	}
	return err
}

// Post process when job is done
func (j *Processor) Post() error {
	var err error
	ctx := client.GetSystemUserContext()
	clusterClient, err := clusterclient.NewClient()
	if err != nil {
		j.JLogger.Error("Executing job post processor failed: %+v", err)
		return err
	}
	switch j.Job.JobAction {
	case constants.ActionCreateCluster:
		providerInterface, err := plugins.GetProviderPlugin(j.Job.Provider, j.JLogger)
		if err != nil {
			j.JLogger.Error("No such provider [%s]. ", j.Job.Provider)
			return err
		}
		err = providerInterface.UpdateClusterStatus(j.Job)
		if err != nil {
			j.JLogger.Error("Executing job post processor failed: %+v", err)
			return err
		}

		err = clusterClient.ModifyClusterStatus(ctx, j.Job.ClusterId, constants.StatusActive)
	case constants.ActionUpgradeCluster:
		providerInterface, err := plugins.GetProviderPlugin(j.Job.Provider, j.JLogger)
		if err != nil {
			j.JLogger.Error("No such provider [%s]. ", j.Job.Provider)
			return err
		}
		err = providerInterface.UpdateClusterStatus(j.Job)
		if err != nil {
			j.JLogger.Error("Executing job post processor failed: %+v", err)
			return err
		}

		err = clusterClient.ModifyClusterStatus(ctx, j.Job.ClusterId, constants.StatusActive)
	case constants.ActionRollbackCluster:
		providerInterface, err := plugins.GetProviderPlugin(j.Job.Provider, j.JLogger)
		if err != nil {
			j.JLogger.Error("No such provider [%s]. ", j.Job.Provider)
			return err
		}
		err = providerInterface.UpdateClusterStatus(j.Job)
		if err != nil {
			j.JLogger.Error("Executing job post processor failed: %+v", err)
			return err
		}

		err = clusterClient.ModifyClusterStatus(ctx, j.Job.ClusterId, constants.StatusActive)
	case constants.ActionResizeCluster:
		err = clusterClient.ModifyClusterStatus(ctx, j.Job.ClusterId, constants.StatusActive)
	case constants.ActionAddClusterNodes:
		// delete node record from db when pre check is failed
		if j.Job.Status == constants.StatusFailed {
			clusterWrappers, err := clusterClient.GetClusterWrappers(ctx, []string{j.Job.ClusterId})
			if err != nil {
				j.JLogger.Error("No such cluster [%s], %+v ", j.Job.ClusterId, err)
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
				clusterClient.DeleteTableClusterNodes(ctx, &pb.DeleteTableClusterNodesRequest{
					NodeId: deleteNodeIds,
				})
			}
		}
		err = clusterClient.ModifyClusterStatus(ctx, j.Job.ClusterId, constants.StatusActive)
	case constants.ActionDeleteClusterNodes:
		err = clusterClient.ModifyClusterStatus(ctx, j.Job.ClusterId, constants.StatusActive)
	case constants.ActionStopClusters:
		err = clusterClient.ModifyClusterStatus(ctx, j.Job.ClusterId, constants.StatusStopped)
		clusterWrappers, err := clusterClient.GetClusterWrappers(ctx, []string{j.Job.ClusterId})
		if err != nil {
			return err
		}
		clusterWrapper := clusterWrappers[0]
		frontgateId := clusterWrapper.Cluster.FrontgateId
		pbClusters, err := clusterClient.DescribeClustersWithFrontgateId(ctx, frontgateId,
			[]string{constants.StatusActive, constants.StatusPending})
		if err != nil {
			return err
		}
		if len(pbClusters) == 0 {
			// need to delete frontgate cluster
			frontgates, err := clusterClient.GetClusterWrappers(ctx, []string{frontgateId})
			if err != nil {
				return err
			}
			frontgate := frontgates[0]
			directive := jsonutil.ToString(frontgate)

			newJob := models.NewJob(
				constants.PlaceHolder,
				frontgate.Cluster.ClusterId,
				frontgate.Cluster.AppId,
				frontgate.Cluster.VersionId,
				constants.ActionStopClusters,
				directive,
				j.Job.Provider,
				j.Job.Owner,
			)

			_, err = jobclient.SendJob(newJob)
			if err != nil {
				return err
			}
		}
	case constants.ActionStartClusters:
		err = clusterClient.ModifyClusterStatus(ctx, j.Job.ClusterId, constants.StatusActive)
	case constants.ActionDeleteClusters:
		err = clusterClient.ModifyClusterStatus(ctx, j.Job.ClusterId, constants.StatusDeleted)
		clusterWrappers, err := clusterClient.GetClusterWrappers(ctx, []string{j.Job.ClusterId})
		if err != nil {
			return err
		}
		clusterWrapper := clusterWrappers[0]
		if clusterWrapper.Cluster.ClusterType == constants.NormalClusterType {
			frontgateId := clusterWrapper.Cluster.FrontgateId
			pbClusters, err := clusterClient.DescribeClustersWithFrontgateId(ctx, frontgateId,
				[]string{constants.StatusStopped, constants.StatusActive, constants.StatusPending})
			if err != nil {
				return err
			}
			if len(pbClusters) == 0 {
				// need to delete frontgate cluster
				frontgates, err := clusterClient.GetClusterWrappers(ctx, []string{frontgateId})
				if err != nil {
					return err
				}
				frontgate := frontgates[0]
				directive := jsonutil.ToString(frontgate)

				newJob := models.NewJob(
					constants.PlaceHolder,
					frontgate.Cluster.ClusterId,
					frontgate.Cluster.AppId,
					frontgate.Cluster.VersionId,
					constants.ActionDeleteClusters,
					directive,
					j.Job.Provider,
					j.Job.Owner,
				)

				_, err = jobclient.SendJob(newJob)
				if err != nil {
					return err
				}
			}
		}
	case constants.ActionRecoverClusters:
		err = clusterClient.ModifyClusterStatus(ctx, j.Job.ClusterId, constants.StatusActive)
	case constants.ActionCeaseClusters:
		err = clusterClient.ModifyClusterStatus(ctx, j.Job.ClusterId, constants.StatusCeased)
	case constants.ActionUpdateClusterEnv:
		err = clusterClient.ModifyClusterStatus(ctx, j.Job.ClusterId, constants.StatusActive)
	default:
		j.JLogger.Error("Unknown job action [%s]", j.Job.JobAction)
	}

	if err != nil {
		j.JLogger.Error("Executing job post processor failed: %+v", err)
	}
	return err
}

func (j *Processor) Final() {
	ctx := client.GetSystemUserContext()
	clusterClient, err := clusterclient.NewClient()
	if err != nil {
		j.JLogger.Error("Executing job final processor failed: %+v", err)
		return
	}
	// TODO: modify cluster status to `active or deleted`
	err = clusterClient.ModifyClusterTransitionStatus(ctx, j.Job.ClusterId, "")
	if err != nil {
		j.JLogger.Error("Executing job final processor failed: %+v", err)
	}
}
