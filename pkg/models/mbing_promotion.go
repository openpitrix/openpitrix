// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import "time"

type CombinationResourceAttribute struct {
	Id             		string
	ResourceVersionIds  []string
	Attributes         	map[string]string //{resourceVersionId: attributeId, ..}
	MeteringAttributes  map[string]string //{resourceVersionId: attributeId, ..}
	CreateTime          time.Time
	UpdateTime         	time.Time
	Status 				string
}

type CombinationSku struct {
	Id             		string
	CRAId             	string
	AttributeValues     map[string]string  //{resourceVersionId: valueId, ..}
	CreateTime          time.Time
	UpdateTime         	time.Time
	Status 				string
}

type CombinationPrice struct {
	Id             		string
	CombinationSkuId    string
	ResourceVersionId   string
	AttributeId 		string
	Prices		 		map[string]float32 //StepPrice: {valueId: price, ..}
	Currency 			string
	CreateTime          time.Time
	UpdateTime         	time.Time
	Status 				string
}

type ProbationSku struct {
	Id             		string
	ResourceAttributeId string
	AttributeValues     []string

	AttributeId 		string
	Prices		 		map[string]float32
	Currency 			string
	CreateTime          time.Time
	UpdateTime         	time.Time
	Status 				string
}

