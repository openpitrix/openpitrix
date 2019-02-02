// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package mbing

import (
	"context"

	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
)

func (s *Server) CreateAttributes(ctx context.Context, req *pb.CreateAttributesRequest) (*pb.CommonResponse, error) {
	return &pb.CommonResponse{Status: pbutil.ToProtoInt32(200), Message: pbutil.ToProtoString("success")}, nil
}

func (s *Server) CreateAttributeUnits(ctx context.Context, req *pb.CreateAttributeUnitsRequest) (*pb.CommonResponse, error) {
	return &pb.CommonResponse{Status: pbutil.ToProtoInt32(200), Message: pbutil.ToProtoString("success")}, nil
}

func (s *Server) CreateAttributeValues(ctx context.Context, req *pb.CreateAttributeValuesRequest) (*pb.CommonResponse, error) {
	return &pb.CommonResponse{Status: pbutil.ToProtoInt32(200), Message: pbutil.ToProtoString("success")}, nil
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
