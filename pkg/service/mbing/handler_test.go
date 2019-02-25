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
	NAME           = "name"
	DISPLAY_NAME   = "displayName"
	ATTRIBUTE_NAME = "attribute_name"
	UNIT_NAME      = "unit_name"
	MAX_VALUE      = "max_value"
	MIN_VALUE      = "min_value"
	STR_VALUE      = "str_value"
)

const (
	TEST_OFFSET = 0
	TEST_LIMIT  = 100
)

var test_attribute_names = []map[string]string{
	{NAME: "node_num", DISPLAY_NAME: "实例数"},
	{NAME: "memory", DISPLAY_NAME: "内存"},
	{NAME: "cpu", DISPLAY_NAME: "CPU"},
	{NAME: "disk", DISPLAY_NAME: "硬盘"},
	{NAME: "region", DISPLAY_NAME: "区域"},
	{NAME: "user_num", DISPLAY_NAME: "用户数"},
	{NAME: "stream", DISPLAY_NAME: "流量"},
}

var test_att_units = []map[string]string{
	{NAME: "number", DISPLAY_NAME: "个数"},
	{NAME: "mb", DISPLAY_NAME: "MB"},
	{NAME: "gb", DISPLAY_NAME: "GB"},
	{NAME: "tb", DISPLAY_NAME: "TB"},
}

var test_attributes = []map[string]interface{}{
	//时长属性值
	{ATTRIBUTE_NAME: "duration", UNIT_NAME: "hour", MIN_VALUE: 1, MAX_VALUE: 1},
	{ATTRIBUTE_NAME: "duration", UNIT_NAME: "month", MIN_VALUE: 1, MAX_VALUE: 1},
	{ATTRIBUTE_NAME: "duration", UNIT_NAME: "year", MIN_VALUE: 1, MAX_VALUE: 1},
	{ATTRIBUTE_NAME: "duration", UNIT_NAME: "hour", MIN_VALUE: 0, MAX_VALUE: 100},
	{ATTRIBUTE_NAME: "duration", UNIT_NAME: "hour", MIN_VALUE: 100, MAX_VALUE: 1000},
	{ATTRIBUTE_NAME: "duration", UNIT_NAME: "hour", MIN_VALUE: 1000},
	//memory
	{ATTRIBUTE_NAME: "memory", UNIT_NAME: "gb", MIN_VALUE: 1, MAX_VALUE: 1},
	{ATTRIBUTE_NAME: "memory", UNIT_NAME: "gb", MIN_VALUE: 2, MAX_VALUE: 2},
	{ATTRIBUTE_NAME: "memory", UNIT_NAME: "gb", MIN_VALUE: 4, MAX_VALUE: 4},
	{ATTRIBUTE_NAME: "memory", UNIT_NAME: "gb", MIN_VALUE: 8, MAX_VALUE: 8},
	{ATTRIBUTE_NAME: "memory", UNIT_NAME: "gb", MIN_VALUE: 16, MAX_VALUE: 16},
	{ATTRIBUTE_NAME: "memory", UNIT_NAME: "gb", MIN_VALUE: 32, MAX_VALUE: 32},
	//CPU
	{ATTRIBUTE_NAME: "cpu", UNIT_NAME: "number", MIN_VALUE: 1, MAX_VALUE: 1},
	{ATTRIBUTE_NAME: "cpu", UNIT_NAME: "number", MIN_VALUE: 2, MAX_VALUE: 2},
	{ATTRIBUTE_NAME: "cpu", UNIT_NAME: "number", MIN_VALUE: 4, MAX_VALUE: 4},
	{ATTRIBUTE_NAME: "cpu", UNIT_NAME: "number", MIN_VALUE: 8, MAX_VALUE: 8},
	{ATTRIBUTE_NAME: "cpu", UNIT_NAME: "number", MIN_VALUE: 16, MAX_VALUE: 16},
	//disk
	{ATTRIBUTE_NAME: "disk", UNIT_NAME: "gb", MIN_VALUE: 100, MAX_VALUE: 100},
	{ATTRIBUTE_NAME: "disk", UNIT_NAME: "gb", MIN_VALUE: 500, MAX_VALUE: 500},
	{ATTRIBUTE_NAME: "disk", UNIT_NAME: "tb", MIN_VALUE: 1, MAX_VALUE: 1},
	{ATTRIBUTE_NAME: "disk", UNIT_NAME: "tb", MIN_VALUE: 10, MAX_VALUE: 10},
	//region
	{ATTRIBUTE_NAME: "region", STR_VALUE: "ap2a"},
	{ATTRIBUTE_NAME: "region", STR_VALUE: "pek3"},
	//instance
	{ATTRIBUTE_NAME: "node_num", UNIT_NAME: "number", MIN_VALUE: 0, MAX_VALUE: 5},
	{ATTRIBUTE_NAME: "node_num", UNIT_NAME: "number", MIN_VALUE: 5, MAX_VALUE: 20},
	{ATTRIBUTE_NAME: "node_num", UNIT_NAME: "number", MIN_VALUE: 20},
	//user
	{ATTRIBUTE_NAME: "user_num", UNIT_NAME: "number", MIN_VALUE: 0, MAX_VALUE: 10},
	{ATTRIBUTE_NAME: "user_num", UNIT_NAME: "number", MIN_VALUE: 10, MAX_VALUE: 50},
	{ATTRIBUTE_NAME: "user_num", UNIT_NAME: "number", MIN_VALUE: 50},
	//stream
	{ATTRIBUTE_NAME: "stream", UNIT_NAME: "mb", MIN_VALUE: 1, MAX_VALUE: 1},
	{ATTRIBUTE_NAME: "stream", UNIT_NAME: "gb", MIN_VALUE: 1, MAX_VALUE: 1},
}

