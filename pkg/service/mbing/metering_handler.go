// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package mbing

import (
	"context"

	"openpitrix.io/openpitrix/pkg/pb"
)

func (s *Server) StartMetering(ctx context.Context, req *pb.MeteringRequest) (*pb.MeteringResponse, error) {
	return &pb.MeteringResponse{Status: 200, Message: "success"}, nil
}

func (s *Server) StopMetering(ctx context.Context, req *pb.MeteringRequest) (*pb.MeteringResponse, error) {
	return &pb.MeteringResponse{Status: 200, Message: "success"}, nil
}
