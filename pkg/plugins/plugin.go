// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package plugins

import (
	"context"
	"fmt"
	"time"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/plugins/aws"
	"openpitrix.io/openpitrix/pkg/plugins/helm"
	"openpitrix.io/openpitrix/pkg/plugins/qingcloud"
)

type ProviderInterface interface {
	// Parse package and conf into cluster which clusterManager will register into db.
	ParseClusterConf(versionId, runtimeId, conf string, clusterWrapper *models.ClusterWrapper) error
	SplitJobIntoTasks(job *models.Job) (*models.TaskLayer, error)
	HandleSubtask(task *models.Task) error
	WaitSubtask(task *models.Task, timeout time.Duration, waitInterval time.Duration) error
	DescribeSubnets(ctx context.Context, req *pb.DescribeSubnetsRequest) (*pb.DescribeSubnetsResponse, error)
	CheckResource(ctx context.Context, clusterWrapper *models.ClusterWrapper) error
	DescribeVpc(runtimeId, vpcId string) (*models.Vpc, error)
	ValidateCredential(url, credential, zone string) error
	DescribeRuntimeProviderZones(url, credential string) ([]string, error)
	UpdateClusterStatus(job *models.Job) error
	DescribeClusterDetails(ctx context.Context, cluster *models.ClusterWrapper) error
}

func GetProviderPlugin(ctx context.Context, provider string) (ProviderInterface, error) {
	var providerInterface ProviderInterface

	switch provider {
	case constants.ProviderQingCloud:
		providerInterface = qingcloud.NewProvider(ctx)
	case constants.ProviderKubernetes:
		providerInterface = helm.NewProvider(ctx)
	case constants.ProviderAWS:
		providerInterface = aws.NewProvider(ctx)
	default:
		return nil, fmt.Errorf("No such provider [%s]. ", provider)
	}
	return providerInterface, nil
}
