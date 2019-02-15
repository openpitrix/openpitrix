// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package mbing

import (

"context"

"github.com/fatih/structs"
"openpitrix.io/openpitrix/pkg/gerr"
	"openpitrix.io/openpitrix/pkg/logger"
)

const (
	CreateFailedCode = 0
	NotExistCode     = 1
)

const (
	EN = "En"
	ZH = "Zh"
)

const (
	InitAttId = "att-000001"
	InitAttUnitHourId = "att-unit-000001"
	InitAttUnitMonthId = "att-unit-000002"
	InitAttUnitYearId = "att-unit-000003"
)

//var structDisName map[string]string

var structDisName = map[string]map[string]string {
	"Attribute": {
		EN: "attribute",
		ZH: "属性",
	},
	"AttributeUnit": {
		EN: "attribute_unit",
		ZH: "属性单位",
	},
	"AttributeValue": {
		EN: "attribute_value",
		ZH: "属性值",

	},
	"ResourceAttribute": {
		EN: "resource_attribute",
		ZH: "资源属性",
	},
	"Sku": {
		EN: "sku",
		ZH: "SKU",
	},
	"Price": {
		EN: "price",
		ZH: "定价",
	},
	"Leasing": {
		EN: "leasing",
		ZH: "合约",
	},
}

//check if existStructName exist when action actionStructName
func checkStructExistById(ctx context.Context, checkStruct, actionStruct interface{}, idValue string, actionErrType int8) error {
	checkStructName := structs.Name(checkStruct)
	exist, err := checkExistById(ctx, checkStructName, idValue)
	if err != nil {
		logger.Error(ctx, "Failed to get %s!", checkStructName)
		return commonInternalErr(ctx, actionStruct, actionErrType)
	}
	if !exist {
		logger.Error(ctx, "The %s that id is %s not exist!", idValue, checkStructName)
		return commonInternalErr(ctx, checkStruct, NotExistCode)
	}
	return nil
}

//CommonInternalErr: return error with gerr.ErrorMessage
func commonInternalErr(ctx context.Context, structObj interface{}, errType int8) error {
	structName := structs.Name(structObj)
	enName := structDisName[structName][EN]
	zhName := structDisName[structName][ZH]
	switch errType {
	case CreateFailedCode:
		return gerr.New(ctx, gerr.Internal, gerr.CreateFailed(enName, zhName))
	case NotExistCode:
		return gerr.New(ctx, gerr.Internal, gerr.NotExistError(enName, zhName))
	default:
		return gerr.New(ctx, gerr.Internal, gerr.ErrorUnknown)
	}
}
