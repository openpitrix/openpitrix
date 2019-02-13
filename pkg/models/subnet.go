// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"time"

	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
)

type Subnet struct {
	Zone        string
	SubnetId    string
	Name        string
	CreateTime  time.Time
	Description string
	InstanceIds []string
	VpcId       string
}

func SubnetToPb(subnet *Subnet) *pb.Subnet {
	pbSubnet := pb.Subnet{}
	pbSubnet.Zone = pbutil.ToProtoString(subnet.Zone)
	pbSubnet.SubnetId = pbutil.ToProtoString(subnet.SubnetId)
	pbSubnet.Name = pbutil.ToProtoString(subnet.Name)
	pbSubnet.CreateTime = pbutil.ToProtoTimestamp(subnet.CreateTime)
	pbSubnet.Description = pbutil.ToProtoString(subnet.Description)
	pbSubnet.InstanceId = subnet.InstanceIds
	pbSubnet.VpcId = pbutil.ToProtoString(subnet.VpcId)
	return &pbSubnet
}

func SubnetsToPbs(subnets []*Subnet) (pbSubnets []*pb.Subnet) {
	for _, subnet := range subnets {
		pbSubnets = append(pbSubnets, SubnetToPb(subnet))
	}
	return
}

func PbToSubnet(pbSubnet *pb.Subnet) *Subnet {
	return &Subnet{
		Zone:        pbSubnet.GetZone().GetValue(),
		SubnetId:    pbSubnet.GetSubnetId().GetValue(),
		Name:        pbSubnet.GetName().GetValue(),
		CreateTime:  pbutil.GetTime(pbSubnet.GetCreateTime()),
		Description: pbSubnet.GetDescription().GetValue(),
		InstanceIds: pbSubnet.GetInstanceId(),
		VpcId:       pbSubnet.GetVpcId().GetValue(),
	}
}
