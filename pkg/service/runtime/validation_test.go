// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package runtime

import (
	"context"
	"testing"
)

var ctx = context.TODO()

func TestValidateURL(t *testing.T) {
	validURLs := []string{
		"http://foo.com/blah_blah",
		"http://userid:password@example.com:8080",
		"http://➡.ws/䨹",
		"http://例子.测试",
		"http://مثال.إختبار",
	}
	invalidURLs := []string{
		"http://??",
		"http://foo.bar?q=Spaces should be encoded",
		"//",
		"rdar://1234",
		"http://224.1.1.1",
	}
	for _, validURL := range validURLs {
		err := ValidateURL(ctx, validURL)
		if err != nil {
			t.Fatalf("%+v should be validURL", validURL)
		}
	}
	for _, invalidURL := range invalidURLs {
		err := ValidateURL(ctx, invalidURL)
		if err == nil {
			t.Fatalf("%+v should be validURL", invalidURL)
		}
	}
}
