// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package job

import (
	"openpitrix.io/openpitrix/pkg/client"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
)

type Client struct {
	pb.JobManagerClient
}

func NewClient() (*Client, error) {
	conn, err := manager.NewClient(constants.JobManagerHost, constants.JobManagerPort)
	if err != nil {
		return nil, err
	}
	return &Client{
		JobManagerClient: pb.NewJobManagerClient(conn),
	}, nil
}

func SendJob(job *models.Job) (string, error) {
	pbJob := models.JobToPb(job)
	jobRequest := &pb.CreateJobRequest{
		ClusterId: pbJob.ClusterId,
		AppId:     pbJob.AppId,
		VersionId: pbJob.VersionId,
		JobAction: pbJob.JobAction,
		Provider:  pbJob.Provider,
		Directive: pbJob.Directive,
	}

	ctx := client.GetSystemUserContext()
	jobClient, err := NewClient()
	if err != nil {
		logger.Error("Connect to job service failed: %+v", err)
		return "", err
	}
	response, err := jobClient.CreateJob(ctx, jobRequest)
	jobId := response.GetJobId().GetValue()
	if err != nil {
		logger.Error("Failed to create job [%s]: %+v", jobId, err)
		return "", err
	}
	return jobId, nil
}
