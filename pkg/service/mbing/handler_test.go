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

var attributes = []map[string]string{
	{NAME: "node_num", DISPLAY_NAME: "实例数"},
	{NAME: "memory", DISPLAY_NAME: "内存"},
	{NAME: "cpu", DISPLAY_NAME: "CPU"},
	{NAME: "disk", DISPLAY_NAME: "硬盘"},
	{NAME: "region", DISPLAY_NAME: "区域"},
	{NAME: "user_num", DISPLAY_NAME: "用户数"},
	{NAME: "stream", DISPLAY_NAME: "流量"},
}

func TestCreateAttribute(t *testing.T) {
	for _, att := range attributes {
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
