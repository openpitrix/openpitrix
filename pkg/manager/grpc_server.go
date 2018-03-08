// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package manager

import (
	"fmt"
	"net"
	"runtime/debug"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/version"
)

type GrpcServer struct {
	ServiceName string
	Port        int
}

type RegisterCallback func(*grpc.Server)

func NewGrpcServer(serviceName string, port int) *GrpcServer {
	return &GrpcServer{serviceName, port}
}

func (g *GrpcServer) Serve(callback RegisterCallback) {
	logger.Infof("Openpitrix %s\n", version.ShortVersion)
	logger.Infof("Service [%s] start listen at port [%d]\n", g.ServiceName, g.Port)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", g.Port))
	if err != nil {
		err = errors.WithStack(err)
		logger.Fatalf("failed to listen: %+v", err)
	}

	grpcServer := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_validator.UnaryServerInterceptor(),
			grpc_recovery.UnaryServerInterceptor(
				grpc_recovery.WithRecoveryHandler(func(p interface{}) error {
					logger.Panic(p)
					logger.Panic(string(debug.Stack()))
					return status.Errorf(codes.Internal, "%+v", p)
				}),
			),
		),
		grpc_middleware.WithStreamServerChain(
			grpc_recovery.StreamServerInterceptor(
				grpc_recovery.WithRecoveryHandler(func(p interface{}) error {
					logger.Panic(p)
					logger.Panic(string(debug.Stack()))
					return status.Errorf(codes.Internal, "%+v", p)
				}),
			),
		),
	)

	callback(grpcServer)

	if err = grpcServer.Serve(lis); err != nil {
		err = errors.WithStack(err)
		logger.Fatalf("%+v", err)
	}
}
