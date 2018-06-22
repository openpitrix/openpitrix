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

type SnapshotService struct {
	Config     *config.Config
	Properties *SnapshotServiceProperties
}

type SnapshotServiceProperties struct {
	// QingCloud Zone ID
	Zone *string `json:"zone" name:"zone"` // Required
}

func (s *QingCloudService) Snapshot(zone string) (*SnapshotService, error) {
	properties := &SnapshotServiceProperties{
		Zone: &zone,
	}

	return &SnapshotService{Config: s.Config, Properties: properties}, nil
}

// Documentation URL: https://docs.qingcloud.com/api/snapshot/apply_snapshots.html
func (s *SnapshotService) ApplySnapshots(i *ApplySnapshotsInput) (*ApplySnapshotsOutput, error) {
	if i == nil {
		i = &ApplySnapshotsInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "ApplySnapshots",
		RequestMethod: "GET",
	}

	x := &ApplySnapshotsOutput{}
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

type ApplySnapshotsInput struct {
	Snapshots []*string `json:"snapshots" name:"snapshots" location:"params"` // Required
}

func (v *ApplySnapshotsInput) Validate() error {

	if len(v.Snapshots) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "Snapshots",
			ParentName:    "ApplySnapshotsInput",
		}
	}

	return nil
}

type ApplySnapshotsOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	JobID   *string `json:"job_id" name:"job_id" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/snapshot/capture_instance_from_snapshot.html
func (s *SnapshotService) CaptureInstanceFromSnapshot(i *CaptureInstanceFromSnapshotInput) (*CaptureInstanceFromSnapshotOutput, error) {
	if i == nil {
		i = &CaptureInstanceFromSnapshotInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "CaptureInstanceFromSnapshot",
		RequestMethod: "GET",
	}

	x := &CaptureInstanceFromSnapshotOutput{}
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

type CaptureInstanceFromSnapshotInput struct {
	ImageName *string `json:"image_name" name:"image_name" location:"params"`
	Snapshot  *string `json:"snapshot" name:"snapshot" location:"params"` // Required
}

func (v *CaptureInstanceFromSnapshotInput) Validate() error {

	if v.Snapshot == nil {
		return errors.ParameterRequiredError{
			ParameterName: "Snapshot",
			ParentName:    "CaptureInstanceFromSnapshotInput",
		}
	}

	return nil
}

type CaptureInstanceFromSnapshotOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	ImageID *string `json:"image_id" name:"image_id" location:"elements"`
	JobID   *string `json:"job_id" name:"job_id" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/snapshot/create_snapshots.html
func (s *SnapshotService) CreateSnapshots(i *CreateSnapshotsInput) (*CreateSnapshotsOutput, error) {
	if i == nil {
		i = &CreateSnapshotsInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "CreateSnapshots",
		RequestMethod: "GET",
	}

	x := &CreateSnapshotsOutput{}
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

type CreateSnapshotsInput struct {

	// IsFull's available values: 0, 1
	IsFull        *int      `json:"is_full" name:"is_full" location:"params"`
	Resources     []*string `json:"resources" name:"resources" location:"params"` // Required
	ServiceParams *string   `json:"service_params" name:"service_params" location:"params"`
	SnapshotName  *string   `json:"snapshot_name" name:"snapshot_name" location:"params"`
}

func (v *CreateSnapshotsInput) Validate() error {

	if v.IsFull != nil {
		isFullValidValues := []string{"0", "1"}
		isFullParameterValue := fmt.Sprint(*v.IsFull)

		isFullIsValid := false
		for _, value := range isFullValidValues {
			if value == isFullParameterValue {
				isFullIsValid = true
			}
		}

		if !isFullIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "IsFull",
				ParameterValue: isFullParameterValue,
				AllowedValues:  isFullValidValues,
			}
		}
	}

	if len(v.Resources) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "Resources",
			ParentName:    "CreateSnapshotsInput",
		}
	}

	return nil
}

