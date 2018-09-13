// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package qingcloud

import (
	"context"
	"fmt"

	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/plugins/vmbased"
	"openpitrix.io/openpitrix/pkg/util/stringutil"
)

type Provider struct{}

func NewProvider() *Provider {
	return new(Provider)
}

func (p *Provider) ParseClusterConf(ctx context.Context, versionId, runtimeId, conf string, clusterWrapper *models.ClusterWrapper) error {
	frameInterface, err := vmbased.GetFrameInterface(ctx, nil)
	if err != nil {
		return err
	}
	return frameInterface.ParseClusterConf(versionId, runtimeId, conf, clusterWrapper)
}

func (p *Provider) SplitJobIntoTasks(ctx context.Context, job *models.Job) (*models.TaskLayer, error) {
	return vmbased.SplitJobIntoTasks(ctx, job)
}

func (p *Provider) HandleSubtask(ctx context.Context, task *models.Task) error {
	handler := GetProviderHandler(ctx)
	return vmbased.HandleSubtask(ctx, task, handler)
}

func (p *Provider) WaitSubtask(ctx context.Context, task *models.Task) error {
	handler := GetProviderHandler(ctx)
	return vmbased.WaitSubtask(ctx, task, handler)
}

func (p *Provider) DescribeSubnets(ctx context.Context, req *pb.DescribeSubnetsRequest) (*pb.DescribeSubnetsResponse, error) {
	handler := GetProviderHandler(ctx)
	return handler.DescribeSubnets(ctx, req)
}

func (p *Provider) CheckResource(ctx context.Context, clusterWrapper *models.ClusterWrapper) error {
	handler := GetProviderHandler(ctx)
	return handler.CheckResourceQuotas(ctx, clusterWrapper)
}

func (p *Provider) DescribeVpc(ctx context.Context, runtimeId, vpcId string) (*models.Vpc, error) {
	handler := GetProviderHandler(ctx)
	return handler.DescribeVpc(runtimeId, vpcId)
}

func (p *Provider) ValidateCredential(ctx context.Context, url, credential, zone string) error {
	handler := GetProviderHandler(ctx)
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

func (p *Provider) UpdateClusterStatus(ctx context.Context, job *models.Job) error {
	return nil
}

func (p *Provider) DescribeRuntimeProviderZones(ctx context.Context, url, credential string) ([]string, error) {
	handler := GetProviderHandler(ctx)
	return handler.DescribeZones(url, credential)
}

func (p *Provider) DescribeClusterDetails(ctx context.Context, cluster *models.ClusterWrapper) error {
	return nil
}
