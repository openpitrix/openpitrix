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

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/status"

	"openpitrix.io/openpitrix/pkg/gerr"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/util/senderutil"
	"openpitrix.io/openpitrix/pkg/version"
)

type checkerT func(ctx context.Context, req interface{}) error

var defaultChecker checkerT

type GrpcServer struct {
	ServiceName    string
	Port           int
	showErrorCause bool
	checker        checkerT
}

type RegisterCallback func(*grpc.Server)

func NewGrpcServer(serviceName string, port int) *GrpcServer {
	return &GrpcServer{serviceName, port, false, defaultChecker}
}

func (g *GrpcServer) ShowErrorCause(b bool) *GrpcServer {
	g.showErrorCause = b
	return g
}

func (g *GrpcServer) WithChecker(c checkerT) *GrpcServer {
	g.checker = c
	return g
}

func (g *GrpcServer) Serve(callback RegisterCallback) {
	version.PrintVersionInfo(logger.Info)
	logger.Info("Service [%s] start listen at port [%d]", g.ServiceName, g.Port)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", g.Port))
	if err != nil {
		err = errors.WithStack(err)
		logger.Critical("failed to listen: %+v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             10 * time.Second,
			PermitWithoutStream: true,
		}),
		grpc_middleware.WithUnaryServerChain(
			grpc_validator.UnaryServerInterceptor(),
			g.unaryServerLogInterceptor(),
			grpc_recovery.UnaryServerInterceptor(
				grpc_recovery.WithRecoveryHandler(func(p interface{}) error {
					logger.Critical("GRPC server recovery with error: %+v", p)
					logger.Critical(string(debug.Stack()))
					if e, ok := p.(error); ok {
						return gerr.NewWithDetail(gerr.Internal, e, gerr.ErrorInternalError)
					}
					return gerr.New(gerr.Internal, gerr.ErrorInternalError)
				}),
			),
		),
		grpc_middleware.WithStreamServerChain(
			grpc_recovery.StreamServerInterceptor(
				grpc_recovery.WithRecoveryHandler(func(p interface{}) error {
					logger.Critical("GRPC server recovery with error: %+v", p)
					logger.Critical(string(debug.Stack()))
					if e, ok := p.(error); ok {
						return gerr.NewWithDetail(gerr.Internal, e, gerr.ErrorInternalError)
					}
					return gerr.New(gerr.Internal, gerr.ErrorInternalError)
				}),
			),
		),
	)

	callback(grpcServer)

	if err = grpcServer.Serve(lis); err != nil {
		err = errors.WithStack(err)
		logger.Critical("%+v", err)
	}
}

var (
	jsonPbMarshaller = &jsonpb.Marshaler{
		OrigName: true,
	}
)

func (g *GrpcServer) unaryServerLogInterceptor() grpc.UnaryServerInterceptor {
	showErrorCause := g.showErrorCause
	checker := g.checker

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
		var err error
		var resp interface{}
		if checker != nil {
			err = checker(ctx, req)
		}
		if err == nil {
			resp, err = handler(ctx, req)
		}
		elapsed := time.Since(start)
		logger.Info("Handled request [%s] [%+v] exec_time is [%s]", action, s, elapsed)
		if e, ok := status.FromError(err); ok {
			if e.Code() != codes.OK {
				logger.Debug("Response is error: %s, %s", e.Code().String(), e.Message())
				if !showErrorCause {
					err = gerr.ClearErrorCause(err)
				}
			}
		}
		return resp, err
	}
}
