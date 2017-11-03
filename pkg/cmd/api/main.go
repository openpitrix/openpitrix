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

	"openpitrix.io/openpitrix/pkg/config"
	pb "openpitrix.io/openpitrix/pkg/service.pb"

	staticSpec "openpitrix.io/openpitrix/pkg/cmd/api/spec"
	staticSwaggerUI "openpitrix.io/openpitrix/pkg/cmd/api/swagger-ui"
)

func Main(cfg *config.Config) {
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
		fmt.Sprintf("%s:%d", cfg.AppService.Host, cfg.AppService.Port),
		opts,
	)
	if err != nil {
		return err
	}

	err = pb.RegisterAppRuntimeServiceHandlerFromEndpoint(
		ctx, gwmux,
		fmt.Sprintf("%s:%d", cfg.AppRuntimeService.Host, cfg.AppRuntimeService.Port),
		opts,
	)
	if err != nil {
		return err
	}

	err = pb.RegisterClusterServiceHandlerFromEndpoint(
		ctx, gwmux,
		fmt.Sprintf("%s:%d", cfg.ClusterService.Host, cfg.ClusterService.Port),
		opts,
	)
	if err != nil {
		return err
	}

	err = pb.RegisterRepoServiceHandlerFromEndpoint(
		ctx, gwmux,
		fmt.Sprintf("%s:%d", cfg.RepoService.Host, cfg.RepoService.Port),
		opts,
	)
	if err != nil {
		return err
	}

	mux := http.NewServeMux()

	ns := vfs.NameSpace{}
	ns.Bind("/", mapfs.New(staticSwaggerUI.Files), "/", vfs.BindReplace)
	ns.Bind("/", mapfs.New(staticSpec.Files), "/", vfs.BindAfter)

	mux.Handle("/", gwmux)
	mux.Handle("/swagger-ui/", http.StripPrefix("/swagger-ui", http.FileServer(httpfs.New(ns))))

	return http.ListenAndServe(fmt.Sprintf(":%d", cfg.ApiService.Port), mux)
}
