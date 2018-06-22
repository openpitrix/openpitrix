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

type InstanceService struct {
	Config     *config.Config
	Properties *InstanceServiceProperties
}

type InstanceServiceProperties struct {
	// QingCloud Zone ID
	Zone *string `json:"zone" name:"zone"` // Required
}

func (s *QingCloudService) Instance(zone string) (*InstanceService, error) {
	properties := &InstanceServiceProperties{
		Zone: &zone,
	}

	return &InstanceService{Config: s.Config, Properties: properties}, nil
}

// Documentation URL: https://docs.qingcloud.com/api/instance/cease_instances.html
func (s *InstanceService) CeaseInstances(i *CeaseInstancesInput) (*CeaseInstancesOutput, error) {
	if i == nil {
		i = &CeaseInstancesInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "CeaseInstances",
		RequestMethod: "GET",
	}

	x := &CeaseInstancesOutput{}
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

type CeaseInstancesInput struct {
	Instances []*string `json:"instances" name:"instances" location:"params"` // Required
}

func (v *CeaseInstancesInput) Validate() error {

	if len(v.Instances) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "Instances",
			ParentName:    "CeaseInstancesInput",
		}
	}

	return nil
}

type CeaseInstancesOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	JobID   *string `json:"job_id" name:"job_id" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/instance/describe_instance_types.html
func (s *InstanceService) DescribeInstanceTypes(i *DescribeInstanceTypesInput) (*DescribeInstanceTypesOutput, error) {
	if i == nil {
		i = &DescribeInstanceTypesInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DescribeInstanceTypes",
		RequestMethod: "GET",
	}

	x := &DescribeInstanceTypesOutput{}
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

type DescribeInstanceTypesInput struct {
	InstanceTypes []*string `json:"instance_types" name:"instance_types" location:"params"`
}

func (v *DescribeInstanceTypesInput) Validate() error {

	return nil
}

type DescribeInstanceTypesOutput struct {
	Message         *string         `json:"message" name:"message"`
	Action          *string         `json:"action" name:"action" location:"elements"`
	InstanceTypeSet []*InstanceType `json:"instance_type_set" name:"instance_type_set" location:"elements"`
	RetCode         *int            `json:"ret_code" name:"ret_code" location:"elements"`
	TotalCount      *int            `json:"total_count" name:"total_count" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/instance/describe_instances.html
