// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"time"

	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/util/idutil"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
)

func NewAttributeId() string {
	return idutil.GetUuid("attribute-")
}

func NewAttributeUnitId() string {
	return idutil.GetUuid("unit-")
}

func NewAttributeValueId() string {
	return idutil.GetUuid("attValue-")
}

func NewResourceAttributeId() string {
	return idutil.GetUuid("resAtt-")
}

func NewSkuId() string {
	return idutil.GetUuid("sku-")
}

func NewPriceId() string {
	return idutil.GetUuid("price-")
}

type Attribute struct {
	AttributeId string
	Name        string
	DisplayName string
	CreateTime  time.Time
	UpdateTime  time.Time
	status      int32
	Remark      string
}

var AttributeColumns = db.GetColumnsFromStruct(&Attribute{})

func PbToAttribute(pbAtt *pb.CreateAttributeRequest) *Attribute {
	return &Attribute{
		AttributeId: NewAttributeId(),
		Name:        pbAtt.GetName().GetValue(),
		DisplayName: pbAtt.GetDisplayName().GetValue(),
		Remark:      pbAtt.GetRemark().GetValue(),
	}
}

type AttributeUnit struct {
	AttributeUnitId string
	Name            string
	DisplayName     string
	CreateTime      time.Time
	UpdateTime      time.Time
	status          int32
}

var AttributeUnitColumns = db.GetColumnsFromStruct(&AttributeUnit{})

func PbToAttUnit(pbAttUnit *pb.CreateAttUnitRequest) *AttributeUnit {
	return &AttributeUnit{
		AttributeUnitId: NewAttributeUnitId(),
		Name:            pbAttUnit.GetName().GetValue(),
		DisplayName:     pbAttUnit.GetDisplayName().GetValue(),
	}
}

type AttributeValue struct {
	AttributeValueId string
	AttributeId      string
	AttributeUnitId  string
	MinValue         int32
	MaxValue         int32
	CreateTime       time.Time
	UpdateTime       time.Time
	status           int32
}

var AttributeValueColumns = db.GetColumnsFromStruct(&AttributeValue{})

func PbToAttValue(pbAttValue *pb.CreateAttValueRequest) *AttributeValue {
	return &AttributeValue{
		AttributeValueId: NewAttributeValueId(),
		AttributeId:      pbAttValue.GetAttributeId().GetValue(),
		AttributeUnitId:  pbAttValue.GetAttributeUnitId().GetValue(),
		MinValue:         pbAttValue.GetMinValue().GetValue(),
		MaxValue:         pbAttValue.GetMaxValue().GetValue(),
	}
}

type ResourceAttribute struct {
	ResourceAttributeId string
	ResourceVersionId   string
	Attributes          []string
	MeteringAttributes  []string
	CreateTime          time.Time
	UpdateTime          time.Time
	status              int32
}

var ResourceAttributeColumns = db.GetColumnsFromStruct(&ResourceAttribute{})

func PbToResAtt(pbResAtt *pb.CreateResAttRequest) *ResourceAttribute {
	return &ResourceAttribute{
		ResourceAttributeId: NewResourceAttributeId(),
		ResourceVersionId:   pbResAtt.GetResourceVersionId().GetValue(),
		Attributes:          pbutil.FromProtoStringSlice(pbResAtt.AttributeIds),
		MeteringAttributes:  pbutil.FromProtoStringSlice(pbResAtt.MeteringAttributeIds),
	}
}

type Sku struct {
	SkuId               string
	ResourceAttributeId string
	Values              []string
	CreateTime          time.Time
	UpdateTime          time.Time
	status              int32
}

func PbToSku(pbSku *pb.CreateSkuRequest) *Sku {
	return &Sku{
		SkuId:               NewSkuId(),
		ResourceAttributeId: pbSku.GetResourceAttributeId().GetValue(),
		Values:              pbutil.FromProtoStringSlice(pbSku.GetAttributeValueIds()),
	}
}

type Price struct {
	PriceId     string
	SkuId       string
	AttributeId string
	Prices      map[string]float64
	currency    string
	CreateTime  time.Time
	UpdateTime  time.Time
	status      int32
}

func PbToPrice(pbPrice *pb.CreatePriceRequest) *Price {
	return &Price{
		PriceId:     NewPriceId(),
		SkuId:       pbPrice.GetSkuId().GetValue(),
		AttributeId: pbPrice.GetAttributeId().GetValue(),
		Prices:      pbPrice.GetPrices(),
		currency:    pbPrice.GetCurrency().String(),
	}
}
