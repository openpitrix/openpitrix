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

type VxNetService struct {
	Config     *config.Config
	Properties *VxNetServiceProperties
}

type VxNetServiceProperties struct {
	// QingCloud Zone ID
	Zone *string `json:"zone" name:"zone"` // Required
}

func (s *QingCloudService) VxNet(zone string) (*VxNetService, error) {
	properties := &VxNetServiceProperties{
		Zone: &zone,
	}

	return &VxNetService{Config: s.Config, Properties: properties}, nil
}

// Documentation URL: https://docs.qingcloud.com/api/vxnet/create_vxnets.html
func (s *VxNetService) CreateVxNets(i *CreateVxNetsInput) (*CreateVxNetsOutput, error) {
	if i == nil {
		i = &CreateVxNetsInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "CreateVxnets",
		RequestMethod: "GET",
	}

	x := &CreateVxNetsOutput{}
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

type CreateVxNetsInput struct {
	Count     *int    `json:"count" name:"count" default:"1" location:"params"`
	VxNetName *string `json:"vxnet_name" name:"vxnet_name" location:"params"`
	// VxNetType's available values: 0, 1
	VxNetType *int `json:"vxnet_type" name:"vxnet_type" location:"params"` // Required
}

func (v *CreateVxNetsInput) Validate() error {

	if v.VxNetType == nil {
		return errors.ParameterRequiredError{
			ParameterName: "VxNetType",
			ParentName:    "CreateVxNetsInput",
		}
	}

	if v.VxNetType != nil {
		vxnetTypeValidValues := []string{"0", "1"}
		vxnetTypeParameterValue := fmt.Sprint(*v.VxNetType)

		vxnetTypeIsValid := false
		for _, value := range vxnetTypeValidValues {
			if value == vxnetTypeParameterValue {
				vxnetTypeIsValid = true
			}
		}

		if !vxnetTypeIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "VxNetType",
				ParameterValue: vxnetTypeParameterValue,
				AllowedValues:  vxnetTypeValidValues,
			}
		}
	}

	return nil
}

