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
	Id               string
	ResouceId        string
	ResouceVersionId string
	UserId           string
	PriceId          string
	CreateTime       time.Time
	LeaseTime        time.Time // action_time
	RenewalTime      time.Time // next update time
	UpdateTime       time.Time
	Status           string //updating / updated / overtime
	Duration         int32
	GroupId          string
}

var LeasingColumns = db.GetColumnsFromStruct(&Leasing{})

func NewLeasing(resourceId, resourceVersionId, userId, priceId, groupId string,
	leaseTime, renewalTime time.Time) *Leasing {
	return &Leasing{
		Id:               NewMeteringId(),
		ResouceId:        resourceId,
		ResouceVersionId: resourceVersionId,
		UserId:           userId,
		CreateTime:       time.Now(),
		LeaseTime:        leaseTime,
		RenewalTime:      renewalTime,
		UpdateTime:       leaseTime,
		Status:           constants.StatusUpdating,
		Duration:         0,
		PriceId:          priceId,
		GroupId:          groupId,
	}
}

func MeteringReq2Leasings(req *pb.MeteringRequest, groupId string) []*Leasing {
	var leasings []*Leasing
	for _, recVersion := range req.ActionResourceList {
		//the renewalTime need to be calculated
		leasing := NewLeasing(req.GetResourceId().GetValue(),
			recVersion.GetResourceVersionId().GetValue(),
			req.GetUserId().GetValue(),
			recVersion.GetPriceId().GetValue(),
			groupId,
			pbutil.FromProtoTimestamp(recVersion.GetActionTime()),
			pbutil.FromProtoTimestamp(recVersion.GetActionTime()),
		)
		leasings = append(leasings, leasing)
	}
	return leasings

}
