// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"time"

	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/util/pbutil"

	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/util/idutil"
)

func NewLeasingId() string {
	return idutil.GetUuid("leasing-")
}

type Leasing struct {
	LeasingId          string
	GroupId            string
	UserId             string
	ResourceId         string
	SkuId              string
	SkuType            string
	OtherInfo          string
	MeteringValues     map[string]float64
	LeaseTime          time.Time  // action_time
	UpdateDurationTime time.Time  // auto current update time
	RenewalTime        *time.Time // next update time
	StatusTime         time.Time
	CreateTime         time.Time
	CloseTime          map[time.Time]time.Time //{closeTime: restartTime, ..}
	Status             string
}

var LeasingColumns = db.GetColumnsFromStruct(&Leasing{})

func pbToMeteringValues(pbMetVals []*pb.MeteringAttributeValue) map[string]float64 {
	metertingValues := map[string]float64{}
	for _, pbMetVal := range pbMetVals {
		attributeId := pbMetVal.GetAttributeId().GetValue()
		metertingValues[attributeId] = pbMetVal.GetValue().Value
	}
	return metertingValues
}

func PbToLeasing(req *pb.MeteringRequest, mSku *pb.MeteringSku, groupId string, renewalTime *time.Time) *Leasing {
	actionTime := pbutil.FromProtoTimestamp(mSku.GetActionTime())
	return &Leasing{
		LeasingId:          NewLeasingId(),
		GroupId:            groupId,
		UserId:             req.GetUserId().GetValue(),
		ResourceId:         req.GetResourceId().GetValue(),
		SkuId:              mSku.GetSkuId().GetValue(),
		SkuType:            mSku.GetType().String(),
		OtherInfo:          mSku.GetOtherInfo().GetValue(),
		MeteringValues:     pbToMeteringValues(mSku.GetAttributeValues()),
		LeaseTime:          actionTime,
		UpdateDurationTime: actionTime,
		RenewalTime:        renewalTime,
	}
}

type Leased struct {
	LeasedId       string
	GroupId        string
	UserId         string
	ResourceId     string
	SkuId          string
	OtherInfo      string
	MeteringValues map[string]float64
	LeaseTime      time.Time // action_time
	CreateTime     time.Time
	CloseTime      map[time.Time]time.Time //{closeTime: restartTime, ..}
}

var LeasedColumns = db.GetColumnsFromStruct(&Leased{})

func (l Leasing) toLeased() Leased {
	return Leased{
		LeasedId:       l.LeasingId,
		GroupId:        l.GroupId,
		UserId:         l.UserId,
		ResourceId:     l.ResourceId,
		SkuId:          l.SkuId,
		OtherInfo:      l.OtherInfo,
		MeteringValues: l.MeteringValues,
		LeaseTime:      l.LeaseTime,
		CloseTime:      l.CloseTime,
	}
}
