// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package drone

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/golang/protobuf/proto"

	"openpitrix.io/openpitrix/pkg/libconfd"
	"openpitrix.io/openpitrix/pkg/libconfd/backends"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/pb/types"
)

type ConfdServer struct {
	mu            sync.Mutex
	cfg           *pbtypes.ConfdConfig
	config        *libconfd.Config
	backendConfig *libconfd.BackendConfig

	processor *libconfd.Processor
	client    libconfd.BackendClient
	running   bool
	err       error
}

func NewConfdServer() *ConfdServer {
	return &ConfdServer{}
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

func (p *ConfdServer) SetConfig(cfg *pbtypes.ConfdConfig) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	config, backendConfig, err := p.parseConfig(cfg)
	if err != nil {
		return err
	}

	p.cfg = proto.Clone(cfg).(*pbtypes.ConfdConfig)
	p.config = config
	p.backendConfig = backendConfig

	return nil
}

func (p *ConfdServer) GetBackendClient() libconfd.BackendClient {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.client
}

func (p *ConfdServer) Start(opts ...libconfd.Options) error {
	logger.Info("ConfdServer: Start")

	p.mu.Lock()
	defer p.mu.Unlock()

	if p.processor == nil {
		p.processor = libconfd.NewProcessor()
	}
	if p.running {
		logger.Error("ConfdServer: confd is running")
		return fmt.Errorf("drone: confd is running")
	}

	if p.cfg == nil {
		logger.Error("ConfdServer: config is nil")
		return fmt.Errorf("drone: config is nil")
	}

	if s := p.backendConfig.Type; s != backends.Etcdv3BackendType {
		logger.Error("ConfdServer: unsupport confd backend: " + s)
		return fmt.Errorf("drone: unsupport confd backend: %s", s)
	}

	backendClient, err := libconfd.NewBackendClient(p.backendConfig)
	if err != nil {
		logger.Error("ConfdServer: NewBackendClient: %v", err)
		return err
	}

	p.client = backendClient
	p.running = true

	go func() {
		logger.Info("ConfdServer: run...")

		var err = p.processor.Run(p.config, backendClient) // blocked

		p.mu.Lock()
		p.running = false
		p.client = nil
		p.err = err
		p.mu.Unlock()

		logger.Info("ConfdServer: stoped")
	}()

	return nil
}

func (p *ConfdServer) Stop() error {
	logger.Info("ConfdServer: Stop")

	p.mu.Lock()
	var processer = p.processor
	p.processor = nil
	p.client = nil
	p.running = false
	p.mu.Unlock()

	if processer != nil {
		return processer.Close()
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
		return nil, nil, err
	}

	sCfgBackend, err := json.Marshal(pbcfg.GetBackendConfig())
	if err != nil {
		return nil, nil, err
	}

	cfg, err := libconfd.LoadConfigFromJsonString(string(sCfg))
	if err != nil {
		return nil, nil, err
	}

	cfgBackend, err := libconfd.LoadBackendConfigFromJsonString(string(sCfgBackend))
	if err != nil {
		return nil, nil, err
	}

	return cfg, cfgBackend, nil
}
