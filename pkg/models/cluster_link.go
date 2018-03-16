// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

const ClusterLinkTableName = "cluster_link"

type ClusterLink struct {
	ClusterId         string
	Name              string
	ExternalClusterId string
	Owner             string
}

var ClusterLinkColumns = GetColumnsFromStruct(&ClusterLink{})
