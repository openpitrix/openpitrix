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

type SharedStorageService struct {
	Config     *config.Config
	Properties *SharedStorageServiceProperties
}

type SharedStorageServiceProperties struct {
	// QingCloud Zone ID
	Zone *string `json:"zone" name:"zone"` // Required
}

func (s *QingCloudService) SharedStorage(zone string) (*SharedStorageService, error) {
	properties := &SharedStorageServiceProperties{
		Zone: &zone,
	}

	return &SharedStorageService{Config: s.Config, Properties: properties}, nil
}

// Documentation URL: https://docs.qingcloud.com/api/vsan/attach_to_s2_shared_target.html
func (s *SharedStorageService) AttachToS2SharedTarget(i *AttachToS2SharedTargetInput) (*AttachToS2SharedTargetOutput, error) {
	if i == nil {
		i = &AttachToS2SharedTargetInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "AttachToS2SharedTarget",
		RequestMethod: "GET",
	}

	x := &AttachToS2SharedTargetOutput{}
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

type AttachToS2SharedTargetInput struct {
	SharedTarget *string   `json:"shared_target" name:"shared_target" location:"params"` // Required
	Volumes      []*string `json:"volumes" name:"volumes" location:"params"`             // Required
}

func (v *AttachToS2SharedTargetInput) Validate() error {

	if v.SharedTarget == nil {
		return errors.ParameterRequiredError{
			ParameterName: "SharedTarget",
			ParentName:    "AttachToS2SharedTargetInput",
		}
	}

	if len(v.Volumes) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "Volumes",
			ParentName:    "AttachToS2SharedTargetInput",
		}
	}

	return nil
}

type AttachToS2SharedTargetOutput struct {
	Message      *string         `json:"message" name:"message"`
	Action       *string         `json:"action" name:"action" location:"elements"`
	RetCode      *int            `json:"ret_code" name:"ret_code" location:"elements"`
	SharedTarget *S2SharedTarget `json:"shared_target" name:"shared_target" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/vsan/change_s2_server_vxnet.html
func (s *SharedStorageService) ChangeS2ServerVxNet(i *ChangeS2ServerVxNetInput) (*ChangeS2ServerVxNetOutput, error) {
	if i == nil {
		i = &ChangeS2ServerVxNetInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "ChangeS2ServerVxnet",
		RequestMethod: "GET",
	}

	x := &ChangeS2ServerVxNetOutput{}
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

type ChangeS2ServerVxNetInput struct {
	PrivateIP *string `json:"private_ip" name:"private_ip" location:"params"`
	S2Server  *string `json:"s2_server" name:"s2_server" location:"params"` // Required
	VxNet     *string `json:"vxnet" name:"vxnet" location:"params"`         // Required
}

func (v *ChangeS2ServerVxNetInput) Validate() error {

	if v.S2Server == nil {
		return errors.ParameterRequiredError{
			ParameterName: "S2Server",
			ParentName:    "ChangeS2ServerVxNetInput",
		}
	}

	if v.VxNet == nil {
		return errors.ParameterRequiredError{
			ParameterName: "VxNet",
			ParentName:    "ChangeS2ServerVxNetInput",
		}
	}

	return nil
}

type ChangeS2ServerVxNetOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	JobID   *string `json:"job_id" name:"job_id" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/vsan/create_s2_server.html
func (s *SharedStorageService) CreateS2Server(i *CreateS2ServerInput) (*CreateS2ServerOutput, error) {
	if i == nil {
		i = &CreateS2ServerInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "CreateS2Server",
		RequestMethod: "GET",
	}

	x := &CreateS2ServerOutput{}
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

type CreateS2ServerInput struct {
	Description *string `json:"description" name:"description" location:"params"`
	PrivateIP   *string `json:"private_ip" name:"private_ip" location:"params"`
	// S2Class's available values: 0, 1
	S2Class      *int    `json:"s2_class" name:"s2_class" location:"params"`
	S2ServerName *string `json:"s2_server_name" name:"s2_server_name" location:"params"`
	ServiceType  *string `json:"service_type" name:"service_type" location:"params"`
	VxNet        *string `json:"vxnet" name:"vxnet" location:"params"`
}

func (v *CreateS2ServerInput) Validate() error {

	if v.S2Class != nil {
		s2ClassValidValues := []string{"0", "1"}
		s2ClassParameterValue := fmt.Sprint(*v.S2Class)

		s2ClassIsValid := false
		for _, value := range s2ClassValidValues {
			if value == s2ClassParameterValue {
				s2ClassIsValid = true
			}
		}

		if !s2ClassIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "S2Class",
				ParameterValue: s2ClassParameterValue,
				AllowedValues:  s2ClassValidValues,
			}
		}
	}

	return nil
}

