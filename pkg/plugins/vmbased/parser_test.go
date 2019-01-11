// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package vmbased

import (
	"testing"

	"openpitrix.io/openpitrix/pkg/devkit/opapp"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/util/jsonutil"
)

const hbaseMustache = `
{
    "app_id": "app-12345678",
    "version_id": "appv-12345678",
    "debug": true,
    "name": "MyHBase",
    "description": "my hbase OpApp",
    "subnet": "vxnet-t8szyjn",
    "links": {
        "zk_service": "cl-w8qh5hf6"
    },
    "nodes": [{
        "role": "hbase-master",
        "container": {
            "type": "docker",
            "image": "img-srf8tv31"
        },
        "count": 1,
        "cpu": 1,
        "memory": 1024,
        "vertical_scaling_policy": "sequential",
        "volume": {
            "size": 10,
            "mount_point": "/data",
            "filesystem": "ext4"
        },
        "advanced_actions": ["change_subnet", "scale_horizontal"],
        "services": {
            "start": {
                "order": 2,
                "cmd": "USER=root /opt/hbase/bin/start.sh"

            },
            "stop": {
                "order": 2,
                "cmd": "USER=root /opt/hbase/bin/stop.sh"
            },
            "scale_out": {
                "order": 2,
                "nodes_to_execute_on": 1,
                "cmd": "USER=root /opt/hbase/bin/start-hbase.sh"
            }
        },
        "health_check": {
            "enable": true,
            "interval_sec": 60,
            "timeout_sec": 45,
            "action_timeout_sec": 45,
            "healthy_threshold": 2,
            "unhealthy_threshold": 2,
            "check_cmd": "/opt/hbase/bin/health-check.sh",
            "action_cmd": "/opt/hbase/bin/health-action.sh"
        },
        "monitor": {
            "enable": true,
            "cmd": "/usr/bin/python /opt/hbase/bin/hbase-monitor.py",
            "items": {
                "ritCount": {
                    "unit": "region",
                    "value_type": "int",
                    "statistics_type": "delta"
                }
            },
            "display": ["ritCount"]
        }
    },
    {
        "role": "hbase-hdfs-master",
        "container": {
            "type": "docker",
            "image": "img-p8hwx35j"
        },
        "count": 1,
        "cpu": 1,
        "memory": 1024,
        "vertical_scaling_policy": "sequential",
        "volume": {
            "size": 10,
            "mount_point": "/data",
            "filesystem": "ext4"
        },
        "advanced_actions": ["change_subnet"],
        "services": {
            "init": {
                "nodes_to_execute_on": 1,
                "cmd": "mkdir -p /data/hadoop;/opt/hadoop/bin/hdfs namenode -format"
            },
            "start": {
                "order": 1,
                "cmd": "USER=root /opt/hadoop/sbin/start-dfs.sh"
            },
            "stop": {
                "order": 2,
                "cmd": "USER=root /opt/hadoop/sbin/stop-dfs.sh"
            },
            "scale_in": {
                "cmd": "USER=root /opt/hadoop/sbin/exclude-node.sh",
                "timeout": 86400
            },
            "scale_out": {
                "order": 1,
                "nodes_to_execute_on": 1,
                "cmd": "USER=root /opt/hadoop/sbin/start-dfs.sh"
            }
        },
        "health_check": {
            "enable": true,
            "interval_sec": 60,
            "timeout_sec": 45,
            "action_timeout_sec": 45,
            "healthy_threshold": 2,
            "unhealthy_threshold": 2,
            "check_cmd": "sh /opt/hadoop/sbin/health-check.sh",
            "action_cmd": "sh /opt/hadoop/sbin/health-action.sh"
        },
        "monitor": {
            "enable": true,
            "cmd": "/usr/bin/python /opt/hbase/bin/hbase-monitor.py",
            "items": {
                "FilesTotal": {
                    "unit": "",
                    "value_type": "int",
                    "statistics_type": "latest"
                },
                "FilesCreated": {
                    "unit": "",
                    "value_type": "int",
                    "statistics_type": "delta"
                },
                "FilesAppended": {
                    "unit": "",
                    "value_type": "int",
                    "statistics_type": "delta"
                },
                "FilesRenamed": {
                    "unit": "",
                    "value_type": "int",
                    "statistics_type": "delta"
                },
                "FilesDeleted": {
                    "unit": "",
                    "value_type": "int",
                    "statistics_type": "delta"
                },
                "UsedGB": {
                    "unit": "",
                    "value_type": "int",
                    "statistics_type": "latest"
                },
                "RemainingGB": {
                    "unit": "",
                    "value_type": "int",
                    "statistics_type": "latest"
                },
                "PercentUsed": {
                    "unit": "",
                    "value_type": "int",
                    "statistics_type": "latest"
                },
                "PercentRemaining": {
                    "unit": "",
                    "value_type": "int",
                    "statistics_type": "latest"
                },
                "LiveNodes": {
                    "unit": "",
                    "value_type": "int",
                    "statistics_type": "latest"
                },
                "DeadNodes": {
                    "unit": "",
                    "value_type": "int",
                    "statistics_type": "latest"
                },
                "DecomLiveNodes": {
                    "unit": "",
                    "value_type": "int",
                    "statistics_type": "latest"
                },
                "DecomDeadNodes": {
                    "unit": "",
                    "value_type": "int",
                    "statistics_type": "latest"
                },
                "DecommissioningNodes": {
                    "unit": "",
                    "value_type": "int",
                    "statistics_type": "latest"
                }
            },
            "groups": {
                "DFS Files": [
                    "FilesTotal",
                    "FilesCreated",
                    "FilesAppended",
                    "FilesRenamed",
                    "FilesDeleted"
                ],
                "DFS Percentage": [
                    "PercentUsed",
                    "PercentRemaining"
                ],
                "DFS Capacity": [
                    "UsedGB",
                    "RemainingGB"
                ],
                "Data Nodes": [
                    "LiveNodes",
                    "DeadNodes",
                    "DecomLiveNodes",
                    "DecomDeadNodes",
                    "DecommissioningNodes"
                ]
            },
            "display": [
                "DFS Files",
                "DFS Percentage",
                "DFS Capacity",
                "Data Nodes"
            ]
        }
    },
    {
        "role": "hbase-slave",
        "container": {
            "type": "docker",
            "image": "img-ornvj6o7"
        },
        "count": 3,
        "cpu": 1,
        "memory": 2048,
        "vertical_scaling_policy": "sequential",
        "advanced_actions": ["change_subnet", "scale_horizontal"],
        "volume": {
            "size": 30,
            "mount_point": "/data",
            "filesystem": "ext4"
        },
        "services": {
            "stop": {
                "order": 3,
                "cmd": "USER=root /opt/hbase/bin/hbase-daemon.sh stop regionserver;USER=root /opt/hadoop/sbin/hadoop-daemon.sh stop datanode"
            },
            "start": {
                "order": 3,
                "cmd": "USER=root /opt/hadoop/sbin/start-hadoop-slave.sh;USER=root /opt/hbase/bin/start-regionserver.sh"
            }
        },
        "health_check": {
            "enable": true,
            "interval_sec": 60,
            "timeout_sec": 45,
            "action_timeout_sec": 45,
            "healthy_threshold": 2,
            "unhealthy_threshold": 2,
            "check_cmd": "sh /opt/hbase/bin/health-check.sh",
            "action_cmd": "sh /opt/hbase/bin/health-action.sh"
        },
        "monitor": {
            "enable": true,
            "cmd": "/usr/bin/python /opt/hbase/bin/hbase-monitor.py",
            "items": {
                "readRequestCount": {
                    "unit": "number per second",
                    "value_type": "int",
                    "statistics_type": "rate"
                },
                "writeRequestCount": {
                    "unit": "number per second",
                    "value_type": "int",
                    "statistics_type": "rate"
                },
                "blockCacheHitCount": {
                    "unit": "number",
                    "value_type": "int",
                    "statistics_type": "delta"
                },
                "blockCacheCountHitPercent": {
                    "unit": "%",
                    "value_type": "int",
                    "statistics_type": "latest",
                    "scale_factor_when_display": 100
                },
                "slowDeleteCount": {
                    "unit": "number",
                    "value_type": "int",
                    "statistics_type": "delta"
                },
                "slowIncrementCount": {
                    "unit": "number",
                    "value_type": "int",
                    "statistics_type": "delta"
                },
                "slowGetCount": {
                    "unit": "number",
                    "value_type": "int",
                    "statistics_type": "delta"
                },
                "slowAppendCount": {
                    "unit": "number",
                    "value_type": "int",
                    "statistics_type": "delta"
                },
                "slowPutCount": {
                    "unit": "number",
                    "value_type": "int",
                    "statistics_type": "delta"
                },
                "GcTimeMillisConcurrentMarkSweep": {
                    "unit": "ms",
                    "value_type": "int",
                    "statistics_type": "delta"
                }
            },
            "groups": {
                "RequestCount": ["readRequestCount", "writeRequestCount"],
                "SlowCount": ["slowDeleteCount", "slowIncrementCount", "slowGetCount", "slowAppendCount", "slowPutCount"]
            },
            "display": ["blockCacheHitCount", "blockCacheCountHitPercent", "RequestCount", "SlowCount", "GcTimeMillisConcurrentMarkSweep"]
        }
    }],
    "env": {
        "fs.trash.interval": 1440,
        "dfs.replication": 2,
        "dfs.namenode.handler.count": 10,
        "dfs.datanode.handler.count": 10,
        "hbase.regionserver.handler.count": 30,
        "hbase.master.handler.count": 25,
        "zookeeper.session.timeout": 60000,
        "hbase.hregion.majorcompaction": 0,
        "hbase.hstore.blockingStoreFiles": 1000000,
        "hbase.regionserver.optionalcacheflushinterval": 0,
        "hfile.block.cache.size": 0.4,
        "hfile.index.block.max.size": 131072,
        "hbase.hregion.max.filesize": 10737418240,
        "hbase.master.logcleaner.ttl": 600000,
        "hbase.ipc.server.callqueue.handler.factor": 0.1,
        "hbase.ipc.server.callqueue.read.ratio": 0,
        "hbase.ipc.server.callqueue.scan.ratio": 0,
        "hbase.regionserver.msginterval": 3000,
        "hbase.regionserver.logroll.period": 3600000,
        "hbase.regionserver.regionSplitLimit": 1000,
        "hbase.balancer.period": 300000,
        "hbase.regions.slop": 0.001,
        "io.storefile.bloom.block.size": 131072,
        "hbase.rpc.timeout": 60000,
        "hbase.regionserver.global.memstore.size": 0.4,
        "hbase.column.max.version": 1,
        "hbase.security.authorization": "true",
        "qingcloud.hbase.major.compact.hour": 3,
        "qingcloud.phoenix.on.hbase.enable": "false",
        "phoenix.functions.allowUserDefinedFunctions": "false",
        "phoenix.transactions.enabled": "false"
    },
    "endpoints": {
        "rest_port": {
            "port": 8000
        },
        "thrift_port": {
            "port": 9090
        }
    }
}
`

