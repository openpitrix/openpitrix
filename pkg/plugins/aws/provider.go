// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package aws

import (
	"context"
	"fmt"
	"time"

	clientutil "openpitrix.io/openpitrix/pkg/client"
	clusterclient "openpitrix.io/openpitrix/pkg/client/cluster"
	runtimeclient "openpitrix.io/openpitrix/pkg/client/runtime"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/pi"
	"openpitrix.io/openpitrix/pkg/plugins/vmbased"
	"openpitrix.io/openpitrix/pkg/util/stringutil"
)

type Provider struct {
	ctx context.Context
}

func NewProvider(ctx context.Context) *Provider {
	return &Provider{
		ctx,
	}
}

func (p *Provider) ParseClusterConf(versionId, runtimeId, conf string, clusterWrapper *models.ClusterWrapper) error {
	frameInterface, err := vmbased.NewFrameInterface(p.ctx, nil)
	if err != nil {
		return err
	}
	err = frameInterface.ParseClusterConf(versionId, runtimeId, conf, clusterWrapper)
	if err != nil {
		return err
	}
	handler := GetProviderHandler(p.ctx)
	availabilityZone, err := handler.DescribeAvailabilityZoneBySubnetId(runtimeId, clusterWrapper.Cluster.SubnetId)
	if err != nil {
		return err
	}
	clusterWrapper.Cluster.Zone = availabilityZone
	return nil
}

func (p *Provider) SplitJobIntoTasks(job *models.Job) (*models.TaskLayer, error) {
	var clusterWrapper *models.ClusterWrapper
	var err error

	switch job.JobAction {
	case constants.ActionAttachKeyPairs, constants.ActionDetachKeyPairs:
		nodeKeyPairDetails, err := models.NewNodeKeyPairDetails(job.Directive)
		if err != nil {
			return nil, err
		}
		clusterId := nodeKeyPairDetails[0].ClusterNode.ClusterId
		clusterClient, err := clusterclient.NewClient()
		if err != nil {
			return nil, err
		}
		ctx := clientutil.SetSystemUserToContext(p.ctx)
		pbClusterWrappers, err := clusterClient.GetClusterWrappers(ctx, []string{clusterId})
		if err != nil {
			return nil, err
		}
		clusterWrapper = pbClusterWrappers[0]
	default:
		clusterWrapper, err = models.NewClusterWrapper(p.ctx, job.Directive)
		if err != nil {
			return nil, err
		}
	}

	runtimeId := clusterWrapper.Cluster.RuntimeId
	runtime, err := runtimeclient.NewRuntime(p.ctx, runtimeId)
	if err != nil {
		return nil, err
	}
	imageConfig, err := pi.Global().GlobalConfig().GetRuntimeImageIdAndUrl(runtime.RuntimeUrl, runtime.Zone)
	if err != nil {
		return nil, err
	}
	if imageConfig.ImageId == "" && imageConfig.ImageName != "" {
		handler := GetProviderHandler(p.ctx)
		imageConfig.ImageId, err = handler.DescribeImage(runtimeId, imageConfig.ImageName)
		if err != nil {
			return nil, err
		}
	}
	frameInterface, err := vmbased.NewFrameInterface(p.ctx, job, imageConfig.ImageId)
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
	case constants.ActionResizeCluster:

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
	case constants.ActionRecoverClusters:
		// not supported yet
		return nil, nil
	case constants.ActionCeaseClusters:
		// not supported yet
		return nil, nil
	case constants.ActionUpdateClusterEnv:
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
		logger.Error(p.ctx, "Unknown job action [%s]", job.JobAction)
		return nil, fmt.Errorf("unknown job action [%s]", job.JobAction)
	}
	return nil, nil
}

