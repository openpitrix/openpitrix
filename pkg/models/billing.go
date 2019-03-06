// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"time"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/util/idutil"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
	"openpitrix.io/openpitrix/pkg/pb"
)

func NewPriceId() string {
	return idutil.GetUuid("price-")
}

type Price struct {
	PriceId     string
	BindingId   string
	Prices      map[int64]float64
	Currency    string
	Status      string
	StartTime   time.Time
	EndTime     time.Time
	CreateTime  time.Time
	StatusTime  time.Time
}

func NewPrice(bindingId, currency string, prices map[int64]float64, startTime, endTime time.Time) *Price {
	now := time.Now()
	if (time.Time{}) == startTime {
		startTime = now
	}
	return &Price{
		PriceId:     NewPriceId(),
		BindingId:   bindingId,
		Prices:      prices,
		Currency:    currency,
		Status:      constants.StatusActive,
		StartTime:   startTime,
		EndTime:     endTime,
		CreateTime:  now,
		StatusTime:  now,
	}
}

func PbToPrice(pbPrice *pb.CreatePriceRequest) *Price {
	return NewPrice(
		pbPrice.GetBindingId().GetValue(),
		pbPrice.GetCurrency().String(),
		pbPrice.GetPrices(),
		pbutil.FromProtoTimestamp(pbPrice.GetStartTime()),
		pbutil.FromProtoTimestamp(pbPrice.GetEndTime()),
	)
}

func PriceToPb(price *Price) *pb.Price {
	return &pb.Price{
		PriceId: pbutil.ToProtoString(price.PriceId),
		BindingId: pbutil.ToProtoString(price.BindingId),
		Prices: price.Prices,
		Currency: pb.Currency(pb.Currency_value[price.Currency]),
		Status: pbutil.ToProtoString(price.Status),
		StartTime: pbutil.ToProtoTimestamp(price.StartTime),
		EndTime: pbutil.ToProtoTimestamp(price.EndTime),
		CreateTime: pbutil.ToProtoTimestamp(price.CreateTime),
		StatusTime: pbutil.ToProtoTimestamp(price.StatusTime),
	}
}




type LeasingContract struct {
	Id             string
	LeasingId      string
	SkuId          string
	UserId         string
	MeteringValues map[string]interface{}
	StartTime      time.Time
	StatusTime     time.Time
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