var hbaseCluster = models.Cluster{
	ClusterId:          "",
	Name:               "MyHBase",
	Description:        "my hbase OpApp",
	AppId:              "app-12345678",
	VersionId:          "appv-12345678",
	SubnetId:           "vxnet-t8szyjn",
	FrontgateId:        "",
	ClusterType:        0,
	Endpoints:          "{\"rest_port\":{\"port\":8000,\"protocol\":\"\"},\"thrift_port\":{\"port\":9090,\"protocol\":\"\"}}",
	Status:             "pending",
	TransitionStatus:   "",
	MetadataRootAccess: false,
	Owner:              "",
	GlobalUuid:         "",
	UpgradeStatus:      "",
	RuntimeId:          "",
}

var hbaseClusterCommons = map[string]models.ClusterCommon{
	"hbase-master": {
		ClusterId:                  "",
		Role:                       "hbase-master",
		ServerIdUpperBound:         0,
		AdvancedActions:            "change_subnet,scale_horizontal",
		InitService:                "",
		StartService:               "{\"cmd\":\"USER=root /opt/hbase/bin/start.sh\",\"order\":2}",
		StopService:                "{\"cmd\":\"USER=root /opt/hbase/bin/stop.sh\",\"order\":2}",
		ScaleOutService:            "{\"cmd\":\"USER=root /opt/hbase/bin/start-hbase.sh\",\"nodes_to_execute_on\":1,\"order\":2}",
		ScaleInService:             "",
		RestartService:             "",
		DestroyService:             "",
		UpgradeService:             "",
		CustomService:              "",
		BackupService:              "",
		RestoreService:             "",
		DeleteSnapshotService:      "",
		HealthCheck:                "{\"enable\":true,\"interval_sec\":60,\"timeout_sec\":45,\"action_timeout_sec\":45,\"healthy_threshold\":2,\"unhealthy_threshold\":2,\"check_cmd\":\"/opt/hbase/bin/health-check.sh\",\"action_cmd\":\"/opt/hbase/bin/health-action.sh\"}",
		Monitor:                    "{\"enable\":true,\"cmd\":\"/usr/bin/python /opt/hbase/bin/hbase-monitor.py\",\"items\":{\"ritCount\":{\"unit\":\"region\",\"value_type\":\"int\",\"statistics_type\":\"delta\",\"scale_factor_when_display\":0,\"enums\":null}},\"groups\":null,\"display\":[\"ritCount\"],\"alarm\":null}",
		Passphraseless:             "",
		VerticalScalingPolicy:      "sequential",
		AgentInstalled:             true,
		CustomMetadataScript:       "",
		ImageId:                    "img-srf8tv31",
		BackupPolicy:               "",
		IncrementalBackupSupported: false,
		Hypervisor:                 "docker",
	},
	"hbase-hdfs-master": {
		ClusterId:                  "",
		Role:                       "hbase-hdfs-master",
		ServerIdUpperBound:         0,
		AdvancedActions:            "change_subnet",
		InitService:                "{\"cmd\":\"mkdir -p /data/hadoop;/opt/hadoop/bin/hdfs namenode -format\",\"nodes_to_execute_on\":1,\"order\":0}",
		StartService:               "{\"cmd\":\"USER=root /opt/hadoop/sbin/start-dfs.sh\",\"order\":1}",
		StopService:                "{\"cmd\":\"USER=root /opt/hadoop/sbin/stop-dfs.sh\",\"order\":2}",
		ScaleOutService:            "{\"cmd\":\"USER=root /opt/hadoop/sbin/start-dfs.sh\",\"nodes_to_execute_on\":1,\"order\":1}",
		ScaleInService:             "{\"cmd\":\"USER=root /opt/hadoop/sbin/exclude-node.sh\",\"order\":0,\"timeout\":86400}",
		RestartService:             "",
		DestroyService:             "",
		UpgradeService:             "",
		CustomService:              "",
		BackupService:              "",
		RestoreService:             "",
		DeleteSnapshotService:      "",
		HealthCheck:                "{\"enable\":true,\"interval_sec\":60,\"timeout_sec\":45,\"action_timeout_sec\":45,\"healthy_threshold\":2,\"unhealthy_threshold\":2,\"check_cmd\":\"sh /opt/hadoop/sbin/health-check.sh\",\"action_cmd\":\"sh /opt/hadoop/sbin/health-action.sh\"}",
		Monitor:                    "{\"enable\":true,\"cmd\":\"/usr/bin/python /opt/hbase/bin/hbase-monitor.py\",\"items\":{\"DeadNodes\":{\"unit\":\"\",\"value_type\":\"int\",\"statistics_type\":\"latest\",\"scale_factor_when_display\":0,\"enums\":null},\"DecomDeadNodes\":{\"unit\":\"\",\"value_type\":\"int\",\"statistics_type\":\"latest\",\"scale_factor_when_display\":0,\"enums\":null},\"DecomLiveNodes\":{\"unit\":\"\",\"value_type\":\"int\",\"statistics_type\":\"latest\",\"scale_factor_when_display\":0,\"enums\":null},\"DecommissioningNodes\":{\"unit\":\"\",\"value_type\":\"int\",\"statistics_type\":\"latest\",\"scale_factor_when_display\":0,\"enums\":null},\"FilesAppended\":{\"unit\":\"\",\"value_type\":\"int\",\"statistics_type\":\"delta\",\"scale_factor_when_display\":0,\"enums\":null},\"FilesCreated\":{\"unit\":\"\",\"value_type\":\"int\",\"statistics_type\":\"delta\",\"scale_factor_when_display\":0,\"enums\":null},\"FilesDeleted\":{\"unit\":\"\",\"value_type\":\"int\",\"statistics_type\":\"delta\",\"scale_factor_when_display\":0,\"enums\":null},\"FilesRenamed\":{\"unit\":\"\",\"value_type\":\"int\",\"statistics_type\":\"delta\",\"scale_factor_when_display\":0,\"enums\":null},\"FilesTotal\":{\"unit\":\"\",\"value_type\":\"int\",\"statistics_type\":\"latest\",\"scale_factor_when_display\":0,\"enums\":null},\"LiveNodes\":{\"unit\":\"\",\"value_type\":\"int\",\"statistics_type\":\"latest\",\"scale_factor_when_display\":0,\"enums\":null},\"PercentRemaining\":{\"unit\":\"\",\"value_type\":\"int\",\"statistics_type\":\"latest\",\"scale_factor_when_display\":0,\"enums\":null},\"PercentUsed\":{\"unit\":\"\",\"value_type\":\"int\",\"statistics_type\":\"latest\",\"scale_factor_when_display\":0,\"enums\":null},\"RemainingGB\":{\"unit\":\"\",\"value_type\":\"int\",\"statistics_type\":\"latest\",\"scale_factor_when_display\":0,\"enums\":null},\"UsedGB\":{\"unit\":\"\",\"value_type\":\"int\",\"statistics_type\":\"latest\",\"scale_factor_when_display\":0,\"enums\":null}},\"groups\":{\"DFS Capacity\":[\"UsedGB\",\"RemainingGB\"],\"DFS Files\":[\"FilesTotal\",\"FilesCreated\",\"FilesAppended\",\"FilesRenamed\",\"FilesDeleted\"],\"DFS Percentage\":[\"PercentUsed\",\"PercentRemaining\"],\"Data Nodes\":[\"LiveNodes\",\"DeadNodes\",\"DecomLiveNodes\",\"DecomDeadNodes\",\"DecommissioningNodes\"]},\"display\":[\"DFS Files\",\"DFS Percentage\",\"DFS Capacity\",\"Data Nodes\"],\"alarm\":null}",
		Passphraseless:             "",
		VerticalScalingPolicy:      "sequential",
		AgentInstalled:             true,
		CustomMetadataScript:       "",
		ImageId:                    "img-p8hwx35j",
		BackupPolicy:               "",
		IncrementalBackupSupported: false,
		Hypervisor:                 "docker",
	},
	"hbase-slave": {
		ClusterId:                  "",
		Role:                       "hbase-slave",
		ServerIdUpperBound:         0,
		AdvancedActions:            "change_subnet,scale_horizontal",
		InitService:                "",
		StartService:               "{\"cmd\":\"USER=root /opt/hadoop/sbin/start-hadoop-slave.sh;USER=root /opt/hbase/bin/start-regionserver.sh\",\"order\":3}",
		StopService:                "{\"cmd\":\"USER=root /opt/hbase/bin/hbase-daemon.sh stop regionserver;USER=root /opt/hadoop/sbin/hadoop-daemon.sh stop datanode\",\"order\":3}",
		ScaleOutService:            "",
		ScaleInService:             "",
		RestartService:             "",
		DestroyService:             "",
		UpgradeService:             "",
		CustomService:              "",
		BackupService:              "",
		RestoreService:             "",
		DeleteSnapshotService:      "",
		HealthCheck:                "{\"enable\":true,\"interval_sec\":60,\"timeout_sec\":45,\"action_timeout_sec\":45,\"healthy_threshold\":2,\"unhealthy_threshold\":2,\"check_cmd\":\"sh /opt/hbase/bin/health-check.sh\",\"action_cmd\":\"sh /opt/hbase/bin/health-action.sh\"}",
		Monitor:                    "{\"enable\":true,\"cmd\":\"/usr/bin/python /opt/hbase/bin/hbase-monitor.py\",\"items\":{\"GcTimeMillisConcurrentMarkSweep\":{\"unit\":\"ms\",\"value_type\":\"int\",\"statistics_type\":\"delta\",\"scale_factor_when_display\":0,\"enums\":null},\"blockCacheCountHitPercent\":{\"unit\":\"%\",\"value_type\":\"int\",\"statistics_type\":\"latest\",\"scale_factor_when_display\":100,\"enums\":null},\"blockCacheHitCount\":{\"unit\":\"number\",\"value_type\":\"int\",\"statistics_type\":\"delta\",\"scale_factor_when_display\":0,\"enums\":null},\"readRequestCount\":{\"unit\":\"number per second\",\"value_type\":\"int\",\"statistics_type\":\"rate\",\"scale_factor_when_display\":0,\"enums\":null},\"slowAppendCount\":{\"unit\":\"number\",\"value_type\":\"int\",\"statistics_type\":\"delta\",\"scale_factor_when_display\":0,\"enums\":null},\"slowDeleteCount\":{\"unit\":\"number\",\"value_type\":\"int\",\"statistics_type\":\"delta\",\"scale_factor_when_display\":0,\"enums\":null},\"slowGetCount\":{\"unit\":\"number\",\"value_type\":\"int\",\"statistics_type\":\"delta\",\"scale_factor_when_display\":0,\"enums\":null},\"slowIncrementCount\":{\"unit\":\"number\",\"value_type\":\"int\",\"statistics_type\":\"delta\",\"scale_factor_when_display\":0,\"enums\":null},\"slowPutCount\":{\"unit\":\"number\",\"value_type\":\"int\",\"statistics_type\":\"delta\",\"scale_factor_when_display\":0,\"enums\":null},\"writeRequestCount\":{\"unit\":\"number per second\",\"value_type\":\"int\",\"statistics_type\":\"rate\",\"scale_factor_when_display\":0,\"enums\":null}},\"groups\":{\"RequestCount\":[\"readRequestCount\",\"writeRequestCount\"],\"SlowCount\":[\"slowDeleteCount\",\"slowIncrementCount\",\"slowGetCount\",\"slowAppendCount\",\"slowPutCount\"]},\"display\":[\"blockCacheHitCount\",\"blockCacheCountHitPercent\",\"RequestCount\",\"SlowCount\",\"GcTimeMillisConcurrentMarkSweep\"],\"alarm\":null}",
		Passphraseless:             "",
		VerticalScalingPolicy:      "sequential",
		AgentInstalled:             true,
		CustomMetadataScript:       "",
		ImageId:                    "img-ornvj6o7",
		BackupPolicy:               "",
		IncrementalBackupSupported: false,
		Hypervisor:                 "docker",
	},
}

