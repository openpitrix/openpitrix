// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package manager

import (
	"context"
	"strings"

	"github.com/gocraft/dbr"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/gerr"
	"openpitrix.io/openpitrix/pkg/util/ctxutil"
)

func CheckOwnerPathPermission(ctx context.Context, req interface{}) error {
	s := ctxutil.GetSender(ctx)
	param := "owner_path"
	if r, ok := req.(RequestWithOwnerPath); ok {
		if !s.AccessPath.CheckOwnerPathPermission(r.GetOwnerPath()...) {
			return gerr.New(ctx, gerr.InvalidArgument, gerr.ErrorUnsupportedParameterValue, param, strings.Join(r.GetOwnerPath(), ","))
		}
	}
	return nil
}

func BuildOwnerPathFilter(ctx context.Context, req Request) dbr.Builder {
	s := ctxutil.GetSender(ctx)
	if s == nil {
		return nil
	}
	accessPath := string(s.GetAccessPath())

	var ownerPaths []string
	if r, ok := req.(RequestWithOwnerPath); ok {
		ownerPaths = r.GetOwnerPath()
	}

	if len(ownerPaths) == 0 {
		return db.Prefix(constants.ColumnOwnerPath, accessPath)
	} else {
		var ops []dbr.Builder
		for _, ownerPath := range ownerPaths {
			ops = append(ops, db.Prefix(constants.ColumnOwnerPath, ownerPath))
		}
		return db.Or(ops...)
	}
}
