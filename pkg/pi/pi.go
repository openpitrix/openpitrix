// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package pi

import (
	"context"
	"strings"
	"sync"

	"github.com/google/gops/agent"

	"openpitrix.io/openpitrix/pkg/config"
	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/etcd"
	"openpitrix.io/openpitrix/pkg/logger"
)

type globalCfgWatcher func(*config.GlobalConfig)

type Pi struct {
	cfg              *config.Config
	globalCfg        *config.GlobalConfig
	globalCfgWatcher []globalCfgWatcher
	database         *db.Database
	etcd             *etcd.Etcd
}

var global *Pi
var mutex sync.RWMutex
var globalOnce sync.Once
var gopsOnce sync.Once

func NewPi(cfg *config.Config) *Pi {
	p := &Pi{cfg: cfg}
	p.openDatabase()
	p.openEtcd()
	p.watchGlobalCfg()

	if !cfg.DisableGops {
		gopsOnce.Do(func() {
			if err := agent.Listen(agent.Options{
				ShutdownCleanup: true,
			}); err != nil {
				logger.Critical(nil, "failed to start gops agent")
			}
		})
	}

	return p
}

func SetGlobal(cfg *config.Config) {
	globalOnce.Do(func() {
		global = NewPi(cfg)
	})
}

func Global() *Pi {
	return global
}

func (p *Pi) DB(ctx context.Context) *db.Conn {
	conn := p.database.New(ctx)
	conn.UpdateHook = p.GetUpdateHook(ctx)
	conn.DeleteHook = p.GetDeleteHook(ctx)
	conn.InsertHook = p.GetInsertHook(ctx)
	return conn
}

func (p *Pi) Etcd(ctx context.Context) *etcd.Etcd {
	return p.etcd
}

func (p *Pi) GlobalConfig() (globalCfg *config.GlobalConfig) {
	mutex.RLock()
	globalCfg = p.globalCfg
	mutex.RUnlock()
	return
}

func (p *Pi) SetGlobalCfg(ctx context.Context) error {
	etcdClient := p.Etcd(ctx)
	err := etcdClient.Dlock(ctx, DlockKey, func() error {
		globalConfig := config.EncodeGlobalConfig(*p.GlobalConfig())
		_, err := etcdClient.Put(ctx, GlobalConfigKey, globalConfig)
		if err != nil {
			return err
		}
		p.setGlobalCfg(p.GlobalConfig())
		return nil
	})
	return err
}

func (p *Pi) RegisterRuntimeProvider(provider, providerConfig string) error {
	ctx := context.Background()
	err := p.GlobalConfig().RegisterRuntimeProviderConfig(provider, providerConfig)
	if err != nil {
		return err
	}
	return p.SetGlobalCfg(ctx)
}

func (p *Pi) setGlobalCfg(globalCfg *config.GlobalConfig) {
	mutex.Lock()
	p.globalCfg = globalCfg
	for _, cb := range p.globalCfgWatcher {
		go cb(globalCfg)
	}
	mutex.Unlock()
}

func (p *Pi) ThreadWatchGlobalConfig(cb globalCfgWatcher) {
	p.globalCfgWatcher = append(p.globalCfgWatcher, cb)
}

func (p *Pi) watchGlobalCfg() *Pi {
	watcher := make(Watcher)

	go func() {
		err := WatchGlobalConfig(p.Etcd(nil), watcher)
		if err != nil {
			logger.Critical(nil, "failed to watch global config")
			panic(err)
		}
	}()

	globalCfg := <-watcher
	p.setGlobalCfg(globalCfg)
	logger.Debug(nil, "Pi got global config: [%+v]", p.globalCfg)

	go func() {
		for globalCfg := range watcher {
			p.setGlobalCfg(globalCfg)
			logger.Debug(nil, "Global config update to [%+v]", globalCfg)
		}
	}()

	return p
}

func (p *Pi) openDatabase() *Pi {
	if p.cfg.Mysql.Disable {
		return p
	}
	dbSession, err := db.OpenDatabase(p.cfg.Mysql)
	if err != nil {
		logger.Critical(nil, "failed to connect mysql")
		panic(err)
	}
	p.database = dbSession
	return p
}

func (p *Pi) openEtcd() *Pi {
	endpoints := strings.Split(p.cfg.Etcd.Endpoints, ",")
	e, err := etcd.Connect(endpoints, EtcdPrefix)
	if err != nil {
		logger.Critical(nil, "failed to connect etcd")
		panic(err)
	}
	p.etcd = e
	return p
}