var hbaseClusterRoles = map[string]models.ClusterRole{
	"hbase-master": {
		ClusterId:    "",
		Role:         "hbase-master",
		Cpu:          1,
		Gpu:          0,
		Memory:       1024,
		InstanceSize: 20,
		StorageSize:  10,
		MountPoint:   "/data",
		MountOptions: "defaults,noatime",
		FileSystem:   "ext4",
		Env:          "{\"dfs.datanode.handler.count\":10,\"dfs.namenode.handler.count\":10,\"dfs.replication\":2,\"fs.trash.interval\":1440,\"hbase.balancer.period\":300000,\"hbase.column.max.version\":1,\"hbase.hregion.majorcompaction\":0,\"hbase.hregion.max.filesize\":10737418240,\"hbase.hstore.blockingStoreFiles\":1000000,\"hbase.ipc.server.callqueue.handler.factor\":0.1,\"hbase.ipc.server.callqueue.read.ratio\":0,\"hbase.ipc.server.callqueue.scan.ratio\":0,\"hbase.master.handler.count\":25,\"hbase.master.logcleaner.ttl\":600000,\"hbase.regions.slop\":0.001,\"hbase.regionserver.global.memstore.size\":0.4,\"hbase.regionserver.handler.count\":30,\"hbase.regionserver.logroll.period\":3600000,\"hbase.regionserver.msginterval\":3000,\"hbase.regionserver.optionalcacheflushinterval\":0,\"hbase.regionserver.regionSplitLimit\":1000,\"hbase.rpc.timeout\":60000,\"hbase.security.authorization\":\"true\",\"hfile.block.cache.size\":0.4,\"hfile.index.block.max.size\":131072,\"io.storefile.bloom.block.size\":131072,\"phoenix.functions.allowUserDefinedFunctions\":\"false\",\"phoenix.transactions.enabled\":\"false\",\"qingcloud.hbase.major.compact.hour\":3,\"qingcloud.phoenix.on.hbase.enable\":\"false\",\"zookeeper.session.timeout\":60000}",
	},
	"hbase-hdfs-master": {
		ClusterId:    "",
		Role:         "hbase-hdfs-master",
		Cpu:          1,
		Gpu:          0,
		Memory:       1024,
		InstanceSize: 20,
		StorageSize:  10,
		MountPoint:   "/data",
		MountOptions: "defaults,noatime",
		FileSystem:   "ext4",
		Env:          "{\"dfs.datanode.handler.count\":10,\"dfs.namenode.handler.count\":10,\"dfs.replication\":2,\"fs.trash.interval\":1440,\"hbase.balancer.period\":300000,\"hbase.column.max.version\":1,\"hbase.hregion.majorcompaction\":0,\"hbase.hregion.max.filesize\":10737418240,\"hbase.hstore.blockingStoreFiles\":1000000,\"hbase.ipc.server.callqueue.handler.factor\":0.1,\"hbase.ipc.server.callqueue.read.ratio\":0,\"hbase.ipc.server.callqueue.scan.ratio\":0,\"hbase.master.handler.count\":25,\"hbase.master.logcleaner.ttl\":600000,\"hbase.regions.slop\":0.001,\"hbase.regionserver.global.memstore.size\":0.4,\"hbase.regionserver.handler.count\":30,\"hbase.regionserver.logroll.period\":3600000,\"hbase.regionserver.msginterval\":3000,\"hbase.regionserver.optionalcacheflushinterval\":0,\"hbase.regionserver.regionSplitLimit\":1000,\"hbase.rpc.timeout\":60000,\"hbase.security.authorization\":\"true\",\"hfile.block.cache.size\":0.4,\"hfile.index.block.max.size\":131072,\"io.storefile.bloom.block.size\":131072,\"phoenix.functions.allowUserDefinedFunctions\":\"false\",\"phoenix.transactions.enabled\":\"false\",\"qingcloud.hbase.major.compact.hour\":3,\"qingcloud.phoenix.on.hbase.enable\":\"false\",\"zookeeper.session.timeout\":60000}",
	},
	"hbase-slave": {
		ClusterId:    "",
		Role:         "hbase-slave",
		Cpu:          1,
		Gpu:          0,
		Memory:       2048,
		InstanceSize: 20,
		StorageSize:  30,
		MountPoint:   "/data",
		MountOptions: "defaults,noatime",
		FileSystem:   "ext4",
		Env:          "{\"dfs.datanode.handler.count\":10,\"dfs.namenode.handler.count\":10,\"dfs.replication\":2,\"fs.trash.interval\":1440,\"hbase.balancer.period\":300000,\"hbase.column.max.version\":1,\"hbase.hregion.majorcompaction\":0,\"hbase.hregion.max.filesize\":10737418240,\"hbase.hstore.blockingStoreFiles\":1000000,\"hbase.ipc.server.callqueue.handler.factor\":0.1,\"hbase.ipc.server.callqueue.read.ratio\":0,\"hbase.ipc.server.callqueue.scan.ratio\":0,\"hbase.master.handler.count\":25,\"hbase.master.logcleaner.ttl\":600000,\"hbase.regions.slop\":0.001,\"hbase.regionserver.global.memstore.size\":0.4,\"hbase.regionserver.handler.count\":30,\"hbase.regionserver.logroll.period\":3600000,\"hbase.regionserver.msginterval\":3000,\"hbase.regionserver.optionalcacheflushinterval\":0,\"hbase.regionserver.regionSplitLimit\":1000,\"hbase.rpc.timeout\":60000,\"hbase.security.authorization\":\"true\",\"hfile.block.cache.size\":0.4,\"hfile.index.block.max.size\":131072,\"io.storefile.bloom.block.size\":131072,\"phoenix.functions.allowUserDefinedFunctions\":\"false\",\"phoenix.transactions.enabled\":\"false\",\"qingcloud.hbase.major.compact.hour\":3,\"qingcloud.phoenix.on.hbase.enable\":\"false\",\"zookeeper.session.timeout\":60000}",
	},
}

