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

func (s *Server) CreateAccount(ctx context.Context, req *pb.CreateAccountRequest) (*pb.CreateAccountResponse, error) {
	//TODO: impl CreateAccount
	return &pb.CreateAccountResponse{}, nil
}

func (s *Server) DescribeAccounts(ctx context.Context, req *pb.DescribeAccountsRequest) (*pb.DescribeAccountsResponse, error) {
	//TODO: impl DescribeAccounts
	return &pb.DescribeAccountsResponse{}, nil
}

func (s *Server) ModifyAccount(ctx context.Context, req *pb.ModifyAccountRequest) (*pb.ModifyAccountResponse, error) {
	//TODO: impl ModifyAccount
	return &pb.ModifyAccountResponse{}, nil
}

func (s *Server) DeleteAccounts(ctx context.Context, req *pb.DeleteAccountsRequest) (*pb.DeleteAccountsResponse, error) {
	//TODO: impl DeleteAccounts
	return &pb.DeleteAccountsResponse{}, nil
}

type Metering struct {
	LeasingId      string
	ReourceId      string
	SkuId          string
	UserId         string
	Action         string //start/update/stop/terminate metering
	MeteringValues map[string]float64
	UpdateTime     time.Time
}

//Loop running...
func EtcdConeumer(ctx context.Context) {
	//TODO: get Metering from etcd
	metering := Metering{}
	Billing(ctx, metering)

}

func CreateLeasingContract(ctx context.Context, metering Metering) *models.LeasingContract {
	//TODO: new LeasingContract and set Status to updating
	contract := &models.LeasingContract{}
	insertLeasingContract(ctx, contract)
	return contract
}

func Billing(ctx context.Context, metering Metering) {
	//get LeasingContract
	var contract *models.LeasingContract
	if metering.Action == "init" {
		contract = CreateLeasingContract(ctx, metering)
	} else {
		//TODO: get LeasingContract by leasingId
		contract, _ = getLeasingContract(ctx, "", metering.LeasingId)
		//TODO: update MeteringValues/Status(updating: incase not finished billing) of LeasingContract and save it
	}

	//calculate due_fee by value
	calculate(ctx, metering, contract)

	//deduct coupon
	if contract.DueFee > 0 {
		deductCoupon(contract)
	}
	if contract.DueFee < 0 {
		refundCoupon(contract)
	}

	//charge due_fee from account
 	if contract.DueFee > 0 {
		_, err := charge(contract)
		if err != nil {
			if err.Error() == "balance not enough" {
				insufficientBalanceToEtcd(contract.ResourceId, contract.SkuId, contract.UserId)
			}
			//TODO: send alert email
			return
		}
	}

	//recharge if attribute_name is duration and attribute_unit is hour
	if contract.DueFee < 0 {
		refund(contract)
	}

	switch metering.Action {
	case "start":
		fallthrough
	case "update":
		//TODO: set Status of LeasingContract to active and save it
	case "stop":
		//TODO: update Status of LeasingContract to stoped and save it
	case "terminate":
		//TODO: update Status of LeasingContract to deleted
		//TODO: Move LeasingContract to LeasedContract
	}

}

//Update DueFee / Probation / Coupon
func calculate(ctx context.Context, metering Metering, contract *models.LeasingContract) error {
	//************************ main process ***********************
	for attId, value := range metering.MeteringValues {
		//get billingValue
		oldValue := contract.MeteringValues[attId]
		billingMeteringValue := value - oldValue
		if attId == "att-00001" { //the attribute is duration
			probationValue := probationFromSku(contract.SkuId, contract.UserId, contract.StartTime)
			//TODO: if the value < probationValue, log for using probation and update Status to active.
			if value < probationValue {
				//loggor.info(...)
				//update status of to active
				billingMeteringValue = 0
			} else if oldValue < probationValue {
				billingMeteringValue = oldValue - probationValue
				//TODO: update status of ProbationRecord to used
			}
		}

		//caculate dueFee
		if billingMeteringValue > 0 {
			//TODO:
			// 1. get Price by skuId and attId, get price by value_interval ----- priceFromSku
			// 2. get discount by spuId/skuId/priceId and startTime/endTime ----- discountFromSku
			// 3. get price by value from Price, price for special value
			//    realPrice = price*Discount.DiscountPercent or price-Discount.DiscountValue
			realPrice := 1.0

			//TODO: calculate dueFee = dueFee + billingMeteringValue * realPrice
			contract.DueFee = contract.DueFee + billingMeteringValue*realPrice
		}

	}
	return nil
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

//TODO: Make sure condition of the coupon
//TODO: Make sure the sequence(eg: Remain, EndTime) of coupons to deduct for due_fee
func deductCoupon(contract *models.LeasingContract) error {
	//TODO: update Status of CouponUsed used by contract from undetermined to done

	//TODO: get CouponReceiveds by UserId and Status (active and using)
	//TODO: get Coupons by CouponIds
	//TODO: if the spuId/skuId/priceId is in Coupon.Limit_ids and equal:
	//          update contract.dueFee and Remain/Status of CouponReceived,
	//          generate CouponUsed and set status to undetermined
	return nil
}

func refundCoupon(contract *models.LeasingContract) error {
	//TODO: update Status of CouponUsed of the contract from undetermined to deleted or done by DueFee and Balance
	//TODO: update The Remian/Status of CouponReceived
	return nil
}
