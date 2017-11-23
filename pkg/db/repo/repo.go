// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db_repo

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
	RepoTableName = "repo"
)

type Repo struct {
	Id           string    `db:"id, size:50, primarykey"`
	Name         string    `db:"name, size:50"`
	Description  string    `db:"description, size:1000"`
	Url          string    `db:"url, size:255"`
	Created      time.Time `db:"created"`
	LastModified time.Time `db:"last_modified"`
}

type RepoDatabase struct {
	db               *sql.DB
	dbMap            *gorp.DbMap
	createTablesDone uint32
}

func OpenRepoDatabase(cfg *config.Database) (p *RepoDatabase, err error) {
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
	dbMap.AddTableWithName(Repo{}, RepoTableName)

	p = &RepoDatabase{
		db:    db,
		dbMap: dbMap,
	}

	p.initTables()
	return
}

func (p *RepoDatabase) Close() error {
	return p.db.Close()
}

func (p *RepoDatabase) GetRepo(ctx context.Context, id string) (*Repo, error) {
	p.initTables()
	if v, err := p.dbMap.Get(Repo{}, id); err == nil && v != nil {
		return v.(*Repo), nil
	} else {
		return nil, err
	}
}

func (p *RepoDatabase) GetRepoList(ctx context.Context) (apps []Repo, err error) {
	p.initTables()
	_, err = p.dbMap.Select(&apps, "select * from "+RepoTableName)
	return
}

func (p *RepoDatabase) CreateRepo(ctx context.Context, app *Repo) error {
	p.initTables()
	return p.dbMap.Insert(app)
}

func (p *RepoDatabase) UpdateRepo(ctx context.Context, app *Repo) error {
	p.initTables()
	_, err := p.dbMap.Update(app)
	return err
}

func (p *RepoDatabase) DeleteRepo(ctx context.Context, id string) error {
	p.initTables()
	_, err := p.dbMap.Delete(&Repo{Id: id})
	return err
}

func (p *RepoDatabase) TruncateTables() error {
	p.initTables()
	err := p.dbMap.TruncateTables()
	return err
}

func (p *RepoDatabase) initTables() {
	if atomic.LoadUint32(&p.createTablesDone) == 1 {
		return
	}
	if err := p.dbMap.CreateTablesIfNotExists(); err != nil {
		log.Printf("CreateTablesIfNotExists: %v", err)
		return
	}
	atomic.StoreUint32(&p.createTablesDone, 1)
}
