// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package plugins

import (
	"fmt"
	"time"

	"openpitrix.io/openpitrix/pkg/models"
)

var providerPlugins map[string]ProviderInterface

func init() {
	providerPlugins = map[string]ProviderInterface{}
}

type ProviderInterface interface {
	// Parse package and conf into cluster which clusterManager will register into db.
	ParseClusterConf(versionId, conf string) (*models.ClusterWrapper, error)
	SplitJobIntoTasks(job *models.Job) (*models.TaskLayer, error)
	HandleSubtask(task *models.Task) error
	WaitSubtask(task *models.Task, timeout time.Duration, waitInterval time.Duration) error
	DescribeSubnet(subnetId string) (*models.Subnet, error)
	DescribeVpc(vpcId string) (*models.Vpc, error)
}

func RegisterProviderPlugin(provider string, providerInterface ProviderInterface) {
	providerPlugins[provider] = providerInterface
}

func GetRuntimePlugin(provider string) (ProviderInterface, error) {
	providerInterface, exists := providerPlugins[provider]
	if exists {
		return providerInterface, nil
	} else {
		return nil, fmt.Errorf("No such provider [%s]. ", provider)
	}
}
