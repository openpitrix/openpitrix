// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package qingcloud

import (
	"time"

	"github.com/yunify/qingcloud-sdk-go/config"
	"github.com/yunify/qingcloud-sdk-go/service"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/plugins"
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

func (p *Runtime) ParseClusterConf(versionId, conf string) (*models.Cluster, error) {
	return nil, nil
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
