// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package drone

import (
	"fmt"
	"sync"

	"openpitrix.io/libconfd"
)

type ConfdServer struct {
	mu            sync.Mutex
	config        *libconfd.Config
	backendConfig *libconfd.BackendConfig
	processor     *libconfd.Processor
	client        libconfd.BackendClient
	running       bool
	err           error
}

func NewConfdServer() *ConfdServer {
	return &ConfdServer{}
}

func (p *ConfdServer) IsRunning() bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.processor != nil && p.running
}

func (p *ConfdServer) GetConfdConfig() *libconfd.Config {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.config != nil {
		return p.config
	}

	return &libconfd.Config{}
}

func (p *ConfdServer) GetBackendConfig() *libconfd.BackendConfig {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.backendConfig != nil {
		return p.backendConfig
	}

	return &libconfd.BackendConfig{}
}

func (p *ConfdServer) GetBackendClient() libconfd.BackendClient {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.client
}

func (p *ConfdServer) Start(config *libconfd.Config, backendConfig *libconfd.BackendConfig) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.processor == nil {
		p.processor = libconfd.NewProcessor()
	}
	if p.running {
		return fmt.Errorf("drone: confd is running!")
	}

	p.config = config.Clone()
	p.backendConfig = backendConfig.Clone()

	etcdClient, err := NewEtcdClient(p.backendConfig)
	if err != nil {
		return err
	}

	p.client = etcdClient
	p.running = true

	go func() {
		var err = p.processor.Run(p.config, etcdClient) // blocked

		p.mu.Lock()
		p.running = false
		p.client = nil
		p.err = err
		p.mu.Unlock()
	}()

	return nil
}

func (p *ConfdServer) Stop() error {
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
