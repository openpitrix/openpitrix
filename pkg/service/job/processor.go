// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package job

import (
	"context"

	clusterclient "openpitrix.io/openpitrix/pkg/client/cluster"
	jobclient "openpitrix.io/openpitrix/pkg/client/job"
	providerclient "openpitrix.io/openpitrix/pkg/client/runtime_provider"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/gerr"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/pi"
	"openpitrix.io/openpitrix/pkg/plugins"
	"openpitrix.io/openpitrix/pkg/util/jsonutil"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
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
func (p *Processor) Pre(ctx context.Context) error {
	var err error
	clusterClient, err := clusterclient.NewClient()
	if err != nil {
		logger.Error(ctx, "Executing job pre processor failed: %+v", err)
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
		logger.Error(ctx, "Unknown job action [%s]", p.Job.JobAction)
	}

	if err != nil {
		logger.Critical(ctx, "Executing job pre processor failed: %+v", err)
	}
	return err
}

// Post process when job is done
func (p *Processor) Post(ctx context.Context) error {
	var err error
	clusterClient, err := clusterclient.NewClient()
	if err != nil {
		logger.Error(ctx, "Executing job post processor failed: %+v", err)
		return err
	}
	switch p.Job.JobAction {
	case constants.ActionCreateCluster:
		err = p.UpdateClusterDetails(ctx)
		if err != nil {
			logger.Error(ctx, "Update cluster details failed: %+v", err)
			return err
		}

		err = clusterClient.ModifyClusterStatus(ctx, p.Job.ClusterId, constants.StatusActive)
	case constants.ActionUpgradeCluster:
		err = p.UpdateClusterDetails(ctx)
		if err != nil {
			logger.Error(ctx, "Update cluster details failed: %+v", err)
			return err
		}

		err = clusterClient.ModifyClusterStatus(ctx, p.Job.ClusterId, constants.StatusActive)
	case constants.ActionRollbackCluster:
		err = p.UpdateClusterDetails(ctx)
		if err != nil {
			logger.Error(ctx, "Update cluster details failed: %+v", err)
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
				logger.Error(ctx, "No such cluster [%s], %+v ", p.Job.ClusterId, err)
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
		if !plugins.IsVmbasedProviders(p.Job.Provider) {
			return nil
		}
		clusterWrapper := clusterWrappers[0]
		if clusterWrapper.Cluster.ClusterType == constants.NormalClusterType {
			frontgateId := clusterWrapper.Cluster.FrontgateId
			pbClusters, err := clusterClient.DescribeClustersWithFrontgateId(
				ctx,
				frontgateId,
				[]string{constants.StatusActive, constants.StatusPending},
				clusterWrapper.Cluster.Debug,
			)
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
					p.Job.OwnerPath,
					frontgate.Cluster.RuntimeId,
				)

				_, err = jobclient.SendJob(ctx, newJob)
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
		if !plugins.IsVmbasedProviders(p.Job.Provider) {
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
			pbClusters, err := clusterClient.DescribeClustersWithFrontgateId(
				ctx,
				frontgateId,
				[]string{constants.StatusStopped, constants.StatusActive, constants.StatusPending},
				clusterWrapper.Cluster.Debug,
			)
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
					p.Job.OwnerPath,
					frontgate.Cluster.RuntimeId,
				)

				_, err = jobclient.SendJob(ctx, newJob)
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
		err = p.UpdateClusterDetails(ctx)
		if err != nil {
			logger.Error(ctx, "Update cluster details failed: %+v", err)
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
		logger.Error(ctx, "Unknown job action [%s]", p.Job.JobAction)
	}

	if err != nil {
		logger.Error(ctx, "Executing job post processor failed: %+v", err)
	}
	return err
}

func (p *Processor) Final(ctx context.Context) {
	clusterClient, err := clusterclient.NewClient()
	if err != nil {
		logger.Error(ctx, "Executing job final processor failed: %+v", err)
		return
	}
	err = clusterClient.ModifyClusterTransitionStatus(ctx, p.Job.ClusterId, "")
	if err != nil {
		logger.Error(ctx, "Executing job final processor failed: %+v", err)
	}
}

func (p *Processor) UpdateClusterDetails(ctx context.Context) error {
	clusterWrapper, err := models.NewClusterWrapper(ctx, p.Job.Directive)
	if err != nil {
		return err
	}

	providerClient, err := providerclient.NewRuntimeProviderManagerClient()
	if err != nil {
		return gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}

	response, err := providerClient.DescribeClusterDetails(ctx, &pb.DescribeClusterDetailsRequest{
		RuntimeId: pbutil.ToProtoString(clusterWrapper.Cluster.RuntimeId),
		Cluster:   models.ClusterWrapperToPb(clusterWrapper),
	})
	if err != nil {
		logger.Error(ctx, "Describe cluster details failed, %+v", err)
		return err
	}
	clusterWrapper = models.PbToClusterWrapper(response.Cluster)

	if clusterWrapper != nil && len(clusterWrapper.Cluster.ClusterId) > 0 {
		var clusterRoles []*models.ClusterRole
		for _, clusterRole := range clusterWrapper.ClusterRoles {
			clusterRoles = append(clusterRoles, clusterRole)
		}

		var clusterNodes []*models.ClusterNodeWithKeyPairs
		for _, clusterNode := range clusterWrapper.ClusterNodesWithKeyPairs {
			clusterNodes = append(clusterNodes, clusterNode)
		}

		clusterClient, err := clusterclient.NewClient()
		if err != nil {
			return err
		}

		modifyClusterRequest := &pb.ModifyClusterRequest{
			Cluster: &pb.Cluster{
				ClusterId:   pbutil.ToProtoString(clusterWrapper.Cluster.ClusterId),
				Description: pbutil.ToProtoString(clusterWrapper.Cluster.Description),
			},
			ClusterRoleSet: models.ClusterRolesToPbs(clusterRoles),
			ClusterNodeSet: models.ClusterNodesWithKeyPairsToPbs(clusterNodes),
		}
		_, err = clusterClient.ModifyCluster(ctx, modifyClusterRequest)
		if err != nil {
			return err
		}
	}

	return nil
}
