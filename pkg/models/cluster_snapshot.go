// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

type ClusterSnapshot struct {
	SnapshotId       string
	Role             string
	ServerIds        string
	count            int32
	AppId            string
	VersionId        string
	ChildSnapshotIds string
	Size             int32
}

var ClusterSnapshotColumns = GetColumnsFromStruct(&ClusterSnapshot{})
