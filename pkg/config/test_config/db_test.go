// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// +build db

package test_config

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/logger"
)

func TestInsertQuery_Exec(t *testing.T) {
	logger.SetLevelByString("debug")
	appId := "app-unittest1"
	appName := "app-unittest1-name"
	dbName := "app"
	tableName := "test_tb"
	dbSession := NewDbTestConfig(dbName).GetDatabaseConn()

	_, err := dbSession.Exec(fmt.Sprintf("create table %s (app_id VARCHAR(50) PRIMARY KEY, name VARCHAR(255), description VARCHAR(255), repo_id VARCHAR(255))", tableName))
	assert.Nil(t, err)

	_, err = dbSession.DeleteFrom(tableName).Where(db.And(db.Eq("app_id", appId), db.Eq("name", appName))).Exec()
	assert.Nil(t, err)

	_, err = dbSession.InsertInto(tableName).Columns("app_id", "name").Values(appId, appName).Exec()
	assert.Nil(t, err)

	_, err = dbSession.Upsert(tableName).
		Where("app_id", appId).
		Where("name", appName).
		Set("description", "app_unittest_desc").
		Set("repo_id", "app_unittest_repo").Exec()
	assert.Nil(t, err)

	_, err = dbSession.Exec(fmt.Sprintf("drop table %s", tableName))
	assert.Nil(t, err)
}
