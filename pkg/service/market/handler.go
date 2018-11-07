// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package market

import (
	"context"
	"time"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/gerr"
	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/pi"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
	"openpitrix.io/openpitrix/pkg/util/senderutil"
)

func (p *Server) CreateMarket(ctx context.Context, req *pb.CreateMarketRequest) (*pb.CreateMarketResponse, error) {
	s := senderutil.GetSenderFromContext(ctx)

	newMarket := models.NewMarket(
		req.GetName().GetValue(),
		req.GetVisibility().GetValue(),
		constants.StatusEnabled,
		req.GetDescription().GetValue(),
		s.UserId,
	)

	_, err := pi.Global().DB(ctx).
		InsertInto(constants.TableMarket).
		Record(newMarket).
		Exec()

	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorCreateResourceFailed)
	}

	newMarketId := newMarket.MarketId
	marketUser := models.NewMarketUser(newMarketId, s.UserId, s.UserId)

	_, err = pi.Global().DB(ctx).
		InsertInto(constants.TableMarketUser).
		Record(marketUser).
		Exec()

	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}

	res := &pb.CreateMarketResponse{
		MarketId: pbutil.ToProtoString(newMarketId),
	}

	return res, nil

}

func (p *Server) DescribeMarkets(ctx context.Context, req *pb.DescribeMarketsRequest) (*pb.DescribeMarketsResponse, error) {
	var markets []*models.Market
	offset := pbutil.GetOffsetFromRequest(req)
	limit := pbutil.GetLimitFromRequest(req)

	query := pi.Global().DB(ctx).
		Select(models.MarketColumns...).
		From(constants.TableMarket).
		Offset(offset).
		Limit(limit).
		Where(manager.BuildFilterConditions(req, constants.TableMarket))
	query = manager.AddQueryOrderDir(query, req, constants.ColumnCreateTime)

	_, err := query.Load(&markets)

	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}

	count, err := query.Count()
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}

	res := &pb.DescribeMarketsResponse{
		MarketSet:  models.MarketToPbs(markets),
		TotalCount: uint32(count),
	}

	return res, nil
}

func (p *Server) ModifyMarket(ctx context.Context, req *pb.ModifyMarketRequest) (*pb.ModifyMarketResponse, error) {
	marketId := req.GetMarketId().GetValue()

	market, err := CheckMarketPermisson(ctx, marketId)
	if err != nil {
		return nil, err
	}

	if market.Status == constants.StatusDeleted {
		return nil, gerr.NewWithDetail(ctx, gerr.FailedPrecondition, err, gerr.ErrorResourceAlreadyDeleted, marketId)
	}

	attributes := manager.BuildUpdateAttributes(req,
		constants.ColumnName, constants.ColumnVisibility, constants.ColumnStatus, constants.ColumnDescription)
	if _, ok := attributes[constants.ColumnStatus]; ok {
		attributes[constants.ColumnStatusTime] = time.Now()
	}
	_, err = pi.Global().DB(ctx).
		Update(constants.TableMarket).
		SetMap(attributes).
		Where(db.Eq(constants.ColumnMarketId, marketId)).
		Exec()

	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorModifyResourceFailed)
	}

	res := &pb.ModifyMarketResponse{
		MarketId: req.GetMarketId(),
	}

	return res, nil
}

func (p *Server) DeleteMarkets(ctx context.Context, req *pb.DeleteMarketsRequest) (*pb.DeleteMarketsResponse, error) {
	marketIds := req.GetMarketId()

	markets, err := CheckMarketsPermission(ctx, marketIds)
	if err != nil {
		return nil, err
	}

	for _, market := range markets {
		if market.Status == constants.StatusDeleted {
			return nil, gerr.NewWithDetail(ctx, gerr.FailedPrecondition, err, gerr.ErrorResourceAlreadyDeleted, market.MarketId)
		}
	}

	_, err = pi.Global().DB(ctx).Update(constants.TableMarket).
		Set(constants.ColumnStatus, constants.StatusDeleted).
		Set(constants.ColumnStatusTime, time.Now()).
		Where(db.Eq(constants.ColumnMarketId, marketIds)).
		Exec()

	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDeleteResourceFailed)
	}

	res := &pb.DeleteMarketsResponse{
		MarketId: marketIds,
	}

	return res, nil
}

