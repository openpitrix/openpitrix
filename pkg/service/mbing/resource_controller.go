// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package mbing

import (
	"context"

	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/util/stringutil"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/pi"
)

func checkExistById(ctx context.Context, structName, idValue string) (bool, error) {
	tableName := stringutil.CamelCaseToUnderscore(structName)
	idName := stringutil.UnderscoreToCamelCase(tableName) + "Id"
	count, err := pi.Global().DB(ctx).
		Select(idName).
		From(tableName).
		Where(db.Eq(idName, idValue)).
		Limit(1).Count()

	exist := false
	if err != nil {
		logger.Error(ctx, "Failed to connect DB, Errors: [%+v]", err)
		return exist, err
	}

	if count > 0 {
		exist = true
	}
	return exist, nil
}

//Sku
func insertAttribute(ctx context.Context, attribute *models.Attribute) error {
	_, err := pi.Global().DB(ctx).
		InsertInto(constants.TableAttribute).
		Record(attribute).
		Exec()
	if err != nil {
		logger.Error(ctx, "Failed to insert attribute, Errors: [%+v]", err)
	}
	return err
}

func getAttributeById(ctx context.Context, attributeId string) (*models.Attribute, error) {
	att := &models.Attribute{}
	err := pi.Global().DB(ctx).
		Select(models.AttributeColumns...).
		From(constants.TableAttribute).
		Where(db.Eq(constants.ColumnAttributeId, attributeId)).
		LoadOne(&att)

	if err != nil {
		logger.Error(ctx, "Failed to get attribute, Errors: [%+v]", err)
	}
	return att, err
}

func insertAttributeUnit(ctx context.Context, attUnit *models.AttributeUnit) error {
	_, err := pi.Global().DB(ctx).
		InsertInto(constants.TableAttUnit).
		Record(attUnit).
		Exec()
	if err != nil {
		logger.Error(ctx, "Failed to insert attribute_unit, Errors: [%+v]", err)
	}
	return err
}

func getAttUnitById(ctx context.Context, attUnitId string) (*models.AttributeUnit, error) {
	attUnit := &models.AttributeUnit{}
	err := pi.Global().DB(ctx).
		Select(models.AttributeUnitColumns...).
		From(constants.TableAttUnit).
		Where(db.Eq(constants.ColumnAttUnitId, attUnitId)).
		LoadOne(&attUnit)

	if err != nil {
		logger.Error(ctx, "Failed to get attribuite_unit, Errors: [%+v]", err)
	}
	return attUnit, err
}

func insertAttributeValue(ctx context.Context, attValue *models.AttributeValue) error {
	_, err := pi.Global().DB(ctx).
		InsertInto(constants.TableAttValue).
		Record(attValue).
		Exec()
	if err != nil {
		logger.Error(ctx, "Failed to insert attribute_value, Errors: [%+v]", err)
	}
	return err
}

func insertResourceAttribute(ctx context.Context, resAtt *models.ResourceAttribute) error {
	_, err := pi.Global().DB(ctx).
		InsertInto(constants.TableResAtt).
		Record(resAtt).
		Exec()
	if err != nil {
		logger.Error(ctx, "Failed to insert resource_attribute, Errors: [%+v]", err)
	}
	return err
}

func insertSku(ctx context.Context, sku *models.Sku) error {
	_, err := pi.Global().DB(ctx).
		InsertInto(constants.TableSku).
		Record(sku).
		Exec()
	if err != nil {
		logger.Error(ctx, "Failed to insert sku, Errors: [%+v]", err)
	}
	return err
}

func getSkuById(ctx context.Context, skuId string) (*models.Sku, error) {
	sku := &models.Sku{}
	err := pi.Global().DB(ctx).
		Select(models.SkuColumns...).
		From(constants.TableSku).
		Where(db.Eq(constants.ColumnSkuId, skuId)).
		LoadOne(sku)

	if err != nil {
		logger.Error(ctx, "Failed to get sku, Errors: [%+v]", err)
	}
	return sku, err
}

func insertPrice(ctx context.Context, price *models.Price) error {
	_, err := pi.Global().DB(ctx).
		InsertInto(constants.TablePrice).
		Record(price).
		Exec()
	if err != nil {
		logger.Error(ctx, "Failed to insert price, Errors: [%+v]", err)
	}
	return err
}

//Metering
func insertLeasings(ctx context.Context, leasings []*models.Leasing) error {
	result, err := pi.Global().DB(ctx).
		InsertInto(constants.TableLeasing).
		Record(leasings).Exec()
	if err != nil {
		logger.Error(ctx, "Failed to insert leasings, Error: [%+v].", err)
	} else {
		count, _ := result.RowsAffected()
		logger.Info(ctx, "Insert %d leasings successfully.", count)
	}
	return err
}

//promotion
func insertCRA(ctx context.Context, cra *models.CombinationResourceAttribute) error {
	_, err := pi.Global().DB(ctx).
		InsertInto(constants.TableCRA).
		Record(cra).Exec()
	if err != nil {
		logger.Error(ctx, "Failed to insert Combination_Resource_Attribute, Error: [%+v].", err)
	} else {
		logger.Info(ctx, "Insert Combination_Resource_Attribute successfully.")
	}
	return err
}

func insertComSku(ctx context.Context, cs *models.CombinationSku) error {
	_, err := pi.Global().DB(ctx).
		InsertInto(constants.TableCS).
		Record(cs).Exec()
	if err != nil {
		logger.Error(ctx, "Failed to insert Combination_Sku, Error: [%+v].", err)
	} else {
		logger.Info(ctx, "Insert Combination_Sku successfully.")
	}
	return err
}

func insertComPrice(ctx context.Context, comPrice *models.CombinationPrice) error {
	_, err := pi.Global().DB(ctx).
		InsertInto(constants.TableComPrice).
		Record(comPrice).Exec()
	if err != nil {
		logger.Error(ctx, "Failed to insert Combination_Price, Error: [%+v].", err)
	} else {
		logger.Info(ctx, "Insert Combination_Price successfully.")
	}
	return err
}

func insertProSku(ctx context.Context, proSku *models.ProbationSku) error {
	_, err := pi.Global().DB(ctx).
		InsertInto(constants.TableProbationSku).
		Record(proSku).Exec()
	if err != nil {
		logger.Error(ctx, "Failed to insert probation_sku, Error: [%+v].", err)
	} else {
		logger.Info(ctx, "Insert probation_sku successfully.")
	}
	return err
}
