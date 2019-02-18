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
	NAME = "name"
	DISPLAY_NAME = "displayName"
)

const (
	TEST_PAGE = 0
	TEST_PAGE_SIZE = 20
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
	{NAME: "num", DISPLAY_NAME: "个数"},
	{NAME: "mb", DISPLAY_NAME: "MB"},
	{NAME: "gb", DISPLAY_NAME: "GB"},
	{NAME: "tb", DISPLAY_NAME: "TB"},
	{NAME: "ap2a", DISPLAY_NAME: "亚洲2区-A"},
	{NAME: "pek3", DISPLAY_NAME: "北京3区"},
}

func TestCreateAttribute(t *testing.T) {
	for _, att := range test_attributes {
		t.Run(fmt.Sprintf("CreateAttribute_%s",att[NAME]),
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
	commonListReq := pb.CommonListRequest{
		Page: 		pbutil.ToProtoUInt32(TEST_PAGE),
		PageSize: 	pbutil.ToProtoUInt32(TEST_PAGE_SIZE),
	}

	res, err := ss.server.ListAttribute(ss.ctx, &commonListReq)
	if err != nil {
		t.Error(err)
	} else {
		attIds := ""
		for _, att := range res.Attributes {
			attIds = attIds + att.GetAttributeId().GetValue() + ", "
		}
		t.Logf("ListAttributes(%d): %s", len(attIds), attIds)
	}
}

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
