// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package config

import (
	"context"
	"fmt"
	"strings"

	"github.com/coreos/etcd/mvcc/mvccpb"

	"openpitrix.io/openpitrix/pkg/etcd"
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
	Plugins       []string `json:"plugins"`
	FrontgateConf string   `json:"frontgate_conf"`
}

type PilotServiceConfig struct {
	Ip string `json:"ip"`
}

type ImageConfig struct {
	ApiServer string `json:"api_server"`
	Zone      string `json:"zone"`
	ImageId   string `json:"image_id"`
	ImageUrl  string `json:"image_url"`
}

func (g *GlobalConfig) GetRuntimeImageIdAndUrl(apiServer, zone string) (string, string, error) {
	if strings.HasPrefix(apiServer, "https://") {
		apiServer = strings.Split(apiServer, "https://")[1]
	}
	for _, imageConfig := range g.Runtime {
		if imageConfig.ApiServer == apiServer && imageConfig.Zone == zone {
			return imageConfig.ImageId, imageConfig.ImageUrl, nil
		}
	}
	for _, imageConfig := range g.Runtime {
		if imageConfig.ApiServer == apiServer && imageConfig.Zone == ".*" {
			return imageConfig.ImageId, imageConfig.ImageUrl, nil
		}
	}
	logger.Error("No such runtime image with api server [%s] zone [%s]. ", apiServer, zone)
	return "", "", fmt.Errorf("no such runtime image with api server [%s] zone [%s]. ", apiServer, zone)
}

const (
	EtcdPrefix      = "openpitrix/"
	GlobalConfigKey = "global_config"
	DlockKey        = "dlock_" + GlobalConfigKey
)

type Watcher chan *GlobalConfig

func WatchGlobalConfig(etcd *etcd.Etcd, watcher Watcher) error {
	ctx := context.Background()
	var globalConfig GlobalConfig
	err := etcd.Dlock(ctx, DlockKey, func() error {
		// get value
		get, err := etcd.Get(ctx, GlobalConfigKey)
		if err != nil {
			return err
		}
		// parse value
		if get.Count == 0 {
			logger.Debug("Cannot get global config, put the initial string. [%s]", InitialGlobalConfig)
			globalConfig = DecodeInitConfig()
			_, err = etcd.Put(ctx, GlobalConfigKey, InitialGlobalConfig)
			if err != nil {
				return err
			}
		} else {
			globalConfig, err = ParseGlobalConfig(get.Kvs[0].Value)
			if err != nil {
				return err
			}
		}
		logger.Debug("Global config update to [%+v]", globalConfig)
		// send it back
		watcher <- &globalConfig
		return nil
	})

	// watch
	go func() {
		logger.Debug("Start watch global config")
		watchRes := etcd.Watch(ctx, GlobalConfigKey)
		for res := range watchRes {
			for _, ev := range res.Events {
				if ev.Type == mvccpb.PUT {
					//logger.Debug("Got updated global config from etcd, try to decode with yaml")
					globalConfig, err := ParseGlobalConfig(ev.Kv.Value)
					if err != nil {
						logger.Error("Watch global config from etcd found error: %+v", err)
					} else {
						watcher <- &globalConfig
					}
				}
			}
		}
	}()
	return err
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
