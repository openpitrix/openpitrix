// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package apps

import (
	"database/sql"
	"fmt"
	"log"
	"sync/atomic"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/gorp.v2"

	"openpitrix.io/openpitrix/pkg/config"
)

const APP_TABLE_NAME = "app"

var _ AppDatabaseInterface = (*AppDatabase)(nil)

type AppDatabaseInterface interface {
	GetApps() (items AppsItems, err error)
	CreateApp(app *AppsItem) error
	GetApp(id string) (item AppsItem, err error)
	DeleteApp(id string) error
	TruncateTables() error
}

type AppDatabase struct {
	Cfg *config.Config

	db    *sql.DB
	dbMap *gorp.DbMap

	createTablesDone uint32
}

func OpenAppDatabase(config *config.Config) (p *AppDatabase, err error) {
	// https://github.com/go-sql-driver/mysql/issues/9
	dbpath := config.Database.GetUrl() + "?parseTime=true"
	db, err := sql.Open(config.Database.Type, dbpath)
	if err != nil {
		return nil, err
	}

	dialect := gorp.MySQLDialect{
		Encoding: config.Database.Encoding,
		Engine:   config.Database.Engine,
	}

	dbMap := &gorp.DbMap{Db: db, Dialect: dialect}
	dbMap.AddTableWithName(AppsItem{}, APP_TABLE_NAME)

	p = &AppDatabase{
		Cfg:   config.Clone(),
		db:    db,
		dbMap: dbMap,
	}

	p.initTables()
	return
}

func (p *AppDatabase) Close() error {
	return p.db.Close()
}

func (p *AppDatabase) GetApps() (items AppsItems, err error) {
	p.initTables()
	_, err = p.dbMap.Select(&items, fmt.Sprintf("select * from %s", APP_TABLE_NAME))
	return
}
func (p *AppDatabase) CreateApp(app *AppsItem) error {
	p.initTables()
	return p.dbMap.Insert(app)
}
func (p *AppDatabase) GetApp(id string) (item AppsItem, err error) {
	p.initTables()
	err = p.dbMap.SelectOne(&item, fmt.Sprintf("select * from %s where id=?", APP_TABLE_NAME), id)
	return
}
func (p *AppDatabase) DeleteApp(id string) error {
	p.initTables()
	println("DeleteApp id:" + id)
	_, err := p.dbMap.Delete(&AppsItem{Id: id})
	if err != nil {
		println("err:" + err.Error())
	}
	return err
}
func (p *AppDatabase) TruncateTables() error {
	p.initTables()
	err := p.dbMap.TruncateTables()
	return err
}

func (p *AppDatabase) initTables() {
	if atomic.LoadUint32(&p.createTablesDone) == 1 {
		return
	}
	if err := p.dbMap.CreateTablesIfNotExists(); err != nil {
		log.Printf("CreateTablesIfNotExists: %v", err)
		return
	}
	atomic.StoreUint32(&p.createTablesDone, 1)
}
