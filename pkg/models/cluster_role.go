// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

const ClusterRoleTableName = "cluster_role"

type ClusterRole struct {
	ClusterId    string
	Role         string
	Cpu          int32
	Gpu          int32
	Memory       int32
	InstanceSize int32
	StorageSize  int32
	MountPoint   string
	MountOptions string
	Env          string
}

var ClusterRoleColumns = GetColumnsFromStruct(&ClusterRole{})
