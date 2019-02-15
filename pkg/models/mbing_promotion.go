// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"time"

	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/util/idutil"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
)

func NewCRAId() string {
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

type CombinationResourceAttribute struct {
	CRAId                string
	ResourceAttributeIds []string //the id slice of ResourceAttribute
	CreateTime           time.Time
	UpdateTime           time.Time
	Status               string
}

func PbToCRA(req *pb.CreateCRARequest) *CombinationResourceAttribute {
	var resAttIds []string
	for _, resAtt := range req.GetResourceAttributes() {
		resAttIds = append(resAttIds, resAtt.GetResourceAttributeId().GetValue())
	}
	return &CombinationResourceAttribute{
		CRAId:                NewCRAId(),
		ResourceAttributeIds: resAttIds,
	}
}

type CombinationSku struct {
	ComSkuId        string
	CRAId           string
	AttributeValues map[string]string //{resourceVersionId: valueId, ..}
	CreateTime      time.Time
	UpdateTime      time.Time
	Status          string
}

func PbToComSku(req *pb.CreateComSkuRequest) *CombinationSku {
	return &CombinationSku{
		ComSkuId:        NewCombinationSkuId(),
		CRAId:           req.GetCraId().GetValue(),
		AttributeValues: req.GetAttributeValues(),
	}
}

type CombinationPrice struct {
	ComPriceId        string
	ComSkuId          string
	ResourceVersionId string
	AttributeId       string
	Prices            map[string]float64 //StepPrice: {valueId: price, ..}
	Currency          string
	CreateTime        time.Time
	UpdateTime        time.Time
	Status            string
}

func PbToComPrice(req *pb.CreateComPriceRequest) *CombinationPrice {
	return &CombinationPrice{
		ComPriceId:        NewCombinationPriceId(),
		ComSkuId:          req.GetComSkuId().GetValue(),
		ResourceVersionId: req.GetResourceVersionId().GetValue(),
		AttributeId:       req.GetAttributeId().GetValue(),
		Prices:            req.GetPrices(),
		Currency:          req.GetCurrency().String(),
	}
}

type ProbationSku struct {
	ProSkuId            string
	ResourceAttributeId string
	AttributeValues     []string
	LimitNum            int32
	CreateTime          time.Time
	UpdateTime          time.Time
	Status              string
}

func PbToProSku(req *pb.CreateProSkuRequest) *ProbationSku {
	return &ProbationSku{
		ProSkuId:            NewProbationSkuId(),
		ResourceAttributeId: req.GetResourceAttributeId().GetValue(),
		AttributeValues:     pbutil.FromProtoStringSlice(req.GetAttributeValueIds()),
		LimitNum:            req.GetLimitNum().GetValue(),
	}
}

type ProbationRecord struct {
	ProbationSkuId string
	UserId         string
	LimitNum       int8
	CreateTime     time.Time
	ProbationTimes []time.Time
}

type Discount struct {
	Id              string
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
	Id         string
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
	Id         string
	CouponId   string
	UserId     string
	Balance    float64
	Status     string
	CreateTime time.Time
}
