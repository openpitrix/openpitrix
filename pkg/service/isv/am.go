// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package isv

import (
	"context"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/pb"
)

var SupportedStatus = []string{
	constants.StatusNew,
	constants.StatusSubmitted,
	constants.StatusPassed,
	constants.StatusRejected,
}

func (s *Server) Checker(ctx context.Context, req interface{}) error {
	switch r := req.(type) {
	case *pb.SubmitVendorVerifyInfoRequest:
		return manager.NewChecker(ctx, r).
			Required(constants.ColumnUserId, constants.ColumnCompanyName, constants.ColumnCompanyWebsite,
				constants.ColumnCompanyProfile, constants.ColumnAuthorizerName, constants.ColumnAuthorizerEmail,
				constants.ColumnAuthorizerPhone, constants.ColumnBankName, constants.ColumnBankAccountName, constants.ColumnBankAccountNumber).
			Exec()
	case *pb.DescribeVendorVerifyInfosRequest:
		return manager.NewChecker(ctx, r).
			StringChosen("status", SupportedStatus).
			Exec()
	case *pb.PassVendorVerifyInfoRequest:
		return manager.NewChecker(ctx, r).
			Required(constants.ColumnUserId).
			Exec()
	case *pb.RejectVendorVerifyInfoRequest:
		return manager.NewChecker(ctx, r).
			Required(constants.ColumnUserId, constants.ColumnRejectMessage).
			Exec()
	}
	return nil
}
