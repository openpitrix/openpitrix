// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package metering

import (
	"context"
	"time"

	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
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

func (s *Server) ModifyAttributeName(ctx context.Context, req *pb.ModifyAttributeNameRequest) (*pb.ModifyAttributeNameResponse, error) {
	//TODO: impl ModifyAttributeName
	return &pb.ModifyAttributeNameResponse{AttributeNameId: pbutil.ToProtoString("tmp")}, nil
}

func (s *Server) DeleteAttributeNames(ctx context.Context, req *pb.DeleteAttributeNamesRequest) (*pb.DeleteAttributeNamesResponse, error) {
	//TODO: impl DeleteAttributeNames
	return &pb.DeleteAttributeNamesResponse{}, nil
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

func (s *Server) ModifyAttributeUnit(ctx context.Context, req *pb.ModifyAttributeUnitRequest) (*pb.ModifyAttributeUnitResponse, error) {
	//TODO: impl ModifyAttributeUnit
	return &pb.ModifyAttributeUnitResponse{}, nil
}

func (s *Server) DeleteAttributeUnits(ctx context.Context, req *pb.DeleteAttributeUnitsRequest) (*pb.DeleteAttributeUnitsResponse, error) {
	//TODO: impl DeleteAttributeUnits
	return &pb.DeleteAttributeUnitsResponse{}, nil
}

func (s *Server) CreateAttribute(ctx context.Context, req *pb.CreateAttributeRequest) (*pb.CreateAttributeResponse, error) {
	//TODO: get id of current user
	attribute := models.PbToAttribute(req, "")

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

func (s *Server) ModifyAttribute(ctx context.Context, req *pb.ModifyAttributeRequest) (*pb.ModifyAttributeResponse, error) {
	//TODO: impl ModifyAttribute
	return &pb.ModifyAttributeResponse{}, nil
}

func (s *Server) DeleteAttributes(ctx context.Context, req *pb.DeleteAttributesRequest) (*pb.DeleteAttributesResponse, error) {
	//TODO: impl DeleteAttributes
	return &pb.DeleteAttributesResponse{}, nil
}

func (s *Server) CreateSpu(ctx context.Context, req *pb.CreateSpuRequest) (*pb.CreateSpuResponse, error) {
	//TODO: get id of current user
	spu := models.PbToSpu(req, "")

	//insert spu
	err := insertSpu(ctx, spu)
	if err != nil {
		return nil, internalError(ctx, err)
	}
	return &pb.CreateSpuResponse{SpuId: pbutil.ToProtoString(spu.SpuId)}, nil
}

func (s *Server) DescribeSpus(ctx context.Context, req *pb.DescribeSpusRequest) (*pb.DescribeSpusResponse, error) {
	//TODO: impl DescribeSpus
	return &pb.DescribeSpusResponse{}, nil
}

func (s *Server) ModifySpu(ctx context.Context, req *pb.ModifySpuRequest) (*pb.ModifySpuResponse, error) {
	//TODO: impl ModifySpu
	return &pb.ModifySpuResponse{}, nil
}

func (s *Server) DeleteSpus(ctx context.Context, req *pb.DeleteSpusRequest) (*pb.DeleteSpusResponse, error) {
	//TODO: impl DeleteSpus
	return &pb.DeleteSpusResponse{}, nil
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

	//TODO: check attribute_ids

	//insert sku
	err = insertSku(ctx, sku)
	if err != nil {
		return nil, internalError(ctx, err)
	}
	return &pb.CreateSkuResponse{SkuId: pbutil.ToProtoString(sku.SkuId)}, nil
}

func (s *Server) DescribeSkus(ctx context.Context, req *pb.DescribeSkusRequest) (*pb.DescribeSkusResponse, error) {
	//TODO: impl DescribeSkus
	return &pb.DescribeSkusResponse{}, nil
}

func (s *Server) ModifySku(ctx context.Context, req *pb.ModifySkuRequest) (*pb.ModifySkuResponse, error) {
	//TODO: impl ModifySku
	return &pb.ModifySkuResponse{}, nil
}

func (s *Server) DeleteSkus(ctx context.Context, req *pb.DeleteSkusRequest) (*pb.DeleteSkusResponse, error) {
	//TODO: impl DeleteSkus
	return &pb.DeleteSkusResponse{}, nil
}

func (s *Server) CreateMeteringAttributeBindings(ctx context.Context, req *pb.CreateMeteringAttributeBindingsRequest) (*pb.CreateMeteringAttributeBindingsResponse, error) {
	//TODO: impl CreateMeteringAttributeBindings
	return &pb.CreateMeteringAttributeBindingsResponse{}, nil
}

func (s *Server) DescribeMeteringAttributeBindings(ctx context.Context, req *pb.DescribeMeteringAttributeBindingsRequest) (*pb.DescribeMeteringAttributeBindingsResponse, error) {
	//TODO: impl DescribeMeteringAttributeBindings
	return &pb.DescribeMeteringAttributeBindingsResponse{}, nil
}

func (s *Server) ModifyMeteringAttributeBinding(ctx context.Context, req *pb.ModifyMeteringAttributeBindingRequest) (*pb.ModifyMeteringAttributeBindingResponse, error) {
	//TODO: impl ModifyMeteringAttributeBinding
	return &pb.ModifyMeteringAttributeBindingResponse{}, nil
}

func (s *Server) DeleteMeteringAttributeBindings(ctx context.Context, req *pb.DeleteMeteringAttributeBindingsRequest) (*pb.DeleteMeteringAttributeBindingsResponse, error) {
	//TODO: impl DeleteMeteringAttributeBindings
	return &pb.DeleteMeteringAttributeBindingsResponse{}, nil
}

func renewalTimeFromSku(ctx context.Context, skuId string, actionTime time.Time) (*time.Time, error) {
	sku, err := getSku(ctx, skuId)

	if err != nil {
		logger.Error(ctx, "Failed to convert renewal time from sku, Error: [%+v]", err)
		return nil, err
	}

	//TODO: calculate renewalTime
	//TODO: check if duration in metering_attributes
	renewalTime := sku.CreateTime

	return &renewalTime, nil
}
