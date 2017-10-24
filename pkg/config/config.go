// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package config

import (
	"io/ioutil"
	"strings"
)

type Config struct {
	DbType     string `json:"db_type" yaml:"db_type"`
	DbHost     string `json:"db_host" yaml:"db_host"`
	DbEncoding string `json:"db_encoding" yaml:"db_encoding"`
	DbEngine   string `json:"db_engine" yaml:"db_engine"`
	Host       string `json:"host" yaml:"host"`
	Port       int    `json:"port" yaml:"port"`
	Protocol   string `json:"protocol" yaml:"protocol"`
	URI        string `json:"uri" yaml:"uri"`
	LogLevel   string `json:"log_level" yaml:"log_level"`
}

func Default() *Config {
	p := new(Config)
	if err := yamlDecode([]byte(DefaultConfigContent), p); err != nil {
		panic(err)
	}
	return p
}

func Load(path string) (*Config, error) {
	if strings.HasPrefix(path, "~/") || strings.HasPrefix(path, `~\`) {
		path = GetHomePath() + path[1:]
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	p := new(Config)
	if err := yamlDecode(data, p); err != nil {
		panic(err)
	}

	return p, nil
}

func MustLoad(path string) *Config {
	p, err := Load(path)
	if err != nil {
		panic(err)
	}
	return p
}

func Parse(content string) (*Config, error) {
	p := new(Config)
	if err := yamlDecode([]byte(content), p); err != nil {
		return nil, err
	}
	return p, nil
}

func (p *Config) Clone() *Config {
	q, _ := Parse(p.String())
	return q
}

func (p *Config) Save(path string) error {
	if strings.HasPrefix(path, "~/") || strings.HasPrefix(path, `~\`) {
		path = GetHomePath() + path[1:]
	}

	data, err := yamlEncode(p)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (p *Config) Valid() bool {
	if p.DbType == "" || p.DbHost == "" {
		return false
	}
	if p.Host == "" {
		return false
	}
	if p.Port <= 0 {
		return false
	}
	if p.Protocol == "" {
		return false
	}
	if p.URI == "" {
		return false
	}
	if p.LogLevel == "" {
		return false
	}

	// OK
	return true
}

func (p *Config) String() string {
	data, _ := yamlEncode(p)
	return string(data)
}
