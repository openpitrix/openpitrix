// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package opapp

import (
	"encoding/json"
	"reflect"
	"testing"

	"openpitrix.io/openpitrix/pkg/util/jsonutil"
)

func TestConfigJson_GetDefault(t *testing.T) {
	var configJson = ConfigTemplate{
		Config: Config{
			Type: TypeArray,
			Properties: []*Config{
				{
					Key:  "cluster",
					Type: TypeArray,
					Properties: []*Config{
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
					Properties: []*Config{
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
		},
	}
	defaultConfig := configJson.GetDefaultConfig()
	j, err := json.Marshal(&defaultConfig)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(j))
}

func TestConfigJson_FillInDefault(t *testing.T) {
	var configJson = ConfigTemplate{
		Config: Config{
			Type: TypeArray,
			Properties: []*Config{
				{
					Key:  "cluster",
					Type: TypeArray,
					Properties: []*Config{
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
					Properties: []*Config{
						{
							Key:  "master",
							Type: TypeArray,
							Properties: []*Config{
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
				},
			},
		},
	}
	var changedConfigJson = ConfigTemplate{
		Config: Config{
			Type: TypeArray,
			Properties: []*Config{
				{
					Key:  "cluster",
					Type: TypeArray,
					Properties: []*Config{
						{
							Key:     "name",
							Default: 11,
						},
						{
							Key:     "description",
							Default: 22,
						},
					},
				},
				{
					Key:  "env",
					Type: TypeArray,
					Properties: []*Config{
						{
							Key:  "master",
							Type: TypeArray,
							Properties: []*Config{
								{
									Key:     "user",
									Default: "test1",
								},
								{
									Key:     "passwd",
									Default: 0.011,
								},
							},
						},
					},
				},
			},
		},
	}
	changedDefaultConfig := changedConfigJson.GetDefaultConfig()
	i, err := json.Marshal(&changedDefaultConfig)
	if err != nil {
		t.Fatal(err)
	}
	configJson.FillInDefaultConfig(changedDefaultConfig)
	defaultConfig := configJson.GetDefaultConfig()
	j, err := json.Marshal(&defaultConfig)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(i, j) {
		t.Logf("%s", i)
		t.Logf("%s", j)
		t.Errorf("Filled config should be equal to changed config.")
	}
}

func TestConfigJson_SpecificConfig(t *testing.T) {
	var configJson = ConfigTemplate{
		Config: Config{
			Type: TypeArray,
			Properties: []*Config{
				{
					Key:  "cluster",
					Type: TypeArray,
					Properties: []*Config{
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
					Properties: []*Config{
						{
							Key:  "master",
							Type: TypeArray,
							Properties: []*Config{
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
				},
			},
		},
	}

	envConfig := map[string]interface{}{
		"env": map[string]interface{}{
			"master": map[string]interface{}{
				"passwd": 0.01,
				"user":   "test",
			},
		},
	}
	i, err := json.Marshal(&envConfig)
	if err != nil {
		t.Fatal(err)
	}
	configJson.SpecificConfig("env")

	specificConfig := configJson.GetDefaultConfig()
	j, err := json.Marshal(&specificConfig)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(i, j) {
		t.Errorf("Specific config should be equal to env config.")
	}
}

func validateConfig(configJson, conf string) error {
	ct, err := DecodeConfigJson([]byte(configJson))
	if err != nil {
		return err
	}
	c, err := jsonutil.NewJson([]byte(conf))
	if err != nil {
		return err
	}
	return ct.Validate(c)
}

func TestValidateConfig(t *testing.T) {
	type args struct {
		configJson string
		conf       string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"empty case",
			args{`{
}`, `{
}`},
			false,
		},
		{
			"success case",
			args{testConfigJson, `{
	"cluster": {
		"subnet": "vxnet-0"
	}
}`},
			false,
		},
		{
			"step error, will be failed",
			args{testConfigJson, `{
	"cluster": {
		"subnet": "vxnet-0",
		"role_name1": {
			"volume_size": 99
		}
	}
}`},
			true,
		},
		{
			"max error, will be failed",
			args{testConfigJson, `{
	"cluster": {
		"subnet": "vxnet-0",
		"role_name1": {
			"count": 200
		}
	}
}`},
			true,
		},
		{
			"type error, will be failed",
			args{testConfigJson, `{
	"cluster": {
		"subnet": "vxnet-0",
		"role_name1": {
			"count": "200"
		}
	}
}`},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateConfig(tt.args.configJson, tt.args.conf); err != nil {
				if !tt.wantErr {
					t.Errorf("validateConfig() error = %v, wantErr %v", err, tt.wantErr)
				} else {
					t.Logf("expect error: %v", err)
				}
			}
		})
	}
}