func (s *InstanceService) DescribeInstances(i *DescribeInstancesInput) (*DescribeInstancesOutput, error) {
	if i == nil {
		i = &DescribeInstancesInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DescribeInstances",
		RequestMethod: "GET",
	}

	x := &DescribeInstancesOutput{}
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

type DescribeInstancesInput struct {
	ImageID []*string `json:"image_id" name:"image_id" location:"params"`
	// InstanceClass's available values: 0, 1
	InstanceClass *int      `json:"instance_class" name:"instance_class" location:"params"`
	InstanceType  []*string `json:"instance_type" name:"instance_type" location:"params"`
	Instances     []*string `json:"instances" name:"instances" location:"params"`
	IsClusterNode *int      `json:"is_cluster_node" name:"is_cluster_node" default:"0" location:"params"`
	Limit         *int      `json:"limit" name:"limit" default:"20" location:"params"`
	Offset        *int      `json:"offset" name:"offset" default:"0" location:"params"`
	SearchWord    *string   `json:"search_word" name:"search_word" location:"params"`
	Status        []*string `json:"status" name:"status" location:"params"`
	Tags          []*string `json:"tags" name:"tags" location:"params"`
	// Verbose's available values: 0, 1
	Verbose *int `json:"verbose" name:"verbose" location:"params"`
}

func (v *DescribeInstancesInput) Validate() error {

	if v.InstanceClass != nil {
		instanceClassValidValues := []string{"0", "1"}
		instanceClassParameterValue := fmt.Sprint(*v.InstanceClass)

		instanceClassIsValid := false
		for _, value := range instanceClassValidValues {
			if value == instanceClassParameterValue {
				instanceClassIsValid = true
			}
		}

		if !instanceClassIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "InstanceClass",
				ParameterValue: instanceClassParameterValue,
				AllowedValues:  instanceClassValidValues,
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

type DescribeInstancesOutput struct {
	Message     *string     `json:"message" name:"message"`
	Action      *string     `json:"action" name:"action" location:"elements"`
	InstanceSet []*Instance `json:"instance_set" name:"instance_set" location:"elements"`
	RetCode     *int        `json:"ret_code" name:"ret_code" location:"elements"`
	TotalCount  *int        `json:"total_count" name:"total_count" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/instance/modify_instance_attributes.html
func (s *InstanceService) ModifyInstanceAttributes(i *ModifyInstanceAttributesInput) (*ModifyInstanceAttributesOutput, error) {
	if i == nil {
		i = &ModifyInstanceAttributesInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "ModifyInstanceAttributes",
		RequestMethod: "GET",
	}

	x := &ModifyInstanceAttributesOutput{}
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

type ModifyInstanceAttributesInput struct {
	Description  *string `json:"description" name:"description" location:"params"`
	Instance     *string `json:"instance" name:"instance" location:"params"` // Required
	InstanceName *string `json:"instance_name" name:"instance_name" location:"params"`
}

func (v *ModifyInstanceAttributesInput) Validate() error {

	if v.Instance == nil {
		return errors.ParameterRequiredError{
			ParameterName: "Instance",
			ParentName:    "ModifyInstanceAttributesInput",
		}
	}

	return nil
}

type ModifyInstanceAttributesOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/instance/reset_instances.html
func (s *InstanceService) ResetInstances(i *ResetInstancesInput) (*ResetInstancesOutput, error) {
	if i == nil {
		i = &ResetInstancesInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "ResetInstances",
		RequestMethod: "GET",
	}

	x := &ResetInstancesOutput{}
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

type ResetInstancesInput struct {
	Instances    []*string `json:"instances" name:"instances" location:"params"` // Required
	LoginKeyPair *string   `json:"login_keypair" name:"login_keypair" location:"params"`
	// LoginMode's available values: keypair, passwd
	LoginMode   *string `json:"login_mode" name:"login_mode" location:"params"` // Required
	LoginPasswd *string `json:"login_passwd" name:"login_passwd" location:"params"`
	// NeedNewSID's available values: 0, 1
	NeedNewSID *int `json:"need_newsid" name:"need_newsid" default:"0" location:"params"`
}

func (v *ResetInstancesInput) Validate() error {

	if len(v.Instances) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "Instances",
			ParentName:    "ResetInstancesInput",
		}
	}

	if v.LoginMode == nil {
		return errors.ParameterRequiredError{
			ParameterName: "LoginMode",
			ParentName:    "ResetInstancesInput",
		}
	}

	if v.LoginMode != nil {
		loginModeValidValues := []string{"keypair", "passwd"}
		loginModeParameterValue := fmt.Sprint(*v.LoginMode)

		loginModeIsValid := false
		for _, value := range loginModeValidValues {
			if value == loginModeParameterValue {
				loginModeIsValid = true
			}
		}

		if !loginModeIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "LoginMode",
				ParameterValue: loginModeParameterValue,
				AllowedValues:  loginModeValidValues,
			}
		}
	}

	if v.NeedNewSID != nil {
		needNewSIDValidValues := []string{"0", "1"}
		needNewSIDParameterValue := fmt.Sprint(*v.NeedNewSID)

		needNewSIDIsValid := false
		for _, value := range needNewSIDValidValues {
			if value == needNewSIDParameterValue {
				needNewSIDIsValid = true
			}
		}

		if !needNewSIDIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "NeedNewSID",
				ParameterValue: needNewSIDParameterValue,
				AllowedValues:  needNewSIDValidValues,
			}
		}
	}

	return nil
}

type ResetInstancesOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	JobID   *string `json:"job_id" name:"job_id" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/instance/resize_instances.html
func (s *InstanceService) ResizeInstances(i *ResizeInstancesInput) (*ResizeInstancesOutput, error) {
	if i == nil {
		i = &ResizeInstancesInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "ResizeInstances",
		RequestMethod: "GET",
	}

	x := &ResizeInstancesOutput{}
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

type ResizeInstancesInput struct {

	// CPU's available values: 1, 2, 4, 8, 16
	CPU          *int      `json:"cpu" name:"cpu" location:"params"`
	InstanceType *string   `json:"instance_type" name:"instance_type" location:"params"`
	Instances    []*string `json:"instances" name:"instances" location:"params"` // Required
	// Memory's available values: 1024, 2048, 4096, 6144, 8192, 12288, 16384, 24576, 32768
	Memory *int `json:"memory" name:"memory" location:"params"`
}

func (v *ResizeInstancesInput) Validate() error {

	if v.CPU != nil {
		cpuValidValues := []string{"1", "2", "4", "8", "16"}
		cpuParameterValue := fmt.Sprint(*v.CPU)

		cpuIsValid := false
		for _, value := range cpuValidValues {
			if value == cpuParameterValue {
				cpuIsValid = true
			}
		}

		if !cpuIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "CPU",
				ParameterValue: cpuParameterValue,
				AllowedValues:  cpuValidValues,
			}
		}
	}

	if len(v.Instances) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "Instances",
			ParentName:    "ResizeInstancesInput",
		}
	}

	if v.Memory != nil {
		memoryValidValues := []string{"1024", "2048", "4096", "6144", "8192", "12288", "16384", "24576", "32768"}
		memoryParameterValue := fmt.Sprint(*v.Memory)

		memoryIsValid := false
		for _, value := range memoryValidValues {
			if value == memoryParameterValue {
				memoryIsValid = true
			}
		}

		if !memoryIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "Memory",
				ParameterValue: memoryParameterValue,
				AllowedValues:  memoryValidValues,
			}
		}
	}

	return nil
}

type ResizeInstancesOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	JobID   *string `json:"job_id" name:"job_id" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/instance/restart_instances.html
func (s *InstanceService) RestartInstances(i *RestartInstancesInput) (*RestartInstancesOutput, error) {
	if i == nil {
		i = &RestartInstancesInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "RestartInstances",
		RequestMethod: "GET",
	}

	x := &RestartInstancesOutput{}
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

type RestartInstancesInput struct {
	Instances []*string `json:"instances" name:"instances" location:"params"` // Required
}

func (v *RestartInstancesInput) Validate() error {

	if len(v.Instances) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "Instances",
			ParentName:    "RestartInstancesInput",
		}
	}

	return nil
}

type RestartInstancesOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	JobID   *string `json:"job_id" name:"job_id" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/instance/run_instances.html
func (s *InstanceService) RunInstances(i *RunInstancesInput) (*RunInstancesOutput, error) {
	if i == nil {
		i = &RunInstancesInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "RunInstances",
		RequestMethod: "GET",
	}

	x := &RunInstancesOutput{}
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

type RunInstancesInput struct {
	BillingID *string `json:"billing_id" name:"billing_id" location:"params"`
	Count     *int    `json:"count" name:"count" default:"1" location:"params"`
	// CPU's available values: 1, 2, 4, 8, 16
	CPU *int `json:"cpu" name:"cpu" default:"1" location:"params"`
	// CPUMax's available values: 1, 2, 4, 8, 16
	CPUMax   *int    `json:"cpu_max" name:"cpu_max" location:"params"`
	Gpu      *int    `json:"gpu" name:"gpu" default:"0" location:"params"`
	Hostname *string `json:"hostname" name:"hostname" location:"params"`
	ImageID  *string `json:"image_id" name:"image_id" location:"params"` // Required
	// InstanceClass's available values: 0, 1
	InstanceClass *int    `json:"instance_class" name:"instance_class" location:"params"`
	InstanceName  *string `json:"instance_name" name:"instance_name" location:"params"`
	InstanceType  *string `json:"instance_type" name:"instance_type" location:"params"`
	LoginKeyPair  *string `json:"login_keypair" name:"login_keypair" location:"params"`
	// LoginMode's available values: keypair, passwd
	LoginMode   *string `json:"login_mode" name:"login_mode" location:"params"` // Required
	LoginPasswd *string `json:"login_passwd" name:"login_passwd" location:"params"`
	// MemMax's available values: 1024, 2048, 4096, 6144, 8192, 12288, 16384, 24576, 32768
	MemMax *int `json:"mem_max" name:"mem_max" location:"params"`
	// Memory's available values: 1024, 2048, 4096, 6144, 8192, 12288, 16384, 24576, 32768
	Memory *int `json:"memory" name:"memory" default:"1024" location:"params"`
	// NeedNewSID's available values: 0, 1
	NeedNewSID *int `json:"need_newsid" name:"need_newsid" default:"0" location:"params"`
	// NeedUserdata's available values: 0, 1
	NeedUserdata  *int    `json:"need_userdata" name:"need_userdata" default:"0" location:"params"`
	SecurityGroup *string `json:"security_group" name:"security_group" location:"params"`
	UIType        *string `json:"ui_type" name:"ui_type" location:"params"`
	UserdataFile  *string `json:"userdata_file" name:"userdata_file" default:"/etc/rc.local" location:"params"`
	UserdataPath  *string `json:"userdata_path" name:"userdata_path" default:"/etc/qingcloud/userdata" location:"params"`
	// UserdataType's available values: plain, exec, tar
	UserdataType  *string   `json:"userdata_type" name:"userdata_type" location:"params"`
	UserdataValue *string   `json:"userdata_value" name:"userdata_value" location:"params"`
	Volumes       []*string `json:"volumes" name:"volumes" location:"params"`
	VxNets        []*string `json:"vxnets" name:"vxnets" location:"params"`
}

func (v *RunInstancesInput) Validate() error {

	if v.CPU != nil {
		cpuValidValues := []string{"1", "2", "4", "8", "16"}
		cpuParameterValue := fmt.Sprint(*v.CPU)

		cpuIsValid := false
		for _, value := range cpuValidValues {
			if value == cpuParameterValue {
				cpuIsValid = true
			}
		}

		if !cpuIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "CPU",
				ParameterValue: cpuParameterValue,
				AllowedValues:  cpuValidValues,
			}
		}
	}

	if v.CPUMax != nil {
		cpuMaxValidValues := []string{"1", "2", "4", "8", "16"}
		cpuMaxParameterValue := fmt.Sprint(*v.CPUMax)

		cpuMaxIsValid := false
		for _, value := range cpuMaxValidValues {
			if value == cpuMaxParameterValue {
				cpuMaxIsValid = true
			}
		}

		if !cpuMaxIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "CPUMax",
				ParameterValue: cpuMaxParameterValue,
				AllowedValues:  cpuMaxValidValues,
			}
		}
	}

	if v.ImageID == nil {
		return errors.ParameterRequiredError{
			ParameterName: "ImageID",
			ParentName:    "RunInstancesInput",
		}
	}

	if v.InstanceClass != nil {
		instanceClassValidValues := []string{"0", "1"}
		instanceClassParameterValue := fmt.Sprint(*v.InstanceClass)

		instanceClassIsValid := false
		for _, value := range instanceClassValidValues {
			if value == instanceClassParameterValue {
				instanceClassIsValid = true
			}
		}

		if !instanceClassIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "InstanceClass",
				ParameterValue: instanceClassParameterValue,
				AllowedValues:  instanceClassValidValues,
			}
		}
	}

	if v.LoginMode == nil {
		return errors.ParameterRequiredError{
			ParameterName: "LoginMode",
			ParentName:    "RunInstancesInput",
		}
	}

	if v.LoginMode != nil {
		loginModeValidValues := []string{"keypair", "passwd"}
		loginModeParameterValue := fmt.Sprint(*v.LoginMode)

		loginModeIsValid := false
		for _, value := range loginModeValidValues {
			if value == loginModeParameterValue {
				loginModeIsValid = true
			}
		}

		if !loginModeIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "LoginMode",
				ParameterValue: loginModeParameterValue,
				AllowedValues:  loginModeValidValues,
			}
		}
	}

	if v.MemMax != nil {
		memMaxValidValues := []string{"1024", "2048", "4096", "6144", "8192", "12288", "16384", "24576", "32768"}
		memMaxParameterValue := fmt.Sprint(*v.MemMax)

		memMaxIsValid := false
		for _, value := range memMaxValidValues {
			if value == memMaxParameterValue {
				memMaxIsValid = true
			}
		}

		if !memMaxIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "MemMax",
				ParameterValue: memMaxParameterValue,
				AllowedValues:  memMaxValidValues,
			}
		}
	}

	if v.Memory != nil {
		memoryValidValues := []string{"1024", "2048", "4096", "6144", "8192", "12288", "16384", "24576", "32768"}
		memoryParameterValue := fmt.Sprint(*v.Memory)

		memoryIsValid := false
		for _, value := range memoryValidValues {
			if value == memoryParameterValue {
				memoryIsValid = true
			}
		}

		if !memoryIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "Memory",
				ParameterValue: memoryParameterValue,
				AllowedValues:  memoryValidValues,
			}
		}
	}

	if v.NeedNewSID != nil {
		needNewSIDValidValues := []string{"0", "1"}
		needNewSIDParameterValue := fmt.Sprint(*v.NeedNewSID)

		needNewSIDIsValid := false
		for _, value := range needNewSIDValidValues {
			if value == needNewSIDParameterValue {
				needNewSIDIsValid = true
			}
		}

		if !needNewSIDIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "NeedNewSID",
				ParameterValue: needNewSIDParameterValue,
				AllowedValues:  needNewSIDValidValues,
			}
		}
	}

	if v.NeedUserdata != nil {
		needUserdataValidValues := []string{"0", "1"}
		needUserdataParameterValue := fmt.Sprint(*v.NeedUserdata)

		needUserdataIsValid := false
		for _, value := range needUserdataValidValues {
			if value == needUserdataParameterValue {
				needUserdataIsValid = true
			}
		}

		if !needUserdataIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "NeedUserdata",
				ParameterValue: needUserdataParameterValue,
				AllowedValues:  needUserdataValidValues,
			}
		}
	}

	if v.UserdataType != nil {
		userdataTypeValidValues := []string{"plain", "exec", "tar"}
		userdataTypeParameterValue := fmt.Sprint(*v.UserdataType)

		userdataTypeIsValid := false
		for _, value := range userdataTypeValidValues {
			if value == userdataTypeParameterValue {
				userdataTypeIsValid = true
			}
		}

		if !userdataTypeIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "UserdataType",
				ParameterValue: userdataTypeParameterValue,
				AllowedValues:  userdataTypeValidValues,
			}
		}
	}

	return nil
}