var hbaseClusterNodes = map[string]models.ClusterNode{
	"hbase-master1": {
		NodeId:           "hbase-master1",
		ClusterId:        "",
		Name:             "",
		InstanceId:       "",
		VolumeId:         "",
		SubnetId:         "vxnet-t8szyjn",
		PrivateIp:        "",
		ServerId:         1,
		Role:             "hbase-master",
		Status:           "pending",
		TransitionStatus: "",
		GroupId:          1,
		Owner:            "",
		GlobalServerId:   0,
		CustomMetadata:   "",
		PubKey:           "",
		HealthStatus:     "",
		IsBackup:         false,
		AutoBackup:       false,
	},
	"hbase-hdfs-master1": {
		NodeId:           "hbase-hdfs-master1",
		ClusterId:        "",
		Name:             "",
		InstanceId:       "",
		VolumeId:         "",
		SubnetId:         "vxnet-t8szyjn",
		PrivateIp:        "",
		ServerId:         1,
		Role:             "hbase-hdfs-master",
		Status:           "pending",
		TransitionStatus: "",
		GroupId:          1,
		Owner:            "",
		GlobalServerId:   0,
		CustomMetadata:   "",
		PubKey:           "",
		HealthStatus:     "",
		IsBackup:         false,
		AutoBackup:       false,
	},
	"hbase-slave1": {
		NodeId:           "hbase-slave1",
		ClusterId:        "",
		Name:             "",
		InstanceId:       "",
		VolumeId:         "",
		SubnetId:         "vxnet-t8szyjn",
		PrivateIp:        "",
		ServerId:         1,
		Role:             "hbase-slave",
		Status:           "pending",
		TransitionStatus: "",
		GroupId:          1,
		Owner:            "",
		GlobalServerId:   0,
		CustomMetadata:   "",
		PubKey:           "",
		HealthStatus:     "",
		IsBackup:         false,
		AutoBackup:       false,
	},
	"hbase-slave2": {
		NodeId:           "hbase-slave2",
		ClusterId:        "",
		Name:             "",
		InstanceId:       "",
		VolumeId:         "",
		SubnetId:         "vxnet-t8szyjn",
		PrivateIp:        "",
		ServerId:         2,
		Role:             "hbase-slave",
		Status:           "pending",
		TransitionStatus: "",
		GroupId:          2,
		Owner:            "",
		GlobalServerId:   0,
		CustomMetadata:   "",
		PubKey:           "",
		HealthStatus:     "",
		IsBackup:         false,
		AutoBackup:       false,
	},
	"hbase-slave3": {
		NodeId:           "hbase-slave3",
		ClusterId:        "",
		Name:             "",
		InstanceId:       "",
		VolumeId:         "",
		SubnetId:         "vxnet-t8szyjn",
		PrivateIp:        "",
		ServerId:         3,
		Role:             "hbase-slave",
		Status:           "pending",
		TransitionStatus: "",
		GroupId:          3,
		Owner:            "",
		GlobalServerId:   0,
		CustomMetadata:   "",
		PubKey:           "",
		HealthStatus:     "",
		IsBackup:         false,
		AutoBackup:       false,
	},
}

