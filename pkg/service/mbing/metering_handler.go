// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package mbing

import (
	"context"

	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
)

func (s *Server) StartMetering(ctx context.Context, req *pb.MeteringRequest) (*pb.MeteringResponse, error) {
	var leasings []*models.Leasing
	//construct leasings
	for _, mSku := range req.GetSkuSet() {
		//check if sku exist and get renewalTime
		renewalTime, err := renewalTimeFromSku(ctx, mSku.GetSkuId().GetValue(), pbutil.FromProtoTimestamp(mSku.GetActionTime()))
		if err != nil {
			if err == db.ErrNotFound {
				return nil, commonInternalErr(ctx, models.Sku{}, NotExistCode)
			} else {
				return nil, commonInternalErr(ctx, models.Leasing{}, CreateFailedCode)
			}
		}
		leasings = append(leasings, models.PbToLeasing(req, mSku, GetGroupId(), renewalTime))
	}

	//insert leasings
	err := insertLeasings(ctx, leasings)
	if err != nil {
		return nil, commonInternalErr(ctx, models.Leasing{}, CreateFailedCode)
	}

	//TODO: Add leasing to REDIS if duration exist.
	//TODO: Add leasing to ETCD.

	//MeteringResponse
	var res pb.MeteringResponse
	for _, l := range leasings {
		res.LeasingIds = append(res.LeasingIds, pbutil.ToProtoString(l.LeasingId))
	}
	return &res, nil
}

func GetGroupId() string {
	return "Group_01"
}
