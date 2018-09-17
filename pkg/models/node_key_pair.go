// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/util/jsonutil"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
)

type NodeKeyPairDetails []NodeKeyPairDetail

func NewNodeKeyPairDetails(data string) (NodeKeyPairDetails, error) {
	var nodeKeyPairDetails NodeKeyPairDetails
	err := jsonutil.Decode([]byte(data), &nodeKeyPairDetails)
	if err != nil {
		logger.Error(nil, "Decode [%s] into node key pair details failed: %+v", data, err)
	}
	return nodeKeyPairDetails, err
}

type NodeKeyPair struct {
	KeyPairId string
	NodeId    string
}

func NodeKeyPairToPb(nodeKeyPair *NodeKeyPair) *pb.NodeKeyPair {
	n := &pb.NodeKeyPair{
		NodeId:    pbutil.ToProtoString(nodeKeyPair.NodeId),
		KeyPairId: pbutil.ToProtoString(nodeKeyPair.KeyPairId),
	}
	return n
}

type NodeKeyPairDetail struct {
	NodeKeyPair *NodeKeyPair
	ClusterNode *ClusterNode
	KeyPair     *KeyPair
}

var NodeKeyPairColumns = db.GetColumnsFromStruct(&NodeKeyPair{})
