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

func NewCombinationId() string {
	return idutil.GetUuid("com-")
}

func NewCombinationSkuId() string {
	return idutil.GetUuid("comSku-")
}

func NewCombinationMABindingId() string {
	return idutil.GetUuid("comBin-")
}

type Combination struct {
	CombinationId string
	Name          string
	Description   string
	Owner         string
	Status        string
	StartTime     time.Time
	EndTime       time.Time
	CreateTime    time.Time
	StatusTime    time.Time
}

func NewCombination(name, description, owner string, startTime, endTime time.Time) *Combination {
	now := time.Now()
	if (time.Time{}) == startTime {
		startTime = now
	}
	return &Combination{
		CombinationId: NewCombinationId(),
		Name:          name,
		Description:   description,
		Owner:         owner,
		Status:        constants.StatusActive,
		StartTime:     startTime,
		EndTime:       endTime,
		CreateTime:    now,
		StatusTime:    now,
	}
}

func PbToCombination(req *pb.CreateCombinationRequest, owner string) *Combination {
	return NewCombination(
		req.GetName().GetValue(),
		req.GetDescription().GetValue(),
		owner,
		pbutil.FromProtoTimestamp(req.GetStartTime()),
		pbutil.FromProtoTimestamp(req.GetEndTime()),
	)
}

type CombinationSku struct {
	CombinationSkuId string
	CombinationId    string
	SkuId            string
	Status           string
	CreateTime       time.Time
	StatusTime       time.Time
}

func NewCombinationSku(comId, skuId string) *CombinationSku {
	now := time.Now()
	return &CombinationSku{
		CombinationSkuId: NewCombinationSkuId(),
		CombinationId:    comId,
		SkuId:            skuId,
		Status:           constants.StatusActive,
		CreateTime:       now,
		StatusTime:       now,
	}
}

func PbToCombinationSku(req *pb.CreateCombinationSkuRequest) *CombinationSku {
	return NewCombinationSku(
		req.GetCombinationId().GetValue(),
		req.GetSkuId().GetValue(),
	)
}
