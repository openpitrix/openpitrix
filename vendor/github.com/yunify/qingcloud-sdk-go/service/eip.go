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

type EIPService struct {
	Config     *config.Config
	Properties *EIPServiceProperties
}

type EIPServiceProperties struct {
	// QingCloud Zone ID
	Zone *string `json:"zone" name:"zone"` // Required
}

func (s *QingCloudService) EIP(zone string) (*EIPService, error) {
	properties := &EIPServiceProperties{
		Zone: &zone,
	}

	return &EIPService{Config: s.Config, Properties: properties}, nil
}

// Documentation URL: https://docs.qingcloud.com/api/eip/allocate_eips.html
func (s *EIPService) AllocateEIPs(i *AllocateEIPsInput) (*AllocateEIPsOutput, error) {
	if i == nil {
		i = &AllocateEIPsInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "AllocateEips",
		RequestMethod: "GET",
	}

	x := &AllocateEIPsOutput{}
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

type AllocateEIPsInput struct {
	Bandwidth *int `json:"bandwidth" name:"bandwidth" location:"params"` // Required
	// BillingMode's available values: bandwidth, traffic
	BillingMode *string `json:"billing_mode" name:"billing_mode" default:"bandwidth" location:"params"`
	Count       *int    `json:"count" name:"count" default:"1" location:"params"`
	EIPName     *string `json:"eip_name" name:"eip_name" location:"params"`
	// NeedICP's available values: 0, 1
	NeedICP *int `json:"need_icp" name:"need_icp" default:"0" location:"params"`
}

func (v *AllocateEIPsInput) Validate() error {

	if v.Bandwidth == nil {
		return errors.ParameterRequiredError{
			ParameterName: "Bandwidth",
			ParentName:    "AllocateEIPsInput",
		}
	}

	if v.BillingMode != nil {
		billingModeValidValues := []string{"bandwidth", "traffic"}
		billingModeParameterValue := fmt.Sprint(*v.BillingMode)

		billingModeIsValid := false
		for _, value := range billingModeValidValues {
			if value == billingModeParameterValue {
				billingModeIsValid = true
			}
		}

		if !billingModeIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "BillingMode",
				ParameterValue: billingModeParameterValue,
				AllowedValues:  billingModeValidValues,
			}
		}
	}

	if v.NeedICP != nil {
		needICPValidValues := []string{"0", "1"}
		needICPParameterValue := fmt.Sprint(*v.NeedICP)

		needICPIsValid := false
		for _, value := range needICPValidValues {
			if value == needICPParameterValue {
				needICPIsValid = true
			}
		}

		if !needICPIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "NeedICP",
				ParameterValue: needICPParameterValue,
				AllowedValues:  needICPValidValues,
			}
		}
	}

	return nil
}

type AllocateEIPsOutput struct {
	Message *string   `json:"message" name:"message"`
	Action  *string   `json:"action" name:"action" location:"elements"`
	EIPs    []*string `json:"eips" name:"eips" location:"elements"`
	RetCode *int      `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/eip/associate_eip.html
func (s *EIPService) AssociateEIP(i *AssociateEIPInput) (*AssociateEIPOutput, error) {
	if i == nil {
		i = &AssociateEIPInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "AssociateEip",
		RequestMethod: "GET",
	}

	x := &AssociateEIPOutput{}
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

type AssociateEIPInput struct {
	EIP      *string `json:"eip" name:"eip" location:"params"`           // Required
	Instance *string `json:"instance" name:"instance" location:"params"` // Required
}

func (v *AssociateEIPInput) Validate() error {

	if v.EIP == nil {
		return errors.ParameterRequiredError{
			ParameterName: "EIP",
			ParentName:    "AssociateEIPInput",
		}
	}

	if v.Instance == nil {
		return errors.ParameterRequiredError{
			ParameterName: "Instance",
			ParentName:    "AssociateEIPInput",
		}
	}

	return nil
}

type AssociateEIPOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	JobID   *string `json:"job_id" name:"job_id" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/eip/dissociate_eips.html
func (s *EIPService) ChangeEIPsBandwidth(i *ChangeEIPsBandwidthInput) (*ChangeEIPsBandwidthOutput, error) {
	if i == nil {
		i = &ChangeEIPsBandwidthInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "ChangeEipsBandwidth",
		RequestMethod: "GET",
	}

	x := &ChangeEIPsBandwidthOutput{}
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

type ChangeEIPsBandwidthInput struct {
	Bandwidth *int      `json:"bandwidth" name:"bandwidth" location:"params"` // Required
	EIPs      []*string `json:"eips" name:"eips" location:"params"`           // Required
}

func (v *ChangeEIPsBandwidthInput) Validate() error {

	if v.Bandwidth == nil {
		return errors.ParameterRequiredError{
			ParameterName: "Bandwidth",
			ParentName:    "ChangeEIPsBandwidthInput",
		}
	}

	if len(v.EIPs) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "EIPs",
			ParentName:    "ChangeEIPsBandwidthInput",
		}
	}

	return nil
}

type ChangeEIPsBandwidthOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	JobID   *string `json:"job_id" name:"job_id" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/eip/change_eips_billing_mode.html
func (s *EIPService) ChangeEIPsBillingMode(i *ChangeEIPsBillingModeInput) (*ChangeEIPsBillingModeOutput, error) {
	if i == nil {
		i = &ChangeEIPsBillingModeInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "ChangeEipsBillingMode",
		RequestMethod: "GET",
	}

	x := &ChangeEIPsBillingModeOutput{}
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

type ChangeEIPsBillingModeInput struct {

	// BillingMode's available values: bandwidth, traffic
	BillingMode *string   `json:"billing_mode" name:"billing_mode" default:"bandwidth" location:"params"` // Required
	EIPGroup    *string   `json:"eip_group" name:"eip_group" location:"params"`
	EIPs        []*string `json:"eips" name:"eips" location:"params"` // Required
}

func (v *ChangeEIPsBillingModeInput) Validate() error {

	if v.BillingMode == nil {
		return errors.ParameterRequiredError{
			ParameterName: "BillingMode",
			ParentName:    "ChangeEIPsBillingModeInput",
		}
	}

	if v.BillingMode != nil {
		billingModeValidValues := []string{"bandwidth", "traffic"}
		billingModeParameterValue := fmt.Sprint(*v.BillingMode)

		billingModeIsValid := false
		for _, value := range billingModeValidValues {
			if value == billingModeParameterValue {
				billingModeIsValid = true
			}
		}

		if !billingModeIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "BillingMode",
				ParameterValue: billingModeParameterValue,
				AllowedValues:  billingModeValidValues,
			}
		}
	}

	if len(v.EIPs) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "EIPs",
			ParentName:    "ChangeEIPsBillingModeInput",
		}
	}

	return nil
}

type ChangeEIPsBillingModeOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	JobID   *string `json:"job_id" name:"job_id" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/eip/describe_eips.html
func (s *EIPService) DescribeEIPs(i *DescribeEIPsInput) (*DescribeEIPsOutput, error) {
	if i == nil {
		i = &DescribeEIPsInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DescribeEips",
		RequestMethod: "GET",
	}

	x := &DescribeEIPsOutput{}
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

type DescribeEIPsInput struct {
	EIPs       []*string `json:"eips" name:"eips" location:"params"`
	InstanceID *string   `json:"instance_id" name:"instance_id" location:"params"`
	Limit      *int      `json:"limit" name:"limit" default:"20" location:"params"`
	Offset     *int      `json:"offset" name:"offset" default:"0" location:"params"`
	SearchWord *string   `json:"search_word" name:"search_word" location:"params"`
	Status     []*string `json:"status" name:"status" location:"params"`
	Tags       []*string `json:"tags" name:"tags" location:"params"`
	Verbose    *int      `json:"verbose" name:"verbose" location:"params"`
}

func (v *DescribeEIPsInput) Validate() error {

	return nil
}

type DescribeEIPsOutput struct {
	Message    *string `json:"message" name:"message"`
	Action     *string `json:"action" name:"action" location:"elements"`
	EIPSet     []*EIP  `json:"eip_set" name:"eip_set" location:"elements"`
	RetCode    *int    `json:"ret_code" name:"ret_code" location:"elements"`
	TotalCount *int    `json:"total_count" name:"total_count" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/eip/dissociate_eips.html
func (s *EIPService) DissociateEIPs(i *DissociateEIPsInput) (*DissociateEIPsOutput, error) {
	if i == nil {
		i = &DissociateEIPsInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DissociateEips",
		RequestMethod: "GET",
	}

	x := &DissociateEIPsOutput{}
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

type DissociateEIPsInput struct {
	EIPs []*string `json:"eips" name:"eips" location:"params"` // Required
}

func (v *DissociateEIPsInput) Validate() error {

	if len(v.EIPs) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "EIPs",
			ParentName:    "DissociateEIPsInput",
		}
	}

	return nil
}

type DissociateEIPsOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	JobID   *string `json:"job_id" name:"job_id" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/monitor/get_monitor.html
func (s *EIPService) GetEIPMonitor(i *GetEIPMonitorInput) (*GetEIPMonitorOutput, error) {
	if i == nil {
		i = &GetEIPMonitorInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "GetMonitor",
		RequestMethod: "GET",
	}

	x := &GetEIPMonitorOutput{}
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

type GetEIPMonitorInput struct {
	EndTime   *time.Time `json:"end_time" name:"end_time" format:"ISO 8601" location:"params"`     // Required
	Meters    []*string  `json:"meters" name:"meters" location:"params"`                           // Required
	Resource  *string    `json:"resource" name:"resource" location:"params"`                       // Required
	StartTime *time.Time `json:"start_time" name:"start_time" format:"ISO 8601" location:"params"` // Required
	// Step's available values: 5m, 15m, 2h, 1d
	Step *string `json:"step" name:"step" location:"params"` // Required
}

func (v *GetEIPMonitorInput) Validate() error {

	if len(v.Meters) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "Meters",
			ParentName:    "GetEIPMonitorInput",
		}
	}

	if v.Resource == nil {
		return errors.ParameterRequiredError{
			ParameterName: "Resource",
			ParentName:    "GetEIPMonitorInput",
		}
	}

	if v.Step == nil {
		return errors.ParameterRequiredError{
			ParameterName: "Step",
			ParentName:    "GetEIPMonitorInput",
		}
	}

	if v.Step != nil {
		stepValidValues := []string{"5m", "15m", "2h", "1d"}
		stepParameterValue := fmt.Sprint(*v.Step)

		stepIsValid := false
		for _, value := range stepValidValues {
			if value == stepParameterValue {
				stepIsValid = true
			}
		}

		if !stepIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "Step",
				ParameterValue: stepParameterValue,
				AllowedValues:  stepValidValues,
			}
		}
	}

	return nil
}

