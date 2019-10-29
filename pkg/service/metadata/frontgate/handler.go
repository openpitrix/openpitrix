// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package frontgate

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"reflect"
	"runtime"
	"time"

	"github.com/chai2010/jsonmap"

	"openpitrix.io/openpitrix/pkg/logger"
	pbfrontgate "openpitrix.io/openpitrix/pkg/pb/metadata/frontgate"
	pbtypes "openpitrix.io/openpitrix/pkg/pb/metadata/types"
	"openpitrix.io/openpitrix/pkg/service/metadata/drone/droneutil"
	"openpitrix.io/openpitrix/pkg/service/metadata/frontgate/frontgateutil"
	"openpitrix.io/openpitrix/pkg/service/metadata/pilot/pilotutil"
	"openpitrix.io/openpitrix/pkg/util/funcutil"
)

var (
	_ pbfrontgate.FrontgateService = (*Server)(nil)
)

func (p *Server) GetPilotVersion(*pbtypes.Empty, *pbtypes.Version) error {
	err := fmt.Errorf("TODO")
	logger.Warn(nil, "%+v", err)
	return err
}
func (p *Server) GetFrontgateVersion(*pbtypes.Empty, *pbtypes.Version) error {
	err := fmt.Errorf("TODO")
	logger.Warn(nil, "%+v", err)
	return err
}
func (p *Server) GetDroneVersion(*pbtypes.DroneEndpoint, *pbtypes.Version) error {
	err := fmt.Errorf("TODO")
	logger.Warn(nil, "%+v", err)
	return err
}

func (p *Server) PingMetadataBackend(*pbtypes.Empty, *pbtypes.Empty) error {

	etcdConfig := p.cfg.Get().GetEtcdConfig()
	etcdClient, err := p.etcd.GetClient(
		pkgGetEtcdEndpointsFromConfig(etcdConfig),
		time.Duration(etcdConfig.GetTimeoutSeconds())*time.Second,
		int(etcdConfig.GetMaxTxnOps()),
	)
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return err
	}

	if err := etcdClient.Ping(); err != nil {
		return err
	}

	return nil
}

func (p *EtcdClient) Ping() error {
	logger.Info(nil, funcutil.CallerName(1))

	if _, err := p.GetAllValues(); err != nil {
		logger.Warn(nil, "%+v", err)
		return err
	}

	return nil
}

func (p *Server) GetPilotConfig(in *pbtypes.Empty, out *pbtypes.PilotConfig) error {
	logger.Info(nil, funcutil.CallerName(1))

	ctx := context.Background()

	cfg := p.cfg.Get()
	client, conn, err := pilotutil.DialPilotServiceForFrontgate_TLS(
		ctx, cfg.PilotHost, int(cfg.PilotPort),
		p.tlsPilotConfig,
	)
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return err
	}
	defer conn.Close()

	info, err := client.GetPilotConfig(ctx, &pbtypes.Empty{})
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return err
	}

	*out = *info
	return nil
}

func (p *Server) GetFrontgateConfig(in *pbtypes.Empty, out *pbtypes.FrontgateConfig) error {
	logger.Info(nil, funcutil.CallerName(1))

	*out = *p.cfg.Get()
	return nil
}

func (p *Server) SetFrontgateConfig(cfg *pbtypes.FrontgateConfig, out *pbtypes.Empty) error {
	logger.Info(nil, funcutil.CallerName(1))

	if reflect.DeepEqual(cfg, p.cfg.Get()) {
		logger.Info(nil, "cfg is same, not changed")
		return nil
	}

	if len(cfg.GetNodeList()) == 0 {
		err := fmt.Errorf("frontgate.SetFrontgateConfig: node list is empty")
		logger.Warn(nil, "%+v", err)
		return err
	}

	var lastErr error
	for _, node := range cfg.GetNodeList() {
		func() {
			client, err := frontgateutil.DialFrontgateService(
				node.NodeIp, int(node.NodePort),
			)
			if err != nil {
				logger.Warn(nil, "%+v", err)
				return
			}
			defer client.Close()

			_, err = client.SetFrontgateNodeConfig(cfg)
			if err != nil {
				logger.Warn(nil, "%+v", err)
				lastErr = err
			}
		}()
	}
	if lastErr != nil {
		return lastErr
	}

	return nil
}

