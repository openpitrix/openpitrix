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

	"openpitrix.io/openpitrix/pkg/devkit/opapp"
)

const defaultClusterJsonTmpl = `
{
    "name": "{{.cluster.name}}",
    "description": "{{.cluster.description}}",
    "subnet": "{{.cluster.subnet}}",
    "nodes": [{
        "role": "role_name1",
        "container": {
            "type": "docker",
            "image": "nginx"
        },
        "count": "{{.cluster.role_name1.count}}",
        "cpu": "{{.cluster.role_name1.cpu}}",
        "memory": "{{.cluster.role_name1.memory}}",
        "volume": {
            "size": "{{.cluster.role_name1.volume_size}}",
            "mount_point": "/test_data",
            "filesystem": "ext4"
        }
    }]
}
`

const defaultConfigJson = `
{
    "type": "array",
    "properties": [{
        "key": "cluster",
        "description": "<APPNAME> cluster properties",
        "type": "array",
        "properties": [{
            "key": "name",
            "label": "name",
            "description": "The name of the <APPNAME> service",
            "type": "string",
            "default": "<APPNAME>",
            "required": false
        }, {
            "key": "description",
            "label": "description",
            "description": "The description of the <APPNAME> service",
            "type": "string",
            "default": "",
            "required": false
        }, {
            "key": "subnet",
            "label": "Subnet",
            "description": "Choose a subnet to join",
            "type": "string",
            "default": "",
            "required": true
        }, {
            "key": "role_name1",
            "label": "role_name1",
            "description": "role-based role_name1 properties",
            "type": "array",
            "properties": [{
                "key": "cpu",
                "label": "CPU",
                "description": "CPUs of each node",
                "type": "integer",
                "default": 1,
                "range": [1, 2, 4, 8, 16],
                "required": true
            }, {
                "key": "memory",
                "label": "Memory",
                "description": "Memory of each node",
                "type": "integer",
                "default": 2048,
                "range": [2048, 8192, 16384, 32768, 49152],
                "required": true
            }, {
                "key": "count",
                "label": "Count",
                "description": "Number of nodes for the cluster to create",
                "type": "integer",
                "default": 3,
                "max": 100,
                "min": 1,
                "required": true
            }, {
                "key": "volume_size",
                "label": "Volume Size",
                "description": "The volume size for each instance",
                "type": "integer",
                "default": 10,
                "min": 10,
                "max": 1000,
                "step": 10,
                "required": true
            }]
        }]
    }]
}
`

func Create(metadata *opapp.Metadata, dir string) (string, error) {
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
			return cdir, fmt.Errorf("failed to write [%s], error: %+v", PackageJson, err)
		}
	}

	files := []struct {
		path    string
		content string
	}{
		{
			// cluster.json.mustache
			path:    filepath.Join(cdir, ClusterJsonTmpl),
			content: strings.Replace(defaultClusterJsonTmpl, "<APPNAME>", metadata.Name, -1),
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
			return cdir, fmt.Errorf("failed to write [%s], error: %+v", file.path, err)
		}
	}
	return cdir, nil
}
