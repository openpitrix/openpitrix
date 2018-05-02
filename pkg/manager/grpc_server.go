// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package manager

import (
	"context"
	"fmt"
	"net"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/golang/protobuf/jsonpb"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/util/senderutil"
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
	logger.Info("Openpitrix %s", version.ShortVersion)
	logger.Info("Service [%s] start listen at port [%d]", g.ServiceName, g.Port)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", g.Port))
	if err != nil {
		err = errors.WithStack(err)
		logger.Fatalf("failed to listen: %+v", err)
	}

	grpcServer := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_validator.UnaryServerInterceptor(),
			UnaryServerLogInterceptor(),
			grpc_recovery.UnaryServerInterceptor(
				grpc_recovery.WithRecoveryHandler(func(p interface{}) error {
					logger.Critical("GRPC server recovery with error: %+v", p)
					logger.Critical(string(debug.Stack()))
					return status.Errorf(codes.Internal, "%+v", p)
				}),
			),
		),
		grpc_middleware.WithStreamServerChain(
			grpc_recovery.StreamServerInterceptor(
				grpc_recovery.WithRecoveryHandler(func(p interface{}) error {
					logger.Critical("GRPC server recovery with error: %+v", p)
					logger.Critical(string(debug.Stack()))
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

var (
	jsonPbMarshaller = &jsonpb.Marshaler{}
)

func UnaryServerLogInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		s := senderutil.GetSenderFromContext(ctx)
		method := strings.Split(info.FullMethod, "/")
		action := method[len(method)-1]
		if p, ok := req.(proto.Message); ok {
			if content, err := jsonPbMarshaller.MarshalToString(p); err != nil {
				logger.Error("Failed to marshal proto message to string [%s] [%+v] [%+v]", action, s, err)
			} else {
				logger.Info("Request received [%s] [%+v] [%s]", action, s, content)
			}
		}
		start := time.Now()
		resp, err := handler(ctx, req)
		elapsed := time.Since(start)
		logger.Info("Handled request [%s] [%+v] exec_time is [%s]", action, s, elapsed)
		return resp, err
	}
}
