// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db_cluster

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
	ClusterTableName = "cluster"
)

type Cluster struct {
	Id               string    `db:"id, size:50, primarykey"`
	Name             string    `db:"name, size:50"`
	Description      string    `db:"description, size:1000"`
	AppId            string    `db:"app_id, size:50"`
	AppVersion       string    `db:"app_version, size:50"`
	Status           string    `db:"status, size:50"`
	TransitionStatus string    `db:"transition_status, size:50"`
	Created          time.Time `db:"created"`
	LastModified     time.Time `db:"last_modified"`
}

type ClusterDatabase struct {
	db               *sql.DB
	dbMap            *gorp.DbMap
	createTablesDone uint32
}

func OpenClusterDatabase(cfg *config.Database) (p *ClusterDatabase, err error) {
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
	dbMap.AddTableWithName(Cluster{}, ClusterTableName)

	p = &ClusterDatabase{
		db:    db,
		dbMap: dbMap,
	}

	p.initTables()
	return
}

func (p *ClusterDatabase) Close() error {
	err := p.db.Close()
	err = errors.WithStack(err)
	return err
}

func (p *ClusterDatabase) GetCluster(ctx context.Context, id string) (*Cluster, error) {
	p.initTables()
	if v, err := p.dbMap.Get(Cluster{}, id); err == nil && v != nil {
		return v.(*Cluster), nil
	} else {
		err = errors.WithStack(err)
		return nil, err
	}
}

func (p *ClusterDatabase) GetClusterList(ctx context.Context) (apps []Cluster, err error) {
	p.initTables()
	_, err = p.dbMap.Select(&apps, "select * from "+ClusterTableName)
	err = errors.WithStack(err)
	return
}

func (p *ClusterDatabase) CreateCluster(ctx context.Context, app *Cluster) error {
	p.initTables()
	err := p.dbMap.Insert(app)
	err = errors.WithStack(err)
	return err
}

func (p *ClusterDatabase) UpdateCluster(ctx context.Context, app *Cluster) error {
	p.initTables()
	_, err := p.dbMap.Update(app)
	err = errors.WithStack(err)
	return err
}

func (p *ClusterDatabase) DeleteCluster(ctx context.Context, id string) error {
	p.initTables()
	_, err := p.dbMap.Delete(&Cluster{Id: id})
	err = errors.WithStack(err)
	return err
}

func (p *ClusterDatabase) TruncateTables() error {
	p.initTables()
	err := p.dbMap.TruncateTables()
	err = errors.WithStack(err)
	return err
}

func (p *ClusterDatabase) initTables() {
	if atomic.LoadUint32(&p.createTablesDone) == 1 {
		return
	}
	if err := p.dbMap.CreateTablesIfNotExists(); err != nil {
		logger.Warningf("CreateTablesIfNotExists: %+v", err)
		return
	}
	atomic.StoreUint32(&p.createTablesDone, 1)
}
