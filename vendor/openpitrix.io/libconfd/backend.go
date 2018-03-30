// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package libconfd

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

type BackendConfig struct {
	Type string   `toml:"type" json:"type"`
	Host []string `toml:"host" json:"host"`

	UserName string `toml:"user" json:"user"`
	Password string `toml:"password" json:"password"`

	ClientCAKeys string `toml:"client-ca-keys" json:"client-ca-keys"`
	ClientCert   string `toml:"client-cert" json:"client-cert"`
	ClientKey    string `toml:"client-key" json:"client-key"`
}

func (p *BackendConfig) Clone() *BackendConfig {
	var q = *p
	q.Host = append([]string{}, p.Host...)
	return &q
}

type BackendClient interface {
	Type() string
	GetValues(keys []string) (map[string]string, error)
	WatchPrefix(prefix string, keys []string, waitIndex uint64, stopChan chan bool) (uint64, error)
	WatchEnabled() bool
}

func MustNewBackendClient(cfg *BackendConfig, opts ...func(*BackendConfig)) BackendClient {
	p, err := NewBackendClient(cfg, opts...)
	if err != nil {
		logger.Panic(err)
	}
	return p
}

func NewBackendClient(cfg *BackendConfig, opts ...func(*BackendConfig)) (BackendClient, error) {
	cfg = cfg.Clone()
	for _, fn := range opts {
		fn(cfg)
	}

	newClient := _BackendClientMap[cfg.Type]
	if newClient == nil {
		return nil, fmt.Errorf("libconfd: unknown backend type %q", cfg.Type)
	}

	return newClient(cfg)
}

func MustLoadBackendConfig(path string) *BackendConfig {
	p, err := LoadBackendConfig(path)
	if err != nil {
		logger.Fatal(err)
	}
	return p
}

func LoadBackendConfig(path string) (p *BackendConfig, err error) {
	p = new(BackendConfig)
	_, err = toml.DecodeFile(path, p)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func RegisterBackendClient(
	typeName string,
	newClient func(cfg *BackendConfig) (BackendClient, error),
) {
	_BackendClientMap[typeName] = newClient
}

var _BackendClientMap = map[string]func(cfg *BackendConfig) (BackendClient, error){}
