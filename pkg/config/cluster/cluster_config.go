// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package cluster_config

import (
	"flag"
	"fmt"
	"os"

	"github.com/golang/glog"
	"github.com/koding/multiconfig"
	"github.com/pkg/errors"
)

type Config struct {
	OpenPitrix_Config
}

type OpenPitrix_Config struct {
	Glog Glog

	DB      ClusterDatabase
	Cluster ClusterService
}

type ClusterService struct {
	Host string `default:"openpitrix-cluster"`
	Port int    `default:"9103"`
}

type ClusterDatabase struct {
	Type         string `default:"mysql"`
	Host         string `default:"openpitrix-db"`
	Port         int    `default:"3306"`
	Encoding     string `default:"utf8"`
	Engine       string `default:"InnoDB"`
	DbName       string `default:"openpitrix"`
	RootName     string `default:"root"`
	RootPassword string `default:"password"`
	UserName     string `default:"openpitrix-user-cluster"`
	UserPassword string `default:"openpitrix-user-cluster-password"`
}

func (p *ClusterDatabase) GetServerAddr() string {
	return fmt.Sprintf("root:%s@tcp(%s:%d)/", p.RootPassword, p.Host, p.Port)
}

func (p *ClusterDatabase) GetUrl() string {
	return fmt.Sprintf("root:%s@tcp(%s:%d)/%s", p.RootPassword, p.Host, p.Port, p.DbName)
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

type Glog struct {
	LogToStderr     bool   `default:"false"`
	AlsoLogToStderr bool   `default:"false"`
	StderrThreshold string `default:"ERROR"` // INFO, WARNING, ERROR, FATAL
	LogDir          string `default:""`

	LogBacktraceAt string `default:""`
	V              int    `default:"0"`
	VModule        string `default:""`

	CopyStandardLogTo string `default:""`
}

func (p *Glog) ActiveFlags() {
	flag.CommandLine.Set("logtostderr", fmt.Sprintf("%v", p.LogToStderr))
	flag.CommandLine.Set("alsologtostderr", fmt.Sprintf("%v", p.AlsoLogToStderr))
	flag.CommandLine.Set("stderrthreshold", p.StderrThreshold)
	flag.CommandLine.Set("log_dir", p.LogDir)

	flag.CommandLine.Set("log_backtrace_at", p.LogBacktraceAt)
	flag.CommandLine.Set("v", fmt.Sprintf("%v", p.V))
	flag.CommandLine.Set("vmodule", p.VModule)

	if p.LogDir != "" {
		os.MkdirAll(p.LogDir, 0666)
	}
	if p.CopyStandardLogTo != "" {
		glog.CopyStandardLogTo(p.CopyStandardLogTo)
	}
}

func loadConfig(conf interface{}) error {
	d := &multiconfig.DefaultLoader{}

	d.Loader = multiconfig.MultiLoader(
		&multiconfig.TagLoader{},
		&multiconfig.EnvironmentLoader{},
		&multiconfig.FlagLoader{},
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
