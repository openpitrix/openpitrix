// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package mbing

import (
	"fmt"
	"openpitrix.io/openpitrix/pkg/models"
)


func Billing() {
	leasing := models.Leasing{}
	contract, err := calculate(leasing)
	fmt.Printf("%v", err)

	_, err = Charge(contract)

	if err.Error() == "balance not enough" {
		addToNoMoney(leasing)
	}

}

func calculate(leasing models.Leasing) (*models.LeasingContract, error) {
	return &models.LeasingContract{}, nil
}

func addToNoMoney(leasing models.Leasing) error {
	return nil
}
