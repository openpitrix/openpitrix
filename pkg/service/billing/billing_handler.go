// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package mbing

import (
	"context"
	"time"

	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
)

func (s *Server) CreatePrice(ctx context.Context, req *pb.CreatePriceRequest) (*pb.CreatePriceResponse, error) {
	price := models.PbToPrice(req)

	//TODO: how to check bindId
	//How about do not check bindId?

	//insert price
	err := insertPrice(ctx, price)
	if err != nil {
		return nil, internalError(ctx, err)
	}
	return &pb.CreatePriceResponse{PriceId: pbutil.ToProtoString(price.PriceId)}, nil
}

func (s *Server) DescribePrices(ctx context.Context, req *pb.DescribePricesRequest) (*pb.DescribePricesResponse, error) {
	//TODO: impl DescribePrices
	return &pb.DescribePricesResponse{}, nil
}

func (s *Server) ModifyPrice(ctx context.Context, req *pb.ModifyPriceRequest) (*pb.ModifyPriceResponse, error) {
	//TODO: impl ModifyPrice
	return &pb.ModifyPriceResponse{}, nil
}

func (s *Server) DeletePrices(ctx context.Context, req *pb.DeletePricesRequest) (*pb.DeletePricesResponse, error) {
	//TODO: impl DeletePrices
	return &pb.DeletePricesResponse{}, nil
}

type Metering struct {
	LeasingId      string
	SkuId          string
	UserId         string
	Action         string //start/update/stop/terminate metering
	MeteringValues map[string]float64
	UpdateTime     time.Time
}

func ConsumeEtcd(ctx context.Context) {
	//TODO: get Metering from etcd
	metering := Metering{}
	Billing(ctx, metering)

}

func Billing(ctx context.Context, metering Metering) error {
	//get LeasingContract
	var contract *models.LeasingContract
	if metering.Action == "start" {
		//TODO: new LeasingContract and set Status to updating
		contract = &models.LeasingContract{}
		insertLeasingContract(ctx, contract)
	} else {
		//TODO: get LeasingContract by leasingId
		contract, _ = getLeasingContract(ctx, "", metering.LeasingId)
		//TODO: update MeteringValues/Status(updating: incase not finished billing) of LeasingContract and save it
	}

	switch metering.Action {
	case "start":
		calculate(ctx, metering, contract)
	case "update":
		calculate(ctx, metering, contract)
	case "stop":
		//TODO: update Status of LeasingContract to stoped
		//TODO: if duration in MeteringValues, reCalculate and reCharge
		//????: How to handle Coupon ?
	case "terminate":
		//TODO: update Status of LeasingContract to deleted
		//TODO: if duration in MeteringValues, reCalculate and reCharge
		//????: How to handle Coupon ?
	}

	//charge due fee
	if contract.DueFee > 0 {
		_, err := charge(contract)
		if err.Error() == "balance not enough" {
			insufficientBalanceToEtcd(contract.ResourceId, contract.SkuId, contract.UserId)
		}
	}

	//recharge
	if contract.DueFee < 0 {
		reCharge(contract)
	}

	return nil
}

func calculate(ctx context.Context, metering Metering, contract *models.LeasingContract){
	//************************ main process ***********************
	for attId, value := range metering.MeteringValues {
		probationValue := probationFromSku(contract.SkuId, contract.UserId, contract.StartTime)
		//TODO: if the value < probationValue, log for using probation and update Status to active.

		//loggor.info(...)
		//update status of to active

		//TODO: if the value > probationValue, get real price and calculate the fee
		//    1. update status of ProbationRecord to used
		//    2. get Price by skuId and attId, get price by value_interval ----- priceFromSku
		//    3. oldValue = contract.MeteringValues[attId]
		//       billingMeteringValue = oldValue>probation?value-oldValue:value-probation
		//    4. get discount by spuId/skuId/priceId and startTime/endTime ----- discountFromSku
		//    5. get price by value from Price, price
		//       realPrice = price*Discount.DiscountPercent or price-Discount.DiscountValue

		//TODO: calculate dueFee = dueFee + billingMeteringValue * realPrice
		deductCoupon(contract.UserId, "", contract.SkuId, "")
	}
}


func probationFromSku(skuId, userId string, endTime time.Time) float64 {
	//TODO: get Probation and ProbationRecord by skuId and userId
	//TODO: if ProbationRecord not exist, create ProbationRecord and set status to using
	//TODO: check if the probation used by the user:
	//      if used, return 0,
	//      if not, get value by attributeId and return it
	return 0.0
}

func priceFromSku(skuId, attributeId string) *models.Price {
	//TODO: get Metering_Attribute_Binding by contract(skuId, meteringAttributeId)
	//TODO: get Price by binding_id of Metering_Attribute_Binding
	return &models.Price{}
}

//TODO: Make sure the discount requirement with PM
func discountFromSku(spuId, skuId, priceId string, startTime, endTime time.Time) (*models.Discount, error) {
	return &models.Discount{}, nil
}

//TODO: Make sure the coupon
func deductCoupon(userId, spuId, skuId, priceId string) error {
	//TODO: get CouponReceived by UserId and Status
	//TODO: get Coupon by CouponId
	//TODO: if the spuId/skuId/priceId in Coupon.Limit_ids,
	//      update contract.dueFee and CouponReceived.Remain,
	//      if CouponReceived.Remain is 0, update Status of CouponReceived to used.
	return nil
}
