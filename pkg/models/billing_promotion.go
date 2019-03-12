// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"time"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/util/idutil"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
)

func NewCombinationPriceId() string {
	return idutil.GetUuid("comPrice-")
}

func NewProbationId() string {
	return idutil.GetUuid("pro-")
}

func NewDiscountId() string {
	return idutil.GetUuid("discount-")
}

func NewCouponId() string {
	return idutil.GetUuid("coupon-")
}

func NewCouponReceivedId() string {
	return idutil.GetUuid("couRec-")
}

type CombinationPrice struct {
	CombinationPriceId   string
	CombinationBindingId string
	Prices               map[int64]float64 //StepPrice: {upto: price, ..}
	Currency             string
	Status               string
	CreateTime           time.Time
	StatusTime           time.Time
}

func NewCombinationPrice(bindingId, currency string, prices map[int64]float64) *CombinationPrice {
	now := time.Now()
	return &CombinationPrice{
		CombinationPriceId:   NewCombinationPriceId(),
		CombinationBindingId: bindingId,
		Prices:               prices,
		Currency:             currency,
		Status:               constants.StatusActive,
		CreateTime:           now,
		StatusTime:           now,
	}
}

func PbToCombinationPrice(req *pb.CreateCombinationPriceRequest) *CombinationPrice {
	return NewCombinationPrice(
		req.GetCombinationBindingId().GetValue(),
		req.GetCurrency().String(),
		req.GetPrices(),
	)
}

//The requirements of this Probation:
//the startTime of using_sku is later than StartTime of Probation,
// and earlier than EndTime of Probation;
type Probation struct {
	ProbationId string
	SkuId       string
	AttributeId string
	Status      string //active/using/used
	StartTime   time.Time
	EndTime     time.Time
	CreateTime  time.Time
	StatusTime  time.Time
}

func NewProbation(skuId, attId string, startTime, endTime time.Time) *Probation {
	now := time.Now()
	if (time.Time{}) == startTime {
		startTime = now
	}
	return &Probation{
		ProbationId: NewProbationId(),
		SkuId:       skuId,
		AttributeId: attId,
		Status:      constants.StatusActive,
		StartTime:   startTime,
		EndTime:     endTime,
		CreateTime:  now,
		StatusTime:  now,
	}
}

func PbToProbation(req *pb.CreateProbationRequest) *Probation {
	return NewProbation(
		req.GetSkuId().GetValue(),
		req.GetAttributeId().GetValue(),
		pbutil.FromProtoTimestamp(req.GetStartTime()),
		pbutil.FromProtoTimestamp(req.GetEndTime()),
	)
}

type ProbationRecord struct {
	ProbationId string
	UserId      string
	Remain      float64
	Status      string
	CreateTime  time.Time
	StatusTime  time.Time
}

func NewProbationRecord(probationId, userId string, remain float64) *ProbationRecord {
	now := time.Now()
	return &ProbationRecord{
		ProbationId: probationId,
		UserId:      userId,
		Remain:      remain,
		Status:      constants.StatusActive,
		CreateTime:  now,
		StatusTime:  now,
	}
}

type Discount struct {
	DiscountId      string
	Name            string
	Owner           string
	LimitIds        []string
	DiscountValue   float64
	DiscountPercent float64
	Status          string
	Description     string
	StartTime       time.Time
	EndTime         time.Time
	CreateTime      time.Time
	StatusTime      time.Time
}

func NewDiscount(name, owner, description string,
	limitIds []string,
	disValue, disPercent float64,
	startTime, endTime time.Time) *Discount {

	now := time.Now()
	if (time.Time{}) == startTime {
		startTime = now
	}
	return &Discount{
		DiscountId:      NewDiscountId(),
		Name:            name,
		Owner:           owner,
		LimitIds:        limitIds,
		DiscountValue:   disValue,
		DiscountPercent: disPercent,
		Description:     description,
		Status:          constants.StatusActive,
		StartTime:       startTime,
		EndTime:         endTime,
		CreateTime:      now,
		StatusTime:      now,
	}
}

func PbToDiscount(req *pb.CreateDiscountRequest, owner string) *Discount {
	return NewDiscount(
		req.GetName().GetValue(),
		owner,
		req.GetDescription().GetValue(),
		req.GetLimitIds(),
		req.GetDiscountValue().GetValue(),
		req.GetDiscountPercent().GetValue(),
		pbutil.FromProtoTimestamp(req.GetStartTime()),
		pbutil.FromProtoTimestamp(req.GetEndTime()),
	)
}

type Coupon struct {
	CouponId    string
	Name        string
	Owner       string
	LimitIds    []string
	Balance     float64
	Count       uint32
	Remain      uint32
	LimitNumPer uint32
	Status      string
	StartTime   time.Time
	EndTime     time.Time
	CreateTime  time.Time
	StatusTime  time.Time
	Description string
}

func NewCoupon(name, owner, description string,
	limitIds []string,
	balance float64,
	count, limitNumPer uint32,
	startTime, endTime time.Time) *Coupon {

	now := time.Now()
	if (time.Time{}) == startTime {
		startTime = now
	}
	return &Coupon{
		CouponId:    NewCouponId(),
		Name:        name,
		Owner:       owner,
		LimitIds:    limitIds,
		Balance:     balance,
		Count:       count,
		Remain:      count,
		LimitNumPer: limitNumPer,
		Status:      constants.StatusActive,
		StartTime:   startTime,
		EndTime:     endTime,
		CreateTime:  now,
		StatusTime:  now,
		Description: description,
	}
}

func PbToCoupon(req *pb.CreateCouponRequest, owner string) *Coupon {
	return NewCoupon(
		req.GetName().GetValue(),
		owner,
		req.GetDescription().GetValue(),
		req.GetLimitIds(),
		req.GetBalance().GetValue(),
		req.GetCount().GetValue(),
		req.GetLimitNumPer().GetValue(),
		pbutil.FromProtoTimestamp(req.GetStartTime()),
		pbutil.FromProtoTimestamp(req.GetEndTime()),
	)
}

type CouponReceived struct {
	CouponReceivedId string
	CouponId         string
	UserId           string
	Remain           float64
	Status           string //active/using/used/overdue
	CreateTime       time.Time
	StatusTime       time.Time
}

func NewCouponReceived(couponId, userId string, remain float64) *CouponReceived {
	now := time.Now()
	return &CouponReceived{
		CouponReceivedId: NewCouponReceivedId(),
		CouponId:         couponId,
		UserId:           userId,
		Remain:           remain,
		Status:           constants.StatusActive,
		CreateTime:       now,
		StatusTime:       now,
	}
}

type CouponUsed struct {
	CouponUsedId     string
	CouponReceivedId string
	ContractId       string
	Balance          float64
	Currency         string
	Status           string //undetermined --> done / refunded
	CreateTime       string
}
