// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package job

import (
	"context"
	"fmt"
	"time"

	"openpitrix.io/openpitrix/pkg/client"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/utils"
)

func NewJobManagerClient(ctx context.Context) (pb.JobManagerClient, error) {
	conn, err := manager.NewClient(ctx, constants.JobManagerHost, constants.JobManagerPort)
	if err != nil {
		return nil, err
	}
	return pb.NewJobManagerClient(conn), err
}

func CreateJob(ctx context.Context, jobRequest *pb.CreateJobRequest) (jobId string, err error) {
	jobManagerClient, err := NewJobManagerClient(ctx)
	if err != nil {
		return
	}
	jobResponse, err := jobManagerClient.CreateJob(ctx, jobRequest)
	if err != nil {
		return
	}
	jobId = jobResponse.GetJobId().GetValue()
	return
}

func DescribeJobs(ctx context.Context, jobRequest *pb.DescribeJobsRequest) (*pb.DescribeJobsResponse, error) {
	jobManagerClient, err := NewJobManagerClient(ctx)
	if err != nil {
		return nil, err
	}
	JobResponse, err := jobManagerClient.DescribeJobs(ctx, jobRequest)
	if err != nil {
		return nil, err
	}
	return JobResponse, err
}

func WaitJob(jobId string, timeout time.Duration, waitInterval time.Duration) error {
	logger.Debug("Waiting for job [%s] finished", jobId)
	return utils.WaitForSpecificOrError(func() (bool, error) {
		jobRequest := &pb.DescribeJobsRequest{
			JobId: []string{jobId},
		}
		jobResponse, err := DescribeJobs(client.GetSystemUserContext(), jobRequest)
		if err != nil {
			//network or api error, not considered job fail.
			return false, nil
		}
		if len(jobResponse.JobSet) == 0 {
			return false, fmt.Errorf("Can not find job [%s]. ", jobId)
		}
		j := jobResponse.JobSet[0]
		if j.Status == nil {
			logger.Errorf("Job [%s] status is nil", jobId)
			return false, nil
		}
		if j.Status.GetValue() == constants.StatusWorking || j.Status.GetValue() == constants.StatusPending {
			return false, nil
		}
		if j.Status.GetValue() == constants.StatusSuccessful {
			return true, nil
		}
		if j.Status.GetValue() == constants.StatusFailed {
			return false, fmt.Errorf("Job [%s] failed. ", jobId)
		}
		logger.Errorf("Unknown status [%s] for job [%s]. ", j.Status.GetValue(), jobId)
		return false, nil
	}, timeout, waitInterval)
}

func SendJob(job *models.Job) (jobId string, err error) {
	pbJob := models.JobToPb(job)
	jobRequest := &pb.CreateJobRequest{
		ClusterId: pbJob.ClusterId,
		AppId:     pbJob.AppId,
		VersionId: pbJob.VersionId,
		JobAction: pbJob.JobAction,
		Provider:  pbJob.Provider,
		Directive: pbJob.Directive,
	}
	jobId, err = CreateJob(client.GetSystemUserContext(), jobRequest)
	if err != nil {
		logger.Errorf("Failed to create job [%s]: %+v", jobId, err)
	}
	return
}
