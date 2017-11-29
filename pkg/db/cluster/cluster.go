// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db_cluster

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"gopkg.in/gorp.v2"

	"openpitrix.io/openpitrix/pkg/config"
	"openpitrix.io/openpitrix/pkg/logger"
)

const (
	ClusterTableName     = "cluster"
	ClusterNodeTableName = "cluster_node"
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

type ClusterNode struct {
	Id               string    `db:"id, size:50, primarykey"`
	InstanceId       string    `db:"instance_id, size:50"`
	Name             string    `db:"name, size:50"`
	Description      string    `db:"description, size:1000"`
	ClusterId        string    `db:"app_id, size:50"`
	PrivateIp        string    `db:"app_version, size:50"`
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
	dbMap.AddTableWithName(ClusterNode{}, ClusterNodeTableName)

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

func (p *ClusterDatabase) GetClusters(ctx context.Context, ids string) (clusters []Cluster, err error) {
	p.initTables()
	parts := strings.Split(ids, ",")
	for _, id := range parts {
		if id == "cl-panic000" {
			panic(id)
		}
	}
	in_ids := fmt.Sprintf("'%s'", strings.Join(parts, "','"))
	_, err = p.dbMap.Select(&clusters, fmt.Sprintf("select * from %s where id in (%s)", ClusterTableName, in_ids))
	err = errors.WithStack(err)
	return
}

func (p *ClusterDatabase) GetClusterList(ctx context.Context) (clusters []Cluster, err error) {
	p.initTables()
	_, err = p.dbMap.Select(&clusters, "select * from "+ClusterTableName)
	err = errors.WithStack(err)
	return
}

func (p *ClusterDatabase) CreateCluster(ctx context.Context, cluster *Cluster) error {
	p.initTables()
	err := p.dbMap.Insert(cluster)
	err = errors.WithStack(err)
	return err
}

func (p *ClusterDatabase) UpdateCluster(ctx context.Context, cluster *Cluster) error {
	p.initTables()
	_, err := p.dbMap.Update(cluster)
	err = errors.WithStack(err)
	return err
}

func (p *ClusterDatabase) DeleteClusters(ctx context.Context, ids string) error {
	p.initTables()
	parts := strings.Split(ids, ",")
	clusters := make([]interface{}, len(parts))
	for i, id := range parts {
		clusters[i] = &Cluster{Id: id}
	}
	_, err := p.dbMap.Delete(clusters...)
	err = errors.WithStack(err)
	return err
}

func (p *ClusterDatabase) GetClusterNodes(ctx context.Context, ids string) (clusterNodes []ClusterNode, err error) {
	p.initTables()
	parts := strings.Split(ids, ",")
	in_ids := fmt.Sprintf("'%s'", strings.Join(parts, "','"))
	_, err = p.dbMap.Select(&clusterNodes, fmt.Sprintf("select * from %s where id in (%s)", ClusterNodeTableName, in_ids))
	err = errors.WithStack(err)
	return
}

func (p *ClusterDatabase) GetClusterNodeList(ctx context.Context) (clusterNodes []ClusterNode, err error) {
	p.initTables()
	_, err = p.dbMap.Select(&clusterNodes, "select * from "+ClusterNodeTableName)
	err = errors.WithStack(err)
	return
}

func (p *ClusterDatabase) CreateClusterNodes(ctx context.Context, clusterNodes []*ClusterNode) error {
	p.initTables()
	dst := make([]interface{}, len(clusterNodes))
	for i, v := range clusterNodes {
		dst[i] = v
	}
	err := p.dbMap.Insert(dst...)
	err = errors.WithStack(err)
	return err
}

func (p *ClusterDatabase) UpdateClusterNode(ctx context.Context, clusterNode *ClusterNode) error {
	p.initTables()
	_, err := p.dbMap.Update(clusterNode)
	err = errors.WithStack(err)
	return err
}

func (p *ClusterDatabase) DeleteClusterNodes(ctx context.Context, ids string) error {
	p.initTables()
	parts := strings.Split(ids, ",")
	clusterNodes := make([]interface{}, len(parts))
	for i, id := range parts {
		clusterNodes[i] = &ClusterNode{Id: id}
	}
	_, err := p.dbMap.Delete(clusterNodes...)
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
