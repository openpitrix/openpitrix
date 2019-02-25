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

func (s *Server) CreateCombinationSpu(ctx context.Context, req *pb.CreateCombinationSpuRequest) (*pb.CreateCombinationSpuResponse, error) {
	comSpu := models.PbToCombinationSpu(req)
	//check if Spu exist
	for _, spuId := range comSpu.SpuIds {
		err := checkStructExist(ctx, models.Spu{}, spuId)
		if err != nil {
			return nil, err
		}
	}

	// insert CombinationSpu
	err := insertCombinationSpu(ctx, comSpu)
	if err != nil {
		return nil, internalError(ctx, err)
	}
	return &pb.CreateCombinationSpuResponse{CombinationSpuId: pbutil.ToProtoString(comSpu.CombinationSpuId)}, nil
}

func (s *Server) CreateCombinationSku(ctx context.Context, req *pb.CreateCombinationSkuRequest) (*pb.CreateCombinationSkuResponse, error) {
	comSku := models.PbToCombinationSku(req)

	//TODO: check if CombinationSpu exist\
	//...
	//TODO: check if attribute exist
	//...

	//insert CombinationSku
	err := insertCombinationSku(ctx, comSku)
	if err != nil {
		return nil, internalError(ctx, err)
	}
	return &pb.CreateCombinationSkuResponse{CombinationSkuId: pbutil.ToProtoString(comSku.CombinationSkuId)}, nil
}

func (s *Server) CreateCombinationPrice(ctx context.Context, req *pb.CreateCombinationPriceRequest) (*pb.CreateCombinationPriceResponse, error) {
	comPrice := models.PbToCombinationPrice(req)

	//TODO:check if combination_sku exist
	//TODO:check if combination_spu exist
	//TODO: check if attribute exist

	//insert com_price
	err := insertCombinationPrice(ctx, comPrice)
	if err != nil {
		return nil, internalError(ctx, err)
	}

	return &pb.CreateCombinationPriceResponse{CombinationPriceId: pbutil.ToProtoString(comPrice.CombinationSkuId)}, nil
}

func (s *Server) CreateProbationSku(ctx context.Context, req *pb.CreateProbationSkuRequest) (*pb.CreateProbationSkuResponse, error) {
	proSku := models.PbToProbationSku(req)

	//TODO:check if spu exist
	//TODO:check if attribute exist

	//insert probation_sku
	err := insertProbationSku(ctx, proSku)
	if err != nil {
		return nil, internalError(ctx, err)
	}
	return &pb.CreateProbationSkuResponse{ProbationSkuId: pbutil.ToProtoString(proSku.ProbationSkuId)}, nil
}
