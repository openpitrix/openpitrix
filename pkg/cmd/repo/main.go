// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package repo

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

	"openpitrix.io/openpitrix/pkg/config"
	db "openpitrix.io/openpitrix/pkg/db/repo"
	pb "openpitrix.io/openpitrix/pkg/service.pb"
)

func Main(cfg *config.Config) {
	cfg.ActiveGlogFlags()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Repo.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterRepoServiceServer(grpcServer, NewRepoServer(&cfg.DB))

	if err = grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}

type RepoServer struct {
	db *db.RepoDatabase
}

func NewRepoServer(cfg *config.Database) *RepoServer {
	db, err := db.OpenRepoDatabase(cfg)
	if err != nil {
		log.Fatal(err)
	}

	return &RepoServer{
		db: db,
	}
}

func (p *RepoServer) GetRepo(ctx context.Context, args *pb.RepoId) (reply *pb.Repo, err error) {
	result, err := p.db.GetRepo(ctx, args.GetId())
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "GetRepo: %v", err)
	}
	reply = To_proto_Repo(nil, result)
	return
}

func (p *RepoServer) GetRepoList(ctx context.Context, args *pb.RepoListRequest) (reply *pb.RepoListResponse, err error) {
	result, err := p.db.GetRepoList(ctx)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "GetRepoList: %v", err)
	}

	items := To_proto_RepoList(result, int(args.GetPageNumber()), int(args.GetPageSize()))
	reply = &pb.RepoListResponse{
		Items:       items,
		TotalItems:  proto.Int32(int32(len(result))),
		TotalPages:  proto.Int32(int32((len(result) + int(args.GetPageSize()) - 1) / int(args.GetPageSize()))),
		PageSize:    proto.Int32(args.GetPageSize()),
		CurrentPage: proto.Int32(int32(len(result)/int(args.GetPageSize())) + 1),
	}

	return
}
func (p *RepoServer) CreateRepo(ctx context.Context, args *pb.Repo) (reply *pbempty.Empty, err error) {
	err = p.db.CreateRepo(ctx, To_database_Repo(nil, args))
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "CreateRepo: %v", err)
	}

	reply = &pbempty.Empty{}
	return
}

func (p *RepoServer) UpdateRepo(ctx context.Context, args *pb.Repo) (reply *pbempty.Empty, err error) {
	err = p.db.UpdateRepo(ctx, To_database_Repo(nil, args))
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "UpdateRepo: %v", err)
	}

	reply = &pbempty.Empty{}
	return
}

func (p *RepoServer) DeleteRepo(ctx context.Context, args *pb.RepoId) (reply *pbempty.Empty, err error) {
	err = p.db.DeleteRepo(ctx, args.GetId())
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "DeleteRepo: %v", err)
	}

	reply = &pbempty.Empty{}
	return
}
