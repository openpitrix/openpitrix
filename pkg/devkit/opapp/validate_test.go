// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package opapp

import (
	"testing"
)

var testConfigJson = `
{
    "type": "array",
    "properties": [{
        "key": "cluster",
        "description": "Sample cluster properties",
        "type": "array",
        "properties": [{
            "key": "name",
            "label": "name",
            "description": "The name of the Sample service",
            "type": "string",
            "default": "Sample",
            "required": false
        }, {
            "key": "description",
            "label": "description",
            "description": "The description of the Sample service",
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

var testErrorClusterTmpl = `
{
    "name": "{{.cluster.name}}",
    "description": "{{.cluster.description}}",
    "subnet": "{{.cluster.subnet}}",
    "nodes": [{
        "role": "role_name1",
        "role2": "role_name1",
        "container": {
            "type": "kvm",
            "zone": "pek3a",
            "image": "img-hlhql5ea"
        },
        "count": "{{.cluster.role_name1.count}}"9999,
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

var testClusterTmpl = `
{
    "name": "{{.cluster.name}}",
    "description": "{{.cluster.description}}",
    "subnet": "{{.cluster.subnet}}",
    "nodes": [{
        "role": "role_name1",
        "container": {
            "type": "kvm",
            "zone": "pek3a",
            "image": "img-hlhql5ea"
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

func TestValidateClusterTmpl(t *testing.T) {
	// normal tmpl
	clusterTmpl := &ClusterConfTemplate{Raw: testClusterTmpl}
	configJson, err := DecodeConfigJson([]byte(testConfigJson))
	if err != nil {
		t.Fatal(err)
	}
	config := configJson.GetDefaultConfig()
	t.Log(config.Interface())
	err = ValidateClusterConfTmpl(clusterTmpl, config)
	if err != nil {
		t.Fatal(err)
	}

	// error tmpl
	clusterTmpl = &ClusterConfTemplate{Raw: testErrorClusterTmpl}
	err = ValidateClusterConfTmpl(clusterTmpl, config)
	if err == nil {
		t.Fatal("error cluster tmpl must failed")
	}
	t.Log(err)
}