func (p *Server) SetFrontgateNodeConfig(cfg *pbtypes.FrontgateConfig, out *pbtypes.Empty) error {
	logger.Info(nil, funcutil.CallerName(1))

	if reflect.DeepEqual(cfg, p.cfg.Get()) {
		return nil
	}

	if err := p.cfg.Set(cfg); err != nil {
		logger.Warn(nil, "%+v", err)
		return err
	}

	if err := p.cfg.Save(); err != nil {
		logger.Warn(nil, "%+v", err)
		return err
	}

	return nil
}

func (p *Server) GetDroneList(in *pbtypes.Empty, out *pbtypes.DroneIdList) error {
	logger.Info(nil, funcutil.CallerName(1))

	err := fmt.Errorf("TODO")
	logger.Warn(nil, "%+v", err)
	return err
}

func (p *Server) GetDroneConfig(in *pbtypes.DroneEndpoint, out *pbtypes.DroneConfig) error {
	logger.Info(nil, funcutil.CallerName(1))

	ctx := context.Background()

	client, conn, err := droneutil.DialDroneService(ctx,
		in.GetDroneIp(),
		int(in.GetDronePort()),
	)
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return err
	}
	defer conn.Close()

	_, err = client.GetDroneConfig(ctx, &pbtypes.Empty{})
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return err
	}

	return nil
}

func (p *Server) SetDroneConfig(in *pbtypes.SetDroneConfigRequest, out *pbtypes.Empty) error {
	logger.Info(nil, funcutil.CallerName(1))

	ctx := context.Background()

	client, conn, err := droneutil.DialDroneService(ctx,
		in.Endpoint.GetDroneIp(),
		int(in.Endpoint.GetDronePort()),
	)
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return err
	}
	defer conn.Close()

	// 1. set drone config
	_, err = client.SetDroneConfig(ctx, in.GetConfig())
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return err
	}

	// 2. set confd config
	_, err = client.SetConfdConfig(ctx, p.cfg.Get().GetConfdConfig())
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return err
	}

	// 3. set frontgate config
	_, err = client.SetFrontgateConfig(ctx, p.cfg.Get())
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return err
	}

	return nil
}

func (p *Server) GetConfdConfig(in *pbtypes.ConfdEndpoint, out *pbtypes.ConfdConfig) error {
	logger.Info(nil, funcutil.CallerName(1))

	*out = *p.cfg.Get().GetConfdConfig()
	return nil
}

func (p *Server) IsConfdRunning(in *pbtypes.ConfdEndpoint, out *pbtypes.Bool) error {
	logger.Info(nil, funcutil.CallerName(1))

	ctx := context.Background()

	client, conn, err := droneutil.DialDroneService(ctx,
		in.GetDroneIp(),
		int(in.GetDronePort()),
	)
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return err
	}
	defer conn.Close()

	reply, err := client.IsConfdRunning(ctx, &pbtypes.Empty{})
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return err
	}

	out.Value = reply.GetValue()
	return nil
}

func (p *Server) StartConfd(in *pbtypes.ConfdEndpoint, out *pbtypes.Empty) error {
	logger.Info(nil, funcutil.CallerName(1))

	ctx := context.Background()

	client, conn, err := droneutil.DialDroneService(ctx,
		in.GetDroneIp(),
		int(in.GetDronePort()),
	)
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return err
	}
	defer conn.Close()

	_, err = client.SetConfdConfig(ctx, p.cfg.Get().GetConfdConfig())
	if err != nil {
		logger.Warn(nil, "%+v", err)
		// donot return
	}

	_, err = client.SetFrontgateConfig(ctx, p.cfg.Get())
	if err != nil {
		logger.Warn(nil, "%+v", err)
		// donot return
	}

	_, err = client.StartConfd(ctx, &pbtypes.Empty{})
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return err
	}

	return nil
}

