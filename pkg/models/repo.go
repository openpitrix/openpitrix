// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"time"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/util/idutil"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
)

func NewRepoId() string {
	return idutil.GetUuid("repo-")
}

type Repo struct {
	RepoId           string
	Name             string
	Description      string
	Type             string
	Url              string
	Credential       string
	Visibility       string
	Owner            string
	AppDefaultStatus string
	Controller       int8

	Status     string
	CreateTime time.Time
	StatusTime time.Time
}

var RepoColumns = db.GetColumnsFromStruct(&Repo{})
var RepoColumnsWithTablePrefix = db.GetColumnsFromStructWithPrefix(constants.TableRepo, &Repo{})

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
	pbRepo.RepoId = pbutil.ToProtoString(repo.RepoId)
	pbRepo.Name = pbutil.ToProtoString(repo.Name)
	pbRepo.Description = pbutil.ToProtoString(repo.Description)
	pbRepo.Type = pbutil.ToProtoString(repo.Type)
	pbRepo.Url = pbutil.ToProtoString(repo.Url)
	pbRepo.Credential = pbutil.ToProtoString(repo.Credential)
	pbRepo.Visibility = pbutil.ToProtoString(repo.Visibility)
	pbRepo.Owner = pbutil.ToProtoString(repo.Owner)
	pbRepo.Status = pbutil.ToProtoString(repo.Status)
	pbRepo.CreateTime = pbutil.ToProtoTimestamp(repo.CreateTime)
	pbRepo.StatusTime = pbutil.ToProtoTimestamp(repo.StatusTime)
	pbRepo.AppDefaultStatus = pbutil.ToProtoString(repo.AppDefaultStatus)
	pbRepo.Controller = pbutil.ToProtoInt32(int32(repo.Controller))
	return &pbRepo
}

func ReposToPbs(repos []*Repo) (pbRepos []*pb.Repo) {
	for _, repo := range repos {
		pbRepos = append(pbRepos, RepoToPb(repo))
	}
	return
}
