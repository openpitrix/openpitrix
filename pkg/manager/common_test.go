// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package manager

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
)

func TestBuildUpdateAttributes(t *testing.T) {
	clusterId := "cl-xxxx"
	req := &pb.Cluster{
		ClusterId:          pbutil.ToProtoString(clusterId),
		TransitionStatus:   pbutil.ToProtoString("creating"),
		MetadataRootAccess: pbutil.ToProtoBool(false),
		CreateTime:         pbutil.ToProtoTimestamp(time.Now()),
		StatusTime:         pbutil.ToProtoTimestamp(time.Now()),
	}
	attributes := BuildUpdateAttributes(req, models.ClusterColumns...)
	t.Log(attributes)
	assert.Equal(t, attributes["cluster_id"], clusterId)
}
