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

func NewJobId() string {
	return idutil.GetUuid("j-")
}

type Job struct {
	JobId      string
	ClusterId  string
	AppId      string
	VersionId  string
	JobAction  string
	Directive  string
	Provider   string
	Owner      string
	OwnerPath  sender.OwnerPath
	Status     string
	ErrorCode  uint32
	Executor   string
	RuntimeId  string
	TaskCount  uint32
	CreateTime time.Time
	StatusTime time.Time
}

var JobColumns = db.GetColumnsFromStruct(&Job{})

func NewJob(jobId, clusterId, appId, versionId, jobAction, directive, provider string, ownerPath sender.OwnerPath, runtimeId string) *Job {
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
		Provider:   provider,
		Owner:      ownerPath.Owner(),
		OwnerPath:  ownerPath,
		RuntimeId:  runtimeId,
		Status:     constants.StatusPending,
		CreateTime: time.Now(),
		StatusTime: time.Now(),
	}
}

func JobToPb(job *Job) *pb.Job {
	pbJob := pb.Job{}
	pbJob.JobId = pbutil.ToProtoString(job.JobId)
	pbJob.ClusterId = pbutil.ToProtoString(job.ClusterId)
	pbJob.AppId = pbutil.ToProtoString(job.AppId)
	pbJob.VersionId = pbutil.ToProtoString(job.VersionId)
	pbJob.JobAction = pbutil.ToProtoString(job.JobAction)
	pbJob.Directive = pbutil.ToProtoString(job.Directive)
	pbJob.Provider = pbutil.ToProtoString(job.Provider)
	pbJob.OwnerPath = job.OwnerPath.ToProtoString()
	pbJob.Owner = pbutil.ToProtoString(job.Owner)
	pbJob.Status = pbutil.ToProtoString(job.Status)
	pbJob.ErrorCode = pbutil.ToProtoUInt32(job.ErrorCode)
	pbJob.Executor = pbutil.ToProtoString(job.Executor)
	pbJob.RuntimeId = pbutil.ToProtoString(job.RuntimeId)
	pbJob.TaskCount = pbutil.ToProtoUInt32(job.TaskCount)
	pbJob.CreateTime = pbutil.ToProtoTimestamp(job.CreateTime)
	pbJob.StatusTime = pbutil.ToProtoTimestamp(job.StatusTime)
	return &pbJob
}

func JobsToPbs(jobs []*Job) (pbJobs []*pb.Job) {
	for _, job := range jobs {
		pbJobs = append(pbJobs, JobToPb(job))
	}
	return
}

func PbToJob(pbJob *pb.Job) *Job {
	ownerPath := sender.OwnerPath(pbJob.GetOwnerPath().GetValue())
	return &Job{
		JobId:      pbJob.GetJobId().GetValue(),
		ClusterId:  pbJob.GetClusterId().GetValue(),
		AppId:      pbJob.GetAppId().GetValue(),
		VersionId:  pbJob.GetVersionId().GetValue(),
		JobAction:  pbJob.GetJobAction().GetValue(),
		Directive:  pbJob.GetDirective().GetValue(),
		Provider:   pbJob.GetProvider().GetValue(),
		OwnerPath:  ownerPath,
		Owner:      ownerPath.Owner(),
		Status:     pbJob.GetStatus().GetValue(),
		ErrorCode:  pbJob.GetErrorCode().GetValue(),
		Executor:   pbJob.GetExecutor().GetValue(),
		RuntimeId:  pbJob.GetRuntimeId().GetValue(),
		TaskCount:  pbJob.GetTaskCount().GetValue(),
		CreateTime: pbutil.GetTime(pbJob.GetCreateTime()),
		StatusTime: pbutil.GetTime(pbJob.GetStatusTime()),
	}
}
