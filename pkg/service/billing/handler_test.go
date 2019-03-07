// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package mbing

import (
	"fmt"
	"testing"

	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
)

const (
	ATTRIBUTE_NAME = "attribute_name"
	UNIT_NAME      = "unit_name"
	VALUE          = "value"
)

const (
	TEST_OFFSET = 0
	TEST_LIMIT  = 100
)

var test_attribute_names = []string{
	"node_num",
	"memory",
	"cpu",
	"disk",
	"region",
	"user_num",
	"stream",
}

var test_att_units = []string{
	"number",
	"mb",
	"gb",
	"tb",
}

var test_attributes = []map[string]string{
	//时长属性值
	{ATTRIBUTE_NAME: "duration", UNIT_NAME: "hour", VALUE: "1"},
	{ATTRIBUTE_NAME: "duration", UNIT_NAME: "month", VALUE: "1"},
	{ATTRIBUTE_NAME: "duration", UNIT_NAME: "year", VALUE: "1"},
	{ATTRIBUTE_NAME: "duration", UNIT_NAME: "hour", VALUE: "(0, 100]"},
	{ATTRIBUTE_NAME: "duration", UNIT_NAME: "hour", VALUE: "(100, 1000]"},
	{ATTRIBUTE_NAME: "duration", UNIT_NAME: "hour", VALUE: "(1000, ]"},
	//memory
	{ATTRIBUTE_NAME: "memory", UNIT_NAME: "gb", VALUE: "1"},
	{ATTRIBUTE_NAME: "memory", UNIT_NAME: "gb", VALUE: "2"},
	{ATTRIBUTE_NAME: "memory", UNIT_NAME: "gb", VALUE: "4"},
	{ATTRIBUTE_NAME: "memory", UNIT_NAME: "gb", VALUE: "8"},
	{ATTRIBUTE_NAME: "memory", UNIT_NAME: "gb", VALUE: "16"},
	{ATTRIBUTE_NAME: "memory", UNIT_NAME: "gb", VALUE: "32"},
	//CPU
	{ATTRIBUTE_NAME: "cpu", UNIT_NAME: "number", VALUE: "1"},
	{ATTRIBUTE_NAME: "cpu", UNIT_NAME: "number", VALUE: "2"},
	{ATTRIBUTE_NAME: "cpu", UNIT_NAME: "number", VALUE: "4"},
	{ATTRIBUTE_NAME: "cpu", UNIT_NAME: "number", VALUE: "8"},
	{ATTRIBUTE_NAME: "cpu", UNIT_NAME: "number", VALUE: "16"},
	//disk
	{ATTRIBUTE_NAME: "disk", UNIT_NAME: "gb", VALUE: "100"},
	{ATTRIBUTE_NAME: "disk", UNIT_NAME: "gb", VALUE: "500"},
	{ATTRIBUTE_NAME: "disk", UNIT_NAME: "tb", VALUE: "1"},
	{ATTRIBUTE_NAME: "disk", UNIT_NAME: "tb", VALUE: "10"},
	//region
	{ATTRIBUTE_NAME: "region", VALUE: "ap2a"},
	{ATTRIBUTE_NAME: "region", VALUE: "pek3"},
	//instance
	{ATTRIBUTE_NAME: "node_num", UNIT_NAME: "number", VALUE: "(0, 5]"},
	{ATTRIBUTE_NAME: "node_num", UNIT_NAME: "number", VALUE: "(5, 20]"},
	{ATTRIBUTE_NAME: "node_num", UNIT_NAME: "number", VALUE: "(20,]"},
	//user
	{ATTRIBUTE_NAME: "user_num", UNIT_NAME: "number", VALUE: "(0, 10]"},
	{ATTRIBUTE_NAME: "user_num", UNIT_NAME: "number", VALUE: "(10, 50]"},
	{ATTRIBUTE_NAME: "user_num", UNIT_NAME: "number", VALUE: "(50, ]"},
	//stream
	{ATTRIBUTE_NAME: "stream", UNIT_NAME: "mb", VALUE: "1"},
	{ATTRIBUTE_NAME: "stream", UNIT_NAME: "gb", VALUE: "1"},
}

//AttributeName
func TestCreateAttributeName(t *testing.T) {
	for _, attName := range test_attribute_names {
		t.Run(fmt.Sprintf("CreateAttribute_%s", attName),
			testCreateAttributeNameFunc(attName))
	}
}

func testCreateAttributeNameFunc(name string) func(t *testing.T) {
	return func(t *testing.T) {
		var attNameReq = &pb.CreateAttributeNameRequest{
			Name: pbutil.ToProtoString(name),
		}

		response, err := ss.server.CreateAttributeName(ss.ctx, attNameReq)
		if err != nil {
			t.Error(err)
		} else {
			t.Logf("TestCreateAttribute: Insert attribute_name(%s) successfully.", response.GetAttributeNameId().GetValue())
		}
	}
}

