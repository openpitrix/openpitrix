// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package market

import (
	"context"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/gerr"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pi"
)

func CheckMarketsPermission(ctx context.Context, resourceIds []string) ([]*models.Market, error) {
	if len(resourceIds) == 0 {
		return nil, nil
	}

	//var sender = ctxutil.GetSender(ctx)
	var markets []*models.Market
	_, err := pi.Global().DB(ctx).Select(models.MarketColumns...).From(constants.TableMarket).
		Where(db.Eq(constants.ColumnMarketId, resourceIds)).Load(&markets)

	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}

	// TODO: check permission
	//if sender != nil && !sender.IsGlobalAdmin() {
	//	for _, market := range markets {
	//		if market.Owner != sender.UserId {
	//			return nil, gerr.New(ctx, gerr.PermissionDenied, gerr.ErrorResourceAccessDenied, market.MarketId)
	//		}
	//	}
	//}

	if len(markets) == 0 {
		return nil, gerr.New(ctx, gerr.NotFound, gerr.ErrorResourceNotFound, resourceIds)
	}

	return markets, nil

}

func CheckMarketPermisson(ctx context.Context, resourceId string) (*models.Market, error) {
	if len(resourceId) == 0 {
		return nil, nil
	}

	//var sender = ctxutil.GetSender(ctx)
	var markets []*models.Market
	_, err := pi.Global().DB(ctx).
		Select(models.MarketColumns...).
		From(constants.TableMarket).
		Where(db.Eq(constants.ColumnMarketId, resourceId)).
		Load(&markets)

	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}

	// TODO: check permission
	//if sender != nil && !sender.IsGlobalAdmin() {
	//	for _, market := range markets {
	//		if market.Owner != sender.UserId {
	//			return nil, gerr.New(ctx, gerr.PermissionDenied, gerr.ErrorResourceAccessDenied, market.MarketId)
	//		}
	//	}
	//}

	if len(markets) == 0 {
		return nil, gerr.New(ctx, gerr.NotFound, gerr.ErrorResourceNotFound, resourceId)
	}

	return markets[0], nil
}
