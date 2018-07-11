// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package vmbased

import (
	"fmt"

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
	ParseClusterConf(versionId, runtimeId, conf string) (*models.ClusterWrapper, error)
}

func NewFrameInterface(job *models.Job, logger *logger.Logger, advancedParam ...string) (FrameInterface, error) {
	if job == nil {
		return &Frame{Logger: logger}, nil
	}
	clusterWrapper, err := models.NewClusterWrapper(job.Directive)
	if err != nil {
		return nil, err
	}

	runtimeId := clusterWrapper.Cluster.RuntimeId
	runtime, err := runtimeclient.NewRuntime(runtimeId)
	if err != nil {
		return nil, err
	}
	imageConfig, err := pi.Global().GlobalConfig().GetRuntimeImageIdAndUrl(runtime.RuntimeUrl, runtime.Zone)
	if err != nil {
		return nil, err
	}

	frame := &Frame{
		Job:            job,
		ClusterWrapper: clusterWrapper,
		Runtime:        runtime,
		Logger:         logger,
		ImageConfig:    imageConfig,
	}

	if len(advancedParam) >= 1 {
		frame.ImageConfig.ImageId = advancedParam[0]
	}

	if frame.ImageConfig.ImageId == "" {
		logger.Error("Failed to find image id for url [%s], zone [%s]", runtime.RuntimeUrl, runtime.Zone)
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