func TestDescribeAttributeNames(t *testing.T) {
	describeReq := pb.DescribeAttributeNamesRequest{
		Offset: TEST_OFFSET,
		Limit:  TEST_LIMIT,
	}

	res, err := ss.server.DescribeAttributeNames(ss.ctx, &describeReq)
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("Describe %d AttributeNames: ", len(res.AttributeNames))
		for _, attName := range res.AttributeNames {
			t.Logf("attribute_name_id: %s, name: %s",
				attName.GetAttributeNameId().GetValue(),
				attName.GetName().GetValue())
		}
	}
}

//Attribute_Unit
func TestCreateAttributeUnit(t *testing.T) {
	for _, attUnit := range test_att_units {
		t.Run(fmt.Sprintf("Create_attribute_unit_%s", attUnit),
			testCreateAttributeUnitFunc(attUnit))
	}
}

func testCreateAttributeUnitFunc(name string) func(t *testing.T) {
	return func(t *testing.T) {
		var attUnitReq = &pb.CreateAttributeUnitRequest{
			Name: pbutil.ToProtoString(name),
		}

		response, err := ss.server.CreateAttributeUnit(ss.ctx, attUnitReq)
		if err != nil {
			t.Error(err)
		} else {
			t.Logf("TestCreateAttributeUnit: Insert attribute_unit(%s) successfully.", response.GetAttributeUnitId().GetValue())
		}
	}
}

func TestDescribeAttributeUnits(t *testing.T) {
	describeReq := pb.DescribeAttributeUnitsRequest{
		Offset: TEST_OFFSET,
		Limit:  TEST_LIMIT,
	}

	res, err := ss.server.DescribeAttributeUnits(ss.ctx, &describeReq)
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("List %d Attribute_units: ", len(res.AttributeUnits))
		for _, attUnit := range res.AttributeUnits {
			t.Logf("attribute_unit_id: %s, unit_name: %s",
				attUnit.GetAttributeUnitId().GetValue(),
				attUnit.GetName().GetValue())
		}
	}
}

//Attribute
//need to run TestCreateAttributeName and TestCreateAttributeUnit
func TestCreateAttribute(t *testing.T) {
	//get attribute_names
	desAttNameReq := pb.DescribeAttributeNamesRequest{
		Offset: TEST_OFFSET,
		Limit:  TEST_LIMIT,
	}
	attNames, err := ss.server.DescribeAttributeNames(ss.ctx, &desAttNameReq)
	if err != nil {
		t.Error(err)
	}

	//get attribute_units
	desAttUnitReq := pb.DescribeAttributeUnitsRequest{
		Offset: TEST_OFFSET,
		Limit:  TEST_LIMIT,
	}
	attUnits, err := ss.server.DescribeAttributeUnits(ss.ctx, &desAttUnitReq)
	if err != nil {
		t.Error(err)
	}

	//generate and create Attribute
	for _, tAtt := range test_attributes {
		t.Run("CreateAttribute", testCreateAttributeFunc(tAtt[ATTRIBUTE_NAME],
			tAtt[UNIT_NAME],
			tAtt[VALUE],
			attNames.AttributeNames,
			attUnits.AttributeUnits))
	}
	t.Logf("Create %d attributes Successfully.", len(test_attributes))

}

func testCreateAttributeFunc(attName, attUnitName, value string,
	attNames []*pb.AttributeName,
	attUnits []*pb.AttributeUnit) func(t *testing.T) {

	return func(t *testing.T) {
		//get attribute_name and attribute_unit
		var attNameId, attUnitId string
		t.Logf("attribute_name: %s", attName)
		for _, attNameObj := range attNames {
			if attName == attNameObj.GetName().GetValue() {
				attNameId = attNameObj.GetAttributeNameId().GetValue()
				break
			}
		}

		if attUnitName != "" {
			t.Logf("attribute_unit_name: %s", attUnitName)
			for _, attUnitObj := range attUnits {
				if attUnitName == attUnitObj.GetName().GetValue() {
					attUnitId = attUnitObj.GetAttributeUnitId().GetValue()
					break
				}
			}
		}

		//generate CreateAttributeRequest
		attReq := &pb.CreateAttributeRequest{
			AttributeNameId: pbutil.ToProtoString(attNameId),
			AttributeUnitId: pbutil.ToProtoString(attUnitId),
			Value:           pbutil.ToProtoString(value),
		}

		//create attribute
		res, err := ss.server.CreateAttribute(ss.ctx, attReq)
		if err != nil {
			t.Skipf("Failed to create attribute(%s), Error: [%+v]", attName, err)
		} else {
			t.Logf("Create attribute(%s) successfully.", res.GetAttributeId().GetValue())
		}
	}
}

func TestDescribeAttributes(t *testing.T) {
	describeReq := pb.DescribeAttributesRequest{
		Offset: TEST_OFFSET,
		Limit:  TEST_LIMIT,
	}

	res, err := ss.server.DescribeAttributes(ss.ctx, &describeReq)
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("List %d Attributes: ", len(res.Attributes))
		for _, att := range res.Attributes {
			t.Logf("attribute_id: %s, attribute_name_id: %s, "+
				"attribute_unit_id: %s, value: %s",
				att.GetAttributeId().GetValue(),
				att.GetAttributeNameId().GetValue(),
				att.GetAttributeUnitId().GetValue(),
				att.GetValue().GetValue())
		}
	}
}
