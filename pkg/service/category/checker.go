// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package category

import (
	"context"

	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/pb"
)

func (p *Server) Checker(ctx context.Context, req interface{}) error {
	switch r := req.(type) {
	case *pb.CreateCategoryRequest:
		return manager.NewChecker(ctx, r).
			Required("type", "locale").
			Exec()
	case *pb.ModifyCategoryRequest:
		return manager.NewChecker(ctx, r).
			Required("category_id").
			Exec()
	case *pb.DeleteCategoriesRequest:
		return manager.NewChecker(ctx, r).
			Required("category_id").
			Exec()
	}
	return nil
}
