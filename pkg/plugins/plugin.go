// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package plugins

import (
	"context"
	"fmt"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/plugins/aliyun"
	"openpitrix.io/openpitrix/pkg/plugins/aws"
	"openpitrix.io/openpitrix/pkg/plugins/helm"
	"openpitrix.io/openpitrix/pkg/plugins/qingcloud"
	"openpitrix.io/openpitrix/pkg/util/stringutil"
)

var SupportedProviders = make(map[string]interface{})

func init() {
	RegisterProvider(constants.ProviderKubernetes, helm.NewProvider())
	RegisterProvider(constants.ProviderQingCloud, qingcloud.NewProvider())
	RegisterProvider(constants.ProviderAWS, aws.NewProvider())
	RegisterProvider(constants.ProviderAliyun, aliyun.NewProvider())
}

type ProviderInterface interface {
	// Parse package and conf into cluster which clusterManager will register into db.
	ParseClusterConf(ctx context.Context, versionId, runtimeId, conf string, clusterWrapper *models.ClusterWrapper) error
	SplitJobIntoTasks(ctx context.Context, job *models.Job) (*models.TaskLayer, error)
	HandleSubtask(ctx context.Context, task *models.Task) error
	WaitSubtask(ctx context.Context, task *models.Task) error
	DescribeSubnets(ctx context.Context, req *pb.DescribeSubnetsRequest) (*pb.DescribeSubnetsResponse, error)
	CheckResource(ctx context.Context, clusterWrapper *models.ClusterWrapper) error
	DescribeVpc(ctx context.Context, runtimeId, vpcId string) (*models.Vpc, error)
	ValidateCredential(ctx context.Context, runtimeId, url, credential, zone string) error
	DescribeRuntimeProviderZones(ctx context.Context, url, credential string) ([]string, error)
	UpdateClusterStatus(ctx context.Context, job *models.Job) error
	DescribeClusterDetails(ctx context.Context, cluster *models.ClusterWrapper) error
}

func RegisterProvider(provider string, providerInterface ProviderInterface) {
	SupportedProviders[provider] = providerInterface
}

func GetProviderPlugin(ctx context.Context, provider string) (ProviderInterface, error) {
	providerInterface, isExist := SupportedProviders[provider]
	if isExist {
		value, ok := providerInterface.(ProviderInterface)
		if !ok {
			logger.Error(ctx, "No such provider interface [%s].", provider)
			return nil, fmt.Errorf("no such provider [%s]", provider)
		}
		return value, nil
	} else {
		logger.Error(ctx, "No such provider [%s].", provider)
		return nil, fmt.Errorf("no such provider [%s]", provider)
	}
}

func GetAvailablePlugins(availablePlugins []string) []string {
	var plugins []string
	for plugin := range SupportedProviders {
		plugins = append(plugins, plugin)
	}
	if len(availablePlugins) == 0 {
		return plugins
	} else {
		var intersection []string
		for _, plugin := range plugins {
			if stringutil.StringIn(plugin, availablePlugins) {
				intersection = append(intersection, plugin)
			}
		}
		return intersection
	}
}

func IsVmbasedProviders(provider string) bool {
	if provider != constants.ProviderKubernetes {
		return true
	}
	return false
}
