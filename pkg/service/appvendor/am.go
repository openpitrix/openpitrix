// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.
package appvendor

import (
	"context"

	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/pb"
)

func (s *Server) Checker(ctx context.Context, req interface{}) error {
	switch r := req.(type) {
	case *pb.SubmitVendorVerifyInfoRequest:
		return manager.NewChecker(ctx, r).
			Required("user_id", "company_name", "company_website", "company_profile", "authorizer_name", "authorizer_email", "authorizer_phone", "bank_name", "bank_account_name", "bank_account_number").
			Exec()
	case *pb.DescribeVendorVerifyInfosRequest:
		return manager.NewChecker(ctx, r).
			Required().
			Exec()
	case *pb.GetVendorVerifyInfoRequest:
		return manager.NewChecker(ctx, r).
			Required("user_id").
			Exec()
	case *pb.PassVendorVerifyInfoRequest:
		return manager.NewChecker(ctx, r).
			Required("user_id").
			Exec()
	case *pb.RejectVendorVerifyInfoRequest:
		return manager.NewChecker(ctx, r).
			Required("user_id").
			Exec()
	}
	return nil
}
