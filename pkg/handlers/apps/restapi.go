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

package apps

import (
	"fmt"
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	"openpitrix.io/openpitrix/pkg/swagger/models"
	"openpitrix.io/openpitrix/pkg/swagger/restapi/operations"
	"openpitrix.io/openpitrix/pkg/swagger/restapi/operations/apps"
)

var _ AppsRestInterface = (*AppsRestService)(nil)

type AppsRestInterface interface {
	GetApps(apps.GetAppsParams) middleware.Responder
	PostApps(apps.PostAppsParams) middleware.Responder
	GetAppsAppID(apps.GetAppsAppIDParams) middleware.Responder
	DeleteAppsAppID(apps.DeleteAppsAppIDParams) middleware.Responder
}

func RegisterHandler(api *operations.OpenPitrixAPI, handler AppsRestInterface) {
	api.AppsGetAppsHandler = apps.GetAppsHandlerFunc(handler.GetApps)
	api.AppsPostAppsHandler = apps.PostAppsHandlerFunc(handler.PostApps)

	api.AppsGetAppsAppIDHandler = apps.GetAppsAppIDHandlerFunc(handler.GetAppsAppID)
	api.AppsDeleteAppsAppIDHandler = apps.DeleteAppsAppIDHandlerFunc(handler.DeleteAppsAppID)
}

type Options struct{}

type AppsRestService struct {
	opt *Options
	db  AppDatabaseInterface
}

func NewAppsRestService(opt *Options) *AppsRestService {
	return &AppsRestService{opt: opt}
}

func (p *AppsRestService) InitAppDatabase(db AppDatabaseInterface) {
	p.db = db
}

func (p *AppsRestService) GetApps(params apps.GetAppsParams) middleware.Responder {
	items, err := p.db.GetApps()
	if err != nil {
		return apps.NewGetAppsInternalServerError().WithPayload(&models.Error{
			Code:    fmt.Sprintf("%d", http.StatusInternalServerError),
			Message: err.Error(),
		})
	}

	modelsAppsItems := items.To_models_AppsItems(int(*params.PageNumber), int(*params.PageSize))
	return apps.NewGetAppsOK().WithPayload(&models.GetAppsOKBody{
		Apps: models.Apps{
			Items: modelsAppsItems,
		},
		Paging: models.Paging{
			TotalItems: int64(len(modelsAppsItems)),
		},
	})
}

func (p *AppsRestService) PostApps(params apps.PostAppsParams) middleware.Responder {
	err := p.db.CreateApp(new(AppsItem).From_models_App(params.App))
	if err != nil {
		return apps.NewGetAppsInternalServerError().WithPayload(&models.Error{
			Code:    fmt.Sprintf("%d", http.StatusInternalServerError),
			Message: err.Error(),
		})
	}
	return apps.NewGetAppsOK()
}

func (p *AppsRestService) GetAppsAppID(params apps.GetAppsAppIDParams) middleware.Responder {
	item, err := p.db.GetApp(params.AppID)
	if err != nil {
		return apps.NewGetAppsInternalServerError().WithPayload(&models.Error{
			Code:    fmt.Sprintf("%d", http.StatusInternalServerError),
			Message: err.Error(),
		})
	}

	return apps.NewGetAppsAppIDOK().WithPayload(item.To_models_App())
}

func (p *AppsRestService) DeleteAppsAppID(params apps.DeleteAppsAppIDParams) middleware.Responder {
	err := p.db.DeleteApp(params.AppID)
	if err != nil {
		return apps.NewGetAppsInternalServerError().WithPayload(&models.Error{
			Code:    fmt.Sprintf("%d", http.StatusInternalServerError),
			Message: err.Error(),
		})
	}
	return apps.NewGetAppsOK()
}
