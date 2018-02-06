// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package unittest_config

import (
	"github.com/koding/multiconfig"
	"github.com/pkg/errors"
)

type Config struct {
	OpenPitrix_Config
}

type OpenPitrix_Config struct {
	Unittest Unittest
}

type Unittest struct {
	K8s  Kubernetes
	QC   QingCloud
	Rest RestAPI
}

type Kubernetes struct {
	Enabled bool `default:"false"`
}

type QingCloud struct {
	Enabled bool `default:"false"`
}

type RestAPI struct {
	Enabled  bool   `default:"false"`
	Host     string `default:"localhost:9100"`
	BasePath string `default:"/"`
}

func PrintEnvs() {
	new(multiconfig.EnvironmentLoader).PrintEnvs(&OpenPitrix_Config{})
}

func MustLoadConfig() *Config {
	cfg, err := LoadConfig()
	if err != nil {
		panic(err)
	}
	return cfg
}

func LoadConfig() (*Config, error) {
	p := new(OpenPitrix_Config)
	if err := loadConfig(p); err != nil {
		return nil, err
	}
	return &Config{*p}, nil
}

func loadConfig(conf interface{}) error {
	d := &multiconfig.DefaultLoader{}

	d.Loader = multiconfig.MultiLoader(
		&multiconfig.TagLoader{},
		&multiconfig.EnvironmentLoader{},
	)
	d.Validator = multiconfig.MultiValidator(
		&multiconfig.RequiredValidator{},
	)

	if err := d.Load(conf); err != nil {
		err = errors.WithStack(err)
		return err
	}
	return nil
}
