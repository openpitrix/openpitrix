// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package utils

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

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

type Cluster struct {
	Cluster_id  string `json:"cluster_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	App_id      string `json:"app_id"`
	App_version string `json:"app_version"`
}

var cluster = Cluster{
	Cluster_id:  "cl-12345678",
	Name:        "cluster",
	Description: "new cluster",
	App_id:      "app-12345678",
	App_version: "appv-12345678",
}

var map_cluster = map[string]string{
	"cluster_id":  "cl-12345678",
	"name":        "cluster",
	"description": "new cluster",
	"app_id":      "app-12345678",
	"app_version": "appv-12345678",
}

var str_cluster = `{
	"cluster_id": "cl-12345678",
	"name": "cluster",
	"description": "new cluster",
	"app_id": "app-12345678",
	"app_version": "appv-12345678"
}
`

func TestLoadToStruct(t *testing.T) {
	var c Cluster
	err := Loads(str_cluster, &c)
	tAssertf(t, err == nil, "Exist Error")
	tAssertf(t, reflect.DeepEqual(c, cluster), "%v != %v", c, cluster)
}

func TestLoadToMap(t *testing.T) {
	m := make(map[string]string)
	err := Loads(str_cluster, &m)
	tAssertf(t, err == nil, "Exist Error")
	tAssertf(t, reflect.DeepEqual(m, map_cluster), "%v != %v", m, map_cluster)
}

func TestDump(t *testing.T) {
	dump, err := Dumps(cluster)
	tAssertf(t, err == nil, "Exist Error")

	var c Cluster
	err = Loads(dump, &c)
	tAssertf(t, err == nil, "Exist Error")
	tAssertf(t, reflect.DeepEqual(c, cluster), "%v != %v", c, cluster)
}

func TestLoadAndDump(t *testing.T) {
	filename := "/tmp/load_and_dump.test"
	err := Dump(filename, cluster)
	tAssertf(t, err == nil, "Exist Error")

	var c Cluster
	err = Load(filename, &c)
	defer os.Remove(filename)
	tAssertf(t, err == nil, "Exist Error")
	tAssertf(t, reflect.DeepEqual(c, cluster), "%v != %v", c, cluster)
}
