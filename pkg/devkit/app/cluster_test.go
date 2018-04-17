// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package app

import (
	"encoding/json"
	"testing"
)

func TestCluster_Render(t *testing.T) {
	tmpl := `
{
	"name": "{{.cluster.name}}"
}
`
	clusterTmpl := ClusterConfTemplate{
		Raw: tmpl,
	}
	configJson := ConfigTemplate{
		Type: TypeArray,
		Properties: []*ConfigTemplate{
			{
				Key:  "cluster",
				Type: TypeArray,
				Properties: []*ConfigTemplate{
					{
						Key:     "name",
						Default: "foobar",
					},
				},
			},
		},
	}
	defaultConfig := configJson.GetDefaultConfig()
	j, err := json.Marshal(&defaultConfig)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(j))
	cluster, err := clusterTmpl.Render(defaultConfig)
	t.Log(cluster.RenderJson)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(cluster.Name, cluster.Description)
}
