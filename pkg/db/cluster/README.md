# db_cluster
`import "openpitrix.io/openpitrix/pkg/db/cluster"`

* [Overview](#pkg-overview)
* [Imported Packages](#pkg-imports)
* [Index](#pkg-index)

## <a name="pkg-overview">Overview</a>

## <a name="pkg-imports">Imported Packages</a>

- [github.com/go-sql-driver/mysql](https://godoc.org/github.com/go-sql-driver/mysql)
- [github.com/pkg/errors](https://godoc.org/github.com/pkg/errors)
- [gopkg.in/gorp.v2](https://godoc.org/gopkg.in/gorp.v2)
- [openpitrix.io/openpitrix/pkg/config](https://godoc.org/openpitrix.io/openpitrix/pkg/config)
- [openpitrix.io/openpitrix/pkg/logger](https://godoc.org/openpitrix.io/openpitrix/pkg/logger)

## <a name="pkg-index">Index</a>
* [Constants](#pkg-constants)
* [type Cluster](#Cluster)
* [type ClusterDatabase](#ClusterDatabase)
  * [func OpenClusterDatabase(cfg \*config.Database) (p \*ClusterDatabase, err error)](#OpenClusterDatabase)
  * [func (p \*ClusterDatabase) Close() error](#ClusterDatabase.Close)
  * [func (p \*ClusterDatabase) CreateCluster(ctx context.Context, cluster \*Cluster) error](#ClusterDatabase.CreateCluster)
  * [func (p \*ClusterDatabase) CreateClusterNodes(ctx context.Context, clusterNodes []\*ClusterNode) error](#ClusterDatabase.CreateClusterNodes)
  * [func (p \*ClusterDatabase) DeleteClusterNodes(ctx context.Context, ids string) error](#ClusterDatabase.DeleteClusterNodes)
  * [func (p \*ClusterDatabase) DeleteClusters(ctx context.Context, ids string) error](#ClusterDatabase.DeleteClusters)
  * [func (p \*ClusterDatabase) GetClusterList(ctx context.Context) (clusters []Cluster, err error)](#ClusterDatabase.GetClusterList)
  * [func (p \*ClusterDatabase) GetClusterNodeList(ctx context.Context) (clusterNodes []ClusterNode, err error)](#ClusterDatabase.GetClusterNodeList)
  * [func (p \*ClusterDatabase) GetClusterNodes(ctx context.Context, ids string) (clusterNodes []ClusterNode, err error)](#ClusterDatabase.GetClusterNodes)
  * [func (p \*ClusterDatabase) GetClusters(ctx context.Context, ids string) (clusters []Cluster, err error)](#ClusterDatabase.GetClusters)
  * [func (p \*ClusterDatabase) TruncateTables() error](#ClusterDatabase.TruncateTables)
  * [func (p \*ClusterDatabase) UpdateCluster(ctx context.Context, cluster \*Cluster) error](#ClusterDatabase.UpdateCluster)
  * [func (p \*ClusterDatabase) UpdateClusterNode(ctx context.Context, clusterNode \*ClusterNode) error](#ClusterDatabase.UpdateClusterNode)
* [type ClusterNode](#ClusterNode)

#### <a name="pkg-files">Package files</a>
[cluster.go](./cluster.go) 

## <a name="pkg-constants">Constants</a>
``` go
const (
    ClusterTableName     = "cluster"
    ClusterNodeTableName = "cluster_node"
)
```

## <a name="Cluster">type</a> [Cluster](./cluster.go#L28-L38)
``` go
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
```

## <a name="ClusterDatabase">type</a> [ClusterDatabase](./cluster.go#L53-L57)
``` go
type ClusterDatabase struct {
    // contains filtered or unexported fields
}
```

### <a name="OpenClusterDatabase">func</a> [OpenClusterDatabase](./cluster.go#L59)
``` go
func OpenClusterDatabase(cfg *config.Database) (p *ClusterDatabase, err error)
```

### <a name="ClusterDatabase.Close">func</a> (\*ClusterDatabase) [Close](./cluster.go#L85)
``` go
func (p *ClusterDatabase) Close() error
```

### <a name="ClusterDatabase.CreateCluster">func</a> (\*ClusterDatabase) [CreateCluster](./cluster.go#L112)
``` go
func (p *ClusterDatabase) CreateCluster(ctx context.Context, cluster *Cluster) error
```

### <a name="ClusterDatabase.CreateClusterNodes">func</a> (\*ClusterDatabase) [CreateClusterNodes](./cluster.go#L154)
``` go
func (p *ClusterDatabase) CreateClusterNodes(ctx context.Context, clusterNodes []*ClusterNode) error
```

### <a name="ClusterDatabase.DeleteClusterNodes">func</a> (\*ClusterDatabase) [DeleteClusterNodes](./cluster.go#L172)
``` go
func (p *ClusterDatabase) DeleteClusterNodes(ctx context.Context, ids string) error
```

### <a name="ClusterDatabase.DeleteClusters">func</a> (\*ClusterDatabase) [DeleteClusters](./cluster.go#L126)
``` go
func (p *ClusterDatabase) DeleteClusters(ctx context.Context, ids string) error
```

### <a name="ClusterDatabase.GetClusterList">func</a> (\*ClusterDatabase) [GetClusterList](./cluster.go#L105)
``` go
func (p *ClusterDatabase) GetClusterList(ctx context.Context) (clusters []Cluster, err error)
```

### <a name="ClusterDatabase.GetClusterNodeList">func</a> (\*ClusterDatabase) [GetClusterNodeList](./cluster.go#L147)
``` go
func (p *ClusterDatabase) GetClusterNodeList(ctx context.Context) (clusterNodes []ClusterNode, err error)
```

### <a name="ClusterDatabase.GetClusterNodes">func</a> (\*ClusterDatabase) [GetClusterNodes](./cluster.go#L138)
``` go
func (p *ClusterDatabase) GetClusterNodes(ctx context.Context, ids string) (clusterNodes []ClusterNode, err error)
```

### <a name="ClusterDatabase.GetClusters">func</a> (\*ClusterDatabase) [GetClusters](./cluster.go#L91)
``` go
func (p *ClusterDatabase) GetClusters(ctx context.Context, ids string) (clusters []Cluster, err error)
```

### <a name="ClusterDatabase.TruncateTables">func</a> (\*ClusterDatabase) [TruncateTables](./cluster.go#L184)
``` go
func (p *ClusterDatabase) TruncateTables() error
```

### <a name="ClusterDatabase.UpdateCluster">func</a> (\*ClusterDatabase) [UpdateCluster](./cluster.go#L119)
``` go
func (p *ClusterDatabase) UpdateCluster(ctx context.Context, cluster *Cluster) error
```

### <a name="ClusterDatabase.UpdateClusterNode">func</a> (\*ClusterDatabase) [UpdateClusterNode](./cluster.go#L165)
``` go
func (p *ClusterDatabase) UpdateClusterNode(ctx context.Context, clusterNode *ClusterNode) error
```

## <a name="ClusterNode">type</a> [ClusterNode](./cluster.go#L40-L51)
``` go
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
```

- - -
Generated by [godoc2ghmd](https://github.com/GandalfUK/godoc2ghmd)