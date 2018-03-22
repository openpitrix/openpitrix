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
const RepoLabelTableName = "repo_label"
const RepoSelectorTableName = "repo_selector"

func NewRepoId() string {
	return utils.GetUuid("repo-")
}

func NewRepoLabelId() string {
	return utils.GetUuid("repol-")
}

func NewRepoSelectorId() string {
	return utils.GetUuid("repos-")
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

var RepoLabelColumns = GetColumnsFromStruct(&RepoLabel{})

func NewRepoLabel(repoId, labelKey, labelValue string) *RepoLabel {
	return &RepoLabel{
		RepoLabelId: NewRepoLabelId(),
		RepoId:      repoId,
		LabelKey:    labelKey,
		LabelValue:  labelValue,

		Status:     constants.StatusActive,
		CreateTime: time.Now(),
		StatusTime: time.Now(),
	}
}

func RepoLabelToPb(repoLabel *RepoLabel) *pb.RepoLabel {
	pbRepoLabel := pb.RepoLabel{}
	pbRepoLabel.RepoId = utils.ToProtoString(repoLabel.RepoId)
	pbRepoLabel.RepoLabelId = utils.ToProtoString(repoLabel.RepoLabelId)
	pbRepoLabel.LabelKey = utils.ToProtoString(repoLabel.LabelKey)
	pbRepoLabel.LabelValue = utils.ToProtoString(repoLabel.LabelValue)

	pbRepoLabel.Status = utils.ToProtoString(repoLabel.Status)
	pbRepoLabel.CreateTime = utils.ToProtoTimestamp(repoLabel.CreateTime)
	pbRepoLabel.StatusTime = utils.ToProtoTimestamp(repoLabel.StatusTime)
	return &pbRepoLabel
}

func RepoLabelsToPbs(repoLabels []*RepoLabel) (pbRepoLabels []*pb.RepoLabel) {
	for _, repoLabel := range repoLabels {
		pbRepoLabels = append(pbRepoLabels, RepoLabelToPb(repoLabel))
	}
	return
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

var RepoSelectorColumns = GetColumnsFromStruct(&RepoSelector{})

func NewRepoSelector(repoId, selectorKey, selectorValue string) *RepoSelector {
	return &RepoSelector{
		RepoSelectorId: NewRepoSelectorId(),
		RepoId:         repoId,
		SelectorKey:    selectorKey,
		SelectorValue:  selectorValue,

		Status:     constants.StatusActive,
		CreateTime: time.Now(),
		StatusTime: time.Now(),
	}
}

func RepoSelectorToPb(repoSelector *RepoSelector) *pb.RepoSelector {
	pbRepoSelector := pb.RepoSelector{}
	pbRepoSelector.RepoId = utils.ToProtoString(repoSelector.RepoId)
	pbRepoSelector.RepoSelectorId = utils.ToProtoString(repoSelector.RepoSelectorId)
	pbRepoSelector.SelectorKey = utils.ToProtoString(repoSelector.SelectorKey)
	pbRepoSelector.SelectorValue = utils.ToProtoString(repoSelector.SelectorValue)

	pbRepoSelector.Status = utils.ToProtoString(repoSelector.Status)
	pbRepoSelector.CreateTime = utils.ToProtoTimestamp(repoSelector.CreateTime)
	pbRepoSelector.StatusTime = utils.ToProtoTimestamp(repoSelector.StatusTime)
	return &pbRepoSelector
}

func RepoSelectorsToPbs(repoSelectors []*RepoSelector) (pbRepoSelectors []*pb.RepoSelector) {
	for _, repoSelector := range repoSelectors {
		pbRepoSelectors = append(pbRepoSelectors, RepoSelectorToPb(repoSelector))
	}
	return
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
