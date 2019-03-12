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

func NewPriceId() string {
	return idutil.GetUuid("price-")
}

func NewContractId() string {
	return idutil.GetUuid("contract-")
}

type Price struct {
	PriceId    string
	BindingId  string
	Prices     map[int64]float64
	Currency   string
	Status     string
	StartTime  time.Time
	EndTime    time.Time
	CreateTime time.Time
	StatusTime time.Time
}

func NewPrice(bindingId, currency string, prices map[int64]float64, startTime, endTime time.Time) *Price {
	now := time.Now()
	if (time.Time{}) == startTime {
		startTime = now
	}
	return &Price{
		PriceId:    NewPriceId(),
		BindingId:  bindingId,
		Prices:     prices,
		Currency:   currency,
		Status:     constants.StatusActive,
		StartTime:  startTime,
		EndTime:    endTime,
		CreateTime: now,
		StatusTime: now,
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
		PriceId:    pbutil.ToProtoString(price.PriceId),
		BindingId:  pbutil.ToProtoString(price.BindingId),
		Prices:     price.Prices,
		Currency:   pb.Currency(pb.Currency_value[price.Currency]),
		Status:     pbutil.ToProtoString(price.Status),
		StartTime:  pbutil.ToProtoTimestamp(price.StartTime),
		EndTime:    pbutil.ToProtoTimestamp(price.EndTime),
		CreateTime: pbutil.ToProtoTimestamp(price.CreateTime),
		StatusTime: pbutil.ToProtoTimestamp(price.StatusTime),
	}
}

type LeasingContract struct {
	ContractId     string
	LeasingId      string
	ResourceId     string
	SkuId          string
	UserId         string
	MeteringValues map[string]float64
	Status         string //active/updating
	StartTime      time.Time
	StatusTime     time.Time
	CreateTime     time.Time
	FeeInfo        string
	Fee            float64
	DueFee         float64
	BeforeBillFee  float64
	CouponFee      float64
	RealFee        float64
	Currency       string
}

func NewLeasingContract(leasingId, resourceId, skuId, userId, currency string,
	meteringValues map[string]float64,
	startTime, updateDurationTime time.Time) *LeasingContract {

	now := time.Now()
	return &LeasingContract{
		ContractId:     NewContractId(),
		LeasingId:      leasingId,
		ResourceId:     resourceId,
		SkuId:          skuId,
		UserId:         userId,
		MeteringValues: meteringValues,
		Status:         constants.StatusActive,
		StartTime:      startTime,
		StatusTime:     now,
		CreateTime:     now,
		Currency:       currency,
	}
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

type Account struct {
	UserId     string
	UserType   string
	Balance    float64
	Currency   string
	Income     map[string]float64
	Status     string
	StartTime  time.Time
	StatusTime time.Time
}
