// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"time"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/sender"
	"openpitrix.io/openpitrix/pkg/util/idutil"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
)

func NewAppId() string {
	return idutil.GetUuid("app-")
}

type App struct {
	AppId       string
	Active      bool
	Name        string
	RepoId      string
	Description string
	Status      string
	Home        string
	Icon        string
	Screenshots string
	Maintainers string
	Keywords    string
	Sources     string
	Readme      string
	Owner       string
	OwnerPath   sender.OwnerPath
	ChartName   string
	Tos         string
	Abstraction string
	Isv         string
	CreateTime  time.Time
	StatusTime  time.Time
	UpdateTime  *time.Time
}

var AppColumns = db.GetColumnsFromStruct(&App{})

func NewApp(name string, ownerPath sender.OwnerPath, isv string) *App {
	return &App{
		AppId:      NewAppId(),
		Active:     false,
		Name:       name,
		Status:     constants.StatusDraft,
		Owner:      ownerPath.Owner(),
		OwnerPath:  ownerPath,
		Isv:        isv,
		CreateTime: time.Now(),
		StatusTime: time.Now(),
	}
}

func AppToPb(app *App) *pb.App {
	pbApp := pb.App{}
	pbApp.AppId = pbutil.ToProtoString(app.AppId)
	pbApp.Active = pbutil.ToProtoBool(app.Active)
	pbApp.Name = pbutil.ToProtoString(app.Name)
	pbApp.RepoId = pbutil.ToProtoString(app.RepoId)
	pbApp.Description = pbutil.ToProtoString(app.Description)
	pbApp.Status = pbutil.ToProtoString(app.Status)
	pbApp.Home = pbutil.ToProtoString(app.Home)
	pbApp.Icon = pbutil.ToProtoString(app.Icon)
	pbApp.Screenshots = pbutil.ToProtoString(app.Screenshots)
	pbApp.Maintainers = pbutil.ToProtoString(app.Maintainers)
	pbApp.Sources = pbutil.ToProtoString(app.Sources)
	pbApp.Readme = pbutil.ToProtoString(app.Readme)
	pbApp.ChartName = pbutil.ToProtoString(app.ChartName)
	pbApp.OwnerPath = app.OwnerPath.ToProtoString()
	pbApp.Owner = pbutil.ToProtoString(app.Owner)
	pbApp.Isv = pbutil.ToProtoString(app.Isv)
	pbApp.Keywords = pbutil.ToProtoString(app.Keywords)
	pbApp.Abstraction = pbutil.ToProtoString(app.Abstraction)
	pbApp.Tos = pbutil.ToProtoString(app.Tos)
	pbApp.CreateTime = pbutil.ToProtoTimestamp(app.CreateTime)
	pbApp.StatusTime = pbutil.ToProtoTimestamp(app.StatusTime)
	if app.UpdateTime != nil {
		pbApp.UpdateTime = pbutil.ToProtoTimestamp(*app.UpdateTime)
	}
	return &pbApp
}
