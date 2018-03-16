// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package config

import (
	"context"
	"fmt"

	"github.com/coreos/etcd/mvcc/mvccpb"

	"openpitrix.io/openpitrix/pkg/etcd"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/utils/yaml"
)

type GlobalConfig struct {
	Repo    RepoServiceConfig    `json:"repo"`
	Cluster ClusterServiceConfig `json:"cluster"`
}

type RepoServiceConfig struct {
	AutoIndex bool `json:"auto-index"`
}

type ClusterServiceConfig struct {
	Plugins []string `json:"plugins"`
}

const InitialGlobalConfig = `
repo:
  auto-index: true
cluster:
  plugins:
    - qingcloud
    - kubernetes
`

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
			logger.Debugf("Cannot get global config, put the initial string. [%s]", InitialGlobalConfig)
			globalConfig = DecodeInitConfig()
			_, err = etcd.Put(ctx, GlobalConfigKey, InitialGlobalConfig)
			if err != nil {
				return err
			}
		} else {
			err = yaml.Decode(get.Kvs[0].Value, &globalConfig)
			if err != nil {
				return err
			}
		}
		logger.Debugf("Got global config [%+v]", globalConfig)
		// send it back
		watcher <- &globalConfig
		return nil
	})

	// watch
	go func() {
		logger.Debugf("Start watch global config")
		watchRes := etcd.Watch(ctx, GlobalConfigKey)
		for res := range watchRes {
			for _, ev := range res.Events {
				if ev.Type == mvccpb.PUT {
					logger.Debugf("Got global config from etcd")
					err = yaml.Decode(ev.Kv.Value, &globalConfig)
					if err != nil {
						logger.Errorf("Watch global config from etcd found error: %+v", err)
					} else {
						watcher <- &globalConfig
					}
				}
			}
		}
	}()
	return err
}

func DecodeInitConfig() GlobalConfig {
	var globalConfig GlobalConfig
	err := yaml.Decode([]byte(InitialGlobalConfig), &globalConfig)
	if err != nil {
		fmt.Print("InitialGlobalConfig is invalid, please fix it")
		panic(err)
	}
	return globalConfig
}

func EncodeGlobalConfig(conf GlobalConfig) string {
	out, err := yaml.Encode(conf)
	if err != nil {
		panic(err)
	}
	return string(out)
}

func init() {
	DecodeInitConfig()
}
