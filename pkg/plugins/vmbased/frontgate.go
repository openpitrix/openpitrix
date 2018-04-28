// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package vmbased

import "openpitrix.io/openpitrix/pkg/models"

type Frontgate struct {
	*Frame
}

func (f *Frontgate) CreateClusterLayer() *models.TaskLayer {
	var nodeIds []string
	for nodeId := range f.ClusterWrapper.ClusterNodes {
		nodeIds = append(nodeIds, nodeId)
	}
	headTaskLayer := new(models.TaskLayer)

	headTaskLayer.
		Append(f.createVolumesLayer(nodeIds)).       // create volume
		Append(f.runInstancesLayer(nodeIds)).        // run instance and attach volume to instance
		Append(f.formatAndMountVolumeLayer(nodeIds)) // format and mount volume to instance

	return headTaskLayer.Child
}
