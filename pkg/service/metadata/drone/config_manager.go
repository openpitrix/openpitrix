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
	"sync"
	"time"

	"github.com/golang/protobuf/proto"

	"openpitrix.io/openpitrix/pkg/pb/types"
)

type ConfigManager struct {
	path string
	cfg  *pbtypes.DroneConfig
	mu   sync.Mutex
}

func NewConfigManager(path string, cfg *pbtypes.DroneConfig, opts ...Options) *ConfigManager {
	if cfg != nil {
		cfg = proto.Clone(cfg).(*pbtypes.DroneConfig)
	} else {
		cfg = NewDefaultConfig()
	}

	for _, fn := range opts {
		fn(cfg)
	}

	return &ConfigManager{
		path: path,
		cfg:  cfg,
	}
}

func (p *ConfigManager) Get() (cfg *pbtypes.DroneConfig) {
	p.mu.Lock()
	defer p.mu.Unlock()

	cfg = proto.Clone(p.cfg).(*pbtypes.DroneConfig)
	return
}

func (p *ConfigManager) Set(cfg *pbtypes.DroneConfig) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if cfg.Id != "" && cfg.Id != p.cfg.Id {
		return fmt.Errorf("drone: config.Id is read only")
	}
	if cfg.ListenPort > 0 && cfg.ListenPort != p.cfg.ListenPort {
		return fmt.Errorf("drone: config.ListenPort is read only")
	}

	cfg.Id = p.cfg.Id
	cfg.ListenPort = p.cfg.ListenPort

	p.cfg = proto.Clone(cfg).(*pbtypes.DroneConfig)
	return nil
}

func (p *ConfigManager) Save() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	data, err := json.MarshalIndent(p.cfg, "", "\t")
	if err != nil {
		return err
	}

	// backup old config
	bakpath := p.path + time.Now().Format(".20060102.bak")
	if err := os.Rename(p.path, bakpath); err != nil {
		return err
	}

	data = bytes.Replace(data, []byte("\n"), []byte("\r\n"), -1)
	err = ioutil.WriteFile(p.path, data, 0666)
	if err != nil {
		os.Rename(bakpath, p.path) // revert
		return err
	}

	return nil
}
