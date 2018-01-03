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

type NicService struct {
	Config     *config.Config
	Properties *NicServiceProperties
}

type NicServiceProperties struct {
	// QingCloud Zone ID
	Zone *string `json:"zone" name:"zone"` // Required
}

func (s *QingCloudService) Nic(zone string) (*NicService, error) {
	properties := &NicServiceProperties{
		Zone: &zone,
	}

	return &NicService{Config: s.Config, Properties: properties}, nil
}

// Documentation URL: https://docs.qingcloud.com/api/nic/attach_nics.html
func (s *NicService) AttachNics(i *AttachNicsInput) (*AttachNicsOutput, error) {
	if i == nil {
		i = &AttachNicsInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "AttachNics",
		RequestMethod: "GET",
	}

	x := &AttachNicsOutput{}
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

type AttachNicsInput struct {
	Instance *string   `json:"instance" name:"instance" location:"params"` // Required
	Nics     []*string `json:"nics" name:"nics" location:"params"`         // Required
}

func (v *AttachNicsInput) Validate() error {

	if v.Instance == nil {
		return errors.ParameterRequiredError{
			ParameterName: "Instance",
			ParentName:    "AttachNicsInput",
		}
	}

	if len(v.Nics) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "Nics",
			ParentName:    "AttachNicsInput",
		}
	}

	return nil
}

type AttachNicsOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	JobID   *string `json:"job_id" name:"job_id" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/nic/create_nics.html
func (s *NicService) CreateNics(i *CreateNicsInput) (*CreateNicsOutput, error) {
	if i == nil {
		i = &CreateNicsInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "CreateNics",
		RequestMethod: "GET",
	}

	x := &CreateNicsOutput{}
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

type CreateNicsInput struct {
	Count      *int      `json:"count" name:"count" default:"1" location:"params"`
	NICName    *string   `json:"nic_name" name:"nic_name" location:"params"`
	PrivateIPs []*string `json:"private_ips" name:"private_ips" location:"params"`
	VxNet      *string   `json:"vxnet" name:"vxnet" location:"params"` // Required
}

func (v *CreateNicsInput) Validate() error {

	if v.VxNet == nil {
		return errors.ParameterRequiredError{
			ParameterName: "VxNet",
			ParentName:    "CreateNicsInput",
		}
	}

	return nil
}

type CreateNicsOutput struct {
	Message *string  `json:"message" name:"message"`
	Action  *string  `json:"action" name:"action" location:"elements"`
	Nics    []*NICIP `json:"nics" name:"nics" location:"elements"`
	RetCode *int     `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/nic/delete_nics.html
func (s *NicService) DeleteNics(i *DeleteNicsInput) (*DeleteNicsOutput, error) {
	if i == nil {
		i = &DeleteNicsInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DeleteNics",
		RequestMethod: "GET",
	}

	x := &DeleteNicsOutput{}
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

type DeleteNicsInput struct {
	Nics []*string `json:"nics" name:"nics" location:"params"` // Required
}

func (v *DeleteNicsInput) Validate() error {

	if len(v.Nics) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "Nics",
			ParentName:    "DeleteNicsInput",
		}
	}

	return nil
}

type DeleteNicsOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/nic/describe_nics.html
func (s *NicService) DescribeNics(i *DescribeNicsInput) (*DescribeNicsOutput, error) {
	if i == nil {
		i = &DescribeNicsInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DescribeNics",
		RequestMethod: "GET",
	}

	x := &DescribeNicsOutput{}
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

type DescribeNicsInput struct {
	Instances []*string `json:"instances" name:"instances" location:"params"`
	Limit     *int      `json:"limit" name:"limit" default:"20" location:"params"`
	NICName   *string   `json:"nic_name" name:"nic_name" location:"params"`
	Nics      []*string `json:"nics" name:"nics" location:"params"`
	Offset    *int      `json:"offset" name:"offset" default:"0" location:"params"`
	// Status's available values: available, in-use
	Status *string `json:"status" name:"status" location:"params"`
	// VxNetType's available values: 0, 1
	VxNetType *int      `json:"vxnet_type" name:"vxnet_type" location:"params"`
	VxNets    []*string `json:"vxnets" name:"vxnets" location:"params"`
}

func (v *DescribeNicsInput) Validate() error {

	if v.Status != nil {
		statusValidValues := []string{"available", "in-use"}
		statusParameterValue := fmt.Sprint(*v.Status)

		statusIsValid := false
		for _, value := range statusValidValues {
			if value == statusParameterValue {
				statusIsValid = true
			}
		}

		if !statusIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "Status",
				ParameterValue: statusParameterValue,
				AllowedValues:  statusValidValues,
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

type DescribeNicsOutput struct {
	Message    *string `json:"message" name:"message"`
	Action     *string `json:"action" name:"action" location:"elements"`
	NICSet     []*NIC  `json:"nic_set" name:"nic_set" location:"elements"`
	RetCode    *int    `json:"ret_code" name:"ret_code" location:"elements"`
	TotalCount *int    `json:"total_count" name:"total_count" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/nic/detach_nics.html
func (s *NicService) DetachNics(i *DetachNicsInput) (*DetachNicsOutput, error) {
	if i == nil {
		i = &DetachNicsInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DetachNics",
		RequestMethod: "GET",
	}

	x := &DetachNicsOutput{}
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

type DetachNicsInput struct {
	Nics []*string `json:"nics" name:"nics" location:"params"` // Required
}

func (v *DetachNicsInput) Validate() error {

	if len(v.Nics) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "Nics",
			ParentName:    "DetachNicsInput",
		}
	}

	return nil
}

type DetachNicsOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	JobID   *string `json:"job_id" name:"job_id" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/nic/modify-nic-attributes.html
func (s *NicService) ModifyNicAttributes(i *ModifyNicAttributesInput) (*ModifyNicAttributesOutput, error) {
	if i == nil {
		i = &ModifyNicAttributesInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "ModifyNicAttributes",
		RequestMethod: "GET",
	}

	x := &ModifyNicAttributesOutput{}
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

type ModifyNicAttributesInput struct {
	NICID     *string `json:"nic_id" name:"nic_id" location:"params"` // Required
	NICName   *string `json:"nic_name" name:"nic_name" location:"params"`
	PrivateIP *string `json:"private_ip" name:"private_ip" location:"params"`
	VxNet     *string `json:"vxnet" name:"vxnet" location:"params"`
}

func (v *ModifyNicAttributesInput) Validate() error {

	if v.NICID == nil {
		return errors.ParameterRequiredError{
			ParameterName: "NICID",
			ParentName:    "ModifyNicAttributesInput",
		}
	}

	return nil
}

type ModifyNicAttributesOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}
