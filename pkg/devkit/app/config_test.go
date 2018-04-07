// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package app

import (
	"encoding/json"
	"testing"
)

func TestConfigJson_GetDefault(t *testing.T) {
	var configJson = ConfigTemplate{
		Type: TypeArray,
		Properties: []*ConfigTemplate{
			{
				Key:  "cluster",
				Type: TypeArray,
				Properties: []*ConfigTemplate{
					{
						Key:     "name",
						Default: 1,
					},
					{
						Key:     "description",
						Default: 2,
					},
				},
			},
			{
				Key:  "env",
				Type: TypeArray,
				Properties: []*ConfigTemplate{
					{
						Key:     "user",
						Default: "test",
					},
					{
						Key:     "passwd",
						Default: 0.01,
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
}
