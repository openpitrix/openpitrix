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
	"time"

	"apphub/src/api/swagger/models"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

type AppsItem struct {
	Id           string    `db:"id, size:50, primarykey"`
	Name         string    `db:"name", size:50`
	Description  string    `db:"description, size:1000"`
	RepoId       string    `db:"repo_id"`
	Url          string    `db:"url, size:255"`
	Created      time.Time `db:"created"`
	LastModified time.Time `db:"last_modified"`
}

type AppsItems []AppsItem

func (p *AppsItem) From_models_App(app *models.App) *AppsItem {
	*p = AppsItem{
		Id:          swag.StringValue(app.AppID),
		Name:        app.Name,
		Description: app.Description,
		Url:         app.URL,
		Created:     time.Now(),
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

func (p AppsItems) To_models_AppsItems(pageNumber, pageSize int) models.AppsItems {
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
