// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package api_config

import (
	"flag"
	"fmt"
	"os"

	"github.com/golang/glog"
	"github.com/koding/multiconfig"
	"github.com/pkg/errors"
)

type ApiConfig struct {
	OpenPitrix_Api_Config
}

type OpenPitrix_Api_Config struct {
	Host string `default:"openpitrix-api"`
	Port int    `default:"9100"`

	App     AppService
	Runtime RuntimeService
	Cluster ClusterService
	Repo    RepoService

	Glog Glog
}

type AppService struct {
	Host string `default:"openpitrix-app"`
	Port int    `default:"9101"`
}

type RuntimeService struct {
	Host string `default:"openpitrix-runtime"`
	Port int    `default:"9102"`
}

type ClusterService struct {
	Host string `default:"openpitrix-cluster"`
	Port int    `default:"9103"`
}

type RepoService struct {
	Host string `default:"openpitrix-repo"`
	Port int    `default:"9104"`
}

func LoadApiConfig() (*ApiConfig, error) {
	p := new(OpenPitrix_Api_Config)
	if err := loadConfig(p); err != nil {
		return nil, err
	}
	return &ApiConfig{*p}, nil
}

type Glog struct {
	LogToStderr     bool   `default:"false"`
	AlsoLogTostderr bool   `default:"false"`
	StderrThreshold string `default:"ERROR"` // INFO, WARNING, ERROR, FATAL
	LogDir          string `default:""`

	LogBacktraceAt string `default:""`
	V              int    `default:"0"`
	VModule        string `default:""`

	CopyStandardLogTo string `default:""`
}

func (p *Glog) ActiveFlags() {
	flag.CommandLine.Set("logtostderr", fmt.Sprintf("%v", p.LogToStderr))
	flag.CommandLine.Set("alsologtostderr", fmt.Sprintf("%v", p.AlsoLogTostderr))
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
