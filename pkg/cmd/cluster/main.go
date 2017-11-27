// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package cluster

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
	db "openpitrix.io/openpitrix/pkg/db/cluster"
	"openpitrix.io/openpitrix/pkg/logger"
	pb "openpitrix.io/openpitrix/pkg/service.pb"
	"openpitrix.io/openpitrix/pkg/version"
)

func Main(cfg *config.Config) {
	cfg.ActiveGlogFlags()

	logger.Printf("openpitrix %s\n", version.ShortVersion)

	logger.Printf("Database %s://tcp(%s:%d)/%s\n", cfg.DB.Type, cfg.DB.Host, cfg.DB.Port, cfg.DB.DbName)
	logger.Printf("Cluster service start http://%s:%d\n", cfg.Cluster.Host, cfg.Cluster.Port)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Cluster.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterClusterServiceServer(grpcServer, NewClusterServer(&cfg.DB))

	if err = grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}

type ClusterServer struct {
	db *db.ClusterDatabase
}

func NewClusterServer(cfg *config.Database) *ClusterServer {
	db, err := db.OpenClusterDatabase(cfg)
	if err != nil {
		log.Fatal(err)
	}

	return &ClusterServer{
		db: db,
	}
}

func (p *ClusterServer) GetCluster(ctx context.Context, args *pb.ClusterId) (reply *pb.Cluster, err error) {
	result, err := p.db.GetCluster(ctx, args.GetId())
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "GetCluster: %v", err)
	}
	if result == nil {
		return nil, grpc.Errorf(codes.NotFound, "Cluster Id %s does not exist", args.GetId())
	}
	reply = To_proto_Cluster(nil, result)
	return
}

func (p *ClusterServer) GetClusterList(ctx context.Context, args *pb.ClusterListRequest) (reply *pb.ClusterListResponse, err error) {
	result, err := p.db.GetClusterList(ctx)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "GetClusterList: %v", err)
	}

	items := To_proto_ClusterList(result, int(args.GetPageNumber()), int(args.GetPageSize()))
	reply = &pb.ClusterListResponse{
		Items:       items,
		TotalItems:  proto.Int32(int32(len(result))),
		TotalPages:  proto.Int32(int32((len(result) + int(args.GetPageSize()) - 1) / int(args.GetPageSize()))),
		PageSize:    proto.Int32(args.GetPageSize()),
		CurrentPage: proto.Int32(int32(len(result)/int(args.GetPageSize())) + 1),
	}

	return
}

func (p *ClusterServer) CreateCluster(ctx context.Context, args *pb.Cluster) (reply *pbempty.Empty, err error) {
	err = p.db.CreateCluster(ctx, To_database_Cluster(nil, args))
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "CreateCluster: %v", err)
	}

	reply = &pbempty.Empty{}
	return
}

func (p *ClusterServer) UpdateCluster(ctx context.Context, args *pb.Cluster) (reply *pbempty.Empty, err error) {
	err = p.db.UpdateCluster(ctx, To_database_Cluster(nil, args))
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "UpdateCluster: %v", err)
	}

	reply = &pbempty.Empty{}
	return
}

func (p *ClusterServer) DeleteCluster(ctx context.Context, args *pb.ClusterId) (reply *pbempty.Empty, err error) {
	err = p.db.DeleteCluster(ctx, args.GetId())
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "DeleteCluster: %v", err)
	}

	reply = &pbempty.Empty{}
	return
}
