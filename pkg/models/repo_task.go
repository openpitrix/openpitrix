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

const RepoTaskTableName = "repo_task"

func NewRepoTaskId() string {
	return utils.GetUuid("rtask-")
}

type RepoTask struct {
	RepoTaskId string
	RepoId     string
	Status     string
	Result     string
	Owner      string
	CreateTime time.Time
	StatusTime time.Time
}

var RepoTaskColumns = GetColumnsFromStruct(&RepoTask{})

func NewRepoTask(repoId, owner string) *RepoTask {
	return &RepoTask{
		RepoTaskId: NewRepoTaskId(),
		RepoId:     repoId,
		Owner:      owner,
		Status:     constants.StatusPending,
		Result:     "",
		CreateTime: time.Now(),
		StatusTime: time.Now(),
	}
}

func RepoTaskToPb(repoTask *RepoTask) *pb.RepoTask {
	pbRepoTask := pb.RepoTask{}
	pbRepoTask.RepoTaskId = utils.ToProtoString(repoTask.RepoTaskId)
	pbRepoTask.RepoId = utils.ToProtoString(repoTask.RepoId)
	pbRepoTask.Status = utils.ToProtoString(repoTask.Status)
	pbRepoTask.Result = utils.ToProtoString(repoTask.Result)
	pbRepoTask.Owner = utils.ToProtoString(repoTask.Owner)
	pbRepoTask.CreateTime = utils.ToProtoTimestamp(repoTask.CreateTime)
	pbRepoTask.StatusTime = utils.ToProtoTimestamp(repoTask.StatusTime)
	return &pbRepoTask
}

func RepoTasksToPbs(repoTasks []*RepoTask) (pbRepoTasks []*pb.RepoTask) {
	for _, repoTask := range repoTasks {
		pbRepoTasks = append(pbRepoTasks, RepoTaskToPb(repoTask))
	}
	return
}