func (p *Server) StopConfd(in *pbtypes.ConfdEndpoint, out *pbtypes.Empty) error {
	logger.Info(nil, funcutil.CallerName(1))

	ctx := context.Background()

	client, conn, err := droneutil.DialDroneService(ctx, in.GetDroneIp(), int(in.GetDronePort()))
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return err
	}
	defer conn.Close()

	_, err = client.StopConfd(ctx, &pbtypes.Empty{})
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return err
	}

	return nil
}

func (p *Server) RegisterMetadata(in *pbtypes.SubTask_RegisterMetadata, out *pbtypes.Empty) error {
	logger.Info(nil, funcutil.CallerName(1))

	etcdConfig := p.cfg.Get().GetEtcdConfig()
	etcdClient, err := p.etcd.GetClient(
		pkgGetEtcdEndpointsFromConfig(etcdConfig),
		time.Duration(etcdConfig.GetTimeoutSeconds())*time.Second,
		int(etcdConfig.GetMaxTxnOps()),
	)
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return err
	}

	err = etcdClient.RegisterMetadata(in)
	if err != nil {
		if etcdClient.isHaltErr(err) {
			p.etcd.ClearClient(etcdClient)
		}
		logger.Warn(nil, "%+v", err)
		return err
	}

	return nil
}

func (p *Server) RegisterMetadataMapping(in *pbtypes.SubTask_RegisterMetadata, out *pbtypes.Empty) error {
	logger.Info(nil, funcutil.CallerName(1))

	etcdConfig := p.cfg.Get().GetEtcdConfig()
	etcdClient, err := p.etcd.GetClient(
		pkgGetEtcdEndpointsFromConfig(etcdConfig),
		time.Duration(etcdConfig.GetTimeoutSeconds())*time.Second,
		int(etcdConfig.GetMaxTxnOps()),
	)
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return err
	}

	err = etcdClient.RegisterMetadataMapping(in)
	if err != nil {
		if etcdClient.isHaltErr(err) {
			p.etcd.ClearClient(etcdClient)
		}
		logger.Warn(nil, "%+v", err)
		return err
	}

	return nil
}

func (p *EtcdClient) RegisterMetadata(in *pbtypes.SubTask_RegisterMetadata) error {
	logger.Info(nil, funcutil.CallerName(1))

	var m = map[string]interface{}{}
	err := json.Unmarshal([]byte(in.Cnodes), &m)
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return err
	}

	err = p.SetValues(jsonmap.JsonMap(m).ToMapString("/"))
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return err
	}

	return nil
}

func (p *EtcdClient) RegisterMetadataMapping(in *pbtypes.SubTask_RegisterMetadata) error {
	logger.Info(nil, funcutil.CallerName(1))

	var m = map[string]interface{}{}
	err := json.Unmarshal([]byte(in.Cnodes), &m)
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return err
	}

	const metadMappingPrefix = "/_metad/mapping/default"
	err = p.SetValuesWithPrefix(metadMappingPrefix, jsonmap.JsonMap(m).ToMapString("/"))
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return err
	}

	return nil
}

func (p *Server) DeregisterMetadata(in *pbtypes.SubTask_DeregisterMetadata, out *pbtypes.Empty) error {
	logger.Info(nil, funcutil.CallerName(1))

	etcdConfig := p.cfg.Get().GetEtcdConfig()
	etcdClient, err := p.etcd.GetClient(
		pkgGetEtcdEndpointsFromConfig(etcdConfig),
		time.Duration(etcdConfig.GetTimeoutSeconds())*time.Second,
		int(etcdConfig.GetMaxTxnOps()),
	)
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return err
	}

	err = etcdClient.DeregisterMetadata(in)
	if err != nil {
		if etcdClient.isHaltErr(err) {
			p.etcd.ClearClient(etcdClient)
		}
		logger.Warn(nil, "%+v", err)
		return err
	}

	return nil
}

