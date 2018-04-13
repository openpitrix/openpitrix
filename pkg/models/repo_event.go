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

const RepoEventTableName = "repo_event"

func NewRepoEventId() string {
	return utils.GetUuid("repoe-")
}

type RepoEvent struct {
	RepoEventId string
	RepoId      string
	Status      string
	Result      string
	Owner       string
	CreateTime  time.Time
	StatusTime  time.Time
}

var RepoEventColumns = GetColumnsFromStruct(&RepoEvent{})

func NewRepoEvent(repoId, owner string) *RepoEvent {
	return &RepoEvent{
		RepoEventId: NewRepoEventId(),
		RepoId:      repoId,
		Owner:       owner,
		Status:      constants.StatusPending,
		Result:      "",
		CreateTime:  time.Now(),
		StatusTime:  time.Now(),
	}
}

func RepoEventToPb(repoTask *RepoEvent) *pb.RepoEvent {
	pbRepoTask := pb.RepoEvent{}
	pbRepoTask.RepoEventId = utils.ToProtoString(repoTask.RepoEventId)
	pbRepoTask.RepoId = utils.ToProtoString(repoTask.RepoId)
	pbRepoTask.Status = utils.ToProtoString(repoTask.Status)
	pbRepoTask.Result = utils.ToProtoString(repoTask.Result)
	pbRepoTask.Owner = utils.ToProtoString(repoTask.Owner)
	pbRepoTask.CreateTime = utils.ToProtoTimestamp(repoTask.CreateTime)
	pbRepoTask.StatusTime = utils.ToProtoTimestamp(repoTask.StatusTime)
	return &pbRepoTask
}

func RepoEventsToPbs(repoTasks []*RepoEvent) (pbRepoTasks []*pb.RepoEvent) {
	for _, repoTask := range repoTasks {
		pbRepoTasks = append(pbRepoTasks, RepoEventToPb(repoTask))
	}
	return
}
