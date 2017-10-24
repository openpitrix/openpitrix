// +-------------------------------------------------------------------------
// | Copyright (C) 2017 Yunify, Inc.
// +-------------------------------------------------------------------------
// | Licensed under the Apache License, Version 2.0 (the "License");
// | you may not use this work except in compliance with the License.
// | You may obtain a copy of the License in the LICENSE file, or at:
// |
// | http://www.apache.org/licenses/LICENSE-2.0
// |
// | Unless required by applicable law or agreed to in writing, software
// | distributed under the License is distributed on an "AS IS" BASIS,
// | WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// | See the License for the specific language governing permissions and
// | limitations under the License.
// +-------------------------------------------------------------------------

package config

import (
	"io/ioutil"
	"strings"
)

type Config struct {
	Database string `json:"database" yaml:"database"`
	Host     string `json:"host" yaml:"host"`
	Port     int    `json:"port" yaml:"port"`
	Protocol string `json:"protocol" yaml:"protocol"`
	URI      string `json:"uri" yaml:"uri"`
	LogLevel string `json:"log_level" yaml:"log_level"`
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
		path = getHome() + path[1:]
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

func (p *Config) Save(path string) error {
	if strings.HasPrefix(path, "~/") || strings.HasPrefix(path, `~\`) {
		path = getHome() + path[1:]
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
	if p.Database == "" {
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
