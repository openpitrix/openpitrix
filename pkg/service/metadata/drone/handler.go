// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package drone

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"time"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/libconfd"
	"openpitrix.io/openpitrix/pkg/logger"
	pbdrone "openpitrix.io/openpitrix/pkg/pb/metadata/drone"
	pbtypes "openpitrix.io/openpitrix/pkg/pb/metadata/types"
	"openpitrix.io/openpitrix/pkg/service/metadata/drone/yunify_confdfunc"
	"openpitrix.io/openpitrix/pkg/util/funcutil"
	"openpitrix.io/openpitrix/pkg/util/httputil"
	"openpitrix.io/openpitrix/pkg/util/retryutil"
	"openpitrix.io/openpitrix/pkg/version"
)

var (
	_                      pbdrone.DroneServiceServer = (*Server)(nil)
	RetryInterval                                     = 3 * time.Second
	RetryCount                                        = 5
	DowloadPathPattern                                = "/opt/openpitrix/bin/%s"
	DowloadFilePathPattern                            = "/opt/openpitrix/bin/%s/%s"
	PilotVersionFilePath                              = "/opt/openpitrix/conf/pilot-version"
)

func (p *Server) getDroneFromFrontgate(pilotVersion, frontgateAddress string) error {
	var droneBinary []byte

	url := fmt.Sprintf("http://%s:%d/%s/drone", frontgateAddress, constants.FrontgateFileServerPort, pilotVersion)
	logger.Info(nil, "Trying to download drone from url [%s]", url)

	err := retryutil.Retry(RetryCount, RetryInterval, func() error {
		resp, err := httputil.HttpGet(url)
		if err != nil {
			return err
		}

		if resp.StatusCode != 200 {
			return fmt.Errorf("download drone from url [%s] failed, status %s", url, resp.Status)
		}

		droneBinary, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	err = os.MkdirAll(fmt.Sprintf(DowloadPathPattern, pilotVersion), os.ModeDir|os.ModePerm)
	if err != nil {
		return err
	}

	filePath := fmt.Sprintf(DowloadFilePathPattern, pilotVersion, "drone")

	logger.Info(nil, "Write drone to [%s]", filePath)
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(droneBinary)
	if err != nil {
		return err
	}

	err = os.Chmod(filePath, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func (p *Server) createPilotVersionFile(pilotVersion string) error {
	f, err := os.Create(PilotVersionFilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(pilotVersion)
	if err != nil {
		return err
	}

	return nil
}

func (p *Server) DistributeDrone(ctx context.Context, req *pbtypes.DistributeDroneRequest) (*pbtypes.Empty, error) {
	pilotVersion := req.GetPilotVersion().GetValue()
	frontgateAddress := req.GetFrontgateAddress().GetValue()

	err := p.getDroneFromFrontgate(pilotVersion, frontgateAddress)
	if err != nil {
		return &pbtypes.Empty{}, err
	}

	err = p.createPilotVersionFile(pilotVersion)
	if err != nil {
		return &pbtypes.Empty{}, err
	}

	logger.Info(nil, "Drone exit")
	os.Exit(0)

	return &pbtypes.Empty{}, nil
}

func (p *Server) GetPilotVersion(context.Context, *pbtypes.Empty) (*pbtypes.Version, error) {
	err := fmt.Errorf("TODO")
	logger.Warn(nil, "%+v", err)
	return nil, err
}
func (p *Server) GetFrontgateVersion(context.Context, *pbtypes.Empty) (*pbtypes.Version, error) {
	err := fmt.Errorf("TODO")
	logger.Warn(nil, "%+v", err)
	return nil, err
}
func (p *Server) GetDroneVersion(context.Context, *pbtypes.DroneEndpoint) (*pbtypes.Version, error) {
	reply := &pbtypes.Version{
		ShortVersion:   version.ShortVersion,
		GitSha1Version: version.GitSha1Version,
		BuildDate:      version.BuildDate,
	}
	return reply, nil
}

func (p *Server) PingMetadataBackend(context.Context, *pbtypes.FrontgateEndpoint) (*pbtypes.Empty, error) {
	err := fmt.Errorf("TODO")
	logger.Warn(nil, "%+v", err)
	return nil, err
}

func (p *Server) GetDroneConfig(context.Context, *pbtypes.Empty) (*pbtypes.DroneConfig, error) {
	logger.Info(nil, funcutil.CallerName(1))

	return p.cfg.Get(), nil
}
func (p *Server) SetDroneConfig(ctx context.Context, cfg *pbtypes.DroneConfig) (*pbtypes.Empty, error) {
	logger.Info(nil, funcutil.CallerName(1))

	if reflect.DeepEqual(cfg, p.cfg.Get()) {
		return &pbtypes.Empty{}, nil
	}

	if err := p.cfg.Set(cfg); err != nil {
		logger.Warn(nil, "%+v", err)
		return &pbtypes.Empty{}, err
	}

	if err := p.cfg.Save(); err != nil {
		logger.Warn(nil, "%+v", err)
		return &pbtypes.Empty{}, err
	}

	return &pbtypes.Empty{}, nil
}

func (p *Server) GetConfdConfig(ctx context.Context, arg *pbtypes.Empty) (*pbtypes.ConfdConfig, error) {
	logger.Info(nil, funcutil.CallerName(1))

	return p.confd.GetConfig(), nil
}

func (p *Server) SetConfdConfig(ctx context.Context, arg *pbtypes.ConfdConfig) (*pbtypes.Empty, error) {
	logger.Info(nil, funcutil.CallerName(1))

	cfg := p.cfg.Get()
	fnHookKeyAdjuster := func(absKey string) (realKey string) {
		if absKey == "/self" {
			return "/" + cfg.ConfdSelfHost
		}
		if strings.HasPrefix(absKey, "/self/") {
			return "/" + cfg.ConfdSelfHost + absKey[len("/self/")-1:]
		}
		return absKey
	}

	if err := p.confd.SetConfig(arg, fnHookKeyAdjuster); err != nil {
		logger.Warn(nil, "%+v", err)
		return nil, err
	}
	if err := p.confd.Save(); err != nil {
		logger.Warn(nil, "%+v", err)
		return nil, err
	}

	return &pbtypes.Empty{}, nil
}

func (p *Server) GetFrontgateConfig(context.Context, *pbtypes.Empty) (*pbtypes.FrontgateConfig, error) {
	logger.Info(nil, funcutil.CallerName(1))

	cfg, err := p.fg.GetConfig()
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return nil, err
	}
	return cfg, nil
}

func (p *Server) SetFrontgateConfig(ctx context.Context, cfg *pbtypes.FrontgateConfig) (*pbtypes.Empty, error) {
	logger.Info(nil, funcutil.CallerName(1))

	err := p.fg.SetConfig(cfg)
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return nil, err
	}
	return &pbtypes.Empty{}, nil
}

func (p *Server) IsConfdRunning(ctx context.Context, arg *pbtypes.Empty) (*pbtypes.Bool, error) {
	logger.Info(nil, funcutil.CallerName(1))

	return &pbtypes.Bool{Value: p.confd.IsRunning()}, nil
}

func (p *Server) StartConfd(ctx context.Context, arg *pbtypes.Empty) (*pbtypes.Empty, error) {
	logger.Info(nil, funcutil.CallerName(1))

	cfg := p.cfg.Get()
	logger.Info(nil, "StartConfd: %v", cfg)

	err := p.confd.Start(func(opt *libconfd.Config) {
		opt.FuncMap = yunify_confdfunc.MakeCustomFuncMap()
		opt.HookAbsKeyAdjuster = func(absKey string) (realKey string) {
			if absKey == "/self" {
				return absKey
			}
			if !strings.HasPrefix(absKey, "/self/") {
				return "/self" + absKey
			}
			return absKey
		}
		opt.HookOnCheckCmdDone = func(trName, cmd string, err error) {
			if err != nil {
				logger.Warn(nil, "%+v", err)
				return
			}
			if trName == "/etc/confd/conf.d/cmd.info.toml" {
				go func() {
					logger.Info(nil, "LoadLastCmdStatus: %v", cfg.CmdInfoLogPath)

					status, isEmpty, err := LoadLastCmdStatus(cfg.CmdInfoLogPath)
					if isEmpty {
						return
					}
					if err != nil {
						logger.Warn(nil, "%+v", err)
						p.fg.ReportSubTaskStatus(&pbtypes.SubTaskStatus{
							TaskId: status.SubtaskId,
							Status: constants.StatusFailed,
						})
					} else {
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
				logger.Warn(nil, "%+v", err)
				return
			}
			if trName == "/etc/confd/conf.d/cmd.info.toml" {
				go func() {
					status, isEmpty, err := LoadLastCmdStatus(cfg.CmdInfoLogPath)
					if isEmpty {
						return
					}
					if err != nil {
						logger.Warn(nil, "%+v", err)
						p.fg.ReportSubTaskStatus(&pbtypes.SubTaskStatus{
							TaskId: status.SubtaskId,
							Status: constants.StatusFailed,
						})
					} else {
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
				logger.Warn(nil, "%+v", err)
				return
			}
		}
	})

	if err != nil {
		logger.Warn(nil, "%+v", err)
		return nil, err
	}

	return &pbtypes.Empty{}, nil
}

func (p *Server) StopConfd(ctx context.Context, arg *pbtypes.Empty) (*pbtypes.Empty, error) {
	logger.Info(nil, funcutil.CallerName(1))

	if err := p.confd.Stop(); err != nil {
		logger.Error(nil, "StopConfd: %v", err)
		return nil, err
	}

	return &pbtypes.Empty{}, nil
}

func (p *Server) GetTemplateFiles(ctx context.Context, arg *pbtypes.Empty) (*pbtypes.StringList, error) {
	logger.Info(nil, funcutil.CallerName(1))

	if !p.confd.IsRunning() {
		err := fmt.Errorf("drone: confd is not started")
		logger.Warn(nil, "%+v", err)
		return nil, err
	}

	cfg := p.confd.GetConfig()
	confdir := cfg.GetProcessorConfig().GetConfdir()
	if confdir == "" {
		err := fmt.Errorf("drone: invaid confdir: %q", confdir)
		logger.Warn(nil, "%+v", err)
		return nil, err
	}

	_, paths, err := libconfd.ListTemplateResource(filepath.Join(confdir, "conf.d"))
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return nil, err
	}

	reply := &pbtypes.StringList{}
	for _, s := range paths {
		reply.ValueList = append(reply.ValueList, filepath.Base(s))
	}

	return reply, nil
}

func (p *Server) GetValues(ctx context.Context, arg *pbtypes.StringList) (*pbtypes.StringMap, error) {
	logger.Info(nil, funcutil.CallerName(1))

	if !p.confd.IsRunning() {
		err := fmt.Errorf("drone: confd is not started")
		logger.Warn(nil, "%+v", err)
		return nil, err
	}

	client := p.confd.GetBackendClient()
	if client == nil {
		return nil, fmt.Errorf("drone: confd is not started")
	}

	keys := arg.GetValueList()
	m, err := client.GetValues(keys)
	if err != nil {
		logger.Warn(nil, "%+v", err)
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
	logger.Info(nil, funcutil.CallerName(1))

	c, err := p.fg.getClient()
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return nil, err
	}

	_, err = c.PingPilot(&pbtypes.Empty{})
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return nil, err
	}

	return &pbtypes.Empty{}, nil
}

func (p *Server) PingFrontgate(ctx context.Context, arg *pbtypes.FrontgateEndpoint) (*pbtypes.Empty, error) {
	logger.Info(nil, funcutil.CallerName(1))

	c, err := p.fg.getClient()
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return nil, err
	}

	_, err = c.PingFrontgate(&pbtypes.Empty{})
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return nil, err
	}

	return &pbtypes.Empty{}, nil
}

func (p *Server) PingDrone(ctx context.Context, arg *pbtypes.Empty) (*pbtypes.Empty, error) {
	logger.Info(nil, funcutil.CallerName(1))
	return &pbtypes.Empty{}, nil
}

func (p *Server) RunCommand(ctx context.Context, arg *pbtypes.RunCommandOnDroneRequest) (*pbtypes.String, error) {
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
		return nil, err
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
		return nil, err
	}

	return &pbtypes.String{Value: string(b.Bytes())}, nil
}
