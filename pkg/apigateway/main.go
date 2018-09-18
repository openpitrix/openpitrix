// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package apigateway

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/pborman/uuid"
	"github.com/pkg/errors"
	"golang.org/x/tools/godoc/vfs"
	"golang.org/x/tools/godoc/vfs/httpfs"
	"golang.org/x/tools/godoc/vfs/mapfs"
	"google.golang.org/grpc"

	staticSpec "openpitrix.io/openpitrix/pkg/apigateway/spec"
	staticSwaggerUI "openpitrix.io/openpitrix/pkg/apigateway/swagger-ui"
	"openpitrix.io/openpitrix/pkg/config"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/pi"
	"openpitrix.io/openpitrix/pkg/topic"
	"openpitrix.io/openpitrix/pkg/util/senderutil"
	"openpitrix.io/openpitrix/pkg/version"
)

type Server struct {
}

type register struct {
	f        func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error)
	endpoint string
}

func Serve(cfg *config.Config) {
	version.PrintVersionInfo(func(s string, i ...interface{}) {
		logger.Info(nil, s, i...)
	})
	logger.Info(nil, "App service http://%s:%d", constants.AppManagerHost, constants.AppManagerPort)
	logger.Info(nil, "Runtime service http://%s:%d", constants.RuntimeManagerHost, constants.RuntimeManagerPort)
	logger.Info(nil, "Cluster service http://%s:%d", constants.ClusterManagerHost, constants.ClusterManagerPort)
	logger.Info(nil, "Repo service http://%s:%d", constants.RepoManagerHost, constants.RepoManagerPort)
	logger.Info(nil, "Job service http://%s:%d", constants.JobManagerHost, constants.JobManagerPort)
	logger.Info(nil, "Task service http://%s:%d", constants.TaskManagerHost, constants.TaskManagerPort)
	logger.Info(nil, "Repo indexer service http://%s:%d", constants.RepoIndexerHost, constants.RepoIndexerPort)
	logger.Info(nil, "Category service http://%s:%d", constants.CategoryManagerHost, constants.CategoryManagerPort)
	logger.Info(nil, "IAM service http://%s:%d", constants.IAMServiceHost, constants.IAMServicePort)
	logger.Info(nil, "Api service start http://%s:%d", constants.ApiGatewayHost, constants.ApiGatewayPort)

	cfg.Mysql.Disable = true
	pi.SetGlobal(cfg)
	s := Server{}

	if err := s.run(); err != nil {
		logger.Critical(nil, "Api gateway run failed: %+v", err)
		panic(err)
	}
}

const RequestIdKey = "X-Request-Id"

func log() gin.HandlerFunc {
	l := logger.NewLogger()
	l.HideCallstack()
	return func(c *gin.Context) {
		requestID := uuid.New()
		c.Request.Header.Set(RequestIdKey, requestID)
		c.Writer.Header().Set(RequestIdKey, requestID)

		t := time.Now()

		// process request
		c.Next()

		latency := time.Since(t)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		path := c.Request.URL.Path

		logStr := fmt.Sprintf("%s | %3d | %v | %s | %s %s %s",
			requestID,
			statusCode,
			latency,
			clientIP, method,
			path,
			c.Errors.String(),
		)

		switch {
		case statusCode >= 400 && statusCode <= 499:
			l.Warn(nil, logStr)
		case statusCode >= 500:
			l.Error(nil, logStr)
		default:
			l.Info(nil, logStr)
		}
	}
}

func recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				httprequest, _ := httputil.DumpRequest(c.Request, false)
				logger.Critical(nil, "Panic recovered: %+v\n%s", err, string(httprequest))
				c.JSON(500, gin.H{
					"title": "Error",
					"err":   err,
				})
			}
		}()
		c.Next() // execute all the handlers
	}
}

func handleSwagger() http.Handler {
	ns := vfs.NameSpace{}
	ns.Bind("/", mapfs.New(staticSwaggerUI.Files), "/", vfs.BindReplace)
	ns.Bind("/", mapfs.New(staticSpec.Files), "/", vfs.BindBefore)
	return http.StripPrefix("/swagger-ui", http.FileServer(httpfs.New(ns)))
}

func (s *Server) run() error {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(log())
	r.Use(recovery())
	r.Any("/swagger-ui/*filepath", gin.WrapH(handleSwagger()))
	r.Any("/v1/*filepath", gin.WrapH(s.mainHandler()))

	return r.Run(fmt.Sprintf(":%d", constants.ApiGatewayPort))
}

func (s *Server) mainHandler() http.Handler {
	var gwmux = runtime.NewServeMux(
		runtime.WithMetadata(senderutil.ServeMuxSetSender),
		runtime.WithIncomingHeaderMatcher(func(s string) (string, bool) {
			if s == RequestIdKey {
				return RequestIdKey, true
			}
			return "", false
		}),
	)
	var opts = manager.ClientOptions
	var err error

	for _, r := range []register{{
		pb.RegisterAppManagerHandlerFromEndpoint,
		fmt.Sprintf("%s:%d", constants.AppManagerHost, constants.AppManagerPort),
	}, {
		pb.RegisterCategoryManagerHandlerFromEndpoint,
		fmt.Sprintf("%s:%d", constants.CategoryManagerHost, constants.CategoryManagerPort),
	}, {
		pb.RegisterRuntimeManagerHandlerFromEndpoint,
		fmt.Sprintf("%s:%d", constants.RuntimeManagerHost, constants.RuntimeManagerPort),
	}, {
		pb.RegisterJobManagerHandlerFromEndpoint,
		fmt.Sprintf("%s:%d", constants.JobManagerHost, constants.JobManagerPort),
	}, {
		pb.RegisterTaskManagerHandlerFromEndpoint,
		fmt.Sprintf("%s:%d", constants.TaskManagerHost, constants.TaskManagerPort),
	}, {
		pb.RegisterRepoManagerHandlerFromEndpoint,
		fmt.Sprintf("%s:%d", constants.RepoManagerHost, constants.RepoManagerPort),
	}, {
		pb.RegisterRepoIndexerHandlerFromEndpoint,
		fmt.Sprintf("%s:%d", constants.RepoIndexerHost, constants.RepoIndexerPort),
	}, {
		pb.RegisterTokenManagerHandlerFromEndpoint,
		fmt.Sprintf("%s:%d", constants.IAMServiceHost, constants.IAMServicePort),
	}, {
		pb.RegisterAccountManagerHandlerFromEndpoint,
		fmt.Sprintf("%s:%d", constants.IAMServiceHost, constants.IAMServicePort),
	}, {
		pb.RegisterClusterManagerHandlerFromEndpoint,
		fmt.Sprintf("%s:%d", constants.ClusterManagerHost, constants.ClusterManagerPort),
	}} {
		err = r.f(context.Background(), gwmux, r.endpoint, opts)
		if err != nil {
			err = errors.WithStack(err)
			logger.Error(nil, "Dial [%s] failed: %+v", r.endpoint, err)
		}
	}

	mux := http.NewServeMux()
	tm := topic.NewTopicManager(pi.Global().Etcd(nil))
	go tm.Run()

	mux.Handle("/", gwmux)
	mux.HandleFunc("/v1/io", tm.HandleEvent)

	return mux
}
