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
)

func (s *Server) CreatePrice(ctx context.Context, req *pb.CreatePriceRequest) (*pb.CreatePriceResponse, error) {
	price := models.PbToPrice(req)

	//TODO: how to check bindId
	//How about do not check bindId?

	//insert price
	err := insertPrice(ctx, price)
	if err != nil {
		return nil, internalError(ctx, err)
	}
	return &pb.CreatePriceResponse{PriceId: pbutil.ToProtoString(price.PriceId)}, nil
}

func (s *Server) DescribePrices(ctx context.Context, req *pb.DescribePricesRequest) (*pb.DescribePricesResponse, error) {
	//TODO: impl DescribePrices
	return &pb.DescribePricesResponse{}, nil
}

func (s *Server) ModifyPrice(ctx context.Context, req *pb.ModifyPriceRequest) (*pb.ModifyPriceResponse, error) {
	//TODO: impl ModifyPrice
	return &pb.ModifyPriceResponse{}, nil
}

func (s *Server) DeletePrices(ctx context.Context, req *pb.DeletePricesRequest) (*pb.DeletePricesResponse, error) {
	//TODO: impl DeletePrices
	return &pb.DeletePricesResponse{}, nil
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
