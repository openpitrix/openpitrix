// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package manager

import (
	"context"
	"fmt"

	"github.com/gocraft/dbr"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/util/ctxutil"
)

func BuildOwnerPathFilter(ctx context.Context, prefix ...string) dbr.Builder {
	s := ctxutil.GetSender(ctx)
	if s == nil {
		return nil
	}
	var column = constants.ColumnOwnerPath
	if len(prefix) > 0 {
		column = fmt.Sprintf("%s.%s", prefix[0], column)
	}
	return db.Prefix(column, string(s.GetAccessPath()))
}
