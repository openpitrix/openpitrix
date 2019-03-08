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
	//TODO: How to guarantee consistency operations.
	//MeteringResponse
	var leasingIds []string
	for _, l := range leasings {
		err = leasingToEtcd(*l)
		leasingIds = append(leasingIds, l.LeasingId)
	}
	return &pb.CommonMeteringResponse{LeasingIds: leasingIds}, nil
}

func (s *Server) UpdateMetering(ctx context.Context, req *pb.UpdateMeteringRequest) (*pb.CommonMeteringResponse, error) {
	userId := req.GetUserId().GetValue()
	for _, metering := range req.GetSkuMeterings() {
		leasing, _ := getLeasing(ctx,
			NIL_STR,
			userId,
			metering.GetResourceId().GetValue(),
			metering.GetSkuId().GetValue(),
		)

		//TODO: Update lesasing metering_values and save
		leasingToEtcd(*leasing)
	}

	return &pb.CommonMeteringResponse{}, nil
}

func (s *Server) StopMetering(ctx context.Context, req *pb.StopMeteringRequest) (*pb.CommonMeteringResponse, error) {
	userId := req.GetUserId().GetValue()

	for resourceId, skuId := range req.GetResourceSkuIds() {
		leasing, _ := getLeasing(ctx, NIL_STR, userId, resourceId, skuId)

		ClearLeasinRedis(leasing.LeasingId)
		//TODO: Update lesasing metering_values and save
		leasingToEtcd(*leasing)
	}
	return &pb.CommonMeteringResponse{}, nil
}

func (s *Server) TerminateMetering(ctx context.Context, req *pb.TerminateMeteringRequest) (*pb.CommonMeteringResponse, error) {
	//TODO: Do same as stopMetering

	//TODO: Move leasing to leased
	return &pb.CommonMeteringResponse{}, nil
}

//meteringValues: map<attributeId>value
func updateMeteringByRedis(ctx context.Context, leasingId string, updateTime time.Time) {

	//TODO: get leasing by leasingId
	leasing, _ := getLeasing(ctx, leasingId, NIL_STR, NIL_STR, NIL_STR)
	//TODO: update updataTIme and next renewalTime
	renewalTime := time.Now()

	//TODO: add to etcd
	leasingToEtcd(*leasing)
	leasingToRedis(leasingId, renewalTime)
	//TODO: guarantee consistency operations
}

func consumeLeasingRedis(ctx context.Context) {
	//TODO: consume due leasing from redis
	leasingId, updateTime := "", time.Now() //updateTIme: current renewalTime
	go updateMeteringByRedis(ctx, leasingId, updateTime)
}

func ClearLeasinRedis(leasingId string) error {
	//TODO: clear leasing in redis
	return nil
}

func GetGroupId() string {
	return "Group_01"
}
