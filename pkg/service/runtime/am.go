// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package runtime

import (
	"context"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/pi"
	"openpitrix.io/openpitrix/pkg/plugins"
	"openpitrix.io/openpitrix/pkg/util/senderutil"
)

func (p *Server) Checker(ctx context.Context, req interface{}) error {
	switch r := req.(type) {
	case *pb.CreateRuntimeRequest:
		return manager.NewChecker(ctx, r).
			Required("name", "provider", "zone", "runtime_credential").
			StringChosen("provider", plugins.GetAvailablePlugins(pi.Global().GlobalConfig().Cluster.Plugins)).
			Exec()
	case *pb.ModifyRuntimeRequest:
		return manager.NewChecker(ctx, r).
			Required("runtime_id").
			Exec()
	case *pb.DeleteRuntimesRequest:
		return manager.NewChecker(ctx, r).
			Required("runtime_id").
			Exec()
	case *pb.DescribeRuntimeProviderZonesRequest:
		return manager.NewChecker(ctx, r).
			Required("provider", "runtime_url", "runtime_credential").
			StringChosen("provider", plugins.GetAvailablePlugins(pi.Global().GlobalConfig().Cluster.Plugins)).
			Exec()
	case *pb.GetRuntimeStatisticsRequest:
		return manager.NewChecker(ctx, r).
			Role(constants.AllAdminRoles).
			Exec()
	}
	return nil
}

func (p *Server) Builder(ctx context.Context, req interface{}) interface{} {
	sender := senderutil.GetSenderFromContext(ctx)
	switch r := req.(type) {
	case *pb.DescribeRuntimesRequest:
		if sender.IsGlobalAdmin() {

		} else if sender.IsUser() {
			r.Owner = []string{sender.UserId}
		}
		return r
	}
	return req
}
