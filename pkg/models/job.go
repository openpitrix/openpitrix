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

const JobTableName = "job"

func NewJobId() string {
	return utils.GetUuid("j-")
}

type Job struct {
	JobId      string
	ClusterId  string
	AppId      string
	VersionId  string
	JobAction  string
	Directive  string
	Runtime    string
	Owner      string
	Status     string
	ErrorCode  uint32
	Executor   string
	TaskCount  uint32
	CreateTime time.Time
	StatusTime time.Time
}

var JobColumns = GetColumnsFromStruct(&Job{})

func NewJob(jobId, clusterId, appId, versionId, jobAction, directive, runtime, userId string) *Job {
	if jobId == "" {
		jobId = NewJobId()
	} else if jobId == constants.PlaceHolder {
		jobId = ""
	}
	return &Job{
		JobId:      jobId,
		ClusterId:  clusterId,
		AppId:      appId,
		VersionId:  versionId,
		JobAction:  jobAction,
		Directive:  directive,
		Runtime:    runtime,
		Owner:      userId,
		Status:     constants.StatusPending,
		CreateTime: time.Now(),
		StatusTime: time.Now(),
	}
}

func JobToPb(job *Job) *pb.Job {
	pbJob := pb.Job{}
	pbJob.JobId = utils.ToProtoString(job.JobId)
	pbJob.ClusterId = utils.ToProtoString(job.ClusterId)
	pbJob.AppId = utils.ToProtoString(job.AppId)
	pbJob.VersionId = utils.ToProtoString(job.VersionId)
	pbJob.JobAction = utils.ToProtoString(job.JobAction)
	pbJob.Directive = utils.ToProtoString(job.Directive)
	pbJob.Runtime = utils.ToProtoString(job.Runtime)
	pbJob.Owner = utils.ToProtoString(job.Owner)
	pbJob.Status = utils.ToProtoString(job.Status)
	pbJob.ErrorCode = utils.ToProtoUInt32(job.ErrorCode)
	pbJob.Executor = utils.ToProtoString(job.Executor)
	pbJob.TaskCount = utils.ToProtoUInt32(job.TaskCount)
	pbJob.CreateTime = utils.ToProtoTimestamp(job.CreateTime)
	pbJob.StatusTime = utils.ToProtoTimestamp(job.StatusTime)
	return &pbJob
}

func JobsToPbs(jobs []*Job) (pbJobs []*pb.Job) {
	for _, job := range jobs {
		pbJobs = append(pbJobs, JobToPb(job))
	}
	return
}
