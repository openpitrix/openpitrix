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

func NewCombinationSpuId() string {
	return idutil.GetUuid("cra-")
}

func NewCombinationSkuId() string {
	return idutil.GetUuid("comSku-")
}

func NewCombinationPriceId() string {
	return idutil.GetUuid("comPrice-")
}

func NewProbationSkuId() string {
	return idutil.GetUuid("proSku-")
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

type CombinationSpu struct {
	CombinationSpuId string
	SpuIds           []string //the id slice of spu
	CreateTime       time.Time
	StatusTime       time.Time
	Status           string
}

func PbToCombinationSpu(req *pb.CreateCombinationSpuRequest) *CombinationSpu {
	return &CombinationSpu{
		CombinationSpuId: NewCombinationSpuId(),
		SpuIds:           pbutil.FromProtoStringSlice(req.GetSpuIds()),
	}
}

type CombinationSku struct {
	CombinationSkuId     string
	CombinationSpuId     string
	AttributeIds         map[string][]string //[resourceVersionId: [attributeId, ..], ..}
	MeteringAttributeIds map[string][]string //[resourceVersionId: [attributeId, ..], ..}
	CreateTime           time.Time
	StatusTime           time.Time
	Status               string
}

func NewCombinationSku(comSpuId string, attributeIds, meteringAttributeIds map[string][]string) *CombinationSku {
	now := time.Now()
	return &CombinationSku{
		CombinationSkuId:     NewCombinationSkuId(),
		CombinationSpuId:     comSpuId,
		AttributeIds:         attributeIds,
		MeteringAttributeIds: meteringAttributeIds,
		CreateTime:           now,
		StatusTime:           now,
		Status:               constants.StatusActive,
	}
}

func PbToCombinationSku(req *pb.CreateCombinationSkuRequest) *CombinationSku {
	var attIds = map[string][]string{}
	var meteringAttIds = map[string][]string{}
	spuId := req.GetCombinationSpuId().GetValue()
	for _, skuReq := range req.GetCreateSkuRequests() {
		v, ok := attIds[spuId]
		if ok {
			attIds[spuId] = append(v, pbutil.FromProtoStringSlice(skuReq.GetAttributeIds())...)
		} else {
			attIds[spuId] = pbutil.FromProtoStringSlice(skuReq.GetAttributeIds())
		}

		v, ok = meteringAttIds[spuId]
		if ok {
			meteringAttIds[spuId] = append(v, pbutil.FromProtoStringSlice(skuReq.GetMeteringAttributeIds())...)
		} else {
			meteringAttIds[spuId] = pbutil.FromProtoStringSlice(skuReq.MeteringAttributeIds)
		}
	}
	return NewCombinationSku(spuId, attIds, meteringAttIds)
}

type CombinationPrice struct {
	CombinationPriceId string
	CombinationSkuId   string
	SpuId              string
	AttributeId        string
	Prices             map[string]float64 //StepPrice: {upto: price, ..}
	Currency           string
	StartTime          time.Time
	EndTime            time.Time
	CreateTime         time.Time
	StatusTime         time.Time
	Status             string
}

func PbToCombinationPrice(req *pb.CreateCombinationPriceRequest) *CombinationPrice {
	return &CombinationPrice{
		CombinationPriceId: NewCombinationPriceId(),
		CombinationSkuId:   req.GetCombinationSkuId().GetValue(),
		SpuId:              req.GetSpuId().GetValue(),
		AttributeId:        req.GetAttributeId().GetValue(),
		Prices:             req.GetPrices(),
		Currency:           req.GetCurrency().String(),
		StartTime:          pbutil.FromProtoTimestamp(req.GetStartTime()),
		EndTime:            pbutil.FromProtoTimestamp(req.GetEndTime()),
	}
}

type ProbationSku struct {
	ProbationSkuId       string
	SpuId                string
	AttributeIds         []string
	MeteringAttributeIds []string //the meaning of attribute value is the limit of attribute
	LimitNum             uint32
	CreateTime           time.Time
	StatusTime           time.Time
	Status               string
}

func PbToProbationSku(req *pb.CreateProbationSkuRequest) *ProbationSku {
	return &ProbationSku{
		ProbationSkuId:       NewProbationSkuId(),
		SpuId:                req.GetSpuId().GetValue(),
		AttributeIds:         pbutil.FromProtoStringSlice(req.GetAttributeIds()),
		MeteringAttributeIds: pbutil.FromProtoStringSlice(req.GetMeteringAttributeIds()),
		LimitNum:             req.GetLimitNum().GetValue(),
	}
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
