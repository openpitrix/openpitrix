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

type Probation struct {
	ProbationId string
	SkuId       string
	AttributeId string
	Status      string
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
	ProbationSkuId string
	UserId         string
	LimitNum       uint32
	CreateTime     time.Time
	ProbationTimes []time.Time
}

type Discount struct {
	DiscountId      string
	Name            string
	Limits          map[string]string
	DiscountValue   float64
	DiscountPercent float64
	StartTime       time.Time
	EndTime         time.Time
	CreateTime      time.Time
	Status          string
	Mark            string
}

type Coupon struct {
	CouponId   string
	Name       string
	Limits     map[string]string
	Balance    float64
	Count      uint32
	LimitNum   uint32
	StartTime  time.Time
	EndTime    time.Time
	CreateTime time.Time
	Status     string
	Mark       string
}

type CouponReceived struct {
	CouponReceivedId string
	CouponId         string
	UserId           string
	Balance          float64
	Status           string
	CreateTime       time.Time
}