func (p *Server) DeregisterMetadataMapping(in *pbtypes.SubTask_DeregisterMetadata, out *pbtypes.Empty) error {
	logger.Info(nil, funcutil.CallerName(1))

	etcdConfig := p.cfg.Get().GetEtcdConfig()
	etcdClient, err := p.etcd.GetClient(
		pkgGetEtcdEndpointsFromConfig(etcdConfig),
		time.Duration(etcdConfig.GetTimeoutSeconds())*time.Second,
		int(etcdConfig.GetMaxTxnOps()),
	)
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return err
	}

	err = etcdClient.DeregisterMetadataMapping(in)
	if err != nil {
		if etcdClient.isHaltErr(err) {
			p.etcd.ClearClient(etcdClient)
		}
		logger.Warn(nil, "%+v", err)
		return err
	}

	return nil
}

func (p *EtcdClient) DeregisterMetadata(in *pbtypes.SubTask_DeregisterMetadata) error {
	logger.Info(nil, funcutil.CallerName(1))

	var m = map[string]interface{}{}
	err := json.Unmarshal([]byte(in.Cnodes), &m)
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return err
	}

	var keyPrefixs []string
	for key := range jsonmap.JsonMap(m).ToMapString("/") {
		keyPrefixs = append(keyPrefixs, key)
	}

	err = p.DelValuesWithPrefix(keyPrefixs...)
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return err
	}

	return nil
}

func (p *EtcdClient) DeregisterMetadataMapping(in *pbtypes.SubTask_DeregisterMetadata) error {
	logger.Info(nil, funcutil.CallerName(1))

	var m = map[string]interface{}{}
	err := json.Unmarshal([]byte(in.Cnodes), &m)
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return err
	}

	var keyPrefixs []string
	const metadMappingPrefix = "/_metad/mapping/default"

	for key := range jsonmap.JsonMap(m).ToMapString("/") {
		keyPrefixs = append(keyPrefixs, metadMappingPrefix+key)
	}

	err = p.DelValuesWithPrefix(keyPrefixs...)
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return err
	}

	return nil
}

func (p *Server) RegisterCmd(in *pbtypes.SubTask_RegisterCmd, out *pbtypes.Empty) error {
	logger.Info(nil, funcutil.CallerName(1))

	etcdConfig := p.cfg.Get().GetEtcdConfig()
	etcdClient, err := p.etcd.GetClient(
		pkgGetEtcdEndpointsFromConfig(etcdConfig),
		time.Duration(etcdConfig.GetTimeoutSeconds())*time.Second,
		int(etcdConfig.GetMaxTxnOps()),
	)
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return err
	}

	err = etcdClient.RegisterCmd(in)
	if err != nil {
		if etcdClient.isHaltErr(err) {
			p.etcd.ClearClient(etcdClient)
		}
		logger.Warn(nil, "%+v", err)
		return err
	}

	return nil
}

func (p *EtcdClient) RegisterCmd(in *pbtypes.SubTask_RegisterCmd) error {
	logger.Info(nil, funcutil.CallerName(1))

	var m = map[string]interface{}{}
	err := json.Unmarshal([]byte(in.Cnodes), &m)
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return err
	}

	err = p.SetValues(jsonmap.JsonMap(m).ToMapString("/"))
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return err
	}

	return nil
}

func (p *Server) DeregisterCmd(in *pbtypes.SubTask_DeregisterCmd, out *pbtypes.Empty) error {
	logger.Info(nil, funcutil.CallerName(1))

	etcdConfig := p.cfg.Get().GetEtcdConfig()
	etcdClient, err := p.etcd.GetClient(
		pkgGetEtcdEndpointsFromConfig(etcdConfig),
		time.Duration(etcdConfig.GetTimeoutSeconds())*time.Second,
		int(etcdConfig.GetMaxTxnOps()),
	)
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return err
	}

	err = etcdClient.DeregisterCmd(in)
	if err != nil {
		if etcdClient.isHaltErr(err) {
			p.etcd.ClearClient(etcdClient)
		}
		logger.Warn(nil, "%+v", err)
		return err
	}

	return nil
}