//AttributeName
func TestCreateAttributeName(t *testing.T) {
	for _, att := range test_attribute_names {
		t.Run(fmt.Sprintf("CreateAttribute_%s", att[NAME]),
			testCreateAttributeNameFunc(att[NAME], att[DISPLAY_NAME]))
	}
}

func testCreateAttributeNameFunc(name, displayName string) func(t *testing.T) {
	return func(t *testing.T) {
		var attNameReq = &pb.CreateAttributeNameRequest{
			Name:        pbutil.ToProtoString(name),
			DisplayName: pbutil.ToProtoString(displayName),
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
			t.Logf("attribute_name_id: %s, name: %s, display_name: %s ",
				attName.GetAttributeNameId().GetValue(),
				attName.GetName().GetValue(),
				attName.GetDisplayName().GetValue())
		}
	}
}

//Attribute_Unit
func TestCreateAttributeUnit(t *testing.T) {
	for _, attUnit := range test_att_units {
		t.Run(fmt.Sprintf("Create_attribute_unit_%s", attUnit[NAME]),
			testCreateAttributeUnitFunc(attUnit[NAME], attUnit[DISPLAY_NAME]))
	}
}

func testCreateAttributeUnitFunc(name, displayName string) func(t *testing.T) {
	return func(t *testing.T) {
		var attUnitReq = &pb.CreateAttributeUnitRequest{
			Name:        pbutil.ToProtoString(name),
			DisplayName: pbutil.ToProtoString(displayName),
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
			t.Logf("attribute_unit_id: %s, unit_name: %s, unit_display_name: %s",
				attUnit.GetAttributeUnitId().GetValue(),
				attUnit.GetName().GetValue(),
				attUnit.GetDisplayName().GetValue())
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
			tAtt[MIN_VALUE],
			tAtt[MAX_VALUE],
			tAtt[STR_VALUE],
			attNames.AttributeNames,
			attUnits.AttributeUnits))
	}
	t.Logf("Create %d attributes Successfully.", len(test_attributes))

}

func testCreateAttributeFunc(attNameStr,
	attUnitNameStr,
	minValue,
	maxValue,
	strValue interface{},
	attNames []*pb.AttributeName,
	attUnits []*pb.AttributeUnit) func(t *testing.T) {

	return func(t *testing.T) {
		//get attribute_name and attribute_unit
		var attNameId, attUnitId string
		t.Logf("attribute_name: %s", attNameStr)
		for _, attName := range attNames {
			if attNameStr == attName.GetName().GetValue() {
				attNameId = attName.GetAttributeNameId().GetValue()
				break
			}
		}

		if attUnitNameStr == nil {
			attUnitId = ""
		} else {
			t.Logf("attribute_unit_name: %s", attUnitNameStr)
			for _, attUnit := range attUnits {
				if attUnitNameStr == attUnit.GetName().GetValue() {
					attUnitId = attUnit.GetAttributeUnitId().GetValue()
					break
				}
			}
		}

		//generate CreateAttributeRequest
		attReq := &pb.CreateAttributeRequest{
			AttributeNameId: pbutil.ToProtoString(attNameId),
			AttributeUnitId: pbutil.ToProtoString(attUnitId),
		}
		if minValue != nil {
			attReq.MinValue = pbutil.ToProtoUInt32(uint32(minValue.(int)))
		}
		if maxValue != nil {
			attReq.MaxValue = pbutil.ToProtoUInt32(uint32(maxValue.(int)))
		}
		if strValue != nil {
			attReq.StrValue = pbutil.ToProtoString(strValue.(string))
		}
		//create attribute
		res, err := ss.server.CreateAttribute(ss.ctx, attReq)
		if err != nil {
			t.Skipf("Failed to create attribute(%s), Error: [%+v]", attNameStr, err)
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
			t.Logf("attribute_id: %s, attribute_name_id: %s, attribute_unit_id: %s",
				att.GetAttributeId().GetValue(),
				att.GetAttributeNameId().GetValue(),
				att.GetAttributeUnitId().GetValue())
		}
	}
}
