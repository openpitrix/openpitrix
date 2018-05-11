// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package frontgate

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/chai2010/jsonmap"

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

	cfg := p.cfg.Get()
	client, conn, err := pilotutil.DialPilotService(ctx, cfg.PilotHost, int(cfg.PilotPort))
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
	*out = *p.cfg.Get()
	return nil
}

func (p *Server) SetFrontgateConfig(cfg *pbtypes.FrontgateConfig, out *pbtypes.Empty) error {
	if reflect.DeepEqual(cfg, p.cfg.Get()) {
		return nil
	}

	if err := p.cfg.Set(cfg); err != nil {
		return err
	}

	if err := p.cfg.Save(); err != nil {
		return err
	}

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

func (p *Server) SetDroneConfig(in *pbtypes.SetDroneConfigRequest, out *pbtypes.Empty) error {
	ctx := context.Background()

	client, conn, err := droneutil.DialDroneService(ctx,
		in.Endpoint.GetDroneIp(),
		int(in.Endpoint.GetDronePort()),
	)
	if err != nil {
		return err
	}
	defer conn.Close()

	// 1. set drone config
	_, err = client.SetDroneConfig(ctx, in.GetConfig())
	if err != nil {
		return err
	}

	// 2. set confd config
	_, err = client.SetConfdConfig(ctx, p.cfg.Get().GetConfdConfig())
	if err != nil {
		return err
	}

	return nil
}

func (p *Server) GetConfdConfig(in *pbtypes.ConfdEndpoint, out *pbtypes.ConfdConfig) error {
	*out = *p.cfg.Get().GetConfdConfig()
	return nil
}
func (p *Server) SetConfdConfig(in *pbtypes.SetConfdConfigRequest, out *pbtypes.Empty) error {
	ctx := context.Background()

	client, conn, err := droneutil.DialDroneService(ctx,
		in.Endpoint.GetDroneIp(),
		int(in.Endpoint.GetDronePort()),
	)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = client.SetConfdConfig(ctx, in.GetConfig())
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

	_, err = client.SetConfdConfig(ctx, p.cfg.Get().GetConfdConfig())
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
	etcdClient, err := p.etcd.GetClient(p.cfg.Get().GetConfdConfig().GetBackendConfig().GetHost(), time.Second)
	if err != nil {
		return err
	}

	return etcdClient.RegisterMetadata(in)
}

func (p *EtcdClient) RegisterMetadata(in *pbtypes.SubTask_RegisterMetadata) error {
	var m = map[string]interface{}{}
	err := json.Unmarshal([]byte(in.Cnodes), &m)
	if err != nil {
		return err
	}

	err = p.SetValues(jsonmap.JsonMap(m).ToMapString("/"))
	if err != nil {
		return err
	}

	return nil
}

func (p *Server) DeregisterMetadata(in *pbtypes.SubTask_DeregisterMetadata, out *pbtypes.Empty) error {
	etcdClient, err := p.etcd.GetClient(p.cfg.Get().GetConfdConfig().GetBackendConfig().GetHost(), time.Second)
	if err != nil {
		return err
	}
	return etcdClient.DeregisterMetadata(in)
}

func (p *EtcdClient) DeregisterMetadata(in *pbtypes.SubTask_DeregisterMetadata) error {
	var m = map[string]interface{}{}
	err := json.Unmarshal([]byte(in.Cnodes), &m)
	if err != nil {
		return err
	}

	var keyPrefixs []string
	for key, _ := range m {
		keyPrefixs = append(keyPrefixs, key)
	}

	err = p.DelValuesWithPrefix(keyPrefixs...)
	if err != nil {
		return err
	}

	return nil
}

func (p *Server) RegisterCmd(in *pbtypes.SubTask_RegisterCmd, out *pbtypes.Empty) error {
	etcdClient, err := p.etcd.GetClient(p.cfg.Get().GetConfdConfig().GetBackendConfig().GetHost(), time.Second)
	if err != nil {
		return err
	}
	return etcdClient.RegisterCmd(in)
}

func (p *EtcdClient) RegisterCmd(in *pbtypes.SubTask_RegisterCmd) error {
	var m = map[string]interface{}{}
	err := json.Unmarshal([]byte(in.Cnodes), &m)
	if err != nil {
		return err
	}

	err = p.SetValues(jsonmap.JsonMap(m).ToMapString("/"))
	if err != nil {
		return err
	}

	return nil
}

func (p *Server) DeregisterCmd(in *pbtypes.SubTask_DeregisterCmd, out *pbtypes.Empty) error {
	etcdClient, err := p.etcd.GetClient(p.cfg.Get().GetConfdConfig().GetBackendConfig().GetHost(), time.Second)
	if err != nil {
		return err
	}
	return etcdClient.DeregisterCmd(in)
}

func (p *EtcdClient) DeregisterCmd(in *pbtypes.SubTask_DeregisterCmd) error {
	var m = map[string]interface{}{}
	err := json.Unmarshal([]byte(in.Cnodes), &m)
	if err != nil {
		return err
	}

	var keyPrefixs []string
	for key, _ := range m {
		keyPrefixs = append(keyPrefixs, key)
	}

	err = p.DelValuesWithPrefix(keyPrefixs...)
	if err != nil {
		return err
	}

	return nil
}

func (p *Server) ReportSubTaskStatus(in *pbtypes.SubTaskStatus, out *pbtypes.Empty) error {
	ctx := context.Background()

	cfg := p.cfg.Get()
	client, conn, err := pilotutil.DialPilotService(ctx, cfg.PilotHost, int(cfg.PilotPort))
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
	etcdClient, err := p.etcd.GetClient(p.cfg.Get().GetConfdConfig().GetBackendConfig().GetHost(), time.Second)
	if err != nil {
		return err
	}

	m, err := etcdClient.GetValuesByPrefix(in.Value)
	if err != nil {
		return err
	}

	out.ValueMap = m
	return nil
}

func (p *Server) GetEtcdValues(in *pbtypes.StringList, out *pbtypes.StringMap) error {
	etcdClient, err := p.etcd.GetClient(p.cfg.Get().GetConfdConfig().GetBackendConfig().GetHost(), time.Second)
	if err != nil {
		return err
	}

	m, err := etcdClient.GetValues(in.ValueList...)
	if err != nil {
		return err
	}

	out.ValueMap = m
	return nil
}

func (p *Server) SetEtcdValues(in *pbtypes.StringMap, out *pbtypes.Empty) error {
	etcdClient, err := p.etcd.GetClient(p.cfg.Get().GetConfdConfig().GetBackendConfig().GetHost(), time.Second)
	if err != nil {
		return err
	}
	return etcdClient.SetValues(in.GetValueMap())
}

func (p *Server) PingPilot(in *pbtypes.Empty, out *pbtypes.Empty) error {
	ctx := context.Background()

	cfg := p.cfg.Get()
	client, conn, err := pilotutil.DialPilotService(ctx, cfg.PilotHost, int(cfg.PilotPort))
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
