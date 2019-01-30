// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package mbing

import (
	"context"

	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
)

func (s *Server) CreateCRA(ctx context.Context, req *pb.CreateCRARequest) (*pb.CommonResponse, error) {
	return &pb.CommonResponse{Status: pbutil.ToProtoInt32(200), Message: pbutil.ToProtoString("success")}, nil
}

func (s *Server) CreateCombinationSku(ctx context.Context, req *pb.CreateCSRequest) (*pb.CommonResponse, error) {
	return &pb.CommonResponse{Status: pbutil.ToProtoInt32(200), Message: pbutil.ToProtoString("success")}, nil
}

func (s *Server) CreateCombinationPrices(ctx context.Context, req *pb.CreateCPRequest) (*pb.CommonResponse, error) {
	return &pb.CommonResponse{Status: pbutil.ToProtoInt32(200), Message: pbutil.ToProtoString("success")}, nil
}

func (s *Server) CreateProbationSku(ctx context.Context, req *pb.CreatePSRequest) (*pb.CommonResponse, error) {
	return &pb.CommonResponse{Status: pbutil.ToProtoInt32(200), Message: pbutil.ToProtoString("success")}, nil
}

func (s *Server) CreateProbationRecord(ctx context.Context, req *pb.CreatePRRequest) (*pb.CommonResponse, error) {
	return &pb.CommonResponse{Status: pbutil.ToProtoInt32(200), Message: pbutil.ToProtoString("success")}, nil
}
