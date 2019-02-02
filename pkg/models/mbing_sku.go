// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"time"
)

type Attribute struct {
	Id          string
	Name        string
	DisplayName string
	CreateTime  time.Time
	UpdateTime  time.Time
	status      int32
	Remark      string
}

type AttributeUnit struct {
	Id         string
	Name       string
	CreateTime time.Time
	UpdateTime time.Time
	status     int32
}

type AttributeValue struct {
	Id              string
	AttributeId     string
	AttributeUnitId string
	MinValue        int32
	MaxValue        int32
	CreateTime      time.Time
	UpdateTime      time.Time
	status          int32
}

type ResourceAttribute struct {
	Id                 string
	ResourceVersionId  string
	Attributes         []string
	MeteringAttributes []string
	BillingAttributes  []string
	CreateTime         time.Time
	UpdateTime         time.Time
	status             int32
}

type Sku struct {
	Id                  string
	ResourceAttributeId string
	Values              []string
	CreateTime          time.Time
	UpdateTime          time.Time
	status              int32
}

type Price struct {
	Id                 string
	SkuId              string
	BillingAttributeId string
	Prices             map[string]float32
	currency           string
	CreateTime         time.Time
	UpdateTime         time.Time
	status             int32
}
