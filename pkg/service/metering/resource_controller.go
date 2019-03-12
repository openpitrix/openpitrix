// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package metering

import (
	"context"
	"time"

	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
	"openpitrix.io/openpitrix/pkg/util/stringutil"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/pi"
)

func checkExistById(ctx context.Context, structName, idValue string) (bool, error) {
	tableName := stringutil.CamelCaseToUnderscore(structName)
	idName := tableName + "_id"
	count, err := pi.Global().DB(ctx).
		Select(idName).
		From(tableName).
		Where(db.Eq(idName, idValue)).
		Limit(1).Count()

	if err != nil {
		logger.Error(ctx, "Failed to connect DB: [%+v]", err)
		return false, err
	}

	if count > 0 {
		return true, nil
	}
	return false, nil
}

//Sku
func insertAttributeName(ctx context.Context, attributeName *models.AttributeName) error {
	_, err := pi.Global().DB(ctx).
		InsertInto(constants.TableAttributeName).
		Record(attributeName).
		Exec()
	if err != nil {
		logger.Error(ctx, "Failed to insert attribute_name: [%+v]", err)
	}
	return err
}

func DescribeAttributeNames(ctx context.Context, req *pb.DescribeAttributeNamesRequest) ([]*models.AttributeName, error) {
	var attributeNames []*models.AttributeName
	_, err := pi.Global().DB(ctx).
		Select(constants.IndexedColumns[constants.TableAttributeName]...).
		From(constants.TableAttributeName).
		Where(db.Eq(constants.ColumnStatus, constants.StatusActive)).
		Where(manager.BuildFilterConditions(req, constants.TableAttributeName)).
		OrderDir(constants.ColumnCreateTime, false).
		Offset(pbutil.GetOffsetFromRequest(req)).
		Limit(pbutil.GetLimitFromRequest(req)).
		Load(&attributeNames)

	if err != nil {
		logger.Error(ctx, "Failed to describe attribute_name: [%+v]", err)
	}
	return attributeNames, err
}

func insertAttributeUnit(ctx context.Context, attUnit *models.AttributeUnit) error {
	_, err := pi.Global().DB(ctx).
		InsertInto(constants.TableAttributeUnit).
		Record(attUnit).
		Exec()
	if err != nil {
		logger.Error(ctx, "Failed to insert attribute_unit: [%+v]", err)
	}
	return err
}

func DescribeAttributeUnits(ctx context.Context, req *pb.DescribeAttributeUnitsRequest) ([]*models.AttributeUnit, error) {
	var attUnits []*models.AttributeUnit
	_, err := pi.Global().DB(ctx).
		Select(constants.IndexedColumns[constants.TableAttributeUnit]...).
		From(constants.TableAttributeUnit).
		Where(db.Eq(constants.ColumnStatus, constants.StatusActive)).
		Where(manager.BuildFilterConditions(req, constants.TableAttributeUnit)).
		OrderDir(constants.ColumnCreateTime, false).
		Offset(pbutil.GetOffsetFromRequest(req)).
		Limit(pbutil.GetLimitFromRequest(req)).
		Load(&attUnits)

	if err != nil {
		logger.Error(ctx, "Failed to describe attribute_unit: [%+v]", err)
	}
	return attUnits, err
}

func insertAttribute(ctx context.Context, attribute *models.Attribute) error {
	_, err := pi.Global().DB(ctx).
		InsertInto(constants.TableAttribute).
		Record(attribute).
		Exec()
	if err != nil {
		logger.Error(ctx, "Failed to insert attribute: [%+v]", err)
	}
	return err
}

func getAttribute(ctx context.Context, attributeId string) (*models.Attribute, error) {
	attribute := &models.Attribute{}
	err := pi.Global().DB(ctx).
		Select(models.AttributeColumns...).
		From(constants.TableAttribute).
		Where(db.Eq(constants.ColumnStatus, constants.StatusActive)).
		Where(db.Eq(constants.ColumnAttributeId, attributeId)).
		LoadOne(&attribute)
	if err != nil {
		logger.Error(ctx, "Failed to get attribute: [%+v]", err)
	}
	if err == db.ErrNotFound {
		return nil, nil
	}
	return attribute, err
}

