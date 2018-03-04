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

	"openpitrix.io/openpitrix/pkg/config"
	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/version"
)

type GrpcServer struct {
	ServiceName string
	Port        int
	MysqlConfig config.MysqlConfig
}

type RegisterCallback func(*grpc.Server, *db.Database)

func (g *GrpcServer) OpenDatabase() *db.Database {
	dbSession, err := db.OpenDatabase(g.MysqlConfig)
	if err != nil {
		err = errors.WithStack(err)
		logger.Fatalf("failed to connect mysql: %+v", err)
	}
	return dbSession
}

func (g *GrpcServer) Serve(callback RegisterCallback) {
	dbSession := g.OpenDatabase()
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

	callback(grpcServer, dbSession)

	if err = grpcServer.Serve(lis); err != nil {
		err = errors.WithStack(err)
		logger.Fatalf("%+v", err)
	}
}
