// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"time"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/db"

	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/util/idutil"
)

func NewLeasingId() string {
	return idutil.GetUuid("leasing-")
}

//SkuMetering
type Leasing struct {
	LeasingId          string
	GroupId            string
	UserId             string
	ResourceId         string
	SkuId              string
	MeteringValues     map[string]float64
	LeaseTime          time.Time //action_time
	UpdateDurationTime time.Time //update time for duration
	RenewalTime        time.Time //next update time
	CreateTime         time.Time
	StatusTime         time.Time               //update time by other services(cluster_manager)
	StopTime           map[time.Time]time.Time //{closeTime: restartTime, ..}
	Status             string
}

var LeasingColumns = db.GetColumnsFromStruct(&Leasing{})

func pbToMeteringValues(pbMetVals []*pb.MeteringValue) map[string]float64 {
	metertingValues := map[string]float64{}
	for _, pbMetVal := range pbMetVals {
		attributeId := pbMetVal.GetAttributeId().GetValue()
		metertingValues[attributeId] = pbMetVal.GetValue().Value
	}
	return metertingValues
}

func NewLeasing(req *pb.SkuMetering, groupId, userId string, actionTime, renewalTime time.Time) *Leasing {
	return &Leasing{
		LeasingId:          NewLeasingId(),
		GroupId:            groupId,
		UserId:             userId,
		ResourceId:         req.GetResourceId().GetValue(),
		SkuId:              req.GetSkuId().GetValue(),
		MeteringValues:     pbToMeteringValues(req.GetMeteringValues()),
		LeaseTime:          actionTime,
		UpdateDurationTime: actionTime,
		RenewalTime:        renewalTime,
		Status:             constants.StatusActive,
		CreateTime:         actionTime,
		StatusTime:         actionTime,
		StopTime:           nil,
	}
}

type Leased struct {
	LeasedId       string
	GroupId        string
	UserId         string
	ResourceId     string
	SkuId          string
	MeteringValues map[string]float64
	LeaseTime      time.Time // action_time
	CreateTime     time.Time
	StopTime       map[time.Time]time.Time //{closeTime: restartTime, ..}
}

func (l Leasing) toLeased() *Leased {
	return &Leased{
		LeasedId:       l.LeasingId,
		GroupId:        l.GroupId,
		UserId:         l.UserId,
		ResourceId:     l.ResourceId,
		SkuId:          l.SkuId,
		MeteringValues: l.MeteringValues,
		LeaseTime:      l.LeaseTime,
		StopTime:       l.StopTime,
		CreateTime:     time.Now(),
	}
}
