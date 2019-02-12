// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package mbing

import (
	"context"

	"github.com/fatih/structs"
	"openpitrix.io/openpitrix/pkg/gerr"
)

const (
	CreateFailedCode = 0
	NotExistCode     = 1
)

var structDisName map[string]string

func init() {
	structDisName = map[string]string{
		"AttributeEn":         "attribute",
		"AttributeZh":         "属性",
		"AttributeUnitEn":     "attribute_unit",
		"AttributeUnitZh":     "属性单位",
		"AttributeValueEn":    "attribute_value",
		"AttributeValueZh":    "属性值",
		"ResourceAttributeEn": "resource_attribute",
		"ResourceAttributeZh": "资源属性",
		"SkuEn":               "sku",
		"SkuZh":               "SKU",
		"PriceEn":             "price",
		"PriceZh":             "定价",
	}
}

//check if existStructName exist when action actionStructName
func checkStructExistById(ctx context.Context, checkStruct, actionStruct interface{}, idValue string, actionErrType int8) error {
	checkStructName := structs.Name(checkStruct)
	exist, err := checkExistById(ctx, checkStructName, idValue)
	if err != nil {
		return commonInternalErr(ctx, actionStruct, actionErrType)
	}
	if !exist {
		return commonInternalErr(ctx, checkStruct, NotExistCode)
	}
	return nil
}

//CommonInternalErr: return error with gerr.ErrorMessage
func commonInternalErr(ctx context.Context, structObj interface{}, errType int8) error {
	structName := structs.Name(structObj)
	enName := structDisName[structName+"En"]
	zhName := structDisName[structName+"Zh"]
	switch errType {
	case CreateFailedCode:
		return gerr.New(ctx, gerr.Internal, gerr.CreateFailed(enName, zhName))
	case NotExistCode:
		return gerr.New(ctx, gerr.Internal, gerr.NotExistError(enName, zhName))
	default:
		return gerr.New(ctx, gerr.Internal, gerr.ErrorUnknown)
	}
}
