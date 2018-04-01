// Copyright 2018 The OpenPitrix Authors. All rights reserved.
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

type Repo struct {
	RepoId      string
	Name        string
	Description string
	Type        string
	Url         string
	Credential  string
	Visibility  string
	Owner       string

	Status     string
	CreateTime time.Time
	StatusTime time.Time
}

var RepoColumns = GetColumnsFromStruct(&Repo{})
var RepoColumnsWithTablePrefix = GetColumnsFromStructWithPrefix(RepoTableName, &Repo{})

func NewRepo(name, description, typ, url, credential, visibility, owner string) *Repo {
	return &Repo{
		RepoId:      NewRepoId(),
		Name:        name,
		Description: description,
		Type:        typ,
		Url:         url,
		Credential:  credential,
		Visibility:  visibility,
		Owner:       owner,
		Status:      constants.StatusActive,
		CreateTime:  time.Now(),
		StatusTime:  time.Now(),
	}
}

func RepoToPb(repo *Repo) *pb.Repo {
	pbRepo := pb.Repo{}
	pbRepo.RepoId = utils.ToProtoString(repo.RepoId)
	pbRepo.Name = utils.ToProtoString(repo.Name)
	pbRepo.Description = utils.ToProtoString(repo.Description)
	pbRepo.Type = utils.ToProtoString(repo.Type)
	pbRepo.Url = utils.ToProtoString(repo.Url)
	pbRepo.Credential = utils.ToProtoString(repo.Credential)
	pbRepo.Visibility = utils.ToProtoString(repo.Visibility)
	pbRepo.Owner = utils.ToProtoString(repo.Owner)
	pbRepo.Status = utils.ToProtoString(repo.Status)
	pbRepo.CreateTime = utils.ToProtoTimestamp(repo.CreateTime)
	pbRepo.StatusTime = utils.ToProtoTimestamp(repo.StatusTime)
	return &pbRepo
}

func ReposToPbs(repos []*Repo) (pbRepos []*pb.Repo) {
	for _, repo := range repos {
		pbRepos = append(pbRepos, RepoToPb(repo))
	}
	return
}
