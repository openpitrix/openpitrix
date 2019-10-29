// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package droneutil

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strings"

	"google.golang.org/grpc"

	"openpitrix.io/openpitrix/pkg/logger"
	pbdrone "openpitrix.io/openpitrix/pkg/pb/metadata/drone"
	pbtypes "openpitrix.io/openpitrix/pkg/pb/metadata/types"
)

func MustLoadConfdConfig(path string) *pbtypes.ConfdConfig {
	p, err := LoadConfdConfig(path)
	if err != nil {
		logger.Critical(nil, "%+v", err)
		os.Exit(1)
	}
	return p
}

func MustLoadDroneConfig(path string) *pbtypes.DroneConfig {
	p, err := LoadDroneConfig(path)
	if err != nil {
		logger.Critical(nil, "%+v", err)
		os.Exit(1)
	}
	return p
}

func LoadConfdConfig(path string) (*pbtypes.ConfdConfig, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return nil, err
	}

	p := new(pbtypes.ConfdConfig)
	if err := json.Unmarshal(data, p); err != nil {
		logger.Warn(nil, "%+v", err)
		return nil, err
	}

	if p.ProcessorConfig == nil {
		p.ProcessorConfig = &pbtypes.ConfdProcessorConfig{}
	}
	if p.BackendConfig == nil {
		p.BackendConfig = &pbtypes.ConfdBackendConfig{}
	}

	p.ProcessorConfig.Confdir = strExtractingEnvValue(
		p.ProcessorConfig.Confdir,
	)

	return p, nil
}

func LoadDroneConfig(path string) (*pbtypes.DroneConfig, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return nil, err
	}

	p := new(pbtypes.DroneConfig)
	if err := json.Unmarshal(data, p); err != nil {
		logger.Warn(nil, "%+v", err)
		return nil, err
	}

	p.CmdInfoLogPath = strExtractingEnvValue(p.CmdInfoLogPath)
	return p, nil
}

func DialDroneService(ctx context.Context, host string, port int) (
	client pbdrone.DroneServiceClient,
	conn *grpc.ClientConn,
	err error,
) {
	conn, err = grpc.Dial(fmt.Sprintf("%s:%d", host, port), grpc.WithInsecure())
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return
	}

	client = pbdrone.NewDroneServiceClient(conn)
	return
}

func strExtractingEnvValue(s string) string {
	if !strings.ContainsAny(s, "${}") {
		return s
	}

	env := os.Environ()
	if runtime.GOOS == "windows" {
		if os.Getenv("HOME") == "" {
			home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
			if home == "" {
				home = os.Getenv("USERPROFILE")
			}

			env = append(env, "HOME="+home)
		}

		if os.Getenv("PWD") == "" {
			pwd, _ := os.Getwd()
			env = append(env, "PWD="+pwd)
		}
	}

	for _, e := range env {
		if i := strings.Index(e, "="); i >= 0 {
			s = strings.Replace(s,
				fmt.Sprintf("${%s}", strings.TrimSpace(e[:i])),
				strings.TrimSpace(e[i+1:]),
				-1,
			)
		}
	}
	return s
}
