// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"time"

	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/util/pbutil"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/util/idutil"
)

func NewMeteringId() string {
	return idutil.GetUuid("metering-")
}

type Leasing struct {
	Id             string
	GroupId        string
	UserId         string
	ResourceId     string
	SkuId          string
	MeteringValues map[string]interface{}
	LeaseTime      time.Time // action_time
	RenewalTime    time.Time // next update time
	UpdateTime     time.Time
	CreateTime     time.Time
	CloseTime      map[time.Time]interface{} //{closeTime: restartTime, ..}
	Status         string
}

var LeasingColumns = db.GetColumnsFromStruct(&Leasing{})

func NewLeasing(groupId, userId, resourceId, skuId string,
	leaseTime, renewalTime time.Time,
	meteringValues map[string]interface{}) *Leasing {
	return &Leasing{
		Id:             NewMeteringId(),
		GroupId:        groupId,
		UserId:         userId,
		ResourceId:     resourceId,
		SkuId:          skuId,
		LeaseTime:      leaseTime,
		RenewalTime:    renewalTime,
		UpdateTime:     leaseTime,
		MeteringValues: meteringValues,
		Status:         constants.StatusRunning,
	}
}

func MeteringReq2Leasings(req *pb.MeteringRequest, groupId string) []*Leasing {
	var leasings []*Leasing
	for _, sku := range req.SkuSet {
		//the renewalTime need to be calculated
		leasing := NewLeasing(
			groupId,
			req.GetUserId().GetValue(),
			req.GetResourceId().GetValue(),
			sku.GetId().GetValue(),
			pbutil.FromProtoTimestamp(sku.GetActionTime()),
			pbutil.FromProtoTimestamp(sku.GetActionTime()), //need to update to real value
			make(map[string]interface{}),
		)
		leasings = append(leasings, leasing)
	}
	return leasings

}

type Leased struct {
	LeasingId      string
	GroupId        string
	UserId         string
	ResourceId     string
	SkuId          string
	MeteringValues map[string]interface{}
	LeaseTime      time.Time // action_time
	UpdateTime     time.Time
	CreateTime     time.Time
	CloseTime      map[time.Time]interface{} //{closeTime: restartTime, ..}
}

var LeasedColumns = db.GetColumnsFromStruct(&Leased{})

func toLeased(leasing Leasing) Leased {
	return Leased{
		LeasingId:      leasing.Id,
		GroupId:        leasing.GroupId,
		UserId:         leasing.UserId,
		ResourceId:     leasing.ResourceId,
		SkuId:          leasing.SkuId,
		MeteringValues: leasing.MeteringValues,
		LeaseTime:      leasing.LeaseTime,
		UpdateTime:     leasing.UpdateTime,
		CloseTime:      leasing.CloseTime,
	}
}
