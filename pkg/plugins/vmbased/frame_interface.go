// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package vmbased

import (
	"context"
	"fmt"

	clusterclient "openpitrix.io/openpitrix/pkg/client/cluster"
	runtimeclient "openpitrix.io/openpitrix/pkg/client/runtime"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pi"
)

type FrameInterface interface {
	CreateClusterLayer() *models.TaskLayer
	StopClusterLayer() *models.TaskLayer
	StartClusterLayer() *models.TaskLayer
	DeleteClusterLayer() *models.TaskLayer
	AddClusterNodesLayer() *models.TaskLayer
	DeleteClusterNodesLayer() *models.TaskLayer
	UpdateClusterEnvLayer() *models.TaskLayer
	ResizeClusterLayer(roleResizeResources models.RoleResizeResources) *models.TaskLayer
	AttachKeyPairsLayer(nodeKeyPairDetails models.NodeKeyPairDetails) *models.TaskLayer
	DetachKeyPairsLayer(nodeKeyPairDetails models.NodeKeyPairDetails) *models.TaskLayer
	ParseClusterConf(ctx context.Context, versionId, runtimeId, conf string, clusterWrapper *models.ClusterWrapper) error
}

func SplitJobIntoTasks(ctx context.Context, job *models.Job, advancedParam ...string) (*models.TaskLayer, error) {
	frameInterface, err := GetFrameInterface(ctx, job, advancedParam...)
	if err != nil {
		return nil, err
	}

	switch job.JobAction {
	case constants.ActionCreateCluster:
		// TODO: vpc, eip, subnet

		return frameInterface.CreateClusterLayer(), nil
	case constants.ActionUpgradeCluster:
		// not supported yet
		return nil, nil
	case constants.ActionRollbackCluster:
		// not supported yet
		return nil, nil
	case constants.ActionAddClusterNodes:
		return frameInterface.AddClusterNodesLayer(), nil
	case constants.ActionDeleteClusterNodes:
		return frameInterface.DeleteClusterNodesLayer(), nil
	case constants.ActionStopClusters:
		return frameInterface.StopClusterLayer(), nil
	case constants.ActionStartClusters:
		return frameInterface.StartClusterLayer(), nil
	case constants.ActionDeleteClusters:
		return frameInterface.DeleteClusterLayer(), nil
	case constants.ActionResizeCluster:
		roleResizeResources, err := models.NewRoleResizeResources(job.Directive)
		if err != nil {
			return nil, err
		}
		return frameInterface.ResizeClusterLayer(roleResizeResources), nil
	case constants.ActionRecoverClusters:
		// not supported yet
		return nil, nil
	case constants.ActionCeaseClusters:
		// not supported yet
		return nil, nil
	case constants.ActionUpdateClusterEnv:
		return frameInterface.UpdateClusterEnvLayer(), nil
	case constants.ActionAttachKeyPairs:
		nodeKeyPairDetails, err := models.NewNodeKeyPairDetails(job.Directive)
		if err != nil {
			return nil, err
		}
		return frameInterface.AttachKeyPairsLayer(nodeKeyPairDetails), nil
	case constants.ActionDetachKeyPairs:
		nodeKeyPairDetails, err := models.NewNodeKeyPairDetails(job.Directive)
		if err != nil {
			return nil, err
		}
		return frameInterface.DetachKeyPairsLayer(nodeKeyPairDetails), nil
	default:
		logger.Error(ctx, "Unknown job action [%s]", job.JobAction)
		return nil, fmt.Errorf("unknown job action [%s]", job.JobAction)
	}
	return nil, nil
}

func GetFrameInterface(ctx context.Context, job *models.Job, advancedParam ...string) (FrameInterface, error) {
	if job == nil {
		return &Frame{Ctx: ctx}, nil
	}

	var clusterWrapper *models.ClusterWrapper
	var err error

	switch job.JobAction {
	case constants.ActionAttachKeyPairs, constants.ActionDetachKeyPairs, constants.ActionResizeCluster:
		clusterId := job.ClusterId
		clusterClient, err := clusterclient.NewClient()
		if err != nil {
			return nil, err
		}
		pbClusterWrappers, err := clusterClient.GetClusterWrappers(ctx, []string{clusterId})
		if err != nil {
			return nil, err
		}
		clusterWrapper = pbClusterWrappers[0]
	default:
		clusterWrapper, err = models.NewClusterWrapper(ctx, job.Directive)
		if err != nil {
			return nil, err
		}
	}

	runtimeId := clusterWrapper.Cluster.RuntimeId
	runtime, err := runtimeclient.NewRuntime(ctx, runtimeId)
	if err != nil {
		return nil, err
	}
	imageConfig, err := pi.Global().GlobalConfig().GetRuntimeImageIdAndUrl(runtime.RuntimeUrl, runtime.Zone)
	if err != nil {
		return nil, err
	}

	var frontgateClusterWrapper *models.ClusterWrapper
	if clusterWrapper.Cluster.ClusterType == constants.NormalClusterType {
		clusterClient, err := clusterclient.NewClient()
		if err != nil {
			return nil, err
		}
		pbClusterWrappers, err := clusterClient.GetClusterWrappers(ctx, []string{clusterWrapper.Cluster.FrontgateId})
		if err != nil {
			return nil, err
		}
		frontgateClusterWrapper = pbClusterWrappers[0]
	}

	frame := &Frame{
		Job:                     job,
		ClusterWrapper:          clusterWrapper,
		FrontgateClusterWrapper: frontgateClusterWrapper,
		Runtime:                 runtime,
		Ctx:                     ctx,
		ImageConfig:             imageConfig,
	}

	if len(advancedParam) >= 1 {
		frame.ImageConfig.ImageId = advancedParam[0]
	}

	if frame.ImageConfig.ImageId == "" {
		logger.Error(ctx, "Failed to find image id for url [%s], zone [%s]", runtime.RuntimeUrl, runtime.Zone)
		return nil, fmt.Errorf("failed to find image id for url [%s], zone [%s]", runtime.RuntimeUrl, runtime.Zone)
	}

	switch clusterWrapper.Cluster.ClusterType {
	case constants.NormalClusterType:
		return frame, nil
	case constants.FrontgateClusterType:
		return &Frontgate{Frame: frame}, nil
	default:
		return frame, nil
	}
}
