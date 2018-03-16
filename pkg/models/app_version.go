// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"time"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/utils"
)

const AppVersionTableName = "app_version"

func NewAppVersionId() string {
	return utils.GetUuid("appv-")
}

type AppVersion struct {
	VersionId   string
	AppId       string
	Owner       string
	Name        string
	Description string
	PackageName string
	Status      string
	CreateTime  time.Time
	StatusTime  time.Time
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
	pbAppVersion := pb.AppVersion{}
	pbAppVersion.VersionId = utils.ToProtoString(appVersion.VersionId)
	pbAppVersion.AppId = utils.ToProtoString(appVersion.AppId)
	pbAppVersion.Name = utils.ToProtoString(appVersion.Name)
	pbAppVersion.Description = utils.ToProtoString(appVersion.Description)
	pbAppVersion.Status = utils.ToProtoString(appVersion.Status)
	pbAppVersion.PackageName = utils.ToProtoString(appVersion.PackageName)
	pbAppVersion.Owner = utils.ToProtoString(appVersion.Owner)
	pbAppVersion.CreateTime = utils.ToProtoTimestamp(appVersion.CreateTime)
	pbAppVersion.StatusTime = utils.ToProtoTimestamp(appVersion.StatusTime)
	return &pbAppVersion
}

func AppVersionsToPbs(appVersions []*AppVersion) (pbAppVersions []*pb.AppVersion) {
	for _, appVersion := range appVersions {
		pbAppVersions = append(pbAppVersions, AppVersionToPb(appVersion))
	}
	return
}
