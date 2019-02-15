// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package mbing

import (
	"testing"

	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
)

func TestServer_CreateAttribute(t *testing.T) {
	var attReq = &pb.CreateAttributeRequest{
		Name:        pbutil.ToProtoString("duration"),
		DisplayName: pbutil.ToProtoString("时长"),
	}

	response, err := ss.server.CreateAttribute(ss.ctx, attReq)
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("TestCreateAttribute Passed, attribute_id: %s", response.GetAttributeId().GetValue())
	}
}
