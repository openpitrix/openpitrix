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

const RepoTableName = "repo"

func NewRepoId() string {
	return utils.GetUuid("repo-")
}

type RepoLabel struct {
	RepoLabelId string
	RepoId      string
	LabelKey    string
	LabelValue  string

	Status     string
	CreateTime time.Time
	StatusTime time.Time
}

type RepoSelector struct {
	RepoSelectorId string
	RepoId         string
	SelectorKey    string
	SelectorValue  string

	Status     string
	CreateTime time.Time
	StatusTime time.Time
}

type Repo struct {
	RepoId      string
	Name        string
	Description string
	Url         string
	Credential  string
	Visibility  string
	Owner       string

	Status     string
	CreateTime time.Time
	StatusTime time.Time
}

var RepoColumns = GetColumnsFromStruct(&Repo{})

func NewRepo(name, description, url, credential, visibility, owner string) *Repo {
	return &Repo{
		RepoId:      NewRepoId(),
		Name:        name,
		Description: description,
		Url:         url,
		Credential:  credential,
		Visibility:  visibility,
		Owner:       owner,
		Status:      constants.StatusActive,
		CreateTime:  time.Now(),
		StatusTime:  time.Now(),
	}
}

func RepoToPb(app *Repo) *pb.Repo {
	pbRepo := pb.Repo{}
	pbRepo.RepoId = utils.ToProtoString(app.RepoId)
	pbRepo.Name = utils.ToProtoString(app.Name)
	pbRepo.Description = utils.ToProtoString(app.Description)
	pbRepo.Url = utils.ToProtoString(app.Url)
	pbRepo.Credential = utils.ToProtoString(app.Credential)
	pbRepo.Visibility = utils.ToProtoString(app.Visibility)
	pbRepo.Owner = utils.ToProtoString(app.Owner)
	pbRepo.Status = utils.ToProtoString(app.Status)
	pbRepo.CreateTime = utils.ToProtoTimestamp(app.CreateTime)
	pbRepo.StatusTime = utils.ToProtoTimestamp(app.StatusTime)
	return &pbRepo
}

func ReposToPbs(apps []*Repo) (pbRepos []*pb.Repo) {
	for _, app := range apps {
		pbRepos = append(pbRepos, RepoToPb(app))
	}
	return
}
