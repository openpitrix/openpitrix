// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"time"
)

type Vpc struct {
	VpcId            string
	Name             string
	CreateTime       time.Time
	Description      string
	Status           string
	TransitionStatus string
	Subnets          []string
	Eip              *Eip
}

type Eip struct {
	EipId string
	Name  string
	Addr  string
}
