// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// +build etcd

package topic

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"openpitrix.io/openpitrix/pkg/config/test_config"
	"openpitrix.io/openpitrix/pkg/etcd"
)

const (
	testUid     = "uid"
	testEventId = uint64(1111111)
)

func TestFormatTopic(t *testing.T) {
	topic := formatTopic(testUid, testEventId)
	uid, eventId := parseTopic(topic)
	require.Equal(t, testUid, uid)
	require.Equal(t, testEventId, eventId)
}

var tc = test_config.NewEtcdTestConfig()

func TestWatchEvents(t *testing.T) {
	tc.CheckEtcdUnitTest(t)
	e, err := etcd.Connect(tc.GetTestEtcdEndpoints(), "test")

	require.NoError(t, err)

	c := watchEvents(e)

	time.Sleep(2 * time.Second)

	err = pushEvent(context.Background(), e, testUid, Message{
		Type: Create,
	})
	require.NoError(t, err)

	event := <-c
	require.Equal(t, Create, event.Message.Type)
	require.Equal(t, testUid, event.UserId)
	t.Log(event)
}
