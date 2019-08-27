// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package plugins

import (
	"context"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/pi"
)

type ProviderInterface interface {
	// cluster
	// Parse package and conf into cluster which clusterManager will register into db.
	ParseClusterConf(ctx context.Context, versionId, runtimeId, conf string, clusterWrapper *models.ClusterWrapper) error
	SplitJobIntoTasks(ctx context.Context, job *models.Job) (*models.TaskLayer, error)
	HandleSubtask(ctx context.Context, task *models.Task) error
	WaitSubtask(ctx context.Context, task *models.Task) error
	DescribeSubnets(ctx context.Context, req *pb.DescribeSubnetsRequest) (*pb.DescribeSubnetsResponse, error)
	CheckResource(ctx context.Context, clusterWrapper *models.ClusterWrapper) error
	DescribeVpc(ctx context.Context, runtimeId, vpcId string) (*models.Vpc, error)
	DescribeClusterDetails(ctx context.Context, cluster *models.ClusterWrapper) (*models.ClusterWrapper, error)
	// runtime
	ValidateRuntime(ctx context.Context, runtimeId, zone string, runtimeCredential *models.RuntimeCredential, needCreate bool) error
	DescribeRuntimeProviderZones(ctx context.Context, runtimeCredential *models.RuntimeCredential) ([]string, error)
}

func GetAvailablePlugins() []string {
	var plugins []string
	for provider := range pi.Global().GlobalConfig().RuntimeProvider {
		plugins = append(plugins, provider)
	}
	return plugins
}

func IsVmbasedProviders(provider string) bool {
	providerConfig, isExist := pi.Global().GlobalConfig().RuntimeProvider[provider]
	if !isExist {
		return false
	}

	if providerConfig.ProviderType == constants.ProviderTypeVmbased {
		return true
	}
	return false
}
