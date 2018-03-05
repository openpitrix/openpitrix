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

const AppTableName = "app"

func NewAppId() string {
	return utils.GetUuid("app-")
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
	Sources     string
	Readme      string
	Owner       string
	ChartName   string
	CreateTime  time.Time
	StatusTime  time.Time
}

var AppColumns = GetColumnsFromStruct(&App{})

func NewApp(name, repoId, description, owner string) *App {
	return &App{
		AppId:       NewAppId(),
		Name:        name,
		RepoId:      repoId,
		Description: description,
		Status:      constants.StatusActive,
		Owner:       owner,
		CreateTime:  time.Now(),
		StatusTime:  time.Now(),
	}
}

func AppToPb(app *App) *pb.App {
	//logger.Infof("%+v", app.Home)
	pbApp := pb.App{}
	pbApp.AppId = utils.ToProtoString(app.AppId)
	pbApp.Name = utils.ToProtoString(app.Name)
	pbApp.RepoId = utils.ToProtoString(app.RepoId)
	pbApp.Description = utils.ToProtoString(app.Description)
	pbApp.Status = utils.ToProtoString(app.Status)
	pbApp.Home = utils.ToProtoString(app.Home)
	pbApp.Icon = utils.ToProtoString(app.Icon)
	pbApp.Screenshots = utils.ToProtoString(app.Screenshots)
	pbApp.Maintainers = utils.ToProtoString(app.Maintainers)
	pbApp.Sources = utils.ToProtoString(app.Sources)
	pbApp.Readme = utils.ToProtoString(app.Readme)
	pbApp.ChartName = utils.ToProtoString(app.ChartName)
	pbApp.Owner = utils.ToProtoString(app.Owner)
	pbApp.CreateTime = utils.ToProtoTimestamp(app.CreateTime)
	pbApp.StatusTime = utils.ToProtoTimestamp(app.StatusTime)
	return &pbApp
}

func AppsToPbs(apps []*App) (pbApps []*pb.App) {
	for _, app := range apps {
		pbApps = append(pbApps, AppToPb(app))
	}
	return
}
