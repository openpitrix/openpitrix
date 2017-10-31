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

func TestOpenPitrix(t *testing.T) {
	conf0 := &Config{
		OpenPitrix: OpenPitrix{
			Host:     "localhost",
			Port:     8443,
			LogLevel: "warn",
		},
	}

	conf1, err := Parse(conf0.String())
	tAssertf(t, err == nil, "err = %v", err)

	tAssertf(t, reflect.DeepEqual(conf0, conf1), "%v != %v", conf0, conf1)
}

func TestOpenPitrix_default(t *testing.T) {
	conf := Default()

	tAssert(t, conf.Host == "127.0.0.1")
	tAssert(t, conf.Port == 8080)
	tAssert(t, conf.LogLevel == "warn")

	tAssert(t, conf.Database.Type == "mysql")
	tAssert(t, conf.Database.Host == "root:password@tcp(127.0.0.1:3306)")
	tAssert(t, conf.Database.Encoding == "utf8")
	tAssert(t, conf.Database.Engine == "InnoDB")
	tAssert(t, conf.Database.DbName == "openpitrix")
}

func TestOpenPitrix_Parse_default(t *testing.T) {
	conf, err := Parse(DefaultConfigContent)
	tAssertf(t, err == nil, "err = %v", err)

	tAssert(t, conf.Host == "127.0.0.1")
	tAssert(t, conf.Port == 8080)
	tAssert(t, conf.LogLevel == "warn")

	tAssert(t, conf.Database.Type == "mysql")
	tAssert(t, conf.Database.Host == "root:password@tcp(127.0.0.1:3306)")
	tAssert(t, conf.Database.Encoding == "utf8")
	tAssert(t, conf.Database.Engine == "InnoDB")
	tAssert(t, conf.Database.DbName == "openpitrix")
}

func TestOpenPitrix_Parse_empty(t *testing.T) {
	conf, err := Parse(``)

	tAssertf(t, err == nil, "err = %v", err)

	tAssert(t, conf.Host == "127.0.0.1")
	tAssert(t, conf.Port == 8080)
	tAssert(t, conf.LogLevel == "warn")

	tAssert(t, conf.Database.Type == "mysql")
	tAssert(t, conf.Database.Host == "root:password@tcp(127.0.0.1:3306)")
	tAssert(t, conf.Database.Encoding == "utf8")
	tAssert(t, conf.Database.Engine == "InnoDB")
	tAssert(t, conf.Database.DbName == "openpitrix")
}

func TestOpenPitrix_Parse(t *testing.T) {
	conf, err := Parse(`
		Host = "localhost"
		Port = 9527
		
		# Valid log levels are "debug", "info", "warn", "error", and "fatal".
		LogLevel = "debug"
		
		[Database]
		Type     = "pq"
		Host     = "user:password@tcp(127.0.0.1:3306)"
		Encoding = "utf8"
		Engine   = "InnoDB"
		DbName   = "openpitrix-dev"
	`)

	tAssertf(t, err == nil, "err = %v", err)

	tAssert(t, conf.Host == "localhost")
	tAssert(t, conf.Port == 9527)
	tAssert(t, conf.LogLevel == "debug")

	tAssert(t, conf.Database.Type == "pq")
	tAssert(t, conf.Database.Host == "user:password@tcp(127.0.0.1:3306)")
	tAssert(t, conf.Database.Encoding == "utf8")
	tAssert(t, conf.Database.Engine == "InnoDB")
	tAssert(t, conf.Database.DbName == "openpitrix-dev")
}
