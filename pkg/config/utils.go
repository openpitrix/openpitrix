// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package config

import (
	"os"
	"runtime"

	"gopkg.in/yaml.v2"
)

func yamlEncode(v interface{}) ([]byte, error) {
	return yaml.Marshal(v)
}

func yamlDecode(s []byte, v interface{}) error {
	return yaml.Unmarshal(s, v)
}

func GetHomePath() string {
	home := os.Getenv("HOME")
	if runtime.GOOS == "windows" {
		home = os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
	}
	if home == "" {
		home = "~"
	}

	return home
}
