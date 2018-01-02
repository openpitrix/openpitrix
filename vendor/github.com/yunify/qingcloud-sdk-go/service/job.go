// +-------------------------------------------------------------------------
// | Copyright (C) 2016 Yunify, Inc.
// +-------------------------------------------------------------------------
// | Licensed under the Apache License, Version 2.0 (the "License");
// | you may not use this work except in compliance with the License.
// | You may obtain a copy of the License in the LICENSE file, or at:
// |
// | http://www.apache.org/licenses/LICENSE-2.0
// |
// | Unless required by applicable law or agreed to in writing, software
// | distributed under the License is distributed on an "AS IS" BASIS,
// | WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// | See the License for the specific language governing permissions and
// | limitations under the License.
// +-------------------------------------------------------------------------

package service

import (
	"fmt"
	"time"

	"github.com/yunify/qingcloud-sdk-go/config"
	"github.com/yunify/qingcloud-sdk-go/request"
	"github.com/yunify/qingcloud-sdk-go/request/data"
	"github.com/yunify/qingcloud-sdk-go/request/errors"
)

var _ fmt.State
var _ time.Time

type JobService struct {
	Config     *config.Config
	Properties *JobServiceProperties
}

type JobServiceProperties struct {
	// QingCloud Zone ID
	Zone *string `json:"zone" name:"zone"` // Required
}

func (s *QingCloudService) Job(zone string) (*JobService, error) {
	properties := &JobServiceProperties{
		Zone: &zone,
	}

	return &JobService{Config: s.Config, Properties: properties}, nil
}

// Documentation URL: https://docs.qingcloud.com/api/job/describe_jobs.html
func (s *JobService) DescribeJobs(i *DescribeJobsInput) (*DescribeJobsOutput, error) {
	if i == nil {
		i = &DescribeJobsInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DescribeJobs",
		RequestMethod: "GET",
	}

	x := &DescribeJobsOutput{}
	r, err := request.New(o, i, x)
	if err != nil {
		return nil, err
	}

	err = r.Send()
	if err != nil {
		return nil, err
	}

	return x, err
}

type DescribeJobsInput struct {
	Jobs   []*string `json:"jobs" name:"jobs" location:"params"`
	Limit  *int      `json:"limit" name:"limit" default:"20" location:"params"`
	Offset *int      `json:"offset" name:"offset" default:"0" location:"params"`
	Status []*string `json:"status" name:"status" location:"params"`
	// Verbose's available values: 0
	Verbose *int `json:"verbose" name:"verbose" default:"0" location:"params"`
}

func (v *DescribeJobsInput) Validate() error {

	if v.Verbose != nil {
		verboseValidValues := []string{"0"}
		verboseParameterValue := fmt.Sprint(*v.Verbose)

		verboseIsValid := false
		for _, value := range verboseValidValues {
			if value == verboseParameterValue {
				verboseIsValid = true
			}
		}

		if !verboseIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "Verbose",
				ParameterValue: verboseParameterValue,
				AllowedValues:  verboseValidValues,
			}
		}
	}

	return nil
}

type DescribeJobsOutput struct {
	Message    *string `json:"message" name:"message"`
	Action     *string `json:"action" name:"action" location:"elements"`
	JobSet     []*Job  `json:"job_set" name:"job_set" location:"elements"`
	RetCode    *int    `json:"ret_code" name:"ret_code" location:"elements"`
	TotalCount *int    `json:"total_count" name:"total_count" location:"elements"`
}
