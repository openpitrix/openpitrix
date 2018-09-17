// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package runtime

import (
	"reflect"
	"testing"
)

func TestSelectorStringToMap(t *testing.T) {
	selectorStringArray := []string{
		"runtime=qingcloud&zone=pk3a&env=test",
		"runtime=kubernetes&env=dev&env=test",
		"runtime=kubernetes&runtime=qingcloud&team=openpitrix",
	}

	selectorMapArray := []map[string][]string{
		{
			"runtime": []string{"qingcloud"},
			"zone":    []string{"pk3a"},
			"env":     []string{"test"},
		},
		{
			"runtime": []string{"kubernetes"},
			"env":     []string{"dev", "test"},
		},
		{
			"runtime": []string{"kubernetes", "qingcloud"},
			"team":    []string{"openpitrix"},
		},
	}
	for n, selectorString := range selectorStringArray {
		selectorMap, err := SelectorStringToMap(selectorString)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(selectorMapArray[n], selectorMap) {
			t.Fatal()
		}
	}

}