type CreateSnapshotsOutput struct {
	Message   *string   `json:"message" name:"message"`
	Action    *string   `json:"action" name:"action" location:"elements"`
	JobID     *string   `json:"job_id" name:"job_id" location:"elements"`
	RetCode   *int      `json:"ret_code" name:"ret_code" location:"elements"`
	Snapshots []*string `json:"snapshots" name:"snapshots" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/snapshot/create_volume_from_snapshot.html
func (s *SnapshotService) CreateVolumeFromSnapshot(i *CreateVolumeFromSnapshotInput) (*CreateVolumeFromSnapshotOutput, error) {
	if i == nil {
		i = &CreateVolumeFromSnapshotInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "CreateVolumeFromSnapshot",
		RequestMethod: "GET",
	}

	x := &CreateVolumeFromSnapshotOutput{}
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

type CreateVolumeFromSnapshotInput struct {
	Snapshot   *string `json:"snapshot" name:"snapshot" location:"params"` // Required
	VolumeName *string `json:"volume_name" name:"volume_name" location:"params"`
}

func (v *CreateVolumeFromSnapshotInput) Validate() error {

	if v.Snapshot == nil {
		return errors.ParameterRequiredError{
			ParameterName: "Snapshot",
			ParentName:    "CreateVolumeFromSnapshotInput",
		}
	}

	return nil
}

type CreateVolumeFromSnapshotOutput struct {
	Message  *string `json:"message" name:"message"`
	Action   *string `json:"action" name:"action" location:"elements"`
	JobID    *string `json:"job_id" name:"job_id" location:"elements"`
	RetCode  *int    `json:"ret_code" name:"ret_code" location:"elements"`
	VolumeID *string `json:"volume_id" name:"volume_id" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/snapshot/delete_snapshots.html
func (s *SnapshotService) DeleteSnapshots(i *DeleteSnapshotsInput) (*DeleteSnapshotsOutput, error) {
	if i == nil {
		i = &DeleteSnapshotsInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DeleteSnapshots",
		RequestMethod: "GET",
	}

	x := &DeleteSnapshotsOutput{}
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

type DeleteSnapshotsInput struct {
	Snapshots []*string `json:"snapshots" name:"snapshots" location:"params"` // Required
}

func (v *DeleteSnapshotsInput) Validate() error {

	if len(v.Snapshots) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "Snapshots",
			ParentName:    "DeleteSnapshotsInput",
		}
	}

	return nil
}

type DeleteSnapshotsOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	JobID   *string `json:"job_id" name:"job_id" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/snapshot/describe_snapshots.html
func (s *SnapshotService) DescribeSnapshots(i *DescribeSnapshotsInput) (*DescribeSnapshotsOutput, error) {
	if i == nil {
		i = &DescribeSnapshotsInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DescribeSnapshots",
		RequestMethod: "GET",
	}

	x := &DescribeSnapshotsOutput{}
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

type DescribeSnapshotsInput struct {
	Limit        *int    `json:"limit" name:"limit" default:"20" location:"params"`
	Offset       *int    `json:"offset" name:"offset" default:"0" location:"params"`
	ResourceID   *string `json:"resource_id" name:"resource_id" location:"params"`
	SearchWord   *string `json:"search_word" name:"search_word" location:"params"`
	SnapshotTime *string `json:"snapshot_time" name:"snapshot_time" location:"params"`
	// SnapshotType's available values: 0, 1
	SnapshotType *int      `json:"snapshot_type" name:"snapshot_type" location:"params"`
	Snapshots    []*string `json:"snapshots" name:"snapshots" location:"params"`
	Status       []*string `json:"status" name:"status" location:"params"`
	Tags         []*string `json:"tags" name:"tags" location:"params"`
	// Verbose's available values: 0, 1
	Verbose *int `json:"verbose" name:"verbose" default:"0" location:"params"`
}

func (v *DescribeSnapshotsInput) Validate() error {

	if v.SnapshotType != nil {
		snapshotTypeValidValues := []string{"0", "1"}
		snapshotTypeParameterValue := fmt.Sprint(*v.SnapshotType)

		snapshotTypeIsValid := false
		for _, value := range snapshotTypeValidValues {
			if value == snapshotTypeParameterValue {
				snapshotTypeIsValid = true
			}
		}

		if !snapshotTypeIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "SnapshotType",
				ParameterValue: snapshotTypeParameterValue,
				AllowedValues:  snapshotTypeValidValues,
			}
		}
	}

	if v.Verbose != nil {
		verboseValidValues := []string{"0", "1"}
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

type DescribeSnapshotsOutput struct {
	Message     *string     `json:"message" name:"message"`
	Action      *string     `json:"action" name:"action" location:"elements"`
	RetCode     *int        `json:"ret_code" name:"ret_code" location:"elements"`
	SnapshotSet []*Snapshot `json:"snapshot_set" name:"snapshot_set" location:"elements"`
	TotalCount  *int        `json:"total_count" name:"total_count" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/snapshot/modify_snapshot_attributes.html
func (s *SnapshotService) ModifySnapshotAttributes(i *ModifySnapshotAttributesInput) (*ModifySnapshotAttributesOutput, error) {
	if i == nil {
		i = &ModifySnapshotAttributesInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "ModifySnapshotAttributes",
		RequestMethod: "GET",
	}

	x := &ModifySnapshotAttributesOutput{}
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

type ModifySnapshotAttributesInput struct {
	Description  *string `json:"description" name:"description" location:"params"`
	Snapshot     *string `json:"snapshot" name:"snapshot" location:"params"` // Required
	SnapshotName *string `json:"snapshot_name" name:"snapshot_name" location:"params"`
}

func (v *ModifySnapshotAttributesInput) Validate() error {

	if v.Snapshot == nil {
		return errors.ParameterRequiredError{
			ParameterName: "Snapshot",
			ParentName:    "ModifySnapshotAttributesInput",
		}
	}

	return nil
}

type ModifySnapshotAttributesOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}
