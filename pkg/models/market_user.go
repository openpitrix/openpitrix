// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"time"

	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/sender"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
)

type MarketUser struct {
	MarketId   string
	UserId     string
	Owner      string
	OwnerPath  sender.OwnerPath
	CreateTime time.Time
}

var MarketUserColumns = db.GetColumnsFromStruct(&MarketUser{})

func NewMarketUser(marketId, userId string, ownerPath sender.OwnerPath) *MarketUser {
	return &MarketUser{
		MarketId:   marketId,
		UserId:     userId,
		Owner:      ownerPath.Owner(),
		OwnerPath:  ownerPath,
		CreateTime: time.Now(),
	}
}

func MarketUserToPb(marketUser *MarketUser) *pb.MarketUser {
	pbMarketUser := pb.MarketUser{}
	pbMarketUser.MarketId = pbutil.ToProtoString(marketUser.MarketId)
	pbMarketUser.UserId = pbutil.ToProtoString(marketUser.UserId)
	pbMarketUser.OwnerPath = marketUser.OwnerPath.ToProtoString()
	pbMarketUser.Owner = pbutil.ToProtoString(marketUser.Owner)
	pbMarketUser.CreateTime = pbutil.ToProtoTimestamp(marketUser.CreateTime)
	return &pbMarketUser
}

func MarketUserToPbs(marketUsers []*MarketUser) (pbMarketUsers []*pb.MarketUser) {
	for _, marketUser := range marketUsers {
		pbMarketUsers = append(pbMarketUsers, MarketUserToPb(marketUser))
	}
	return
}
