// Copyright 2018 The OpenPitrix Authors. All rights reserved.
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

func NewRepoEventId() string {
	return idutil.GetUuid("repoe-")
}

type RepoEvent struct {
	RepoEventId string
	RepoId      string
	Status      string
	Result      string
	Owner       string
	OwnerPath   sender.OwnerPath
	CreateTime  time.Time
	StatusTime  time.Time
}

var RepoEventColumns = db.GetColumnsFromStruct(&RepoEvent{})

func NewRepoEvent(repoId string, ownerPath sender.OwnerPath) *RepoEvent {
	return &RepoEvent{
		RepoEventId: NewRepoEventId(),
		RepoId:      repoId,
		Owner:       ownerPath.Owner(),
		OwnerPath:   ownerPath,
		Status:      constants.StatusPending,
		Result:      "",
		CreateTime:  time.Now(),
		StatusTime:  time.Now(),
	}
}

func RepoEventToPb(repoTask *RepoEvent) *pb.RepoEvent {
	pbRepoTask := pb.RepoEvent{}
	pbRepoTask.RepoEventId = pbutil.ToProtoString(repoTask.RepoEventId)
	pbRepoTask.RepoId = pbutil.ToProtoString(repoTask.RepoId)
	pbRepoTask.Status = pbutil.ToProtoString(repoTask.Status)
	pbRepoTask.Result = pbutil.ToProtoString(repoTask.Result)
	pbRepoTask.OwnerPath = repoTask.OwnerPath.ToProtoString()
	pbRepoTask.CreateTime = pbutil.ToProtoTimestamp(repoTask.CreateTime)
	pbRepoTask.StatusTime = pbutil.ToProtoTimestamp(repoTask.StatusTime)
	return &pbRepoTask
}

func RepoEventsToPbs(repoTasks []*RepoEvent) (pbRepoTasks []*pb.RepoEvent) {
	for _, repoTask := range repoTasks {
		pbRepoTasks = append(pbRepoTasks, RepoEventToPb(repoTask))
	}
	return
}
