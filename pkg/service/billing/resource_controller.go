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
	idName := tableName + "_id"
	count, err := pi.Global().DB(ctx).
		Select(idName).
		From(tableName).
		Where(db.Eq(idName, idValue)).
		Limit(1).Count()

	if err != nil {
		logger.Error(ctx, "Failed to connect DB, Errors: [%+v]", err)
		return false, err
	}

	if count > 0 {
		return true, nil
	}
	return false, nil
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

//promotion
func insertCombinationPrice(ctx context.Context, comPrice *models.CombinationPrice) error {
	_, err := pi.Global().DB(ctx).
		InsertInto(constants.TableCombinationPrice).
		Record(comPrice).Exec()
	if err != nil {
		logger.Error(ctx, "Failed to insert Combination_Price, Error: [%+v].", err)
	} else {
		logger.Info(ctx, "Insert Combination_Price successfully.")
	}
	return err
}

func insertProbation(ctx context.Context, probation *models.Probation) error {
	_, err := pi.Global().DB(ctx).
		InsertInto(constants.TableProbation).
		Record(probation).Exec()
	if err != nil {
		logger.Error(ctx, "Failed to insert probation, Error: [%+v].", err)
	} else {
		logger.Info(ctx, "Insert probation successfully.")
	}
	return err
}

func insertLeasingContract(ctx context.Context, contract *models.LeasingContract) error {
	//TODO: impl insert
	return nil
}

func getLeasingContract(ctx context.Context, contractId, leasingId string) (*models.LeasingContract, error){
	//TODO: impl
	return &models.LeasingContract{}, nil
}

func insufficientBalanceToEtcd(resourceId, skuId, userId string) error {
	//TODO: add resourceId, skuId, userId to Etcd(insufficient queue)
	return nil
}