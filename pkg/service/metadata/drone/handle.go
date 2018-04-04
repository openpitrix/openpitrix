// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package drone

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/golang/protobuf/ptypes"

	"openpitrix.io/libconfd"
	"openpitrix.io/openpitrix/pkg/constants"
	pbdrone "openpitrix.io/openpitrix/pkg/pb/drone"
)

var (
	_ pbdrone.DroneServiceServer             = (*Server)(nil)
	_ pbdrone.DroneServiceForFrontgateServer = (*Server)(nil)
)

func (p *Server) GetConfdConfig(context.Context, *pbdrone.Empty) (*pbdrone.ConfdConfig, error) {
	cfg := p.confd.GetConfdConfig()
	reply := To_pbdrone_ConfdConfig(cfg)
	return reply, nil
}

func (p *Server) GetBackendConfig(context.Context, *pbdrone.Empty) (*pbdrone.ConfdBackendConfig, error) {
	bcfg := p.confd.GetBackendConfig()
	reply := To_pbdrone_ConfdBackendConfig(bcfg)
	return reply, nil
}

func (p *Server) GetInfo(context.Context, *pbdrone.Empty) (*pbdrone.Info, error) {
	cfg := p.confd.GetConfdConfig()
	bcfg := p.confd.GetBackendConfig()

	reply := &pbdrone.Info{
		DroneIp:            getLocalIP(),
		ConfdConfig:        To_pbdrone_ConfdConfig(cfg),
		ConfdBackendConfig: To_pbdrone_ConfdBackendConfig(bcfg),
	}
	return reply, nil
}

func (p *Server) StartConfd(ctx context.Context, arg *pbdrone.StartConfdRequest) (*pbdrone.Empty, error) {
	cfg := To_libconfd_Config(arg.GetConfdConfig())
	bcfg := To_libconfd_BackendConfig(arg.GetConfdBackendConfig())

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
	status := constants.StatusStopped
	if p.confd.IsRunning() {
		status = constants.StatusWorking
	}

	reply := &pbdrone.ConfdStatus{
		ProcessId: int32(os.Getpid()),
		Status:    status,
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

func (p *Server) GetTemplateFiles(ctx context.Context, arg *pbdrone.GetTemplateFilesRequest) (*pbdrone.GetTemplateFilesResponse, error) {
	if !p.confd.IsRunning() {
		return nil, fmt.Errorf("drone: confd is not started!")
	}

	confdConfig := p.confd.GetConfdConfig()
	_, paths, err := libconfd.ListTemplateResource(confdConfig.ConfDir)
	if err != nil {
		return nil, err
	}

	re := arg.GetRegexp()
	reply := &pbdrone.GetTemplateFilesResponse{}

	for _, s := range paths {
		basename := filepath.Base(s)
		if re == "" {
			reply.Files = append(reply.Files, basename)
			continue
		}
		matched, err := regexp.MatchString(re, basename)
		if err != nil {
			return nil, err
		}
		if matched {
			reply.Files = append(reply.Files, basename)
		}
	}

	return reply, nil
}
func (p *Server) GetValues(ctx context.Context, arg *pbdrone.GetValuesRequest) (*pbdrone.GetValuesResponse, error) {
	if !p.confd.IsRunning() {
		return nil, fmt.Errorf("drone: confd is not started!")
	}

	client := p.confd.GetBackendClient()
	if client == nil {
		return nil, fmt.Errorf("drone: confd is not started!")
	}

	keys := arg.GetKeys()
	m, err := client.GetValues(keys)
	if err != nil {
		return nil, err
	}

	reply := &pbdrone.GetValuesResponse{
		Values: make(map[string]string),
	}

	for _, k := range keys {
		reply.Values[k] = m[k]
	}

	return reply, nil
}
