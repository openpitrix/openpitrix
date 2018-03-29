// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package devkit

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"openpitrix.io/openpitrix/pkg/devkit/app"
)

const defaultClusterJsonMustache = `
{}
`

const defaultConfigJson = `
{}
`

func Create(metadata *app.Metadata, dir string) (string, error) {
	path, err := filepath.Abs(dir)
	if err != nil {
		return path, err
	}

	if fi, err := os.Stat(path); err != nil {
		return path, err
	} else if !fi.IsDir() {
		return path, fmt.Errorf("no such directory [%s]", path)
	}

	n := metadata.Name
	cdir := filepath.Join(path, n)
	if fi, err := os.Stat(cdir); err == nil && !fi.IsDir() {
		return cdir, fmt.Errorf("file [%s] already exists and is not a directory", cdir)
	}
	if err := os.MkdirAll(cdir, 0755); err != nil {
		return cdir, err
	}

	cf := filepath.Join(cdir, PackageJson)
	if _, err := os.Stat(cf); err != nil {
		if err := savePackageJson(cf, metadata); err != nil {
			return cdir, err
		}
	}

	files := []struct {
		path    string
		content string
	}{
		{
			// cluster.json.mustache
			path:    filepath.Join(cdir, ClusterJsonTemplate),
			content: strings.Replace(defaultClusterJsonMustache, "<APPNAME>", metadata.Name, -1),
		},
		{
			// config.json
			path:    filepath.Join(cdir, ConfigJson),
			content: strings.Replace(defaultConfigJson, "<APPNAME>", metadata.Name, -1),
		},
	}

	for _, file := range files {
		if _, err := os.Stat(file.path); err == nil {
			// File exists and is okay. Skip it.
			continue
		}
		if err := ioutil.WriteFile(file.path, []byte(file.content), 0644); err != nil {
			return cdir, err
		}
	}
	return cdir, nil
}
