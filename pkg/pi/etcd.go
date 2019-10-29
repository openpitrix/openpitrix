// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package pi

import (
	"context"

	"github.com/coreos/etcd/mvcc/mvccpb"

	"openpitrix.io/openpitrix/pkg/config"
	"openpitrix.io/openpitrix/pkg/etcd"
	"openpitrix.io/openpitrix/pkg/logger"
)

const (
	EtcdPrefix      = "openpitrix/"
	GlobalConfigKey = "global_config"
	DlockKey        = "dlock_" + GlobalConfigKey
)

type Watcher chan *config.GlobalConfig

func WatchGlobalConfig(etcd *etcd.Etcd, watcher Watcher) error {
	ctx := context.Background()
	var globalConfig config.GlobalConfig
	err := etcd.Dlock(ctx, DlockKey, func() error {
		// get value
		get, err := etcd.Get(ctx, GlobalConfigKey)
		if err != nil {
			return err
		}
		// parse value
		if get.Count == 0 {
			logger.Debug(nil, "Cannot get global config, put the initial string. [%s]", config.InitialGlobalConfig)
			globalConfig = config.DecodeInitConfig()
			_, err = etcd.Put(ctx, GlobalConfigKey, config.InitialGlobalConfig)
			if err != nil {
				return err
			}
		} else {
			globalConfig, err = config.ParseGlobalConfig(get.Kvs[0].Value)
			if err != nil {
				return err
			}
		}
		logger.Debug(nil, "Global config update to [%+v]", globalConfig)
		// send it back
		watcher <- &globalConfig
		return nil
	})

	// watch
	go func() {
		logger.Debug(nil, "Start watch global config")
		watchRes := etcd.Watch(ctx, GlobalConfigKey)
		for res := range watchRes {
			for _, ev := range res.Events {
				if ev.Type == mvccpb.PUT {
					//logger.Debug(nil, "Got updated global config from etcd, try to decode with yaml")
					globalConfig, err := config.ParseGlobalConfig(ev.Kv.Value)
					if err != nil {
						logger.Error(nil, "Watch global config from etcd found error: %+v", err)
					} else {
						watcher <- &globalConfig
					}
				}
			}
		}
	}()
	return err
}
