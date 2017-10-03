// +-------------------------------------------------------------------------
// | Copyright (C) 2017 Yunify, Inc.
// +-------------------------------------------------------------------------
// | Licensed under the Apache License, Version 2.0 (the "License");
// | you may not use this work except in compliance with the License.
// | You may obtain a copy of the License in the LICENSE file, or at:
// |
// | http://www.apache.org/licenses/LICENSE-2.0
// |
// | Unless required by applicable law or agreed to in writing, software
// | distributed under the License is distributed on an "AS IS" BASIS,
// | WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// | See the License for the specific language governing permissions and
// | limitations under the License.
// +-------------------------------------------------------------------------

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