func DescribeAttributes(ctx context.Context, req *pb.DescribeAttributesRequest) ([]*models.Attribute, error) {
	var attributes []*models.Attribute
	_, err := pi.Global().DB(ctx).
		Select(constants.IndexedColumns[constants.TableAttribute]...).
		From(constants.TableAttribute).
		Where(db.Eq(constants.ColumnStatus, constants.StatusActive)).
		Where(manager.BuildFilterConditions(req, constants.TableAttribute)).
		OrderDir(constants.ColumnCreateTime, false).
		Offset(pbutil.GetOffsetFromRequest(req)).
		Limit(pbutil.GetLimitFromRequest(req)).
		Load(&attributes)

	if err != nil {
		logger.Error(ctx, "Failed to describe attribute: [%+v]", err)
	}
	return attributes, err
}

func insertSpu(ctx context.Context, spu *models.Spu) error {
	_, err := pi.Global().DB(ctx).
		InsertInto(constants.TableSpu).
		Record(spu).
		Exec()
	if err != nil {
		logger.Error(ctx, "Failed to insert spu: [%+v]", err)
	}
	return err
}

func getSpu(ctx context.Context, spuId string) (*models.Spu, error) {
	spu := &models.Spu{}
	err := pi.Global().DB(ctx).
		Select(constants.IndexedColumns[constants.TableSpu]...).
		From(constants.TableSpu).
		Where(db.Eq(constants.ColumnStatus, constants.StatusActive)).
		Where(db.Eq(constants.ColumnSpuId, spuId)).
		LoadOne(&spu)
	if err != nil {
		logger.Error(ctx, "Failed to get spu: [%+v]", err)
	}
	if err == db.ErrNotFound {
		return nil, nil
	}
	return spu, err
}

func insertSku(ctx context.Context, sku *models.Sku) error {
	_, err := pi.Global().DB(ctx).
		InsertInto(constants.TableSku).
		Record(sku).
		Exec()
	if err != nil {
		logger.Error(ctx, "Failed to insert sku: [%+v]", err)
	}
	return err
}

func getSku(ctx context.Context, skuId string) (*models.Sku, error) {
	sku := &models.Sku{}
	err := pi.Global().DB(ctx).
		Select(constants.IndexedColumns[constants.TableSku]...).
		From(constants.TableSku).
		Where(db.Eq(constants.ColumnStatus, constants.StatusActive)).
		Where(db.Eq(constants.ColumnSkuId, skuId)).
		LoadOne(sku)

	if err != nil {
		logger.Error(ctx, "Failed to get sku: [%+v]", err)
	}
	if err == db.ErrNotFound {
		return nil, nil
	}
	return sku, err
}

func insertPrice(ctx context.Context, price *models.Price) error {
	_, err := pi.Global().DB(ctx).
		InsertInto(constants.TablePrice).
		Record(price).
		Exec()
	if err != nil {
		logger.Error(ctx, "Failed to insert price: [%+v]", err)
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

func getLeasing(ctx context.Context, leasingID, resourceId, skuId string) (*models.Leasing, error) {
	//ids: leasingId, resourceId, skuId
	//TODO: impl get leasing
	return &models.Leasing{}, nil
}

//promotion
func insertCombination(ctx context.Context, com *models.Combination) error {
	_, err := pi.Global().DB(ctx).
		InsertInto(constants.TableCombination).
		Record(com).Exec()
	if err != nil {
		logger.Error(ctx, "Failed to insert combination, Error: [%+v].", err)
	} else {
		logger.Info(ctx, "Insert combination successfully.")
	}
	return err
}

func insertCombinationSku(ctx context.Context, comSku *models.CombinationSku) error {
	_, err := pi.Global().DB(ctx).
		InsertInto(constants.TableCombinationSku).
		Record(comSku).Exec()
	if err != nil {
		logger.Error(ctx, "Failed to insert combination_sku: [%+v].", err)
	} else {
		logger.Info(ctx, "Insert combination_sku successfully.")
	}
	return err
}

func leasingToEtcd(leasing models.Leasing) error {
	//TODO: impl add leasing to etcd
	return nil
}

func leasingToRedis(leasingId string, renewalTime time.Time) error {
	//TODO: impl add leasing to redis
	return nil
}
