// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"time"

	"openpitrix.io/openpitrix/pkg/util/idutil"
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
	Id                 string
	ResourceVersionIds []string
	Attributes         map[string]string //{resourceVersionId: attributeId, ..}
	MeteringAttributes map[string]string //{resourceVersionId: attributeId, ..}
	CreateTime         time.Time
	UpdateTime         time.Time
	Status             string
}

type CombinationSku struct {
	Id              string
	CRAId           string
	AttributeValues map[string]string //{resourceVersionId: valueId, ..}
	CreateTime      time.Time
	UpdateTime      time.Time
	Status          string
}

type CombinationPrice struct {
	Id                string
	CombinationSkuId  string
	ResourceVersionId string
	AttributeId       string
	Prices            map[string]float64 //StepPrice: {valueId: price, ..}
	Currency          string
	CreateTime        time.Time
	UpdateTime        time.Time
	Status            string
}

type ProbationSku struct {
	Id                  string
	ResourceAttributeId string
	AttributeValues     []string
	LimitNum            int8
	CreateTime          time.Time
	UpdateTime          time.Time
	Status              string
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
