// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package metering

import (
	"context"

	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
)

func (s *Server) CreateCombination(ctx context.Context, req *pb.CreateCombinationRequest) (*pb.CreateCombinationResponse, error) {
	//TODO: get id of current user
	comSpu := models.PbToCombination(req, "")

	// insert Combination
	err := insertCombination(ctx, comSpu)
	if err != nil {
		return nil, internalError(ctx, err)
	}
	return &pb.CreateCombinationResponse{CombinationId: pbutil.ToProtoString(comSpu.CombinationId)}, nil
}

func (s *Server) DescribeCombinations(ctx context.Context, req *pb.DescribeCombinationsRequest) (*pb.DescribeCombinationsResponse, error) {
	//TODO: impl DescribeCombinations
	return &pb.DescribeCombinationsResponse{}, nil
}

func (s *Server) ModifyCombination(ctx context.Context, req *pb.ModifyCombinationRequest) (*pb.ModifyCombinationResponse, error) {
	//TODO: impl ModifyCombination
	return &pb.ModifyCombinationResponse{}, nil
}

func (s *Server) DeleteCombinations(ctx context.Context, req *pb.DeleteCombinationsRequest) (*pb.DeleteCombinationsResponse, error) {
	//TODO: impl DeleteCombinations
	return &pb.DeleteCombinationsResponse{}, nil
}

func (s *Server) CreateCombinationSku(ctx context.Context, req *pb.CreateCombinationSkuRequest) (*pb.CreateCombinationSkuResponse, error) {
	comSku := models.PbToCombinationSku(req)

	//TODO: Re-impl CreateCombinationSku

	//insert CombinationSku
	err := insertCombinationSku(ctx, comSku)
	if err != nil {
		return nil, internalError(ctx, err)
	}
	return &pb.CreateCombinationSkuResponse{CombinationSkuId: pbutil.ToProtoString(comSku.CombinationSkuId)}, nil
}

func (s *Server) DescribeCombinationSkus(ctx context.Context, req *pb.DescribeCombinationSkusRequest) (*pb.DescribeCombinationSkusResponse, error) {
	//TODO: impl DescribeCombinationSkus
	return &pb.DescribeCombinationSkusResponse{}, nil
}

func (s *Server) ModifyCombinationSku(ctx context.Context, req *pb.ModifyCombinationSkuRequest) (*pb.ModifyCombinationSkuResponse, error) {
	//TODO: impl ModifyCombinationSku
	return &pb.ModifyCombinationSkuResponse{}, nil
}

func (s *Server) DeleteCombinationSkus(ctx context.Context, req *pb.DeleteCombinationSkusRequest) (*pb.DeleteCombinationSkusResponse, error) {
	//TODO: impl DeleteCombinationSkus
	return &pb.DeleteCombinationSkusResponse{}, nil
}

func (s *Server) CreateCombinationMABinding(ctx context.Context, req *pb.CreateCombinationMABindingRequest) (*pb.CreateCombinationMABindingResponse, error) {
	//TODO: impl CreateCombinationMABinding
	return &pb.CreateCombinationMABindingResponse{}, nil
}

func (s *Server) DescribeCombinationMABindings(ctx context.Context, req *pb.DescribeCombinationMABindingsRequest) (*pb.DescribeCombinationMABindingsResponse, error) {
	//TODO: impl DescribeCombinationMABindings
	return &pb.DescribeCombinationMABindingsResponse{}, nil
}

func (s *Server) ModifyCombinationMABinding(ctx context.Context, req *pb.ModifyCombinationMABindingRequest) (*pb.ModifyCombinationMABindingResponse, error) {
	//TODO: impl ModifyCombinationMABinding
	return &pb.ModifyCombinationMABindingResponse{}, nil
}

func (s *Server) DeleteCombinationMABindings(ctx context.Context, req *pb.DeleteCombinationMABindingsRequest) (*pb.DeleteCombinationMABindingsResponse, error) {
	//TODO: impl DeleteCombinationMABindings
	return &pb.DeleteCombinationMABindingsResponse{}, nil
}
