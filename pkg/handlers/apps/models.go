// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package apps

import (
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"openpitrix.io/openpitrix/pkg/swagger/models"
)

type AppsItem struct {
	Id          string    `db:"id, size:50, primarykey"`
	Name        string    `db:"name", size:50`
	Description string    `db:"description, size:1000"`
	Url         string    `db:"url, size:255"`
	Created     time.Time `db:"created"`
}

type AppsItems []AppsItem

func (p *AppsItem) From_models_App(app *models.App) *AppsItem {
	*p = AppsItem{
		Id:          swag.StringValue(app.AppID),
		Name:        app.Name,
		Description: app.Description,
		Url:         app.URL,
		Created:     time.Time(app.CreateTime),
	}
	return p
}

func (p *AppsItem) To_models_App() *models.App {
	return &models.App{
		AppID:       swag.String(p.Id),
		CreateTime:  strfmt.DateTime(p.Created),
		Description: p.Description,
		Name:        p.Name,
		URL:         p.Url,
	}
}

func (p *AppsItems) From_models_AppsItems(items models.AppsItems) AppsItems {
	q := make(AppsItems, len(items))
	for i, v := range items {
		q[i].From_models_App(v)
	}
	*p = q
	return *p
}

func (p AppsItems) To_models_AppsItems(pageNumber, pageSize int) models.AppsItems {
	if pageNumber > 0 {
		pageNumber = pageNumber - 1 // start with 1
	}

	start := pageNumber * pageSize
	end := start + pageSize

	if start >= len(p) {
		return nil
	}
	if end > len(p) {
		end = len(p)
	}

	q := make(models.AppsItems, end-start)
	for i := start; i < end; i++ {
		q[i-start] = p[i].To_models_App()
	}
	return q
}
