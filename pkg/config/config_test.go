// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package config

import (
	"fmt"
	"reflect"
	"testing"
)

func tAssert(tb testing.TB, condition bool) {
	tb.Helper()
	if !condition {
		tb.Fatal("Assert failed")
	}
}

func tAssertf(tb testing.TB, condition bool, format string, a ...interface{}) {
	tb.Helper()
	if !condition {
		if msg := fmt.Sprintf(format, a...); msg != "" {
			tb.Fatal("Assert failed: " + msg)
		} else {
			tb.Fatal("Assert failed")
		}
	}
}

func TestConfig(t *testing.T) {
	conf0 := &Config{
		DbType:     "mysql",
		DbHost:     "root:password@tcp(127.0.0.1:3306)/openpitrix",
		DbEncoding: "utf8",
		DbEngine:   "InnoDB",
		Host:       "127.0.0.1",
		Port:       8443,
		Protocol:   "https",
		URI:        "/openpitrix/api/v1",
		LogLevel:   "warn",
	}

	conf1, err := Parse(conf0.String())
	tAssertf(t, err == nil, "err = %v", err)

	tAssertf(t, reflect.DeepEqual(conf0, conf1), "%v != %v", conf0, conf1)
}

func TestConfig_default(t *testing.T) {
	_, err := Parse(DefaultConfigContent)
	tAssertf(t, err == nil, "err = %v", err)
}
