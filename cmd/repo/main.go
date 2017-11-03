// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// openpitrix repo server
package main

import (
	"log"

	"github.com/go-openapi/loads"

	"openpitrix.io/openpitrix/pkg/config"
	"openpitrix.io/openpitrix/pkg/swagger/restapi"
	"openpitrix.io/openpitrix/pkg/swagger/restapi/operations"
)

func main() {
	cfg := config.MustLoadUserConfig()
	ListenAndServe(cfg)
}

type RestServer struct {
	*restapi.Server
	Spec *loads.Document
	Api  *operations.OpenPitrixAPI
	Cfg  *config.Config
}

func NewRestServer(config *config.Config) *RestServer {
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		log.Fatal(err)
	}

	p := new(RestServer)
	p.Api = operations.NewOpenPitrixAPI(swaggerSpec)
	p.Spec = swaggerSpec
	p.Server = restapi.NewServer(p.Api)
	p.Cfg = config.Clone()

	return p
}

func (p *RestServer) ConfigureFlags() {
	p.Server.ConfigureFlags()

	p.Server.Host = "0.0.0.0"
	p.Server.Port = p.Cfg.RepoService.Port
}

func (p *RestServer) ConfigureAPI() {
	// TODO
}

func (p *RestServer) Serve() error {
	p.ConfigureFlags()
	p.ConfigureAPI()

	return p.Server.Serve()
}

func ListenAndServe(config *config.Config) error {
	server := NewRestServer(config)
	defer server.Shutdown()

	if err := server.Serve(); err != nil {
		return err
	}
	return nil
}
