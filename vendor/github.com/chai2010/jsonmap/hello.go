// Copyright 2018 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"fmt"
	"sort"

	"github.com/chai2010/jsonmap"
)

func main() {
	var jsonMap = jsonmap.JsonMap{
		"a": map[string]interface{}{
			"sub-a": "value-sub-a",
		},
		"b": map[string]interface{}{
			"sub-b": "value-sub-b",
		},
		"c": 123,
		"d": 3.14,
		"e": true,

		"x": map[string]interface{}{
			"a": map[string]interface{}{
				"sub-a": "value-sub-a",
			},
			"b": map[string]interface{}{
				"sub-b": "value-sub-b",
			},
			"c": 123,
			"d": 3.14,
			"e": true,

			"z": map[string]interface{}{
				"zz": "value-zz",
			},
		},
	}

	m := jsonMap.ToMapString("/")

	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		fmt.Println(k, m[k])
	}

	// Output:
	// /a/sub-a value-sub-a
	// /b/sub-b value-sub-b
	// /c 123
	// /d 3.14
	// /e true
	// /x/a/sub-a value-sub-a
	// /x/b/sub-b value-sub-b
	// /x/c 123
	// /x/d 3.14
	// /x/e true
	// /x/z/zz value-zz
}