func (p *EtcdClient) DeregisterCmd(in *pbtypes.SubTask_DeregisterCmd) error {
	logger.Info(nil, funcutil.CallerName(1))

	var m = map[string]interface{}{}
	err := json.Unmarshal([]byte(in.Cnodes), &m)
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return err
	}

	var keyPrefixs []string
	for key := range jsonmap.JsonMap(m).ToMapString("/") {
		keyPrefixs = append(keyPrefixs, key)
	}

	err = p.DelValuesWithPrefix(keyPrefixs...)
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return err
	}

	return nil
}

func (p *Server) ReportSubTaskStatus(in *pbtypes.SubTaskStatus, out *pbtypes.Empty) error {
	logger.Info(nil, "%s taskId: %s", funcutil.CallerName(1), in.TaskId)

	ctx := context.Background()

	cfg := p.cfg.Get()
	client, conn, err := pilotutil.DialPilotServiceForFrontgate_TLS(
		ctx, cfg.PilotHost, int(cfg.PilotPort),
		p.tlsPilotConfig,
	)
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return err
	}
	defer conn.Close()

	_, err = client.ReportSubTaskStatus(ctx, in)
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return err
	}

	return nil
}

func (p *Server) ClosePilotChannel(in *pbtypes.Empty, out *pbtypes.Empty) error {
	logger.Info(nil, funcutil.CallerName(1))

	err := p.conn.Close()
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return err
	}

	return nil
}

func (p *Server) GetEtcdValuesByPrefix(in *pbtypes.String, out *pbtypes.StringMap) error {
	logger.Info(nil, funcutil.CallerName(1))

	etcdConfig := p.cfg.Get().GetEtcdConfig()
	etcdClient, err := p.etcd.GetClient(
		pkgGetEtcdEndpointsFromConfig(etcdConfig),
		time.Duration(etcdConfig.GetTimeoutSeconds())*time.Second,
		int(etcdConfig.GetMaxTxnOps()),
	)
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return err
	}

	m, err := etcdClient.GetValuesByPrefix(in.Value)
	if err != nil {
		if etcdClient.isHaltErr(err) {
			p.etcd.ClearClient(etcdClient)
		}
		logger.Warn(nil, "%+v", err)
		return err
	}

	out.ValueMap = m
	return nil
}

func (p *Server) GetEtcdValues(in *pbtypes.StringList, out *pbtypes.StringMap) error {
	logger.Info(nil, funcutil.CallerName(1))

	etcdConfig := p.cfg.Get().GetEtcdConfig()
	etcdClient, err := p.etcd.GetClient(
		pkgGetEtcdEndpointsFromConfig(etcdConfig),
		time.Duration(etcdConfig.GetTimeoutSeconds())*time.Second,
		int(etcdConfig.GetMaxTxnOps()),
	)
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return err
	}

	m, err := etcdClient.GetValues(in.ValueList...)
	if err != nil {
		if etcdClient.isHaltErr(err) {
			p.etcd.ClearClient(etcdClient)
		}
		logger.Warn(nil, "%+v", err)
		return err
	}

	out.ValueMap = m
	return nil
}

func (p *Server) SetEtcdValues(in *pbtypes.StringMap, out *pbtypes.Empty) error {
	logger.Info(nil, funcutil.CallerName(1))

	etcdConfig := p.cfg.Get().GetEtcdConfig()
	etcdClient, err := p.etcd.GetClient(
		pkgGetEtcdEndpointsFromConfig(etcdConfig),
		time.Duration(etcdConfig.GetTimeoutSeconds())*time.Second,
		int(etcdConfig.GetMaxTxnOps()),
	)
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return err
	}

	err = etcdClient.SetValues(in.GetValueMap())
	if err != nil {
		if etcdClient.isHaltErr(err) {
			p.etcd.ClearClient(etcdClient)
		}
		logger.Warn(nil, "%+v", err)
		return err
	}

	return nil
}

