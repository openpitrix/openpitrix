// Copyright 2019 The OpenPitrix Authors. All rights reserved.
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
	PriceId     string
	SkuId       string
	AttributeId string
	Prices      map[int64]float64
	Currency    string
	Status      string
	StartTime   time.Time
	EndTime     time.Time
	CreateTime  time.Time
	StatusTime  time.Time
}

func NewPrice(skuId, attributeId, currency string, prices map[int64]float64, startTime, endTime time.Time) *Price {
	now := time.Now()
	if (time.Time{}) == startTime {
		startTime = now
	}
	return &Price{
		PriceId:     NewPriceId(),
		SkuId:       skuId,
		AttributeId: attributeId,
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
		pbPrice.GetSkuId().GetValue(),
		pbPrice.GetAttributeId().GetValue(),
		pbPrice.GetCurrency().String(),
		pbPrice.GetPrices(),
		pbutil.FromProtoTimestamp(pbPrice.GetStartTime()),
		pbutil.FromProtoTimestamp(pbPrice.GetEndTime()),
	)
}

func PriceToPb(price *Price) *pb.Price {
	return &pb.Price{
		PriceId:     pbutil.ToProtoString(price.PriceId),
		SkuId:       pbutil.ToProtoString(price.SkuId),
		AttributeId: pbutil.ToProtoString(price.AttributeId),
		Prices:      price.Prices,
		Currency:    pb.Currency(pb.Currency_value[price.Currency]),
		Status:      pbutil.ToProtoString(price.Status),
		StartTime:   pbutil.ToProtoTimestamp(price.StartTime),
		EndTime:     pbutil.ToProtoTimestamp(price.EndTime),
		CreateTime:  pbutil.ToProtoTimestamp(price.CreateTime),
		StatusTime:  pbutil.ToProtoTimestamp(price.StatusTime),
	}
}

type LeasingContract struct {
	ContractId       string
	LeasingId        string
	ResourceId       string
	SkuId            string
	UserId           string
	MeteringValues   map[string]float64
	FeeInfo          string
	Fee              float64
	DueFee           float64
	OtherContractFee float64
	CouponFee        float64
	RealFee          float64
	Currency         string
	Status           string //active/updating/deleted
	StartTime        time.Time
	StatusTime       time.Time
	CreateTime       time.Time
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
	ContractId       string
	LeasingId        string
	ResourceId       string
	SkuId            string
	UserId           string
	MeteringValues   map[string]float64
	FeeInfo          string
	Fee              float64
	DueFee           float64
	OtherContractFee float64
	CouponFee        float64
	RealFee          float64
	Currency         string
	StartTime        time.Time
	EndTime          time.Time
	CreateTime       time.Time
}

func (leasingContract LeasingContract) ToLeasedContract() *LeasedContract {
	return &LeasedContract{
		ContractId:       leasingContract.ContractId,
		LeasingId:        leasingContract.LeasingId,
		ResourceId:       leasingContract.ResourceId,
		SkuId:            leasingContract.SkuId,
		UserId:           leasingContract.UserId,
		MeteringValues:   leasingContract.MeteringValues,
		FeeInfo:          leasingContract.FeeInfo,
		Fee:              leasingContract.Fee,
		DueFee:           leasingContract.DueFee,
		OtherContractFee: leasingContract.OtherContractFee,
		CouponFee:        leasingContract.CouponFee,
		RealFee:          leasingContract.RealFee,
		Currency:         leasingContract.Currency,
		StartTime:        leasingContract.StartTime,
		EndTime:          leasingContract.StatusTime,
		CreateTime:       time.Now(),
	}
}

type Account struct {
	UserId     string
	UserType   string
	Balance    float64
	Currency   string
	Income     map[string]float64
	Status     string
	CreateTime time.Time
	StatusTime time.Time
}
