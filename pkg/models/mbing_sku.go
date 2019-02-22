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

func NewSpuId() string {
	return idutil.GetUuid("spu-")
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
		AttributeUnitId: NewAttributeUnitId(),
		Name:            name,
		DisplayName:     display,
		CreateTime:      now,
		UpdateTime:      now,
		Status:          constants.StatusInUse2,
	}
}

func PbToAttributeUnit(pbAttUnit *pb.CreateAttributeUnitRequest) *AttributeUnit {
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

type Attribute struct {
	AttributeId     string
	AttributeNameId string
	AttributeUnitId string
	MinValue        uint32
	MaxValue        uint32
	CreateTime      time.Time
	UpdateTime      time.Time
	Status          string
}

func NewAttribute(attNameId, attUnitId string, minValue, maxValue uint32) *Attribute {
	now := time.Now()
	return &Attribute{
		AttributeId:     NewAttributeId(),
		AttributeNameId: attNameId,
		AttributeUnitId: attUnitId,
		MinValue:        minValue,
		MaxValue:        maxValue,
		CreateTime:      now,
		UpdateTime:      now,
		Status:          constants.StatusInUse2,
	}
}

func PbToAttribute(pbAttribute *pb.CreateAttributeRequest) *Attribute {
	return NewAttribute(
		pbAttribute.GetAttributeNameId().GetValue(),
		pbAttribute.GetAttributeUnitId().GetValue(),
		pbAttribute.GetMinValue().GetValue(),
		pbAttribute.GetMaxValue().GetValue(),
	)
}

//SPU: standard product unit
type Spu struct {
	SpuId                    string
	ResourceVersionId        string
	AttributeNameIds         []string
	MeteringAttributeNameIds []string
	CreateTime               time.Time
	UpdateTime               time.Time
	Status                   string
}

var SpuColumns = db.GetColumnsFromStruct(&Spu{})

func NewSpu(resourceVersionId string, attNameIds, meteringAttNameIds []string) *Spu {
	now := time.Now()
	return &Spu{
		SpuId:                    NewSpuId(),
		ResourceVersionId:        resourceVersionId,
		AttributeNameIds:         attNameIds,
		MeteringAttributeNameIds: meteringAttNameIds,
		CreateTime:               now,
		UpdateTime:               now,
		Status:                   constants.StatusInUse2,
	}
}

func PbToSpu(pbSpu *pb.CreateSpuRequest) *Spu {
	return NewSpu(
		pbSpu.GetResourceVersionId().GetValue(),
		pbutil.FromProtoStringSlice(pbSpu.AttributeNameIds),
		pbutil.FromProtoStringSlice(pbSpu.MeteringAttributeNameIds),
	)
}

//SKU: stock keeping unit
type Sku struct {
	SkuId                string
	SpuId                string
	AttributeIds         []string
	MeteringAttributeIds []string
	CreateTime           time.Time
	UpdateTime           time.Time
	Status               string
}

var SkuColumns = db.GetColumnsFromStruct(&Sku{})

func NewSku(spuId string, attributeIds, meteringAttIds []string) *Sku {
	now := time.Now()
	return &Sku{
		SkuId:                NewSkuId(),
		SpuId:                spuId,
		AttributeIds:         attributeIds,
		MeteringAttributeIds: meteringAttIds,
		CreateTime:           now,
		UpdateTime:           now,
		Status:               constants.StatusInUse2,
	}
}

func PbToSku(pbSku *pb.CreateSkuRequest) *Sku {
	return NewSku(
		pbSku.GetSpuId().GetValue(),
		pbutil.FromProtoStringSlice(pbSku.GetAttributeIds()),
		pbutil.FromProtoStringSlice(pbSku.GetMeteringAttributeIds()),
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
