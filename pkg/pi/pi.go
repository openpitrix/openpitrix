// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package pi

import (
	"strings"
	"sync"

	"openpitrix.io/openpitrix/pkg/config"
	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/etcd"
	"openpitrix.io/openpitrix/pkg/logger"
)

type Pi struct {
	cfg       *config.Config
	globalCfg *config.GlobalConfig
	Db        *db.Database
	Etcd      *etcd.Etcd
}

func NewPi(cfg *config.Config) *Pi {
	p := &Pi{cfg: cfg}
	p.openDatabase()
	p.openEtcd()
	p.watchGlobalCfg()
	return p
}

var mutex sync.RWMutex

func (p *Pi) GlobalConfig() (globalCfg *config.GlobalConfig) {
	mutex.RLock()
	globalCfg = p.globalCfg
	mutex.RUnlock()
	return
}

func (p *Pi) setGlobalCfg(globalCfg *config.GlobalConfig) {
	mutex.Lock()
	p.globalCfg = globalCfg
	mutex.Unlock()
}

func (p *Pi) watchGlobalCfg() *Pi {
	watcher := make(config.Watcher)

	go func() {
		err := config.WatchGlobalConfig(p.Etcd, watcher)
		if err != nil {
			logger.Fatalf("failed to watch global config")
			panic(err)
		}
	}()

	globalCfg := <-watcher
	p.setGlobalCfg(globalCfg)
	logger.Debugf("Pi got global config: [%+v]", p.globalCfg)

	go func() {
		for globalCfg := range watcher {
			p.setGlobalCfg(globalCfg)
			logger.Debugf("Got global config [%+v]", globalCfg)
		}
	}()

	return p
}

func (p *Pi) openDatabase() *Pi {
	dbSession, err := db.OpenDatabase(p.cfg.Mysql)
	if err != nil {
		logger.Fatalf("failed to connect mysql")
		panic(err)
	}
	p.Db = dbSession
	return p
}

func (p *Pi) openEtcd() *Pi {
	endpoints := strings.Split(p.cfg.Etcd.Endpoints, ",")
	e, err := etcd.Connect(endpoints, config.EtcdPrefix)
	if err != nil {
		logger.Fatalf("failed to connect etcd")
		panic(err)
	}
	p.Etcd = e
	return p
}
