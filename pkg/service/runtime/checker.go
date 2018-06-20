// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package runtime

import (
	"context"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/pb"
)

func (p *Server) Checker(ctx context.Context, req interface{}) error {
	switch r := req.(type) {
	case *pb.CreateRuntimeRequest:
		return manager.NewChecker(ctx, r).
			Required("name", "provider", "runtime_url", "zone", "runtime_credential").
			StringChosen("provider", constants.SupportedProvider).
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
			StringChosen("provider", constants.SupportedProvider).
			Exec()

	}
	return nil
}
