// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db_cluster

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
	return p.db.Close()
}

func (p *ClusterDatabase) GetCluster(ctx context.Context, id string) (*Cluster, error) {
	p.initTables()
	if v, err := p.dbMap.Get(Cluster{}, id); err == nil {
		return v.(*Cluster), nil
	} else {
		return nil, err
	}
}

func (p *ClusterDatabase) GetClusterList(ctx context.Context) (apps []Cluster, err error) {
	p.initTables()
	_, err = p.dbMap.Select(&apps, "select * from "+ClusterTableName)
	return
}

func (p *ClusterDatabase) CreateCluster(ctx context.Context, app *Cluster) error {
	p.initTables()
	return p.dbMap.Insert(app)
}

func (p *ClusterDatabase) UpdateCluster(ctx context.Context, app *Cluster) error {
	p.initTables()
	_, err := p.dbMap.Update(app)
	return err
}

func (p *ClusterDatabase) DeleteCluster(ctx context.Context, id string) error {
	p.initTables()
	_, err := p.dbMap.Delete(&Cluster{Id: id})
	return err
}

func (p *ClusterDatabase) TruncateTables() error {
	p.initTables()
	err := p.dbMap.TruncateTables()
	return err
}

func (p *ClusterDatabase) initTables() {
	if atomic.LoadUint32(&p.createTablesDone) == 1 {
		return
	}
	if err := p.dbMap.CreateTablesIfNotExists(); err != nil {
		log.Printf("CreateTablesIfNotExists: %v", err)
		return
	}
	atomic.StoreUint32(&p.createTablesDone, 1)
}
