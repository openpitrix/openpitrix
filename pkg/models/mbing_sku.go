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
	"openpitrix.io/openpitrix/pkg/util/pbutil"
)

func NewAttributeNameId() string {
	return idutil.GetUuid("att-name-")
}

func NewAttributeUnitId() string {
	return idutil.GetUuid("att-unit-")
}

func NewAttributeId() string {
	return idutil.GetUuid("att-")
}

func NewResourceAttributeId() string {
	return idutil.GetUuid("res-att-")
}

func NewSkuId() string {
	return idutil.GetUuid("sku-")
}

func NewPriceId() string {
	return idutil.GetUuid("price-")
}

type AttributeName struct {
	AttributeNameId string
	Name            string
	DisplayName     string
	CreateTime      time.Time
	UpdateTime      time.Time
	Status          string
	Remark          string
}

var AttributeNameColumns = db.GetColumnsFromStruct(&AttributeName{})

func NewAttributeName(name, displayName, remark string) *AttributeName {
	now := time.Now()
	return &AttributeName{
		AttributeNameId: NewAttributeNameId(),
		Name:            name,
		DisplayName:     displayName,
		Remark:          remark,
		Status:          constants.StatusInUse2,
		CreateTime:      now,
		UpdateTime:      now,
	}
}

func PbToAttributeName(pbAttName *pb.CreateAttributeNameRequest) *AttributeName {
	return NewAttributeName(
		pbAttName.GetName().GetValue(),
		pbAttName.GetDisplayName().GetValue(),
		pbAttName.GetRemark().GetValue(),
	)
}

func AttributeNameToPb(attName *AttributeName) *pb.AttributeName {
	return &pb.AttributeName{
		AttributeNameId: pbutil.ToProtoString(attName.AttributeNameId),
		Name:            pbutil.ToProtoString(attName.Name),
		DisplayName:     pbutil.ToProtoString(attName.DisplayName),
		Remark:          pbutil.ToProtoString(attName.Remark),
	}
}

type AttributeUnit struct {
	AttributeUnitId string
	Name            string
	DisplayName     string
	CreateTime      time.Time
	UpdateTime      time.Time
	Status          string
}

var AttributeUnitColumns = db.GetColumnsFromStruct(&AttributeUnit{})

func NewAttributeUnit(name, display string) *AttributeUnit {
	now := time.Now()
	return &AttributeUnit{
		AttributeUnitId: NewAttUnitId(),
		Name:            name,
		DisplayName:     display,
		CreateTime:      now,
		UpdateTime:      now,
		Status:          constants.StatusInUse2,
	}
}

func PbToAttUnit(pbAttUnit *pb.CreateAttUnitRequest) *AttributeUnit {
	return NewAttributeUnit(
		pbAttUnit.GetName().GetValue(),
		pbAttUnit.GetDisplayName().GetValue(),
	)
}

func AttributeUnitToPb(attUnit *AttributeUnit) *pb.AttributeUnit {
	return &pb.AttributeUnit{
		AttributeUnitId: pbutil.ToProtoString(attUnit.AttributeUnitId),
		Name:            pbutil.ToProtoString(attUnit.Name),
		DisplayName:     pbutil.ToProtoString(attUnit.DisplayName),
	}
}

type AttributeValue struct {
	AttributeValueId string
	AttributeId      string
	AttributeUnitId  string
	MinValue         uint32
	MaxValue         uint32
	CreateTime       time.Time
	UpdateTime       time.Time
	Status           string
}

func NewAttributeValue(attId, attUnitId string, minValue, maxValue uint32) *AttributeValue {
	now := time.Now()
	return &AttributeValue{
		AttributeValueId: NewAttValueId(),
		AttributeId:      attId,
		AttributeUnitId:  attUnitId,
		MinValue:         minValue,
		MaxValue:         maxValue,
		CreateTime:       now,
		UpdateTime:       now,
		Status:           constants.StatusInUse2,
	}
}

func PbToAttValue(pbAttValue *pb.CreateAttValueRequest) *AttributeValue {
	return NewAttributeValue(
		pbAttValue.GetAttributeId().GetValue(),
		pbAttValue.GetAttributeUnitId().GetValue(),
		pbAttValue.GetMinValue().GetValue(),
		pbAttValue.GetMaxValue().GetValue(),
	)
}

type ResourceAttribute struct {
	ResourceAttributeId string
	ResourceVersionId   string
	Attributes          []string
	MeteringAttributes  []string
	CreateTime          time.Time
	UpdateTime          time.Time
	Status              string
}

func NewResourceAttribute(resVerId string, atts, metAtts []string) *ResourceAttribute {
	now := time.Now()
	return &ResourceAttribute{
		ResourceAttributeId: NewResAttId(),
		ResourceVersionId:   resVerId,
		Attributes:          atts,
		MeteringAttributes:  metAtts,
		CreateTime:          now,
		UpdateTime:          now,
		Status:              constants.StatusInUse2,
	}
}

func PbToResAtt(pbResAtt *pb.CreateResAttRequest) *ResourceAttribute {
	return NewResourceAttribute(
		pbResAtt.GetResourceVersionId().GetValue(),
		pbutil.FromProtoStringSlice(pbResAtt.AttributeIds),
		pbutil.FromProtoStringSlice(pbResAtt.MeteringAttributeIds),
	)
}

type Sku struct {
	SkuId               string
	ResourceAttributeId string
	Values              []string
	CreateTime          time.Time
	UpdateTime          time.Time
	Status              string
}

var SkuColumns = db.GetColumnsFromStruct(&ResourceAttribute{})

func NewSku(resAttId string, values []string) *Sku {
	now := time.Now()
	return &Sku{
		SkuId:               NewSkuId(),
		ResourceAttributeId: resAttId,
		Values:              values,
		CreateTime:          now,
		UpdateTime:          now,
		Status:              constants.StatusInUse2,
	}
}

func PbToSku(pbSku *pb.CreateSkuRequest) *Sku {
	return NewSku(
		pbSku.GetResourceAttributeId().GetValue(),
		pbutil.FromProtoStringSlice(pbSku.GetAttributeValueIds()),
	)
}

type Price struct {
	PriceId     string
	SkuId       string
	AttributeId string
	Prices      map[string]float64
	currency    string
	CreateTime  time.Time
	UpdateTime  time.Time
	Status      string
}

func NewPrice(skuId, attId, currency string, prices map[string]float64) *Price {
	now := time.Now()
	return &Price{
		PriceId:     NewPriceId(),
		SkuId:       skuId,
		AttributeId: attId,
		Prices:      prices,
		currency:    currency,
		CreateTime:  now,
		UpdateTime:  now,
		Status:      constants.StatusInUse2,
	}
}

func PbToPrice(pbPrice *pb.CreatePriceRequest) *Price {
	return NewPrice(
		pbPrice.GetSkuId().GetValue(),
		pbPrice.GetAttributeId().GetValue(),
		pbPrice.GetCurrency().String(),
		pbPrice.GetPrices(),
	)
}
