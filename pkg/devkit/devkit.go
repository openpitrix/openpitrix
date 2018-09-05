// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package devkit

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"openpitrix.io/openpitrix/pkg/devkit/opapp"
)

const (
	ApiVersionV1    = "v1"
	PackageJson     = "package.json"
	ClusterJsonTmpl = "cluster.json.tmpl"
	ConfigJson      = "config.json"
)

// IsAppDir validate a app directory.
//
// Checks for a valid package.json.
func IsAppDir(dirName string) (bool, error) {
	if fi, err := os.Stat(dirName); err != nil {
		return false, err
	} else if !fi.IsDir() {
		return false, fmt.Errorf("[%s] is not a directory", dirName)
	}

	packageJson := filepath.Join(dirName, PackageJson)
	if _, err := os.Stat(packageJson); os.IsNotExist(err) {
		return false, fmt.Errorf("no %s exists in directory [%s]", PackageJson, dirName)
	}

	packageJsonContent, err := ioutil.ReadFile(packageJson)
	if err != nil {
		return false, fmt.Errorf("cannot read %s in directory [%s]", PackageJson, dirName)
	}

	packageContent, err := opapp.DecodePackageJson(packageJsonContent)
	if err != nil {
		return false, err
	}
	if packageContent == nil {
		return false, fmt.Errorf("app metadata [%s] missing", PackageJson)
	}
	if packageContent.Name == "" {
		return false, fmt.Errorf("invalid app [%s]: name must not be empty", PackageJson)
	}

	return true, nil
}
