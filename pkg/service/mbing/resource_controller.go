// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package mbing

import (
	"context"

	"openpitrix.io/openpitrix/pkg/models"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/pi"
)

func insertLeasingsToDB(ctx context.Context, leasings []*models.Leasing) error {

	dbConn := pi.Global().DB(ctx)
	tx, err := dbConn.Session.BeginTx(ctx, nil)

	if err != nil {
		logger.Error(ctx, "Failed to connect mysql, Errors: [%+v]", err)
		return err
	}
	defer tx.RollbackUnlessCommitted()

	err = nil
	for _, leasing := range leasings {
		_, err := tx.InsertInto(constants.TableMeteringLeasing).Record(leasing).Exec()
		if err != nil {
			logger.Error(ctx, "Failed to save leasing [%+v].", leasing)
			tx.Rollback()
			break
		}
	}
	tx.Commit()
	return err
}
