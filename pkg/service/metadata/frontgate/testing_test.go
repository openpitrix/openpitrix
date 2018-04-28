// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// go test -test-etcd-enabled

package frontgate

import (
	"flag"
	"fmt"
	"os"
	"testing"
)

var (
	tEtcdEnabled = flag.Bool("test-etcd-enabled", false, "enable etcd server")
	tEtcdHost    = flag.String("test-ectd-host", "localhost:2379", "set etcd nodes")
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func Assert(tb testing.TB, condition bool, a ...interface{}) {
	tb.Helper()
	if !condition {
		if msg := fmt.Sprint(a...); msg != "" {
			tb.Fatal("Assert failed: " + msg)
		} else {
			tb.Fatal("Assert failed")
		}
	}
}

func Assertf(tb testing.TB, condition bool, format string, a ...interface{}) {
	tb.Helper()
	if !condition {
		if msg := fmt.Sprintf(format, a...); msg != "" {
			tb.Fatal("Assert failed: " + msg)
		} else {
			tb.Fatal("Assert failed")
		}
	}
}
