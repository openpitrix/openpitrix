// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package drone

import (
	"context"

	"github.com/golang/protobuf/ptypes"

	"openpitrix.io/libconfd"
	pbdrone "openpitrix.io/openpitrix/pkg/pb/drone"
)

var (
	_ pbdrone.DroneServiceServer             = (*Server)(nil)
	_ pbdrone.DroneServiceForFrontgateServer = (*Server)(nil)
)

func (p *Server) GetConfdConfig(context.Context, *pbdrone.Empty) (*pbdrone.ConfdConfig, error) {
	panic("TODO")
}

func (p *Server) GetBackendConfig(context.Context, *pbdrone.Empty) (*pbdrone.ConfdBackendConfig, error) {
	panic("TODO")
}

func (p *Server) GetInfo(context.Context, *pbdrone.Empty) (*pbdrone.Info, error) {
	cfg := p.confd.GetConfdConfig()
	bcfg := p.confd.GetBackendConfig()

	reply := &pbdrone.Info{
		DroneIp: getLocalIP(),
		ConfdConfig: &pbdrone.ConfdConfig{
			ConfDir:  cfg.ConfDir,
			Interval: int32(cfg.Interval),
			Prefix:   cfg.Prefix,
			SyncOnly: cfg.SyncOnly,
			LogLevel: cfg.LogLevel,
		},
		ConfdBackendConfig: &pbdrone.ConfdBackendConfig{
			Type:         bcfg.Type,
			Host:         append([]string{}, bcfg.Host...),
			Username:     bcfg.UserName,
			Password:     bcfg.Password,
			ClientCaKeys: bcfg.ClientCAKeys,
			ClientCert:   bcfg.ClientCert,
			ClientKey:    bcfg.ClientKey,
		},
	}
	return reply, nil
}

func (p *Server) StartConfd(ctx context.Context, arg *pbdrone.StartConfdRequest) (*pbdrone.Empty, error) {
	cfg := &libconfd.Config{}         // todo
	bcfg := &libconfd.BackendConfig{} // todo

	if err := p.confd.Start(cfg, bcfg); err != nil {
		return nil, err
	}

	return &pbdrone.Empty{}, nil
}

func (p *Server) StopConfd(ctx context.Context, arg *pbdrone.Empty) (*pbdrone.Empty, error) {
	if err := p.confd.Stop(); err != nil {
		return nil, err
	}

	return &pbdrone.Empty{}, nil
}

func (p *Server) GetConfdStatus(ctx context.Context, arg *pbdrone.Empty) (*pbdrone.ConfdStatus, error) {
	reply := &pbdrone.ConfdStatus{
		ProcessId: "",
		Status:    "",
		UpTime:    ptypes.TimestampNow(),
	}
	return reply, nil
}

func (p *Server) SubscribeCmdStatus(*pbdrone.SubscribeCmdStatusRequest, pbdrone.DroneServiceForFrontgate_SubscribeCmdStatusServer) error {
	panic("todo")
}

func (p *Server) GetRegisterCmdStatus(context.Context, *pbdrone.GetRegisterCmdStatusRequest) (*pbdrone.GetRegisterCmdStatusResponse, error) {
	panic("todo")
}

func (p *Server) GetTemplateFiles(context.Context, *pbdrone.GetTemplateFilesRequest) (*pbdrone.GetTemplateFilesResponse, error) {
	panic("todo")
}
func (p *Server) GetValues(context.Context, *pbdrone.GetValuesRequest) (*pbdrone.GetValuesResponse, error) {
	panic("todo")
}
