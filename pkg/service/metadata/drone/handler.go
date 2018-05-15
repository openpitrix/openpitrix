// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package drone

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"

	"openpitrix.io/openpitrix/pkg/libconfd"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/pb/drone"
	"openpitrix.io/openpitrix/pkg/pb/types"
)

var (
	_ pbdrone.DroneServiceServer = (*Server)(nil)
)

func (p *Server) GetDroneConfig(context.Context, *pbtypes.Empty) (*pbtypes.DroneConfig, error) {
	return p.cfg.Get(), nil
}
func (p *Server) SetDroneConfig(ctx context.Context, cfg *pbtypes.DroneConfig) (*pbtypes.Empty, error) {
	if reflect.DeepEqual(cfg, p.cfg.Get()) {
		return &pbtypes.Empty{}, nil
	}

	if err := p.cfg.Set(cfg); err != nil {
		return &pbtypes.Empty{}, err
	}

	if err := p.cfg.Save(); err != nil {
		return &pbtypes.Empty{}, err
	}

	return &pbtypes.Empty{}, nil
}

func (p *Server) GetConfdConfig(ctx context.Context, arg *pbtypes.Empty) (*pbtypes.ConfdConfig, error) {
	return p.confd.GetConfig(), nil
}

func (p *Server) SetConfdConfig(ctx context.Context, arg *pbtypes.ConfdConfig) (*pbtypes.Empty, error) {
	if err := p.confd.SetConfig(arg); err != nil {
		return nil, err
	}
	return &pbtypes.Empty{}, nil
}

func (p *Server) GetFrontgateConfig(context.Context, *pbtypes.Empty) (*pbtypes.FrontgateConfig, error) {
	cfg, err := p.fg.GetConfig()
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func (p *Server) SetFrontgateConfig(ctx context.Context, cfg *pbtypes.FrontgateConfig) (*pbtypes.Empty, error) {
	err := p.fg.SetConfig(cfg)
	if err != nil {
		return nil, err
	}
	return &pbtypes.Empty{}, nil
}

func (p *Server) IsConfdRunning(ctx context.Context, arg *pbtypes.Empty) (*pbtypes.Bool, error) {
	return &pbtypes.Bool{Value: p.confd.IsRunning()}, nil
}

func (p *Server) StartConfd(ctx context.Context, arg *pbtypes.Empty) (*pbtypes.Empty, error) {
	err := p.confd.Start(func(opt *libconfd.Config) {
		cfg := p.cfg.Get()
		opt.HookAbsKeyAdjuster = func(absKey string) (realKey string) {
			if absKey == "/self" {
				return "/" + cfg.ConfdSelfHost
			}
			if strings.HasPrefix(absKey, "/self/") {
				return "/" + cfg.ConfdSelfHost + absKey[len("/self/")-1:]
			}
			return absKey
		}
		opt.HookOnCheckCmdDone = func(trName, cmd string, err error) {
			if err != nil {
				logger.Warn("%+v", err)
				return
			}
			if trName == "cmd.info" {
				go func() {
					status, err := LoadLastCmdStatus(cfg.CmdInfoLogPath)
					if err == nil {
						p.fg.ReportSubTaskStatus(&pbtypes.SubTaskStatus{
							TaskId: status.SubtaskId,
							Status: status.Status,
						})
					}
				}()
			}
		}
		opt.HookOnReloadCmdDone = func(trName, cmd string, err error) {
			if err != nil {
				logger.Warn("%+v", err)
				return
			}
			if trName == "cmd.info" {
				go func() {
					status, err := LoadLastCmdStatus(cfg.CmdInfoLogPath)
					if err == nil {
						p.fg.ReportSubTaskStatus(&pbtypes.SubTaskStatus{
							TaskId: status.SubtaskId,
							Status: status.Status,
						})
					}
				}()
			}
		}
		opt.HookOnUpdateDone = func(trName string, err error) {
			if err != nil {
				logger.Warn("%+v", err)
				return
			}
		}
	})

	if err != nil {
		return nil, err
	}

	return &pbtypes.Empty{}, nil
}

func (p *Server) StopConfd(ctx context.Context, arg *pbtypes.Empty) (*pbtypes.Empty, error) {
	if err := p.confd.Stop(); err != nil {
		return nil, err
	}

	return &pbtypes.Empty{}, nil
}

func (p *Server) GetTemplateFiles(ctx context.Context, arg *pbtypes.Empty) (*pbtypes.StringList, error) {
	if !p.confd.IsRunning() {
		return nil, fmt.Errorf("drone: confd is not started")
	}

	cfg := p.confd.GetConfig()
	confdir := cfg.GetProcessorConfig().GetConfdir()
	if confdir == "" {
		return nil, fmt.Errorf("drone: invaid confdir: %q", confdir)
	}

	_, paths, err := libconfd.ListTemplateResource(filepath.Join(confdir, "conf.d"))
	if err != nil {
		return nil, err
	}

	reply := &pbtypes.StringList{}
	for _, s := range paths {
		reply.ValueList = append(reply.ValueList, filepath.Base(s))
	}

	return reply, nil
}

func (p *Server) GetValues(ctx context.Context, arg *pbtypes.StringList) (*pbtypes.StringMap, error) {
	if !p.confd.IsRunning() {
		return nil, fmt.Errorf("drone: confd is not started")
	}

	client := p.confd.GetBackendClient()
	if client == nil {
		return nil, fmt.Errorf("drone: confd is not started")
	}

	keys := arg.GetValueList()
	m, err := client.GetValues(keys)
	if err != nil {
		return nil, err
	}

	reply := &pbtypes.StringMap{
		ValueMap: make(map[string]string),
	}

	for _, k := range keys {
		reply.ValueMap[k] = m[k]
	}

	return reply, nil
}

func (p *Server) PingPilot(ctx context.Context, arg *pbtypes.FrontgateEndpoint) (*pbtypes.Empty, error) {
	c, err := p.fg.getClient()
	if err != nil {
		return nil, err
	}

	_, err = c.PingPilot(&pbtypes.Empty{})
	if err != nil {
		return nil, err
	}

	return &pbtypes.Empty{}, nil
}

func (p *Server) PingFrontgate(ctx context.Context, arg *pbtypes.FrontgateEndpoint) (*pbtypes.Empty, error) {
	c, err := p.fg.getClient()
	if err != nil {
		return nil, err
	}

	_, err = c.PingFrontgate(&pbtypes.Empty{})
	if err != nil {
		return nil, err
	}

	return &pbtypes.Empty{}, nil
}

func (p *Server) PingDrone(ctx context.Context, arg *pbtypes.Empty) (*pbtypes.Empty, error) {
	logger.Info("PingDrone: ok")
	return &pbtypes.Empty{}, nil
}

func (p *Server) RunCommand(ctx context.Context, arg *pbtypes.String) (*pbtypes.String, error) {
	var c *exec.Cmd
	if runtime.GOOS == "windows" {
		c = exec.Command("cmd", "/C", arg.Value)
	} else {
		c = exec.Command("/bin/sh", "-c", arg.Value)
	}

	output, err := c.CombinedOutput()
	if err != nil {
		return nil, err
	}

	return &pbtypes.String{Value: string(output)}, nil
}
