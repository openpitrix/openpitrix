// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package config

import (
	"bytes"
	"encoding/gob"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	pathpkg "path"
	"runtime"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/koding/multiconfig"
	"github.com/pkg/errors"

	"openpitrix.io/openpitrix/pkg/logger"
)

type Config struct {
	OpenPitrix_Config
}

type OpenPitrix_Config struct {
	Glog Glog

	DB      Database
	Api     ApiService
	App     AppService
	Runtime RuntimeService
	Cluster ClusterService
	Repo    RepoService

	Unittest Unittest
}

type ApiService struct {
	Host string `default:"127.0.0.1"`
	Port int    `default:"8080"`
}

type AppService struct {
	Host string `default:"127.0.0.1"`
	Port int    `default:"8081"`
}

type RuntimeService struct {
	Host string `default:"127.0.0.1"`
	Port int    `default:"8082"`
}

type ClusterService struct {
	Host string `default:"127.0.0.1"`
	Port int    `default:"8083"`
}

type RepoService struct {
	Host string `default:"127.0.0.1"`
	Port int    `default:"8084"`
}

type Database struct {
	Type         string `default:"mysql"`
	Host         string `default:"127.0.0.1"`
	Port         int    `default:"3306"`
	Encoding     string `default:"utf8"`
	Engine       string `default:"InnoDB"`
	DbName       string `default:"openpitrix"`
	RootPassword string `default:"password"`
}

type Unittest struct {
	EnableDbTest bool `default:"false"`
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

func (p *Database) GetUrl() string {
	return fmt.Sprintf("root:%s@tcp(%s:%d)/%s", p.RootPassword, p.Host, p.Port, p.DbName)
}

func Default() *Config {
	p, err := Parse(DefaultConfigContent)
	if err != nil {
		if err == flag.ErrHelp {
			fmt.Println("See https://openpitrix.io")
			os.Exit(0)
		}
		logger.Fatalf("%+v", err)
	}

	return p
}

func Load(path string) (*Config, error) {
	p := new(OpenPitrix_Config)

	if strings.HasPrefix(path, "~/") || strings.HasPrefix(path, `~\`) {
		path = GetHomePath() + path[1:]
	}

	if err := multiconfig.NewWithPath(path).Load(p); err != nil {
		if err == flag.ErrHelp {
			fmt.Println("See https://openpitrix.io")
			os.Exit(0)
		}
		err = errors.WithStack(err)
		return nil, err
	}

	return &Config{*p}, nil
}

func MustLoad(path string) *Config {
	p, err := Load(path)
	if err != nil {
		if err == flag.ErrHelp {
			fmt.Println("See https://openpitrix.io")
			os.Exit(0)
		}
		logger.Fatalf("%+v", err)
	}
	return p
}

func MustLoadUserConfig() *Config {
	path := DefaultConfigPath
	if strings.HasPrefix(path, "~/") || strings.HasPrefix(path, `~\`) {
		path = GetHomePath() + path[1:]
	}
	if _, err := os.Stat(path); err != nil {
		os.MkdirAll(pathpkg.Dir(path), 0755)
		ioutil.WriteFile(path, []byte(DefaultConfigContent), 0644)
	}

	return MustLoad(path)
}

func MustLoadUnittestConfig() *Config {
	path := UnittestConfigPath
	if strings.HasPrefix(path, "~/") || strings.HasPrefix(path, `~\`) {
		path = GetHomePath() + path[1:]
	}
	if _, err := os.Stat(path); err != nil {
		if err := os.MkdirAll(pathpkg.Dir(path), 0755); err != nil {
			err = errors.WithStack(err)
			logger.Warningf("%+v", err)
		}
		if err := ioutil.WriteFile(path, []byte(UnittestConfigContent), 0644); err != nil {
			err = errors.WithStack(err)
			logger.Warningf("%+v", err)
		}
	}

	return MustLoad(path)
}

func Parse(content string) (*Config, error) {
	p := new(OpenPitrix_Config)

	d := &multiconfig.DefaultLoader{}
	d.Loader = multiconfig.MultiLoader(
		&multiconfig.TagLoader{},
		&multiconfig.TOMLLoader{Reader: strings.NewReader(content)},
		&multiconfig.EnvironmentLoader{},
		&multiconfig.FlagLoader{},
	)
	d.Validator = multiconfig.MultiValidator(&multiconfig.RequiredValidator{})

	if err := d.Load(p); err != nil {
		if err == flag.ErrHelp {
			fmt.Println("See https://openpitrix.io")
			os.Exit(0)
		}
		err = errors.WithStack(err)
		return nil, err
	}

	return &Config{*p}, nil
}

func (p *Config) Clone() *Config {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(p); err != nil {
		err = errors.WithStack(err)
		logger.Fatalf("%+v", err)
	}

	var q Config
	if err := gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(&q); err != nil {
		err = errors.WithStack(err)
		logger.Fatalf("%+v", err)
	}

	return &q
}

func (p *Config) Save(path string) error {
	if strings.HasPrefix(path, "~/") || strings.HasPrefix(path, `~\`) {
		path = GetHomePath() + path[1:]
	}

	buf := new(bytes.Buffer)
	err := toml.NewEncoder(buf).Encode(p)
	if err != nil {
		err = errors.WithStack(err)
		return err
	}

	err = ioutil.WriteFile(path, buf.Bytes(), 0644)
	if err != nil {
		err = errors.WithStack(err)
		return err
	}

	return nil
}

func (p *Config) String() string {
	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(p); err != nil {
		return ""
	}
	return (buf.String())
}

func GetHomePath() string {
	home := os.Getenv("HOME")
	if runtime.GOOS == "windows" {
		home = os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
	}
	if home == "" {
		home = "~"
	}

	return home
}

func IsUserConfigExists() bool {
	cfgpath := DefaultConfigPath
	if strings.HasPrefix(cfgpath, "~/") || strings.HasPrefix(cfgpath, `~\`) {
		cfgpath = GetHomePath() + cfgpath[1:]
	}

	fi, err := os.Stat(cfgpath)
	if err != nil {
		return false
	}
	if fi.IsDir() {
		return false
	}

	return true
}
