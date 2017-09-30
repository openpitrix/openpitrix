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
	"github.com/go-openapi/runtime/middleware"

	"apphub/src/api/swagger/restapi/operations/apps"
)

var _ Handler = (*AppsRestService)(nil)

type Options struct {
	// "gopkg.in/gorp.v2"
	// *gorp.DbMap
}

type AppsRestService struct {
	//
}

func NewAppsRestService(opt *Options) *AppsRestService {
	return &AppsRestService{}
}

func (p *AppsRestService) GetApps(apps.GetAppsParams) middleware.Responder {
	return middleware.NotImplemented("TODO")
}

func (p *AppsRestService) PostApps(apps.PostAppsParams) middleware.Responder {
	return middleware.NotImplemented("TODO")
}

func (p *AppsRestService) GetAppsAppID(apps.GetAppsAppIDParams) middleware.Responder {
	return middleware.NotImplemented("TODO")
}

func (p *AppsRestService) DeleteAppsAppID(apps.DeleteAppsAppIDParams) middleware.Responder {
	return middleware.NotImplemented("TODO")
}
