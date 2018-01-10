// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package utils

import (
	"encoding/json"
	"io/ioutil"

	"openpitrix.io/openpitrix/pkg/logger"
)

func Loads(src string, dest interface{}) error {
	err := json.Unmarshal([]byte(src), &dest)
	if err != nil {
		logger.Errorf("Failed to load json: %v", err)
	}
	return err
}

func Load(filename string, dest interface{}) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		logger.Errorf("Failed to load json from file [%s]: %v", filename, err)
		return err
	}
	return Loads(string(data), &dest)
}

func Dumps(src interface{}) (string, error) {
	dest, err := json.MarshalIndent(src, "", "  ")
	if err != nil {
		logger.Errorf("Failed to dump json: %v", err)
	}
	return string(dest), err
}

func Dump(filename string, src interface{}) error {
	dest, err := Dumps(src)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filename, []byte(dest), 0644)
	if err != nil {
		logger.Errorf("Failed to load json to file [%s]: %v", filename, err)
		return err
	}
	return err
}
