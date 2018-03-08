// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package config

import (
	"context"
	"fmt"

	"github.com/coreos/etcd/mvcc/mvccpb"
	"gopkg.in/yaml.v2"

	"openpitrix.io/openpitrix/pkg/etcd"
	"openpitrix.io/openpitrix/pkg/logger"
)

type GlobalConfig struct {
	Repo    RepoServiceConfig    `yaml:"repo"`
	Cluster ClusterServiceConfig `yaml:"cluster"`
}

type RepoServiceConfig struct {
	AutoIndex bool `yaml:"auto-index"`
}

type ClusterServiceConfig struct {
	Plugins []string `yaml:"plugins"`
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
			globalConfig = UnmarshalInitConfig()
			_, err = etcd.Put(ctx, GlobalConfigKey, InitialGlobalConfig)
			if err != nil {
				return err
			}
		} else {
			err = yaml.Unmarshal(get.Kvs[0].Value, &globalConfig)
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
					err = yaml.Unmarshal(ev.Kv.Value, &globalConfig)
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

func UnmarshalInitConfig() GlobalConfig {
	var globalConfig GlobalConfig
	err := yaml.Unmarshal([]byte(InitialGlobalConfig), &globalConfig)
	if err != nil {
		fmt.Print("InitialGlobalConfig is invalid, please fix it")
		panic(err)
	}
	return globalConfig
}

func MarshalGlobalConfig(conf GlobalConfig) string {
	out, err := yaml.Marshal(conf)
	if err != nil {
		panic(err)
	}
	return string(out)
}

func init() {
	UnmarshalInitConfig()
}
