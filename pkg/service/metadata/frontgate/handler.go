// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package frontgate

import (
	"context"
	"fmt"

	"github.com/gogo/protobuf/proto"

	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/pb/frontgate"
	"openpitrix.io/openpitrix/pkg/pb/types"
	"openpitrix.io/openpitrix/pkg/service/metadata/drone/droneutil"
	"openpitrix.io/openpitrix/pkg/service/pilot/pilotutil"
)

var (
	_ pbfrontgate.FrontgateService = (*Server)(nil)
)

func (p *Server) GetPilotConfig(in *pbtypes.Empty, out *pbtypes.PilotConfig) error {
	ctx := context.Background()

	client, conn, err := pilotutil.DialPilotService(ctx, p.cfg.PilotHost, int(p.cfg.PilotPort))
	if err != nil {
		return err
	}
	defer conn.Close()

	info, err := client.GetPilotConfig(ctx, &pbtypes.Empty{})
	if err != nil {
		return err
	}

	*out = *info
	return nil
}

func (p *Server) GetFrontgateConfig(in *pbtypes.Empty, out *pbtypes.FrontgateConfig) error {
	*out = *proto.Clone(p.cfg).(*pbtypes.FrontgateConfig)
	return nil
}

func (p *Server) GetDroneList(in *pbtypes.Empty, out *pbtypes.DroneIdList) error {
	return fmt.Errorf("TODO")
}

func (p *Server) GetDroneConfig(in *pbtypes.DroneEndpoint, out *pbtypes.DroneConfig) error {
	ctx := context.Background()

	client, conn, err := droneutil.DialDroneService(ctx,
		in.GetDroneIp(),
		int(in.GetDronePort()),
	)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = client.GetDroneConfig(ctx, &pbtypes.Empty{})
	if err != nil {
		return err
	}

	return nil
}

func (p *Server) IsConfdRunning(in *pbtypes.ConfdEndpoint, out *pbtypes.Bool) error {
	ctx := context.Background()

	client, conn, err := droneutil.DialDroneService(ctx,
		in.GetDroneIp(),
		int(in.GetDronePort()),
	)
	if err != nil {
		return err
	}
	defer conn.Close()

	reply, err := client.IsConfdRunning(ctx, &pbtypes.Empty{})
	if err != nil {
		return err
	}

	out.Value = reply.GetValue()
	return nil
}

func (p *Server) StartConfd(in *pbtypes.ConfdEndpoint, out *pbtypes.Empty) error {
	ctx := context.Background()

	client, conn, err := droneutil.DialDroneService(ctx,
		in.GetDroneIp(),
		int(in.GetDronePort()),
	)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = client.SetConfdConfig(ctx, p.cfg.GetConfdConfig())
	if err != nil {
		return err
	}

	_, err = client.StartConfd(ctx, &pbtypes.Empty{})
	if err != nil {
		return err
	}

	return nil
}

func (p *Server) StopConfd(in *pbtypes.ConfdEndpoint, out *pbtypes.Empty) error {
	ctx := context.Background()

	client, conn, err := droneutil.DialDroneService(ctx, in.GetDroneIp(), int(in.GetDronePort()))
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = client.StopConfd(ctx, &pbtypes.Empty{})
	if err != nil {
		return err
	}

	return nil
}

func (p *Server) RegisterMetadata(in *pbtypes.SubTask_RegisterMetadata, out *pbtypes.Empty) error {
	return p.etcd.RegisterMetadata(in)
}

func (p *EtcdClient) RegisterMetadata(in *pbtypes.SubTask_RegisterMetadata) error {
	err := p.SetValues(in.Cnodes)
	if err != nil {
		return err
	}

	return nil
}

func (p *Server) DeregisterMetadata(in *pbtypes.SubTask_DeregisterMetadata, out *pbtypes.Empty) error {
	return p.etcd.DeregisterMetadata(in)
}

func (p *EtcdClient) DeregisterMetadata(in *pbtypes.SubTask_DeregisterMetadata) error {
	var keyPrefixs []string
	for key, _ := range in.Cnodes {
		keyPrefixs = append(keyPrefixs, key)
	}

	err := p.DelValuesWithPrefix(keyPrefixs...)
	if err != nil {
		return err
	}

	return nil
}

func (p *Server) RegisterCmd(in *pbtypes.SubTask_RegisterCmd, out *pbtypes.Empty) error {
	return p.etcd.RegisterCmd(in)
}

func (p *EtcdClient) RegisterCmd(in *pbtypes.SubTask_RegisterCmd) error {
	err := p.SetValues(in.Cnodes)
	if err != nil {
		return err
	}

	return nil
}

func (p *Server) DeregisterCmd(in *pbtypes.SubTask_DeregisterCmd, out *pbtypes.Empty) error {
	return p.etcd.DeregisterCmd(in)
}

func (p *EtcdClient) DeregisterCmd(in *pbtypes.SubTask_DeregisterCmd) error {
	var keyPrefixs []string
	for key, _ := range in.Cnodes {
		keyPrefixs = append(keyPrefixs, key)
	}

	err := p.DelValuesWithPrefix(keyPrefixs...)
	if err != nil {
		return err
	}

	return nil
}

func (p *Server) ReportSubTaskStatus(in *pbtypes.SubTaskStatus, out *pbtypes.Empty) error {
	ctx := context.Background()

	client, conn, err := pilotutil.DialPilotService(ctx, p.cfg.PilotHost, int(p.cfg.PilotPort))
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = client.ReportSubTaskStatus(ctx, &pbtypes.SubTaskStatus{}) // todo
	if err != nil {
		return err
	}

	return nil
}

func (p *Server) ClosePilotChannel(in *pbtypes.Empty, out *pbtypes.Empty) error {
	return p.conn.Close()
}

func (p *Server) GetEtcdValuesByPrefix(in *pbtypes.String, out *pbtypes.StringMap) error {
	m, err := p.etcd.GetValuesByPrefix(in.Value)
	if err != nil {
		return err
	}

	out.ValueMap = m
	return nil
}

func (p *Server) GetEtcdValues(in *pbtypes.StringList, out *pbtypes.StringMap) error {
	m, err := p.etcd.GetValues(in.ValueList...)
	if err != nil {
		return err
	}

	out.ValueMap = m
	return nil
}

func (p *Server) SetEtcdValues(in *pbtypes.StringMap, out *pbtypes.Empty) error {
	return p.etcd.SetValues(in.GetValueMap())
}

func (p *Server) PingPilot(in *pbtypes.Empty, out *pbtypes.Empty) error {
	ctx := context.Background()

	client, conn, err := pilotutil.DialPilotService(ctx, p.cfg.PilotHost, int(p.cfg.PilotPort))
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = client.PingPilot(ctx, &pbtypes.Empty{}) // todo
	if err != nil {
		return err
	}

	return nil
}

func (p *Server) PingFrontgate(in *pbtypes.Empty, out *pbtypes.Empty) error {
	logger.Info("PingFrontgate: ok")
	return nil // OK
}

func (p *Server) PingDrone(in *pbtypes.DroneEndpoint, out *pbtypes.Empty) error {
	ctx := context.Background()

	client, conn, err := droneutil.DialDroneService(ctx,
		in.GetDroneIp(),
		int(in.GetDronePort()),
	)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = client.PingDrone(ctx, &pbtypes.Empty{})
	if err != nil {
		return err
	}

	return nil
}

func (p *Server) HeartBeat(in *pbtypes.Empty, out *pbtypes.Empty) error {
	return nil // OK
}
