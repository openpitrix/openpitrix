// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package repo

import (
	"context"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/pi"
	"openpitrix.io/openpitrix/pkg/plugins"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
)

var SupportedVisibility = []string{
	constants.VisibilityPrivate,
	constants.VisibilityPublic,
}

func (p *Server) Checker(ctx context.Context, req interface{}) error {
	switch r := req.(type) {
	case *pb.CreateRepoRequest:
		return manager.NewChecker(ctx, r).
			Required("type", "name", "url", "credential", "visibility", "providers").
			StringChosen("providers", plugins.GetAvailablePlugins()).
			StringChosen("visibility", SupportedVisibility).
			StringChosen("app_default_status", constants.AllowedAppDefaultStatus).
			Exec()
	case *pb.ModifyRepoRequest:
		return manager.NewChecker(ctx, r).
			Required("repo_id").
			StringChosen("providers", plugins.GetAvailablePlugins()).
			StringChosen("visibility", SupportedVisibility).
			StringChosen("app_default_status", constants.AllowedAppDefaultStatus).
			Exec()
	case *pb.DeleteReposRequest:
		return manager.NewChecker(ctx, r).
			Required("repo_id").
			Exec()
	case *pb.ValidateRepoRequest:
		return manager.NewChecker(ctx, r).
			Required("type", "url", "credential").
			Exec()
	}
	return nil
}

func (p *Server) Builder(ctx context.Context, req interface{}) interface{} {
	switch r := req.(type) {
	case *pb.ModifyRepoRequest:
		if len(r.GetAppDefaultStatus().GetValue()) == 0 {
			r.AppDefaultStatus = pbutil.ToProtoString(constants.StatusDraft)
		}
		if r.AppDefaultStatus == nil {
			r.AppDefaultStatus = pbutil.ToProtoString(pi.Global().GlobalConfig().GetAppDefaultStatus())
		}
		return r

	case *pb.CreateRepoRequest:
		if len(r.GetAppDefaultStatus().GetValue()) == 0 {
			r.AppDefaultStatus = pbutil.ToProtoString(constants.StatusDraft)
		}
		if r.AppDefaultStatus == nil {
			r.AppDefaultStatus = pbutil.ToProtoString(pi.Global().GlobalConfig().GetAppDefaultStatus())
		}
		return r

	}
	return req
}
