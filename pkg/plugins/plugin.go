// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package plugins

import (
	"context"
	"fmt"
	"time"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/plugins/helm"
	"openpitrix.io/openpitrix/pkg/plugins/qingcloud"
)

var providerPlugins = make(map[string]ProviderInterface)

func init() {
	RegisterProviderPlugin(constants.ProviderQingCloud, qingcloud.NewProvider())
	RegisterProviderPlugin(constants.ProviderKubernetes, helm.NewProvider())
}

type ProviderInterface interface {
	SetLogger(logger *logger.Logger)
	// Parse package and conf into cluster which clusterManager will register into db.
	ParseClusterConf(versionId, conf string) (*models.ClusterWrapper, error)
	SplitJobIntoTasks(job *models.Job) (*models.TaskLayer, error)
	HandleSubtask(task *models.Task) error
	WaitSubtask(task *models.Task, timeout time.Duration, waitInterval time.Duration) error
	DescribeSubnets(ctx context.Context, req *pb.DescribeSubnetsRequest) (*pb.DescribeSubnetsResponse, error)
	CheckResourceQuotas(ctx context.Context, clusterWrapper *models.ClusterWrapper) (string, error)
	DescribeVpc(runtimeId, vpcId string) (*models.Vpc, error)
	ValidateCredential(url, credential, zone string) error
	DescribeRuntimeProviderZones(url, credential string) ([]string, error)
	UpdateClusterStatus(job *models.Job) error
}

func RegisterProviderPlugin(provider string, providerInterface ProviderInterface) {
	providerPlugins[provider] = providerInterface
}

func GetProviderPlugin(provider string, logger *logger.Logger) (ProviderInterface, error) {
	providerInterface, exists := providerPlugins[provider]
	if exists {
		providerInterface.SetLogger(logger)
		return providerInterface, nil
	} else {
		return nil, fmt.Errorf("No such provider [%s]. ", provider)
	}
}
