// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/context"
	"golang.org/x/tools/godoc/vfs"
	"golang.org/x/tools/godoc/vfs/httpfs"
	"golang.org/x/tools/godoc/vfs/mapfs"
	"google.golang.org/grpc"
	_ "google.golang.org/grpc/grpclog/glogger"

	"openpitrix.io/openpitrix/pkg/config"
	"openpitrix.io/openpitrix/pkg/logger"
	pb "openpitrix.io/openpitrix/pkg/service.pb"

	staticSpec "openpitrix.io/openpitrix/pkg/cmd/api/spec"
	staticSwaggerUI "openpitrix.io/openpitrix/pkg/cmd/api/swagger-ui"
)

func Main(cfg *config.Config) {
	cfg.ActiveGlogFlags()

	logger.Printf("Database %s://tcp(%s:%d)/%s\n", cfg.DB.Type, cfg.DB.Host, cfg.DB.Port, cfg.DB.DbName)
	logger.Printf("App service http://%s:%d\n", cfg.App.Host, cfg.App.Port)
	logger.Printf("Runtime service http://%s:%d\n", cfg.Runtime.Host, cfg.Runtime.Port)
	logger.Printf("Cluster service http://%s:%d\n", cfg.Cluster.Host, cfg.Cluster.Port)
	logger.Printf("Repo service http://%s:%d\n", cfg.Repo.Host, cfg.Repo.Port)
	logger.Printf("Api service start http://%s:%d\n", cfg.Api.Host, cfg.Api.Port)

	if err := run(cfg); err != nil {
		log.Fatal(err)
	}
}

func run(cfg *config.Config) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var gwmux = runtime.NewServeMux()
	var opts = []grpc.DialOption{grpc.WithInsecure()}
	var err error

	err = pb.RegisterAppServiceHandlerFromEndpoint(
		ctx, gwmux,
		fmt.Sprintf("%s:%d", cfg.App.Host, cfg.App.Port),
		opts,
	)
	if err != nil {
		return err
	}

	err = pb.RegisterAppRuntimeServiceHandlerFromEndpoint(
		ctx, gwmux,
		fmt.Sprintf("%s:%d", cfg.Runtime.Host, cfg.Runtime.Port),
		opts,
	)
	if err != nil {
		return err
	}

	err = pb.RegisterClusterServiceHandlerFromEndpoint(
		ctx, gwmux,
		fmt.Sprintf("%s:%d", cfg.Cluster.Host, cfg.Cluster.Port),
		opts,
	)
	if err != nil {
		return err
	}

	err = pb.RegisterRepoServiceHandlerFromEndpoint(
		ctx, gwmux,
		fmt.Sprintf("%s:%d", cfg.Repo.Host, cfg.Repo.Port),
		opts,
	)
	if err != nil {
		return err
	}

	mux := http.NewServeMux()

	ns := vfs.NameSpace{}
	ns.Bind("/", mapfs.New(staticSwaggerUI.Files), "/", vfs.BindReplace)
	ns.Bind("/", mapfs.New(staticSpec.Files), "/", vfs.BindBefore)

	mux.Handle("/", gwmux)
	mux.Handle("/swagger-ui/", http.StripPrefix("/swagger-ui", http.FileServer(httpfs.New(ns))))

	return http.ListenAndServe(fmt.Sprintf(":%d", cfg.Api.Port), mux)
}
