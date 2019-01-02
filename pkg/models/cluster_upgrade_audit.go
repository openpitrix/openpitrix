// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"time"

	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/sender"
)

type ClusterUpgradeAudit struct {
	ClusterUpgradeAuditId string
	ClusterId             string
	FromVersionId         string
	ToVersionId           string
	ServiceParams         string
	CreateTime            time.Time
	StatusTime            time.Time
	Status                string
	Owner                 string
	OwnerPath             sender.OwnerPath
}

var ClusterUpgradeAuditColumns = db.GetColumnsFromStruct(&ClusterUpgradeAudit{})
