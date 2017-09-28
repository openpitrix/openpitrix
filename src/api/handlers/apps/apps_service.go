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

var _ AppsServiceHander = (*AppsService)(nil)

type Options struct {
	// "gopkg.in/gorp.v2"
	// *gorp.DbMap
}

type AppsService struct {
	//
}

func DefaultAppsService(opt *Options) *AppsService {
	return &AppsService{}
}

func (p *AppsService) AppsGetAppsHandler(apps.GetAppsParams) middleware.Responder {
	return middleware.NotImplemented("TODO")
}

func (p *AppsService) AppsPostAppsHandler(apps.PostAppsParams) middleware.Responder {
	return middleware.NotImplemented("TODO")
}

func (p *AppsService) AppsGetAppsAppIDHandler(apps.GetAppsAppIDParams) middleware.Responder {
	return middleware.NotImplemented("TODO")
}

func (p *AppsService) AppsDeleteAppsAppIDHandler(apps.DeleteAppsAppIDParams) middleware.Responder {
	return middleware.NotImplemented("TODO")
}
