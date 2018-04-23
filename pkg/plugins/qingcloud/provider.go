// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package qingcloud

import (
	"bytes"
	"context"
	"encoding/json"
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
	"openpitrix.io/openpitrix/pkg/utils"
)

type Provider struct {
}

func (p *Provider) ParseClusterConf(versionId, conf string) (*models.ClusterWrapper, error) {
	clusterConf := app.ClusterConf{}
	// Normal cluster need package to generate final conf
	if versionId != constants.FrontgateVersionId {
		ctx := context.Background()
		appManagerClient, err := appclient.NewAppManagerClient(ctx)
		if err != nil {
			logger.Errorf("Connect to app manager failed: %+v", err)
			return nil, err
		}

		req := &pb.GetAppVersionPackageRequest{
			VersionId: utils.ToProtoString(versionId),
		}

		resp, err := appManagerClient.GetAppVersionPackage(ctx, req)
		if err != nil {
			logger.Errorf("Get app version [%s] package failed: %+v", versionId, err)
			return nil, err
		}

		appPackage, err := devkit.LoadArchive(bytes.NewReader(resp.GetPackage()))
		if err != nil {
			logger.Errorf("Load app version [%s] package failed: %+v", versionId, err)
			return nil, err
		}
		var confJson app.ClusterUserConfig
		err = json.Unmarshal([]byte(conf), &confJson)
		if err != nil {
			logger.Errorf("Parse conf [%s] failed: %+v", conf, err)
			return nil, err
		}
		clusterConf, err = appPackage.ClusterConfTemplate.Render(confJson)
		if err != nil {
			logger.Errorf("Render app version [%s] cluster template failed: %+v", versionId, err)
			return nil, err
		}
		err = clusterConf.Validate()
		if err != nil {
			logger.Errorf("Validate app version [%s] conf [%s] failed: %+v", versionId, conf, err)
			return nil, err
		}

	} else {
		err := json.Unmarshal([]byte(conf), &clusterConf)
		if err != nil {
			logger.Errorf("Parse conf [%s] to cluster failed: %+v", conf, err)
			return nil, err
		}
	}

	parser := Parser{}
	clusterWrapper, err := parser.Parse(clusterConf)
	if err != nil {
		return nil, err
	}
	return clusterWrapper, nil
}

func (p *Provider) SplitJobIntoTasks(job *models.Job) (*models.TaskLayer, error) {
	frameInterface, err := vmbased.NewFrameInterface(job)
	if err != nil {
		return nil, err
	}

	switch job.JobAction {
	case constants.ActionCreateCluster:
		// TODO: vpc, eip, vxnet

		return frameInterface.CreateClusterLayer(), nil
	case constants.ActionUpgradeCluster:
		// not supported yet
		return nil, nil
	case constants.ActionRollbackCluster:
		// not supported yet
		return nil, nil
	case constants.ActionResizeCluster:

	case constants.ActionAddClusterNodes:

	case constants.ActionDeleteClusterNodes:

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
		logger.Errorf("Unknown job action [%s]", job.JobAction)
		return nil, fmt.Errorf("unknown job action [%s]", job.JobAction)
	}
	return nil, nil
}

func (p *Provider) HandleSubtask(task *models.Task) error {
	handler := new(ProviderHandler)

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
		logger.Errorf("Unknown task action [%s]", task.TaskAction)
		return fmt.Errorf("unknown task action [%s]", task.TaskAction)
	}
}
func (p *Provider) WaitSubtask(task *models.Task, timeout time.Duration, waitInterval time.Duration) error {
	logger.Debugf("Wait sub task [%s] timeout [%s] interval [%s]", task.TaskId, timeout, waitInterval)
	handler := new(ProviderHandler)

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
		logger.Errorf("Unknown task action [%s]", task.TaskAction)
		return fmt.Errorf("unknown task action [%s]", task.TaskAction)
	}
}

func (p *Provider) DescribeSubnet(runtimeId, subnetId string) (*models.Subnet, error) {
	handler := new(ProviderHandler)
	return handler.DescribeSubnet(runtimeId, subnetId)
}

func (p *Provider) DescribeVpc(runtimeId, vpcId string) (*models.Vpc, error) {
	handler := new(ProviderHandler)
	return handler.DescribeVpc(runtimeId, vpcId)
}

func (p *Provider) ValidateCredential(url, credential string) error {
	handler := new(ProviderHandler)
	_, err := handler.DescribeZones(url, credential)
	return err
}

func (p *Provider) UpdateClusterStatus(job *models.Job) error {
	return nil
}

func (p *Provider) DescribeRuntimeProviderZones(url, credential string) []string {
	handler := new(ProviderHandler)
	zones, _ := handler.DescribeZones(url, credential)
	return zones
}
