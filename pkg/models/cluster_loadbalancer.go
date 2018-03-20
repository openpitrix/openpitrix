// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

const ClusterLoadbalancerTableName = "cluster_loadbalancer"

type ClusterLoadbalancer struct {
	ClusterId              string
	Role                   string
	LoadbalancerListenerId string
	LoadbalancerPort       int32
	LoadbalancerPolicyId   string
}

var ClusterLoadbalancerColumns = GetColumnsFromStruct(&ClusterLoadbalancer{})
