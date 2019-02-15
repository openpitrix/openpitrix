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

func (s *Server) CreateCRA(ctx context.Context, req *pb.CreateCRARequest) (*pb.CreateCRAResponse, error) {
	cra := models.PbToCRA(req)
	//check if resourceAttributes exist
	for _, resAttId := range cra.ResourceAttributeIds {
		err := checkStructExistById(ctx, models.ResourceAttribute{}, cra, resAttId, CreateFailedCode)
		if err != nil {
			return nil, err
		}
	}

	// insert CRA
	err := insertCRA(ctx, cra)
	if err != nil {
		return nil, commonInternalErr(ctx, cra, CreateFailedCode)
	}
	return &pb.CreateCRAResponse{CraId: pbutil.ToProtoString(cra.CRAId)}, nil
}

func (s *Server) CreateCombinationSku(ctx context.Context, req *pb.CreateComSkuRequest) (*pb.CreateComSkuResponse, error) {
	comSku := models.PbToComSku(req)

	//check if CRA exist
	err := checkStructExistById(ctx, models.CombinationResourceAttribute{}, comSku, comSku.ComSkuId, CreateFailedCode)
	if err != nil {
		return nil, err
	}

	//TODO: check if attribute value exist
	//...

	//insert ComSku
	err = insertComSku(ctx, comSku)
	if err != nil {
		return nil, commonInternalErr(ctx, comSku, CreateFailedCode)
	}
	return &pb.CreateComSkuResponse{ComSkuId: pbutil.ToProtoString(comSku.ComSkuId)}, nil
}

func (s *Server) CreateCombinationPrice(ctx context.Context, req *pb.CreateComPriceRequest) (*pb.CreateComPriceResponse, error) {
	comPrice := models.PbToComPrice(req)

	//check if com_sku exist
	err := checkStructExistById(ctx, models.CombinationSku{}, comPrice, comPrice.ComSkuId, CreateFailedCode)
	if err != nil {
		return nil, err
	}

	//check if attribute_value exist
	attValTmp := models.AttributeValue{}
	for k, _ := range comPrice.Prices {
		err = checkStructExistById(ctx, attValTmp, comPrice, k, CreateFailedCode)
		if err != nil {
			return nil, err
		}
	}

	//insert com_price
	err = insertComPrice(ctx, comPrice)
	if err != nil {
		return nil, commonInternalErr(ctx, comPrice, CreateFailedCode)
	}

	return &pb.CreateComPriceResponse{ComPriceId: pbutil.ToProtoString(comPrice.ComPriceId)}, nil
}

func (s *Server) CreateProbationSku(ctx context.Context, req *pb.CreateProSkuRequest) (*pb.CreateProSkuResponse, error) {
	proSku := models.PbToProSku(req)

	//check if resource_attribute exist
	err := checkStructExistById(ctx, models.ResourceAttribute{}, proSku, proSku.ResourceAttributeId, CreateFailedCode)
	if err != nil {
		return nil, err
	}

	//check if attribute_values exist
	attValTmp := models.AttributeValue{}
	for _, val := range proSku.AttributeValues {
		err := checkStructExistById(ctx, attValTmp, proSku, val, CreateFailedCode)
		if err != nil {
			return nil, err
		}
	}

	//insert probation_sku
	err = insertProSku(ctx, proSku)
	if err != nil {
		return nil, commonInternalErr(ctx, proSku, CreateFailedCode)
	}
	return &pb.CreateProSkuResponse{ProSkuId: pbutil.ToProtoString(proSku.ProSkuId)}, nil
}
