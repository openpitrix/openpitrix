// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import "time"

type Charge struct {
	ChargeId   string
	ContractId string
	Fee        float64
	Currency   string
	Status     string
	CreateTime time.Time
}

type Refund struct {
	RefundId   string
	ContractId string
	Fee        float64
	Currency   string
	Status     string
	CreateTime time.Time
}

type Recharge struct {
	ReChargeId  string
	Balance     float64
	Currency    string
	Status      string
	CreateTime  time.Time
	Description string
}

type Income struct {
	IncomeId   string
	ContractId string
	Balance    string
	Currency   string
	CreateTime time.Time
}
