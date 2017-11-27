// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package app

import (
	"fmt"
	"log"
	"net"

	"github.com/golang/protobuf/proto"
	pbempty "github.com/golang/protobuf/ptypes/empty"
	context "golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	_ "google.golang.org/grpc/grpclog/glogger"

	"github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"openpitrix.io/openpitrix/pkg/config"
	db "openpitrix.io/openpitrix/pkg/db/app"
	"openpitrix.io/openpitrix/pkg/logger"
	pb "openpitrix.io/openpitrix/pkg/service.pb"
)

func Main(cfg *config.Config) {
	cfg.ActiveGlogFlags()

	logger.Printf("Database %s://tcp(%s:%d)/%s\n", cfg.DB.Type, cfg.DB.Host, cfg.DB.Port, cfg.DB.DbName)
	logger.Printf("App service start http://%s:%d\n", cfg.App.Host, cfg.App.Port)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.App.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var opts = []grpc.ServerOption{grpc.UnaryInterceptor(grpc_validator.UnaryServerInterceptor())}
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterAppServiceServer(grpcServer, NewAppServer(&cfg.DB))

	if err = grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}

type AppServer struct {
	db *db.AppDatabase
}

func NewAppServer(cfg *config.Database) *AppServer {
	db, err := db.OpenAppDatabase(cfg)
	if err != nil {
		log.Fatal(err)
	}

	return &AppServer{
		db: db,
	}
}

func (p *AppServer) GetApp(ctx context.Context, args *pb.AppId) (reply *pb.App, err error) {
	result, err := p.db.GetApp(ctx, args.GetId())
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "GetApp: %v", err)
	}
	if result == nil {
		return nil, grpc.Errorf(codes.NotFound, "App Id %s dose not exist", args.GetId())
	}
	reply = To_proto_App(nil, result)
	return
}

func (p *AppServer) GetAppList(ctx context.Context, args *pb.AppListRequest) (reply *pb.AppListResponse, err error) {
	result, err := p.db.GetAppList(ctx)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "GetAppList: %v", err)
	}

	items := To_proto_AppList(result, int(args.GetPageNumber()), int(args.GetPageSize()))
	reply = &pb.AppListResponse{
		Items:       items,
		TotalItems:  proto.Int32(int32(len(result))),
		TotalPages:  proto.Int32(int32((len(result) + int(args.GetPageSize()) - 1) / int(args.GetPageSize()))),
		PageSize:    proto.Int32(args.GetPageSize()),
		CurrentPage: proto.Int32(int32(len(result)/int(args.GetPageSize())) + 1),
	}

	return
}

func (p *AppServer) CreateApp(ctx context.Context, args *pb.App) (reply *pbempty.Empty, err error) {
	err = p.db.CreateApp(ctx, To_database_App(nil, args))
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "CreateApp: %v", err)
	}

	reply = &pbempty.Empty{}
	return
}

func (p *AppServer) UpdateApp(ctx context.Context, args *pb.App) (reply *pbempty.Empty, err error) {
	err = p.db.UpdateApp(ctx, To_database_App(nil, args))
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "UpdateApp: %v", err)
	}

	reply = &pbempty.Empty{}
	return
}

func (p *AppServer) DeleteApp(ctx context.Context, args *pb.AppId) (reply *pbempty.Empty, err error) {
	err = p.db.DeleteApp(ctx, args.GetId())
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "DeleteApp: %v", err)
	}

	reply = &pbempty.Empty{}
	return
}