type GetEIPMonitorOutput struct {
	Message    *string  `json:"message" name:"message"`
	Action     *string  `json:"action" name:"action" location:"elements"`
	MeterSet   []*Meter `json:"meter_set" name:"meter_set" location:"elements"`
	ResourceID *string  `json:"resource_id" name:"resource_id" location:"elements"`
	RetCode    *int     `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/eip/modify_eip_attributes.html
func (s *EIPService) ModifyEIPAttributes(i *ModifyEIPAttributesInput) (*ModifyEIPAttributesOutput, error) {
	if i == nil {
		i = &ModifyEIPAttributesInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "ModifyEipAttributes",
		RequestMethod: "GET",
	}

	x := &ModifyEIPAttributesOutput{}
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

type ModifyEIPAttributesInput struct {
	Description *string `json:"description" name:"description" location:"params"`
	EIP         *string `json:"eip" name:"eip" location:"params"` // Required
	EIPName     *string `json:"eip_name" name:"eip_name" location:"params"`
}

func (v *ModifyEIPAttributesInput) Validate() error {

	if v.EIP == nil {
		return errors.ParameterRequiredError{
			ParameterName: "EIP",
			ParentName:    "ModifyEIPAttributesInput",
		}
	}

	return nil
}

type ModifyEIPAttributesOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	EIPID   *string `json:"eip_id" name:"eip_id" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/eip/release_eips.html
func (s *EIPService) ReleaseEIPs(i *ReleaseEIPsInput) (*ReleaseEIPsOutput, error) {
	if i == nil {
		i = &ReleaseEIPsInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "ReleaseEips",
		RequestMethod: "GET",
	}

	x := &ReleaseEIPsOutput{}
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

type ReleaseEIPsInput struct {
	EIPs []*string `json:"eips" name:"eips" location:"params"` // Required
}

func (v *ReleaseEIPsInput) Validate() error {

	if len(v.EIPs) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "EIPs",
			ParentName:    "ReleaseEIPsInput",
		}
	}

	return nil
}

type ReleaseEIPsOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	JobID   *string `json:"job_id" name:"job_id" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}
