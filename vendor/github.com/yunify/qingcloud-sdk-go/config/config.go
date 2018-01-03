// +-------------------------------------------------------------------------
// | Copyright (C) 2016 Yunify, Inc.
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
	"net/http"
	"os"
	"strings"

	"github.com/yunify/qingcloud-sdk-go/logger"
	"github.com/yunify/qingcloud-sdk-go/utils"
	"net"
	"time"
)

// A Config stores a configuration of this sdk.
type Config struct {
	AccessKeyID     string `yaml:"qy_access_key_id"`
	SecretAccessKey string `yaml:"qy_secret_access_key"`

	Host              string `yaml:"host"`
	Port              int    `yaml:"port"`
	Protocol          string `yaml:"protocol"`
	URI               string `yaml:"uri"`
	ConnectionRetries int    `yaml:"connection_retries"`
	ConnectionTimeout int    `yaml:"connection_timeout"`

	LogLevel string `yaml:"log_level"`

	Zone string `yaml:"zone"`

	Connection *http.Client
}

// New create a Config with given AccessKeyID and SecretAccessKey.
func New(accessKeyID, secretAccessKey string) (*Config, error) {
	config, err := NewDefault()
	if err != nil {
		return nil, err
	}

	config.AccessKeyID = accessKeyID
	config.SecretAccessKey = secretAccessKey

	config.Connection = &http.Client{}

	return config, nil
}

// NewDefault create a Config with default configuration.
func NewDefault() (*Config, error) {
	config := &Config{}
	err := config.LoadDefaultConfig()
	if err != nil {
		return nil, err
	}

	timeout := time.Duration(config.ConnectionTimeout) * time.Second
	transport := &http.Transport{
		Dial: func(network, addr string) (net.Conn, error) {
			return net.DialTimeout(network, addr, timeout)
		},
	}
	config.Connection = &http.Client{
		Transport: transport,
	}

	return config, nil
}

// LoadDefaultConfig loads the default configuration for Config.
// It returns error if yaml decode failed.
func (c *Config) LoadDefaultConfig() error {
	_, err := utils.YAMLDecode([]byte(DefaultConfigFileContent), c)
	if err != nil {
		logger.Error("Config parse error: " + err.Error())
		return err
	}

	logger.SetLevel(c.LogLevel)

	return nil
}

// LoadUserConfig loads user configuration in ~/.qingcloud/config.yaml for Config.
// It returns error if file not found.
func (c *Config) LoadUserConfig() error {
	_, err := os.Stat(GetUserConfigFilePath())
	if err != nil {
		logger.Warn("Installing default config file to \"" + GetUserConfigFilePath() + "\"")
		InstallDefaultUserConfig()
	}

	return c.LoadConfigFromFilepath(GetUserConfigFilePath())
}

// LoadConfigFromFilepath loads configuration from a specified local path.
// It returns error if file not found or yaml decode failed.
func (c *Config) LoadConfigFromFilepath(filepath string) error {
	if strings.Index(filepath, "~/") == 0 {
		filepath = strings.Replace(filepath, "~/", getHome()+"/", 1)
	}

	configYAML, err := ioutil.ReadFile(filepath)
	if err != nil {
		logger.Error("File not found: " + filepath)
		return err
	}

	return c.LoadConfigFromContent(configYAML)
}

// LoadConfigFromContent loads configuration from a given byte slice.
// It returns error if yaml decode failed.
func (c *Config) LoadConfigFromContent(content []byte) error {
	c.LoadDefaultConfig()

	_, err := utils.YAMLDecode(content, c)
	if err != nil {
		logger.Error("Config parse error: " + err.Error())
		return err
	}

	logger.SetLevel(c.LogLevel)

	timeout := time.Duration(c.ConnectionTimeout) * time.Second
	transport := &http.Transport{
		Dial: func(network, addr string) (net.Conn, error) {
			return net.DialTimeout(network, addr, timeout)
		},
	}
	c.Connection = &http.Client{
		Transport: transport,
	}

	return nil
}
