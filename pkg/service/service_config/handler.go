// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package service_config

import (
	"context"
	"fmt"

	nfpb "openpitrix.io/notification/pkg/pb"
	nfclient "openpitrix.io/openpitrix/pkg/client/notification"
	"openpitrix.io/openpitrix/pkg/config"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/gerr"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/pi"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
)

func OpToNfConfig(opConfig *pb.NotificationConfig) *nfpb.ServiceConfig {
	return &nfpb.ServiceConfig{
		EmailServiceConfig: &nfpb.EmailServiceConfig{
			Protocol:      opConfig.EmailServiceConfig.Protocol,
			EmailHost:     opConfig.EmailServiceConfig.EmailHost,
			Port:          opConfig.EmailServiceConfig.Port,
			DisplaySender: opConfig.EmailServiceConfig.DisplaySender,
			Email:         opConfig.EmailServiceConfig.Email,
			Password:      opConfig.EmailServiceConfig.Password,
			SslEnable:     opConfig.EmailServiceConfig.SslEnable,
		},
	}
}

func NfToOpConfig(nfConfig *nfpb.ServiceConfig) *pb.NotificationConfig {
	return &pb.NotificationConfig{
		EmailServiceConfig: &pb.EmailServiceConfig{
			Protocol:      nfConfig.EmailServiceConfig.Protocol,
			EmailHost:     nfConfig.EmailServiceConfig.EmailHost,
			Port:          nfConfig.EmailServiceConfig.Port,
			DisplaySender: nfConfig.EmailServiceConfig.DisplaySender,
			Email:         nfConfig.EmailServiceConfig.Email,
			Password:      nfConfig.EmailServiceConfig.Password,
			SslEnable:     nfConfig.EmailServiceConfig.SslEnable,
		},
	}
}

func (p *Server) SetServiceConfig(ctx context.Context, req *pb.SetServiceConfigRequest) (*pb.SetServiceConfigResponse, error) {
	if req.NotificationConfig != nil && req.NotificationConfig.EmailServiceConfig != nil {
		nfClient, err := nfclient.NewClient()
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorSetNotificationConfig)
		}
		response, err := nfClient.SetServiceConfig(ctx, OpToNfConfig(req.NotificationConfig))
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorSetNotificationConfig)
		}
		if !response.GetIsSucc().GetValue() {
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorSetNotificationConfig)
		}
	} else if req.RuntimeConfig != nil {
		for _, config := range req.RuntimeConfig.ConfigSet {
			name := config.GetName().GetValue()
			enable := config.GetEnable().GetValue()
			runtimeProviderConfig, isExist := pi.Global().GlobalConfig().Runtime[name]
			if !isExist {
				return nil, gerr.New(ctx, gerr.InvalidArgument, gerr.ErrorUnsupportedRuntimeProvider, name)
			}
			runtimeProviderConfig.Enable = enable
		}
		err := pi.Global().SetGlobalCfg(ctx)
		if err != nil {
			return nil, gerr.New(ctx, gerr.Internal, gerr.ErrorSetServiceConfig)
		}
	} else if req.BasicConfig != nil {
		basicCfg := config.BasicConfig{
			PlatformName: req.BasicConfig.GetPlatformName().GetValue(),
			PlatformUrl:  req.BasicConfig.GetPlatformUrl().GetValue(),
		}
		err := pi.Global().SetBasicCfg(basicCfg)
		if err != nil {
			return nil, gerr.New(ctx, gerr.Internal, gerr.ErrorSetServiceConfig)
		}

	} else {
		err := fmt.Errorf("need service config to set")
		return nil, gerr.NewWithDetail(ctx, gerr.InvalidArgument, err, gerr.ErrorSetServiceConfig)
	}

	return &pb.SetServiceConfigResponse{
		IsSucc: pbutil.ToProtoBool(true),
	}, nil
}

func (p *Server) GetServiceConfig(ctx context.Context, req *pb.GetServiceConfigRequest) (*pb.GetServiceConfigResponse, error) {
	if len(req.ServiceType) == 0 {
		req.ServiceType = constants.ServiceTypes
	}

	serviceConfigResponse := new(pb.GetServiceConfigResponse)
	for _, serviceType := range req.GetServiceType() {
		switch serviceType {
		case constants.ServiceTypeNotification:
			nfClient, err := nfclient.NewClient()
			if err != nil {
				return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorGetNotificationConfig)
			}
			// empty means all config
			response, err := nfClient.GetServiceConfig(ctx, &nfpb.GetServiceConfigRequest{
				ServiceType: []string{},
			})
			if err != nil {
				return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorGetNotificationConfig)
			}
			serviceConfigResponse.NotificationConfig = NfToOpConfig(response)
		case constants.ServiceTypeRuntime:
			var configs []*pb.RuntimeItemConfig
			for name, runtimeProviderConfig := range pi.Global().GlobalConfig().Runtime {
				configs = append(configs, &pb.RuntimeItemConfig{
					Name:   pbutil.ToProtoString(name),
					Enable: pbutil.ToProtoBool(runtimeProviderConfig.Enable),
				})
			}
			serviceConfigResponse.RuntimeConfig = &pb.RuntimeConfig{
				ConfigSet: configs,
			}
		case constants.ServiceTypeBasicConfig:
			basicCfg := pi.Global().GetBasicCfg()
			serviceConfigResponse.BasicConfig = &pb.BasicConfig{
				PlatformName: pbutil.ToProtoString(basicCfg.PlatformName),
				PlatformUrl:  pbutil.ToProtoString(basicCfg.PlatformUrl),
			}
		}
	}
	return serviceConfigResponse, nil
}
