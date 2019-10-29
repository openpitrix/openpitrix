// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package drone

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"

	"openpitrix.io/openpitrix/pkg/libconfd"
	"openpitrix.io/openpitrix/pkg/libconfd/backends"
	"openpitrix.io/openpitrix/pkg/logger"
	pbtypes "openpitrix.io/openpitrix/pkg/pb/metadata/types"
)

type ConfdServer struct {
	cfgpath string

	mu            sync.Mutex
	cfg           *pbtypes.ConfdConfig
	config        *libconfd.Config
	backendConfig *libconfd.BackendConfig

	processor *libconfd.Processor
	client    libconfd.BackendClient
	running   bool
	err       error
}

func NewConfdServer(cfgpath string) *ConfdServer {
	if !filepath.IsAbs(cfgpath) {
		logger.Error(nil, "NewConfdServer: cfgpath is not abs path: %s", cfgpath)
	}

	return &ConfdServer{
		cfgpath: cfgpath,
	}
}

func (p *ConfdServer) IsRunning() bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.processor != nil && p.running
}

func (p *ConfdServer) GetConfig() *pbtypes.ConfdConfig {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.cfg != nil {
		return p.cfg
	}

	return &pbtypes.ConfdConfig{}
}

func (p *ConfdServer) SetConfig(cfg *pbtypes.ConfdConfig, fnHookKeyAdjuster func(key string) (realKey string)) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	config, backendConfig, err := p.parseConfig(cfg)
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return err
	}

	p.cfg = proto.Clone(cfg).(*pbtypes.ConfdConfig)
	p.config = config
	p.backendConfig = backendConfig
	p.backendConfig.HookKeyAdjuster = fnHookKeyAdjuster

	return nil
}

func (p *ConfdServer) GetBackendClient() libconfd.BackendClient {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.client
}

func (p *ConfdServer) Start(opts ...libconfd.Options) error {
	logger.Info(nil, "ConfdServer: Start")

	p.mu.Lock()
	defer p.mu.Unlock()

	if p.processor == nil {
		p.processor = libconfd.NewProcessor()
	}
	if p.running {
		logger.Error(nil, "ConfdServer: confd is running")
		return fmt.Errorf("drone: confd is running")
	}

	if p.cfg == nil {
		logger.Error(nil, "ConfdServer: config is nil")
		return fmt.Errorf("drone: config is nil")
	}

	switch p.backendConfig.Type {
	case backends.Etcdv3BackendType:
		// etcd: OK
	case backends.MetadBackendType:
		// metad: OK
	default:
		s := p.backendConfig.Type
		logger.Error(nil, "ConfdServer: unsupport confd backend: "+s)
		return fmt.Errorf("drone: unsupport confd backend: %s", s)
	}

	// apply opts
	for _, fn := range opts {
		fn(p.config)
	}

	logger.Info(nil, "Confd backend is [%s]", p.backendConfig.Type)
	backendClient, err := libconfd.NewBackendClient(p.backendConfig)
	if err != nil {
		logger.Error(nil, "ConfdServer: NewBackendClient: %v", err)
		return err
	}

	p.client = backendClient
	p.running = true

	go func() {
		logger.Info(nil, "ConfdServer: run...")

		// set log level
		libconfd.GetLogger().SetLevel(p.config.LogLevel)

		var err = p.processor.Run(p.config, backendClient) // blocked
		if err != nil {
			logger.Warn(nil, "%+v", err)
		}

		p.mu.Lock()
		p.running = false
		p.client = nil
		p.err = err
		p.mu.Unlock()

		logger.Info(nil, "ConfdServer: stoped")
	}()

	return nil
}

func (p *ConfdServer) Stop() error {
	logger.Info(nil, "ConfdServer: Stop")

	p.mu.Lock()
	var processer = p.processor
	p.processor = nil
	p.client = nil
	p.running = false
	p.mu.Unlock()

	if processer != nil {
		if err := processer.Close(); err != nil {
			logger.Warn(nil, "%+v", err)
			return err
		}
	}
	return nil
}

func (p *ConfdServer) Err() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.err
}

func (p *ConfdServer) parseConfig(pbcfg *pbtypes.ConfdConfig) (*libconfd.Config, *libconfd.BackendConfig, error) {
	sCfg, err := json.Marshal(pbcfg.GetProcessorConfig())
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return nil, nil, err
	}

	sCfgBackend, err := json.Marshal(pbcfg.GetBackendConfig())
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return nil, nil, err
	}

	cfg, err := libconfd.LoadConfigFromJsonString(string(sCfg))
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return nil, nil, err
	}

	cfgBackend, err := libconfd.LoadBackendConfigFromJsonString(string(sCfgBackend))
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return nil, nil, err
	}

	return cfg, cfgBackend, nil
}

func (p *ConfdServer) Save() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	data, err := json.MarshalIndent(p.cfg, "", "\t")
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return err
	}

	// backup old config
	bakpath := p.cfgpath + time.Now().Format(".20060102.bak")
	if false {
		if err := os.Rename(p.cfgpath, bakpath); err != nil && !os.IsNotExist(err) {
			logger.Warn(nil, "%+v", err)
			return err
		}
	}

	data = bytes.Replace(data, []byte("\n"), []byte("\r\n"), -1)
	err = ioutil.WriteFile(p.cfgpath, data, 0666)
	if err != nil {
		if false {
			os.Rename(bakpath, p.cfgpath) // revert
		}
		logger.Warn(nil, "%+v", err)
		return err
	}

	return nil
}
