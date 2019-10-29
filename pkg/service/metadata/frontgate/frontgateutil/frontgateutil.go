// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package frontgateutil

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strings"

	"openpitrix.io/openpitrix/pkg/logger"
	pbfrontgate "openpitrix.io/openpitrix/pkg/pb/metadata/frontgate"
	pbtypes "openpitrix.io/openpitrix/pkg/pb/metadata/types"
)

func MustLoadFrontgateConfig(path string) *pbtypes.FrontgateConfig {
	p, err := LoadFrontgateConfig(path)
	if err != nil {
		logger.Critical(nil, "%+v", err)
		os.Exit(1)
	}
	return p
}

func LoadFrontgateConfig(path string) (*pbtypes.FrontgateConfig, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return nil, err
	}

	p := new(pbtypes.FrontgateConfig)
	if err := json.Unmarshal(data, p); err != nil {
		logger.Warn(nil, "%+v", err)
		return nil, err
	}

	if p.ConfdConfig == nil {
		p.ConfdConfig = &pbtypes.ConfdConfig{}
	}

	if p.ConfdConfig.ProcessorConfig == nil {
		p.ConfdConfig.ProcessorConfig = &pbtypes.ConfdProcessorConfig{}
	}
	if p.ConfdConfig.BackendConfig == nil {
		p.ConfdConfig.BackendConfig = &pbtypes.ConfdBackendConfig{}
	}

	p.ConfdConfig.ProcessorConfig.Confdir = strExtractingEnvValue(
		p.ConfdConfig.ProcessorConfig.Confdir,
	)

	return p, nil
}

func DialFrontgateService(host string, port int) (
	client *pbfrontgate.FrontgateServiceClient,
	err error,
) {
	c, err := pbfrontgate.DialFrontgateService(
		"tcp", fmt.Sprintf("%s:%d", host, port),
	)
	if err != nil {
		logger.Warn(nil, "%+v", err)
		return nil, err
	}

	return c, nil
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
