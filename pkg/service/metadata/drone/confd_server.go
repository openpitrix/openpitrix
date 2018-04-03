// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package drone

import (
	"sync"

	"openpitrix.io/libconfd"
)

type ConfdServer struct {
	confdProcessor *libconfd.Processor
	confdCall      *libconfd.Call
	confdMutex     sync.Mutex
}

func NewConfdServer() *ConfdServer {
	return &ConfdServer{
		confdProcessor: libconfd.NewProcessor(),
	}
}

func (p *ConfdServer) GetConfdConfig() *libconfd.Config {
	return &libconfd.Config{}
}

func (p *ConfdServer) GetBackendConfig() *libconfd.BackendConfig {
	return &libconfd.BackendConfig{}
}

func (p *ConfdServer) Start(cfg *libconfd.Config, cfgBcakend *libconfd.BackendConfig) error {
	return nil
}

func (p *ConfdServer) Stop() error {
	return nil
}

func (p *ConfdServer) Err() error {
	return nil
}
