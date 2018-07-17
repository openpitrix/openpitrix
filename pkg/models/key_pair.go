// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"time"

	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/util/idutil"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
)

const KeyPairTableName = "key_pair"

func NewKeyPairId() string {
	return idutil.GetUuid("kp-")
}

type KeyPair struct {
	KeyPairId   string
	Name        string
	Description string
	PubKey      string
	Owner       string
	CreateTime  time.Time
	StatusTime  time.Time
}

var KeyPairColumns = GetColumnsFromStruct(&KeyPair{})

func KeyPairToPb(keyPair *KeyPair) *pb.KeyPair {
	pbKeyPair := pb.KeyPair{}
	pbKeyPair.KeyPairId = pbutil.ToProtoString(keyPair.KeyPairId)
	pbKeyPair.Name = pbutil.ToProtoString(keyPair.Name)
	pbKeyPair.Description = pbutil.ToProtoString(keyPair.Description)
	pbKeyPair.PubKey = pbutil.ToProtoString(keyPair.PubKey)
	pbKeyPair.Owner = pbutil.ToProtoString(keyPair.Owner)
	pbKeyPair.CreateTime = pbutil.ToProtoTimestamp(keyPair.CreateTime)
	pbKeyPair.StatusTime = pbutil.ToProtoTimestamp(keyPair.StatusTime)
	return &pbKeyPair
}

func KeyPairsToPbs(keyPairs []*KeyPair) (pbKeyPairs []*pb.KeyPair) {
	for _, keyPair := range keyPairs {
		pbKeyPairs = append(pbKeyPairs, KeyPairToPb(keyPair))
	}
	return
}
