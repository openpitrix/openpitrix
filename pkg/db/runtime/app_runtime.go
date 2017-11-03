// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db_app_runtime

import (
	"database/sql"
	"log"
	"sync/atomic"
	"time"

	_ "github.com/go-sql-driver/mysql"
	context "golang.org/x/net/context"
	"gopkg.in/gorp.v2"

	"openpitrix.io/openpitrix/pkg/config"
)

const (
	AppRuntimeTableName = "runtime"
)

type AppRuntime struct {
	Id           string    `db:"id, size:50, primarykey"`
	Name         string    `db:"name, size:50"`
	Description  string    `db:"description, size:1000"`
	Url          string    `db:"url, size:255"`
	Created      time.Time `db:"created"`
	LastModified time.Time `db:"last_modified"`
}

type AppRuntimeDatabase struct {
	db               *sql.DB
	dbMap            *gorp.DbMap
	createTablesDone uint32
}

func OpenAppRuntimeDatabase(cfg *config.Database) (p *AppRuntimeDatabase, err error) {
	// https://github.com/go-sql-driver/mysql/issues/9
	db, err := sql.Open(cfg.Type, cfg.GetUrl()+"?parseTime=true")
	if err != nil {
		return nil, err
	}

	dialect := gorp.MySQLDialect{
		Encoding: cfg.Encoding,
		Engine:   cfg.Engine,
	}

	dbMap := &gorp.DbMap{Db: db, Dialect: dialect}
	dbMap.AddTableWithName(AppRuntime{}, AppRuntimeTableName)

	p = &AppRuntimeDatabase{
		db:    db,
		dbMap: dbMap,
	}

	p.initTables()
	return
}

func (p *AppRuntimeDatabase) Close() error {
	return p.db.Close()
}

func (p *AppRuntimeDatabase) GetAppRuntime(ctx context.Context, id string) (*AppRuntime, error) {
	p.initTables()
	if v, err := p.dbMap.Get(AppRuntime{}, id); err == nil {
		return v.(*AppRuntime), nil
	} else {
		return nil, err
	}
}

func (p *AppRuntimeDatabase) GetAppRuntimeList(ctx context.Context) (apps []AppRuntime, err error) {
	p.initTables()
	_, err = p.dbMap.Select(&apps, "select * from "+AppRuntimeTableName)
	return
}

func (p *AppRuntimeDatabase) CreateAppRuntime(ctx context.Context, app *AppRuntime) error {
	p.initTables()
	return p.dbMap.Insert(app)
}

func (p *AppRuntimeDatabase) UpdateAppRuntime(ctx context.Context, app *AppRuntime) error {
	p.initTables()
	_, err := p.dbMap.Update(app)
	return err
}

func (p *AppRuntimeDatabase) DeleteAppRuntime(ctx context.Context, id string) error {
	p.initTables()
	_, err := p.dbMap.Delete(&AppRuntime{Id: id})
	return err
}

func (p *AppRuntimeDatabase) TruncateTables() error {
	p.initTables()
	err := p.dbMap.TruncateTables()
	return err
}

func (p *AppRuntimeDatabase) initTables() {
	if atomic.LoadUint32(&p.createTablesDone) == 1 {
		return
	}
	if err := p.dbMap.CreateTablesIfNotExists(); err != nil {
		log.Printf("CreateTablesIfNotExists: %v", err)
		return
	}
	atomic.StoreUint32(&p.createTablesDone, 1)
}
