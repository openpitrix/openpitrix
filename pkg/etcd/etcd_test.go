// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package etcd_test

import (
	"fmt"
	"math/rand"
	"testing"

	"openpitrix.io/openpitrix/pkg/etcd"
)

func TestEtcd(t *testing.T) {
	e, err := etcd.Connect([]string{"localhost:2379"}, "test")
	if err != nil {
		t.Fatal(err)
	}
	queue := e.NewQueue(fmt.Sprintf("test-queue-%d", rand.Intn(10000)))
	go func() {
		for i := 0; i < 100; i++ {
			err := queue.Enqueue(fmt.Sprintf("%d", i))
			if err != nil {
				t.Fatal(err)
			}
			t.Logf("Push message to queue, worker number [%d]", i)
		}

	}()
	for i := 0; i < 100; i++ {
		n, err := queue.Dequeue()
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("Got message [%s] from queue, worker number [%d]", n, i)
	}
}
