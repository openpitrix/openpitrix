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
	return idutil.GetUuid("attn-")
}

func NewAttributeUnitId() string {
	return idutil.GetUuid("attu-")
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

func NewMeteringAttributeBindingId() string {
	return idutil.GetUuid("binding-")
}

type AttributeName struct {
	AttributeNameId string
	Name            string
	Description     string
	Type            string
	Status          string
	CreateTime      time.Time
	StatusTime      time.Time
}

func NewAttributeName(name, description, attType string) *AttributeName {
	now := time.Now()
	return &AttributeName{
		AttributeNameId: NewAttributeNameId(),
		Name:            name,
		Description:     description,
		Type:            attType,
		Status:          constants.StatusActive,
		CreateTime:      now,
		StatusTime:      now,
	}
}

func PbToAttributeName(pbAttName *pb.CreateAttributeNameRequest) *AttributeName {
	return NewAttributeName(
		pbAttName.GetName().GetValue(),
		pbAttName.GetDescription().GetValue(),
		pbAttName.GetType().String(),
	)
}

func AttributeNameToPb(attName *AttributeName) *pb.AttributeName {
	return &pb.AttributeName{
		AttributeNameId: pbutil.ToProtoString(attName.AttributeNameId),
		Name:            pbutil.ToProtoString(attName.Name),
		Description:     pbutil.ToProtoString(attName.Description),
		Type:            pb.AttributeType(pb.AttributeType_value[attName.Type]),
		Status:          pbutil.ToProtoString(attName.Status),
		CreateTime:      pbutil.ToProtoTimestamp(attName.CreateTime),
		StatusTime:      pbutil.ToProtoTimestamp(attName.StatusTime),
	}
}

type AttributeUnit struct {
	AttributeUnitId string
	Name            string
	Status          string
	CreateTime      time.Time
	StatusTime      time.Time
}

func NewAttributeUnit(name string) *AttributeUnit {
	now := time.Now()
	return &AttributeUnit{
		AttributeUnitId: NewAttributeUnitId(),
		Name:            name,
		Status:          constants.StatusActive,
		CreateTime:      now,
		StatusTime:      now,
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
		Status:          pbutil.ToProtoString(attUnit.Status),
		CreateTime:      pbutil.ToProtoTimestamp(attUnit.CreateTime),
		StatusTime:      pbutil.ToProtoTimestamp(attUnit.StatusTime),
	}
}

type Attribute struct {
	AttributeId     string
	AttributeNameId string
	AttributeUnitId string
	Value           string
	Owner           string
	Status          string
	CreateTime      time.Time
	StatusTime      time.Time
}

var AttributeColumns = db.GetColumnsFromStruct(&Attribute{})

func NewAttribute(attNameId, attUnitId, value, owner string) *Attribute {
	now := time.Now()
	return &Attribute{
		AttributeId:     NewAttributeId(),
		AttributeNameId: attNameId,
		AttributeUnitId: attUnitId,
		Value:           value,
		Owner:           owner,
		CreateTime:      now,
		StatusTime:      now,
		Status:          constants.StatusActive,
	}
}

func PbToAttribute(pbAttribute *pb.CreateAttributeRequest, owner string) *Attribute {
	return NewAttribute(
		pbAttribute.GetAttributeNameId().GetValue(),
		pbAttribute.GetAttributeUnitId().GetValue(),
		pbAttribute.GetValue().GetValue(),
		owner,
	)
}

func AttributeToPb(att *Attribute) *pb.Attribute {
	return &pb.Attribute{
		AttributeId:     pbutil.ToProtoString(att.AttributeId),
		AttributeNameId: pbutil.ToProtoString(att.AttributeNameId),
		AttributeUnitId: pbutil.ToProtoString(att.AttributeUnitId),
		Value:           pbutil.ToProtoString(att.Value),
		Owner:           pbutil.ToProtoString(att.Owner),
		Status:          pbutil.ToProtoString(att.Status),
		CreateTime:      pbutil.ToProtoTimestamp(att.CreateTime),
		StatusTime:      pbutil.ToProtoTimestamp(att.StatusTime),
	}
}

//SPU: standard product unit
type Spu struct {
	SpuId      string
	ProductId  string
	Owner      string
	Status     string
	CreateTime time.Time
	StatusTime time.Time
}

func NewSpu(productId, owner string) *Spu {
	now := time.Now()
	return &Spu{
		SpuId:      NewSpuId(),
		ProductId:  productId,
		Owner:      owner,
		Status:     constants.StatusActive,
		CreateTime: now,
		StatusTime: now,
	}
}

func PbToSpu(pbSpu *pb.CreateSpuRequest, owner string) *Spu {
	return NewSpu(pbSpu.GetProductId().GetValue(), owner)
}

func SpuToPb(spu *Spu) *pb.Spu {
	return &pb.Spu{
		SpuId:      pbutil.ToProtoString(spu.SpuId),
		ProductId:  pbutil.ToProtoString(spu.ProductId),
		Owner:      pbutil.ToProtoString(spu.Owner),
		Status:     pbutil.ToProtoString(spu.Status),
		CreateTime: pbutil.ToProtoTimestamp(spu.CreateTime),
		StatusTime: pbutil.ToProtoTimestamp(spu.StatusTime),
	}
}

//SKU: stock keeping unit
type Sku struct {
	SkuId        string
	SpuId        string
	AttributeIds []string
	Status       string
	CreateTime   time.Time
	StatusTime   time.Time
}

func NewSku(spuId string, attributeIds []string) *Sku {
	now := time.Now()
	return &Sku{
		SkuId:        NewSkuId(),
		SpuId:        spuId,
		AttributeIds: attributeIds,
		Status:       constants.StatusActive,
		CreateTime:   now,
		StatusTime:   now,
	}
}

func PbToSku(pbSku *pb.CreateSkuRequest) *Sku {
	return NewSku(
		pbSku.GetSpuId().GetValue(),
		pbSku.GetAttributeIds(),
	)
}

func SkuToPb(sku *Sku) *pb.Sku {
	return &pb.Sku{
		SkuId:        pbutil.ToProtoString(sku.SkuId),
		SpuId:        pbutil.ToProtoString(sku.SpuId),
		AttributeIds: sku.AttributeIds,
		Status:       pbutil.ToProtoString(sku.Status),
		CreateTime:   pbutil.ToProtoTimestamp(sku.CreateTime),
		StatusTime:   pbutil.ToProtoTimestamp(sku.StatusTime),
	}
}

type MeteringAttributeBinding struct {
	BindingId   string
	SkuId       string
	AttributeId string
	Status      string
	CreateTime  time.Time
	StatusTime  time.Time
}

func NewMeteringAttributeBinding(skuId, attributeId string) *MeteringAttributeBinding {
	now := time.Now()
	return &MeteringAttributeBinding{
		BindingId:   NewMeteringAttributeBindingId(),
		SkuId:       skuId,
		AttributeId: attributeId,
		Status:      constants.StatusActive,
		CreateTime:  now,
		StatusTime:  now,
	}
}

func PbToMeteringAttributeBindings(pbMab *pb.CreateMeteringAttributeBindingsRequest) []*MeteringAttributeBinding {
	var mabs []*MeteringAttributeBinding
	for _, attId := range pbMab.GetAttributeIds() {
		mab := NewMeteringAttributeBinding(
			pbMab.GetSkuId().GetValue(),
			attId,
		)
		mabs = append(mabs, mab)
	}
	return mabs
}

func MeteringAttributeBindingToPb(mab *MeteringAttributeBinding) *pb.MeteringAttributeBinding {
	return &pb.MeteringAttributeBinding{
		BindingId:   pbutil.ToProtoString(mab.BindingId),
		SkuId:       pbutil.ToProtoString(mab.SkuId),
		AttributeId: pbutil.ToProtoString(mab.AttributeId),
		Status:      pbutil.ToProtoString(mab.Status),
		CreateTime:  pbutil.ToProtoTimestamp(mab.CreateTime),
		StatusTime:  pbutil.ToProtoTimestamp(mab.StatusTime),
	}
}
