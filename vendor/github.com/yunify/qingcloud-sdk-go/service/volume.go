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

type VolumeService struct {
	Config     *config.Config
	Properties *VolumeServiceProperties
}

type VolumeServiceProperties struct {
	// QingCloud Zone ID
	Zone *string `json:"zone" name:"zone"` // Required
}

func (s *QingCloudService) Volume(zone string) (*VolumeService, error) {
	properties := &VolumeServiceProperties{
		Zone: &zone,
	}

	return &VolumeService{Config: s.Config, Properties: properties}, nil
}

// Documentation URL: https://docs.qingcloud.com/api/volume/attach_volumes.html
func (s *VolumeService) AttachVolumes(i *AttachVolumesInput) (*AttachVolumesOutput, error) {
	if i == nil {
		i = &AttachVolumesInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "AttachVolumes",
		RequestMethod: "GET",
	}

	x := &AttachVolumesOutput{}
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

type AttachVolumesInput struct {
	Instance *string   `json:"instance" name:"instance" location:"params"` // Required
	Volumes  []*string `json:"volumes" name:"volumes" location:"params"`   // Required
}

func (v *AttachVolumesInput) Validate() error {

	if v.Instance == nil {
		return errors.ParameterRequiredError{
			ParameterName: "Instance",
			ParentName:    "AttachVolumesInput",
		}
	}

	if len(v.Volumes) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "Volumes",
			ParentName:    "AttachVolumesInput",
		}
	}

	return nil
}

type AttachVolumesOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	JobID   *string `json:"job_id" name:"job_id" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/volume/create_volumes.html
func (s *VolumeService) CreateVolumes(i *CreateVolumesInput) (*CreateVolumesOutput, error) {
	if i == nil {
		i = &CreateVolumesInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "CreateVolumes",
		RequestMethod: "GET",
	}

	x := &CreateVolumesOutput{}
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

type CreateVolumesInput struct {
	Count      *int    `json:"count" name:"count" default:"1" location:"params"`
	Size       *int    `json:"size" name:"size" location:"params"` // Required
	VolumeName *string `json:"volume_name" name:"volume_name" location:"params"`
	// VolumeType's available values: 0, 1, 2, 3
	VolumeType *int `json:"volume_type" name:"volume_type" default:"0" location:"params"`
}

func (v *CreateVolumesInput) Validate() error {

	if v.Size == nil {
		return errors.ParameterRequiredError{
			ParameterName: "Size",
			ParentName:    "CreateVolumesInput",
		}
	}

	if v.VolumeType != nil {
		volumeTypeValidValues := []string{"0", "1", "2", "3"}
		volumeTypeParameterValue := fmt.Sprint(*v.VolumeType)

		volumeTypeIsValid := false
		for _, value := range volumeTypeValidValues {
			if value == volumeTypeParameterValue {
				volumeTypeIsValid = true
			}
		}

		if !volumeTypeIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "VolumeType",
				ParameterValue: volumeTypeParameterValue,
				AllowedValues:  volumeTypeValidValues,
			}
		}
	}

	return nil
}

type CreateVolumesOutput struct {
	Message *string   `json:"message" name:"message"`
	Action  *string   `json:"action" name:"action" location:"elements"`
	JobID   *string   `json:"job_id" name:"job_id" location:"elements"`
	RetCode *int      `json:"ret_code" name:"ret_code" location:"elements"`
	Volumes []*string `json:"volumes" name:"volumes" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/volume/delete_volumes.html
func (s *VolumeService) DeleteVolumes(i *DeleteVolumesInput) (*DeleteVolumesOutput, error) {
	if i == nil {
		i = &DeleteVolumesInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DeleteVolumes",
		RequestMethod: "GET",
	}

	x := &DeleteVolumesOutput{}
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

type DeleteVolumesInput struct {
	Volumes []*string `json:"volumes" name:"volumes" location:"params"` // Required
}

func (v *DeleteVolumesInput) Validate() error {

	if len(v.Volumes) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "Volumes",
			ParentName:    "DeleteVolumesInput",
		}
	}

	return nil
}

type DeleteVolumesOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	JobID   *string `json:"job_id" name:"job_id" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/volume/describe_volumes.html
func (s *VolumeService) DescribeVolumes(i *DescribeVolumesInput) (*DescribeVolumesOutput, error) {
	if i == nil {
		i = &DescribeVolumesInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DescribeVolumes",
		RequestMethod: "GET",
	}

	x := &DescribeVolumesOutput{}
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

type DescribeVolumesInput struct {
	Limit      *int      `json:"limit" name:"limit" default:"20" location:"params"`
	Offset     *int      `json:"offset" name:"offset" default:"0" location:"params"`
	SearchWord *string   `json:"search_word" name:"search_word" location:"params"`
	Status     []*string `json:"status" name:"status" location:"params"`
	Tags       []*string `json:"tags" name:"tags" location:"params"`
	// Verbose's available values: 0, 1
	Verbose *int `json:"verbose" name:"verbose" default:"0" location:"params"`
	// VolumeType's available values: 0, 1, 2, 3
	VolumeType *int      `json:"volume_type" name:"volume_type" location:"params"`
	Volumes    []*string `json:"volumes" name:"volumes" location:"params"`
}

func (v *DescribeVolumesInput) Validate() error {

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

	if v.VolumeType != nil {
		volumeTypeValidValues := []string{"0", "1", "2", "3"}
		volumeTypeParameterValue := fmt.Sprint(*v.VolumeType)

		volumeTypeIsValid := false
		for _, value := range volumeTypeValidValues {
			if value == volumeTypeParameterValue {
				volumeTypeIsValid = true
			}
		}

		if !volumeTypeIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "VolumeType",
				ParameterValue: volumeTypeParameterValue,
				AllowedValues:  volumeTypeValidValues,
			}
		}
	}

	return nil
}

type DescribeVolumesOutput struct {
	Message    *string   `json:"message" name:"message"`
	Action     *string   `json:"action" name:"action" location:"elements"`
	RetCode    *int      `json:"ret_code" name:"ret_code" location:"elements"`
	TotalCount *int      `json:"total_count" name:"total_count" location:"elements"`
	VolumeSet  []*Volume `json:"volume_set" name:"volume_set" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/volume/detach_volumes.html
func (s *VolumeService) DetachVolumes(i *DetachVolumesInput) (*DetachVolumesOutput, error) {
	if i == nil {
		i = &DetachVolumesInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DetachVolumes",
		RequestMethod: "GET",
	}

	x := &DetachVolumesOutput{}
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

type DetachVolumesInput struct {
	Instance *string   `json:"instance" name:"instance" location:"params"` // Required
	Volumes  []*string `json:"volumes" name:"volumes" location:"params"`   // Required
}

func (v *DetachVolumesInput) Validate() error {

	if v.Instance == nil {
		return errors.ParameterRequiredError{
			ParameterName: "Instance",
			ParentName:    "DetachVolumesInput",
		}
	}

	if len(v.Volumes) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "Volumes",
			ParentName:    "DetachVolumesInput",
		}
	}

	return nil
}

type DetachVolumesOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	JobID   *string `json:"job_id" name:"job_id" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/volume/modify_volume_attributes.html
func (s *VolumeService) ModifyVolumeAttributes(i *ModifyVolumeAttributesInput) (*ModifyVolumeAttributesOutput, error) {
	if i == nil {
		i = &ModifyVolumeAttributesInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "ModifyVolumeAttributes",
		RequestMethod: "GET",
	}

	x := &ModifyVolumeAttributesOutput{}
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

type ModifyVolumeAttributesInput struct {
	Description *string `json:"description" name:"description" location:"params"`
	Volume      *string `json:"volume" name:"volume" location:"params"` // Required
	VolumeName  *string `json:"volume_name" name:"volume_name" location:"params"`
}

func (v *ModifyVolumeAttributesInput) Validate() error {

	if v.Volume == nil {
		return errors.ParameterRequiredError{
			ParameterName: "Volume",
			ParentName:    "ModifyVolumeAttributesInput",
		}
	}

	return nil
}

type ModifyVolumeAttributesOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/volume/resize_volumes.html
func (s *VolumeService) ResizeVolumes(i *ResizeVolumesInput) (*ResizeVolumesOutput, error) {
	if i == nil {
		i = &ResizeVolumesInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "ResizeVolumes",
		RequestMethod: "GET",
	}

	x := &ResizeVolumesOutput{}
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

type ResizeVolumesInput struct {
	Size    *int      `json:"size" name:"size" location:"params"`       // Required
	Volumes []*string `json:"volumes" name:"volumes" location:"params"` // Required
}

func (v *ResizeVolumesInput) Validate() error {

	if v.Size == nil {
		return errors.ParameterRequiredError{
			ParameterName: "Size",
			ParentName:    "ResizeVolumesInput",
		}
	}

	if len(v.Volumes) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "Volumes",
			ParentName:    "ResizeVolumesInput",
		}
	}

	return nil
}

type ResizeVolumesOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	JobID   *string `json:"job_id" name:"job_id" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}
