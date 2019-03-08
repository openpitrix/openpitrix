// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package metering

import (
	"context"

	"github.com/fatih/structs"
	"openpitrix.io/openpitrix/pkg/gerr"
	"openpitrix.io/openpitrix/pkg/logger"
)

const (
	EN = "En"
	ZH = "Zh"

	NIL_STR = ""
)

const (
	InitAttId          = "att-000001"
	InitAttUnitHourId  = "att-unit-000001"
	InitAttUnitMonthId = "att-unit-000002"
	InitAttUnitYearId  = "att-unit-000003"
)

//struct english name and chinese name
var structDisName = map[string]map[string]string{
	"AttributeName": {
		EN: "attribute_name",
		ZH: "属性名",
	},
	"AttributeUnit": {
		EN: "attribute_unit",
		ZH: "属性单位",
	},
	"Attribute": {
		EN: "attribute",
		ZH: "属性",
	},
	"spu": {
		EN: "spu",
		ZH: "SPU",
	},
	"Sku": {
		EN: "sku",
		ZH: "SKU",
	},
	"Price": {
		EN: "price",
		ZH: "价格",
	},
	"Leasing": {
		EN: "leasing",
		ZH: "合约",
	},
}

func internalError(ctx context.Context, err error) error {
	return gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
}

//check if structObj exist
func checkStructExist(ctx context.Context, structObj interface{}, id string) error {
	structName := structs.Name(structObj)
	exist, err := checkExistById(ctx, structName, id)
	if err != nil {
		logger.Error(ctx, "Failed to get %s(%s), Error: [%+v]!", structName, id, err)
		return internalError(ctx, err)
	}
	if !exist {
		return notExistError(ctx, structObj, id)
	}
	return nil
}

func notExistError(ctx context.Context, structObj interface{}, id string) error {
	structName := structs.Name(structObj)
	logger.Error(ctx, "The %s(%s) not exist!", structName, id)
	a := []string{
		structDisName[structName][EN],
		structDisName[structName][EN],
		id,
		structDisName[structName][ZH],
		id,
	}
	return gerr.New(ctx, gerr.NotFound, gerr.ErrorNotExist, a)
}

func notExistInOtherError(ctx context.Context, structObj, targetStructObj interface{}) error {
	//SN: Struct Name
	structObjSN := structs.Name(structObj)
	targetStructObjSn := structs.Name(targetStructObj)
	logger.Error(ctx, "The %s not exist in %s!", structObjSN, targetStructObjSn)
	a := []string{
		structDisName[structObjSN][EN],
		structDisName[structObjSN][EN],
		structDisName[targetStructObjSn][EN],
		structDisName[targetStructObjSn][ZH],
		structDisName[structObjSN][ZH],
	}
	return gerr.New(ctx, gerr.NotFound, gerr.ErrorNotExistInOther, a)
}
