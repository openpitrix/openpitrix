// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package qingcloud

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
$ cat ~/.qingcloud/config.yaml

# access key info
qy_access_key_id: 'CDUSSVOZCVYCVVQDFYUE'
qy_secret_access_key: 'TGltimaj5i7DspZE5i0QBKljMpM87IRvKX6SLIAa'

# message expired time in seconds
msg_time_out: 3600

# http socket timeout in seconds
http_socket_timeout: 30

# retry time after message send failed
retry_time: 3

# remote host
host: "api.alphacloud.com"
#host: "webservice0"

# remote host port, 443 for https
port: 7777
#port: 8882

# protocol, "http" or "https"
protocol: "http"

zone: "test"
*/

const test_appConf = `
{
  "app_id": "app-ZooKeeper81S8qik7spNTE39hAtZEmtjhUU5",
  "app_version": "1.0",
  "vxnet": "vxnet-djxab0n",
  "debug": true,
  "nodes": {
     "container": {
        "type": "kvm",
        "image": "img-k97x35rx",
        "zone": "test"
     },
     "instance_class": 0,
     "count": 3,
     "cpu": 1,
     "vertical_scaling_policy": "sequential",
     "memory": 512,
     "volume": {
         "size": 3,
         "mount_point": "/zk_data",
         "filesystem": "ext4",
         "class": 0
     },
     "server_id_upper_bound":255,
     "services": {
        "start": {
           "cmd": "/opt/zookeeper/bin/zkServer.sh start"
        },
        "stop": {
           "cmd": "/opt/zookeeper/bin/zkServer.sh stop"
        }
     }
  },
  "advanced_actions": ["change_vxnet", "scale_horizontal"],
  "endpoints": {
    "client_port": {
      "port": 2181,
      "protocol": "tcp"
    }
  },
  "health_check": {
    "enable": true,
    "interval_sec": 60,
    "timeout_sec": 10,
    "action_timeout_sec": 10,
    "healthy_threshold": 3,
    "unhealthy_threshold": 3,
    "check_cmd": "/opt/bin/check.sh",
    "action_cmd": "/opt/zookeeper/bin/zkServer.sh start"
  },
  "monitor": {
    "enable": true,
    "cmd": "/opt/bin/monitor.sh",
    "items": {
      "mode": {
        "value_type": "str",
        "statistics_type": "mode"
      },
      "received": {
        "unit": "",
        "value_type": "int",
        "statistics_type": "latest"
      },
      "sent": {
        "statistics_type": "latest"
      },
      "latency_min": {
        "statistics_type": "min",
        "scale_factor_when_display": 1
      },
      "latency_avg": {
      },
      "latency_max": {
        "statistics_type": "max"
      },
      "connections": {
        "statistics_type": "median"
      },
      "outstanding": {
        "statistics_type": "median"
      },
      "node_count": {
        "statistics_type": "mode"
      }
    },
    "groups": {
      "Throughput": ["sent", "received"],
      "ConnectionNum": ["connections", "outstanding"],
      "Latency": ["latency_min", "latency_avg", "latency_max"]
    },
    "display": ["ConnectionNum", "Throughput", "Latency", "mode", "node_count"]
  }
}
`

func TestQingCloudRuntime(t *testing.T) {
	clientConf := "~/.qingcloud/config.yaml"
	_, err := os.Stat(strings.Replace(clientConf, "~/", os.Getenv("HOME")+"/", 1))
	if err != nil {
		t.Skipf("QingCloud runtime test skipped because no [%s], err: %v", clientConf, err)
	}

	runtime := QingcloudRuntime{}

	clusterId, err := runtime.CreateCluster(test_appConf, true)
	assert.Empty(t, err)
	err = runtime.StopClusters(clusterId, true, "0")
	assert.Empty(t, err)
	err = runtime.StartClusters(clusterId, true)
	assert.Empty(t, err)
	err = runtime.DeleteClusters(clusterId, true, "1")
	assert.Empty(t, err)
	err = runtime.RecoverClusters(clusterId, true)
	assert.Empty(t, err)
	err = runtime.DeleteClusters(clusterId, true, "1")
	assert.Empty(t, err)
	err = runtime.CeaseClusters(clusterId, true)
	assert.Empty(t, err)
	_, err = runtime.DescribeClusters(clusterId)
	_, err = runtime.DescribeClusterNodes(clusterId)
}
