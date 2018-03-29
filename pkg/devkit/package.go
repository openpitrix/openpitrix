// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package devkit

import (
	"encoding/json"
	"io/ioutil"

	"openpitrix.io/openpitrix/pkg/devkit/app"
)

func UnmarshalPackageJson(data []byte) (*app.Metadata, error) {
	y := &app.Metadata{}
	err := json.Unmarshal(data, y)
	if err != nil {
		return nil, err
	}
	return y, nil
}

func UnmarshalConfigJson(data []byte) (*app.Config, error) {
	y := &app.Config{}
	err := json.Unmarshal(data, y)
	if err != nil {
		return nil, err
	}
	y.Raw = string(data)
	return y, nil
}

func savePackageJson(filename string, metadata *app.Metadata) error {
	out, err := json.MarshalIndent(metadata, "", "    ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, out, 0644)
}
