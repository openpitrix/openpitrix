// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"time"

	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
)

type AppVersionAudit struct {
	VersionId  string
	AppId      string
	Status     string
	Operator   string
	Role       string
	Message    string
	StatusTime time.Time
}

var AppVersionAuditColumns = db.GetColumnsFromStruct(&AppVersionAudit{})

func NewAppVersionAudit(versionId, appId, status, operator, role string) *AppVersionAudit {
	return &AppVersionAudit{
		VersionId:  versionId,
		AppId:      appId,
		Status:     status,
		Operator:   operator,
		Role:       role,
		StatusTime: time.Now(),
	}
}

func AppVersionAuditToPb(appVersionAudit *AppVersionAudit) *pb.AppVersionAudit {
	if appVersionAudit == nil {
		return nil
	}
	pbAppVersionAudit := pb.AppVersionAudit{}
	pbAppVersionAudit.VersionId = pbutil.ToProtoString(appVersionAudit.VersionId)
	pbAppVersionAudit.AppId = pbutil.ToProtoString(appVersionAudit.AppId)
	pbAppVersionAudit.Status = pbutil.ToProtoString(appVersionAudit.Status)
	pbAppVersionAudit.Operator = pbutil.ToProtoString(appVersionAudit.Operator)
	pbAppVersionAudit.Role = pbutil.ToProtoString(appVersionAudit.Role)
	pbAppVersionAudit.Message = pbutil.ToProtoString(appVersionAudit.Message)
	pbAppVersionAudit.StatusTime = pbutil.ToProtoTimestamp(appVersionAudit.StatusTime)
	return &pbAppVersionAudit
}

func AppVersionAuditsToPbs(appVersionAudits []*AppVersionAudit) (pbAppVersionAudits []*pb.AppVersionAudit) {
	for _, appVersionAudit := range appVersionAudits {
		pbAppVersionAudits = append(pbAppVersionAudits, AppVersionAuditToPb(appVersionAudit))
	}
	return
}
