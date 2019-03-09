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

func NewMarketId() string {
	return idutil.GetUuid("mkt-")
}

type Market struct {
	MarketId    string
	Name        string
	Visibility  string
	Status      string
	Owner       string
	OwnerPath   sender.OwnerPath
	Description string
	CreateTime  time.Time
	StatusTime  time.Time
}

var MarketColumns = db.GetColumnsFromStruct(&Market{})

func NewMarket(name, visibility, status, description string, ownerPath sender.OwnerPath) *Market {
	return &Market{
		MarketId:    NewMarketId(),
		Name:        name,
		Visibility:  visibility,
		Status:      status,
		Owner:       ownerPath.Owner(),
		OwnerPath:   ownerPath,
		Description: description,
		CreateTime:  time.Now(),
		StatusTime:  time.Now(),
	}
}

func MarketToPb(market *Market) *pb.Market {
	pbMarket := pb.Market{}
	pbMarket.MarketId = pbutil.ToProtoString(market.MarketId)
	pbMarket.Name = pbutil.ToProtoString(market.Name)
	pbMarket.Visibility = pbutil.ToProtoString(market.Visibility)
	pbMarket.Status = pbutil.ToProtoString(market.Status)
	pbMarket.OwnerPath = market.OwnerPath.ToProtoString()
	pbMarket.Owner = pbutil.ToProtoString(market.Owner)
	pbMarket.Description = pbutil.ToProtoString(market.Description)
	pbMarket.CreateTime = pbutil.ToProtoTimestamp(market.CreateTime)
	pbMarket.StatusTime = pbutil.ToProtoTimestamp(market.StatusTime)

	return &pbMarket
}

func MarketToPbs(markets []*Market) (pbMarkets []*pb.Market) {
	for _, market := range markets {
		pbMarkets = append(pbMarkets, MarketToPb(market))
	}

	return
}
