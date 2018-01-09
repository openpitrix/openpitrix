// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package runtime

import (
	"fmt"
	"net"

	"github.com/golang/protobuf/proto"
	pbempty "github.com/golang/protobuf/ptypes/empty"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"github.com/pkg/errors"
	context "golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	_ "google.golang.org/grpc/grpclog/glogger"

	config "openpitrix.io/openpitrix/pkg/config/runtime"
	db "openpitrix.io/openpitrix/pkg/db/runtime"
	"openpitrix.io/openpitrix/pkg/logger"
	pb "openpitrix.io/openpitrix/pkg/service.pb"
	"openpitrix.io/openpitrix/pkg/version"
)

func Main(cfg *config.Config) {
	cfg.Glog.ActiveFlags()

	logger.Printf("openpitrix %s\n", version.ShortVersion)
	logger.Printf("Database %s://tcp(%s:%d)/%s\n", cfg.DB.Type, cfg.DB.Host, cfg.DB.Port, cfg.DB.DbName)
	logger.Printf("Runtime service start http://%s:%d\n", cfg.Runtime.Host, cfg.Runtime.Port)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Runtime.Port))
	if err != nil {
		err = errors.WithStack(err)
		logger.Fatalf("failed to listen: %+v", err)
	}

	grpcServer := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_validator.UnaryServerInterceptor(),
			grpc_recovery.UnaryServerInterceptor(
				grpc_recovery.WithRecoveryHandler(func(p interface{}) error {
					return grpc.Errorf(codes.Internal, "%+v", p)
				}),
			),
		),
		grpc_middleware.WithStreamServerChain(
			grpc_recovery.StreamServerInterceptor(
				grpc_recovery.WithRecoveryHandler(func(p interface{}) error {
					return grpc.Errorf(codes.Internal, "%+v", p)
				}),
			),
		),
	)
	pb.RegisterAppRuntimeServiceServer(grpcServer, NewAppRuntimeServer(&cfg.DB))

	if err = grpcServer.Serve(lis); err != nil {
		err = errors.WithStack(err)
		logger.Fatalf("%+v", err)
	}
}

type AppRuntimeServer struct {
	db *db.AppRuntimeDatabase
}

func NewAppRuntimeServer(cfg *config.RuntimeDatabase) *AppRuntimeServer {
	db, err := db.OpenAppRuntimeDatabase(cfg)
	if err != nil {
		logger.Fatalf("%+v", err)
	}

	return &AppRuntimeServer{
		db: db,
	}
}

func (p *AppRuntimeServer) GetAppRuntime(ctx context.Context, args *pb.AppRuntimeId) (reply *pb.AppRuntime, err error) {
	if id := args.GetId(); id == "rt-panic000" {
		panic(id) // only for test
	}

	result, err := p.db.GetAppRuntime(ctx, args.GetId())
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "GetAppRuntime: %+v", err)
	}
	if result == nil {
		return nil, grpc.Errorf(codes.NotFound, "App Runtime Id %s does not exist", args.GetId())
	}
	reply = To_proto_AppRuntime(nil, result)
	return
}

func (p *AppRuntimeServer) GetAppRuntimeList(ctx context.Context, args *pb.AppRuntimeListRequest) (reply *pb.AppRuntimeListResponse, err error) {
	result, err := p.db.GetAppRuntimeList(ctx)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "GetAppRuntimeList: %+v", err)
	}

	items := To_proto_AppRuntimeList(result, int(args.GetPageNumber()), int(args.GetPageSize()))
	reply = &pb.AppRuntimeListResponse{
		Items:       items,
		TotalItems:  proto.Int32(int32(len(result))),
		TotalPages:  proto.Int32(int32((len(result) + int(args.GetPageSize()) - 1) / int(args.GetPageSize()))),
		PageSize:    proto.Int32(args.GetPageSize()),
		CurrentPage: proto.Int32(int32(len(result)/int(args.GetPageSize())) + 1),
	}

	return
}

func (p *AppRuntimeServer) CreateAppRuntime(ctx context.Context, args *pb.AppRuntime) (reply *pbempty.Empty, err error) {
	err = p.db.CreateAppRuntime(ctx, To_database_AppRuntime(nil, args))
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "CreateAppRuntime: %+v", err)
	}

	reply = &pbempty.Empty{}
	return
}

func (p *AppRuntimeServer) UpdateAppRuntime(ctx context.Context, args *pb.AppRuntime) (reply *pbempty.Empty, err error) {
	err = p.db.UpdateAppRuntime(ctx, To_database_AppRuntime(nil, args))
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "UpdateAppRuntime: %+v", err)
	}

	reply = &pbempty.Empty{}
	return
}

func (p *AppRuntimeServer) DeleteAppRuntime(ctx context.Context, args *pb.AppRuntimeId) (reply *pbempty.Empty, err error) {
	err = p.db.DeleteAppRuntime(ctx, args.GetId())
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "DeleteAppRuntime: %+v", err)
	}

	reply = &pbempty.Empty{}
	return
}
