// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package runtime

import (
	"reflect"
	"testing"
)

func TestLabelStringToMap(t *testing.T) {
	labelStringArray := []string{
		"runtime=qingcloud&zone=pk3a&env=test",
		"runtime=kubernetes&env=dev",
		"runtime=kubernetes&team=openpitrix",
	}
	labelMapArray := []map[string]string{
		{
			"runtime": "qingcloud",
			"zone":    "pk3a",
			"env":     "test",
		},
		{
			"runtime": "kubernetes",
			"env":     "dev",
		},
		{
			"runtime": "kubernetes",
			"team":    "openpitrix",
		},
	}
	for n, lableString := range labelStringArray {
		labelMap, err := LabelStringToMap(lableString)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(labelMapArray[n], labelMap) {
			t.Fatal()
		}
	}
}

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

func TestLabelMapDiff(t *testing.T) {
	oldLabelMapArray := []map[string]string{
		{
			"runtime": "qingcloud",
			"zone":    "pk3a",
			"env":     "test",
		},
		{
			"runtime": "kubernetes",
			"env":     "dev",
		},
		{
			"runtime": "kubernetes",
			"team":    "openpitrix",
		},
	}
	newLabelMapArray := []map[string]string{
		{
			"runtime": "qingcloud",
			"zone":    "pk3a",
		},
		{
			"runtime": "kubernetes",
			"env":     "dev",
			"team":    "openpitrix",
		},
		{
			"runtime": "kubernetes",
			"env":     "dev",
		},
	}
	additionsArray := []map[string]string{
		{},
		{
			"team": "openpitrix",
		},
		{
			"env": "dev",
		},
	}
	deletionsArray := []map[string]string{
		{
			"env": "test",
		},
		{},
		{
			"team": "openpitrix",
		},
	}

	for n := range oldLabelMapArray {
		additions, deletions := LabelMapDiff(oldLabelMapArray[n], newLabelMapArray[n])
		if !reflect.DeepEqual(additions, additionsArray[n]) || !reflect.DeepEqual(deletions, deletionsArray[n]) {
			t.Fatal()
		}
	}
}
