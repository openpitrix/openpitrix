// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package qingcloud

import (
	"bytes"
	"context"
	"fmt"
	"time"

	appclient "openpitrix.io/openpitrix/pkg/client/app"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/devkit"
	"openpitrix.io/openpitrix/pkg/devkit/app"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/plugins/vmbased"
	"openpitrix.io/openpitrix/pkg/util/jsonutil"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
	"openpitrix.io/openpitrix/pkg/util/stringutil"
)

type Provider struct {
	Logger *logger.Logger
}

func NewProvider(l *logger.Logger) *Provider {
	return &Provider{
		Logger: l,
	}
}

func (p *Provider) ParseClusterConf(versionId, runtimeId, conf string) (*models.ClusterWrapper, error) {
	clusterConf := app.ClusterConf{}
	// Normal cluster need package to generate final conf
	if versionId != constants.FrontgateVersionId {
		ctx := context.Background()
		appManagerClient, err := appclient.NewAppManagerClient()
		if err != nil {
			p.Logger.Error("Connect to app manager failed: %+v", err)
			return nil, err
		}

		req := &pb.GetAppVersionPackageRequest{
			VersionId: pbutil.ToProtoString(versionId),
		}

		resp, err := appManagerClient.GetAppVersionPackage(ctx, req)
		if err != nil {
			p.Logger.Error("Get app version [%s] package failed: %+v", versionId, err)
			return nil, err
		}

		appPackage, err := devkit.LoadArchive(bytes.NewReader(resp.GetPackage()))
		if err != nil {
			p.Logger.Error("Load app version [%s] package failed: %+v", versionId, err)
			return nil, err
		}
		var confJson app.ClusterUserConfig
		err = jsonutil.Decode([]byte(conf), &confJson)
		if err != nil {
			p.Logger.Error("Parse conf [%s] failed: %+v", conf, err)
			return nil, err
		}
		clusterConf, err = appPackage.ClusterConfTemplate.Render(confJson)
		if err != nil {
			p.Logger.Error("Render app version [%s] cluster template failed: %+v", versionId, err)
			return nil, err
		}
		err = clusterConf.Validate()
		if err != nil {
			p.Logger.Error("Validate app version [%s] conf [%s] failed: %+v", versionId, conf, err)
			return nil, err
		}

	} else {
		err := jsonutil.Decode([]byte(conf), &clusterConf)
		if err != nil {
			p.Logger.Error("Parse conf [%s] to cluster failed: %+v", conf, err)
			return nil, err
		}
	}

	parser := Parser{Logger: p.Logger}
	clusterWrapper, err := parser.Parse(clusterConf)
	if err != nil {
		return nil, err
	}
	return clusterWrapper, nil
}

func (p *Provider) SplitJobIntoTasks(job *models.Job) (*models.TaskLayer, error) {
	frameInterface, err := vmbased.NewFrameInterface(job, p.Logger)
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

	default:
		p.Logger.Error("Unknown job action [%s]", job.JobAction)
		return nil, fmt.Errorf("unknown job action [%s]", job.JobAction)
	}
	return nil, nil
}

func (p *Provider) HandleSubtask(task *models.Task) error {
	handler := GetProviderHandler(p.Logger)

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
		p.Logger.Error("Unknown task action [%s]", task.TaskAction)
		return fmt.Errorf("unknown task action [%s]", task.TaskAction)
	}
}
func (p *Provider) WaitSubtask(task *models.Task, timeout time.Duration, waitInterval time.Duration) error {
	p.Logger.Debug("Wait sub task timeout [%s] interval [%s]", timeout, waitInterval)
	handler := GetProviderHandler(p.Logger)

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
		p.Logger.Error("Unknown task action [%s]", task.TaskAction)
		return fmt.Errorf("unknown task action [%s]", task.TaskAction)
	}
}

func (p *Provider) DescribeSubnets(ctx context.Context, req *pb.DescribeSubnetsRequest) (*pb.DescribeSubnetsResponse, error) {
	handler := GetProviderHandler(p.Logger)
	return handler.DescribeSubnets(ctx, req)
}

func (p *Provider) CheckResourceQuotas(ctx context.Context, clusterWrapper *models.ClusterWrapper) error {
	handler := GetProviderHandler(p.Logger)
	return handler.CheckResourceQuotas(ctx, clusterWrapper)
}

func (p *Provider) DescribeVpc(runtimeId, vpcId string) (*models.Vpc, error) {
	handler := GetProviderHandler(p.Logger)
	return handler.DescribeVpc(runtimeId, vpcId)
}

func (p *Provider) ValidateCredential(url, credential, zone string) error {
	handler := GetProviderHandler(p.Logger)
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

func (p *Provider) DescribeRuntimeProviderZones(url, credential string) ([]string, error) {
	handler := GetProviderHandler(p.Logger)
	return handler.DescribeZones(url, credential)
}