func (p *Server) PingPilot(in *pbtypes.Empty, out *pbtypes.Empty) error {
	logger.Info(nil, funcutil.CallerName(1))

	ctx := context.Background()

	cfg := p.cfg.Get()
	client, conn, err := pilotutil.DialPilotServiceForFrontgate_TLS(
		ctx, cfg.PilotHost, int(cfg.PilotPort),
		p.tlsPilotConfig,
	)
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return err
	}
	defer conn.Close()

	_, err = client.PingPilot(ctx, &pbtypes.Empty{}) // todo
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return err
	}

	return nil
}

func (p *Server) PingFrontgate(in *pbtypes.Empty, out *pbtypes.Empty) error {
	logger.Info(nil, funcutil.CallerName(1))

	cfg := p.cfg.Get()
	if len(cfg.GetNodeList()) == 0 {
		logger.Info(nil, "PingFrontgate: ok")
		return nil // OK
	}

	var lastErr error
	for _, node := range cfg.GetNodeList() {
		func() {
			client, err := frontgateutil.DialFrontgateService(
				node.NodeIp, int(node.NodePort),
			)
			if err != nil {
				logger.Warn(nil, "%+v", err)
				return
			}
			defer client.Close()

			_, err = client.PingFrontgateNode(&pbtypes.Empty{})
			if err != nil {
				logger.Warn(nil, "%+v", err)
				lastErr = err
			}
		}()
	}
	if lastErr != nil {
		return lastErr
	}

	logger.Info(nil, "PingFrontgate: all nodes ok")
	return nil // OK
}

func (p *Server) PingFrontgateNode(in *pbtypes.Empty, out *pbtypes.Empty) error {
	logger.Info(nil, funcutil.CallerName(1))
	return nil // OK
}

func (p *Server) PingDrone(in *pbtypes.DroneEndpoint, out *pbtypes.Empty) error {
	logger.Info(nil, funcutil.CallerName(1))

	ctx := context.Background()

	client, conn, err := droneutil.DialDroneService(ctx,
		in.GetDroneIp(),
		int(in.GetDronePort()),
	)
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return err
	}
	defer conn.Close()

	_, err = client.PingDrone(ctx, &pbtypes.Empty{})
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return err
	}

	return nil
}

func (p *Server) RunCommand(arg *pbtypes.RunCommandOnFrontgateRequest, out *pbtypes.String) error {
	logger.Info(nil, funcutil.CallerName(1))

	var c *exec.Cmd
	if runtime.GOOS == "windows" {
		c = exec.Command("cmd", "/C", arg.GetCommand())
	} else {
		c = exec.Command("/bin/sh", "-c", arg.GetCommand())
	}

	var b bytes.Buffer
	c.Stdout = &b
	c.Stderr = &b

	if err := c.Start(); err != nil {
		logger.Warn(nil, "%+v", err)
		return err
	}

	var timeout = time.Second * 3
	if x := arg.GetTimeoutSeconds(); x > 0 {
		timeout = time.Duration(x) * time.Second
	}

	timer := time.AfterFunc(timeout, func() {
		c.Process.Kill()
	})

	err := c.Wait()
	timer.Stop()

	if err != nil {
		logger.Warn(nil, "%+v", err)
		return err
	}

	out.Value = string(b.Bytes())
	return nil // OK
}
func (p *Server) RunCommandOnDrone(in *pbtypes.RunCommandOnDroneRequest, out *pbtypes.String) error {
	logger.Info(nil, funcutil.CallerName(1))

	ctx := context.Background()

	client, conn, err := droneutil.DialDroneService(ctx,
		in.GetEndpoint().GetDroneIp(),
		int(in.GetEndpoint().GetDronePort()),
	)
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return err
	}
	defer conn.Close()

	_, err = client.RunCommand(ctx, in)
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return err
	}

	return nil
}

func (p *Server) HeartBeat(in *pbtypes.Empty, out *pbtypes.Empty) error {
	return nil // OK
}
