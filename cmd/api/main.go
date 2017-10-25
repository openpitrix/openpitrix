// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// install httpie:
// - macOS: brew install httpie
// - Ubuntu: apt-get install httpie
// - CentOS: yum install httpie

// http get :9527/v1/apps

// curl http://localhost:9527/v1/apps

// openpitrix server
package main

import (
	"log"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"

	"openpitrix.io/openpitrix/pkg/swagger/models"
	"openpitrix.io/openpitrix/pkg/swagger/restapi"
	"openpitrix.io/openpitrix/pkg/swagger/restapi/operations"
	"openpitrix.io/openpitrix/pkg/swagger/restapi/operations/apps"
)

type Server struct {
	*restapi.Server
	Spec *loads.Document
	Api  *operations.OpenPitrixAPI
}

func NewServer() *Server {
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		log.Fatal(err)
	}

	p := new(Server)
	p.Api = operations.NewOpenPitrixAPI(swaggerSpec)
	p.Spec = swaggerSpec
	p.Server = restapi.NewServer(p.Api)

	return p
}

func (p *Server) ConfigureFlags() {
	p.Server.ConfigureFlags()
	p.Port = 9527
}

func (p *Server) ConfigureAPI() {
	p.Api.AppsGetAppsHandler = apps.GetAppsHandlerFunc(
		func(params apps.GetAppsParams) middleware.Responder {
			return apps.NewGetAppsOK().WithPayload(&models.GetAppsOKBody{
				Apps: models.Apps{
					Items: models.AppsItems{
						&models.App{
							AppID: swag.String("app-id-1"),
						},
						&models.App{
							AppID: swag.String("app-id-2"),
						},
					},
				},
				Paging: models.Paging{
					TotalItems: 2,
				},
			})
		},
	)
}

func main() {
	server := NewServer()
	defer server.Shutdown()

	server.ConfigureFlags()
	server.ConfigureAPI()

	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}
}
