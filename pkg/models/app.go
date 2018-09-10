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

const AppTableName = "app"

func NewAppId() string {
	return idutil.GetUuid("app-")
}

type App struct {
	AppId       string
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
	ChartName   string
	CreateTime  time.Time
	StatusTime  time.Time
	UpdateTime  *time.Time
}

var AppColumns = GetColumnsFromStruct(&App{})

func NewApp(name, repoId, description, owner, chartName string) *App {
	return &App{
		AppId:       NewAppId(),
		Name:        name,
		RepoId:      repoId,
		Description: description,
		Status:      constants.StatusDraft,
		Owner:       owner,
		ChartName:   chartName,
		CreateTime:  time.Now(),
		StatusTime:  time.Now(),
	}
}

func AppToPb(app *App) *pb.App {
	pbApp := pb.App{}
	pbApp.AppId = pbutil.ToProtoString(app.AppId)
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
	pbApp.Owner = pbutil.ToProtoString(app.Owner)
	pbApp.Keywords = pbutil.ToProtoString(app.Keywords)
	pbApp.CreateTime = pbutil.ToProtoTimestamp(app.CreateTime)
	pbApp.StatusTime = pbutil.ToProtoTimestamp(app.StatusTime)
	if app.UpdateTime != nil {
		pbApp.UpdateTime = pbutil.ToProtoTimestamp(*app.UpdateTime)
	}
	return &pbApp
}
