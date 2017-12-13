# cluster
`import "openpitrix.io/openpitrix/pkg/cmd/cluster"`

* [Overview](#pkg-overview)
* [Imported Packages](#pkg-imports)
* [Index](#pkg-index)

## <a name="pkg-overview">Overview</a>

## <a name="pkg-imports">Imported Packages</a>

- [github.com/golang/protobuf/proto](https://godoc.org/github.com/golang/protobuf/proto)
- [github.com/golang/protobuf/ptypes](https://godoc.org/github.com/golang/protobuf/ptypes)
- [github.com/golang/protobuf/ptypes/empty](https://godoc.org/github.com/golang/protobuf/ptypes/empty)
- [github.com/grpc-ecosystem/go-grpc-middleware](https://godoc.org/github.com/grpc-ecosystem/go-grpc-middleware)
- [github.com/grpc-ecosystem/go-grpc-middleware/recovery](https://godoc.org/github.com/grpc-ecosystem/go-grpc-middleware/recovery)
- [github.com/grpc-ecosystem/go-grpc-middleware/validator](https://godoc.org/github.com/grpc-ecosystem/go-grpc-middleware/validator)
- [github.com/pkg/errors](https://godoc.org/github.com/pkg/errors)
- [google.golang.org/grpc](https://godoc.org/google.golang.org/grpc)
- [google.golang.org/grpc/codes](https://godoc.org/google.golang.org/grpc/codes)
- [google.golang.org/grpc/grpclog/glogger](https://godoc.org/google.golang.org/grpc/grpclog/glogger)
- [openpitrix.io/openpitrix/pkg/config](https://godoc.org/openpitrix.io/openpitrix/pkg/config)
- [openpitrix.io/openpitrix/pkg/db/cluster](https://godoc.org/openpitrix.io/openpitrix/pkg/db/cluster)
- [openpitrix.io/openpitrix/pkg/logger](https://godoc.org/openpitrix.io/openpitrix/pkg/logger)
- [openpitrix.io/openpitrix/pkg/service.pb](https://godoc.org/openpitrix.io/openpitrix/pkg/service.pb)
- [openpitrix.io/openpitrix/pkg/version](https://godoc.org/openpitrix.io/openpitrix/pkg/version)

## <a name="pkg-index">Index</a>
* [func Main(cfg \*config.Config)](#Main)
* [func To\_database\_Cluster(dst \*db.Cluster, src \*pb.Cluster) \*db.Cluster](#To_database_Cluster)
* [func To\_database\_ClusterNode(dst \*db.ClusterNode, src \*pb.ClusterNode) \*db.ClusterNode](#To_database_ClusterNode)
* [func To\_database\_ClusterNodes(src \*pb.ClusterNodes) []\*db.ClusterNode](#To_database_ClusterNodes)
* [func To\_database\_Clusters(src \*pb.Clusters) []\*db.Cluster](#To_database_Clusters)
* [func To\_proto\_Cluster(dst \*pb.Cluster, src \*db.Cluster) \*pb.Cluster](#To_proto_Cluster)
* [func To\_proto\_ClusterList(p []db.Cluster, pageNumber, pageSize int) []\*pb.Cluster](#To_proto_ClusterList)
* [func To\_proto\_ClusterNode(dst \*pb.ClusterNode, src \*db.ClusterNode) \*pb.ClusterNode](#To_proto_ClusterNode)
* [func To\_proto\_ClusterNodeList(p []db.ClusterNode, pageNumber, pageSize int) []\*pb.ClusterNode](#To_proto_ClusterNodeList)
* [func To\_proto\_ClusterNodes(src []db.ClusterNode) \*pb.ClusterNodes](#To_proto_ClusterNodes)
* [func To\_proto\_Clusters(src []db.Cluster) \*pb.Clusters](#To_proto_Clusters)
* [type ClusterServer](#ClusterServer)
  * [func NewClusterServer(cfg \*config.Database) \*ClusterServer](#NewClusterServer)
  * [func (p \*ClusterServer) CreateCluster(ctx context.Context, args \*pb.Cluster) (reply \*pbempty.Empty, err error)](#ClusterServer.CreateCluster)
  * [func (p \*ClusterServer) CreateClusterNodes(ctx context.Context, args \*pb.ClusterNodes) (reply \*pbempty.Empty, err error)](#ClusterServer.CreateClusterNodes)
  * [func (p \*ClusterServer) DeleteClusterNodes(ctx context.Context, args \*pb.ClusterNodeIds) (reply \*pbempty.Empty, err error)](#ClusterServer.DeleteClusterNodes)
  * [func (p \*ClusterServer) DeleteClusters(ctx context.Context, args \*pb.ClusterIds) (reply \*pbempty.Empty, err error)](#ClusterServer.DeleteClusters)
  * [func (p \*ClusterServer) GetClusterList(ctx context.Context, args \*pb.ClusterListRequest) (reply \*pb.ClusterListResponse, err error)](#ClusterServer.GetClusterList)
  * [func (p \*ClusterServer) GetClusterNodeList(ctx context.Context, args \*pb.ClusterNodeListRequest) (reply \*pb.ClusterNodeListResponse, err error)](#ClusterServer.GetClusterNodeList)
  * [func (p \*ClusterServer) GetClusterNodes(ctx context.Context, args \*pb.ClusterNodeIds) (reply \*pb.ClusterNodes, err error)](#ClusterServer.GetClusterNodes)
  * [func (p \*ClusterServer) GetClusters(ctx context.Context, args \*pb.ClusterIds) (reply \*pb.Clusters, err error)](#ClusterServer.GetClusters)
  * [func (p \*ClusterServer) UpdateCluster(ctx context.Context, args \*pb.Cluster) (reply \*pbempty.Empty, err error)](#ClusterServer.UpdateCluster)
  * [func (p \*ClusterServer) UpdateClusterNode(ctx context.Context, args \*pb.ClusterNode) (reply \*pbempty.Empty, err error)](#ClusterServer.UpdateClusterNode)

#### <a name="pkg-files">Package files</a>
[main.go](./main.go) [types.go](./types.go) 

## <a name="Main">func</a> [Main](./main.go#L29)
``` go
func Main(cfg *config.Config)
```

## <a name="To_database_Cluster">func</a> [To_database_Cluster](./types.go#L15)
``` go
func To_database_Cluster(dst *db.Cluster, src *pb.Cluster) *db.Cluster
```

## <a name="To_database_ClusterNode">func</a> [To_database_ClusterNode](./types.go#L94)
``` go
func To_database_ClusterNode(dst *db.ClusterNode, src *pb.ClusterNode) *db.ClusterNode
```

## <a name="To_database_ClusterNodes">func</a> [To_database_ClusterNodes](./types.go#L167)
``` go
func To_database_ClusterNodes(src *pb.ClusterNodes) []*db.ClusterNode
```

## <a name="To_database_Clusters">func</a> [To_database_Clusters](./types.go#L86)
``` go
func To_database_Clusters(src *pb.Clusters) []*db.Cluster
```

## <a name="To_proto_Cluster">func</a> [To_proto_Cluster](./types.go#L34)
``` go
func To_proto_Cluster(dst *pb.Cluster, src *db.Cluster) *pb.Cluster
```

## <a name="To_proto_ClusterList">func</a> [To_proto_ClusterList](./types.go#L53)
``` go
func To_proto_ClusterList(p []db.Cluster, pageNumber, pageSize int) []*pb.Cluster
```

## <a name="To_proto_ClusterNode">func</a> [To_proto_ClusterNode](./types.go#L114)
``` go
func To_proto_ClusterNode(dst *pb.ClusterNode, src *db.ClusterNode) *pb.ClusterNode
```

## <a name="To_proto_ClusterNodeList">func</a> [To_proto_ClusterNodeList](./types.go#L134)
``` go
func To_proto_ClusterNodeList(p []db.ClusterNode, pageNumber, pageSize int) []*pb.ClusterNode
```

## <a name="To_proto_ClusterNodes">func</a> [To_proto_ClusterNodes](./types.go#L156)
``` go
func To_proto_ClusterNodes(src []db.ClusterNode) *pb.ClusterNodes
```

## <a name="To_proto_Clusters">func</a> [To_proto_Clusters](./types.go#L75)
``` go
func To_proto_Clusters(src []db.Cluster) *pb.Clusters
```

## <a name="ClusterServer">type</a> [ClusterServer](./main.go#L72-L74)
``` go
type ClusterServer struct {
    // contains filtered or unexported fields
}
```

### <a name="NewClusterServer">func</a> [NewClusterServer](./main.go#L76)
``` go
func NewClusterServer(cfg *config.Database) *ClusterServer
```

### <a name="ClusterServer.CreateCluster">func</a> (\*ClusterServer) [CreateCluster](./main.go#L119)
``` go
func (p *ClusterServer) CreateCluster(ctx context.Context, args *pb.Cluster) (reply *pbempty.Empty, err error)
```

### <a name="ClusterServer.CreateClusterNodes">func</a> (\*ClusterServer) [CreateClusterNodes](./main.go#L181)
``` go
func (p *ClusterServer) CreateClusterNodes(ctx context.Context, args *pb.ClusterNodes) (reply *pbempty.Empty, err error)
```

### <a name="ClusterServer.DeleteClusterNodes">func</a> (\*ClusterServer) [DeleteClusterNodes](./main.go#L201)
``` go
func (p *ClusterServer) DeleteClusterNodes(ctx context.Context, args *pb.ClusterNodeIds) (reply *pbempty.Empty, err error)
```

### <a name="ClusterServer.DeleteClusters">func</a> (\*ClusterServer) [DeleteClusters](./main.go#L139)
``` go
func (p *ClusterServer) DeleteClusters(ctx context.Context, args *pb.ClusterIds) (reply *pbempty.Empty, err error)
```

### <a name="ClusterServer.GetClusterList">func</a> (\*ClusterServer) [GetClusterList](./main.go#L101)
``` go
func (p *ClusterServer) GetClusterList(ctx context.Context, args *pb.ClusterListRequest) (reply *pb.ClusterListResponse, err error)
```

### <a name="ClusterServer.GetClusterNodeList">func</a> (\*ClusterServer) [GetClusterNodeList](./main.go#L163)
``` go
func (p *ClusterServer) GetClusterNodeList(ctx context.Context, args *pb.ClusterNodeListRequest) (reply *pb.ClusterNodeListResponse, err error)
```

### <a name="ClusterServer.GetClusterNodes">func</a> (\*ClusterServer) [GetClusterNodes](./main.go#L149)
``` go
func (p *ClusterServer) GetClusterNodes(ctx context.Context, args *pb.ClusterNodeIds) (reply *pb.ClusterNodes, err error)
```

### <a name="ClusterServer.GetClusters">func</a> (\*ClusterServer) [GetClusters](./main.go#L87)
``` go
func (p *ClusterServer) GetClusters(ctx context.Context, args *pb.ClusterIds) (reply *pb.Clusters, err error)
```

### <a name="ClusterServer.UpdateCluster">func</a> (\*ClusterServer) [UpdateCluster](./main.go#L129)
``` go
func (p *ClusterServer) UpdateCluster(ctx context.Context, args *pb.Cluster) (reply *pbempty.Empty, err error)
```

### <a name="ClusterServer.UpdateClusterNode">func</a> (\*ClusterServer) [UpdateClusterNode](./main.go#L191)
``` go
func (p *ClusterServer) UpdateClusterNode(ctx context.Context, args *pb.ClusterNode) (reply *pbempty.Empty, err error)
```

- - -
Generated by [godoc2ghmd](https://github.com/GandalfUK/godoc2ghmd)