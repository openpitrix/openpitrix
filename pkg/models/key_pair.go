// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"time"

	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/sender"
	"openpitrix.io/openpitrix/pkg/util/idutil"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
)

func NewKeyPairId() string {
	return idutil.GetUuid("kp-")
}

type KeyPair struct {
	KeyPairId   string
	Name        string
	Description string
	PubKey      string
	Owner       string
	OwnerPath   sender.OwnerPath
	CreateTime  time.Time
	StatusTime  time.Time
}

type KeyPairWithNodes struct {
	*KeyPair
	NodeId []string
}

var KeyPairColumns = db.GetColumnsFromStruct(&KeyPair{})

func KeyPairNodesToPb(keyPairNodes *KeyPairWithNodes) *pb.KeyPair {
	pbKeyPair := pb.KeyPair{}
	pbKeyPair.KeyPairId = pbutil.ToProtoString(keyPairNodes.KeyPairId)
	pbKeyPair.Name = pbutil.ToProtoString(keyPairNodes.Name)
	pbKeyPair.Description = pbutil.ToProtoString(keyPairNodes.Description)
	pbKeyPair.PubKey = pbutil.ToProtoString(keyPairNodes.PubKey)
	pbKeyPair.OwnerPath = keyPairNodes.OwnerPath.ToProtoString()
	pbKeyPair.CreateTime = pbutil.ToProtoTimestamp(keyPairNodes.CreateTime)
	pbKeyPair.StatusTime = pbutil.ToProtoTimestamp(keyPairNodes.StatusTime)
	pbKeyPair.NodeId = keyPairNodes.NodeId
	return &pbKeyPair
}

func KeyPairNodesToPbs(keyPairNodes []*KeyPairWithNodes) (pbKeyPairs []*pb.KeyPair) {
	for _, keyPairNodesItem := range keyPairNodes {
		pbKeyPairs = append(pbKeyPairs, KeyPairNodesToPb(keyPairNodesItem))
	}
	return
}
