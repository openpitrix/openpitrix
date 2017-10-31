// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package apps

import (
	"fmt"
	"net/http"

	"github.com/go-openapi/runtime/middleware"

	"openpitrix.io/openpitrix/pkg/config-v2"
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

type AppsRestService struct {
	Cfg *config.Config
	db  AppDatabaseInterface
}

func NewAppsRestService(config *config.Config, db ...AppDatabaseInterface) *AppsRestService {
	if len(db) > 0 {
		return &AppsRestService{Cfg: config.Clone(), db: db[0]}
	} else {
		return &AppsRestService{Cfg: config.Clone()}
	}
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
	return apps.NewPostAppsNoContent()
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
	if _, err := p.db.GetApp(params.AppID); err != nil {
		return apps.NewDeleteAppsAppIDNotFound()
	}

	err := p.db.DeleteApp(params.AppID)
	if err != nil {
		return apps.NewGetAppsInternalServerError().WithPayload(&models.Error{
			Code:    fmt.Sprintf("%d", http.StatusInternalServerError),
			Message: err.Error(),
		})
	}
	return apps.NewDeleteAppsAppIDNoContent()
}
