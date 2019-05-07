// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"time"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/db"

	"openpitrix.io/openpitrix/pkg/util/idutil"
)

func NewLeasingId() string {
	return idutil.GetUuid("leasing-")
}

//SkuMetering
type Leasing struct {
	LeasingId          string
	UserId             string
	ResourceId         string
	SkuId              string
	MeteringValues     map[string]float64 //
	LeaseTime          time.Time          //action_time
	UpdateDurationTime time.Time          //update time for duration
	RenewalTime        time.Time          //next update time
	CreateTime         time.Time
	StatusTime         time.Time               //update time by other services(cluster_manager)
	StopTimes          map[time.Time]time.Time //{closeTime: restartTime, ..}
	Status             string
}

var LeasingColumns = db.GetColumnsFromStruct(&Leasing{})

func NewLeasing(values map[string]float64, userId, resourceId, skuId string, actionTime, renewalTime time.Time) *Leasing {
	return &Leasing{
		LeasingId:          NewLeasingId(),
		UserId:             userId,
		ResourceId:         resourceId,
		SkuId:              skuId,
		MeteringValues:     values,
		LeaseTime:          actionTime,
		UpdateDurationTime: actionTime,
		RenewalTime:        renewalTime,
		Status:             constants.StatusActive,
		CreateTime:         actionTime,
		StatusTime:         actionTime,
		StopTimes:          nil,
	}
}

type Leased struct {
	LeasedId       string
	UserId         string
	ResourceId     string
	SkuId          string
	MeteringValues map[string]float64
	LeaseTime      time.Time // action_time
	CreateTime     time.Time
	StopTime       map[time.Time]time.Time //{closeTime: restartTime, ..}
}

func (l *Leasing) ToLeased() *Leased {
	return &Leased{
		LeasedId:       l.LeasingId,
		UserId:         l.UserId,
		ResourceId:     l.ResourceId,
		SkuId:          l.SkuId,
		MeteringValues: l.MeteringValues,
		LeaseTime:      l.LeaseTime,
		StopTime:       l.StopTimes,
		CreateTime:     time.Now(),
	}
}
