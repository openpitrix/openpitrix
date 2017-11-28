// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db_app

import (
	"database/sql"
	"sync/atomic"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	context "golang.org/x/net/context"
	"gopkg.in/gorp.v2"

	"openpitrix.io/openpitrix/pkg/config"
	"openpitrix.io/openpitrix/pkg/logger"
)

const (
	AppTableName = "app"
)

type App struct {
	Id           string    `db:"id, size:50, primarykey"`
	Name         string    `db:"name, size:50"`
	Description  string    `db:"description, size:1000"`
	RepoId       string    `db:"repo_id, size:50"`
	Created      time.Time `db:"created"`
	LastModified time.Time `db:"last_modified"`
}

type AppDatabase struct {
	db               *sql.DB
	dbMap            *gorp.DbMap
	createTablesDone uint32
}

func OpenAppDatabase(cfg *config.Database) (p *AppDatabase, err error) {
	// https://github.com/go-sql-driver/mysql/issues/9
	db, err := sql.Open(cfg.Type, cfg.GetUrl()+"?parseTime=true")
	if err != nil {
		err = errors.WithStack(err)
		return nil, err
	}

	dialect := gorp.MySQLDialect{
		Encoding: cfg.Encoding,
		Engine:   cfg.Engine,
	}

	dbMap := &gorp.DbMap{Db: db, Dialect: dialect}
	dbMap.AddTableWithName(App{}, AppTableName)

	p = &AppDatabase{
		db:    db,
		dbMap: dbMap,
	}

	p.initTables()
	return
}

func (p *AppDatabase) Close() error {
	err := p.db.Close()
	err = errors.WithStack(err)
	return err
}

func (p *AppDatabase) GetApp(ctx context.Context, id string) (*App, error) {
	p.initTables()
	if v, err := p.dbMap.Get(App{}, id); err == nil && v != nil {
		return v.(*App), nil
	} else {
		err = errors.WithStack(err)
		return nil, err
	}
}

func (p *AppDatabase) GetAppList(ctx context.Context) (apps []App, err error) {
	p.initTables()
	_, err = p.dbMap.Select(&apps, "select * from "+AppTableName)
	err = errors.WithStack(err)
	return
}

func (p *AppDatabase) CreateApp(ctx context.Context, app *App) error {
	p.initTables()
	err := p.dbMap.Insert(app)
	err = errors.WithStack(err)
	return err
}

func (p *AppDatabase) UpdateApp(ctx context.Context, app *App) error {
	p.initTables()
	_, err := p.dbMap.Update(app)
	err = errors.WithStack(err)
	return err
}

func (p *AppDatabase) DeleteApp(ctx context.Context, id string) error {
	p.initTables()
	_, err := p.dbMap.Delete(&App{Id: id})
	err = errors.WithStack(err)
	return err
}

func (p *AppDatabase) TruncateTables() error {
	p.initTables()
	err := p.dbMap.TruncateTables()
	err = errors.WithStack(err)
	return err
}

func (p *AppDatabase) initTables() {
	if atomic.LoadUint32(&p.createTablesDone) == 1 {
		return
	}
	if err := p.dbMap.CreateTablesIfNotExists(); err != nil {
		logger.Warnf("CreateTablesIfNotExists: %+v", err)
		return
	}
	atomic.StoreUint32(&p.createTablesDone, 1)
}
