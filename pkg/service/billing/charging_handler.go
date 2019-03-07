// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package mbing

import (
	"openpitrix.io/openpitrix/pkg/models"
)

func Charge(contract *models.LeasingContract) (string, error) {
	return "Charge.id", nil
}

func ReChargeFromSys(contract *models.LeasingContract) (string, error) {
	return "ReCharge.id", nil
}

func ReCharge(userId string, currency string, fee float32) (string, error) {
	return "ReCharge.id", nil
}
