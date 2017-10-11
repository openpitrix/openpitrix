// +-------------------------------------------------------------------------
// | Copyright (C) 2017 Yunify, Inc.
// +-------------------------------------------------------------------------
// | Licensed under the Apache License, Version 2.0 (the "License");
// | you may not use this work except in compliance with the License.
// | You may obtain a copy of the License in the LICENSE file, or at:
// |
// | http://www.apache.org/licenses/LICENSE-2.0
// |
// | Unless required by applicable law or agreed to in writing, software
// | distributed under the License is distributed on an "AS IS" BASIS,
// | WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// | See the License for the specific language governing permissions and
// | limitations under the License.
// +-------------------------------------------------------------------------

// install httpie:
// - macOS: brew install httpie
// - Ubuntu: apt-get install httpie
// - CentOS: yum install httpie

// http get :9527/v1/apps

// curl http://localhost:9527/v1/apps

// apphub server
package main

import (
	"log"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"

	"apphub/src/api/swagger/models"
	"apphub/src/api/swagger/restapi"
	"apphub/src/api/swagger/restapi/operations"
	"apphub/src/api/swagger/restapi/operations/apps"
)

type Server struct {
	*restapi.Server
	Spec *loads.Document
	Api  *operations.AppHubAPI
}

func NewServer() *Server {
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		log.Fatal(err)
	}

	p := new(Server)
	p.Api = operations.NewAppHubAPI(swaggerSpec)
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
