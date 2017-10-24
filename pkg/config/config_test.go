// +-------------------------------------------------------------------------
// | Copyright (C) 2017 Yunify, Inc.
// +-------------------------------------------------------------------------
// | Licensed under the Apache License, Version 2.0 (the "License");
// | you may not use this work except in compliance with the License.
// | You may obtain a copy of the License in the LICENSE file, or at:
// |
// | http://www.apache.org/licenses/LICENSE-2.0
// |
// | Unless required by applicable law or agreed to in writing, software
// | distributed under the License is distributed on an "AS IS" BASIS,
// | WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// | See the License for the specific language governing permissions and
// | limitations under the License.
// +-------------------------------------------------------------------------

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
		Database: "a",
		Host:     "b",
		Port:     123,
		Protocol: "abc",
		URI:      "/??",
		LogLevel: "warn",
	}

	conf1, err := Parse(conf0.String())
	tAssertf(t, err == nil, "err = %v", err)

	tAssertf(t, reflect.DeepEqual(conf0, conf1), "%v != %v", conf0, conf1)
}

func TestConfig_default(t *testing.T) {
	_, err := Parse(DefaultConfigContent)
	tAssertf(t, err == nil, "err = %v", err)
}
