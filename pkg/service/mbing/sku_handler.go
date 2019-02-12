// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package mbing

import (
	"context"

	"github.com/fatih/structs"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
)

func (s *Server) CreateAttribute(ctx context.Context, req *pb.CreateAttributeRequest) (*pb.CreateAttributeResponse, error) {
	att := models.PbToAttribute(req)
	err := insertAttribute(ctx, att)
	if err != nil {
		return nil, commonInternalErr(ctx, structs.Name(att), CreateFailedCode)
	}
	return &pb.CreateAttributeResponse{AttributeId: pbutil.ToProtoString(att.AttributeId)}, nil
}

func (s *Server) CreateAttributeUnit(ctx context.Context, req *pb.CreateAttUnitRequest) (*pb.CreateAttUnitResponse, error) {
	attUnit := models.PbToAttUnit(req)
	err := insertAttributeUnit(ctx, attUnit)
	if err != nil {
		return nil, commonInternalErr(ctx, structs.Name(attUnit), CreateFailedCode)
	}
	return &pb.CreateAttUnitResponse{AttributeUnitId: pbutil.ToProtoString(attUnit.AttributeUnitId)}, nil
}

func (s *Server) CreateAttributeValue(ctx context.Context, req *pb.CreateAttValueRequest) (*pb.CreateAttValueResponse, error) {
	attValue := models.PbToAttValue(req)

	actionStructName := structs.Name(attValue)
	//check if attribute exist
	err := checkStructExistById(ctx, structs.Name(models.Attribute{}), attValue.AttributeId, actionStructName, CreateFailedCode)
	if err != nil {
		return nil, err
	}

	//check if attribute_unit exist
	err = checkStructExistById(ctx, structs.Name(models.AttributeUnit{}), attValue.AttributeUnitId, actionStructName, CreateFailedCode)
	if err != nil {
		return nil, err
	}

	//insert into attribute_value
	err = insertAttributeValue(ctx, attValue)
	if err != nil {
		return nil, commonInternalErr(ctx, actionStructName, CreateFailedCode)
	}
	return &pb.CreateAttValueResponse{AttributeValueId: pbutil.ToProtoString(attValue.AttributeValueId)}, nil
}

func (s *Server) CreateResourceAttribute(ctx context.Context, req *pb.CreateResAttRequest) (*pb.CreateResAttResponse, error) {
	resAtt := models.PbToResAtt(req)

	actionStructName := structs.Name(resAtt)
	attName := structs.Name(models.Attribute{})
	//check if attributes exist
	for _, attId := range resAtt.Attributes {
		err := checkStructExistById(ctx, attName, attId, actionStructName, CreateFailedCode)
		if err != nil {
			return nil, err
		}
	}

	//insert resource_attribute
	err := insertResourceAttribute(ctx, resAtt)
	if err != nil {
		return nil, commonInternalErr(ctx, actionStructName, CreateFailedCode)
	}
	return &pb.CreateResAttResponse{ResourceAtrributeId: pbutil.ToProtoString(resAtt.ResourceAttributeId)}, nil
}

func (s *Server) CreateSku(ctx context.Context, req *pb.CreateSkuRequest) (*pb.CreateSkuResponse, error) {
	sku := models.PbToSku(req)

	actionStructName := structs.Name(sku)
	//check if resource_attribute exist
	err := checkStructExistById(
		ctx,
		structs.Name(models.ResourceAttribute{}),
		sku.ResourceAttributeId,
		actionStructName,
		CreateFailedCode)
	if err != nil {
		return nil, err
	}

	attValueName := structs.Name(models.AttributeValue{})
	//check if attribute_values exist
	for _, v := range sku.Values {
		err = checkStructExistById(ctx, attValueName, v, actionStructName, CreateFailedCode)
		if err != nil {
			return nil, err
		}
	}

	//insert resource_attribute
	err = insertSku(ctx, sku)
	if err != nil {
		return nil, commonInternalErr(ctx, actionStructName, CreateFailedCode)
	}
	return &pb.CreateSkuResponse{SkuId: pbutil.ToProtoString(sku.SkuId)}, nil
}

func (s *Server) CreatePrice(ctx context.Context, req *pb.CreatePriceRequest) (*pb.CreatePriceResponse, error) {
	price := models.PbToPrice(req)

	actionStructName := structs.Name(price)
	//check if resource_attribute exist
	err := checkStructExistById(
		ctx,
		structs.Name(models.ResourceAttribute{}),
		sku.ResourceAttributeId,
		actionStructName,
		CreateFailedCode)
	if err != nil {
		return nil, err
	}

	attValueName := structs.Name(models.AttributeValue{})
	//check if attribute_values exist
	for _, v := range sku.Values {
		err = checkStructExistById(ctx, attValueName, v, actionStructName, CreateFailedCode)
		if err != nil {
			return nil, err
		}
	}

	//insert resource_attribute
	err = insertSku(ctx, sku)
	if err != nil {
		return nil, commonInternalErr(ctx, actionStructName, CreateFailedCode)
	}
	return &pb.CreateSkuResponse{SkuId: pbutil.ToProtoString(sku.SkuId)}, nil
}
