// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package apps

import (
	"log"

	"github.com/go-openapi/loads"

	"openpitrix.io/openpitrix/pkg/config"
	"openpitrix.io/openpitrix/pkg/swagger/restapi"
	"openpitrix.io/openpitrix/pkg/swagger/restapi/operations"
)

type AppsServer struct {
	*restapi.Server
	Spec *loads.Document
	Api  *operations.OpenPitrixAPI
	Cfg  *config.Config

	service *AppsRestService
}

func NewAppsServer(config *config.Config) *AppsServer {
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		log.Fatal(err)
	}

	p := new(AppsServer)
	p.Api = operations.NewOpenPitrixAPI(swaggerSpec)
	p.Spec = swaggerSpec
	p.Server = restapi.NewServer(p.Api)
	p.Cfg = config.Clone()

	p.service = NewAppsRestService(config.Clone())

	return p
}

func (p *AppsServer) ConfigureFlags() {
	p.Server.ConfigureFlags()

	p.Server.Host = p.Cfg.Host
	p.Server.Port = p.Cfg.Port
}

func (p *AppsServer) ConfigureAPI() {
	RegisterHandler(p.Api, p.service)
}

func (p *AppsServer) Serve() error {
	p.ConfigureFlags()
	p.ConfigureAPI()

	db, err := OpenAppDatabase(p.Cfg)
	if err != nil {
		return err
	}

	p.service.InitAppDatabase(db)
	return p.Server.Serve()
}

func ListenAndServeAppsServer(config *config.Config) error {
	server := NewAppsServer(config)
	defer server.Shutdown()

	if err := server.Serve(); err != nil {
		return err
	}
	return nil
}
