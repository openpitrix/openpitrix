// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// +build etcd

package etcd_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"openpitrix.io/openpitrix/pkg/config/test_config"
	"openpitrix.io/openpitrix/pkg/etcd"
)

var tc = test_config.NewEtcdTestConfig()

func TestEtcdQueue(t *testing.T) {
	tc.CheckEtcdUnitTest(t)
	e, err := etcd.Connect(tc.GetTestEtcdEndpoints(), "test")
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

func TestDlockWithTimeout(t *testing.T) {
	tc.CheckEtcdUnitTest(t)
	e, err := etcd.Connect(tc.GetTestEtcdEndpoints(), "test")
	if err != nil {
		t.Fatal(err)
	}
	dlockKey := "test-dlock"

	t.Logf("start first dlock [%s]", time.Now())
	err = e.DlockWithTimeout(dlockKey, 100*time.Second, func() error {

		go func() {
			t.Logf("start second dlock [%s]", time.Now())
			err := e.DlockWithTimeout(dlockKey, 3*time.Second, func() error {
				time.Sleep(2 * time.Second)
				t.Logf("sleep 2 second finish [%s]", time.Now())
				return nil
			})
			if err == nil {
				t.Fatalf("second dlock should fail [%s]", time.Now())
			} else {
				t.Logf("second dlock fail as except [%s]", time.Now())
			}
		}()

		time.Sleep(10 * time.Second)
		t.Logf("sleep 10 second finish [%s]", time.Now())
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

}
