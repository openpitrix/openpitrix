// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package qingcloud

import (
	"context"
	"time"

	"github.com/yunify/qingcloud-sdk-go/config"
	"github.com/yunify/qingcloud-sdk-go/service"

	appClient "openpitrix.io/openpitrix/pkg/client/app"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/plugins"
	"openpitrix.io/openpitrix/pkg/utils"
)

func init() {
	plugins.RegisterRuntimePlugin(constants.RuntimeQingCloud, new(Runtime))
}

type Runtime struct {
}

func (p *Runtime) initService() (qingCloudService *service.QingCloudService, err error) {
	userConf, err := config.NewDefault()
	if err != nil {
		return
	}
	err = userConf.LoadUserConfig()
	if err != nil {
		return
	}
	qingCloudService, err = service.Init(userConf)
	if err != nil {
		return
	}
	return
}

func (p *Runtime) initJobService() (jobService *service.JobService, err error) {
	qingcloudService, err := p.initService()
	if err != nil {
		logger.Errorf("Failed to init qingcloud api service: %v", err)
		return
	}
	jobService, err = qingcloudService.Job(qingcloudService.Config.Zone)
	return
}

func (p *Runtime) ParseClusterConf(versionId, conf string) (*models.ClusterWrapper, error) {
	// Normal cluster need package to generate final conf
	if versionId != constants.FrontgateVersionId {
		ctx := context.Background()
		appManagerClient, err := appClient.NewAppManagerClient(ctx)
		if err != nil {
			logger.Errorf("Connect to app manager failed: %v", err)
			return nil, err
		}

		req := &pb.GetAppVersionPackageRequest{
			VersionId: utils.ToProtoString(versionId),
		}

		_, err = appManagerClient.GetAppVersionPackage(ctx, req)
		if err != nil {
			logger.Errorf("Get app version [%s] package failed: %v", versionId, err)
			return nil, err
		}

		// TODO after rendered, got the final conf
	}

	parser := Parser{}
	clusterWrapper, err := parser.Parse([]byte(conf))
	if err != nil {
		return nil, err
	}
	return clusterWrapper, nil
}

func (p *Runtime) SplitJobIntoTasks(job *models.Job) (*models.TaskLayer, error) {
	return nil, nil
}
func (p *Runtime) HandleSubtask(task *models.Task) error {
	return nil
}
func (p *Runtime) WaitSubtask(taskId string, timeout time.Duration, waitInterval time.Duration) error {
	return nil
}

func (p *Runtime) DescribeSubnet(subnetId string) (*models.Subnet, error) {
	return nil, nil
}
func (p *Runtime) DescribeVpc(vpcId string) (*models.Vpc, error) {
	return nil, nil
}
