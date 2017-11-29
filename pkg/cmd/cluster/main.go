// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package cluster

import (
	"context"
	"fmt"
	"net"

	"github.com/golang/protobuf/proto"
	pbempty "github.com/golang/protobuf/ptypes/empty"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"github.com/pkg/errors"
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
	pb.RegisterClusterServiceServer(grpcServer, NewClusterServer(&cfg.DB))

	if err = grpcServer.Serve(lis); err != nil {
		err = errors.WithStack(err)
		logger.Fatalf("%+v", err)
	}
}

type ClusterServer struct {
	db *db.ClusterDatabase
}

func NewClusterServer(cfg *config.Database) *ClusterServer {
	cluster_db, err := db.OpenClusterDatabase(cfg)
	if err != nil {
		logger.Fatalf("%+v", err)
	}

	return &ClusterServer{
		db: cluster_db,
	}
}

func (p *ClusterServer) GetClusters(ctx context.Context, args *pb.ClusterIds) (reply *pb.Clusters, err error) {
	result, err := p.db.GetClusters(ctx, args.GetIds())
	if err != nil {
		err = errors.WithStack(err)
		return nil, grpc.Errorf(codes.Internal, "GetClusters: %+v", err)
	}
	if result == nil {
		err = errors.WithStack(err)
		return nil, grpc.Errorf(codes.NotFound, "Cluster Ids %s do not exist", args.GetIds())
	}
	reply = To_proto_Clusters(result)
	return
}

func (p *ClusterServer) GetClusterList(ctx context.Context, args *pb.ClusterListRequest) (reply *pb.ClusterListResponse, err error) {
	result, err := p.db.GetClusterList(ctx)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "GetClusterList: %+v", err)
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
		return nil, grpc.Errorf(codes.Internal, "CreateCluster: %+v", err)
	}

	reply = &pbempty.Empty{}
	return
}

func (p *ClusterServer) UpdateCluster(ctx context.Context, args *pb.Cluster) (reply *pbempty.Empty, err error) {
	err = p.db.UpdateCluster(ctx, To_database_Cluster(nil, args))
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "UpdateCluster: %+v", err)
	}

	reply = &pbempty.Empty{}
	return
}

func (p *ClusterServer) DeleteClusters(ctx context.Context, args *pb.ClusterIds) (reply *pbempty.Empty, err error) {
	err = p.db.DeleteClusters(ctx, args.GetIds())
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "DeleteClusters: %+v", err)
	}

	reply = &pbempty.Empty{}
	return
}

func (p *ClusterServer) GetClusterNodes(ctx context.Context, args *pb.ClusterNodeIds) (reply *pb.ClusterNodes, err error) {
	result, err := p.db.GetClusterNodes(ctx, args.GetIds())
	if err != nil {
		err = errors.WithStack(err)
		return nil, grpc.Errorf(codes.Internal, "GetClusterNodes: %+v", err)
	}
	if result == nil {
		err = errors.WithStack(err)
		return nil, grpc.Errorf(codes.NotFound, "ClusterNode Ids %s do not exist", args.GetIds())
	}
	reply = To_proto_ClusterNodes(result)
	return
}

func (p *ClusterServer) GetClusterNodeList(ctx context.Context, args *pb.ClusterNodeListRequest) (reply *pb.ClusterNodeListResponse, err error) {
	result, err := p.db.GetClusterNodeList(ctx)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "GetClusterNodeList: %+v", err)
	}

	items := To_proto_ClusterNodeList(result, int(args.GetPageNumber()), int(args.GetPageSize()))
	reply = &pb.ClusterNodeListResponse{
		Items:       items,
		TotalItems:  proto.Int32(int32(len(result))),
		TotalPages:  proto.Int32(int32((len(result) + int(args.GetPageSize()) - 1) / int(args.GetPageSize()))),
		PageSize:    proto.Int32(args.GetPageSize()),
		CurrentPage: proto.Int32(int32(len(result)/int(args.GetPageSize())) + 1),
	}

	return
}

func (p *ClusterServer) CreateClusterNodes(ctx context.Context, args *pb.ClusterNodes) (reply *pbempty.Empty, err error) {
	err = p.db.CreateClusterNodes(ctx, To_database_ClusterNodes(args))
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "CreateClusterNodes: %+v", err)
	}

	reply = &pbempty.Empty{}
	return
}

func (p *ClusterServer) UpdateClusterNode(ctx context.Context, args *pb.ClusterNode) (reply *pbempty.Empty, err error) {
	err = p.db.UpdateClusterNode(ctx, To_database_ClusterNode(nil, args))
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "UpdateClusterNode: %+v", err)
	}

	reply = &pbempty.Empty{}
	return
}

func (p *ClusterServer) DeleteClusterNodes(ctx context.Context, args *pb.ClusterNodeIds) (reply *pbempty.Empty, err error) {
	err = p.db.DeleteClusterNodes(ctx, args.GetIds())
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "DeleteClusterNodes: %+v", err)
	}

	reply = &pbempty.Empty{}
	return
}
