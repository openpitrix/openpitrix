// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package runtime

import (
	"fmt"
	"log"
	"net"

	"github.com/golang/protobuf/proto"
	pbempty "github.com/golang/protobuf/ptypes/empty"
	context "golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"openpitrix.io/openpitrix/pkg/config"
	db "openpitrix.io/openpitrix/pkg/db/runtime"
	pb "openpitrix.io/openpitrix/pkg/service.pb"
)

func Main(cfg *config.Config) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.AppRuntimeService.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAppRuntimeServiceServer(grpcServer, NewAppRuntimeServer(&cfg.Database))

	if err = grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}

type AppRuntimeServer struct {
	db *db.AppRuntimeDatabase
}

func NewAppRuntimeServer(cfg *config.Database) *AppRuntimeServer {
	db, err := db.OpenAppRuntimeDatabase(cfg)
	if err != nil {
		log.Fatal(err)
	}

	return &AppRuntimeServer{
		db: db,
	}
}

func (p *AppRuntimeServer) GetAppRuntime(ctx context.Context, args *pb.AppRuntimeId) (reply *pb.AppRuntime, err error) {
	result, err := p.db.GetAppRuntime(ctx, args.GetId())
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "GetAppRuntime: %v", err)
	}
	reply = To_proto_AppRuntime(nil, result)
	return
}

func (p *AppRuntimeServer) GetAppRuntimeList(ctx context.Context, args *pb.AppRuntimeListRequest) (reply *pb.AppRuntimeListResponse, err error) {
	result, err := p.db.GetAppRuntimeList(ctx)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "GetAppRuntimeList: %v", err)
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
		return nil, grpc.Errorf(codes.Internal, "CreateAppRuntime: %v", err)
	}

	reply = &pbempty.Empty{}
	return
}

func (p *AppRuntimeServer) UpdateAppRuntime(ctx context.Context, args *pb.AppRuntime) (reply *pbempty.Empty, err error) {
	err = p.db.UpdateAppRuntime(ctx, To_database_AppRuntime(nil, args))
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "UpdateAppRuntime: %v", err)
	}

	reply = &pbempty.Empty{}
	return
}

func (p *AppRuntimeServer) DeleteAppRuntime(ctx context.Context, args *pb.AppRuntimeId) (reply *pbempty.Empty, err error) {
	err = p.db.DeleteAppRuntime(ctx, args.GetId())
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "DeleteAppRuntime: %v", err)
	}

	reply = &pbempty.Empty{}
	return
}
