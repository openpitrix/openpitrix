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
	NAME         = "name"
	DISPLAY_NAME = "displayName"
	ATT_NAME     = "att_name"
	UNIT_NAME    = "unit_name"
	MAX_VALUE    = "max_value"
	MIN_VALUE    = "min_value"
)

const (
	TEST_OFFSET = 0
	TEST_LIMIT  = 30
)

var test_attributes = []map[string]string{
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
	{NAME: "ap2a", DISPLAY_NAME: "亚洲2区-A"},
	{NAME: "pek3", DISPLAY_NAME: "北京3区"},
}

var test_att_values = []map[string]interface{}{
	//时长属性值
	{ATT_NAME: "duration", UNIT_NAME: "hour", MIN_VALUE: 1, MAX_VALUE: 1},
	{ATT_NAME: "duration", UNIT_NAME: "month", MIN_VALUE: 1, MAX_VALUE: 1},
	{ATT_NAME: "duration", UNIT_NAME: "year", MIN_VALUE: 1, MAX_VALUE: 1},
	{ATT_NAME: "duration", UNIT_NAME: "hour", MIN_VALUE: 0, MAX_VALUE: 100},
	{ATT_NAME: "duration", UNIT_NAME: "hour", MIN_VALUE: 100, MAX_VALUE: 1000},
	{ATT_NAME: "duration", UNIT_NAME: "hour", MIN_VALUE: 1000},
	//memory
	{ATT_NAME: "memory", UNIT_NAME: "gb", MIN_VALUE: 1, MAX_VALUE: 1},
	{ATT_NAME: "memory", UNIT_NAME: "gb", MIN_VALUE: 2, MAX_VALUE: 2},
	{ATT_NAME: "memory", UNIT_NAME: "gb", MIN_VALUE: 4, MAX_VALUE: 4},
	{ATT_NAME: "memory", UNIT_NAME: "gb", MIN_VALUE: 8, MAX_VALUE: 8},
	{ATT_NAME: "memory", UNIT_NAME: "gb", MIN_VALUE: 16, MAX_VALUE: 16},
	{ATT_NAME: "memory", UNIT_NAME: "gb", MIN_VALUE: 32, MAX_VALUE: 32},
	//CPU
	{ATT_NAME: "cpu", UNIT_NAME: "number", MIN_VALUE: 1, MAX_VALUE: 1},
	{ATT_NAME: "cpu", UNIT_NAME: "number", MIN_VALUE: 2, MAX_VALUE: 2},
	{ATT_NAME: "cpu", UNIT_NAME: "number", MIN_VALUE: 4, MAX_VALUE: 4},
	{ATT_NAME: "cpu", UNIT_NAME: "number", MIN_VALUE: 8, MAX_VALUE: 8},
	{ATT_NAME: "cpu", UNIT_NAME: "number", MIN_VALUE: 16, MAX_VALUE: 16},
	//disk
	{ATT_NAME: "disk", UNIT_NAME: "gb", MIN_VALUE: 100, MAX_VALUE: 100},
	{ATT_NAME: "disk", UNIT_NAME: "gb", MIN_VALUE: 500, MAX_VALUE: 500},
	{ATT_NAME: "disk", UNIT_NAME: "tb", MIN_VALUE: 1, MAX_VALUE: 1},
	{ATT_NAME: "disk", UNIT_NAME: "tb", MIN_VALUE: 10, MAX_VALUE: 10},
	//region
	{ATT_NAME: "region", UNIT_NAME: "ap2a"},
	{ATT_NAME: "region", UNIT_NAME: "pek3"},
	//instance
	{ATT_NAME: "node_num", UNIT_NAME: "number", MIN_VALUE: 0, MAX_VALUE: 5},
	{ATT_NAME: "node_num", UNIT_NAME: "number", MIN_VALUE: 5, MAX_VALUE: 20},
	{ATT_NAME: "node_num", UNIT_NAME: "number", MIN_VALUE: 20},
	//user
	{ATT_NAME: "user_num", UNIT_NAME: "number", MIN_VALUE: 0, MAX_VALUE: 10},
	{ATT_NAME: "user_num", UNIT_NAME: "number", MIN_VALUE: 10, MAX_VALUE: 50},
	{ATT_NAME: "user_num", UNIT_NAME: "number", MIN_VALUE: 50},
	//stream
	{ATT_NAME: "stream", UNIT_NAME: "mb", MIN_VALUE: 1, MAX_VALUE: 1},
	{ATT_NAME: "stream", UNIT_NAME: "gb", MIN_VALUE: 1, MAX_VALUE: 1},
}

//Attribute
func TestCreateAttribute(t *testing.T) {
	for _, att := range test_attributes {
		t.Run(fmt.Sprintf("CreateAttribute_%s", att[NAME]),
			testCreateAttributeFunc(att[NAME], att[DISPLAY_NAME]))
	}
}