type RunInstancesOutput struct {
	Message   *string   `json:"message" name:"message"`
	Action    *string   `json:"action" name:"action" location:"elements"`
	Instances []*string `json:"instances" name:"instances" location:"elements"`
	JobID     *string   `json:"job_id" name:"job_id" location:"elements"`
	RetCode   *int      `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/instance/start_instances.html
func (s *InstanceService) StartInstances(i *StartInstancesInput) (*StartInstancesOutput, error) {
	if i == nil {
		i = &StartInstancesInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "StartInstances",
		RequestMethod: "GET",
	}

	x := &StartInstancesOutput{}
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

type StartInstancesInput struct {
	Instances []*string `json:"instances" name:"instances" location:"params"` // Required
}

func (v *StartInstancesInput) Validate() error {

	if len(v.Instances) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "Instances",
			ParentName:    "StartInstancesInput",
		}
	}

	return nil
}

type StartInstancesOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	JobID   *string `json:"job_id" name:"job_id" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/instance/stop_instances.html
func (s *InstanceService) StopInstances(i *StopInstancesInput) (*StopInstancesOutput, error) {
	if i == nil {
		i = &StopInstancesInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "StopInstances",
		RequestMethod: "GET",
	}

	x := &StopInstancesOutput{}
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

type StopInstancesInput struct {

	// Force's available values: 0, 1
	Force     *int      `json:"force" name:"force" default:"0" location:"params"`
	Instances []*string `json:"instances" name:"instances" location:"params"` // Required
}

func (v *StopInstancesInput) Validate() error {

	if v.Force != nil {
		forceValidValues := []string{"0", "1"}
		forceParameterValue := fmt.Sprint(*v.Force)

		forceIsValid := false
		for _, value := range forceValidValues {
			if value == forceParameterValue {
				forceIsValid = true
			}
		}

		if !forceIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "Force",
				ParameterValue: forceParameterValue,
				AllowedValues:  forceValidValues,
			}
		}
	}

	if len(v.Instances) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "Instances",
			ParentName:    "StopInstancesInput",
		}
	}

	return nil
}

type StopInstancesOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	JobID   *string `json:"job_id" name:"job_id" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/instance/terminate_instances.html
func (s *InstanceService) TerminateInstances(i *TerminateInstancesInput) (*TerminateInstancesOutput, error) {
	if i == nil {
		i = &TerminateInstancesInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "TerminateInstances",
		RequestMethod: "GET",
	}

	x := &TerminateInstancesOutput{}
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

type TerminateInstancesInput struct {
	Instances []*string `json:"instances" name:"instances" location:"params"` // Required
}

func (v *TerminateInstancesInput) Validate() error {

	if len(v.Instances) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "Instances",
			ParentName:    "TerminateInstancesInput",
		}
	}

	return nil
}

type TerminateInstancesOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	JobID   *string `json:"job_id" name:"job_id" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}
