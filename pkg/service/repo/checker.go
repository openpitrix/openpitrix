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
			StringChosen("providers", plugins.GetAvailablePlugins(pi.Global().GlobalConfig().Cluster.Plugins)).
			StringChosen("visibility", SupportedVisibility).
			Exec()
	case *pb.ModifyRepoRequest:
		return manager.NewChecker(ctx, r).
			Required("repo_id").
			StringChosen("providers", plugins.GetAvailablePlugins(pi.Global().GlobalConfig().Cluster.Plugins)).
			StringChosen("visibility", SupportedVisibility).
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
