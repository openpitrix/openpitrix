// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package mbing

import (
	"context"

	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
)

func (s *Server) CreateCombinationPrice(ctx context.Context, req *pb.CreateCombinationPriceRequest) (*pb.CreateCombinationPriceResponse, error) {
	comPrice := models.PbToCombinationPrice(req)

	//TODO: How to check combination_binding_id
	//Same as Price, how about do not check ?

	//insert com_price
	err := insertCombinationPrice(ctx, comPrice)
	if err != nil {
		return nil, internalError(ctx, err)
	}

	return &pb.CreateCombinationPriceResponse{CombinationPriceId: pbutil.ToProtoString(comPrice.CombinationSkuId)}, nil
}

func (s *Server) DescribeCombinationPrices(ctx context.Context, req *pb.DescribeCombinationPricesRequest) (*pb.DescribeCombinationPricesResponse, error) {
	//TODO: impl DescribeCombinationPrices
	return &pb.DescribeCombinationPricesResponse{}, nil
}

func (s *Server) ModifyCombinationPrice(ctx context.Context, req *pb.ModifyCombinationPriceRequest) (*pb.ModifyCombinationPriceResponse, error) {
	//TODO: impl ModifyCombinationPrice
	return &pb.ModifyCombinationPriceResponse{}, nil
}

func (s *Server) DeleteCombinationPrices(ctx context.Context, req *pb.DeleteCombinationPricesRequest) (*pb.DeleteCombinationPricesResponse, error) {
	//TODO: impl DeleteCombinationPrices
	return &pb.DeleteCombinationPricesResponse{}, nil
}

func (s *Server) CreateProbation(ctx context.Context, req *pb.CreateProbationRequest) (*pb.CreateProbationResponse, error) {
	proSku := models.PbToProbation(req)

	//TODO:check if spu exist
	//TODO:check if attribute exist

	//insert probation_sku
	err := insertProbation(ctx, proSku)
	if err != nil {
		return nil, internalError(ctx, err)
	}
	return &pb.CreateProbationResponse{ProbationId: pbutil.ToProtoString(proSku.ProbationId)}, nil
}

func (s *Server) DescribeProbations(ctx context.Context, req *pb.DescribeProbationsRequest) (*pb.DescribeProbationsResponse, error) {
	//TODO: impl DescribeProbations
	return &pb.DescribeProbationsResponse{}, nil
}

func (s *Server) ModifyProbation(ctx context.Context, req *pb.ModifyProbationRequest) (*pb.ModifyProbationResponse, error) {
	//TODO: impl ModifyProbation
	return &pb.ModifyProbationResponse{}, nil
}

func (s *Server) DeleteProbations(ctx context.Context, req *pb.DeleteProbationsRequest) (*pb.DeleteProbationsResponse, error) {
	//TODO: impl DeleteProbations
	return &pb.DeleteProbationsResponse{}, nil
}

func (s *Server) CreateDiscount(ctx context.Context, req *pb.CreateDiscountRequest) (*pb.CreateDiscountResponse, error) {
	//TODO: impl CreateDiscount
	return &pb.CreateDiscountResponse{}, nil
}

func (s *Server) DescribeDiscounts(ctx context.Context, req *pb.DescribeDiscountsRequest) (*pb.DescribeDiscountsResponse, error) {
	//TODO: impl DescribeDiscounts
	return &pb.DescribeDiscountsResponse{}, nil
}

func (s *Server) ModifyDiscount(ctx context.Context, req *pb.ModifyDiscountRequest) (*pb.ModifyDiscountResponse, error) {
	//TODO: impl ModifyDiscount
	return &pb.ModifyDiscountResponse{}, nil
}

func (s *Server) DeleteDiscounts(ctx context.Context, req *pb.DeleteDiscountsRequest) (*pb.DeleteDiscountsResponse, error) {
	//TODO: impl DeleteDiscounts
	return &pb.DeleteDiscountsResponse{}, nil
}

func (s *Server) CreateCoupon(ctx context.Context, req *pb.CreateCouponRequest) (*pb.CreateCouponResponse, error) {
	//TODO: impl CreateCoupon
	return &pb.CreateCouponResponse{}, nil
}

func (s *Server) DescribeCoupons(ctx context.Context, req *pb.DescribeCouponsRequest) (*pb.DescribeCouponsResponse, error) {
	//TODO: impl DescribeCoupons
	return &pb.DescribeCouponsResponse{}, nil
}

func (s *Server) ModifyCoupon(ctx context.Context, req *pb.ModifyCouponRequest) (*pb.ModifyCouponResponse, error) {
	//TODO: impl ModifyCoupon
	return &pb.ModifyCouponResponse{}, nil
}

func (s *Server) DeleteCoupons(ctx context.Context, req *pb.DeleteCouponsRequest) (*pb.DeleteCouponsResponse, error) {
	//TODO: impl DeleteCoupons
	return &pb.DeleteCouponsResponse{}, nil
}

func (s *Server) CreateCouponReceived(ctx context.Context, req *pb.CreateCouponReceivedRequest) (*pb.CreateCouponReceivedResponse, error) {
	//TODO: impl CreateCouponReceived
	return &pb.CreateCouponReceivedResponse{}, nil
}

func (s *Server) DescribeCouponReceiveds(ctx context.Context, req *pb.DescribeCouponReceivedsRequest) (*pb.DescribeCouponReceivedsResponse, error) {
	//TODO: impl DescribeCouponReceiveds
	return &pb.DescribeCouponReceivedsResponse{}, nil
}

func (s *Server) DeleteCouponReceiveds(ctx context.Context, req *pb.DeleteCouponReceivedsRequest) (*pb.DeleteCouponReceivedsResponse, error) {
	//TODO: impl DeleteCouponReceiveds
	return &pb.DeleteCouponReceivedsResponse{}, nil
}
