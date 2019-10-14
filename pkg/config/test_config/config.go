// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package test_config

import (
	"os"
	"testing"

	"openpitrix.io/openpitrix/pkg/config"
	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/logger"
)

type DbTestConfig struct {
	openDbUnitTests string
	envConfig       *config.Config
}

var _ = func() bool {
	testing.Init()
	return true
}()

func NewDbTestConfig(database string) DbTestConfig {
	tc := DbTestConfig{
		openDbUnitTests: os.Getenv("OP_DB_UNIT_TEST"),
		envConfig:       config.GetConf(),
	}
	tc.envConfig.Mysql.Database = database
	return tc
}

func (tc DbTestConfig) GetDatabaseConn() *db.Database {
	if tc.openDbUnitTests == "1" {
		d, err := db.OpenDatabase(tc.envConfig.Mysql)
		if err != nil {
			logger.Critical(nil, "failed to open database %+v", tc.envConfig.Mysql)
		}
		return d
	}
	return nil
}

func (tc DbTestConfig) CheckDbUnitTest(t *testing.T) {
	if tc.openDbUnitTests != "1" {
		t.Skipf("if you want run unit tests with db,set OP_DB_UNIT_TEST=1")
	}
}

type EtcdTestConfig struct {
	openEtcdUnitTests string
	envConfig         *config.Config
}

func NewEtcdTestConfig() EtcdTestConfig {
	tc := EtcdTestConfig{
		openEtcdUnitTests: os.Getenv("OP_ETCD_UNIT_TEST"),
		envConfig:         config.GetConf(),
	}
	return tc
}

func (tc EtcdTestConfig) GetTestEtcdEndpoints() []string {
	if tc.openEtcdUnitTests == "1" {
		return []string{tc.envConfig.Etcd.Endpoints}
	}
	return nil
}

func (tc EtcdTestConfig) CheckEtcdUnitTest(t *testing.T) {
	if tc.openEtcdUnitTests != "1" {
		t.Skipf("if you want run unit tests with db,set OP_ETCD_UNIT_TEST=1")
	}
}
