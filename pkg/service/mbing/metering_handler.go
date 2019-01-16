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

func (s *Server) StartMetering(ctx context.Context, req *pb.MeteringRequest) (*pb.MeteringResponse, error) {

	err := insertLeasingsToDB(ctx, models.MeteringReq2Leasings(req, getGroupId()))

	if err != nil {
		return &pb.MeteringResponse{Status: pbutil.ToProtoInt32(400), Message: pbutil.ToProtoString(err.Error())}, nil
	}
	return &pb.MeteringResponse{Status: pbutil.ToProtoInt32(200), Message: pbutil.ToProtoString("success")}, nil
}

func (s *Server) StopMetering(ctx context.Context, req *pb.MeteringRequest) (*pb.MeteringResponse, error) {
	return &pb.MeteringResponse{Status: pbutil.ToProtoInt32(200), Message: pbutil.ToProtoString("success")}, nil
}

func getGroupId() string {
	return "Group-1"
}
