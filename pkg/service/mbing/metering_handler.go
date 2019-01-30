// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package mbing

import (
	"context"

	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
)

func (s *Server) StartMetering(ctx context.Context, req *pb.MeteringRequest) (*pb.CommonResponse, error) {
	skuAttribute := "in sku"
	if "duration" == skuAttribute {
	}
	return &pb.CommonResponse{Status: pbutil.ToProtoInt32(200), Message: pbutil.ToProtoString("success")}, nil
}

func (s *Server) StopMetering(ctx context.Context, req *pb.MeteringRequest) (*pb.CommonResponse, error) {
	return &pb.CommonResponse{Status: pbutil.ToProtoInt32(200), Message: pbutil.ToProtoString("success")}, nil
}

func (s *Server) UpdateMetering(ctx context.Context, req *pb.MeteringRequest) (*pb.CommonResponse, error) {
	return &pb.CommonResponse{Status: pbutil.ToProtoInt32(200), Message: pbutil.ToProtoString("success")}, nil
}

func (s *Server) CloseMetering(ctx context.Context, req *pb.MeteringRequest) (*pb.CommonResponse, error) {
	return &pb.CommonResponse{Status: pbutil.ToProtoInt32(200), Message: pbutil.ToProtoString("success")}, nil
}

func GetGroupId() string {
	return "Group_01"
}
