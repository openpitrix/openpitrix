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
)

func (s *Server) CreateAttribute(ctx context.Context, req *pb.CreateAttributeRequest) (*pb.CreateAttributeResponse, error) {
	att := models.PbToAttribute(req)
	err := insertAttribute(ctx, att)
	if err != nil {
		return nil, commonInternalErr(ctx, att, CreateFailedCode)
	}
	return &pb.CreateAttributeResponse{AttributeId: pbutil.ToProtoString(att.AttributeId)}, nil
}

func (s *Server) ListAttribute(ctx context.Context, req *pb.ListAttributeRequest) (*pb.ListAttributeResponse, error) {
	atts, err := listAttribute(ctx, req)
	if err != nil {
		return nil, commonInternalErr(ctx, models.Attribute{}, ListFailedCode)
	}

	var pbAtts []*pb.Attribute
	for _, att := range atts {
		pbAtts = append(pbAtts, models.AttributeToPb(att))
	}

	return &pb.ListAttributeResponse{Attributes: pbAtts}, nil
}

func (s *Server) CreateAttributeUnit(ctx context.Context, req *pb.CreateAttUnitRequest) (*pb.CreateAttUnitResponse, error) {
	attUnit := models.PbToAttUnit(req)
	err := insertAttributeUnit(ctx, attUnit)
	if err != nil {
		return nil, commonInternalErr(ctx, attUnit, CreateFailedCode)
	}
	return &pb.CreateAttUnitResponse{AttributeUnitId: pbutil.ToProtoString(attUnit.AttributeUnitId)}, nil
}

func (s *Server) ListAttributeUnit(ctx context.Context, req *pb.ListAttUnitRequest) (*pb.ListAttUnitResponse, error) {
	attUnits, err := listAttributeUnit(ctx, req)
	if err != nil {
		return nil, commonInternalErr(ctx, models.Attribute{}, ListFailedCode)
	}

	var pbAttUnits []*pb.AttributeUnit
	for _, attUnit := range attUnits {
		pbAttUnits = append(pbAttUnits, models.AttributeUnitToPb(attUnit))
	}

	return &pb.ListAttUnitResponse{AttributeUnits: pbAttUnits}, nil
}

func (s *Server) CreateAttributeValue(ctx context.Context, req *pb.CreateAttValueRequest) (*pb.CreateAttValueResponse, error) {
	attValue := models.PbToAttValue(req)

	//check if attribute exist
	err := checkStructExistById(ctx, models.Attribute{}, attValue, attValue.AttributeId, CreateFailedCode)
	if err != nil {
		return nil, err
	}

	//check if attribute_unit exist
	err = checkStructExistById(ctx, models.AttributeUnit{}, attValue, attValue.AttributeUnitId, CreateFailedCode)
	if err != nil {
		return nil, err
	}

	//insert into attribute_value
	err = insertAttributeValue(ctx, attValue)
	if err != nil {
		return nil, commonInternalErr(ctx, attValue, CreateFailedCode)
	}
	return &pb.CreateAttValueResponse{AttributeValueId: pbutil.ToProtoString(attValue.AttributeValueId)}, nil
}

func (s *Server) CreateResourceAttribute(ctx context.Context, req *pb.CreateResAttRequest) (*pb.CreateResAttResponse, error) {
	resAtt := models.PbToResAtt(req)

	att := models.Attribute{}
	//check if attributes exist
	for _, attId := range resAtt.Attributes {
		err := checkStructExistById(ctx, att, resAtt, attId, CreateFailedCode)
		if err != nil {
			return nil, err
		}
	}

	//insert resource_attribute
	err := insertResourceAttribute(ctx, resAtt)
	if err != nil {
		return nil, commonInternalErr(ctx, resAtt, CreateFailedCode)
	}
	return &pb.CreateResAttResponse{ResourceAtrributeId: pbutil.ToProtoString(resAtt.ResourceAttributeId)}, nil
}

func (s *Server) CreateSku(ctx context.Context, req *pb.CreateSkuRequest) (*pb.CreateSkuResponse, error) {
	sku := models.PbToSku(req)

	//check if resource_attribute exist
	err := checkStructExistById(
		ctx,
		models.ResourceAttribute{},
		sku,
		sku.ResourceAttributeId,
		CreateFailedCode)
	if err != nil {
		return nil, err
	}

	attValue := models.AttributeValue{}
	//check if attribute_values exist
	for _, VId := range sku.Values {
		err = checkStructExistById(ctx, attValue, sku, VId, CreateFailedCode)
		if err != nil {
			return nil, err
		}
	}

	//insert sku
	err = insertSku(ctx, sku)
	if err != nil {
		return nil, commonInternalErr(ctx, sku, CreateFailedCode)
	}
	return &pb.CreateSkuResponse{SkuId: pbutil.ToProtoString(sku.SkuId)}, nil
}

func (s *Server) CreatePrice(ctx context.Context, req *pb.CreatePriceRequest) (*pb.CreatePriceResponse, error) {
	price := models.PbToPrice(req)

	//check if sku exist
	err := checkStructExistById(ctx, models.Sku{}, price, price.SkuId, CreateFailedCode)
	if err != nil {
		return nil, err
	}

	//check if attribute exist
	err = checkStructExistById(ctx, models.Attribute{}, price, price.AttributeId, CreateFailedCode)
	if err != nil {
		return nil, err
	}

	//check if attribute_values exist
	attValue := models.AttributeValue{}
	for k := range price.Prices {
		err = checkStructExistById(ctx, attValue, price, k, CreateFailedCode)
		if err != nil {
			return nil, err
		}
	}

	//insert price
	err = insertPrice(ctx, price)
	if err != nil {
		return nil, commonInternalErr(ctx, price, CreateFailedCode)
	}
	return &pb.CreatePriceResponse{PriceId: pbutil.ToProtoString(price.PriceId)}, nil
}

func renewalTimeFromSku(ctx context.Context, skuId string, actionTime time.Time) (*time.Time, error) {
	sku, err := getSkuById(ctx, skuId)

	if err != nil {
		logger.Error(ctx, "Failed to convert renewal time from sku, Error: [%+v]", err)
		return nil, err
	}

	//TODO: calculate renewalTime
	renewalTime := sku.CreateTime

	return &renewalTime, nil
}