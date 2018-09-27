// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package job

import (
	"context"

	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/util/senderutil"
)

func (p *Server) Builder(ctx context.Context, req interface{}) interface{} {
	sender := senderutil.GetSenderFromContext(ctx)
	switch r := req.(type) {
	case *pb.DescribeJobsRequest:
		if sender.IsGlobalAdmin() {

		} else {
			r.Owner = []string{sender.UserId}
		}
		return r
	}
	return req
}
