// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package apps

import (
	"log"
	"net/url"
	"strconv"

	"github.com/go-openapi/loads"

	"openpitrix.io/openpitrix/pkg/swagger/restapi"
	"openpitrix.io/openpitrix/pkg/swagger/restapi/operations"
)

type AppsServer struct {
	*restapi.Server
	Spec *loads.Document
	Api  *operations.OpenPitrixAPI

	addr    string
	dbpath  string
	service *AppsRestService
}

func NewAppsServer(addr, dbpath string) *AppsServer {
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		log.Fatal(err)
	}

	p := new(AppsServer)
	p.Api = operations.NewOpenPitrixAPI(swaggerSpec)
	p.Spec = swaggerSpec
	p.Server = restapi.NewServer(p.Api)

	p.addr = addr
	p.dbpath = dbpath

	p.service = NewAppsRestService(&Options{})

	return p
}

func (p *AppsServer) ConfigureFlags() {
	p.Server.ConfigureFlags()
}

func (p *AppsServer) ConfigureAPI() {
	p.Server.ConfigureAPI()

	RegisterHandler(p.Api, p.service)
}

func (p *AppsServer) Serve() error {
	url, err := url.Parse(p.addr)
	if err != nil {
		return err
	}

	if s := url.Hostname(); s != "" {
		p.Host = s
	}
	if s := url.Port(); s != "" {
		p.Port, _ = strconv.Atoi(s)
	}

	db, err := OpenAppDatabase(p.dbpath, &DbOptions{})
	if err != nil {
		return err
	}

	p.service.InitAppDatabase(db)
	return p.Server.Serve()
}

func ListenAndServeAppsServer(addr, dbpath string) error {
	server := NewAppsServer(addr, dbpath)
	defer server.Shutdown()

	server.ConfigureFlags()
	server.ConfigureAPI()

	if err := server.Serve(); err != nil {
		return err
	}
	return nil
}
