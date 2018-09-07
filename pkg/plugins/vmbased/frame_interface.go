// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package vmbased

import (
	"context"
	"fmt"

	clientutil "openpitrix.io/openpitrix/pkg/client"
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
	ResizeClusterLayer(roleResizeResources models.RoleResizeResources) *models.TaskLayer
	AttachKeyPairsLayer(nodeKeyPairDetails models.NodeKeyPairDetails) *models.TaskLayer
	DetachKeyPairsLayer(nodeKeyPairDetails models.NodeKeyPairDetails) *models.TaskLayer
	ParseClusterConf(versionId, runtimeId, conf string, clusterWrapper *models.ClusterWrapper) error
}

func NewFrameInterface(ctx context.Context, job *models.Job, advancedParam ...string) (FrameInterface, error) {
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
		ctx := clientutil.SetSystemUserToContext(ctx)
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
		ctx := clientutil.SetSystemUserToContext(ctx)
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