type CreateS2ServerOutput struct {
	Message  *string `json:"message" name:"message"`
	Action   *string `json:"action" name:"action" location:"elements"`
	JobID    *string `json:"job_id" name:"job_id" location:"elements"`
	RetCode  *int    `json:"ret_code" name:"ret_code" location:"elements"`
	S2Server *string `json:"s2_server" name:"s2_server" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/vsan/create_s2_shared_target.html
func (s *SharedStorageService) CreateS2SharedTarget(i *CreateS2SharedTargetInput) (*CreateS2SharedTargetOutput, error) {
	if i == nil {
		i = &CreateS2SharedTargetInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "CreateS2SharedTarget",
		RequestMethod: "GET",
	}

	x := &CreateS2SharedTargetOutput{}
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

type CreateS2SharedTargetInput struct {
	Description    *string   `json:"description" name:"description" location:"params"`
	ExportName     *string   `json:"export_name" name:"export_name" location:"params"` // Required
	ExportNameNfs  *string   `json:"export_name_nfs" name:"export_name_nfs" location:"params"`
	InitiatorNames []*string `json:"initiator_names" name:"initiator_names" location:"params"`
	S2Group        *string   `json:"s2_group" name:"s2_group" location:"params"`
	S2ServerID     *string   `json:"s2_server_id" name:"s2_server_id" location:"params"` // Required
	// TargetType's available values: ISCSI, NFS
	TargetType *string   `json:"target_type" name:"target_type" location:"params"` // Required
	Volumes    []*string `json:"volumes" name:"volumes" location:"params"`
}

func (v *CreateS2SharedTargetInput) Validate() error {

	if v.ExportName == nil {
		return errors.ParameterRequiredError{
			ParameterName: "ExportName",
			ParentName:    "CreateS2SharedTargetInput",
		}
	}

	if v.S2ServerID == nil {
		return errors.ParameterRequiredError{
			ParameterName: "S2ServerID",
			ParentName:    "CreateS2SharedTargetInput",
		}
	}

	if v.TargetType == nil {
		return errors.ParameterRequiredError{
			ParameterName: "TargetType",
			ParentName:    "CreateS2SharedTargetInput",
		}
	}

	if v.TargetType != nil {
		targetTypeValidValues := []string{"ISCSI", "NFS"}
		targetTypeParameterValue := fmt.Sprint(*v.TargetType)

		targetTypeIsValid := false
		for _, value := range targetTypeValidValues {
			if value == targetTypeParameterValue {
				targetTypeIsValid = true
			}
		}

		if !targetTypeIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "TargetType",
				ParameterValue: targetTypeParameterValue,
				AllowedValues:  targetTypeValidValues,
			}
		}
	}

	return nil
}

type CreateS2SharedTargetOutput struct {
	Message        *string `json:"message" name:"message"`
	Action         *string `json:"action" name:"action" location:"elements"`
	RetCode        *int    `json:"ret_code" name:"ret_code" location:"elements"`
	S2SharedTarget *string `json:"s2_shared_target" name:"s2_shared_target" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/vsan/delete_s2_servers.html
func (s *SharedStorageService) DeleteS2Servers(i *DeleteS2ServersInput) (*DeleteS2ServersOutput, error) {
	if i == nil {
		i = &DeleteS2ServersInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DeleteS2Servers",
		RequestMethod: "GET",
	}

	x := &DeleteS2ServersOutput{}
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

type DeleteS2ServersInput struct {
	S2Servers []*string `json:"s2_servers" name:"s2_servers" location:"params"` // Required
}

func (v *DeleteS2ServersInput) Validate() error {

	if len(v.S2Servers) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "S2Servers",
			ParentName:    "DeleteS2ServersInput",
		}
	}

	return nil
}

