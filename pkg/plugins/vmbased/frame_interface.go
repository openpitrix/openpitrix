// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package vmbased

import (
	"bytes"
	"context"
	"fmt"

	appclient "openpitrix.io/openpitrix/pkg/client/app"
	clusterclient "openpitrix.io/openpitrix/pkg/client/cluster"
	runtimeclient "openpitrix.io/openpitrix/pkg/client/runtime"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/devkit"
	"openpitrix.io/openpitrix/pkg/devkit/opapp"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/pi"
	"openpitrix.io/openpitrix/pkg/util/jsonutil"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
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

	frame := &Frame{
		Job:                   job,
		ClusterWrapper:        clusterWrapper,
		Runtime:               runtime,
		Ctx:                   ctx,
		RuntimeProviderConfig: imageConfig,
	}

	if len(advancedParam) >= 1 {
		frame.RuntimeProviderConfig.ImageId = advancedParam[0]
	}

	if frame.RuntimeProviderConfig.ImageId == "" {
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

func ParseClusterConf(ctx context.Context, versionId, runtimeId, conf string, clusterWrapper *models.ClusterWrapper) (*models.ClusterWrapper, error) {
	clusterConf := opapp.ClusterConf{}
	clusterEnv := ""
	// Normal cluster need package to generate final conf
	if versionId != constants.FrontgateVersionId {
		appManagerClient, err := appclient.NewAppManagerClient()
		if err != nil {
			logger.Error(ctx, "Connect to app manager failed: %+v", err)
			return clusterWrapper, err
		}

		req := &pb.GetAppVersionPackageRequest{
			VersionId: pbutil.ToProtoString(versionId),
		}

		resp, err := appManagerClient.GetAppVersionPackage(ctx, req)
		if err != nil {
			logger.Error(ctx, "Get app version [%s] package failed: %+v", versionId, err)
			return clusterWrapper, err
		}

		appPackage, err := devkit.LoadArchive(bytes.NewReader(resp.GetPackage()))
		if err != nil {
			logger.Error(ctx, "Load app version [%s] package failed: %+v", versionId, err)
			return clusterWrapper, err
		}
		var confJson jsonutil.Json
		if len(conf) != 0 {
			confJson, err = jsonutil.NewJson([]byte(conf))
			if err != nil {
				logger.Error(ctx, "Parse conf [%s] failed: %+v", conf, err)
				return clusterWrapper, err
			}
			appPackage.ConfigTemplate.FillInDefaultConfig(confJson)
		}
		confJson = appPackage.ConfigTemplate.GetDefaultConfig()
		err = appPackage.Validate(confJson)
		if err != nil {
			logger.Error(ctx, "Validate conf [%s] failed: %+v", conf, err)
			return clusterWrapper, err
		}
		clusterConf, err = appPackage.ClusterConfTemplate.Render(confJson)
		if err != nil {
			logger.Error(ctx, "Render app version [%s] cluster template failed: %+v", versionId, err)
			return clusterWrapper, err
		}
		err = clusterConf.Validate()
		if err != nil {
			logger.Error(ctx, "Validate app version [%s] conf [%s] failed: %+v", versionId, conf, err)
			return clusterWrapper, err
		}
		clusterConf.AppId = resp.GetAppId().GetValue()
		clusterConf.VersionId = resp.GetVersionId().GetValue()

		// Set cluster env
		appPackage.ConfigTemplate.SpecificConfig("env")
		clusterEnv = jsonutil.ToString(appPackage.ConfigTemplate.Config)
	} else {
		err := jsonutil.Decode([]byte(conf), &clusterConf)
		if err != nil {
			logger.Error(ctx, "Parse conf [%s] to cluster failed: %+v", conf, err)
			return clusterWrapper, err
		}
	}

	parser := Parser{Ctx: ctx}
	err := parser.Parse(clusterConf, clusterWrapper, clusterEnv)
	if err != nil {
		logger.Error(ctx, "Parse app version [%s] failed: %+v", versionId, err)
		return clusterWrapper, err
	}

	return clusterWrapper, nil
}
