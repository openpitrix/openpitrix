// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package config

import (
	"fmt"
	"strings"

	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/util/yamlutil"
)

type GlobalConfig struct {
	Repo    RepoServiceConfig      `json:"repo"`
	Cluster ClusterServiceConfig   `json:"cluster"`
	Runtime map[string]ImageConfig `json:"runtime"`
	Pilot   PilotServiceConfig     `json:"pilot"`
}

type RepoServiceConfig struct {
	Cron string `json:"cron"`
}

type ClusterServiceConfig struct {
	Plugins             []string `json:"plugins"`
	FrontgateConf       string   `json:"frontgate_conf"`
	FrontgateAutoDelete bool     `json:"frontgate_auto_delete"`
}

type PilotServiceConfig struct {
	Ip   string `json:"ip"`
	Port int32  `json:"port"`
}

type ImageConfig struct {
	ApiServer string `json:"api_server"`
	Zone      string `json:"zone"`
	ImageId   string `json:"image_id"`
	ImageUrl  string `json:"image_url"`
	ImageName string `json:"image_name"`
}

func (g *GlobalConfig) GetFrontgateAutoDelete() bool {
	if g.Cluster.FrontgateAutoDelete == false {
		return false
	} else {
		return true
	}
}

func (g *GlobalConfig) GetRuntimeImageIdAndUrl(apiServer, zone string) (*ImageConfig, error) {
	if strings.HasPrefix(apiServer, "https://") {
		apiServer = strings.Split(apiServer, "https://")[1]
	}
	for _, imageConfig := range g.Runtime {
		if imageConfig.ApiServer == apiServer && imageConfig.Zone == zone {
			return &imageConfig, nil
		}
	}
	for _, imageConfig := range g.Runtime {
		if imageConfig.ApiServer == apiServer && imageConfig.Zone == ".*" {
			return &imageConfig, nil
		}
	}
	logger.Error(nil, "No such runtime image with api server [%s] zone [%s]. ", apiServer, zone)
	return nil, fmt.Errorf("no such runtime image with api server [%s] zone [%s]. ", apiServer, zone)
}

func ParseGlobalConfig(data []byte) (GlobalConfig, error) {
	var globalConfig GlobalConfig
	err := yamlutil.Decode(data, &globalConfig)
	if err != nil {
		return globalConfig, err
	}
	return globalConfig, nil
}

func DecodeInitConfig() GlobalConfig {
	globalConfig, err := ParseGlobalConfig([]byte(InitialGlobalConfig))
	if err != nil {
		fmt.Print("InitialGlobalConfig is invalid, please fix it")
		panic(err)
	}
	return globalConfig
}

func EncodeGlobalConfig(conf GlobalConfig) string {
	out, err := yamlutil.Encode(conf)
	if err != nil {
		panic(err)
	}
	return string(out)
}

func init() {
	DecodeInitConfig()
}