type DeleteS2ServersOutput struct {
	Message   *string   `json:"message" name:"message"`
	Action    *string   `json:"action" name:"action" location:"elements"`
	JobID     *string   `json:"job_id" name:"job_id" location:"elements"`
	RetCode   *int      `json:"ret_code" name:"ret_code" location:"elements"`
	S2Servers []*string `json:"s2_servers" name:"s2_servers" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/vsan/delete_s2_shared_targets.html
func (s *SharedStorageService) DeleteS2SharedTargets(i *DeleteS2SharedTargetsInput) (*DeleteS2SharedTargetsOutput, error) {
	if i == nil {
		i = &DeleteS2SharedTargetsInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DeleteS2SharedTargets",
		RequestMethod: "GET",
	}

	x := &DeleteS2SharedTargetsOutput{}
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

type DeleteS2SharedTargetsInput struct {
	SharedTargets []*string `json:"shared_targets" name:"shared_targets" location:"params"` // Required
}

func (v *DeleteS2SharedTargetsInput) Validate() error {

	if len(v.SharedTargets) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "SharedTargets",
			ParentName:    "DeleteS2SharedTargetsInput",
		}
	}

	return nil
}

type DeleteS2SharedTargetsOutput struct {
	Message       *string   `json:"message" name:"message"`
	Action        *string   `json:"action" name:"action" location:"elements"`
	RetCode       *int      `json:"ret_code" name:"ret_code" location:"elements"`
	SharedTargets []*string `json:"shared_targets" name:"shared_targets" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/vsan/describle_s2_default_parameters.html
func (s *SharedStorageService) DescribeS2DefaultParameters(i *DescribeS2DefaultParametersInput) (*DescribeS2DefaultParametersOutput, error) {
	if i == nil {
		i = &DescribeS2DefaultParametersInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DescribeS2DefaultParameters",
		RequestMethod: "GET",
	}

	x := &DescribeS2DefaultParametersOutput{}
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

type DescribeS2DefaultParametersInput struct {
	Limit  *int `json:"limit" name:"limit" default:"20" location:"params"`
	Offset *int `json:"offset" name:"offset" default:"0" location:"params"`
	// ServiceType's available values: vsan
	ServiceType *string `json:"service_type" name:"service_type" location:"params"`
	// TargetType's available values: ISCSI
	TargetType *string `json:"target_type" name:"target_type" location:"params"`
}

func (v *DescribeS2DefaultParametersInput) Validate() error {

	if v.ServiceType != nil {
		serviceTypeValidValues := []string{"vsan"}
		serviceTypeParameterValue := fmt.Sprint(*v.ServiceType)

		serviceTypeIsValid := false
		for _, value := range serviceTypeValidValues {
			if value == serviceTypeParameterValue {
				serviceTypeIsValid = true
			}
		}

		if !serviceTypeIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "ServiceType",
				ParameterValue: serviceTypeParameterValue,
				AllowedValues:  serviceTypeValidValues,
			}
		}
	}

	if v.TargetType != nil {
		targetTypeValidValues := []string{"ISCSI"}
		targetTypeParameterValue := fmt.Sprint(*v.TargetType)

		targetTypeIsValid := false
		for _, value := range targetTypeValidValues {
			if value == targetTypeParameterValue {
				targetTypeIsValid = true
			}
		}

		if !targetTypeIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "TargetType",
				ParameterValue: targetTypeParameterValue,
				AllowedValues:  targetTypeValidValues,
			}
		}
	}

	return nil
}

