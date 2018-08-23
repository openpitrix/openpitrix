// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db

import (
	"time"

	"github.com/gocraft/dbr"

	"openpitrix.io/openpitrix/pkg/config"
)

func OpenDatabase(cfg config.MysqlConfig) (*Database, error) {
	// https://github.com/go-sql-driver/mysql/issues/9
	conn, err := dbr.Open("mysql", cfg.GetUrl()+"?parseTime=1&multiStatements=1&charset=utf8mb4&collation=utf8mb4_unicode_ci", nil)
	if err != nil {
		return nil, err
	}
	conn.SetMaxIdleConns(100)
	conn.SetMaxOpenConns(100)
	conn.SetConnMaxLifetime(10 * time.Second)
	return &Database{
		Conn: conn,
	}, nil
}