var hbaseClusterLinks = map[string]models.ClusterLink{
	"zk_service": {
		ClusterId:         "",
		Name:              "zk_service",
		ExternalClusterId: "cl-w8qh5hf6",
		Owner:             "",
	},
}

func getTestClusterWrapper(t *testing.T) *models.ClusterWrapper {
	cluster := opapp.ClusterConf{}
	err := jsonutil.Decode([]byte(hbaseMustache), &cluster)
	if err != nil {
		t.Fatalf("Parse mustache failed: %+v", err)
	}
	// TODO: add validate to test
	//cluster.RenderJson = hbaseMustache
	//err = cluster.Validate()
	//if err != nil {
	//	t.Fatalf("Validate cluster failed: %+v", err)
	//}

	parser := Parser{}
	clusterWrapper := new(models.ClusterWrapper)
	err = parser.Parse(cluster, clusterWrapper, "")
	if err != nil {
		t.Fatalf("Parse mustache failed: %+v", err)
	}
	return clusterWrapper
}

func TestParse(t *testing.T) {
	clusterWrapper := getTestClusterWrapper(t)

	// check cluster
	if hbaseCluster != *clusterWrapper.Cluster {
		t.Errorf("ClusterConf not equal")
		t.Logf("ori: %+v", hbaseCluster)
		t.Logf("dst: %+v", *clusterWrapper.Cluster)
	}

	// check cluster common
	if len(hbaseClusterCommons) != len(clusterWrapper.ClusterCommons) {
		t.Errorf("ClusterConf common length not equal, ori: %d, dst: %d",
			len(hbaseClusterCommons), len(clusterWrapper.ClusterCommons))
	}
	for index := range clusterWrapper.ClusterCommons {
		if hbaseClusterCommons[index] != *clusterWrapper.ClusterCommons[index] {
			t.Errorf("ClusterConf common [%s] not equal.", index)
			t.Logf("ori: %+v", hbaseClusterCommons[index])
			t.Logf("dst: %+v", *clusterWrapper.ClusterCommons[index])
		}
	}

	// check cluser role
	if len(hbaseClusterRoles) != len(clusterWrapper.ClusterRoles) {
		t.Errorf("ClusterConf role length not equal, ori: %d, dst: %d",
			len(hbaseClusterRoles), len(clusterWrapper.ClusterRoles))
	}
	for index := range clusterWrapper.ClusterRoles {
		if hbaseClusterRoles[index] != *clusterWrapper.ClusterRoles[index] {
			t.Errorf("ClusterConf role [%s] not equal.", index)
			t.Logf("ori: %+v", hbaseClusterRoles[index])
			t.Logf("dst: %+v", *clusterWrapper.ClusterRoles[index])
		}
	}

	// check cluser node
	if len(hbaseClusterNodes) != len(clusterWrapper.ClusterNodesWithKeyPairs) {
		t.Errorf("ClusterConf node length not equal, ori: %d, dst: %d",
			len(hbaseClusterNodes), len(clusterWrapper.ClusterNodesWithKeyPairs))
	}
	for index := range clusterWrapper.ClusterNodesWithKeyPairs {
		if hbaseClusterNodes[index] != *clusterWrapper.ClusterNodesWithKeyPairs[index].ClusterNode {
			t.Errorf("ClusterConf node [%s] not equal.", index)
			t.Logf("ori: %+v", hbaseClusterNodes[index])
			t.Logf("dst: %+v", clusterWrapper.ClusterNodesWithKeyPairs[index].ClusterNode)
		}
	}

	// check cluser link
	if len(hbaseClusterLinks) != len(clusterWrapper.ClusterLinks) {
		t.Errorf("ClusterConf link length not equal, ori: %d, dst: %d",
			len(hbaseClusterLinks), len(clusterWrapper.ClusterLinks))
	}
	for index := range clusterWrapper.ClusterLinks {
		if hbaseClusterLinks[index] != *clusterWrapper.ClusterLinks[index] {
			t.Errorf("ClusterConf link [%s] not equal", index)
			t.Logf("ori: %+v", hbaseClusterLinks[index])
			t.Logf("dst: %+v", *clusterWrapper.ClusterLinks[index])
		}
	}

	// check cluser loadbalancer
	if 0 != len(clusterWrapper.ClusterLoadbalancers) {
		t.Errorf("ClusterConf loadbalancer length not equal, ori: %d, dst: %d",
			0, len(clusterWrapper.ClusterLoadbalancers))
	}
}