type DescribeS2DefaultParametersOutput struct {
	Message                *string                `json:"message" name:"message"`
	Action                 *string                `json:"action" name:"action" location:"elements"`
	RetCode                *int                   `json:"ret_code" name:"ret_code" location:"elements"`
	S2DefaultParametersSet []*S2DefaultParameters `json:"s2_default_parameters_set" name:"s2_default_parameters_set" location:"elements"`
	TotalCount             *int                   `json:"total_count" name:"total_count" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/vsan/describe_s2_servers.html
func (s *SharedStorageService) DescribeS2Servers(i *DescribeS2ServersInput) (*DescribeS2ServersOutput, error) {
	if i == nil {
		i = &DescribeS2ServersInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DescribeS2Servers",
		RequestMethod: "GET",
	}

	x := &DescribeS2ServersOutput{}
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

type DescribeS2ServersInput struct {
	Limit      *int      `json:"limit" name:"limit" default:"20" location:"params"`
	Offset     *int      `json:"offset" name:"offset" default:"0" location:"params"`
	S2Servers  []*string `json:"s2_servers" name:"s2_servers" location:"params"`
	SearchWord *string   `json:"search_word" name:"search_word" location:"params"`
	Status     []*string `json:"status" name:"status" location:"params"`
	Tags       []*string `json:"tags" name:"tags" location:"params"`
	Verbose    *int      `json:"verbose" name:"verbose" location:"params"`
}

func (v *DescribeS2ServersInput) Validate() error {

	return nil
}

type DescribeS2ServersOutput struct {
	Message     *string     `json:"message" name:"message"`
	Action      *string     `json:"action" name:"action" location:"elements"`
	RetCode     *int        `json:"ret_code" name:"ret_code" location:"elements"`
	S2ServerSet []*S2Server `json:"s2_server_set" name:"s2_server_set" location:"elements"`
	TotalCount  *int        `json:"total_count" name:"total_count" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/vsan/describe_s2_shared_targets.html
func (s *SharedStorageService) DescribeS2SharedTargets(i *DescribeS2SharedTargetsInput) (*DescribeS2SharedTargetsOutput, error) {
	if i == nil {
		i = &DescribeS2SharedTargetsInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DescribeS2SharedTargets",
		RequestMethod: "GET",
	}

	x := &DescribeS2SharedTargetsOutput{}
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

type DescribeS2SharedTargetsInput struct {
	Limit         *int      `json:"limit" name:"limit" default:"20" location:"params"`
	Offset        *int      `json:"offset" name:"offset" default:"0" location:"params"`
	S2ServerID    *string   `json:"s2_server_id" name:"s2_server_id" location:"params"`
	SearchWord    *string   `json:"search_word" name:"search_word" location:"params"`
	SharedTargets []*string `json:"shared_targets" name:"shared_targets" location:"params"`
	Verbose       *int      `json:"verbose" name:"verbose" location:"params"`
}

func (v *DescribeS2SharedTargetsInput) Validate() error {

	return nil
}

type DescribeS2SharedTargetsOutput struct {
	Message         *string           `json:"message" name:"message"`
	Action          *string           `json:"action" name:"action" location:"elements"`
	RetCode         *int              `json:"ret_code" name:"ret_code" location:"elements"`
	SharedTargetSet []*S2SharedTarget `json:"shared_target_set" name:"shared_target_set" location:"elements"`
	TotalCount      *int              `json:"total_count" name:"total_count" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/vsan/detach_from_s2_shared_target.html
func (s *SharedStorageService) DetachFromS2SharedTarget(i *DetachFromS2SharedTargetInput) (*DetachFromS2SharedTargetOutput, error) {
	if i == nil {
		i = &DetachFromS2SharedTargetInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DetachFromS2SharedTarget",
		RequestMethod: "GET",
	}

	x := &DetachFromS2SharedTargetOutput{}
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

type DetachFromS2SharedTargetInput struct {
	SharedTarget *string   `json:"shared_target" name:"shared_target" location:"params"` // Required
	Volumes      []*string `json:"volumes" name:"volumes" location:"params"`             // Required
}

func (v *DetachFromS2SharedTargetInput) Validate() error {

	if v.SharedTarget == nil {
		return errors.ParameterRequiredError{
			ParameterName: "SharedTarget",
			ParentName:    "DetachFromS2SharedTargetInput",
		}
	}

	if len(v.Volumes) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "Volumes",
			ParentName:    "DetachFromS2SharedTargetInput",
		}
	}

	return nil
}

type DetachFromS2SharedTargetOutput struct {
	Message      *string         `json:"message" name:"message"`
	Action       *string         `json:"action" name:"action" location:"elements"`
	RetCode      *int            `json:"ret_code" name:"ret_code" location:"elements"`
	SharedTarget *S2SharedTarget `json:"shared_target" name:"shared_target" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/vsan/disable_s2_shared_targets.html
func (s *SharedStorageService) DisableS2SharedTargets(i *DisableS2SharedTargetsInput) (*DisableS2SharedTargetsOutput, error) {
	if i == nil {
		i = &DisableS2SharedTargetsInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DisableS2SharedTargets",
		RequestMethod: "GET",
	}

	x := &DisableS2SharedTargetsOutput{}
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

type DisableS2SharedTargetsInput struct {
	SharedTargets []*string `json:"shared_targets" name:"shared_targets" location:"params"` // Required
}

func (v *DisableS2SharedTargetsInput) Validate() error {

	if len(v.SharedTargets) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "SharedTargets",
			ParentName:    "DisableS2SharedTargetsInput",
		}
	}

	return nil
}

type DisableS2SharedTargetsOutput struct {
	Message       *string   `json:"message" name:"message"`
	Action        *string   `json:"action" name:"action" location:"elements"`
	RetCode       *int      `json:"ret_code" name:"ret_code" location:"elements"`
	SharedTargets []*string `json:"shared_targets" name:"shared_targets" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/vsan/enable_s2_shared_targets.html
func (s *SharedStorageService) EnableS2SharedTargets(i *EnableS2SharedTargetsInput) (*EnableS2SharedTargetsOutput, error) {
	if i == nil {
		i = &EnableS2SharedTargetsInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "EnableS2SharedTargets",
		RequestMethod: "GET",
	}

	x := &EnableS2SharedTargetsOutput{}
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

type EnableS2SharedTargetsInput struct {
	SharedTargets []*string `json:"shared_targets" name:"shared_targets" location:"params"` // Required
}

func (v *EnableS2SharedTargetsInput) Validate() error {

	if len(v.SharedTargets) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "SharedTargets",
			ParentName:    "EnableS2SharedTargetsInput",
		}
	}

	return nil
}

type EnableS2SharedTargetsOutput struct {
	Message       *string   `json:"message" name:"message"`
	Action        *string   `json:"action" name:"action" location:"elements"`
	RetCode       *int      `json:"ret_code" name:"ret_code" location:"elements"`
	SharedTargets []*string `json:"shared_targets" name:"shared_targets" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/vsan/modify_s2_server.html
func (s *SharedStorageService) ModifyS2Server(i *ModifyS2ServerInput) (*ModifyS2ServerOutput, error) {
	if i == nil {
		i = &ModifyS2ServerInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "ModifyS2Server",
		RequestMethod: "GET",
	}

	x := &ModifyS2ServerOutput{}
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

type ModifyS2ServerInput struct {
	Description  *string `json:"description" name:"description" location:"params"`
	S2Server     *string `json:"s2_server" name:"s2_server" location:"params"` // Required
	S2ServerName *string `json:"s2_server_name" name:"s2_server_name" location:"params"`
}

func (v *ModifyS2ServerInput) Validate() error {

	if v.S2Server == nil {
		return errors.ParameterRequiredError{
			ParameterName: "S2Server",
			ParentName:    "ModifyS2ServerInput",
		}
	}

	return nil
}

type ModifyS2ServerOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/vsan/modify_s2_shared_target.html
func (s *SharedStorageService) ModifyS2SharedTargets(i *ModifyS2SharedTargetsInput) (*ModifyS2SharedTargetsOutput, error) {
	if i == nil {
		i = &ModifyS2SharedTargetsInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "ModifyS2SharedTargets",
		RequestMethod: "GET",
	}

	x := &ModifyS2SharedTargetsOutput{}
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

type ModifyS2SharedTargetsInput struct {
	InitiatorNames []*string `json:"initiator_names" name:"initiator_names" location:"params"`
	Operation      *string   `json:"operation" name:"operation" location:"params"`           // Required
	Parameters     []*string `json:"parameters" name:"parameters" location:"params"`         // Required
	SharedTargets  []*string `json:"shared_targets" name:"shared_targets" location:"params"` // Required
}

func (v *ModifyS2SharedTargetsInput) Validate() error {

	if v.Operation == nil {
		return errors.ParameterRequiredError{
			ParameterName: "Operation",
			ParentName:    "ModifyS2SharedTargetsInput",
		}
	}

	if len(v.Parameters) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "Parameters",
			ParentName:    "ModifyS2SharedTargetsInput",
		}
	}

	if len(v.SharedTargets) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "SharedTargets",
			ParentName:    "ModifyS2SharedTargetsInput",
		}
	}

	return nil
}

type ModifyS2SharedTargetsOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/vsan/poweroff_s2_servers.html
func (s *SharedStorageService) PowerOffS2Servers(i *PowerOffS2ServersInput) (*PowerOffS2ServersOutput, error) {
	if i == nil {
		i = &PowerOffS2ServersInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "PowerOffS2Servers",
		RequestMethod: "GET",
	}

	x := &PowerOffS2ServersOutput{}
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

type PowerOffS2ServersInput struct {
	S2Servers *string `json:"s2_servers" name:"s2_servers" location:"params"` // Required
}

func (v *PowerOffS2ServersInput) Validate() error {

	if v.S2Servers == nil {
		return errors.ParameterRequiredError{
			ParameterName: "S2Servers",
			ParentName:    "PowerOffS2ServersInput",
		}
	}

	return nil
}

type PowerOffS2ServersOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	JobID   *string `json:"job_id" name:"job_id" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/vsan/poweron_s2_servers.html
func (s *SharedStorageService) PowerOnS2Servers(i *PowerOnS2ServersInput) (*PowerOnS2ServersOutput, error) {
	if i == nil {
		i = &PowerOnS2ServersInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "PowerOnS2Servers",
		RequestMethod: "GET",
	}

	x := &PowerOnS2ServersOutput{}
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

type PowerOnS2ServersInput struct {
	S2Servers []*string `json:"s2_servers" name:"s2_servers" location:"params"` // Required
}

func (v *PowerOnS2ServersInput) Validate() error {

	if len(v.S2Servers) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "S2Servers",
			ParentName:    "PowerOnS2ServersInput",
		}
	}

	return nil
}

type PowerOnS2ServersOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	JobID   *string `json:"job_id" name:"job_id" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/vsan/resize_s2_servers.html
func (s *SharedStorageService) ResizeS2Servers(i *ResizeS2ServersInput) (*ResizeS2ServersOutput, error) {
	if i == nil {
		i = &ResizeS2ServersInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "ResizeS2Servers",
		RequestMethod: "GET",
	}

	x := &ResizeS2ServersOutput{}
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

type ResizeS2ServersInput struct {
	S2Server     *string `json:"s2_server" name:"s2_server" location:"params"`           // Required
	S2ServerType *int    `json:"s2_server_type" name:"s2_server_type" location:"params"` // Required
}

func (v *ResizeS2ServersInput) Validate() error {

	if v.S2Server == nil {
		return errors.ParameterRequiredError{
			ParameterName: "S2Server",
			ParentName:    "ResizeS2ServersInput",
		}
	}

	if v.S2ServerType == nil {
		return errors.ParameterRequiredError{
			ParameterName: "S2ServerType",
			ParentName:    "ResizeS2ServersInput",
		}
	}

	return nil
}

type ResizeS2ServersOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	JobID   *string `json:"job_id" name:"job_id" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/vsan/update_s2_servers.html
func (s *SharedStorageService) UpdateS2Servers(i *UpdateS2ServersInput) (*UpdateS2ServersOutput, error) {
	if i == nil {
		i = &UpdateS2ServersInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "UpdateS2Servers",
		RequestMethod: "GET",
	}

	x := &UpdateS2ServersOutput{}
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

type UpdateS2ServersInput struct {
	S2Servers []*string `json:"s2_servers" name:"s2_servers" location:"params"` // Required
}

func (v *UpdateS2ServersInput) Validate() error {

	if len(v.S2Servers) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "S2Servers",
			ParentName:    "UpdateS2ServersInput",
		}
	}

	return nil
}

type UpdateS2ServersOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	JobID   *string `json:"job_id" name:"job_id" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}