func (p *Provider) HandleSubtask(task *models.Task) error {
	handler := GetProviderHandler(p.ctx)

	switch task.TaskAction {
	case vmbased.ActionRunInstances:
		return handler.RunInstances(task)
	case vmbased.ActionStopInstances:
		return handler.StopInstances(task)
	case vmbased.ActionStartInstances:
		return handler.StartInstances(task)
	case vmbased.ActionTerminateInstances:
		return handler.DeleteInstances(task)
	case vmbased.ActionCreateVolumes:
		return handler.CreateVolumes(task)
	case vmbased.ActionDetachVolumes:
		return handler.DetachVolumes(task)
	case vmbased.ActionAttachVolumes:
		return handler.AttachVolumes(task)
	case vmbased.ActionDeleteVolumes:
		return handler.DeleteVolumes(task)

	case vmbased.ActionWaitFrontgateAvailable:
		// do nothing
		return nil
	default:
		logger.Error(p.ctx, "Unknown task action [%s]", task.TaskAction)
		return fmt.Errorf("unknown task action [%s]", task.TaskAction)
	}
}
func (p *Provider) WaitSubtask(task *models.Task, timeout time.Duration, waitInterval time.Duration) error {
	logger.Debug(p.ctx, "Wait sub task timeout [%s] interval [%s]", timeout, waitInterval)
	handler := GetProviderHandler(p.ctx)

	switch task.TaskAction {
	case vmbased.ActionRunInstances:
		return handler.WaitRunInstances(task)
	case vmbased.ActionStopInstances:
		return handler.WaitStopInstances(task)
	case vmbased.ActionStartInstances:
		return handler.WaitStartInstances(task)
	case vmbased.ActionTerminateInstances:
		return handler.WaitDeleteInstances(task)
	case vmbased.ActionCreateVolumes:
		return handler.WaitCreateVolumes(task)
	case vmbased.ActionDetachVolumes:
		return handler.WaitDetachVolumes(task)
	case vmbased.ActionAttachVolumes:
		return handler.WaitAttachVolumes(task)
	case vmbased.ActionDeleteVolumes:
		return handler.WaitDeleteVolumes(task)
	case vmbased.ActionWaitFrontgateAvailable:
		return handler.WaitFrontgateAvailable(task)
	default:
		logger.Error(p.ctx, "Unknown task action [%s]", task.TaskAction)
		return fmt.Errorf("unknown task action [%s]", task.TaskAction)
	}
}

func (p *Provider) DescribeSubnets(ctx context.Context, req *pb.DescribeSubnetsRequest) (*pb.DescribeSubnetsResponse, error) {
	handler := GetProviderHandler(p.ctx)
	return handler.DescribeSubnets(ctx, req)
}

func (p *Provider) CheckResource(ctx context.Context, clusterWrapper *models.ClusterWrapper) error {
	handler := GetProviderHandler(p.ctx)
	return handler.CheckResourceQuotas(ctx, clusterWrapper)
}

func (p *Provider) DescribeVpc(runtimeId, vpcId string) (*models.Vpc, error) {
	handler := GetProviderHandler(p.ctx)
	return handler.DescribeVpc(runtimeId, vpcId)
}

func (p *Provider) ValidateCredential(url, credential, zone string) error {
	handler := GetProviderHandler(p.ctx)
	zones, err := handler.DescribeZones(url, credential)
	if err != nil {
		return err
	}
	if zone == "" {
		return nil
	}
	if !stringutil.StringIn(zone, zones) {
		return fmt.Errorf("cannot access zone [%s]", zone)
	}
	return nil
}

func (p *Provider) UpdateClusterStatus(job *models.Job) error {
	return nil
}

func (p *Provider) DescribeRuntimeProviderAvailabilityZones(url, credential, zone string) ([]string, error) {
	handler := GetProviderHandler(p.ctx)
	return handler.DescribeAvailabilityZones(url, credential, zone)
}

func (p *Provider) DescribeRuntimeProviderZones(url, credential string) ([]string, error) {
	handler := GetProviderHandler(p.ctx)
	return handler.DescribeZones(url, credential)
}

func (p *Provider) DescribeClusterDetails(ctx context.Context, cluster *models.ClusterWrapper) error {
	return nil
}
