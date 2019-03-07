// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package mbing

import (
	"context"
	"fmt"

	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
	"openpitrix.io/openpitrix/pkg/util/stringutil"
)

func (s *Server) CreatePrice(ctx context.Context, req *pb.CreatePriceRequest) (*pb.CreatePriceResponse, error) {
	price := models.PbToPrice(req)

	//get sku
	sku, err := getSku(ctx, price.SkuId)
	if err != nil {
		return nil, err
	}
	if sku == nil {
		return nil, notExistError(ctx, models.Sku{}, price.SkuId)
	}

	//check if price.AttributeId exist in sku.MeteringAttributeIds
	if !stringutil.StringIn(price.AttributeId, sku.MeteringAttributeIds) {
		return nil, notExistInOtherError(ctx, models.Attribute{}, models.Sku{})
	}

	//insert price
	err = insertPrice(ctx, price)
	if err != nil {
		return nil, internalError(ctx, err)
	}
	return &pb.CreatePriceResponse{PriceId: pbutil.ToProtoString(price.PriceId)}, nil
}

func Billing() {
	leasing := models.Leasing{}
	contract, err := calculate(leasing)
	fmt.Printf("%v", err)

	_, err = Charge(contract)

	if err.Error() == "balance not enough" {
		addToNoMoney(leasing)
	}

}

func calculate(leasing models.Leasing) (*models.LeasingContract, error) {
	return &models.LeasingContract{}, nil
}

func addToNoMoney(leasing models.Leasing) error {
	return nil
}
