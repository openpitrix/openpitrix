// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package metering

import (
	"context"
	"time"

	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
)

func (s *Server) StartMetering(ctx context.Context, req *pb.StartMeteringRequest) (*pb.CommonMeteringResponse, error) {
	var leasings []*models.Leasing
	now := time.Now()
	for _, metering := range req.GetSkuMeterings() {
		renewaltime, _ := renewalTimeFromSku(ctx, metering.GetSkuId().GetValue(), now)
		leasings = append(leasings, models.NewLeasing(metering, GetGroupId(), req.GetUserId().GetValue(), now, *renewaltime))
	}

	//insert leasings
	err := insertLeasings(ctx, leasings)
	if err != nil {
		return nil, internalError(ctx, err)
	}

	//TODO: Add leasing to REDIS if duration exist.

	//MeteringResponse
	var leasingIds []string
	for _, l := range leasings {
		err = leasingToEtcd(*l)
		leasingIds = append(leasingIds, l.LeasingId)
	}
	return &pb.CommonMeteringResponse{LeasingIds: leasingIds}, nil
}

func (s *Server) UpdateMetering(ctx context.Context, req *pb.UpdateMeteringRequest) (*pb.CommonMeteringResponse, error) {
	return &pb.CommonMeteringResponse{}, nil
}

func (s *Server) StopMetering(ctx context.Context, req *pb.StopMeteringRequest) (*pb.CommonMeteringResponse, error) {
	return &pb.CommonMeteringResponse{}, nil
}

func (s *Server) TerminateMetering(ctx context.Context, req *pb.TerminateMeteringRequest) (*pb.CommonMeteringResponse, error) {
	return &pb.CommonMeteringResponse{}, nil
}

func GetGroupId() string {
	return "Group_01"
}
