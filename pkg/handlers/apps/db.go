// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package apps

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/gorp.v2"
)

var _ AppDatabaseInterface = (*AppDatabase)(nil)

type AppDatabaseInterface interface {
	GetApps() (items AppsItems, err error)
	CreateApp(app *AppsItem) error
	GetApp(id string) (item AppsItem, err error)
	DeleteApp(id string) error
}

const APP_TABLE_NAME = "app"

type DbOptions struct {
	MySQLDialect *gorp.MySQLDialect
}

type AppDatabase struct {
	db    *sql.DB
	dbMap *gorp.DbMap
}

func OpenAppDatabase(dbname string, opt *DbOptions) (p *AppDatabase, err error) {
	db, err := sql.Open("mysql", dbname)
	if err != nil {
		return nil, err
	}

	dialect := gorp.MySQLDialect{
		Encoding: "utf8",
	}
	if opt != nil && opt.MySQLDialect != nil {
		dialect = *opt.MySQLDialect
	}

	dbMap := &gorp.DbMap{Db: db, Dialect: dialect}
	dbMap.AddTableWithName(AppsItem{}, APP_TABLE_NAME)

	err = dbMap.CreateTablesIfNotExists()
	if err != nil {
		db.Close()
		return nil, err
	}

	p = &AppDatabase{
		db:    db,
		dbMap: dbMap,
	}
	return
}

func (p *AppDatabase) Close() error {
	return p.db.Close()
}

func (p *AppDatabase) GetApps() (items AppsItems, err error) {
	_, err = p.dbMap.Select(&items, fmt.Sprintf("select * from %s", APP_TABLE_NAME))
	return
}
func (p *AppDatabase) CreateApp(app *AppsItem) error {
	return p.dbMap.Insert(app)
}
func (p *AppDatabase) GetApp(id string) (item AppsItem, err error) {
	err = p.dbMap.SelectOne(&item, fmt.Sprintf("select * from %s where id=?", APP_TABLE_NAME), id)
	return
}
func (p *AppDatabase) DeleteApp(id string) error {
	_, err := p.dbMap.Delete(&AppsItem{Id: id})
	return err
}
