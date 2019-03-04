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
	CreateTime      time.Time
	StatusTime      time.Time
	Status          string
	Description     string
}

func NewAttributeName(name, description string) *AttributeName {
	now := time.Now()
	return &AttributeName{
		AttributeNameId: NewAttributeNameId(),
		Name:            name,
		Description:     description,
		Status:          constants.StatusActive,
		CreateTime:      now,
		StatusTime:      now,
	}
}

func PbToAttributeName(pbAttName *pb.CreateAttributeNameRequest) *AttributeName {
	return NewAttributeName(
		pbAttName.GetName().GetValue(),
		pbAttName.GetDescription().GetValue(),
	)
}

func AttributeNameToPb(attName *AttributeName) *pb.AttributeName {
	return &pb.AttributeName{
		AttributeNameId: pbutil.ToProtoString(attName.AttributeNameId),
		Name:            pbutil.ToProtoString(attName.Name),
		Description:     pbutil.ToProtoString(attName.Description),
	}
}

type AttributeUnit struct {
	AttributeUnitId string
	Name            string
	CreateTime      time.Time
	StatusTime      time.Time
	Status          string
}

func NewAttributeUnit(name string) *AttributeUnit {
	now := time.Now()
	return &AttributeUnit{
		AttributeUnitId: NewAttributeUnitId(),
		Name:            name,
		CreateTime:      now,
		StatusTime:      now,
		Status:          constants.StatusActive,
	}
}

func PbToAttributeUnit(pbAttUnit *pb.CreateAttributeUnitRequest) *AttributeUnit {
	return NewAttributeUnit(
		pbAttUnit.GetName().GetValue(),
	)
}

func AttributeUnitToPb(attUnit *AttributeUnit) *pb.AttributeUnit {
	return &pb.AttributeUnit{
		AttributeUnitId: pbutil.ToProtoString(attUnit.AttributeUnitId),
		Name:            pbutil.ToProtoString(attUnit.Name),
	}
}

type Attribute struct {
	AttributeId     string
	AttributeNameId string
	AttributeUnitId string
	Value           string
	CreateTime      time.Time
	StatusTime      time.Time
	Status          string
}

var AttributeColumns = db.GetColumnsFromStruct(&Attribute{})

func NewAttribute(attNameId, attUnitId, value string) *Attribute {
	now := time.Now()
	return &Attribute{
		AttributeId:     NewAttributeId(),
		AttributeNameId: attNameId,
		AttributeUnitId: attUnitId,
		Value:           value,
		CreateTime:      now,
		StatusTime:      now,
		Status:          constants.StatusActive,
	}
}

func PbToAttribute(pbAttribute *pb.CreateAttributeRequest) *Attribute {
	return NewAttribute(
		pbAttribute.GetAttributeNameId().GetValue(),
		pbAttribute.GetAttributeUnitId().GetValue(),
		pbAttribute.GetValue().GetValue(),
	)
}

func AttributeToPb(att *Attribute) *pb.Attribute {
	return &pb.Attribute{
		AttributeId:     pbutil.ToProtoString(att.AttributeId),
		AttributeNameId: pbutil.ToProtoString(att.AttributeNameId),
		AttributeUnitId: pbutil.ToProtoString(att.AttributeUnitId),
		Value:           pbutil.ToProtoString(att.Value),
	}
}

//SPU: standard product unit
type Spu struct {
	SpuId                    string
	ResourceVersionId        string
	AttributeNameIds         []string
	MeteringAttributeNameIds []string
	CreateTime               time.Time
	StatusTime               time.Time
	Status                   string
}

func NewSpu(resourceVersionId string, attNameIds, meteringAttNameIds []string) *Spu {
	now := time.Now()
	return &Spu{
		SpuId:                    NewSpuId(),
		ResourceVersionId:        resourceVersionId,
		AttributeNameIds:         attNameIds,
		MeteringAttributeNameIds: meteringAttNameIds,
		CreateTime:               now,
		StatusTime:               now,
		Status:                   constants.StatusActive,
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
	StatusTime           time.Time
	Status               string
}

func NewSku(spuId string, attributeIds, meteringAttIds []string) *Sku {
	now := time.Now()
	return &Sku{
		SkuId:                NewSkuId(),
		SpuId:                spuId,
		AttributeIds:         attributeIds,
		MeteringAttributeIds: meteringAttIds,
		CreateTime:           now,
		StatusTime:           now,
		Status:               constants.StatusActive,
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
	Currency    string
	StartTime   time.Time
	EndTime     time.Time
	CreateTime  time.Time
	StatusTime  time.Time
	Status      string
}

func NewPrice(skuId, attId, currency string, prices map[string]float64, startTime, endTime time.Time) *Price {
	now := time.Now()
	if (time.Time{}) == startTime {
		startTime = now
	}
	return &Price{
		PriceId:     NewPriceId(),
		SkuId:       skuId,
		AttributeId: attId,
		Prices:      prices,
		Currency:    currency,
		StartTime:   startTime,
		EndTime:     endTime,
		CreateTime:  now,
		StatusTime:  now,
		Status:      constants.StatusActive,
	}
}

func PbToPrice(pbPrice *pb.CreatePriceRequest) *Price {
	return NewPrice(
		pbPrice.GetSkuId().GetValue(),
		pbPrice.GetAttributeId().GetValue(),
		pbPrice.GetCurrency().String(),
		pbPrice.GetPrices(),
		pbutil.FromProtoTimestamp(pbPrice.GetStartTime()),
		pbutil.FromProtoTimestamp(pbPrice.GetEndTime()),
	)
}
