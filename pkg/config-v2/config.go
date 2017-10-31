// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package config

import (
	"bytes"
	"encoding/gob"
	"io/ioutil"
	"os"
	"runtime"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/koding/multiconfig"
)

type OpenPitrix struct {
	Database Database

	Host string `default:"127.0.0.1"`
	Port int    `default:"8080"`

	// Valid log levels are "debug", "info", "warn", "error", and "fatal".
	LogLevel string `default:"warn"`
}

type Database struct {
	Type     string `default:"mysql"`
	Host     string `default:"root:password@tcp(127.0.0.1:3306)"`
	Encoding string `default:"utf8"`
	Engine   string `default:"InnoDB"`
	DbName   string `default:"openpitrix"`
}

func Default() *OpenPitrix {
	p := new(OpenPitrix)

	loader := &multiconfig.TOMLLoader{
		Reader: strings.NewReader(DefaultConfigContent),
	}
	if err := loader.Load(p); err != nil {
		panic(err)
	}

	return p
}

func Load(path string) (*OpenPitrix, error) {
	p := new(OpenPitrix)

	if strings.HasPrefix(path, "~/") || strings.HasPrefix(path, `~\`) {
		path = GetHomePath() + path[1:]
	}

	if err := multiconfig.NewWithPath(path).Load(p); err != nil {
		return nil, err
	}

	return p, nil
}

func MustLoad(path string) *OpenPitrix {
	p, err := Load(path)
	if err != nil {
		panic(err)
	}
	return p
}

func Parse(content string) (*OpenPitrix, error) {
	p := new(OpenPitrix)

	d := &multiconfig.DefaultLoader{}
	d.Loader = multiconfig.MultiLoader(
		&multiconfig.TagLoader{},
		&multiconfig.TOMLLoader{Reader: strings.NewReader(content)},
		&multiconfig.EnvironmentLoader{},
		&multiconfig.FlagLoader{},
	)
	d.Validator = multiconfig.MultiValidator(&multiconfig.RequiredValidator{})

	if err := d.Load(p); err != nil {
		return nil, err
	}

	return p, nil
}

func (p *OpenPitrix) Clone() *OpenPitrix {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(p); err != nil {
		panic(err)
	}

	var q OpenPitrix
	if err := gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(&q); err != nil {
		panic(err)
	}

	return &q
}

func (p *OpenPitrix) Save(path string) error {
	if strings.HasPrefix(path, "~/") || strings.HasPrefix(path, `~\`) {
		path = GetHomePath() + path[1:]
	}

	buf := new(bytes.Buffer)
	err := toml.NewEncoder(buf).Encode(p)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, buf.Bytes(), 0644)
	if err != nil {
		return err
	}

	return nil
}

func (p *OpenPitrix) String() string {
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
