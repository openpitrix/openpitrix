// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package mbing

import (
	"context"
	"time"

	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
	"openpitrix.io/openpitrix/pkg/util/stringutil"
)

func (s *Server) CreateAttributeName(ctx context.Context, req *pb.CreateAttributeNameRequest) (*pb.CreateAttributeNameResponse, error) {
	attName := models.PbToAttributeName(req)
	err := insertAttributeName(ctx, attName)
	if err != nil {
		return nil, internalError(ctx, err)
	}
	return &pb.CreateAttributeNameResponse{AttributeNameId: pbutil.ToProtoString(attName.AttributeNameId)}, nil
}

func (s *Server) DescribeAttributeNames(ctx context.Context, req *pb.DescribeAttributeNamesRequest) (*pb.DescribeAttributeNamesResponse, error) {
	attributeNames, err := DescribeAttributeNames(ctx, req)
	if err != nil {
		return nil, internalError(ctx, err)
	}

	var pbAttributeNames []*pb.AttributeName
	for _, attName := range attributeNames {
		pbAttributeNames = append(pbAttributeNames, models.AttributeNameToPb(attName))
	}

	return &pb.DescribeAttributeNamesResponse{AttributeNames: pbAttributeNames}, nil
}

func (s *Server) CreateAttributeUnit(ctx context.Context, req *pb.CreateAttributeUnitRequest) (*pb.CreateAttributeUnitResponse, error) {
	attributeUnit := models.PbToAttributeUnit(req)
	err := insertAttributeUnit(ctx, attributeUnit)
	if err != nil {
		return nil, internalError(ctx, err)
	}
	return &pb.CreateAttributeUnitResponse{AttributeUnitId: pbutil.ToProtoString(attributeUnit.AttributeUnitId)}, nil
}

func (s *Server) DescribeAttributeUnits(ctx context.Context, req *pb.DescribeAttributeUnitsRequest) (*pb.DescribeAttributeUnitsResponse, error) {
	attributeUnits, err := DescribeAttributeUnits(ctx, req)
	if err != nil {
		return nil, internalError(ctx, err)
	}

	var pbAttributeUnits []*pb.AttributeUnit
	for _, attUnit := range attributeUnits {
		pbAttributeUnits = append(pbAttributeUnits, models.AttributeUnitToPb(attUnit))
	}

	return &pb.DescribeAttributeUnitsResponse{AttributeUnits: pbAttributeUnits}, nil
}

func (s *Server) CreateAttribute(ctx context.Context, req *pb.CreateAttributeRequest) (*pb.CreateAttributeResponse, error) {
	attribute := models.PbToAttribute(req)

	//check if attribute_name exist
	err := checkStructExist(ctx, models.AttributeName{}, attribute.AttributeNameId)
	if err != nil {
		return nil, err
	}

	//check if attribute_unit exist
	if attribute.AttributeUnitId != "" {
		err = checkStructExist(ctx, models.AttributeUnit{}, attribute.AttributeUnitId)
		if err != nil {
			return nil, err
		}
	}

	//insert into attribute
	err = insertAttribute(ctx, attribute)
	if err != nil {
		return nil, internalError(ctx, err)
	}
	return &pb.CreateAttributeResponse{AttributeId: pbutil.ToProtoString(attribute.AttributeId)}, nil
}

func (s *Server) DescribeAttributes(ctx context.Context, req *pb.DescribeAttributesRequest) (*pb.DescribeAttributesResponse, error) {
	attributes, err := DescribeAttributes(ctx, req)
	if err != nil {
		return nil, internalError(ctx, err)
	}

	var pbAttributes []*pb.Attribute
	for _, att := range attributes {
		pbAttributes = append(pbAttributes, models.AttributeToPb(att))
	}

	return &pb.DescribeAttributesResponse{Attributes: pbAttributes}, nil
}

func (s *Server) CreateSpu(ctx context.Context, req *pb.CreateSpuRequest) (*pb.CreateSpuResponse, error) {
	spu := models.PbToSpu(req)

	attNameTmp := models.AttributeName{}
	//check if attribute_names exist
	for _, attNameId := range spu.AttributeNameIds {
		err := checkStructExist(ctx, attNameTmp, attNameId)
		if err != nil {
			return nil, err
		}
	}

	//insert spu
	err := insertSpu(ctx, spu)
	if err != nil {
		return nil, internalError(ctx, err)
	}
	return &pb.CreateSpuResponse{SpuId: pbutil.ToProtoString(spu.SpuId)}, nil
}

func (s *Server) CreateSku(ctx context.Context, req *pb.CreateSkuRequest) (*pb.CreateSkuResponse, error) {
	sku := models.PbToSku(req)

	//get spu
	spu, err := getSpu(ctx, sku.SpuId)
	if err != nil {
		return nil, internalError(ctx, err)
	}
	if spu == nil {
		return nil, notExistError(ctx, models.Spu{}, sku.SpuId)
	}

	//check if attributeNameId in sku.AttributeIds exist in spu.AttributeNameIds
	for _, attId := range sku.AttributeIds {
		attribute, err := getAttribute(ctx, attId)
		if err != nil {
			return nil, internalError(ctx, err)
		}
		if attribute == nil {
			return nil, notExistError(ctx, models.Attribute{}, attId)
		}
		if !stringutil.StringIn(attribute.AttributeNameId, spu.AttributeNameIds) {
			//attribute_name in sku.AttributeIds not exist in spu.AttribuateNameIds
			return nil, notExistInOtherError(ctx, models.AttributeName{}, models.Spu{})
		}

	}
	//check if attributeNameId in sku.MeteringAttributeIds exist in spu.MeteringAttributeNameIds
	for _, attId := range sku.MeteringAttributeIds {
		attribute, err := getAttribute(ctx, attId)
		if err != nil {
			return nil, internalError(ctx, err)
		}
		if attribute == nil {
			return nil, notExistError(ctx, models.Attribute{}, attId)
		}
		if !stringutil.StringIn(attribute.AttributeNameId, spu.MeteringAttributeNameIds) {
			//attribute_name in sku.MeteringAttributeIds not exist in spu.MeteringAttribuateNameIds
			return nil, notExistInOtherError(ctx, models.AttributeName{}, models.Spu{})
		}

	}

	//insert sku
	err = insertSku(ctx, sku)
	if err != nil {
		return nil, internalError(ctx, err)
	}
	return &pb.CreateSkuResponse{SkuId: pbutil.ToProtoString(sku.SkuId)}, nil
}

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

func renewalTimeFromSku(ctx context.Context, skuId string, actionTime time.Time) (*time.Time, error) {
	sku, err := getSku(ctx, skuId)

	if err != nil {
		logger.Error(ctx, "Failed to convert renewal time from sku, Error: [%+v]", err)
		return nil, err
	}

	//TODO: calculate renewalTime
	renewalTime := sku.CreateTime

	return &renewalTime, nil
}