func (p *Server) UserJoinMarket(ctx context.Context, req *pb.UserJoinMarketRequest) (*pb.UserJoinMarketResponse, error) {
	userIds := req.GetUserId()
	marketIds := req.GetMarketId()

	markets, err := CheckMarketsPermission(ctx, marketIds)
	if err != nil {
		return nil, err
	}

	insert := pi.Global().DB(ctx).InsertInto(constants.TableMarketUser)
	var ok bool
	for _, market := range markets {
		marketId := market.MarketId
		owner := market.Owner
		if market.Status == constants.StatusDeleted {
			return nil, gerr.NewWithDetail(ctx, gerr.FailedPrecondition, err, gerr.ErrorResourceAlreadyDeleted, marketId)
		}

		for _, userId := range userIds {
			count, err := pi.Global().DB(ctx).
				Select(models.MarketUserColumns...).
				From(constants.TableMarketUser).
				Where(db.Eq(constants.ColumnMarketId, marketId)).
				Where(db.Eq(constants.ColumnUserId, userId)).
				Count()
			if err != nil {
				return nil, err
			}

			if count == 0 {
				ok = true
				record := models.NewMarketUser(marketId, userId, owner)
				insert = insert.Record(record)
			}
		}
	}

	if ok {
		_, err = insert.Exec()

		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
		}
	}

	res := &pb.UserJoinMarketResponse{
		MarketId: marketIds,
		UserId:   userIds,
	}

	return res, nil
}

func (p *Server) UserLeaveMarket(ctx context.Context, req *pb.UserLeaveMarketRequest) (*pb.UserLeaveMarketResponse, error) {

	marketIds := req.GetMarketId()
	userIds := req.GetUserId()

	//Check if the user can access those markets or users
	markets, err := CheckMarketsPermission(ctx, marketIds)
	if err != nil {
		return nil, err
	}

	for _, market := range markets {
		if market.Status == constants.StatusDeleted {
			return nil, gerr.NewWithDetail(ctx, gerr.FailedPrecondition, err, gerr.ErrorResourceAlreadyDeleted, market.MarketId)
		}
	}

	_, err = pi.Global().DB(ctx).
		DeleteFrom(constants.TableMarketUser).
		Where(db.Eq(constants.ColumnMarketId, marketIds)).
		Where(db.Eq(constants.ColumnUserId, userIds)).
		Exec()

	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}

	res := &pb.UserLeaveMarketResponse{
		MarketId: req.GetMarketId(),
		UserId:   req.GetUserId(),
	}

	return res, nil
}

func (p *Server) DescribeMarketUsers(ctx context.Context, req *pb.DescribeMarketUsersRequest) (*pb.DescribeMarketUsersResponse, error) {
	var marketUsers []*models.MarketUser
	s := senderutil.GetSenderFromContext(ctx)
	offset := pbutil.GetOffsetFromRequest(req)
	limit := pbutil.GetLimitFromRequest(req)

	if len(req.GetOwner()) > 0 {
		req.Owner = []string{s.UserId}
	} else {
		req.UserId = []string{s.UserId}
	}

	query := pi.Global().DB(ctx).Select(models.MarketUserColumns...).
		From(constants.TableMarketUser).
		Offset(offset).Limit(limit).
		Where(manager.BuildFilterConditions(req, constants.TableMarketUser))

	query = manager.AddQueryOrderDir(query, req, constants.ColumnCreateTime)

	_, err := query.Load(&marketUsers)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}

	count, err := query.Count()
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}

	res := &pb.DescribeMarketUsersResponse{
		MarketUserSet: models.MarketUserToPbs(marketUsers),
		TotalCount:    count,
	}

	return res, nil
}