type CreateVxNetsOutput struct {
	Message *string   `json:"message" name:"message"`
	Action  *string   `json:"action" name:"action" location:"elements"`
	RetCode *int      `json:"ret_code" name:"ret_code" location:"elements"`
	VxNets  []*string `json:"vxnets" name:"vxnets" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/vxnet/delete_vxnets.html
func (s *VxNetService) DeleteVxNets(i *DeleteVxNetsInput) (*DeleteVxNetsOutput, error) {
	if i == nil {
		i = &DeleteVxNetsInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DeleteVxnets",
		RequestMethod: "GET",
	}

	x := &DeleteVxNetsOutput{}
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

type DeleteVxNetsInput struct {
	VxNets []*string `json:"vxnets" name:"vxnets" location:"params"` // Required
}

func (v *DeleteVxNetsInput) Validate() error {

	if len(v.VxNets) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "VxNets",
			ParentName:    "DeleteVxNetsInput",
		}
	}

	return nil
}

type DeleteVxNetsOutput struct {
	Message *string   `json:"message" name:"message"`
	Action  *string   `json:"action" name:"action" location:"elements"`
	RetCode *int      `json:"ret_code" name:"ret_code" location:"elements"`
	VxNets  []*string `json:"vxnets" name:"vxnets" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/vxnet/describe_vxnet_instances.html
func (s *VxNetService) DescribeVxNetInstances(i *DescribeVxNetInstancesInput) (*DescribeVxNetInstancesOutput, error) {
	if i == nil {
		i = &DescribeVxNetInstancesInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DescribeVxnetInstances",
		RequestMethod: "GET",
	}

	x := &DescribeVxNetInstancesOutput{}
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

type DescribeVxNetInstancesInput struct {
	Image        *string   `json:"image" name:"image" location:"params"`
	InstanceType *string   `json:"instance_type" name:"instance_type" location:"params"`
	Instances    []*string `json:"instances" name:"instances" location:"params"`
	Limit        *int      `json:"limit" name:"limit" default:"20" location:"params"`
	Offset       *int      `json:"offset" name:"offset" default:"0" location:"params"`
	Status       *string   `json:"status" name:"status" location:"params"`
	VxNet        *string   `json:"vxnet" name:"vxnet" location:"params"` // Required
}

func (v *DescribeVxNetInstancesInput) Validate() error {

	if v.VxNet == nil {
		return errors.ParameterRequiredError{
			ParameterName: "VxNet",
			ParentName:    "DescribeVxNetInstancesInput",
		}
	}

	return nil
}

type DescribeVxNetInstancesOutput struct {
	Message     *string     `json:"message" name:"message"`
	Action      *string     `json:"action" name:"action" location:"elements"`
	InstanceSet []*Instance `json:"instance_set" name:"instance_set" location:"elements"`
	RetCode     *int        `json:"ret_code" name:"ret_code" location:"elements"`
	TotalCount  *int        `json:"total_count" name:"total_count" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/vxnet/describe_vxnets.html
func (s *VxNetService) DescribeVxNets(i *DescribeVxNetsInput) (*DescribeVxNetsOutput, error) {
	if i == nil {
		i = &DescribeVxNetsInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DescribeVxnets",
		RequestMethod: "GET",
	}

	x := &DescribeVxNetsOutput{}
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

type DescribeVxNetsInput struct {
	Limit      *int      `json:"limit" name:"limit" default:"20" location:"params"`
	Offset     *int      `json:"offset" name:"offset" default:"0" location:"params"`
	SearchWord *string   `json:"search_word" name:"search_word" location:"params"`
	Tags       []*string `json:"tags" name:"tags" location:"params"`
	// Verbose's available values: 0, 1
	Verbose *int `json:"verbose" name:"verbose" default:"0" location:"params"`
	// VxNetType's available values: 0, 1
	VxNetType *int      `json:"vxnet_type" name:"vxnet_type" location:"params"`
	VxNets    []*string `json:"vxnets" name:"vxnets" location:"params"`
}

func (v *DescribeVxNetsInput) Validate() error {

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

	if v.VxNetType != nil {
		vxnetTypeValidValues := []string{"0", "1"}
		vxnetTypeParameterValue := fmt.Sprint(*v.VxNetType)

		vxnetTypeIsValid := false
		for _, value := range vxnetTypeValidValues {
			if value == vxnetTypeParameterValue {
				vxnetTypeIsValid = true
			}
		}

		if !vxnetTypeIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "VxNetType",
				ParameterValue: vxnetTypeParameterValue,
				AllowedValues:  vxnetTypeValidValues,
			}
		}
	}

	return nil
}

type DescribeVxNetsOutput struct {
	Message    *string  `json:"message" name:"message"`
	Action     *string  `json:"action" name:"action" location:"elements"`
	RetCode    *int     `json:"ret_code" name:"ret_code" location:"elements"`
	TotalCount *int     `json:"total_count" name:"total_count" location:"elements"`
	VxNetSet   []*VxNet `json:"vxnet_set" name:"vxnet_set" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/vxnet/join_vxnet.html
func (s *VxNetService) JoinVxNet(i *JoinVxNetInput) (*JoinVxNetOutput, error) {
	if i == nil {
		i = &JoinVxNetInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "JoinVxnet",
		RequestMethod: "GET",
	}

	x := &JoinVxNetOutput{}
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

type JoinVxNetInput struct {
	Instances []*string `json:"instances" name:"instances" location:"params"` // Required
	VxNet     *string   `json:"vxnet" name:"vxnet" location:"params"`         // Required
}

func (v *JoinVxNetInput) Validate() error {

	if len(v.Instances) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "Instances",
			ParentName:    "JoinVxNetInput",
		}
	}

	if v.VxNet == nil {
		return errors.ParameterRequiredError{
			ParameterName: "VxNet",
			ParentName:    "JoinVxNetInput",
		}
	}

	return nil
}

type JoinVxNetOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	JobID   *string `json:"job_id" name:"job_id" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/vxnet/leave_vxnet.html
func (s *VxNetService) LeaveVxNet(i *LeaveVxNetInput) (*LeaveVxNetOutput, error) {
	if i == nil {
		i = &LeaveVxNetInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "LeaveVxnet",
		RequestMethod: "GET",
	}

	x := &LeaveVxNetOutput{}
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

type LeaveVxNetInput struct {
	Instances []*string `json:"instances" name:"instances" location:"params"` // Required
	VxNet     *string   `json:"vxnet" name:"vxnet" location:"params"`         // Required
}

func (v *LeaveVxNetInput) Validate() error {

	if len(v.Instances) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "Instances",
			ParentName:    "LeaveVxNetInput",
		}
	}

	if v.VxNet == nil {
		return errors.ParameterRequiredError{
			ParameterName: "VxNet",
			ParentName:    "LeaveVxNetInput",
		}
	}

	return nil
}

type LeaveVxNetOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	JobID   *string `json:"job_id" name:"job_id" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/vxnet/modify_vxnet_attributes.html
func (s *VxNetService) ModifyVxNetAttributes(i *ModifyVxNetAttributesInput) (*ModifyVxNetAttributesOutput, error) {
	if i == nil {
		i = &ModifyVxNetAttributesInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "ModifyVxnetAttributes",
		RequestMethod: "GET",
	}

	x := &ModifyVxNetAttributesOutput{}
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

type ModifyVxNetAttributesInput struct {
	Description *string `json:"description" name:"description" location:"params"`
	VxNet       *string `json:"vxnet" name:"vxnet" location:"params"` // Required
	VxNetName   *string `json:"vxnet_name" name:"vxnet_name" location:"params"`
}

func (v *ModifyVxNetAttributesInput) Validate() error {

	if v.VxNet == nil {
		return errors.ParameterRequiredError{
			ParameterName: "VxNet",
			ParentName:    "ModifyVxNetAttributesInput",
		}
	}

	return nil
}

type ModifyVxNetAttributesOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}