func testCreateAttributeFunc(name, displayName string) func(t *testing.T) {
	return func(t *testing.T) {
		var attReq = &pb.CreateAttributeRequest{
			Name:        pbutil.ToProtoString(name),
			DisplayName: pbutil.ToProtoString(displayName),
		}

		response, err := ss.server.CreateAttribute(ss.ctx, attReq)
		if err != nil {
			t.Error(err)
		} else {
			t.Logf("TestCreateAttribute: Insert attribute(%s) successfully.", response.GetAttributeId().GetValue())
		}
	}
}

func TestListAttribute(t *testing.T) {
	listReq := pb.ListAttributeRequest{
		Offset: TEST_OFFSET,
		Limit:  TEST_LIMIT,
	}

	res, err := ss.server.ListAttribute(ss.ctx, &listReq)
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("List %d Attributes: ", len(res.Attributes))
		for _, att := range res.Attributes {
			t.Logf("%s, %s, %s ",
				att.GetAttributeId().GetValue(),
				att.GetName().GetValue(),
				att.GetDisplayName().GetValue())
		}
	}
}

//Attribute_Unit
func TestCreateAttributeUnit(t *testing.T) {
	for _, attUnit := range test_att_units {
		t.Run(fmt.Sprintf("Create_att_unit_%s", attUnit[NAME]),
			testCreateAttUnitFunc(attUnit[NAME], attUnit[DISPLAY_NAME]))
	}
}

func testCreateAttUnitFunc(name, displayName string) func(t *testing.T) {
	return func(t *testing.T) {
		var attUnitReq = &pb.CreateAttUnitRequest{
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

func TestListAttributeUnit(t *testing.T) {
	listReq := pb.ListAttUnitRequest{
		Offset: TEST_OFFSET,
		Limit:  TEST_LIMIT,
	}

	res, err := ss.server.ListAttributeUnit(ss.ctx, &listReq)
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("List %d Attribute_units: ", len(res.AttributeUnits))
		for _, attUnit := range res.AttributeUnits {
		t.Logf("%s, %s, %s",
			attUnit.GetAttributeUnitId().GetValue(),
			attUnit.GetName().GetValue(),
			attUnit.GetDisplayName().GetValue())
		}
	}
}

//Attribute_Value
//need to run TestCreateAttribute and TestCreateAttributeUnit
func TestCreateAttributeValue(t *testing.T) {
	//get attributes
	listAttReq := pb.ListAttributeRequest{
		Offset: TEST_OFFSET,
		Limit:  TEST_LIMIT,
	}
	atts, err := ss.server.ListAttribute(ss.ctx, &listAttReq)
	if err != nil {
		t.Error(err)
	}

	//get attribute_units
	listAttUnitReq := pb.ListAttUnitRequest{
		Offset: TEST_OFFSET,
		Limit:  TEST_LIMIT,
	}
	attUnits, err := ss.server.ListAttributeUnit(ss.ctx, &listAttUnitReq)
	if err != nil {
		t.Error(err)
	}

	//generate and create AttributeValue
	for _, tAttValue := range test_att_values {
		//get attributes and attributeUnits
		var attId, attUnitId string
		t.Logf("attName: %s", tAttValue[ATT_NAME])
		t.Logf("attUnitName: %s", tAttValue[UNIT_NAME])
		for _, att := range atts.Attributes {
			if tAttValue[ATT_NAME] == att.GetName().GetValue() {
				attId = att.GetAttributeId().GetValue()
				t.Logf("attID: %s", attId)
				break
			}
		}
		for _, attUnit := range attUnits.AttributeUnits {
			if tAttValue[UNIT_NAME] == attUnit.GetName().GetValue() {
				attUnitId = attUnit.GetAttributeUnitId().GetValue()
				t.Logf("attUnitID: %s", attUnitId)
				break
			}
		}

		//generate
		attValueReq := &pb.CreateAttValueRequest{
			AttributeId:     pbutil.ToProtoString(attId),
			AttributeUnitId: pbutil.ToProtoString(attUnitId),
		}
		if tAttValue[MIN_VALUE] != nil {
			attValueReq.MinValue = pbutil.ToProtoUInt32(uint32(tAttValue[MIN_VALUE].(int)))
		}
		if tAttValue[MAX_VALUE] != nil {
			attValueReq.MaxValue = pbutil.ToProtoUInt32(uint32(tAttValue[MAX_VALUE].(int)))
		}
		//create
		valRes, err := ss.server.CreateAttributeValue(ss.ctx, attValueReq)
		if err != nil {
			t.Skip(err)
		} else {
			t.Logf("Create attribute_value(%s) successfully.", valRes.GetAttributeValueId().GetValue())
		}

	}

}
