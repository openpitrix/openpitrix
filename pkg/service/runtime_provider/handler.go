// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package runtime_provider

import (
	"context"
	"fmt"

	runtimeclient "openpitrix.io/openpitrix/pkg/client/runtime"
	providerclient "openpitrix.io/openpitrix/pkg/client/runtime_provider"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/pi"
	"openpitrix.io/openpitrix/pkg/plugins"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
)

func getProviderClientFromProvider(ctx context.Context, provider string) (pb.RuntimeProviderManagerClient, error) {
	providerConfig, isExist := pi.Global().GlobalConfig().RuntimeProvider[provider]
	if !isExist {
		return nil, fmt.Errorf("provider [%s] has not been registered yet", provider)
	}

	host := providerConfig.GetHost(provider)
	port := providerConfig.GetPort()
	return providerclient.NewRuntimeProviderClient(host, port)
}

func getProviderClient(ctx context.Context, runtimeId string) (pb.RuntimeProviderManagerClient, error) {
	runtime, err := runtimeclient.NewRuntime(ctx, runtimeId)
	if err != nil {
		return nil, err
	}

	provider := runtime.Runtime.Provider
	return getProviderClientFromProvider(ctx, provider)
}

func registerRuntimeProvider(ctx context.Context, provider, config string) error {
	err := pi.Global().RegisterRuntimeProvider(provider, config)
	if err != nil {
		return err
	}

	logger.Debug(ctx, "Available plugins: %+v", plugins.GetAvailablePlugins())
	return nil
}

func (p *Server) RegisterRuntimeProvider(ctx context.Context, req *pb.RegisterRuntimeProviderRequest) (*pb.RegisterRuntimeProviderResponse, error) {
	err := registerRuntimeProvider(ctx, req.GetProvider().GetValue(), req.GetConfig().GetValue())

	if err != nil {
		return &pb.RegisterRuntimeProviderResponse{
			Ok: pbutil.ToProtoBool(false),
		}, err
	} else {
		return &pb.RegisterRuntimeProviderResponse{
			Ok: pbutil.ToProtoBool(true),
		}, nil
	}
}

func (p *Server) ParseClusterConf(ctx context.Context, req *pb.ParseClusterConfRequest) (*pb.ParseClusterConfResponse, error) {
	runtimeId := req.GetRuntimeId().GetValue()
	providerClient, err := getProviderClient(ctx, runtimeId)
	if err != nil {
		return nil, err
	}

	return providerClient.ParseClusterConf(ctx, req)
}

func (p *Server) SplitJobIntoTasks(ctx context.Context, req *pb.SplitJobIntoTasksRequest) (*pb.SplitJobIntoTasksResponse, error) {
	runtimeId := req.GetRuntimeId().GetValue()
	providerClient, err := getProviderClient(ctx, runtimeId)
	if err != nil {
		return nil, err
	}

	return providerClient.SplitJobIntoTasks(ctx, req)
}

func (p *Server) HandleSubtask(ctx context.Context, req *pb.HandleSubtaskRequest) (*pb.HandleSubtaskResponse, error) {
	runtimeId := req.GetRuntimeId().GetValue()
	providerClient, err := getProviderClient(ctx, runtimeId)
	if err != nil {
		return nil, err
	}

	return providerClient.HandleSubtask(ctx, req)
}

func (p *Server) WaitSubtask(ctx context.Context, req *pb.WaitSubtaskRequest) (*pb.WaitSubtaskResponse, error) {
	runtimeId := req.GetRuntimeId().GetValue()
	providerClient, err := getProviderClient(ctx, runtimeId)
	if err != nil {
		return nil, err
	}

	return providerClient.WaitSubtask(ctx, req)
}

func (p *Server) DescribeSubnets(ctx context.Context, req *pb.DescribeSubnetsRequest) (*pb.DescribeSubnetsResponse, error) {
	runtimeId := req.GetRuntimeId().GetValue()
	providerClient, err := getProviderClient(ctx, runtimeId)
	if err != nil {
		return nil, err
	}

	return providerClient.DescribeSubnets(ctx, req)
}

func (p *Server) CheckResource(ctx context.Context, req *pb.CheckResourceRequest) (*pb.CheckResourceResponse, error) {
	runtimeId := req.GetRuntimeId().GetValue()
	providerClient, err := getProviderClient(ctx, runtimeId)
	if err != nil {
		return nil, err
	}

	return providerClient.CheckResource(ctx, req)
}

func (p *Server) DescribeVpc(ctx context.Context, req *pb.DescribeVpcRequest) (*pb.DescribeVpcResponse, error) {
	runtimeId := req.GetRuntimeId().GetValue()
	providerClient, err := getProviderClient(ctx, runtimeId)
	if err != nil {
		return nil, err
	}

	return providerClient.DescribeVpc(ctx, req)
}

func (p *Server) DescribeClusterDetails(ctx context.Context, req *pb.DescribeClusterDetailsRequest) (*pb.DescribeClusterDetailsResponse, error) {
	runtimeId := req.GetRuntimeId().GetValue()
	providerClient, err := getProviderClient(ctx, runtimeId)
	if err != nil {
		return nil, err
	}

	return providerClient.DescribeClusterDetails(ctx, req)
}

func (p *Server) ValidateRuntime(ctx context.Context, req *pb.ValidateRuntimeRequest) (*pb.ValidateRuntimeResponse, error) {
	providerClient, err := getProviderClientFromProvider(ctx, req.GetRuntimeCredential().GetProvider().GetValue())
	if err != nil {
		return nil, err
	}

	return providerClient.ValidateRuntime(ctx, req)
}

func (p *Server) DescribeZones(ctx context.Context, req *pb.DescribeZonesRequest) (*pb.DescribeZonesResponse, error) {
	providerClient, err := getProviderClientFromProvider(ctx, req.GetProvider().GetValue())
	if err != nil {
		return nil, err
	}

	return providerClient.DescribeZones(ctx, req)
}
