// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package mbing

import (
	"context"
	"testing"
	"time"

	"openpitrix.io/openpitrix/pkg/logger"
	mbing "openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
)

func common(t *testing.T) (*Server, context.Context, context.CancelFunc) {
	if !*tTestingEnvEnabled {
		t.Skip("testing env disabled")
	}
	InitGlobelSetting()
	s, _ := NewServer()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	return s, ctx, cancel
}

func TestNewServer(t *testing.T) {
	s, ctx, cancle := common(t)
	t.Logf("TestNewServer Passed, server: %v, context: %v, cancleFunc: %v", s, ctx, cancle)
}

func TestStartMetering(t *testing.T) {
	s, ctx, cancel := common(t)
	defer cancel()

	var resourceList []*mbing.ResourceVersion
	for i := 0; i < 3; i++ {
		resourceList = append(resourceList, &mbing.ResourceVersion{
			ResourceVersionId: pbutil.ToProtoString("testResourceVersionId" + string(i)),
			PriceId:           pbutil.ToProtoString("PriceId" + string(i)),
			ActionTime:        pbutil.ToProtoTimestamp(time.Now()),
		})
	}

	var req = &mbing.MeteringRequest{
		ResourceId:         pbutil.ToProtoString("testResourceID"),
		UserId:             pbutil.ToProtoString("testUserID"),
		ActionResourceList: resourceList,
	}
	resp, _ := s.StartMetering(ctx, req)
	logger.Info(nil, "Test Passed, StartMetering status %s, message %s", resp.GetStatus(), resp.GetMessage())

}
