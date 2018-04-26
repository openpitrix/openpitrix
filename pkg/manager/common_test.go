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
	"openpitrix.io/openpitrix/pkg/utils"
)

func TestBuildUpdateAttributes(t *testing.T) {
	clusterId := "cl-xxxx"
	req := &pb.Cluster{
		ClusterId:          utils.ToProtoString(clusterId),
		TransitionStatus:   utils.ToProtoString("creating"),
		MetadataRootAccess: utils.ToProtoBool(false),
		CreateTime:         utils.ToProtoTimestamp(time.Now()),
		StatusTime:         utils.ToProtoTimestamp(time.Now()),
	}
	attributes := BuildUpdateAttributes(req, models.ClusterColumns...)
	t.Log(attributes)
	assert.Equal(t, attributes["cluster_id"], clusterId)
}
