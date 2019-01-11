// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"time"

	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
)

type Vpc struct {
	VpcId            string
	Name             string
	CreateTime       time.Time
	Description      string
	Status           string
	TransitionStatus string
	Subnets          []string
	Eip              *Eip
}

type Eip struct {
	EipId string
	Name  string
	Addr  string
}

func EipToPb(eip *Eip) *pb.Eip {
	pbEip := pb.Eip{}
	pbEip.EipId = pbutil.ToProtoString(eip.EipId)
	pbEip.Name = pbutil.ToProtoString(eip.Name)
	pbEip.Addr = pbutil.ToProtoString(eip.Addr)
	return &pbEip
}

func PbToEip(pbEip *pb.Eip) *Eip {
	return &Eip{
		EipId: pbEip.GetEipId().GetValue(),
		Name:  pbEip.GetName().GetValue(),
		Addr:  pbEip.GetAddr().GetValue(),
	}
}

func VpcToPb(vpc *Vpc) *pb.Vpc {
	pbVpc := pb.Vpc{}
	pbVpc.VpcId = pbutil.ToProtoString(vpc.VpcId)
	pbVpc.Name = pbutil.ToProtoString(vpc.Name)
	pbVpc.CreateTime = pbutil.ToProtoTimestamp(vpc.CreateTime)
	pbVpc.Description = pbutil.ToProtoString(vpc.Description)
	pbVpc.Status = pbutil.ToProtoString(vpc.Status)
	pbVpc.TransitionStatus = pbutil.ToProtoString(vpc.TransitionStatus)
	pbVpc.Subnets = vpc.Subnets
	pbVpc.Eip = EipToPb(vpc.Eip)

	return &pbVpc
}

func PbToVpc(pbVpc *pb.Vpc) *Vpc {
	return &Vpc{
		VpcId:            pbVpc.GetVpcId().GetValue(),
		Name:             pbVpc.GetName().GetValue(),
		CreateTime:       pbutil.GetTime(pbVpc.GetCreateTime()),
		Description:      pbVpc.GetDescription().GetValue(),
		Status:           pbVpc.GetStatus().GetValue(),
		TransitionStatus: pbVpc.GetTransitionStatus().GetValue(),
		Subnets:          pbVpc.GetSubnets(),
		Eip:              PbToEip(pbVpc.GetEip()),
	}
}
