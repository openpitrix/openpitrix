// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package market

import (
	"context"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/util/senderutil"
)

var SupportedStatus = []string{
	constants.StatusDisabled,
	constants.StatusEnabled,
}

var SupportedVisibility = []string{
	constants.VisibilityPrivate,
	constants.VisibilityPublic,
}

func (s *Server) Checker(ctx context.Context, req interface{}) error {
	switch r := req.(type) {
	case *pb.CreateMarketRequest:
		return manager.NewChecker(ctx, r).
			Role(constants.AllRoles).
			Required("name", "visibility").
			StringChosen("visibility", SupportedVisibility).
			Exec()

	case *pb.ModifyMarketRequest:
		return manager.NewChecker(ctx, r).
			Role(constants.AllRoles).
			StringChosen("status", SupportedStatus).
			StringChosen("visibility", SupportedVisibility).
			Required("market_id").
			Exec()

	case *pb.DeleteMarketsRequest:
		return manager.NewChecker(ctx, r).
			Role(constants.AllRoles).
			Required("market_id").
			Exec()

	case *pb.UserJoinMarketRequest:
		return manager.NewChecker(ctx, r).
			Role(constants.AllRoles).
			Required("market_id", "user_id").
			Exec()

	case *pb.UserLeaveMarketRequest:
		return manager.NewChecker(ctx, r).
			Role(constants.AllRoles).
			Required("market_id", "user_id").
			Exec()

	case *pb.DescribeMarketUsersRequest:
		return manager.NewChecker(ctx, r).
			Role(constants.AllRoles).
			Exec()
	}
	return nil
}

func (p *Server) Builder(ctx context.Context, req interface{}) interface{} {
	sender := senderutil.GetSenderFromContext(ctx)
	switch r := req.(type) {
	case *pb.DescribeMarketsRequest:
		if sender.IsGlobalAdmin() {

		} else if sender.IsDeveloper() {
			r.Owner = []string{sender.UserId}
		}
		return r
	case *pb.DescribeMarketUsersRequest:
		if sender.IsGlobalAdmin() {

		} else if sender.IsDeveloper() {
			r.Owner = []string{sender.UserId}
		}
		return r
	}
	return req
}
