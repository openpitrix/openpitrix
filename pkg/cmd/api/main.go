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

	staticSpec "openpitrix.io/openpitrix/pkg/cmd/api/spec"
	staticSwaggerUI "openpitrix.io/openpitrix/pkg/cmd/api/swagger-ui"
	"openpitrix.io/openpitrix/pkg/config"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/utils/sender"
	"openpitrix.io/openpitrix/pkg/version"
)

func Serve() {
	config.LoadConf()

	logger.Infof("Openpitrix %s\n", version.ShortVersion)
	logger.Infof("App service http://%s:%s\n", constants.AppManagerHost, constants.AppManagerPort)
	logger.Infof("Runtime service http://%s:%s\n", constants.RuntimeEnvManagerHost, constants.RuntimeEnvManagerPort)
	logger.Infof("Cluster service http://%s:%s\n", constants.ClusterManagerHost, constants.ClusterManagerPort)
	logger.Infof("Repo service http://%s:%s\n", constants.RepoManagerHost, constants.RepoManagerPort)
	logger.Infof("Api service start http://%s:%s\n", constants.ApiGatewayHost, constants.ApiGatewayPort)

	if err := run(); err != nil {
		logger.Fatalf("%+v", err)
		panic(err)
	}
}

func run() error {
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

	return r.Run(":" + constants.ApiGatewayPort)
}

func mainHandler(ctx context.Context) http.Handler {
	var gwmux = runtime.NewServeMux(runtime.WithMetadata(sender.ServeMuxSetSender))
	var opts = []grpc.DialOption{grpc.WithInsecure()}
	var err error

	err = pb.RegisterAppManagerHandlerFromEndpoint(
		ctx, gwmux,
		fmt.Sprintf("%s:%s", constants.AppManagerHost, constants.AppManagerPort),
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
