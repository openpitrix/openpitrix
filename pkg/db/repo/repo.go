// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db_repo

import (
	"database/sql"
	"fmt"
	"sync"
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
	cfg               config.Database
	db                *sql.DB
	dbMap             *gorp.DbMap
	createTablesDone  uint32
	createTablesMutex sync.Mutex
}

func OpenRepoDatabase(cfg *config.Database) (p *RepoDatabase, err error) {
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
	dbMap.AddTableWithName(Repo{}, RepoTableName)

	p = &RepoDatabase{
		cfg:   *cfg,
		db:    db,
		dbMap: dbMap,
	}

	p.initTables()
	return
}

func (p *RepoDatabase) Close() error {
	err := p.db.Close()
	err = errors.WithStack(err)
	return err
}

func (p *RepoDatabase) GetRepo(ctx context.Context, id string) (*Repo, error) {
	p.initTables()
	if v, err := p.dbMap.Get(Repo{}, id); err == nil && v != nil {
		return v.(*Repo), nil
	} else {
		err = errors.WithStack(err)
		return nil, err
	}
}

func (p *RepoDatabase) GetRepoList(ctx context.Context) (apps []Repo, err error) {
	p.initTables()
	_, err = p.dbMap.Select(&apps, "select * from "+RepoTableName)
	err = errors.WithStack(err)
	return
}

func (p *RepoDatabase) CreateRepo(ctx context.Context, app *Repo) error {
	p.initTables()
	err := p.dbMap.Insert(app)
	err = errors.WithStack(err)
	return err
}

func (p *RepoDatabase) UpdateRepo(ctx context.Context, app *Repo) error {
	p.initTables()
	_, err := p.dbMap.Update(app)
	err = errors.WithStack(err)
	return err
}

func (p *RepoDatabase) DeleteRepo(ctx context.Context, id string) error {
	p.initTables()
	_, err := p.dbMap.Delete(&Repo{Id: id})
	err = errors.WithStack(err)
	return err
}

func (p *RepoDatabase) TruncateTables() error {
	p.initTables()
	err := p.dbMap.TruncateTables()
	err = errors.WithStack(err)
	return err
}

func (p *RepoDatabase) initTables() {
	if atomic.LoadUint32(&p.createTablesDone) == 1 {
		return
	}

	// Slow-path.
	p.createTablesMutex.Lock()
	defer p.createTablesMutex.Unlock()

	if atomic.LoadUint32(&p.createTablesDone) == 1 {
		return
	}

	db, err := sql.Open(p.cfg.Type, p.cfg.GetServerAddr())
	if err != nil {
		logger.Warningf("sql.Open failed: %+v", err)
		return
	}
	defer db.Close()

	sqlCreateDbIfNotExists := fmt.Sprintf(
		"CREATE DATABASE IF NOT EXISTS %s DEFAULT CHARACTER SET utf8;",
		p.cfg.DbName,
	)
	if _, err = db.Exec(sqlCreateDbIfNotExists); err != nil {
		logger.Warningf("CREATE DATABASE failed: %+v", err)
		return
	}
	if err := p.dbMap.CreateTablesIfNotExists(); err != nil {
		logger.Warningf("CreateTablesIfNotExists: %+v", err)
		return
	}

	atomic.StoreUint32(&p.createTablesDone, 1)
}
