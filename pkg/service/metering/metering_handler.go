// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package metering

import (
	"context"
	"fmt"
	"time"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
)

func (s *Server) InitMetering(ctx context.Context, req *pb.InitMeteringRequest) (*pb.CommonMeteringResponse, error) {
	var leasings []*models.Leasing
	now := time.Now()
	for _, skuM := range req.GetSkuMeterings() {
		skuId := skuM.GetSkuId().GetValue()
		renewalTime, _ := renewalTimeFromSku(ctx, skuId, now)
		leasings = append(leasings, models.NewLeasing(skuM.GetValue(), req.GetUserId().GetValue(), req.GetResourceId().GetValue(), skuId, now, *renewalTime))
	}

	//insert leasings
	err := insertLeasings(ctx, leasings)
	if err != nil {
		return nil, internalError(ctx, err)
	}

	//TODO: How to guarantee consistency operations.
	for _, l := range leasings {
		err = leasingToRedis(*l)
		//TODO: check if BillingService exist, then add billing task for pre-charging by curl TaskService
	}
	return &pb.CommonMeteringResponse{ResourceId: req.GetResourceId()}, nil
}

//Can not update duration
func (s *Server) UpdateMetering(ctx context.Context, req *pb.UpdateMeteringRequest) (*pb.CommonMeteringResponse, error) {
	for _, metering := range req.GetUpdateSkuMeterings() {
		leasing, _ := getLeasing(ctx,
			NIL_STR,
			req.GetResourceId().GetValue(),
			metering.GetSkuId().GetValue(),
		)

		//TODO: Update lesasing metering_values and save leasing
		//      check attribute_name, make sure not duration
		leasingToRedis(*leasing)
		//TODO: check if BillingService exist and if need to charging, then add billing task for pre-charging by curl TaskService
	}

	return &pb.CommonMeteringResponse{}, nil
}

func (s *Server) StartMeterings(ctx context.Context, req *pb.StartMeteringsRequest) (*pb.CommonMeteringsResponse, error) {
	//TODO: get leasings by resource.ResourceId and resource.skuId(if skuIds is nil, get all leasings of resourceId)
	//      update Status of skus to active
	//TODO: Add leasing to REDIS

	return &pb.CommonMeteringsResponse{}, nil
}

//Before StopMetering, need to call UpdateMetering if needed.
func (s *Server) StopMeterings(ctx context.Context, req *pb.StopMeteringsRequest) (*pb.CommonMeteringsResponse, error) {
	var leasings []*models.Leasing

	for _, skuId := range req.GetSkuIds() {
		leasing, _ := getLeasing(ctx, NIL_STR, req.GetResourceId().GetValue(), skuId)
		leasings = append(leasings, leasing)
	}

	for _, leasing := range leasings {

		clearLeasingRedis(leasing.LeasingId)
		//TODO: Update UpdateTime/renewalTime/StopTimes/Status of leasing and save it
		//TODO: if pre-charging, add refund task by curl TaskService
	}
	return &pb.CommonMeteringsResponse{}, nil
}

func (s *Server) TerminateMeterings(ctx context.Context, req *pb.TerminateMeteringsRequest) (*pb.CommonMeteringsResponse, error) {
	var leasings []*models.Leasing

	for _, skuId := range req.GetSkuIds() {
		leasing, _ := getLeasing(ctx, NIL_STR, req.GetResourceId().GetValue(), skuId)
		leasings = append(leasings, leasing)
	}

	for _, leasing := range leasings {

		clearLeasingRedis(leasing.LeasingId)
		//TODO: Update UpdateTime/renewalTime/StopTimes/Status of leasing and save it
		//TODO: if pre-charging, add refund task by curl TaskService

		toLeased(leasing)
	}
	return &pb.CommonMeteringsResponse{}, nil
}

//meteringValues: map<attributeId>value
func updateMeteringByRedis(ctx context.Context, leasingId string, updateTime time.Time) {
	//TODO: get leasing by leasingId
	leasing, _ := getLeasing(ctx, leasingId, NIL_STR, NIL_STR)
	//TODO: update updataTIme and next renewalTime
	//TODO: add billing task by curl task service
	leasingToRedis(*leasing)
	//TODO: guarantee consistency operations
}

func ConsumeRedis(ctx context.Context) {
	//TODO: consume due leasing from redis
	leasingId, updateTime := "", time.Now() //updateTIme: current renewalTime
	go updateMeteringByRedis(ctx, leasingId, updateTime)
}

func clearLeasingRedis(leasingId string) error {
	//TODO: clear leasing in redis
	return nil
}

func toLeased(leasing *models.Leasing) error {
	leased := leasing.ToLeased()
	leasing.Status = constants.StatusDeleted
	//TODO: save leasing and leased
	fmt.Println(leased.LeasedId)
	return nil
}

func (s *Server) DescribeLeasings(ctx context.Context, req *pb.DescribeLeasingsRequest) (*pb.DescribeLeasingsResponse, error) {
	var leasings []*pb.Leasing

	//TODO: get leasings by DescribeLeasingsRequest

	return &pb.DescribeLeasingsResponse{LeasingSet: leasings}, nil
}
