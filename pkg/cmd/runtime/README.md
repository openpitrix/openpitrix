# runtime
`import "openpitrix.io/openpitrix/pkg/cmd/runtime"`

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
- [golang.org/x/net/context](https://godoc.org/golang.org/x/net/context)
- [google.golang.org/grpc](https://godoc.org/google.golang.org/grpc)
- [google.golang.org/grpc/codes](https://godoc.org/google.golang.org/grpc/codes)
- [google.golang.org/grpc/grpclog/glogger](https://godoc.org/google.golang.org/grpc/grpclog/glogger)
- [openpitrix.io/openpitrix/pkg/config](https://godoc.org/openpitrix.io/openpitrix/pkg/config)
- [openpitrix.io/openpitrix/pkg/db/runtime](https://godoc.org/openpitrix.io/openpitrix/pkg/db/runtime)
- [openpitrix.io/openpitrix/pkg/logger](https://godoc.org/openpitrix.io/openpitrix/pkg/logger)
- [openpitrix.io/openpitrix/pkg/service.pb](https://godoc.org/openpitrix.io/openpitrix/pkg/service.pb)
- [openpitrix.io/openpitrix/pkg/version](https://godoc.org/openpitrix.io/openpitrix/pkg/version)

## <a name="pkg-index">Index</a>
* [func Main(cfg \*config.Config)](#Main)
* [func RegisterRuntime(runtime RuntimeInterface)](#RegisterRuntime)
* [func To\_database\_AppRuntime(dst \*db.AppRuntime, src \*pb.AppRuntime) \*db.AppRuntime](#To_database_AppRuntime)
* [func To\_proto\_AppRuntime(dst \*pb.AppRuntime, src \*db.AppRuntime) \*pb.AppRuntime](#To_proto_AppRuntime)
* [func To\_proto\_AppRuntimeList(p []db.AppRuntime, pageNumber, pageSize int) []\*pb.AppRuntime](#To_proto_AppRuntimeList)
* [type AppRuntimeServer](#AppRuntimeServer)
  * [func NewAppRuntimeServer(cfg \*config.Database) \*AppRuntimeServer](#NewAppRuntimeServer)
  * [func (p \*AppRuntimeServer) CreateAppRuntime(ctx context.Context, args \*pb.AppRuntime) (reply \*pbempty.Empty, err error)](#AppRuntimeServer.CreateAppRuntime)
  * [func (p \*AppRuntimeServer) DeleteAppRuntime(ctx context.Context, args \*pb.AppRuntimeId) (reply \*pbempty.Empty, err error)](#AppRuntimeServer.DeleteAppRuntime)
  * [func (p \*AppRuntimeServer) GetAppRuntime(ctx context.Context, args \*pb.AppRuntimeId) (reply \*pb.AppRuntime, err error)](#AppRuntimeServer.GetAppRuntime)
  * [func (p \*AppRuntimeServer) GetAppRuntimeList(ctx context.Context, args \*pb.AppRuntimeListRequest) (reply \*pb.AppRuntimeListResponse, err error)](#AppRuntimeServer.GetAppRuntimeList)
  * [func (p \*AppRuntimeServer) UpdateAppRuntime(ctx context.Context, args \*pb.AppRuntime) (reply \*pbempty.Empty, err error)](#AppRuntimeServer.UpdateAppRuntime)
* [type RuntimeInterface](#RuntimeInterface)

#### <a name="pkg-files">Package files</a>
[main.go](./main.go) [plugin.go](./plugin.go) [types.go](./types.go) 

## <a name="Main">func</a> [Main](./main.go#L29)
``` go
func Main(cfg *config.Config)
```

## <a name="RegisterRuntime">func</a> [RegisterRuntime](./plugin.go#L15)
``` go
func RegisterRuntime(runtime RuntimeInterface)
```

## <a name="To_database_AppRuntime">func</a> [To_database_AppRuntime](./types.go#L15)
``` go
func To_database_AppRuntime(dst *db.AppRuntime, src *pb.AppRuntime) *db.AppRuntime
```

## <a name="To_proto_AppRuntime">func</a> [To_proto_AppRuntime](./types.go#L31)
``` go
func To_proto_AppRuntime(dst *pb.AppRuntime, src *db.AppRuntime) *pb.AppRuntime
```

## <a name="To_proto_AppRuntimeList">func</a> [To_proto_AppRuntimeList](./types.go#L47)
``` go
func To_proto_AppRuntimeList(p []db.AppRuntime, pageNumber, pageSize int) []*pb.AppRuntime
```

## <a name="AppRuntimeServer">type</a> [AppRuntimeServer](./main.go#L72-L74)
``` go
type AppRuntimeServer struct {
    // contains filtered or unexported fields
}
```

### <a name="NewAppRuntimeServer">func</a> [NewAppRuntimeServer](./main.go#L76)
``` go
func NewAppRuntimeServer(cfg *config.Database) *AppRuntimeServer
```

### <a name="AppRuntimeServer.CreateAppRuntime">func</a> (\*AppRuntimeServer) [CreateAppRuntime](./main.go#L121)
``` go
func (p *AppRuntimeServer) CreateAppRuntime(ctx context.Context, args *pb.AppRuntime) (reply *pbempty.Empty, err error)
```

### <a name="AppRuntimeServer.DeleteAppRuntime">func</a> (\*AppRuntimeServer) [DeleteAppRuntime](./main.go#L141)
``` go
func (p *AppRuntimeServer) DeleteAppRuntime(ctx context.Context, args *pb.AppRuntimeId) (reply *pbempty.Empty, err error)
```

### <a name="AppRuntimeServer.GetAppRuntime">func</a> (\*AppRuntimeServer) [GetAppRuntime](./main.go#L87)
``` go
func (p *AppRuntimeServer) GetAppRuntime(ctx context.Context, args *pb.AppRuntimeId) (reply *pb.AppRuntime, err error)
```

### <a name="AppRuntimeServer.GetAppRuntimeList">func</a> (\*AppRuntimeServer) [GetAppRuntimeList](./main.go#L103)
``` go
func (p *AppRuntimeServer) GetAppRuntimeList(ctx context.Context, args *pb.AppRuntimeListRequest) (reply *pb.AppRuntimeListResponse, err error)
```

### <a name="AppRuntimeServer.UpdateAppRuntime">func</a> (\*AppRuntimeServer) [UpdateAppRuntime](./main.go#L131)
``` go
func (p *AppRuntimeServer) UpdateAppRuntime(ctx context.Context, args *pb.AppRuntime) (reply *pbempty.Empty, err error)
```

## <a name="RuntimeInterface">type</a> [RuntimeInterface](./plugin.go#L9-L13)
``` go
type RuntimeInterface interface {
    Name() string

    Run(app string, args ...string) error
}
```

- - -
Generated by [godoc2ghmd](https://github.com/GandalfUK/godoc2ghmd)