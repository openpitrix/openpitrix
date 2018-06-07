// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package api

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/pkg/errors"
	"golang.org/x/tools/godoc/vfs"
	"golang.org/x/tools/godoc/vfs/httpfs"
	"golang.org/x/tools/godoc/vfs/mapfs"
	"google.golang.org/grpc"

	staticSpec "openpitrix.io/openpitrix/pkg/cmd/api/spec"
	staticSwaggerUI "openpitrix.io/openpitrix/pkg/cmd/api/swagger-ui"
	"openpitrix.io/openpitrix/pkg/config"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/util/senderutil"
	"openpitrix.io/openpitrix/pkg/version"
)

func Serve() {
	config.LoadConf()

	logger.Info("Openpitrix %s", version.ShortVersion)
	logger.Info("App service http://%s:%d", constants.AppManagerHost, constants.AppManagerPort)
	logger.Info("Runtime service http://%s:%d", constants.RuntimeManagerHost, constants.RuntimeManagerPort)
	logger.Info("Cluster service http://%s:%d", constants.ClusterManagerHost, constants.ClusterManagerPort)
	logger.Info("Repo service http://%s:%d", constants.RepoManagerHost, constants.RepoManagerPort)
	logger.Info("Job service http://%s:%d", constants.JobManagerHost, constants.JobManagerPort)
	logger.Info("Task service http://%s:%d", constants.TaskManagerHost, constants.TaskManagerPort)
	logger.Info("Repo indexer service http://%s:%d", constants.RepoIndexerHost, constants.RepoIndexerPort)
	logger.Info("Category service http://%s:%d", constants.CategoryManagerHost, constants.CategoryManagerPort)
	logger.Info("Api service start http://%s:%d", constants.ApiGatewayHost, constants.ApiGatewayPort)

	if err := run(); err != nil {
		logger.Critical("%+v", err)
		panic(err)
	}
}

func log() gin.HandlerFunc {
	l := logger.NewLogger()
	l.HideCallstack()
	return func(c *gin.Context) {
		t := time.Now()

		// process request
		c.Next()

		latency := time.Since(t)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		path := c.Request.URL.Path

		logStr := fmt.Sprintf("| %3d | %v | %s | %s %s %s",
			statusCode,
			latency,
			clientIP, method,
			path,
			c.Errors.String(),
		)

		switch {
		case statusCode >= 400 && statusCode <= 499:
			l.Warn(logStr)
		case statusCode >= 500:
			l.Error(logStr)
		default:
			l.Info(logStr)
		}
	}
}

func recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				httprequest, _ := httputil.DumpRequest(c.Request, false)
				logger.Critical("Panic recovered: %+v\n%s", err, string(httprequest))
				c.JSON(500, gin.H{
					"title": "Error",
					"err":   err,
				})
			}
		}()
		c.Next() // execute all the handlers
	}
}

func run() error {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(log())
	r.Use(recovery())

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	allHandler := gin.WrapH(mainHandler(ctx))
	r.Any("/v1/*filepath", allHandler)
	r.Any("/swagger-ui/*filepath", allHandler)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/panic", func(c *gin.Context) {
		panic("this is a panic")
	})

	return r.Run(fmt.Sprintf(":%d", constants.ApiGatewayPort))
}

func mainHandler(ctx context.Context) http.Handler {
	var gwmux = runtime.NewServeMux(runtime.WithMetadata(senderutil.ServeMuxSetSender))
	var opts = []grpc.DialOption{grpc.WithInsecure()}
	var err error

	err = pb.RegisterAppManagerHandlerFromEndpoint(
		ctx, gwmux,
		fmt.Sprintf("%s:%d", constants.AppManagerHost, constants.AppManagerPort),
		opts,
	)
	err = pb.RegisterCategoryManagerHandlerFromEndpoint(
		ctx, gwmux,
		fmt.Sprintf("%s:%d", constants.CategoryManagerHost, constants.CategoryManagerPort),
		opts,
	)
	err = pb.RegisterRuntimeManagerHandlerFromEndpoint(
		ctx, gwmux,
		fmt.Sprintf("%s:%d", constants.RuntimeManagerHost, constants.RuntimeManagerPort),
		opts,
	)
	err = pb.RegisterJobManagerHandlerFromEndpoint(
		ctx, gwmux,
		fmt.Sprintf("%s:%d", constants.JobManagerHost, constants.JobManagerPort),
		opts,
	)
	err = pb.RegisterTaskManagerHandlerFromEndpoint(
		ctx, gwmux,
		fmt.Sprintf("%s:%d", constants.TaskManagerHost, constants.TaskManagerPort),
		opts,
	)
	err = pb.RegisterRepoManagerHandlerFromEndpoint(
		ctx, gwmux,
		fmt.Sprintf("%s:%d", constants.RepoManagerHost, constants.RepoManagerPort),
		opts,
	)
	err = pb.RegisterRepoIndexerHandlerFromEndpoint(
		ctx, gwmux,
		fmt.Sprintf("%s:%d", constants.RepoIndexerHost, constants.RepoIndexerPort),
		opts,
	)
	err = pb.RegisterClusterManagerHandlerFromEndpoint(
		ctx, gwmux,
		fmt.Sprintf("%s:%d", constants.ClusterManagerHost, constants.ClusterManagerPort),
		opts,
	)
	if err != nil {
		err = errors.WithStack(err)
		logger.Error("%+v", err)
	}

	mux := http.NewServeMux()

	ns := vfs.NameSpace{}
	ns.Bind("/", mapfs.New(staticSwaggerUI.Files), "/", vfs.BindReplace)
	ns.Bind("/", mapfs.New(staticSpec.Files), "/", vfs.BindBefore)

	mux.Handle("/", gwmux)
	mux.Handle("/swagger-ui/", http.StripPrefix("/swagger-ui", http.FileServer(httpfs.New(ns))))

	return mux
}
