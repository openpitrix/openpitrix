// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/ekyoung/gin-nice-recovery"
	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/pkg/errors"
	"github.com/szuecs/gin-glog"
	"golang.org/x/tools/godoc/vfs"
	"golang.org/x/tools/godoc/vfs/httpfs"
	"golang.org/x/tools/godoc/vfs/mapfs"
	"google.golang.org/grpc"
	_ "google.golang.org/grpc/grpclog/glogger"

	config "openpitrix.io/openpitrix/pkg/config/api"
	"openpitrix.io/openpitrix/pkg/logger"
	pb "openpitrix.io/openpitrix/pkg/service.pb"
	"openpitrix.io/openpitrix/pkg/version"

	staticSpec "openpitrix.io/openpitrix/pkg/cmd/api/spec"
	staticSwaggerUI "openpitrix.io/openpitrix/pkg/cmd/api/swagger-ui"
)

func Main(cfg *config.Config) {
	cfg.Glog.ActiveFlags()

	logger.Printf("openpitrix %s\n", version.ShortVersion)
	logger.Printf("App service http://%s:%d\n", cfg.App.Host, cfg.App.Port)
	logger.Printf("Runtime service http://%s:%d\n", cfg.Runtime.Host, cfg.Runtime.Port)
	logger.Printf("Cluster service http://%s:%d\n", cfg.Cluster.Host, cfg.Cluster.Port)
	logger.Printf("Repo service http://%s:%d\n", cfg.Repo.Host, cfg.Repo.Port)
	logger.Printf("Api service start http://%s:%d\n", cfg.Api.Host, cfg.Api.Port)

	if err := run(cfg); err != nil {
		logger.Fatalf("%+v", err)
	}
}

func run(cfg *config.Config) error {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(ginglog.Logger(3 * time.Second))
	r.Use(nice.Recovery(func(c *gin.Context, err interface{}) {
		c.JSON(500, gin.H{
			"title": "Error",
			"err":   err,
		})
	}))

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	allHandler := gin.WrapH(mainHandler(cfg, ctx))
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

	return r.Run(fmt.Sprintf(":%d", cfg.Api.Port))
}

func mainHandler(cfg *config.Config, ctx context.Context) http.Handler {
	var gwmux = runtime.NewServeMux()
	var opts = []grpc.DialOption{grpc.WithInsecure()}
	var err error

	err = pb.RegisterAppServiceHandlerFromEndpoint(
		ctx, gwmux,
		fmt.Sprintf("%s:%d", cfg.App.Host, cfg.App.Port),
		opts,
	)
	if err != nil {
		err = errors.WithStack(err)
		logger.Fatalf("%+v", err)
	}

	err = pb.RegisterAppRuntimeServiceHandlerFromEndpoint(
		ctx, gwmux,
		fmt.Sprintf("%s:%d", cfg.Runtime.Host, cfg.Runtime.Port),
		opts,
	)
	if err != nil {
		err = errors.WithStack(err)
		logger.Fatalf("%+v", err)
	}

	err = pb.RegisterClusterServiceHandlerFromEndpoint(
		ctx, gwmux,
		fmt.Sprintf("%s:%d", cfg.Cluster.Host, cfg.Cluster.Port),
		opts,
	)
	if err != nil {
		err = errors.WithStack(err)
		logger.Fatalf("%+v", err)
	}

	err = pb.RegisterRepoServiceHandlerFromEndpoint(
		ctx, gwmux,
		fmt.Sprintf("%s:%d", cfg.Repo.Host, cfg.Repo.Port),
		opts,
	)
	if err != nil {
		err = errors.WithStack(err)
		logger.Fatalf("%+v", err)
	}

	mux := http.NewServeMux()

	ns := vfs.NameSpace{}
	ns.Bind("/", mapfs.New(staticSwaggerUI.Files), "/", vfs.BindReplace)
	ns.Bind("/", mapfs.New(staticSpec.Files), "/", vfs.BindBefore)

	mux.Handle("/", gwmux)
	mux.Handle("/swagger-ui/", http.StripPrefix("/swagger-ui", http.FileServer(httpfs.New(ns))))

	return mux
}
