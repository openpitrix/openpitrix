// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package job

import (
	"context"

	"openpitrix.io/openpitrix/pkg/client"
	clusterclient "openpitrix.io/openpitrix/pkg/client/cluster"
	jobclient "openpitrix.io/openpitrix/pkg/client/job"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/pi"
	"openpitrix.io/openpitrix/pkg/plugins"
	"openpitrix.io/openpitrix/pkg/util/jsonutil"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
	"openpitrix.io/openpitrix/pkg/util/reflectutil"
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
func (p *Processor) Pre() error {
	ctx := client.GetSystemUserContext()
	var err error
	clusterClient, err := clusterclient.NewClient()
	if err != nil {
		p.JLogger.Error("Executing job pre processor failed: %+v", err)
		return err
	}
	switch p.Job.JobAction {
	case constants.ActionCreateCluster:
		err = clusterClient.ModifyClusterTransitionStatus(ctx, p.Job.ClusterId, constants.StatusCreating)
	case constants.ActionUpgradeCluster:
		err = clusterClient.ModifyClusterTransitionStatus(ctx, p.Job.ClusterId, constants.StatusUpgrading)
	case constants.ActionRollbackCluster:
		err = clusterClient.ModifyClusterTransitionStatus(ctx, p.Job.ClusterId, constants.StatusRollbacking)
	case constants.ActionResizeCluster:
		err = clusterClient.ModifyClusterTransitionStatus(ctx, p.Job.ClusterId, constants.StatusResizing)
	case constants.ActionAddClusterNodes:
		err = clusterClient.ModifyClusterTransitionStatus(ctx, p.Job.ClusterId, constants.StatusScaling)
	case constants.ActionDeleteClusterNodes:
		err = clusterClient.ModifyClusterTransitionStatus(ctx, p.Job.ClusterId, constants.StatusScaling)
	case constants.ActionStopClusters:
		err = clusterClient.ModifyClusterTransitionStatus(ctx, p.Job.ClusterId, constants.StatusStopping)
	case constants.ActionStartClusters:
		err = clusterClient.ModifyClusterTransitionStatus(ctx, p.Job.ClusterId, constants.StatusStarting)
	case constants.ActionDeleteClusters:
		err = clusterClient.ModifyClusterTransitionStatus(ctx, p.Job.ClusterId, constants.StatusDeleting)
	case constants.ActionRecoverClusters:
		err = clusterClient.ModifyClusterTransitionStatus(ctx, p.Job.ClusterId, constants.StatusRecovering)
	case constants.ActionCeaseClusters:
		err = clusterClient.ModifyClusterTransitionStatus(ctx, p.Job.ClusterId, constants.StatusCeasing)
	case constants.ActionUpdateClusterEnv, constants.ActionAttachKeyPairs, constants.ActionDetachKeyPairs:
		err = clusterClient.ModifyClusterTransitionStatus(ctx, p.Job.ClusterId, constants.StatusUpdating)
	default:
		p.JLogger.Error("Unknown job action [%s]", p.Job.JobAction)
	}

	if err != nil {
		p.JLogger.Critical("Executing job pre processor failed: %+v", err)
	}
	return err
}

// Post process when job is done
func (p *Processor) Post() error {
	ctx := client.GetSystemUserContext()
	var err error
	clusterClient, err := clusterclient.NewClient()
	if err != nil {
		p.JLogger.Error("Executing job post processor failed: %+v", err)
		return err
	}
	switch p.Job.JobAction {
	case constants.ActionCreateCluster:
		providerInterface, err := plugins.GetProviderPlugin(p.Job.Provider, p.JLogger)
		if err != nil {
			p.JLogger.Error("No such provider [%s]. ", p.Job.Provider)
			return err
		}
		err = providerInterface.UpdateClusterStatus(p.Job)
		if err != nil {
			p.JLogger.Error("Executing job post processor failed: %+v", err)
			return err
		}

		err = clusterClient.ModifyClusterStatus(ctx, p.Job.ClusterId, constants.StatusActive)
	case constants.ActionUpgradeCluster:
		providerInterface, err := plugins.GetProviderPlugin(p.Job.Provider, p.JLogger)
		if err != nil {
			p.JLogger.Error("No such provider [%s]. ", p.Job.Provider)
			return err
		}
		err = providerInterface.UpdateClusterStatus(p.Job)
		if err != nil {
			p.JLogger.Error("Executing job post processor failed: %+v", err)
			return err
		}

		err = clusterClient.ModifyClusterStatus(ctx, p.Job.ClusterId, constants.StatusActive)
	case constants.ActionRollbackCluster:
		providerInterface, err := plugins.GetProviderPlugin(p.Job.Provider, p.JLogger)
		if err != nil {
			p.JLogger.Error("No such provider [%s]. ", p.Job.Provider)
			return err
		}
		err = providerInterface.UpdateClusterStatus(p.Job)
		if err != nil {
			p.JLogger.Error("Executing job post processor failed: %+v", err)
			return err
		}

		err = clusterClient.ModifyClusterStatus(ctx, p.Job.ClusterId, constants.StatusActive)
	case constants.ActionResizeCluster:
		err = clusterClient.ModifyClusterStatus(ctx, p.Job.ClusterId, constants.StatusActive)
	case constants.ActionAddClusterNodes:
		// delete node record from db when pre check is failed
		if p.Job.Status == constants.StatusFailed {
			clusterWrappers, err := clusterClient.GetClusterWrappers(ctx, []string{p.Job.ClusterId})
			if err != nil {
				p.JLogger.Error("No such cluster [%s], %+v ", p.Job.ClusterId, err)
				return err
			}
			clusterWrapper := clusterWrappers[0]
			var deleteNodeIds []string
			for _, clusterNode := range clusterWrapper.ClusterNodesWithKeyPairs {
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
		err = clusterClient.ModifyClusterStatus(ctx, p.Job.ClusterId, constants.StatusActive)
	case constants.ActionDeleteClusterNodes:
		err = clusterClient.ModifyClusterStatus(ctx, p.Job.ClusterId, constants.StatusActive)
	case constants.ActionStopClusters:
		err = clusterClient.ModifyClusterStatus(ctx, p.Job.ClusterId, constants.StatusStopped)
		clusterWrappers, err := clusterClient.GetClusterWrappers(ctx, []string{p.Job.ClusterId})
		if err != nil {
			return err
		}
		if !reflectutil.In(p.Job.Provider, constants.VmBaseProviders) {
			return nil
		}
		clusterWrapper := clusterWrappers[0]
		if clusterWrapper.Cluster.ClusterType == constants.NormalClusterType {
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
					p.Job.Provider,
					p.Job.Owner,
				)

				_, err = jobclient.SendJob(newJob)
				if err != nil {
					return err
				}
			}
		}
	case constants.ActionStartClusters:
		err = clusterClient.ModifyClusterStatus(ctx, p.Job.ClusterId, constants.StatusActive)
	case constants.ActionDeleteClusters:
		err = clusterClient.ModifyClusterStatus(ctx, p.Job.ClusterId, constants.StatusDeleted)
		if err != nil {
			return err
		}
		if !reflectutil.In(p.Job.Provider, constants.VmBaseProviders) {
			return nil
		}
		clusterWrappers, err := clusterClient.GetClusterWrappers(ctx, []string{p.Job.ClusterId})
		if err != nil {
			return err
		}
		clusterWrapper := clusterWrappers[0]

		// delete node key pairs
		var pbNodeKeyPairs []*pb.NodeKeyPair
		for _, clusterNode := range clusterWrapper.ClusterNodesWithKeyPairs {
			for _, keyPairId := range clusterNode.KeyPairId {
				pbNodeKeyPairs = append(pbNodeKeyPairs, &pb.NodeKeyPair{
					NodeId:    pbutil.ToProtoString(clusterNode.NodeId),
					KeyPairId: pbutil.ToProtoString(keyPairId),
				})
			}
		}
		_, err = clusterClient.DeleteNodeKeyPairs(ctx, &pb.DeleteNodeKeyPairsRequest{
			NodeKeyPair: pbNodeKeyPairs,
		})
		if err != nil {
			return err
		}

		if clusterWrapper.Cluster.ClusterType == constants.NormalClusterType && pi.Global().GlobalConfig().Cluster.FrontgateAutoDelete {
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
					p.Job.Provider,
					p.Job.Owner,
				)

				_, err = jobclient.SendJob(newJob)
				if err != nil {
					return err
				}
			}
		}
	case constants.ActionRecoverClusters:
		err = clusterClient.ModifyClusterStatus(ctx, p.Job.ClusterId, constants.StatusActive)
	case constants.ActionCeaseClusters:
		err = clusterClient.ModifyClusterStatus(ctx, p.Job.ClusterId, constants.StatusCeased)
	case constants.ActionUpdateClusterEnv:
		providerInterface, err := plugins.GetProviderPlugin(p.Job.Provider, p.JLogger)
		if err != nil {
			p.JLogger.Error("No such provider [%s]. ", p.Job.Provider)
			return err
		}
		err = providerInterface.UpdateClusterStatus(p.Job)
		if err != nil {
			p.JLogger.Error("Executing job post processor failed: %+v", err)
			return err
		}

		err = clusterClient.ModifyClusterStatus(ctx, p.Job.ClusterId, constants.StatusActive)
	case constants.ActionAttachKeyPairs:
		nodeKeyPairDetails, err := models.NewNodeKeyPairDetails(p.Job.Directive)
		if err != nil {
			return err
		}
		var pbNodeKeyPairs []*pb.NodeKeyPair
		for _, nodeKeyPairDetail := range nodeKeyPairDetails {
			pbNodeKeyPair := models.NodeKeyPairToPb(nodeKeyPairDetail.NodeKeyPair)
			pbNodeKeyPairs = append(pbNodeKeyPairs, pbNodeKeyPair)
		}
		_, err = clusterClient.AddNodeKeyPairs(ctx, &pb.AddNodeKeyPairsRequest{
			NodeKeyPair: pbNodeKeyPairs,
		})
	case constants.ActionDetachKeyPairs:
		nodeKeyPairDetails, err := models.NewNodeKeyPairDetails(p.Job.Directive)
		if err != nil {
			return err
		}
		var pbNodeKeyPairs []*pb.NodeKeyPair
		for _, nodeKeyPairDetail := range nodeKeyPairDetails {
			pbNodeKeyPair := models.NodeKeyPairToPb(nodeKeyPairDetail.NodeKeyPair)
			pbNodeKeyPairs = append(pbNodeKeyPairs, pbNodeKeyPair)
		}
		_, err = clusterClient.DeleteNodeKeyPairs(ctx, &pb.DeleteNodeKeyPairsRequest{
			NodeKeyPair: pbNodeKeyPairs,
		})
	default:
		p.JLogger.Error("Unknown job action [%s]", p.Job.JobAction)
	}

	if err != nil {
		p.JLogger.Error("Executing job post processor failed: %+v", err)
	}
	return err
}

func (p *Processor) Final() {
	ctx := context.WithValue(client.GetSystemUserContext(), "owner", p.Job.Owner)
	clusterClient, err := clusterclient.NewClient()
	if err != nil {
		p.JLogger.Error("Executing job final processor failed: %+v", err)
		return
	}
	err = clusterClient.ModifyClusterTransitionStatus(ctx, p.Job.ClusterId, "")
	if err != nil {
		p.JLogger.Error("Executing job final processor failed: %+v", err)
	}
}
