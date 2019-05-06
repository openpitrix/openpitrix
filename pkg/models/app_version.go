// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"strings"
	"time"

	"github.com/Masterminds/semver"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/sender"
	"openpitrix.io/openpitrix/pkg/util/idutil"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
)

func NewAppVersionId() string {
	return idutil.GetUuid("appv-")
}

type AppVersion struct {
	VersionId   string
	Active      bool
	AppId       string
	Owner       string
	OwnerPath   sender.OwnerPath
	Name        string
	Description string
	PackageName string
	Home        string
	Icon        string
	Screenshots string
	Maintainers string
	Keywords    string
	Sources     string
	Readme      string
	Status      string
	ReviewId    string
	Message     string
	Type        string
	Sequence    uint32
	CreateTime  time.Time
	StatusTime  time.Time
	UpdateTime  *time.Time
}

var AppVersionColumns = db.GetColumnsFromStruct(&AppVersion{})

func (v AppVersion) GetSemver() string {
	return strings.Split(v.Name, " ")[0]
}

type AppVersions []*AppVersion

// Len returns the length.
func (c AppVersions) Len() int { return len(c) }

// Swap swaps the position of two items in the versions slice.
func (c AppVersions) Swap(i, j int) { c[i], c[j] = c[j], c[i] }

// Less returns true if the version of entry a is less than the version of entry b.
func (c AppVersions) Less(a, b int) bool {
	// Failed parse pushes to the back.
	aVersion := c[a]
	bVersion := c[b]
	i, err := semver.NewVersion(aVersion.GetSemver())
	if err != nil {
		return true
	}
	j, err := semver.NewVersion(bVersion.GetSemver())
	if err != nil {
		return false
	}
	if i.Equal(j) {
		return aVersion.CreateTime.Before(bVersion.CreateTime)
	}
	return i.LessThan(j)
}

func NewAppVersion(appId, name, description string, ownerPath sender.OwnerPath) *AppVersion {
	return &AppVersion{
		VersionId:   NewAppVersionId(),
		Active:      false,
		AppId:       appId,
		Name:        name,
		Owner:       ownerPath.Owner(),
		OwnerPath:   ownerPath,
		Description: description,
		Status:      constants.StatusDraft,
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
	pbAppVersion.Active = pbutil.ToProtoBool(appVersion.Active)
	pbAppVersion.AppId = pbutil.ToProtoString(appVersion.AppId)
	pbAppVersion.Name = pbutil.ToProtoString(appVersion.Name)
	pbAppVersion.Description = pbutil.ToProtoString(appVersion.Description)
	pbAppVersion.Status = pbutil.ToProtoString(appVersion.Status)
	pbAppVersion.PackageName = pbutil.ToProtoString(appVersion.PackageName)
	pbAppVersion.OwnerPath = appVersion.OwnerPath.ToProtoString()
	pbAppVersion.Owner = pbutil.ToProtoString(appVersion.Owner)
	pbAppVersion.CreateTime = pbutil.ToProtoTimestamp(appVersion.CreateTime)
	pbAppVersion.StatusTime = pbutil.ToProtoTimestamp(appVersion.StatusTime)
	pbAppVersion.Sequence = pbutil.ToProtoUInt32(appVersion.Sequence)
	pbAppVersion.Message = pbutil.ToProtoString(appVersion.Message)
	pbAppVersion.Type = pbutil.ToProtoString(appVersion.Type)
	pbAppVersion.ReviewId = pbutil.ToProtoString(appVersion.ReviewId)
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
