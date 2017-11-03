// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package app

import (
	"fmt"
	"log"
	"net"

	pbempty "github.com/golang/protobuf/ptypes/empty"
	context "golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"openpitrix.io/openpitrix/pkg/config"
	db "openpitrix.io/openpitrix/pkg/db/app"
	pb "openpitrix.io/openpitrix/pkg/service.pb"
)

func Main(cfg *config.Config) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.AppService.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAppServiceServer(grpcServer, NewAppServer(&cfg.Database))

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
	result, err := p.db.GetApp(ctx, args.Id)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "GetApp: %v", err)
	}
	reply = To_proto_App(nil, result)
	return
}

func (p *AppServer) GetAppList(ctx context.Context, args *pb.AppListRequest) (reply *pb.AppListResponse, err error) {
	if args.PageNumber <= 0 {
		args.PageNumber = 1
	}
	if args.PageSize <= 0 {
		args.PageSize = 10
	}

	result, err := p.db.GetAppList(ctx)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "GetAppList: %v", err)
	}

	items := To_proto_AppList(result, int(args.PageNumber), int(args.PageSize))
	reply = &pb.AppListResponse{
		Items:       items,
		TotalItems:  int32(len(result)),
		TotalPages:  int32((len(result) + int(args.PageSize) - 1) / int(args.PageSize)),
		PageSize:    args.PageSize,
		CurrentPage: int32(len(result)/int(args.PageSize)) + 1,
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
	err = p.db.DeleteApp(ctx, args.Id)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "DeleteApp: %v", err)
	}

	reply = &pbempty.Empty{}
	return
}
