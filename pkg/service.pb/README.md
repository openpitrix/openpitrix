# openpitrix
`import "openpitrix.io/openpitrix/pkg/service.pb"`

* [Overview](#pkg-overview)
* [Imported Packages](#pkg-imports)
* [Index](#pkg-index)

## <a name="pkg-overview">Overview</a>
App
	AppId
	AppListRequest
	AppListResponse
	AppRuntime
	AppRuntimeLabel
	AppRuntimeId
	AppRuntimeListRequest
	AppRuntimeListResponse
	Cluster
	Clusters
	ClusterNode
	ClusterNodes
	ClusterId
	ClusterIds
	ClusterListRequest
	ClusterListResponse
	ClusterNodeId
	ClusterNodeIds
	ClusterNodeListRequest
	ClusterNodeListResponse
	Repo
	RepoLabel
	RepoSelector
	RepoId
	RepoListRequest
	RepoListResponse

Package openpitrix is a generated protocol buffer package.

It is generated from these files:

	annotations.proto
	app.proto
	app_runtime.proto
	cluster.proto
	repo.proto

It has these top-level messages:

	App
	AppId
	AppListRequest
	AppListResponse
	AppRuntime
	AppRuntimeLabel
	AppRuntimeId
	AppRuntimeListRequest
	AppRuntimeListResponse
	Cluster
	Clusters
	ClusterNode
	ClusterNodes
	ClusterId
	ClusterIds
	ClusterListRequest
	ClusterListResponse
	ClusterNodeId
	ClusterNodeIds
	ClusterNodeListRequest
	ClusterNodeListResponse
	Repo
	RepoLabel
	RepoSelector
	RepoId
	RepoListRequest
	RepoListResponse

Package openpitrix is a reverse proxy.

It translates gRPC into RESTful JSON APIs.

Package openpitrix is a reverse proxy.

It translates gRPC into RESTful JSON APIs.

Package openpitrix is a reverse proxy.

It translates gRPC into RESTful JSON APIs.

Package openpitrix is a reverse proxy.

It translates gRPC into RESTful JSON APIs.

**Package openpitrix is a generated [Protobuf](https://developers.google.com/protocol-buffers/)-compatible package.**

It is generated from these files:

- [annotations.proto](./annotations.proto)
- [app.proto](./app.proto)
- [app_runtime.proto](./app_runtime.proto)
- [cluster.proto](./cluster.proto)
- [repo.proto](./repo.proto)

## <a name="pkg-imports">Imported Packages</a>

- [github.com/golang/protobuf/proto](https://godoc.org/github.com/golang/protobuf/proto)
- [github.com/golang/protobuf/protoc-gen-go/descriptor](https://godoc.org/github.com/golang/protobuf/protoc-gen-go/descriptor)
- [github.com/golang/protobuf/ptypes/empty](https://godoc.org/github.com/golang/protobuf/ptypes/empty)
- [github.com/golang/protobuf/ptypes/timestamp](https://godoc.org/github.com/golang/protobuf/ptypes/timestamp)
- [github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger/options](https://godoc.org/github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger/options)
- [github.com/grpc-ecosystem/grpc-gateway/runtime](https://godoc.org/github.com/grpc-ecosystem/grpc-gateway/runtime)
- [github.com/grpc-ecosystem/grpc-gateway/utilities](https://godoc.org/github.com/grpc-ecosystem/grpc-gateway/utilities)
- [github.com/mwitkow/go-proto-validators](https://godoc.org/github.com/mwitkow/go-proto-validators)
- [golang.org/x/net/context](https://godoc.org/golang.org/x/net/context)
- [google.golang.org/genproto/googleapis/api/annotations](https://godoc.org/google.golang.org/genproto/googleapis/api/annotations)
- [google.golang.org/grpc](https://godoc.org/google.golang.org/grpc)
- [google.golang.org/grpc/codes](https://godoc.org/google.golang.org/grpc/codes)
- [google.golang.org/grpc/grpclog](https://godoc.org/google.golang.org/grpc/grpclog)
- [google.golang.org/grpc/status](https://godoc.org/google.golang.org/grpc/status)

## <a name="pkg-index">Index</a>
* [Constants](#pkg-constants)
* [Variables](#pkg-variables)
* [func RegisterAppRuntimeServiceHandler(ctx context.Context, mux \*runtime.ServeMux, conn \*grpc.ClientConn) error](#RegisterAppRuntimeServiceHandler)
* [func RegisterAppRuntimeServiceHandlerClient(ctx context.Context, mux \*runtime.ServeMux, client AppRuntimeServiceClient) error](#RegisterAppRuntimeServiceHandlerClient)
* [func RegisterAppRuntimeServiceHandlerFromEndpoint(ctx context.Context, mux \*runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error)](#RegisterAppRuntimeServiceHandlerFromEndpoint)
* [func RegisterAppRuntimeServiceServer(s \*grpc.Server, srv AppRuntimeServiceServer)](#RegisterAppRuntimeServiceServer)
* [func RegisterAppServiceHandler(ctx context.Context, mux \*runtime.ServeMux, conn \*grpc.ClientConn) error](#RegisterAppServiceHandler)
* [func RegisterAppServiceHandlerClient(ctx context.Context, mux \*runtime.ServeMux, client AppServiceClient) error](#RegisterAppServiceHandlerClient)
* [func RegisterAppServiceHandlerFromEndpoint(ctx context.Context, mux \*runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error)](#RegisterAppServiceHandlerFromEndpoint)
* [func RegisterAppServiceServer(s \*grpc.Server, srv AppServiceServer)](#RegisterAppServiceServer)
* [func RegisterClusterServiceHandler(ctx context.Context, mux \*runtime.ServeMux, conn \*grpc.ClientConn) error](#RegisterClusterServiceHandler)
* [func RegisterClusterServiceHandlerClient(ctx context.Context, mux \*runtime.ServeMux, client ClusterServiceClient) error](#RegisterClusterServiceHandlerClient)
* [func RegisterClusterServiceHandlerFromEndpoint(ctx context.Context, mux \*runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error)](#RegisterClusterServiceHandlerFromEndpoint)
* [func RegisterClusterServiceServer(s \*grpc.Server, srv ClusterServiceServer)](#RegisterClusterServiceServer)
* [func RegisterRepoServiceHandler(ctx context.Context, mux \*runtime.ServeMux, conn \*grpc.ClientConn) error](#RegisterRepoServiceHandler)
* [func RegisterRepoServiceHandlerClient(ctx context.Context, mux \*runtime.ServeMux, client RepoServiceClient) error](#RegisterRepoServiceHandlerClient)
* [func RegisterRepoServiceHandlerFromEndpoint(ctx context.Context, mux \*runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error)](#RegisterRepoServiceHandlerFromEndpoint)
* [func RegisterRepoServiceServer(s \*grpc.Server, srv RepoServiceServer)](#RegisterRepoServiceServer)
* [type App](#App)
  * [func (\*App) Descriptor() ([]byte, []int)](#App.Descriptor)
  * [func (m \*App) GetCreated() \*google\_protobuf3.Timestamp](#App.GetCreated)
  * [func (m \*App) GetDescription() string](#App.GetDescription)
  * [func (m \*App) GetId() string](#App.GetId)
  * [func (m \*App) GetLastModified() \*google\_protobuf3.Timestamp](#App.GetLastModified)
  * [func (m \*App) GetName() string](#App.GetName)
  * [func (m \*App) GetRepoId() string](#App.GetRepoId)
  * [func (\*App) ProtoMessage()](#App.ProtoMessage)
  * [func (m \*App) Reset()](#App.Reset)
  * [func (m \*App) String() string](#App.String)
  * [func (this \*App) Validate() error](#App.Validate)
* [type AppId](#AppId)
  * [func (\*AppId) Descriptor() ([]byte, []int)](#AppId.Descriptor)
  * [func (m \*AppId) GetId() string](#AppId.GetId)
  * [func (\*AppId) ProtoMessage()](#AppId.ProtoMessage)
  * [func (m \*AppId) Reset()](#AppId.Reset)
  * [func (m \*AppId) String() string](#AppId.String)
  * [func (this \*AppId) Validate() error](#AppId.Validate)
* [type AppListRequest](#AppListRequest)
  * [func (\*AppListRequest) Descriptor() ([]byte, []int)](#AppListRequest.Descriptor)
  * [func (m \*AppListRequest) GetPageNumber() int32](#AppListRequest.GetPageNumber)
  * [func (m \*AppListRequest) GetPageSize() int32](#AppListRequest.GetPageSize)
  * [func (\*AppListRequest) ProtoMessage()](#AppListRequest.ProtoMessage)
  * [func (m \*AppListRequest) Reset()](#AppListRequest.Reset)
  * [func (m \*AppListRequest) String() string](#AppListRequest.String)
  * [func (this \*AppListRequest) Validate() error](#AppListRequest.Validate)
* [type AppListResponse](#AppListResponse)
  * [func (\*AppListResponse) Descriptor() ([]byte, []int)](#AppListResponse.Descriptor)
  * [func (m \*AppListResponse) GetCurrentPage() int32](#AppListResponse.GetCurrentPage)
  * [func (m \*AppListResponse) GetItems() []\*App](#AppListResponse.GetItems)
  * [func (m \*AppListResponse) GetPageSize() int32](#AppListResponse.GetPageSize)
  * [func (m \*AppListResponse) GetTotalItems() int32](#AppListResponse.GetTotalItems)
  * [func (m \*AppListResponse) GetTotalPages() int32](#AppListResponse.GetTotalPages)
  * [func (\*AppListResponse) ProtoMessage()](#AppListResponse.ProtoMessage)
  * [func (m \*AppListResponse) Reset()](#AppListResponse.Reset)
  * [func (m \*AppListResponse) String() string](#AppListResponse.String)
  * [func (this \*AppListResponse) Validate() error](#AppListResponse.Validate)
* [type AppRuntime](#AppRuntime)
  * [func (\*AppRuntime) Descriptor() ([]byte, []int)](#AppRuntime.Descriptor)
  * [func (m \*AppRuntime) GetCreated() \*google\_protobuf3.Timestamp](#AppRuntime.GetCreated)
  * [func (m \*AppRuntime) GetDescription() string](#AppRuntime.GetDescription)
  * [func (m \*AppRuntime) GetId() string](#AppRuntime.GetId)
  * [func (m \*AppRuntime) GetLastModified() \*google\_protobuf3.Timestamp](#AppRuntime.GetLastModified)
  * [func (m \*AppRuntime) GetName() string](#AppRuntime.GetName)
  * [func (m \*AppRuntime) GetUrl() string](#AppRuntime.GetUrl)
  * [func (\*AppRuntime) ProtoMessage()](#AppRuntime.ProtoMessage)
  * [func (m \*AppRuntime) Reset()](#AppRuntime.Reset)
  * [func (m \*AppRuntime) String() string](#AppRuntime.String)
  * [func (this \*AppRuntime) Validate() error](#AppRuntime.Validate)
* [type AppRuntimeId](#AppRuntimeId)
  * [func (\*AppRuntimeId) Descriptor() ([]byte, []int)](#AppRuntimeId.Descriptor)
  * [func (m \*AppRuntimeId) GetId() string](#AppRuntimeId.GetId)
  * [func (\*AppRuntimeId) ProtoMessage()](#AppRuntimeId.ProtoMessage)
  * [func (m \*AppRuntimeId) Reset()](#AppRuntimeId.Reset)
  * [func (m \*AppRuntimeId) String() string](#AppRuntimeId.String)
  * [func (this \*AppRuntimeId) Validate() error](#AppRuntimeId.Validate)
* [type AppRuntimeLabel](#AppRuntimeLabel)
  * [func (\*AppRuntimeLabel) Descriptor() ([]byte, []int)](#AppRuntimeLabel.Descriptor)
  * [func (m \*AppRuntimeLabel) GetAppRuntimeId() string](#AppRuntimeLabel.GetAppRuntimeId)
  * [func (m \*AppRuntimeLabel) GetLabelKey() string](#AppRuntimeLabel.GetLabelKey)
  * [func (m \*AppRuntimeLabel) GetLabelValue() string](#AppRuntimeLabel.GetLabelValue)
  * [func (\*AppRuntimeLabel) ProtoMessage()](#AppRuntimeLabel.ProtoMessage)
  * [func (m \*AppRuntimeLabel) Reset()](#AppRuntimeLabel.Reset)
  * [func (m \*AppRuntimeLabel) String() string](#AppRuntimeLabel.String)
  * [func (this \*AppRuntimeLabel) Validate() error](#AppRuntimeLabel.Validate)
* [type AppRuntimeListRequest](#AppRuntimeListRequest)
  * [func (\*AppRuntimeListRequest) Descriptor() ([]byte, []int)](#AppRuntimeListRequest.Descriptor)
  * [func (m \*AppRuntimeListRequest) GetPageNumber() int32](#AppRuntimeListRequest.GetPageNumber)
  * [func (m \*AppRuntimeListRequest) GetPageSize() int32](#AppRuntimeListRequest.GetPageSize)
  * [func (\*AppRuntimeListRequest) ProtoMessage()](#AppRuntimeListRequest.ProtoMessage)
  * [func (m \*AppRuntimeListRequest) Reset()](#AppRuntimeListRequest.Reset)
  * [func (m \*AppRuntimeListRequest) String() string](#AppRuntimeListRequest.String)
  * [func (this \*AppRuntimeListRequest) Validate() error](#AppRuntimeListRequest.Validate)
* [type AppRuntimeListResponse](#AppRuntimeListResponse)
  * [func (\*AppRuntimeListResponse) Descriptor() ([]byte, []int)](#AppRuntimeListResponse.Descriptor)
  * [func (m \*AppRuntimeListResponse) GetCurrentPage() int32](#AppRuntimeListResponse.GetCurrentPage)
  * [func (m \*AppRuntimeListResponse) GetItems() []\*AppRuntime](#AppRuntimeListResponse.GetItems)
  * [func (m \*AppRuntimeListResponse) GetPageSize() int32](#AppRuntimeListResponse.GetPageSize)
  * [func (m \*AppRuntimeListResponse) GetTotalItems() int32](#AppRuntimeListResponse.GetTotalItems)
  * [func (m \*AppRuntimeListResponse) GetTotalPages() int32](#AppRuntimeListResponse.GetTotalPages)
  * [func (\*AppRuntimeListResponse) ProtoMessage()](#AppRuntimeListResponse.ProtoMessage)
  * [func (m \*AppRuntimeListResponse) Reset()](#AppRuntimeListResponse.Reset)
  * [func (m \*AppRuntimeListResponse) String() string](#AppRuntimeListResponse.String)
  * [func (this \*AppRuntimeListResponse) Validate() error](#AppRuntimeListResponse.Validate)
* [type AppRuntimeServiceClient](#AppRuntimeServiceClient)
  * [func NewAppRuntimeServiceClient(cc \*grpc.ClientConn) AppRuntimeServiceClient](#NewAppRuntimeServiceClient)
* [type AppRuntimeServiceServer](#AppRuntimeServiceServer)
* [type AppServiceClient](#AppServiceClient)
  * [func NewAppServiceClient(cc \*grpc.ClientConn) AppServiceClient](#NewAppServiceClient)
* [type AppServiceServer](#AppServiceServer)
* [type Cluster](#Cluster)
  * [func (\*Cluster) Descriptor() ([]byte, []int)](#Cluster.Descriptor)
  * [func (m \*Cluster) GetAppId() string](#Cluster.GetAppId)
  * [func (m \*Cluster) GetAppVersion() string](#Cluster.GetAppVersion)
  * [func (m \*Cluster) GetCreated() \*google\_protobuf3.Timestamp](#Cluster.GetCreated)
  * [func (m \*Cluster) GetDescription() string](#Cluster.GetDescription)
  * [func (m \*Cluster) GetId() string](#Cluster.GetId)
  * [func (m \*Cluster) GetLastModified() \*google\_protobuf3.Timestamp](#Cluster.GetLastModified)
  * [func (m \*Cluster) GetName() string](#Cluster.GetName)
  * [func (m \*Cluster) GetStatus() string](#Cluster.GetStatus)
  * [func (m \*Cluster) GetTransitionStatus() string](#Cluster.GetTransitionStatus)
  * [func (\*Cluster) ProtoMessage()](#Cluster.ProtoMessage)
  * [func (m \*Cluster) Reset()](#Cluster.Reset)
  * [func (m \*Cluster) String() string](#Cluster.String)
  * [func (this \*Cluster) Validate() error](#Cluster.Validate)
* [type ClusterId](#ClusterId)
  * [func (\*ClusterId) Descriptor() ([]byte, []int)](#ClusterId.Descriptor)
  * [func (m \*ClusterId) GetId() string](#ClusterId.GetId)
  * [func (\*ClusterId) ProtoMessage()](#ClusterId.ProtoMessage)
  * [func (m \*ClusterId) Reset()](#ClusterId.Reset)
  * [func (m \*ClusterId) String() string](#ClusterId.String)
  * [func (this \*ClusterId) Validate() error](#ClusterId.Validate)
* [type ClusterIds](#ClusterIds)
  * [func (\*ClusterIds) Descriptor() ([]byte, []int)](#ClusterIds.Descriptor)
  * [func (m \*ClusterIds) GetIds() string](#ClusterIds.GetIds)
  * [func (\*ClusterIds) ProtoMessage()](#ClusterIds.ProtoMessage)
  * [func (m \*ClusterIds) Reset()](#ClusterIds.Reset)
  * [func (m \*ClusterIds) String() string](#ClusterIds.String)
  * [func (this \*ClusterIds) Validate() error](#ClusterIds.Validate)
* [type ClusterListRequest](#ClusterListRequest)
  * [func (\*ClusterListRequest) Descriptor() ([]byte, []int)](#ClusterListRequest.Descriptor)
  * [func (m \*ClusterListRequest) GetPageNumber() int32](#ClusterListRequest.GetPageNumber)
  * [func (m \*ClusterListRequest) GetPageSize() int32](#ClusterListRequest.GetPageSize)
  * [func (\*ClusterListRequest) ProtoMessage()](#ClusterListRequest.ProtoMessage)
  * [func (m \*ClusterListRequest) Reset()](#ClusterListRequest.Reset)
  * [func (m \*ClusterListRequest) String() string](#ClusterListRequest.String)
  * [func (this \*ClusterListRequest) Validate() error](#ClusterListRequest.Validate)
* [type ClusterListResponse](#ClusterListResponse)
  * [func (\*ClusterListResponse) Descriptor() ([]byte, []int)](#ClusterListResponse.Descriptor)
  * [func (m \*ClusterListResponse) GetCurrentPage() int32](#ClusterListResponse.GetCurrentPage)
  * [func (m \*ClusterListResponse) GetItems() []\*Cluster](#ClusterListResponse.GetItems)
  * [func (m \*ClusterListResponse) GetPageSize() int32](#ClusterListResponse.GetPageSize)
  * [func (m \*ClusterListResponse) GetTotalItems() int32](#ClusterListResponse.GetTotalItems)
  * [func (m \*ClusterListResponse) GetTotalPages() int32](#ClusterListResponse.GetTotalPages)
  * [func (\*ClusterListResponse) ProtoMessage()](#ClusterListResponse.ProtoMessage)
  * [func (m \*ClusterListResponse) Reset()](#ClusterListResponse.Reset)
  * [func (m \*ClusterListResponse) String() string](#ClusterListResponse.String)
  * [func (this \*ClusterListResponse) Validate() error](#ClusterListResponse.Validate)
* [type ClusterNode](#ClusterNode)
  * [func (\*ClusterNode) Descriptor() ([]byte, []int)](#ClusterNode.Descriptor)
  * [func (m \*ClusterNode) GetClusterId() string](#ClusterNode.GetClusterId)
  * [func (m \*ClusterNode) GetCreated() \*google\_protobuf3.Timestamp](#ClusterNode.GetCreated)
  * [func (m \*ClusterNode) GetDescription() string](#ClusterNode.GetDescription)
  * [func (m \*ClusterNode) GetId() string](#ClusterNode.GetId)
  * [func (m \*ClusterNode) GetInstanceId() string](#ClusterNode.GetInstanceId)
  * [func (m \*ClusterNode) GetLastModified() \*google\_protobuf3.Timestamp](#ClusterNode.GetLastModified)
  * [func (m \*ClusterNode) GetName() string](#ClusterNode.GetName)
  * [func (m \*ClusterNode) GetPrivateIp() string](#ClusterNode.GetPrivateIp)
  * [func (m \*ClusterNode) GetStatus() string](#ClusterNode.GetStatus)
  * [func (m \*ClusterNode) GetTransitionStatus() string](#ClusterNode.GetTransitionStatus)
  * [func (\*ClusterNode) ProtoMessage()](#ClusterNode.ProtoMessage)
  * [func (m \*ClusterNode) Reset()](#ClusterNode.Reset)
  * [func (m \*ClusterNode) String() string](#ClusterNode.String)
  * [func (this \*ClusterNode) Validate() error](#ClusterNode.Validate)
* [type ClusterNodeId](#ClusterNodeId)
  * [func (\*ClusterNodeId) Descriptor() ([]byte, []int)](#ClusterNodeId.Descriptor)
  * [func (m \*ClusterNodeId) GetId() string](#ClusterNodeId.GetId)
  * [func (\*ClusterNodeId) ProtoMessage()](#ClusterNodeId.ProtoMessage)
  * [func (m \*ClusterNodeId) Reset()](#ClusterNodeId.Reset)
  * [func (m \*ClusterNodeId) String() string](#ClusterNodeId.String)
  * [func (this \*ClusterNodeId) Validate() error](#ClusterNodeId.Validate)
* [type ClusterNodeIds](#ClusterNodeIds)
  * [func (\*ClusterNodeIds) Descriptor() ([]byte, []int)](#ClusterNodeIds.Descriptor)
  * [func (m \*ClusterNodeIds) GetIds() string](#ClusterNodeIds.GetIds)
  * [func (\*ClusterNodeIds) ProtoMessage()](#ClusterNodeIds.ProtoMessage)
  * [func (m \*ClusterNodeIds) Reset()](#ClusterNodeIds.Reset)
  * [func (m \*ClusterNodeIds) String() string](#ClusterNodeIds.String)
  * [func (this \*ClusterNodeIds) Validate() error](#ClusterNodeIds.Validate)
* [type ClusterNodeListRequest](#ClusterNodeListRequest)
  * [func (\*ClusterNodeListRequest) Descriptor() ([]byte, []int)](#ClusterNodeListRequest.Descriptor)
  * [func (m \*ClusterNodeListRequest) GetPageNumber() int32](#ClusterNodeListRequest.GetPageNumber)
  * [func (m \*ClusterNodeListRequest) GetPageSize() int32](#ClusterNodeListRequest.GetPageSize)
  * [func (\*ClusterNodeListRequest) ProtoMessage()](#ClusterNodeListRequest.ProtoMessage)
  * [func (m \*ClusterNodeListRequest) Reset()](#ClusterNodeListRequest.Reset)
  * [func (m \*ClusterNodeListRequest) String() string](#ClusterNodeListRequest.String)
  * [func (this \*ClusterNodeListRequest) Validate() error](#ClusterNodeListRequest.Validate)
* [type ClusterNodeListResponse](#ClusterNodeListResponse)
  * [func (\*ClusterNodeListResponse) Descriptor() ([]byte, []int)](#ClusterNodeListResponse.Descriptor)
  * [func (m \*ClusterNodeListResponse) GetCurrentPage() int32](#ClusterNodeListResponse.GetCurrentPage)
  * [func (m \*ClusterNodeListResponse) GetItems() []\*ClusterNode](#ClusterNodeListResponse.GetItems)
  * [func (m \*ClusterNodeListResponse) GetPageSize() int32](#ClusterNodeListResponse.GetPageSize)
  * [func (m \*ClusterNodeListResponse) GetTotalItems() int32](#ClusterNodeListResponse.GetTotalItems)
  * [func (m \*ClusterNodeListResponse) GetTotalPages() int32](#ClusterNodeListResponse.GetTotalPages)
  * [func (\*ClusterNodeListResponse) ProtoMessage()](#ClusterNodeListResponse.ProtoMessage)
  * [func (m \*ClusterNodeListResponse) Reset()](#ClusterNodeListResponse.Reset)
  * [func (m \*ClusterNodeListResponse) String() string](#ClusterNodeListResponse.String)
  * [func (this \*ClusterNodeListResponse) Validate() error](#ClusterNodeListResponse.Validate)
* [type ClusterNodes](#ClusterNodes)
  * [func (\*ClusterNodes) Descriptor() ([]byte, []int)](#ClusterNodes.Descriptor)
  * [func (m \*ClusterNodes) GetItems() []\*ClusterNode](#ClusterNodes.GetItems)
  * [func (\*ClusterNodes) ProtoMessage()](#ClusterNodes.ProtoMessage)
  * [func (m \*ClusterNodes) Reset()](#ClusterNodes.Reset)
  * [func (m \*ClusterNodes) String() string](#ClusterNodes.String)
  * [func (this \*ClusterNodes) Validate() error](#ClusterNodes.Validate)
* [type ClusterServiceClient](#ClusterServiceClient)
  * [func NewClusterServiceClient(cc \*grpc.ClientConn) ClusterServiceClient](#NewClusterServiceClient)
* [type ClusterServiceServer](#ClusterServiceServer)
* [type Clusters](#Clusters)
  * [func (\*Clusters) Descriptor() ([]byte, []int)](#Clusters.Descriptor)
  * [func (m \*Clusters) GetItems() []\*Cluster](#Clusters.GetItems)
  * [func (\*Clusters) ProtoMessage()](#Clusters.ProtoMessage)
  * [func (m \*Clusters) Reset()](#Clusters.Reset)
  * [func (m \*Clusters) String() string](#Clusters.String)
  * [func (this \*Clusters) Validate() error](#Clusters.Validate)
* [type Repo](#Repo)
  * [func (\*Repo) Descriptor() ([]byte, []int)](#Repo.Descriptor)
  * [func (m \*Repo) GetCreated() \*google\_protobuf3.Timestamp](#Repo.GetCreated)
  * [func (m \*Repo) GetDescription() string](#Repo.GetDescription)
  * [func (m \*Repo) GetId() string](#Repo.GetId)
  * [func (m \*Repo) GetLastModified() \*google\_protobuf3.Timestamp](#Repo.GetLastModified)
  * [func (m \*Repo) GetName() string](#Repo.GetName)
  * [func (m \*Repo) GetUrl() string](#Repo.GetUrl)
  * [func (\*Repo) ProtoMessage()](#Repo.ProtoMessage)
  * [func (m \*Repo) Reset()](#Repo.Reset)
  * [func (m \*Repo) String() string](#Repo.String)
  * [func (this \*Repo) Validate() error](#Repo.Validate)
* [type RepoId](#RepoId)
  * [func (\*RepoId) Descriptor() ([]byte, []int)](#RepoId.Descriptor)
  * [func (m \*RepoId) GetId() string](#RepoId.GetId)
  * [func (\*RepoId) ProtoMessage()](#RepoId.ProtoMessage)
  * [func (m \*RepoId) Reset()](#RepoId.Reset)
  * [func (m \*RepoId) String() string](#RepoId.String)
  * [func (this \*RepoId) Validate() error](#RepoId.Validate)
* [type RepoLabel](#RepoLabel)
  * [func (\*RepoLabel) Descriptor() ([]byte, []int)](#RepoLabel.Descriptor)
  * [func (m \*RepoLabel) GetLabelKey() string](#RepoLabel.GetLabelKey)
  * [func (m \*RepoLabel) GetLabelValue() string](#RepoLabel.GetLabelValue)
  * [func (m \*RepoLabel) GetRepoId() string](#RepoLabel.GetRepoId)
  * [func (\*RepoLabel) ProtoMessage()](#RepoLabel.ProtoMessage)
  * [func (m \*RepoLabel) Reset()](#RepoLabel.Reset)
  * [func (m \*RepoLabel) String() string](#RepoLabel.String)
  * [func (this \*RepoLabel) Validate() error](#RepoLabel.Validate)
* [type RepoListRequest](#RepoListRequest)
  * [func (\*RepoListRequest) Descriptor() ([]byte, []int)](#RepoListRequest.Descriptor)
  * [func (m \*RepoListRequest) GetPageNumber() int32](#RepoListRequest.GetPageNumber)
  * [func (m \*RepoListRequest) GetPageSize() int32](#RepoListRequest.GetPageSize)
  * [func (\*RepoListRequest) ProtoMessage()](#RepoListRequest.ProtoMessage)
  * [func (m \*RepoListRequest) Reset()](#RepoListRequest.Reset)
  * [func (m \*RepoListRequest) String() string](#RepoListRequest.String)
  * [func (this \*RepoListRequest) Validate() error](#RepoListRequest.Validate)
* [type RepoListResponse](#RepoListResponse)
  * [func (\*RepoListResponse) Descriptor() ([]byte, []int)](#RepoListResponse.Descriptor)
  * [func (m \*RepoListResponse) GetCurrentPage() int32](#RepoListResponse.GetCurrentPage)
  * [func (m \*RepoListResponse) GetItems() []\*Repo](#RepoListResponse.GetItems)
  * [func (m \*RepoListResponse) GetPageSize() int32](#RepoListResponse.GetPageSize)
  * [func (m \*RepoListResponse) GetTotalItems() int32](#RepoListResponse.GetTotalItems)
  * [func (m \*RepoListResponse) GetTotalPages() int32](#RepoListResponse.GetTotalPages)
  * [func (\*RepoListResponse) ProtoMessage()](#RepoListResponse.ProtoMessage)
  * [func (m \*RepoListResponse) Reset()](#RepoListResponse.Reset)
  * [func (m \*RepoListResponse) String() string](#RepoListResponse.String)
  * [func (this \*RepoListResponse) Validate() error](#RepoListResponse.Validate)
* [type RepoSelector](#RepoSelector)
  * [func (\*RepoSelector) Descriptor() ([]byte, []int)](#RepoSelector.Descriptor)
  * [func (m \*RepoSelector) GetRepoId() string](#RepoSelector.GetRepoId)
  * [func (m \*RepoSelector) GetSelectorKey() string](#RepoSelector.GetSelectorKey)
  * [func (m \*RepoSelector) GetSelectorValue() string](#RepoSelector.GetSelectorValue)
  * [func (\*RepoSelector) ProtoMessage()](#RepoSelector.ProtoMessage)
  * [func (m \*RepoSelector) Reset()](#RepoSelector.Reset)
  * [func (m \*RepoSelector) String() string](#RepoSelector.String)
  * [func (this \*RepoSelector) Validate() error](#RepoSelector.Validate)
* [type RepoServiceClient](#RepoServiceClient)
  * [func NewRepoServiceClient(cc \*grpc.ClientConn) RepoServiceClient](#NewRepoServiceClient)
* [type RepoServiceServer](#RepoServiceServer)

#### <a name="pkg-files">Package files</a>
[annotations.pb.go](./annotations.pb.go) [annotations.validator.pb.go](./annotations.validator.pb.go) [app.pb.go](./app.pb.go) [app.pb.gw.go](./app.pb.gw.go) [app.validator.pb.go](./app.validator.pb.go) [app_runtime.pb.go](./app_runtime.pb.go) [app_runtime.pb.gw.go](./app_runtime.pb.gw.go) [app_runtime.validator.pb.go](./app_runtime.validator.pb.go) [cluster.pb.go](./cluster.pb.go) [cluster.pb.gw.go](./cluster.pb.gw.go) [cluster.validator.pb.go](./cluster.validator.pb.go) [repo.pb.go](./repo.pb.go) [repo.pb.gw.go](./repo.pb.gw.go) [repo.validator.pb.go](./repo.validator.pb.go) 

## <a name="pkg-constants">Constants</a>
``` go
const Default_AppListRequest_PageNumber int32 = 1
```
``` go
const Default_AppListRequest_PageSize int32 = 10
```
``` go
const Default_AppRuntimeListRequest_PageNumber int32 = 1
```
``` go
const Default_AppRuntimeListRequest_PageSize int32 = 10
```
``` go
const Default_ClusterListRequest_PageNumber int32 = 1
```
``` go
const Default_ClusterListRequest_PageSize int32 = 10
```
``` go
const Default_ClusterNodeListRequest_PageNumber int32 = 1
```
``` go
const Default_ClusterNodeListRequest_PageSize int32 = 10
```
``` go
const Default_RepoListRequest_PageNumber int32 = 1
```
``` go
const Default_RepoListRequest_PageSize int32 = 10
```

## <a name="pkg-variables">Variables</a>
``` go
var E_Openapiv2FieldSchema = &proto.ExtensionDesc{
    ExtendedType:  (*google_protobuf1.FieldOptions)(nil),
    ExtensionType: (*grpc_gateway_protoc_gen_swagger_options.JSONSchema)(nil),
    Field:         1042,
    Name:          "openpitrix.openapiv2_field_schema",
    Tag:           "bytes,1042,opt,name=openapiv2_field_schema,json=openapiv2FieldSchema",
    Filename:      "annotations.proto",
}
```

## <a name="RegisterAppRuntimeServiceHandler">func</a> [RegisterAppRuntimeServiceHandler](./app_runtime.pb.gw.go#L173)
``` go
func RegisterAppRuntimeServiceHandler(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error
```
RegisterAppRuntimeServiceHandler registers the http handlers for service AppRuntimeService to "mux".
The handlers forward requests to the grpc endpoint over "conn".

## <a name="RegisterAppRuntimeServiceHandlerClient">func</a> [RegisterAppRuntimeServiceHandlerClient](./app_runtime.pb.gw.go#L182)
``` go
func RegisterAppRuntimeServiceHandlerClient(ctx context.Context, mux *runtime.ServeMux, client AppRuntimeServiceClient) error
```
RegisterAppRuntimeServiceHandler registers the http handlers for service AppRuntimeService to "mux".
The handlers forward requests to the grpc endpoint over the given implementation of "AppRuntimeServiceClient".
Note: the gRPC framework executes interceptors within the gRPC handler. If the passed in "AppRuntimeServiceClient"
doesn't go through the normal gRPC flow (creating a gRPC client etc.) then it will be up to the passed in
"AppRuntimeServiceClient" to call the correct interceptors.

## <a name="RegisterAppRuntimeServiceHandlerFromEndpoint">func</a> [RegisterAppRuntimeServiceHandlerFromEndpoint](./app_runtime.pb.gw.go#L148)
``` go
func RegisterAppRuntimeServiceHandlerFromEndpoint(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error)
```
RegisterAppRuntimeServiceHandlerFromEndpoint is same as RegisterAppRuntimeServiceHandler but
automatically dials to "endpoint" and closes the connection when "ctx" gets done.

## <a name="RegisterAppRuntimeServiceServer">func</a> [RegisterAppRuntimeServiceServer](./app_runtime.pb.go#L300)
``` go
func RegisterAppRuntimeServiceServer(s *grpc.Server, srv AppRuntimeServiceServer)
```

## <a name="RegisterAppServiceHandler">func</a> [RegisterAppServiceHandler](./app.pb.gw.go#L173)
``` go
func RegisterAppServiceHandler(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error
```
RegisterAppServiceHandler registers the http handlers for service AppService to "mux".
The handlers forward requests to the grpc endpoint over "conn".

## <a name="RegisterAppServiceHandlerClient">func</a> [RegisterAppServiceHandlerClient](./app.pb.gw.go#L182)
``` go
func RegisterAppServiceHandlerClient(ctx context.Context, mux *runtime.ServeMux, client AppServiceClient) error
```
RegisterAppServiceHandler registers the http handlers for service AppService to "mux".
The handlers forward requests to the grpc endpoint over the given implementation of "AppServiceClient".
Note: the gRPC framework executes interceptors within the gRPC handler. If the passed in "AppServiceClient"
doesn't go through the normal gRPC flow (creating a gRPC client etc.) then it will be up to the passed in
"AppServiceClient" to call the correct interceptors.

## <a name="RegisterAppServiceHandlerFromEndpoint">func</a> [RegisterAppServiceHandlerFromEndpoint](./app.pb.gw.go#L148)
``` go
func RegisterAppServiceHandlerFromEndpoint(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error)
```
RegisterAppServiceHandlerFromEndpoint is same as RegisterAppServiceHandler but
automatically dials to "endpoint" and closes the connection when "ctx" gets done.

## <a name="RegisterAppServiceServer">func</a> [RegisterAppServiceServer](./app.pb.go#L264)
``` go
func RegisterAppServiceServer(s *grpc.Server, srv AppServiceServer)
```

## <a name="RegisterClusterServiceHandler">func</a> [RegisterClusterServiceHandler](./cluster.pb.gw.go#L288)
``` go
func RegisterClusterServiceHandler(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error
```
RegisterClusterServiceHandler registers the http handlers for service ClusterService to "mux".
The handlers forward requests to the grpc endpoint over "conn".

## <a name="RegisterClusterServiceHandlerClient">func</a> [RegisterClusterServiceHandlerClient](./cluster.pb.gw.go#L297)
``` go
func RegisterClusterServiceHandlerClient(ctx context.Context, mux *runtime.ServeMux, client ClusterServiceClient) error
```
RegisterClusterServiceHandler registers the http handlers for service ClusterService to "mux".
The handlers forward requests to the grpc endpoint over the given implementation of "ClusterServiceClient".
Note: the gRPC framework executes interceptors within the gRPC handler. If the passed in "ClusterServiceClient"
doesn't go through the normal gRPC flow (creating a gRPC client etc.) then it will be up to the passed in
"ClusterServiceClient" to call the correct interceptors.

## <a name="RegisterClusterServiceHandlerFromEndpoint">func</a> [RegisterClusterServiceHandlerFromEndpoint](./cluster.pb.gw.go#L263)
``` go
func RegisterClusterServiceHandlerFromEndpoint(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error)
```
RegisterClusterServiceHandlerFromEndpoint is same as RegisterClusterServiceHandler but
automatically dials to "endpoint" and closes the connection when "ctx" gets done.

## <a name="RegisterClusterServiceServer">func</a> [RegisterClusterServiceServer](./cluster.pb.go#L602)
``` go
func RegisterClusterServiceServer(s *grpc.Server, srv ClusterServiceServer)
```

## <a name="RegisterRepoServiceHandler">func</a> [RegisterRepoServiceHandler](./repo.pb.gw.go#L173)
``` go
func RegisterRepoServiceHandler(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error
```
RegisterRepoServiceHandler registers the http handlers for service RepoService to "mux".
The handlers forward requests to the grpc endpoint over "conn".

## <a name="RegisterRepoServiceHandlerClient">func</a> [RegisterRepoServiceHandlerClient](./repo.pb.gw.go#L182)
``` go
func RegisterRepoServiceHandlerClient(ctx context.Context, mux *runtime.ServeMux, client RepoServiceClient) error
```
RegisterRepoServiceHandler registers the http handlers for service RepoService to "mux".
The handlers forward requests to the grpc endpoint over the given implementation of "RepoServiceClient".
Note: the gRPC framework executes interceptors within the gRPC handler. If the passed in "RepoServiceClient"
doesn't go through the normal gRPC flow (creating a gRPC client etc.) then it will be up to the passed in
"RepoServiceClient" to call the correct interceptors.

## <a name="RegisterRepoServiceHandlerFromEndpoint">func</a> [RegisterRepoServiceHandlerFromEndpoint](./repo.pb.gw.go#L148)
``` go
func RegisterRepoServiceHandlerFromEndpoint(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error)
```
RegisterRepoServiceHandlerFromEndpoint is same as RegisterRepoServiceHandler but
automatically dials to "endpoint" and closes the connection when "ctx" gets done.

## <a name="RegisterRepoServiceServer">func</a> [RegisterRepoServiceServer](./repo.pb.go#L332)
``` go
func RegisterRepoServiceServer(s *grpc.Server, srv RepoServiceServer)
```

## <a name="App">type</a> [App](./app.pb.go#L25-L33)
``` go
type App struct {
    Id               *string                     `protobuf:"bytes,1,req,name=id" json:"id,omitempty"`
    Name             *string                     `protobuf:"bytes,2,opt,name=name" json:"name,omitempty"`
    Description      *string                     `protobuf:"bytes,3,opt,name=description" json:"description,omitempty"`
    RepoId           *string                     `protobuf:"bytes,4,opt,name=repo_id,json=repoId" json:"repo_id,omitempty"`
    Created          *google_protobuf3.Timestamp `protobuf:"bytes,5,opt,name=created" json:"created,omitempty"`
    LastModified     *google_protobuf3.Timestamp `protobuf:"bytes,6,opt,name=last_modified,json=lastModified" json:"last_modified,omitempty"`
    XXX_unrecognized []byte                      `json:"-"`
}
```

### <a name="App.Descriptor">func</a> (\*App) [Descriptor](./app.pb.go#L38)
``` go
func (*App) Descriptor() ([]byte, []int)
```

### <a name="App.GetCreated">func</a> (\*App) [GetCreated](./app.pb.go#L68)
``` go
func (m *App) GetCreated() *google_protobuf3.Timestamp
```

### <a name="App.GetDescription">func</a> (\*App) [GetDescription](./app.pb.go#L54)
``` go
func (m *App) GetDescription() string
```

### <a name="App.GetId">func</a> (\*App) [GetId](./app.pb.go#L40)
``` go
func (m *App) GetId() string
```

### <a name="App.GetLastModified">func</a> (\*App) [GetLastModified](./app.pb.go#L75)
``` go
func (m *App) GetLastModified() *google_protobuf3.Timestamp
```

### <a name="App.GetName">func</a> (\*App) [GetName](./app.pb.go#L47)
``` go
func (m *App) GetName() string
```

### <a name="App.GetRepoId">func</a> (\*App) [GetRepoId](./app.pb.go#L61)
``` go
func (m *App) GetRepoId() string
```

### <a name="App.ProtoMessage">func</a> (\*App) [ProtoMessage](./app.pb.go#L37)
``` go
func (*App) ProtoMessage()
```

### <a name="App.Reset">func</a> (\*App) [Reset](./app.pb.go#L35)
``` go
func (m *App) Reset()
```

### <a name="App.String">func</a> (\*App) [String](./app.pb.go#L36)
``` go
func (m *App) String() string
```

### <a name="App.Validate">func</a> (\*App) [Validate](./app.validator.pb.go#L24)
``` go
func (this *App) Validate() error
```

## <a name="AppId">type</a> [AppId](./app.pb.go#L82-L85)
``` go
type AppId struct {
    Id               *string `protobuf:"bytes,1,req,name=id" json:"id,omitempty"`
    XXX_unrecognized []byte  `json:"-"`
}
```

### <a name="AppId.Descriptor">func</a> (\*AppId) [Descriptor](./app.pb.go#L90)
``` go
func (*AppId) Descriptor() ([]byte, []int)
```

### <a name="AppId.GetId">func</a> (\*AppId) [GetId](./app.pb.go#L92)
``` go
func (m *AppId) GetId() string
```

### <a name="AppId.ProtoMessage">func</a> (\*AppId) [ProtoMessage](./app.pb.go#L89)
``` go
func (*AppId) ProtoMessage()
```

### <a name="AppId.Reset">func</a> (\*AppId) [Reset](./app.pb.go#L87)
``` go
func (m *AppId) Reset()
```

### <a name="AppId.String">func</a> (\*AppId) [String](./app.pb.go#L88)
``` go
func (m *AppId) String() string
```

### <a name="AppId.Validate">func</a> (\*AppId) [Validate](./app.validator.pb.go#L45)
``` go
func (this *AppId) Validate() error
```

## <a name="AppListRequest">type</a> [AppListRequest](./app.pb.go#L99-L103)
``` go
type AppListRequest struct {
    PageSize         *int32 `protobuf:"varint,1,opt,name=page_size,json=pageSize,def=10" json:"page_size,omitempty"`
    PageNumber       *int32 `protobuf:"varint,2,opt,name=page_number,json=pageNumber,def=1" json:"page_number,omitempty"`
    XXX_unrecognized []byte `json:"-"`
}
```

### <a name="AppListRequest.Descriptor">func</a> (\*AppListRequest) [Descriptor](./app.pb.go#L108)
``` go
func (*AppListRequest) Descriptor() ([]byte, []int)
```

### <a name="AppListRequest.GetPageNumber">func</a> (\*AppListRequest) [GetPageNumber](./app.pb.go#L120)
``` go
func (m *AppListRequest) GetPageNumber() int32
```

### <a name="AppListRequest.GetPageSize">func</a> (\*AppListRequest) [GetPageSize](./app.pb.go#L113)
``` go
func (m *AppListRequest) GetPageSize() int32
```

### <a name="AppListRequest.ProtoMessage">func</a> (\*AppListRequest) [ProtoMessage](./app.pb.go#L107)
``` go
func (*AppListRequest) ProtoMessage()
```

### <a name="AppListRequest.Reset">func</a> (\*AppListRequest) [Reset](./app.pb.go#L105)
``` go
func (m *AppListRequest) Reset()
```

### <a name="AppListRequest.String">func</a> (\*AppListRequest) [String](./app.pb.go#L106)
``` go
func (m *AppListRequest) String() string
```

### <a name="AppListRequest.Validate">func</a> (\*AppListRequest) [Validate](./app.validator.pb.go#L53)
``` go
func (this *AppListRequest) Validate() error
```

## <a name="AppListResponse">type</a> [AppListResponse](./app.pb.go#L127-L134)
``` go
type AppListResponse struct {
    TotalItems       *int32 `protobuf:"varint,1,opt,name=total_items,json=totalItems" json:"total_items,omitempty"`
    TotalPages       *int32 `protobuf:"varint,2,opt,name=total_pages,json=totalPages" json:"total_pages,omitempty"`
    PageSize         *int32 `protobuf:"varint,3,opt,name=page_size,json=pageSize" json:"page_size,omitempty"`
    CurrentPage      *int32 `protobuf:"varint,4,opt,name=current_page,json=currentPage" json:"current_page,omitempty"`
    Items            []*App `protobuf:"bytes,5,rep,name=items" json:"items,omitempty"`
    XXX_unrecognized []byte `json:"-"`
}
```

### <a name="AppListResponse.Descriptor">func</a> (\*AppListResponse) [Descriptor](./app.pb.go#L139)
``` go
func (*AppListResponse) Descriptor() ([]byte, []int)
```

### <a name="AppListResponse.GetCurrentPage">func</a> (\*AppListResponse) [GetCurrentPage](./app.pb.go#L162)
``` go
func (m *AppListResponse) GetCurrentPage() int32
```

### <a name="AppListResponse.GetItems">func</a> (\*AppListResponse) [GetItems](./app.pb.go#L169)
``` go
func (m *AppListResponse) GetItems() []*App
```

### <a name="AppListResponse.GetPageSize">func</a> (\*AppListResponse) [GetPageSize](./app.pb.go#L155)
``` go
func (m *AppListResponse) GetPageSize() int32
```

### <a name="AppListResponse.GetTotalItems">func</a> (\*AppListResponse) [GetTotalItems](./app.pb.go#L141)
``` go
func (m *AppListResponse) GetTotalItems() int32
```

### <a name="AppListResponse.GetTotalPages">func</a> (\*AppListResponse) [GetTotalPages](./app.pb.go#L148)
``` go
func (m *AppListResponse) GetTotalPages() int32
```

### <a name="AppListResponse.ProtoMessage">func</a> (\*AppListResponse) [ProtoMessage](./app.pb.go#L138)
``` go
func (*AppListResponse) ProtoMessage()
```

### <a name="AppListResponse.Reset">func</a> (\*AppListResponse) [Reset](./app.pb.go#L136)
``` go
func (m *AppListResponse) Reset()
```

### <a name="AppListResponse.String">func</a> (\*AppListResponse) [String](./app.pb.go#L137)
``` go
func (m *AppListResponse) String() string
```

### <a name="AppListResponse.Validate">func</a> (\*AppListResponse) [Validate](./app.validator.pb.go#L56)
``` go
func (this *AppListResponse) Validate() error
```

## <a name="AppRuntime">type</a> [AppRuntime](./app_runtime.pb.go#L25-L33)
``` go
type AppRuntime struct {
    Id               *string                     `protobuf:"bytes,1,req,name=id" json:"id,omitempty"`
    Name             *string                     `protobuf:"bytes,2,opt,name=name" json:"name,omitempty"`
    Description      *string                     `protobuf:"bytes,3,opt,name=description" json:"description,omitempty"`
    Url              *string                     `protobuf:"bytes,4,opt,name=url" json:"url,omitempty"`
    Created          *google_protobuf3.Timestamp `protobuf:"bytes,5,opt,name=created" json:"created,omitempty"`
    LastModified     *google_protobuf3.Timestamp `protobuf:"bytes,6,opt,name=last_modified,json=lastModified" json:"last_modified,omitempty"`
    XXX_unrecognized []byte                      `json:"-"`
}
```

### <a name="AppRuntime.Descriptor">func</a> (\*AppRuntime) [Descriptor](./app_runtime.pb.go#L38)
``` go
func (*AppRuntime) Descriptor() ([]byte, []int)
```

### <a name="AppRuntime.GetCreated">func</a> (\*AppRuntime) [GetCreated](./app_runtime.pb.go#L68)
``` go
func (m *AppRuntime) GetCreated() *google_protobuf3.Timestamp
```

### <a name="AppRuntime.GetDescription">func</a> (\*AppRuntime) [GetDescription](./app_runtime.pb.go#L54)
``` go
func (m *AppRuntime) GetDescription() string
```

### <a name="AppRuntime.GetId">func</a> (\*AppRuntime) [GetId](./app_runtime.pb.go#L40)
``` go
func (m *AppRuntime) GetId() string
```

### <a name="AppRuntime.GetLastModified">func</a> (\*AppRuntime) [GetLastModified](./app_runtime.pb.go#L75)
``` go
func (m *AppRuntime) GetLastModified() *google_protobuf3.Timestamp
```

### <a name="AppRuntime.GetName">func</a> (\*AppRuntime) [GetName](./app_runtime.pb.go#L47)
``` go
func (m *AppRuntime) GetName() string
```

### <a name="AppRuntime.GetUrl">func</a> (\*AppRuntime) [GetUrl](./app_runtime.pb.go#L61)
``` go
func (m *AppRuntime) GetUrl() string
```

### <a name="AppRuntime.ProtoMessage">func</a> (\*AppRuntime) [ProtoMessage](./app_runtime.pb.go#L37)
``` go
func (*AppRuntime) ProtoMessage()
```

### <a name="AppRuntime.Reset">func</a> (\*AppRuntime) [Reset](./app_runtime.pb.go#L35)
``` go
func (m *AppRuntime) Reset()
```

### <a name="AppRuntime.String">func</a> (\*AppRuntime) [String](./app_runtime.pb.go#L36)
``` go
func (m *AppRuntime) String() string
```

### <a name="AppRuntime.Validate">func</a> (\*AppRuntime) [Validate](./app_runtime.validator.pb.go#L24)
``` go
func (this *AppRuntime) Validate() error
```

## <a name="AppRuntimeId">type</a> [AppRuntimeId](./app_runtime.pb.go#L115-L118)
``` go
type AppRuntimeId struct {
    Id               *string `protobuf:"bytes,1,req,name=id" json:"id,omitempty"`
    XXX_unrecognized []byte  `json:"-"`
}
```

### <a name="AppRuntimeId.Descriptor">func</a> (\*AppRuntimeId) [Descriptor](./app_runtime.pb.go#L123)
``` go
func (*AppRuntimeId) Descriptor() ([]byte, []int)
```

### <a name="AppRuntimeId.GetId">func</a> (\*AppRuntimeId) [GetId](./app_runtime.pb.go#L125)
``` go
func (m *AppRuntimeId) GetId() string
```

### <a name="AppRuntimeId.ProtoMessage">func</a> (\*AppRuntimeId) [ProtoMessage](./app_runtime.pb.go#L122)
``` go
func (*AppRuntimeId) ProtoMessage()
```

### <a name="AppRuntimeId.Reset">func</a> (\*AppRuntimeId) [Reset](./app_runtime.pb.go#L120)
``` go
func (m *AppRuntimeId) Reset()
```

### <a name="AppRuntimeId.String">func</a> (\*AppRuntimeId) [String](./app_runtime.pb.go#L121)
``` go
func (m *AppRuntimeId) String() string
```

### <a name="AppRuntimeId.Validate">func</a> (\*AppRuntimeId) [Validate](./app_runtime.validator.pb.go#L56)
``` go
func (this *AppRuntimeId) Validate() error
```

## <a name="AppRuntimeLabel">type</a> [AppRuntimeLabel](./app_runtime.pb.go#L82-L87)
``` go
type AppRuntimeLabel struct {
    AppRuntimeId     *string `protobuf:"bytes,1,req,name=app_runtime_id,json=appRuntimeId" json:"app_runtime_id,omitempty"`
    LabelKey         *string `protobuf:"bytes,2,req,name=label_key,json=labelKey" json:"label_key,omitempty"`
    LabelValue       *string `protobuf:"bytes,3,req,name=label_value,json=labelValue" json:"label_value,omitempty"`
    XXX_unrecognized []byte  `json:"-"`
}
```

### <a name="AppRuntimeLabel.Descriptor">func</a> (\*AppRuntimeLabel) [Descriptor](./app_runtime.pb.go#L92)
``` go
func (*AppRuntimeLabel) Descriptor() ([]byte, []int)
```

### <a name="AppRuntimeLabel.GetAppRuntimeId">func</a> (\*AppRuntimeLabel) [GetAppRuntimeId](./app_runtime.pb.go#L94)
``` go
func (m *AppRuntimeLabel) GetAppRuntimeId() string
```

### <a name="AppRuntimeLabel.GetLabelKey">func</a> (\*AppRuntimeLabel) [GetLabelKey](./app_runtime.pb.go#L101)
``` go
func (m *AppRuntimeLabel) GetLabelKey() string
```

### <a name="AppRuntimeLabel.GetLabelValue">func</a> (\*AppRuntimeLabel) [GetLabelValue](./app_runtime.pb.go#L108)
``` go
func (m *AppRuntimeLabel) GetLabelValue() string
```

### <a name="AppRuntimeLabel.ProtoMessage">func</a> (\*AppRuntimeLabel) [ProtoMessage](./app_runtime.pb.go#L91)
``` go
func (*AppRuntimeLabel) ProtoMessage()
```

### <a name="AppRuntimeLabel.Reset">func</a> (\*AppRuntimeLabel) [Reset](./app_runtime.pb.go#L89)
``` go
func (m *AppRuntimeLabel) Reset()
```

### <a name="AppRuntimeLabel.String">func</a> (\*AppRuntimeLabel) [String](./app_runtime.pb.go#L90)
``` go
func (m *AppRuntimeLabel) String() string
```

### <a name="AppRuntimeLabel.Validate">func</a> (\*AppRuntimeLabel) [Validate](./app_runtime.validator.pb.go#L45)
``` go
func (this *AppRuntimeLabel) Validate() error
```

## <a name="AppRuntimeListRequest">type</a> [AppRuntimeListRequest](./app_runtime.pb.go#L132-L136)
``` go
type AppRuntimeListRequest struct {
    PageSize         *int32 `protobuf:"varint,1,opt,name=page_size,json=pageSize,def=10" json:"page_size,omitempty"`
    PageNumber       *int32 `protobuf:"varint,2,opt,name=page_number,json=pageNumber,def=1" json:"page_number,omitempty"`
    XXX_unrecognized []byte `json:"-"`
}
```

### <a name="AppRuntimeListRequest.Descriptor">func</a> (\*AppRuntimeListRequest) [Descriptor](./app_runtime.pb.go#L141)
``` go
func (*AppRuntimeListRequest) Descriptor() ([]byte, []int)
```

### <a name="AppRuntimeListRequest.GetPageNumber">func</a> (\*AppRuntimeListRequest) [GetPageNumber](./app_runtime.pb.go#L153)
``` go
func (m *AppRuntimeListRequest) GetPageNumber() int32
```

### <a name="AppRuntimeListRequest.GetPageSize">func</a> (\*AppRuntimeListRequest) [GetPageSize](./app_runtime.pb.go#L146)
``` go
func (m *AppRuntimeListRequest) GetPageSize() int32
```

### <a name="AppRuntimeListRequest.ProtoMessage">func</a> (\*AppRuntimeListRequest) [ProtoMessage](./app_runtime.pb.go#L140)
``` go
func (*AppRuntimeListRequest) ProtoMessage()
```

### <a name="AppRuntimeListRequest.Reset">func</a> (\*AppRuntimeListRequest) [Reset](./app_runtime.pb.go#L138)
``` go
func (m *AppRuntimeListRequest) Reset()
```

### <a name="AppRuntimeListRequest.String">func</a> (\*AppRuntimeListRequest) [String](./app_runtime.pb.go#L139)
``` go
func (m *AppRuntimeListRequest) String() string
```

### <a name="AppRuntimeListRequest.Validate">func</a> (\*AppRuntimeListRequest) [Validate](./app_runtime.validator.pb.go#L64)
``` go
func (this *AppRuntimeListRequest) Validate() error
```

## <a name="AppRuntimeListResponse">type</a> [AppRuntimeListResponse](./app_runtime.pb.go#L160-L167)
``` go
type AppRuntimeListResponse struct {
    TotalItems       *int32        `protobuf:"varint,1,opt,name=total_items,json=totalItems" json:"total_items,omitempty"`
    TotalPages       *int32        `protobuf:"varint,2,opt,name=total_pages,json=totalPages" json:"total_pages,omitempty"`
    PageSize         *int32        `protobuf:"varint,3,opt,name=page_size,json=pageSize" json:"page_size,omitempty"`
    CurrentPage      *int32        `protobuf:"varint,4,opt,name=current_page,json=currentPage" json:"current_page,omitempty"`
    Items            []*AppRuntime `protobuf:"bytes,5,rep,name=items" json:"items,omitempty"`
    XXX_unrecognized []byte        `json:"-"`
}
```

### <a name="AppRuntimeListResponse.Descriptor">func</a> (\*AppRuntimeListResponse) [Descriptor](./app_runtime.pb.go#L172)
``` go
func (*AppRuntimeListResponse) Descriptor() ([]byte, []int)
```

### <a name="AppRuntimeListResponse.GetCurrentPage">func</a> (\*AppRuntimeListResponse) [GetCurrentPage](./app_runtime.pb.go#L195)
``` go
func (m *AppRuntimeListResponse) GetCurrentPage() int32
```

### <a name="AppRuntimeListResponse.GetItems">func</a> (\*AppRuntimeListResponse) [GetItems](./app_runtime.pb.go#L202)
``` go
func (m *AppRuntimeListResponse) GetItems() []*AppRuntime
```

### <a name="AppRuntimeListResponse.GetPageSize">func</a> (\*AppRuntimeListResponse) [GetPageSize](./app_runtime.pb.go#L188)
``` go
func (m *AppRuntimeListResponse) GetPageSize() int32
```

### <a name="AppRuntimeListResponse.GetTotalItems">func</a> (\*AppRuntimeListResponse) [GetTotalItems](./app_runtime.pb.go#L174)
``` go
func (m *AppRuntimeListResponse) GetTotalItems() int32
```

### <a name="AppRuntimeListResponse.GetTotalPages">func</a> (\*AppRuntimeListResponse) [GetTotalPages](./app_runtime.pb.go#L181)
``` go
func (m *AppRuntimeListResponse) GetTotalPages() int32
```

### <a name="AppRuntimeListResponse.ProtoMessage">func</a> (\*AppRuntimeListResponse) [ProtoMessage](./app_runtime.pb.go#L171)
``` go
func (*AppRuntimeListResponse) ProtoMessage()
```

### <a name="AppRuntimeListResponse.Reset">func</a> (\*AppRuntimeListResponse) [Reset](./app_runtime.pb.go#L169)
``` go
func (m *AppRuntimeListResponse) Reset()
```

### <a name="AppRuntimeListResponse.String">func</a> (\*AppRuntimeListResponse) [String](./app_runtime.pb.go#L170)
``` go
func (m *AppRuntimeListResponse) String() string
```

### <a name="AppRuntimeListResponse.Validate">func</a> (\*AppRuntimeListResponse) [Validate](./app_runtime.validator.pb.go#L67)
``` go
func (this *AppRuntimeListResponse) Validate() error
```

## <a name="AppRuntimeServiceClient">type</a> [AppRuntimeServiceClient](./app_runtime.pb.go#L227-L234)
``` go
type AppRuntimeServiceClient interface {
    GetAppRuntime(ctx context.Context, in *AppRuntimeId, opts ...grpc.CallOption) (*AppRuntime, error)
    // Returns a list containing all app runtimes.
    GetAppRuntimeList(ctx context.Context, in *AppRuntimeListRequest, opts ...grpc.CallOption) (*AppRuntimeListResponse, error)
    CreateAppRuntime(ctx context.Context, in *AppRuntime, opts ...grpc.CallOption) (*google_protobuf2.Empty, error)
    UpdateAppRuntime(ctx context.Context, in *AppRuntime, opts ...grpc.CallOption) (*google_protobuf2.Empty, error)
    DeleteAppRuntime(ctx context.Context, in *AppRuntimeId, opts ...grpc.CallOption) (*google_protobuf2.Empty, error)
}
```

### <a name="NewAppRuntimeServiceClient">func</a> [NewAppRuntimeServiceClient](./app_runtime.pb.go#L240)
``` go
func NewAppRuntimeServiceClient(cc *grpc.ClientConn) AppRuntimeServiceClient
```

## <a name="AppRuntimeServiceServer">type</a> [AppRuntimeServiceServer](./app_runtime.pb.go#L291-L298)
``` go
type AppRuntimeServiceServer interface {
    GetAppRuntime(context.Context, *AppRuntimeId) (*AppRuntime, error)
    // Returns a list containing all app runtimes.
    GetAppRuntimeList(context.Context, *AppRuntimeListRequest) (*AppRuntimeListResponse, error)
    CreateAppRuntime(context.Context, *AppRuntime) (*google_protobuf2.Empty, error)
    UpdateAppRuntime(context.Context, *AppRuntime) (*google_protobuf2.Empty, error)
    DeleteAppRuntime(context.Context, *AppRuntimeId) (*google_protobuf2.Empty, error)
}
```

## <a name="AppServiceClient">type</a> [AppServiceClient](./app.pb.go#L193-L199)
``` go
type AppServiceClient interface {
    GetApp(ctx context.Context, in *AppId, opts ...grpc.CallOption) (*App, error)
    GetAppList(ctx context.Context, in *AppListRequest, opts ...grpc.CallOption) (*AppListResponse, error)
    CreateApp(ctx context.Context, in *App, opts ...grpc.CallOption) (*google_protobuf2.Empty, error)
    UpdateApp(ctx context.Context, in *App, opts ...grpc.CallOption) (*google_protobuf2.Empty, error)
    DeleteApp(ctx context.Context, in *AppId, opts ...grpc.CallOption) (*google_protobuf2.Empty, error)
}
```

### <a name="NewAppServiceClient">func</a> [NewAppServiceClient](./app.pb.go#L205)
``` go
func NewAppServiceClient(cc *grpc.ClientConn) AppServiceClient
```

## <a name="AppServiceServer">type</a> [AppServiceServer](./app.pb.go#L256-L262)
``` go
type AppServiceServer interface {
    GetApp(context.Context, *AppId) (*App, error)
    GetAppList(context.Context, *AppListRequest) (*AppListResponse, error)
    CreateApp(context.Context, *App) (*google_protobuf2.Empty, error)
    UpdateApp(context.Context, *App) (*google_protobuf2.Empty, error)
    DeleteApp(context.Context, *AppId) (*google_protobuf2.Empty, error)
}
```

## <a name="Cluster">type</a> [Cluster](./cluster.pb.go#L25-L36)
``` go
type Cluster struct {
    Id               *string                     `protobuf:"bytes,1,req,name=id" json:"id,omitempty"`
    Name             *string                     `protobuf:"bytes,2,opt,name=name" json:"name,omitempty"`
    Description      *string                     `protobuf:"bytes,3,opt,name=description" json:"description,omitempty"`
    AppId            *string                     `protobuf:"bytes,4,opt,name=app_id,json=appId" json:"app_id,omitempty"`
    AppVersion       *string                     `protobuf:"bytes,5,opt,name=app_version,json=appVersion" json:"app_version,omitempty"`
    Status           *string                     `protobuf:"bytes,6,opt,name=status" json:"status,omitempty"`
    TransitionStatus *string                     `protobuf:"bytes,7,opt,name=transition_status,json=transitionStatus" json:"transition_status,omitempty"`
    Created          *google_protobuf3.Timestamp `protobuf:"bytes,8,opt,name=created" json:"created,omitempty"`
    LastModified     *google_protobuf3.Timestamp `protobuf:"bytes,9,opt,name=last_modified,json=lastModified" json:"last_modified,omitempty"`
    XXX_unrecognized []byte                      `json:"-"`
}
```

### <a name="Cluster.Descriptor">func</a> (\*Cluster) [Descriptor](./cluster.pb.go#L41)
``` go
func (*Cluster) Descriptor() ([]byte, []int)
```

### <a name="Cluster.GetAppId">func</a> (\*Cluster) [GetAppId](./cluster.pb.go#L64)
``` go
func (m *Cluster) GetAppId() string
```

### <a name="Cluster.GetAppVersion">func</a> (\*Cluster) [GetAppVersion](./cluster.pb.go#L71)
``` go
func (m *Cluster) GetAppVersion() string
```

### <a name="Cluster.GetCreated">func</a> (\*Cluster) [GetCreated](./cluster.pb.go#L92)
``` go
func (m *Cluster) GetCreated() *google_protobuf3.Timestamp
```

### <a name="Cluster.GetDescription">func</a> (\*Cluster) [GetDescription](./cluster.pb.go#L57)
``` go
func (m *Cluster) GetDescription() string
```

### <a name="Cluster.GetId">func</a> (\*Cluster) [GetId](./cluster.pb.go#L43)
``` go
func (m *Cluster) GetId() string
```

### <a name="Cluster.GetLastModified">func</a> (\*Cluster) [GetLastModified](./cluster.pb.go#L99)
``` go
func (m *Cluster) GetLastModified() *google_protobuf3.Timestamp
```

### <a name="Cluster.GetName">func</a> (\*Cluster) [GetName](./cluster.pb.go#L50)
``` go
func (m *Cluster) GetName() string
```

### <a name="Cluster.GetStatus">func</a> (\*Cluster) [GetStatus](./cluster.pb.go#L78)
``` go
func (m *Cluster) GetStatus() string
```

### <a name="Cluster.GetTransitionStatus">func</a> (\*Cluster) [GetTransitionStatus](./cluster.pb.go#L85)
``` go
func (m *Cluster) GetTransitionStatus() string
```

### <a name="Cluster.ProtoMessage">func</a> (\*Cluster) [ProtoMessage](./cluster.pb.go#L40)
``` go
func (*Cluster) ProtoMessage()
```

### <a name="Cluster.Reset">func</a> (\*Cluster) [Reset](./cluster.pb.go#L38)
``` go
func (m *Cluster) Reset()
```

### <a name="Cluster.String">func</a> (\*Cluster) [String](./cluster.pb.go#L39)
``` go
func (m *Cluster) String() string
```

### <a name="Cluster.Validate">func</a> (\*Cluster) [Validate](./cluster.validator.pb.go#L24)
``` go
func (this *Cluster) Validate() error
```

## <a name="ClusterId">type</a> [ClusterId](./cluster.pb.go#L229-L232)
``` go
type ClusterId struct {
    Id               *string `protobuf:"bytes,1,req,name=id" json:"id,omitempty"`
    XXX_unrecognized []byte  `json:"-"`
}
```

### <a name="ClusterId.Descriptor">func</a> (\*ClusterId) [Descriptor](./cluster.pb.go#L237)
``` go
func (*ClusterId) Descriptor() ([]byte, []int)
```

### <a name="ClusterId.GetId">func</a> (\*ClusterId) [GetId](./cluster.pb.go#L239)
``` go
func (m *ClusterId) GetId() string
```

### <a name="ClusterId.ProtoMessage">func</a> (\*ClusterId) [ProtoMessage](./cluster.pb.go#L236)
``` go
func (*ClusterId) ProtoMessage()
```

### <a name="ClusterId.Reset">func</a> (\*ClusterId) [Reset](./cluster.pb.go#L234)
``` go
func (m *ClusterId) Reset()
```

### <a name="ClusterId.String">func</a> (\*ClusterId) [String](./cluster.pb.go#L235)
``` go
func (m *ClusterId) String() string
```

### <a name="ClusterId.Validate">func</a> (\*ClusterId) [Validate](./cluster.validator.pb.go#L82)
``` go
func (this *ClusterId) Validate() error
```

## <a name="ClusterIds">type</a> [ClusterIds](./cluster.pb.go#L246-L249)
``` go
type ClusterIds struct {
    Ids              *string `protobuf:"bytes,1,req,name=ids" json:"ids,omitempty"`
    XXX_unrecognized []byte  `json:"-"`
}
```

### <a name="ClusterIds.Descriptor">func</a> (\*ClusterIds) [Descriptor](./cluster.pb.go#L254)
``` go
func (*ClusterIds) Descriptor() ([]byte, []int)
```

### <a name="ClusterIds.GetIds">func</a> (\*ClusterIds) [GetIds](./cluster.pb.go#L256)
``` go
func (m *ClusterIds) GetIds() string
```

### <a name="ClusterIds.ProtoMessage">func</a> (\*ClusterIds) [ProtoMessage](./cluster.pb.go#L253)
``` go
func (*ClusterIds) ProtoMessage()
```

### <a name="ClusterIds.Reset">func</a> (\*ClusterIds) [Reset](./cluster.pb.go#L251)
``` go
func (m *ClusterIds) Reset()
```

### <a name="ClusterIds.String">func</a> (\*ClusterIds) [String](./cluster.pb.go#L252)
``` go
func (m *ClusterIds) String() string
```

### <a name="ClusterIds.Validate">func</a> (\*ClusterIds) [Validate](./cluster.validator.pb.go#L93)
``` go
func (this *ClusterIds) Validate() error
```

## <a name="ClusterListRequest">type</a> [ClusterListRequest](./cluster.pb.go#L263-L267)
``` go
type ClusterListRequest struct {
    PageSize         *int32 `protobuf:"varint,1,opt,name=page_size,json=pageSize,def=10" json:"page_size,omitempty"`
    PageNumber       *int32 `protobuf:"varint,2,opt,name=page_number,json=pageNumber,def=1" json:"page_number,omitempty"`
    XXX_unrecognized []byte `json:"-"`
}
```

### <a name="ClusterListRequest.Descriptor">func</a> (\*ClusterListRequest) [Descriptor](./cluster.pb.go#L272)
``` go
func (*ClusterListRequest) Descriptor() ([]byte, []int)
```

### <a name="ClusterListRequest.GetPageNumber">func</a> (\*ClusterListRequest) [GetPageNumber](./cluster.pb.go#L284)
``` go
func (m *ClusterListRequest) GetPageNumber() int32
```

### <a name="ClusterListRequest.GetPageSize">func</a> (\*ClusterListRequest) [GetPageSize](./cluster.pb.go#L277)
``` go
func (m *ClusterListRequest) GetPageSize() int32
```

### <a name="ClusterListRequest.ProtoMessage">func</a> (\*ClusterListRequest) [ProtoMessage](./cluster.pb.go#L271)
``` go
func (*ClusterListRequest) ProtoMessage()
```

### <a name="ClusterListRequest.Reset">func</a> (\*ClusterListRequest) [Reset](./cluster.pb.go#L269)
``` go
func (m *ClusterListRequest) Reset()
```

### <a name="ClusterListRequest.String">func</a> (\*ClusterListRequest) [String](./cluster.pb.go#L270)
``` go
func (m *ClusterListRequest) String() string
```

### <a name="ClusterListRequest.Validate">func</a> (\*ClusterListRequest) [Validate](./cluster.validator.pb.go#L101)
``` go
func (this *ClusterListRequest) Validate() error
```

## <a name="ClusterListResponse">type</a> [ClusterListResponse](./cluster.pb.go#L291-L298)
``` go
type ClusterListResponse struct {
    TotalItems       *int32     `protobuf:"varint,1,opt,name=total_items,json=totalItems" json:"total_items,omitempty"`
    TotalPages       *int32     `protobuf:"varint,2,opt,name=total_pages,json=totalPages" json:"total_pages,omitempty"`
    PageSize         *int32     `protobuf:"varint,3,opt,name=page_size,json=pageSize" json:"page_size,omitempty"`
    CurrentPage      *int32     `protobuf:"varint,4,opt,name=current_page,json=currentPage" json:"current_page,omitempty"`
    Items            []*Cluster `protobuf:"bytes,5,rep,name=items" json:"items,omitempty"`
    XXX_unrecognized []byte     `json:"-"`
}
```

### <a name="ClusterListResponse.Descriptor">func</a> (\*ClusterListResponse) [Descriptor](./cluster.pb.go#L303)
``` go
func (*ClusterListResponse) Descriptor() ([]byte, []int)
```

### <a name="ClusterListResponse.GetCurrentPage">func</a> (\*ClusterListResponse) [GetCurrentPage](./cluster.pb.go#L326)
``` go
func (m *ClusterListResponse) GetCurrentPage() int32
```

### <a name="ClusterListResponse.GetItems">func</a> (\*ClusterListResponse) [GetItems](./cluster.pb.go#L333)
``` go
func (m *ClusterListResponse) GetItems() []*Cluster
```

### <a name="ClusterListResponse.GetPageSize">func</a> (\*ClusterListResponse) [GetPageSize](./cluster.pb.go#L319)
``` go
func (m *ClusterListResponse) GetPageSize() int32
```

### <a name="ClusterListResponse.GetTotalItems">func</a> (\*ClusterListResponse) [GetTotalItems](./cluster.pb.go#L305)
``` go
func (m *ClusterListResponse) GetTotalItems() int32
```

### <a name="ClusterListResponse.GetTotalPages">func</a> (\*ClusterListResponse) [GetTotalPages](./cluster.pb.go#L312)
``` go
func (m *ClusterListResponse) GetTotalPages() int32
```

### <a name="ClusterListResponse.ProtoMessage">func</a> (\*ClusterListResponse) [ProtoMessage](./cluster.pb.go#L302)
``` go
func (*ClusterListResponse) ProtoMessage()
```

### <a name="ClusterListResponse.Reset">func</a> (\*ClusterListResponse) [Reset](./cluster.pb.go#L300)
``` go
func (m *ClusterListResponse) Reset()
```

### <a name="ClusterListResponse.String">func</a> (\*ClusterListResponse) [String](./cluster.pb.go#L301)
``` go
func (m *ClusterListResponse) String() string
```

### <a name="ClusterListResponse.Validate">func</a> (\*ClusterListResponse) [Validate](./cluster.validator.pb.go#L104)
``` go
func (this *ClusterListResponse) Validate() error
```

## <a name="ClusterNode">type</a> [ClusterNode](./cluster.pb.go#L123-L135)
``` go
type ClusterNode struct {
    Id               *string                     `protobuf:"bytes,1,req,name=id" json:"id,omitempty"`
    InstanceId       *string                     `protobuf:"bytes,2,req,name=instance_id,json=instanceId" json:"instance_id,omitempty"`
    Name             *string                     `protobuf:"bytes,3,opt,name=name" json:"name,omitempty"`
    Description      *string                     `protobuf:"bytes,4,opt,name=description" json:"description,omitempty"`
    ClusterId        *string                     `protobuf:"bytes,5,opt,name=cluster_id,json=clusterId" json:"cluster_id,omitempty"`
    PrivateIp        *string                     `protobuf:"bytes,6,opt,name=private_ip,json=privateIp" json:"private_ip,omitempty"`
    Status           *string                     `protobuf:"bytes,7,opt,name=status" json:"status,omitempty"`
    TransitionStatus *string                     `protobuf:"bytes,8,opt,name=transition_status,json=transitionStatus" json:"transition_status,omitempty"`
    Created          *google_protobuf3.Timestamp `protobuf:"bytes,9,opt,name=created" json:"created,omitempty"`
    LastModified     *google_protobuf3.Timestamp `protobuf:"bytes,10,opt,name=last_modified,json=lastModified" json:"last_modified,omitempty"`
    XXX_unrecognized []byte                      `json:"-"`
}
```

### <a name="ClusterNode.Descriptor">func</a> (\*ClusterNode) [Descriptor](./cluster.pb.go#L140)
``` go
func (*ClusterNode) Descriptor() ([]byte, []int)
```

### <a name="ClusterNode.GetClusterId">func</a> (\*ClusterNode) [GetClusterId](./cluster.pb.go#L170)
``` go
func (m *ClusterNode) GetClusterId() string
```

### <a name="ClusterNode.GetCreated">func</a> (\*ClusterNode) [GetCreated](./cluster.pb.go#L198)
``` go
func (m *ClusterNode) GetCreated() *google_protobuf3.Timestamp
```

### <a name="ClusterNode.GetDescription">func</a> (\*ClusterNode) [GetDescription](./cluster.pb.go#L163)
``` go
func (m *ClusterNode) GetDescription() string
```

### <a name="ClusterNode.GetId">func</a> (\*ClusterNode) [GetId](./cluster.pb.go#L142)
``` go
func (m *ClusterNode) GetId() string
```

### <a name="ClusterNode.GetInstanceId">func</a> (\*ClusterNode) [GetInstanceId](./cluster.pb.go#L149)
``` go
func (m *ClusterNode) GetInstanceId() string
```

### <a name="ClusterNode.GetLastModified">func</a> (\*ClusterNode) [GetLastModified](./cluster.pb.go#L205)
``` go
func (m *ClusterNode) GetLastModified() *google_protobuf3.Timestamp
```

### <a name="ClusterNode.GetName">func</a> (\*ClusterNode) [GetName](./cluster.pb.go#L156)
``` go
func (m *ClusterNode) GetName() string
```

### <a name="ClusterNode.GetPrivateIp">func</a> (\*ClusterNode) [GetPrivateIp](./cluster.pb.go#L177)
``` go
func (m *ClusterNode) GetPrivateIp() string
```

### <a name="ClusterNode.GetStatus">func</a> (\*ClusterNode) [GetStatus](./cluster.pb.go#L184)
``` go
func (m *ClusterNode) GetStatus() string
```

### <a name="ClusterNode.GetTransitionStatus">func</a> (\*ClusterNode) [GetTransitionStatus](./cluster.pb.go#L191)
``` go
func (m *ClusterNode) GetTransitionStatus() string
```

### <a name="ClusterNode.ProtoMessage">func</a> (\*ClusterNode) [ProtoMessage](./cluster.pb.go#L139)
``` go
func (*ClusterNode) ProtoMessage()
```

### <a name="ClusterNode.Reset">func</a> (\*ClusterNode) [Reset](./cluster.pb.go#L137)
``` go
func (m *ClusterNode) Reset()
```

### <a name="ClusterNode.String">func</a> (\*ClusterNode) [String](./cluster.pb.go#L138)
``` go
func (m *ClusterNode) String() string
```

### <a name="ClusterNode.Validate">func</a> (\*ClusterNode) [Validate](./cluster.validator.pb.go#L53)
``` go
func (this *ClusterNode) Validate() error
```

## <a name="ClusterNodeId">type</a> [ClusterNodeId](./cluster.pb.go#L340-L343)
``` go
type ClusterNodeId struct {
    Id               *string `protobuf:"bytes,1,req,name=id" json:"id,omitempty"`
    XXX_unrecognized []byte  `json:"-"`
}
```

### <a name="ClusterNodeId.Descriptor">func</a> (\*ClusterNodeId) [Descriptor](./cluster.pb.go#L348)
``` go
func (*ClusterNodeId) Descriptor() ([]byte, []int)
```

### <a name="ClusterNodeId.GetId">func</a> (\*ClusterNodeId) [GetId](./cluster.pb.go#L350)
``` go
func (m *ClusterNodeId) GetId() string
```

### <a name="ClusterNodeId.ProtoMessage">func</a> (\*ClusterNodeId) [ProtoMessage](./cluster.pb.go#L347)
``` go
func (*ClusterNodeId) ProtoMessage()
```

### <a name="ClusterNodeId.Reset">func</a> (\*ClusterNodeId) [Reset](./cluster.pb.go#L345)
``` go
func (m *ClusterNodeId) Reset()
```

### <a name="ClusterNodeId.String">func</a> (\*ClusterNodeId) [String](./cluster.pb.go#L346)
``` go
func (m *ClusterNodeId) String() string
```

### <a name="ClusterNodeId.Validate">func</a> (\*ClusterNodeId) [Validate](./cluster.validator.pb.go#L115)
``` go
func (this *ClusterNodeId) Validate() error
```

## <a name="ClusterNodeIds">type</a> [ClusterNodeIds](./cluster.pb.go#L357-L360)
``` go
type ClusterNodeIds struct {
    Ids              *string `protobuf:"bytes,1,req,name=ids" json:"ids,omitempty"`
    XXX_unrecognized []byte  `json:"-"`
}
```

### <a name="ClusterNodeIds.Descriptor">func</a> (\*ClusterNodeIds) [Descriptor](./cluster.pb.go#L365)
``` go
func (*ClusterNodeIds) Descriptor() ([]byte, []int)
```

### <a name="ClusterNodeIds.GetIds">func</a> (\*ClusterNodeIds) [GetIds](./cluster.pb.go#L367)
``` go
func (m *ClusterNodeIds) GetIds() string
```

### <a name="ClusterNodeIds.ProtoMessage">func</a> (\*ClusterNodeIds) [ProtoMessage](./cluster.pb.go#L364)
``` go
func (*ClusterNodeIds) ProtoMessage()
```

### <a name="ClusterNodeIds.Reset">func</a> (\*ClusterNodeIds) [Reset](./cluster.pb.go#L362)
``` go
func (m *ClusterNodeIds) Reset()
```

### <a name="ClusterNodeIds.String">func</a> (\*ClusterNodeIds) [String](./cluster.pb.go#L363)
``` go
func (m *ClusterNodeIds) String() string
```

### <a name="ClusterNodeIds.Validate">func</a> (\*ClusterNodeIds) [Validate](./cluster.validator.pb.go#L126)
``` go
func (this *ClusterNodeIds) Validate() error
```

## <a name="ClusterNodeListRequest">type</a> [ClusterNodeListRequest](./cluster.pb.go#L374-L378)
``` go
type ClusterNodeListRequest struct {
    PageSize         *int32 `protobuf:"varint,1,opt,name=page_size,json=pageSize,def=10" json:"page_size,omitempty"`
    PageNumber       *int32 `protobuf:"varint,2,opt,name=page_number,json=pageNumber,def=1" json:"page_number,omitempty"`
    XXX_unrecognized []byte `json:"-"`
}
```

### <a name="ClusterNodeListRequest.Descriptor">func</a> (\*ClusterNodeListRequest) [Descriptor](./cluster.pb.go#L383)
``` go
func (*ClusterNodeListRequest) Descriptor() ([]byte, []int)
```

### <a name="ClusterNodeListRequest.GetPageNumber">func</a> (\*ClusterNodeListRequest) [GetPageNumber](./cluster.pb.go#L395)
``` go
func (m *ClusterNodeListRequest) GetPageNumber() int32
```

### <a name="ClusterNodeListRequest.GetPageSize">func</a> (\*ClusterNodeListRequest) [GetPageSize](./cluster.pb.go#L388)
``` go
func (m *ClusterNodeListRequest) GetPageSize() int32
```

### <a name="ClusterNodeListRequest.ProtoMessage">func</a> (\*ClusterNodeListRequest) [ProtoMessage](./cluster.pb.go#L382)
``` go
func (*ClusterNodeListRequest) ProtoMessage()
```

### <a name="ClusterNodeListRequest.Reset">func</a> (\*ClusterNodeListRequest) [Reset](./cluster.pb.go#L380)
``` go
func (m *ClusterNodeListRequest) Reset()
```

### <a name="ClusterNodeListRequest.String">func</a> (\*ClusterNodeListRequest) [String](./cluster.pb.go#L381)
``` go
func (m *ClusterNodeListRequest) String() string
```

### <a name="ClusterNodeListRequest.Validate">func</a> (\*ClusterNodeListRequest) [Validate](./cluster.validator.pb.go#L134)
``` go
func (this *ClusterNodeListRequest) Validate() error
```

## <a name="ClusterNodeListResponse">type</a> [ClusterNodeListResponse](./cluster.pb.go#L402-L409)
``` go
type ClusterNodeListResponse struct {
    TotalItems       *int32         `protobuf:"varint,1,opt,name=total_items,json=totalItems" json:"total_items,omitempty"`
    TotalPages       *int32         `protobuf:"varint,2,opt,name=total_pages,json=totalPages" json:"total_pages,omitempty"`
    PageSize         *int32         `protobuf:"varint,3,opt,name=page_size,json=pageSize" json:"page_size,omitempty"`
    CurrentPage      *int32         `protobuf:"varint,4,opt,name=current_page,json=currentPage" json:"current_page,omitempty"`
    Items            []*ClusterNode `protobuf:"bytes,5,rep,name=items" json:"items,omitempty"`
    XXX_unrecognized []byte         `json:"-"`
}
```

### <a name="ClusterNodeListResponse.Descriptor">func</a> (\*ClusterNodeListResponse) [Descriptor](./cluster.pb.go#L414)
``` go
func (*ClusterNodeListResponse) Descriptor() ([]byte, []int)
```

### <a name="ClusterNodeListResponse.GetCurrentPage">func</a> (\*ClusterNodeListResponse) [GetCurrentPage](./cluster.pb.go#L437)
``` go
func (m *ClusterNodeListResponse) GetCurrentPage() int32
```

### <a name="ClusterNodeListResponse.GetItems">func</a> (\*ClusterNodeListResponse) [GetItems](./cluster.pb.go#L444)
``` go
func (m *ClusterNodeListResponse) GetItems() []*ClusterNode
```

### <a name="ClusterNodeListResponse.GetPageSize">func</a> (\*ClusterNodeListResponse) [GetPageSize](./cluster.pb.go#L430)
``` go
func (m *ClusterNodeListResponse) GetPageSize() int32
```

### <a name="ClusterNodeListResponse.GetTotalItems">func</a> (\*ClusterNodeListResponse) [GetTotalItems](./cluster.pb.go#L416)
``` go
func (m *ClusterNodeListResponse) GetTotalItems() int32
```

### <a name="ClusterNodeListResponse.GetTotalPages">func</a> (\*ClusterNodeListResponse) [GetTotalPages](./cluster.pb.go#L423)
``` go
func (m *ClusterNodeListResponse) GetTotalPages() int32
```

### <a name="ClusterNodeListResponse.ProtoMessage">func</a> (\*ClusterNodeListResponse) [ProtoMessage](./cluster.pb.go#L413)
``` go
func (*ClusterNodeListResponse) ProtoMessage()
```

### <a name="ClusterNodeListResponse.Reset">func</a> (\*ClusterNodeListResponse) [Reset](./cluster.pb.go#L411)
``` go
func (m *ClusterNodeListResponse) Reset()
```

### <a name="ClusterNodeListResponse.String">func</a> (\*ClusterNodeListResponse) [String](./cluster.pb.go#L412)
``` go
func (m *ClusterNodeListResponse) String() string
```

### <a name="ClusterNodeListResponse.Validate">func</a> (\*ClusterNodeListResponse) [Validate](./cluster.validator.pb.go#L137)
``` go
func (this *ClusterNodeListResponse) Validate() error
```

## <a name="ClusterNodes">type</a> [ClusterNodes](./cluster.pb.go#L212-L215)
``` go
type ClusterNodes struct {
    Items            []*ClusterNode `protobuf:"bytes,1,rep,name=items" json:"items,omitempty"`
    XXX_unrecognized []byte         `json:"-"`
}
```

### <a name="ClusterNodes.Descriptor">func</a> (\*ClusterNodes) [Descriptor](./cluster.pb.go#L220)
``` go
func (*ClusterNodes) Descriptor() ([]byte, []int)
```

### <a name="ClusterNodes.GetItems">func</a> (\*ClusterNodes) [GetItems](./cluster.pb.go#L222)
``` go
func (m *ClusterNodes) GetItems() []*ClusterNode
```

### <a name="ClusterNodes.ProtoMessage">func</a> (\*ClusterNodes) [ProtoMessage](./cluster.pb.go#L219)
``` go
func (*ClusterNodes) ProtoMessage()
```

### <a name="ClusterNodes.Reset">func</a> (\*ClusterNodes) [Reset](./cluster.pb.go#L217)
``` go
func (m *ClusterNodes) Reset()
```

### <a name="ClusterNodes.String">func</a> (\*ClusterNodes) [String](./cluster.pb.go#L218)
``` go
func (m *ClusterNodes) String() string
```

### <a name="ClusterNodes.Validate">func</a> (\*ClusterNodes) [Validate](./cluster.validator.pb.go#L71)
``` go
func (this *ClusterNodes) Validate() error
```

## <a name="ClusterServiceClient">type</a> [ClusterServiceClient](./cluster.pb.go#L476-L487)
``` go
type ClusterServiceClient interface {
    GetClusters(ctx context.Context, in *ClusterIds, opts ...grpc.CallOption) (*Clusters, error)
    GetClusterList(ctx context.Context, in *ClusterListRequest, opts ...grpc.CallOption) (*ClusterListResponse, error)
    CreateCluster(ctx context.Context, in *Cluster, opts ...grpc.CallOption) (*google_protobuf2.Empty, error)
    UpdateCluster(ctx context.Context, in *Cluster, opts ...grpc.CallOption) (*google_protobuf2.Empty, error)
    DeleteClusters(ctx context.Context, in *ClusterIds, opts ...grpc.CallOption) (*google_protobuf2.Empty, error)
    GetClusterNodes(ctx context.Context, in *ClusterNodeIds, opts ...grpc.CallOption) (*ClusterNodes, error)
    GetClusterNodeList(ctx context.Context, in *ClusterNodeListRequest, opts ...grpc.CallOption) (*ClusterNodeListResponse, error)
    CreateClusterNodes(ctx context.Context, in *ClusterNodes, opts ...grpc.CallOption) (*google_protobuf2.Empty, error)
    UpdateClusterNode(ctx context.Context, in *ClusterNode, opts ...grpc.CallOption) (*google_protobuf2.Empty, error)
    DeleteClusterNodes(ctx context.Context, in *ClusterNodeIds, opts ...grpc.CallOption) (*google_protobuf2.Empty, error)
}
```

### <a name="NewClusterServiceClient">func</a> [NewClusterServiceClient](./cluster.pb.go#L493)
``` go
func NewClusterServiceClient(cc *grpc.ClientConn) ClusterServiceClient
```

## <a name="ClusterServiceServer">type</a> [ClusterServiceServer](./cluster.pb.go#L589-L600)
``` go
type ClusterServiceServer interface {
    GetClusters(context.Context, *ClusterIds) (*Clusters, error)
    GetClusterList(context.Context, *ClusterListRequest) (*ClusterListResponse, error)
    CreateCluster(context.Context, *Cluster) (*google_protobuf2.Empty, error)
    UpdateCluster(context.Context, *Cluster) (*google_protobuf2.Empty, error)
    DeleteClusters(context.Context, *ClusterIds) (*google_protobuf2.Empty, error)
    GetClusterNodes(context.Context, *ClusterNodeIds) (*ClusterNodes, error)
    GetClusterNodeList(context.Context, *ClusterNodeListRequest) (*ClusterNodeListResponse, error)
    CreateClusterNodes(context.Context, *ClusterNodes) (*google_protobuf2.Empty, error)
    UpdateClusterNode(context.Context, *ClusterNode) (*google_protobuf2.Empty, error)
    DeleteClusterNodes(context.Context, *ClusterNodeIds) (*google_protobuf2.Empty, error)
}
```

## <a name="Clusters">type</a> [Clusters](./cluster.pb.go#L106-L109)
``` go
type Clusters struct {
    Items            []*Cluster `protobuf:"bytes,1,rep,name=items" json:"items,omitempty"`
    XXX_unrecognized []byte     `json:"-"`
}
```

### <a name="Clusters.Descriptor">func</a> (\*Clusters) [Descriptor](./cluster.pb.go#L114)
``` go
func (*Clusters) Descriptor() ([]byte, []int)
```

### <a name="Clusters.GetItems">func</a> (\*Clusters) [GetItems](./cluster.pb.go#L116)
``` go
func (m *Clusters) GetItems() []*Cluster
```

### <a name="Clusters.ProtoMessage">func</a> (\*Clusters) [ProtoMessage](./cluster.pb.go#L113)
``` go
func (*Clusters) ProtoMessage()
```

### <a name="Clusters.Reset">func</a> (\*Clusters) [Reset](./cluster.pb.go#L111)
``` go
func (m *Clusters) Reset()
```

### <a name="Clusters.String">func</a> (\*Clusters) [String](./cluster.pb.go#L112)
``` go
func (m *Clusters) String() string
```

### <a name="Clusters.Validate">func</a> (\*Clusters) [Validate](./cluster.validator.pb.go#L42)
``` go
func (this *Clusters) Validate() error
```

## <a name="Repo">type</a> [Repo](./repo.pb.go#L25-L33)
``` go
type Repo struct {
    Id               *string                     `protobuf:"bytes,1,req,name=id" json:"id,omitempty"`
    Name             *string                     `protobuf:"bytes,2,opt,name=name" json:"name,omitempty"`
    Description      *string                     `protobuf:"bytes,3,opt,name=description" json:"description,omitempty"`
    Url              *string                     `protobuf:"bytes,4,opt,name=url" json:"url,omitempty"`
    Created          *google_protobuf3.Timestamp `protobuf:"bytes,5,opt,name=created" json:"created,omitempty"`
    LastModified     *google_protobuf3.Timestamp `protobuf:"bytes,6,opt,name=last_modified,json=lastModified" json:"last_modified,omitempty"`
    XXX_unrecognized []byte                      `json:"-"`
}
```

### <a name="Repo.Descriptor">func</a> (\*Repo) [Descriptor](./repo.pb.go#L38)
``` go
func (*Repo) Descriptor() ([]byte, []int)
```

### <a name="Repo.GetCreated">func</a> (\*Repo) [GetCreated](./repo.pb.go#L68)
``` go
func (m *Repo) GetCreated() *google_protobuf3.Timestamp
```

### <a name="Repo.GetDescription">func</a> (\*Repo) [GetDescription](./repo.pb.go#L54)
``` go
func (m *Repo) GetDescription() string
```

### <a name="Repo.GetId">func</a> (\*Repo) [GetId](./repo.pb.go#L40)
``` go
func (m *Repo) GetId() string
```

### <a name="Repo.GetLastModified">func</a> (\*Repo) [GetLastModified](./repo.pb.go#L75)
``` go
func (m *Repo) GetLastModified() *google_protobuf3.Timestamp
```

### <a name="Repo.GetName">func</a> (\*Repo) [GetName](./repo.pb.go#L47)
``` go
func (m *Repo) GetName() string
```

### <a name="Repo.GetUrl">func</a> (\*Repo) [GetUrl](./repo.pb.go#L61)
``` go
func (m *Repo) GetUrl() string
```

### <a name="Repo.ProtoMessage">func</a> (\*Repo) [ProtoMessage](./repo.pb.go#L37)
``` go
func (*Repo) ProtoMessage()
```

### <a name="Repo.Reset">func</a> (\*Repo) [Reset](./repo.pb.go#L35)
``` go
func (m *Repo) Reset()
```

### <a name="Repo.String">func</a> (\*Repo) [String](./repo.pb.go#L36)
``` go
func (m *Repo) String() string
```

### <a name="Repo.Validate">func</a> (\*Repo) [Validate](./repo.validator.pb.go#L24)
``` go
func (this *Repo) Validate() error
```

## <a name="RepoId">type</a> [RepoId](./repo.pb.go#L148-L151)
``` go
type RepoId struct {
    Id               *string `protobuf:"bytes,1,req,name=id" json:"id,omitempty"`
    XXX_unrecognized []byte  `json:"-"`
}
```

### <a name="RepoId.Descriptor">func</a> (\*RepoId) [Descriptor](./repo.pb.go#L156)
``` go
func (*RepoId) Descriptor() ([]byte, []int)
```

### <a name="RepoId.GetId">func</a> (\*RepoId) [GetId](./repo.pb.go#L158)
``` go
func (m *RepoId) GetId() string
```

### <a name="RepoId.ProtoMessage">func</a> (\*RepoId) [ProtoMessage](./repo.pb.go#L155)
``` go
func (*RepoId) ProtoMessage()
```

### <a name="RepoId.Reset">func</a> (\*RepoId) [Reset](./repo.pb.go#L153)
``` go
func (m *RepoId) Reset()
```

### <a name="RepoId.String">func</a> (\*RepoId) [String](./repo.pb.go#L154)
``` go
func (m *RepoId) String() string
```

### <a name="RepoId.Validate">func</a> (\*RepoId) [Validate](./repo.validator.pb.go#L67)
``` go
func (this *RepoId) Validate() error
```

## <a name="RepoLabel">type</a> [RepoLabel](./repo.pb.go#L82-L87)
``` go
type RepoLabel struct {
    RepoId           *string `protobuf:"bytes,1,req,name=repo_id,json=repoId" json:"repo_id,omitempty"`
    LabelKey         *string `protobuf:"bytes,2,req,name=label_key,json=labelKey" json:"label_key,omitempty"`
    LabelValue       *string `protobuf:"bytes,3,req,name=label_value,json=labelValue" json:"label_value,omitempty"`
    XXX_unrecognized []byte  `json:"-"`
}
```

### <a name="RepoLabel.Descriptor">func</a> (\*RepoLabel) [Descriptor](./repo.pb.go#L92)
``` go
func (*RepoLabel) Descriptor() ([]byte, []int)
```

### <a name="RepoLabel.GetLabelKey">func</a> (\*RepoLabel) [GetLabelKey](./repo.pb.go#L101)
``` go
func (m *RepoLabel) GetLabelKey() string
```

### <a name="RepoLabel.GetLabelValue">func</a> (\*RepoLabel) [GetLabelValue](./repo.pb.go#L108)
``` go
func (m *RepoLabel) GetLabelValue() string
```

### <a name="RepoLabel.GetRepoId">func</a> (\*RepoLabel) [GetRepoId](./repo.pb.go#L94)
``` go
func (m *RepoLabel) GetRepoId() string
```

### <a name="RepoLabel.ProtoMessage">func</a> (\*RepoLabel) [ProtoMessage](./repo.pb.go#L91)
``` go
func (*RepoLabel) ProtoMessage()
```

### <a name="RepoLabel.Reset">func</a> (\*RepoLabel) [Reset](./repo.pb.go#L89)
``` go
func (m *RepoLabel) Reset()
```

### <a name="RepoLabel.String">func</a> (\*RepoLabel) [String](./repo.pb.go#L90)
``` go
func (m *RepoLabel) String() string
```

### <a name="RepoLabel.Validate">func</a> (\*RepoLabel) [Validate](./repo.validator.pb.go#L45)
``` go
func (this *RepoLabel) Validate() error
```

## <a name="RepoListRequest">type</a> [RepoListRequest](./repo.pb.go#L165-L169)
``` go
type RepoListRequest struct {
    PageSize         *int32 `protobuf:"varint,1,opt,name=page_size,json=pageSize,def=10" json:"page_size,omitempty"`
    PageNumber       *int32 `protobuf:"varint,2,opt,name=page_number,json=pageNumber,def=1" json:"page_number,omitempty"`
    XXX_unrecognized []byte `json:"-"`
}
```

### <a name="RepoListRequest.Descriptor">func</a> (\*RepoListRequest) [Descriptor](./repo.pb.go#L174)
``` go
func (*RepoListRequest) Descriptor() ([]byte, []int)
```

### <a name="RepoListRequest.GetPageNumber">func</a> (\*RepoListRequest) [GetPageNumber](./repo.pb.go#L186)
``` go
func (m *RepoListRequest) GetPageNumber() int32
```

### <a name="RepoListRequest.GetPageSize">func</a> (\*RepoListRequest) [GetPageSize](./repo.pb.go#L179)
``` go
func (m *RepoListRequest) GetPageSize() int32
```

### <a name="RepoListRequest.ProtoMessage">func</a> (\*RepoListRequest) [ProtoMessage](./repo.pb.go#L173)
``` go
func (*RepoListRequest) ProtoMessage()
```

### <a name="RepoListRequest.Reset">func</a> (\*RepoListRequest) [Reset](./repo.pb.go#L171)
``` go
func (m *RepoListRequest) Reset()
```

### <a name="RepoListRequest.String">func</a> (\*RepoListRequest) [String](./repo.pb.go#L172)
``` go
func (m *RepoListRequest) String() string
```

### <a name="RepoListRequest.Validate">func</a> (\*RepoListRequest) [Validate](./repo.validator.pb.go#L75)
``` go
func (this *RepoListRequest) Validate() error
```

## <a name="RepoListResponse">type</a> [RepoListResponse](./repo.pb.go#L193-L200)
``` go
type RepoListResponse struct {
    TotalItems       *int32  `protobuf:"varint,1,opt,name=total_items,json=totalItems" json:"total_items,omitempty"`
    TotalPages       *int32  `protobuf:"varint,2,opt,name=total_pages,json=totalPages" json:"total_pages,omitempty"`
    PageSize         *int32  `protobuf:"varint,3,opt,name=page_size,json=pageSize" json:"page_size,omitempty"`
    CurrentPage      *int32  `protobuf:"varint,4,opt,name=current_page,json=currentPage" json:"current_page,omitempty"`
    Items            []*Repo `protobuf:"bytes,5,rep,name=items" json:"items,omitempty"`
    XXX_unrecognized []byte  `json:"-"`
}
```

### <a name="RepoListResponse.Descriptor">func</a> (\*RepoListResponse) [Descriptor](./repo.pb.go#L205)
``` go
func (*RepoListResponse) Descriptor() ([]byte, []int)
```

### <a name="RepoListResponse.GetCurrentPage">func</a> (\*RepoListResponse) [GetCurrentPage](./repo.pb.go#L228)
``` go
func (m *RepoListResponse) GetCurrentPage() int32
```

### <a name="RepoListResponse.GetItems">func</a> (\*RepoListResponse) [GetItems](./repo.pb.go#L235)
``` go
func (m *RepoListResponse) GetItems() []*Repo
```

### <a name="RepoListResponse.GetPageSize">func</a> (\*RepoListResponse) [GetPageSize](./repo.pb.go#L221)
``` go
func (m *RepoListResponse) GetPageSize() int32
```

### <a name="RepoListResponse.GetTotalItems">func</a> (\*RepoListResponse) [GetTotalItems](./repo.pb.go#L207)
``` go
func (m *RepoListResponse) GetTotalItems() int32
```

### <a name="RepoListResponse.GetTotalPages">func</a> (\*RepoListResponse) [GetTotalPages](./repo.pb.go#L214)
``` go
func (m *RepoListResponse) GetTotalPages() int32
```

### <a name="RepoListResponse.ProtoMessage">func</a> (\*RepoListResponse) [ProtoMessage](./repo.pb.go#L204)
``` go
func (*RepoListResponse) ProtoMessage()
```

### <a name="RepoListResponse.Reset">func</a> (\*RepoListResponse) [Reset](./repo.pb.go#L202)
``` go
func (m *RepoListResponse) Reset()
```

### <a name="RepoListResponse.String">func</a> (\*RepoListResponse) [String](./repo.pb.go#L203)
``` go
func (m *RepoListResponse) String() string
```

### <a name="RepoListResponse.Validate">func</a> (\*RepoListResponse) [Validate](./repo.validator.pb.go#L78)
``` go
func (this *RepoListResponse) Validate() error
```

## <a name="RepoSelector">type</a> [RepoSelector](./repo.pb.go#L115-L120)
``` go
type RepoSelector struct {
    RepoId           *string `protobuf:"bytes,1,req,name=repo_id,json=repoId" json:"repo_id,omitempty"`
    SelectorKey      *string `protobuf:"bytes,2,req,name=selector_key,json=selectorKey" json:"selector_key,omitempty"`
    SelectorValue    *string `protobuf:"bytes,3,req,name=selector_value,json=selectorValue" json:"selector_value,omitempty"`
    XXX_unrecognized []byte  `json:"-"`
}
```

### <a name="RepoSelector.Descriptor">func</a> (\*RepoSelector) [Descriptor](./repo.pb.go#L125)
``` go
func (*RepoSelector) Descriptor() ([]byte, []int)
```

### <a name="RepoSelector.GetRepoId">func</a> (\*RepoSelector) [GetRepoId](./repo.pb.go#L127)
``` go
func (m *RepoSelector) GetRepoId() string
```

### <a name="RepoSelector.GetSelectorKey">func</a> (\*RepoSelector) [GetSelectorKey](./repo.pb.go#L134)
``` go
func (m *RepoSelector) GetSelectorKey() string
```

### <a name="RepoSelector.GetSelectorValue">func</a> (\*RepoSelector) [GetSelectorValue](./repo.pb.go#L141)
``` go
func (m *RepoSelector) GetSelectorValue() string
```

### <a name="RepoSelector.ProtoMessage">func</a> (\*RepoSelector) [ProtoMessage](./repo.pb.go#L124)
``` go
func (*RepoSelector) ProtoMessage()
```

### <a name="RepoSelector.Reset">func</a> (\*RepoSelector) [Reset](./repo.pb.go#L122)
``` go
func (m *RepoSelector) Reset()
```

### <a name="RepoSelector.String">func</a> (\*RepoSelector) [String](./repo.pb.go#L123)
``` go
func (m *RepoSelector) String() string
```

### <a name="RepoSelector.Validate">func</a> (\*RepoSelector) [Validate](./repo.validator.pb.go#L56)
``` go
func (this *RepoSelector) Validate() error
```

## <a name="RepoServiceClient">type</a> [RepoServiceClient](./repo.pb.go#L261-L267)
``` go
type RepoServiceClient interface {
    GetRepo(ctx context.Context, in *RepoId, opts ...grpc.CallOption) (*Repo, error)
    GetRepoList(ctx context.Context, in *RepoListRequest, opts ...grpc.CallOption) (*RepoListResponse, error)
    CreateRepo(ctx context.Context, in *Repo, opts ...grpc.CallOption) (*google_protobuf2.Empty, error)
    UpdateRepo(ctx context.Context, in *Repo, opts ...grpc.CallOption) (*google_protobuf2.Empty, error)
    DeleteRepo(ctx context.Context, in *RepoId, opts ...grpc.CallOption) (*google_protobuf2.Empty, error)
}
```

### <a name="NewRepoServiceClient">func</a> [NewRepoServiceClient](./repo.pb.go#L273)
``` go
func NewRepoServiceClient(cc *grpc.ClientConn) RepoServiceClient
```

## <a name="RepoServiceServer">type</a> [RepoServiceServer](./repo.pb.go#L324-L330)
``` go
type RepoServiceServer interface {
    GetRepo(context.Context, *RepoId) (*Repo, error)
    GetRepoList(context.Context, *RepoListRequest) (*RepoListResponse, error)
    CreateRepo(context.Context, *Repo) (*google_protobuf2.Empty, error)
    UpdateRepo(context.Context, *Repo) (*google_protobuf2.Empty, error)
    DeleteRepo(context.Context, *RepoId) (*google_protobuf2.Empty, error)
}
```

- - -
Generated by [godoc2ghmd](https://github.com/GandalfUK/godoc2ghmd)