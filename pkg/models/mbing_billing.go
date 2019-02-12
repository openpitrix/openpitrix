// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"time"
)

type LeasingContract struct {
	Id             string
	LeasingId      string
	SkuId          string
	UserId         string
	MeteringValues map[string]interface{}
	StartTime      time.Time
	UpdateTime     time.Time
	CreateTime     time.Time
	FeeInfo        string
	Fee            float32
	DueFee         float32
	BeforeBillFee  float32
	CouponFee      float32
	RealFee        float32
	currency       string
}

type LeasedContract struct {
	ContractId     string
	LeasingId      string
	SkuId          string
	UserId         string
	MeteringValues map[string]interface{}
	StartTime      time.Time
	EndTime        time.Time
	CreateTime     time.Time
	FeeInfo        string
	Fee            float32
	DueFee         float32
	BeforeBillFee  float32
	CouponFee      float32
	RealFee        float32
	currency       string
}
