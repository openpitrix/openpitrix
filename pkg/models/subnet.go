// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import "time"

type Subnet struct {
	Zone        string
	SubnetId    string
	Name        string
	CreateTime  time.Time
	Description string
	InstanceIds []string
	VpcId       string
}
