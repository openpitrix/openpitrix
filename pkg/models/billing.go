// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"time"

	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/util/idutil"
)

func NewPriceId() string {
	return idutil.GetUuid("price-")
}

func NewDiscountId() string {
	return idutil.GetUuid("discount-")
}

func NewCouponId() string {
	return idutil.GetUuid("coupon-")
}

type Price struct {
	Id                string
	ResourceVersionId string
	ChargeMode        string
	Price             int32 //单位是分, 即：currency * 100
	Currency          string
	Duration          int32
	Count             int32
	Rule              string
	FreeTime          int32 //时长单位是小时
}

var PriceColumns = db.GetColumnsFromStruct(&Price{})

type Discount struct {
	Id           string
	Name         string
	PriceId      string
	NewPrice     int32 //和price一样
	Discount     float32
	DiscountType string
	StartTime    time.Time
	EndTime      time.Time
	Mark         string
	UserId       string
}

var DiscountColumns = db.GetColumnsFromStruct(&Discount{})

type Coupon struct {
	Id                string
	Name              string
	Sn                string
	Quota             float64
	Balance           float64
	CouponType        string
	ResourceVersionId string
	Region            string
	Status            string
	StartTime         time.Time
	EndTime           time.Time
	CreateTime        time.Time
	Mark              string
}

var CouponColumns = db.GetColumnsFromStruct(&Coupon{})
