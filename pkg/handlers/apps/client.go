// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package apps

import (
	"fmt"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"golang.org/x/net/context"

	"openpitrix.io/openpitrix/pkg/swagger/client"
	"openpitrix.io/openpitrix/pkg/swagger/client/apps"
)

type AppClientInterface interface {
	GetApps() (items AppsItems, err error)
	CreateApp(app *AppsItem) error
	GetApp(id string) (item AppsItem, err error)
	DeleteApp(id string) error
}

type AppClient struct {
	*apps.Client
}

func NewAppClient(host string, port int) *AppClient {
	transport := httptransport.New(fmt.Sprintf("%s:%d", host, port), client.DefaultBasePath, nil)
	return &AppClient{Client: client.New(transport, strfmt.Default).Apps}
}

func (p *AppClient) GetApps(ctx context.Context, pageNumber, pageSize int) (items AppsItems, err error) {
	params := &apps.GetAppsParams{}

	if ctx == nil {
		params.Context = context.Background()
	} else {
		params.Context = ctx
	}

	if pageNumber > 0 {
		params.PageNumber = swag.Int64(int64(pageNumber))
	}
	if pageSize > 0 {
		params.PageSize = swag.Int32(int32(pageSize))
	}

	resp, err := p.Client.GetApps(params)
	if err != nil {
		return nil, err
	}

	items.From_models_AppsItems(resp.Payload.Items)
	return
}

func (p *AppClient) CreateApp(ctx context.Context, app AppsItem) error {
	params := &apps.PostAppsParams{
		App: app.To_models_App(),
	}

	if err := params.App.Validate(strfmt.Default); err != nil {
		return err
	}

	if ctx == nil {
		params.Context = context.Background()
	} else {
		params.Context = ctx
	}

	_, err := p.Client.PostApps(params)
	if err != nil {
		return err
	}

	return nil
}

func (p *AppClient) GetApp(ctx context.Context, id string) (item AppsItem, err error) {
	params := &apps.GetAppsAppIDParams{
		AppID: id,
	}

	if ctx == nil {
		params.Context = context.Background()
	} else {
		params.Context = ctx
	}

	resp, err := p.Client.GetAppsAppID(params)
	if err != nil {
		return AppsItem{}, err
	}

	item.From_models_App(resp.Payload)
	return
}

func (p *AppClient) DeleteApp(ctx context.Context, id string) error {
	params := &apps.DeleteAppsAppIDParams{
		AppID: id,
	}

	if ctx == nil {
		params.Context = context.Background()
	} else {
		params.Context = ctx
	}

	_, err := p.Client.DeleteAppsAppID(params)
	if err != nil {
		return err
	}

	return nil
}
