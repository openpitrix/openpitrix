// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"time"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/util/idutil"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
)

const AppVersionTableName = "app_version"

func NewAppVersionId() string {
	return idutil.GetUuid("appv-")
}

type AppVersion struct {
	VersionId   string
	AppId       string
	Owner       string
	Name        string
	Description string
	PackageName string
	Status      string
	Sequence    uint32
	CreateTime  time.Time
	StatusTime  time.Time
	UpdateTime  *time.Time
}

var AppVersionColumns = GetColumnsFromStruct(&AppVersion{})

func NewAppVersion(appId, name, description, owner, packageName string) *AppVersion {
	return &AppVersion{
		VersionId:   NewAppVersionId(),
		AppId:       appId,
		Name:        name,
		Owner:       owner,
		PackageName: packageName,
		Description: description,
		Status:      constants.StatusActive,
		CreateTime:  time.Now(),
		StatusTime:  time.Now(),
	}
}

func AppVersionToPb(appVersion *AppVersion) *pb.AppVersion {
	if appVersion == nil {
		return nil
	}
	pbAppVersion := pb.AppVersion{}
	pbAppVersion.VersionId = pbutil.ToProtoString(appVersion.VersionId)
	pbAppVersion.AppId = pbutil.ToProtoString(appVersion.AppId)
	pbAppVersion.Name = pbutil.ToProtoString(appVersion.Name)
	pbAppVersion.Description = pbutil.ToProtoString(appVersion.Description)
	pbAppVersion.Status = pbutil.ToProtoString(appVersion.Status)
	pbAppVersion.PackageName = pbutil.ToProtoString(appVersion.PackageName)
	pbAppVersion.Owner = pbutil.ToProtoString(appVersion.Owner)
	pbAppVersion.CreateTime = pbutil.ToProtoTimestamp(appVersion.CreateTime)
	pbAppVersion.StatusTime = pbutil.ToProtoTimestamp(appVersion.StatusTime)
	pbAppVersion.Sequence = uint32(appVersion.Sequence)
	if appVersion.UpdateTime != nil {
		pbAppVersion.UpdateTime = pbutil.ToProtoTimestamp(*appVersion.UpdateTime)
	}
	return &pbAppVersion
}

func AppVersionsToPbs(appVersions []*AppVersion) (pbAppVersions []*pb.AppVersion) {
	for _, appVersion := range appVersions {
		pbAppVersions = append(pbAppVersions, AppVersionToPb(appVersion))
	}
	return
}
