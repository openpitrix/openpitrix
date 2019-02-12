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

func insertLeasingsToDB(ctx context.Context, leasings []*models.Leasing) error {

	dbConn := pi.Global().DB(ctx)
	//tx, err := dbConn.Session.BeginTx(ctx, nil)
	//
	//if err != nil {
	//	logger.Error(ctx, "Failed to connect mysql, Errors: [%+v]", err)
	//	return err
	//}
	//defer tx.RollbackUnlessCommitted()

	var err error
	for _, leasing := range leasings {
		_, err := dbConn.InsertInto(constants.TableMbing).Record(leasing).Exec()
		if err != nil {
			logger.Error(ctx, "Failed to save leasing: [%+v].", leasing)
			break
		}
	}
	return err
}
