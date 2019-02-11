// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package mbing

import (
	"context"

	"openpitrix.io/openpitrix/pkg/gerr"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
)

func (s *Server) CreateAttribute(ctx context.Context, req *pb.CreateAttributeRequest) (*pb.CreateAttributeResponse, error) {
	att := models.PbToAttribute(req)
	err := insertAttribute(ctx, att)
	if err != nil {
		return nil, gerr.New(ctx, gerr.Internal, gerr.CreateFailed("attribute", "资源属性"))
	}
	return &pb.CreateAttributeResponse{AttributeId: pbutil.ToProtoString(att.AttributeId)}, nil
}

func (s *Server) CreateAttributeUnit(ctx context.Context, req *pb.CreateAttUnitRequest) (*pb.CreateAttUnitResponse, error) {
	attUnit := models.PbToAttUnit(req)
	err := insertAttributeUnit(ctx, attUnit)
	if err != nil {
		return nil, gerr.New(ctx, gerr.Internal, gerr.CreateFailed("attribute_unit", "属性单位"))
	}
	return &pb.CreateAttUnitResponse{AttributeUnitId: pbutil.ToProtoString(attUnit.AttributeUnitId)}, nil
}

func (s *Server) CreateAttributeValue(ctx context.Context, req *pb.CreateAttValueRequest) (*pb.CreateAttValueResponse, error) {
	attValue := models.PbToAttValue(req)

	_, err := getAttributeById(ctx, attValue.AttributeId)
	if err != nil {
		return nil, gerr.New(ctx, gerr.Internal, gerr.ErrorAttributeNotExist)
	}

	_, err = getAttUnitById(ctx, attValue.AttributeUnitId)
	if err != nil {
		return nil, gerr.New(ctx, gerr.Internal, gerr.ErrorAttUnitNotExist)
	}

	err = insertAttributeValue(ctx, attValue)
	if err != nil {
		return nil, gerr.New(ctx, gerr.Internal, gerr.CreateFailed("attribute_value", "属性值"))
	}
	return &pb.CreateAttValueResponse{AttributeValueId: pbutil.ToProtoString(attValue.AttributeValueId)}, nil
}



func (s *Server) CreateResourceAttributes(ctx context.Context, req *pb.CreateResourceAttributesRequest) (*pb.CommonResponse, error) {
	return &pb.CommonResponse{Status: pbutil.ToProtoInt32(200), Message: pbutil.ToProtoString("success")}, nil
}

func (s *Server) CreateSkus(ctx context.Context, req *pb.CreateSkusRequest) (*pb.CommonResponse, error) {
	return &pb.CommonResponse{Status: pbutil.ToProtoInt32(200), Message: pbutil.ToProtoString("success")}, nil
}

func (s *Server) CreatePrices(ctx context.Context, req *pb.CreatePricesRequest) (*pb.CommonResponse, error) {
	return &pb.CommonResponse{Status: pbutil.ToProtoInt32(200), Message: pbutil.ToProtoString("success")}, nil
}
