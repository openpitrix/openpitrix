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

type SecurityGroupService struct {
	Config     *config.Config
	Properties *SecurityGroupServiceProperties
}

type SecurityGroupServiceProperties struct {
	// QingCloud Zone ID
	Zone *string `json:"zone" name:"zone"` // Required
}

func (s *QingCloudService) SecurityGroup(zone string) (*SecurityGroupService, error) {
	properties := &SecurityGroupServiceProperties{
		Zone: &zone,
	}

	return &SecurityGroupService{Config: s.Config, Properties: properties}, nil
}

// Documentation URL: https://docs.qingcloud.com/api/sg/add_security_group_rules.html
func (s *SecurityGroupService) AddSecurityGroupRules(i *AddSecurityGroupRulesInput) (*AddSecurityGroupRulesOutput, error) {
	if i == nil {
		i = &AddSecurityGroupRulesInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "AddSecurityGroupRules",
		RequestMethod: "GET",
	}

	x := &AddSecurityGroupRulesOutput{}
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

type AddSecurityGroupRulesInput struct {
	Rules         []*SecurityGroupRule `json:"rules" name:"rules" location:"params"`                   // Required
	SecurityGroup *string              `json:"security_group" name:"security_group" location:"params"` // Required
}

func (v *AddSecurityGroupRulesInput) Validate() error {

	if len(v.Rules) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "Rules",
			ParentName:    "AddSecurityGroupRulesInput",
		}
	}

	if len(v.Rules) > 0 {
		for _, property := range v.Rules {
			if err := property.Validate(); err != nil {
				return err
			}
		}
	}

	if v.SecurityGroup == nil {
		return errors.ParameterRequiredError{
			ParameterName: "SecurityGroup",
			ParentName:    "AddSecurityGroupRulesInput",
		}
	}

	return nil
}

type AddSecurityGroupRulesOutput struct {
	Message            *string   `json:"message" name:"message"`
	Action             *string   `json:"action" name:"action" location:"elements"`
	RetCode            *int      `json:"ret_code" name:"ret_code" location:"elements"`
	SecurityGroupRules []*string `json:"security_group_rules" name:"security_group_rules" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/sg/apply_security_group.html
func (s *SecurityGroupService) ApplySecurityGroup(i *ApplySecurityGroupInput) (*ApplySecurityGroupOutput, error) {
	if i == nil {
		i = &ApplySecurityGroupInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "ApplySecurityGroup",
		RequestMethod: "GET",
	}

	x := &ApplySecurityGroupOutput{}
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

type ApplySecurityGroupInput struct {
	Instances     []*string `json:"instances" name:"instances" location:"params"`
	SecurityGroup *string   `json:"security_group" name:"security_group" location:"params"` // Required
}

func (v *ApplySecurityGroupInput) Validate() error {

	if v.SecurityGroup == nil {
		return errors.ParameterRequiredError{
			ParameterName: "SecurityGroup",
			ParentName:    "ApplySecurityGroupInput",
		}
	}

	return nil
}

type ApplySecurityGroupOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	JobID   *string `json:"job_id" name:"job_id" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/sg/create_security_group.html
func (s *SecurityGroupService) CreateSecurityGroup(i *CreateSecurityGroupInput) (*CreateSecurityGroupOutput, error) {
	if i == nil {
		i = &CreateSecurityGroupInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "CreateSecurityGroup",
		RequestMethod: "GET",
	}

	x := &CreateSecurityGroupOutput{}
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

type CreateSecurityGroupInput struct {
	SecurityGroupName *string `json:"security_group_name" name:"security_group_name" location:"params"`
}

func (v *CreateSecurityGroupInput) Validate() error {

	return nil
}

type CreateSecurityGroupOutput struct {
	Message         *string `json:"message" name:"message"`
	Action          *string `json:"action" name:"action" location:"elements"`
	RetCode         *int    `json:"ret_code" name:"ret_code" location:"elements"`
	SecurityGroupID *string `json:"security_group_id" name:"security_group_id" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/sg/create_security_group_ipset.html
func (s *SecurityGroupService) CreateSecurityGroupIPSet(i *CreateSecurityGroupIPSetInput) (*CreateSecurityGroupIPSetOutput, error) {
	if i == nil {
		i = &CreateSecurityGroupIPSetInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "CreateSecurityGroupIPSet",
		RequestMethod: "GET",
	}

	x := &CreateSecurityGroupIPSetOutput{}
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

type CreateSecurityGroupIPSetInput struct {

	// IPSetType's available values: 0, 1
	IPSetType              *int    `json:"ipset_type" name:"ipset_type" location:"params"` // Required
	SecurityGroupIPSetName *string `json:"security_group_ipset_name" name:"security_group_ipset_name" location:"params"`
	Val                    *string `json:"val" name:"val" location:"params"` // Required
}

func (v *CreateSecurityGroupIPSetInput) Validate() error {

	if v.IPSetType == nil {
		return errors.ParameterRequiredError{
			ParameterName: "IPSetType",
			ParentName:    "CreateSecurityGroupIPSetInput",
		}
	}

	if v.IPSetType != nil {
		ipSetTypeValidValues := []string{"0", "1"}
		ipSetTypeParameterValue := fmt.Sprint(*v.IPSetType)

		ipSetTypeIsValid := false
		for _, value := range ipSetTypeValidValues {
			if value == ipSetTypeParameterValue {
				ipSetTypeIsValid = true
			}
		}

		if !ipSetTypeIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "IPSetType",
				ParameterValue: ipSetTypeParameterValue,
				AllowedValues:  ipSetTypeValidValues,
			}
		}
	}

	if v.Val == nil {
		return errors.ParameterRequiredError{
			ParameterName: "Val",
			ParentName:    "CreateSecurityGroupIPSetInput",
		}
	}

	return nil
}

type CreateSecurityGroupIPSetOutput struct {
	Message              *string `json:"message" name:"message"`
	Action               *string `json:"action" name:"action" location:"elements"`
	RetCode              *int    `json:"ret_code" name:"ret_code" location:"elements"`
	SecurityGroupIPSetID *string `json:"security_group_ipset_id" name:"security_group_ipset_id" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/sg/create_security_group_snapshot.html
func (s *SecurityGroupService) CreateSecurityGroupSnapshot(i *CreateSecurityGroupSnapshotInput) (*CreateSecurityGroupSnapshotOutput, error) {
	if i == nil {
		i = &CreateSecurityGroupSnapshotInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "CreateSecurityGroupSnapshot",
		RequestMethod: "GET",
	}

	x := &CreateSecurityGroupSnapshotOutput{}
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

type CreateSecurityGroupSnapshotInput struct {
	Name          *string `json:"name" name:"name" location:"params"`
	SecurityGroup *string `json:"security_group" name:"security_group" location:"params"` // Required
}

func (v *CreateSecurityGroupSnapshotInput) Validate() error {

	if v.SecurityGroup == nil {
		return errors.ParameterRequiredError{
			ParameterName: "SecurityGroup",
			ParentName:    "CreateSecurityGroupSnapshotInput",
		}
	}

	return nil
}

type CreateSecurityGroupSnapshotOutput struct {
	Message                 *string `json:"message" name:"message"`
	Action                  *string `json:"action" name:"action" location:"elements"`
	RetCode                 *int    `json:"ret_code" name:"ret_code" location:"elements"`
	SecurityGroupID         *string `json:"security_group_id" name:"security_group_id" location:"elements"`
	SecurityGroupSnapshotID *string `json:"security_group_snapshot_id" name:"security_group_snapshot_id" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/sg/delete_security_group_ipsets.html
func (s *SecurityGroupService) DeleteSecurityGroupIPSets(i *DeleteSecurityGroupIPSetsInput) (*DeleteSecurityGroupIPSetsOutput, error) {
	if i == nil {
		i = &DeleteSecurityGroupIPSetsInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DeleteSecurityGroupIPSets",
		RequestMethod: "GET",
	}

	x := &DeleteSecurityGroupIPSetsOutput{}
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

type DeleteSecurityGroupIPSetsInput struct {
	SecurityGroupIPSets []*string `json:"security_group_ipsets" name:"security_group_ipsets" location:"params"` // Required
}

func (v *DeleteSecurityGroupIPSetsInput) Validate() error {

	if len(v.SecurityGroupIPSets) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "SecurityGroupIPSets",
			ParentName:    "DeleteSecurityGroupIPSetsInput",
		}
	}

	return nil
}

type DeleteSecurityGroupIPSetsOutput struct {
	Message             *string   `json:"message" name:"message"`
	Action              *string   `json:"action" name:"action" location:"elements"`
	RetCode             *int      `json:"ret_code" name:"ret_code" location:"elements"`
	SecurityGroupIPSets []*string `json:"security_group_ipsets" name:"security_group_ipsets" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/sg/delete_security_group_rules.html
func (s *SecurityGroupService) DeleteSecurityGroupRules(i *DeleteSecurityGroupRulesInput) (*DeleteSecurityGroupRulesOutput, error) {
	if i == nil {
		i = &DeleteSecurityGroupRulesInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DeleteSecurityGroupRules",
		RequestMethod: "GET",
	}

	x := &DeleteSecurityGroupRulesOutput{}
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

type DeleteSecurityGroupRulesInput struct {
	SecurityGroupRules []*string `json:"security_group_rules" name:"security_group_rules" location:"params"` // Required
}

func (v *DeleteSecurityGroupRulesInput) Validate() error {

	if len(v.SecurityGroupRules) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "SecurityGroupRules",
			ParentName:    "DeleteSecurityGroupRulesInput",
		}
	}

	return nil
}

type DeleteSecurityGroupRulesOutput struct {
	Message            *string   `json:"message" name:"message"`
	Action             *string   `json:"action" name:"action" location:"elements"`
	RetCode            *int      `json:"ret_code" name:"ret_code" location:"elements"`
	SecurityGroupRules []*string `json:"security_group_rules" name:"security_group_rules" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/sg/delete_security_group_snapshots.html
func (s *SecurityGroupService) DeleteSecurityGroupSnapshots(i *DeleteSecurityGroupSnapshotsInput) (*DeleteSecurityGroupSnapshotsOutput, error) {
	if i == nil {
		i = &DeleteSecurityGroupSnapshotsInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DeleteSecurityGroupSnapshots",
		RequestMethod: "GET",
	}

	x := &DeleteSecurityGroupSnapshotsOutput{}
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

type DeleteSecurityGroupSnapshotsInput struct {
	SecurityGroupSnapshots []*string `json:"security_group_snapshots" name:"security_group_snapshots" location:"params"` // Required
}

func (v *DeleteSecurityGroupSnapshotsInput) Validate() error {

	if len(v.SecurityGroupSnapshots) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "SecurityGroupSnapshots",
			ParentName:    "DeleteSecurityGroupSnapshotsInput",
		}
	}

	return nil
}

type DeleteSecurityGroupSnapshotsOutput struct {
	Message                *string   `json:"message" name:"message"`
	Action                 *string   `json:"action" name:"action" location:"elements"`
	RetCode                *int      `json:"ret_code" name:"ret_code" location:"elements"`
	SecurityGroupSnapshots []*string `json:"security_group_snapshots" name:"security_group_snapshots" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/sg/delete_security_groups.html
func (s *SecurityGroupService) DeleteSecurityGroups(i *DeleteSecurityGroupsInput) (*DeleteSecurityGroupsOutput, error) {
	if i == nil {
		i = &DeleteSecurityGroupsInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DeleteSecurityGroups",
		RequestMethod: "GET",
	}

	x := &DeleteSecurityGroupsOutput{}
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

type DeleteSecurityGroupsInput struct {
	SecurityGroups []*string `json:"security_groups" name:"security_groups" location:"params"` // Required
}

func (v *DeleteSecurityGroupsInput) Validate() error {

	if len(v.SecurityGroups) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "SecurityGroups",
			ParentName:    "DeleteSecurityGroupsInput",
		}
	}

	return nil
}

type DeleteSecurityGroupsOutput struct {
	Message        *string   `json:"message" name:"message"`
	Action         *string   `json:"action" name:"action" location:"elements"`
	RetCode        *int      `json:"ret_code" name:"ret_code" location:"elements"`
	SecurityGroups []*string `json:"security_groups" name:"security_groups" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/sg/describe_security_group_ipsets.html
func (s *SecurityGroupService) DescribeSecurityGroupIPSets(i *DescribeSecurityGroupIPSetsInput) (*DescribeSecurityGroupIPSetsOutput, error) {
	if i == nil {
		i = &DescribeSecurityGroupIPSetsInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DescribeSecurityGroupIPSets",
		RequestMethod: "GET",
	}

	x := &DescribeSecurityGroupIPSetsOutput{}
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

type DescribeSecurityGroupIPSetsInput struct {

	// IPSetType's available values: 0, 1
	IPSetType              *int      `json:"ipset_type" name:"ipset_type" location:"params"`
	Limit                  *int      `json:"limit" name:"limit" default:"20" location:"params"`
	Offset                 *int      `json:"offset" name:"offset" default:"0" location:"params"`
	SecurityGroupIPSetName *string   `json:"security_group_ipset_name" name:"security_group_ipset_name" location:"params"`
	SecurityGroupIPSets    []*string `json:"security_group_ipsets" name:"security_group_ipsets" location:"params"`
	Tags                   []*string `json:"tags" name:"tags" location:"params"`
	Verbose                *int      `json:"verbose" name:"verbose" default:"0" location:"params"`
}

func (v *DescribeSecurityGroupIPSetsInput) Validate() error {

	if v.IPSetType != nil {
		ipSetTypeValidValues := []string{"0", "1"}
		ipSetTypeParameterValue := fmt.Sprint(*v.IPSetType)

		ipSetTypeIsValid := false
		for _, value := range ipSetTypeValidValues {
			if value == ipSetTypeParameterValue {
				ipSetTypeIsValid = true
			}
		}

		if !ipSetTypeIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "IPSetType",
				ParameterValue: ipSetTypeParameterValue,
				AllowedValues:  ipSetTypeValidValues,
			}
		}
	}

	return nil
}

type DescribeSecurityGroupIPSetsOutput struct {
	Message               *string               `json:"message" name:"message"`
	Action                *string               `json:"action" name:"action" location:"elements"`
	RetCode               *int                  `json:"ret_code" name:"ret_code" location:"elements"`
	SecurityGroupIPSetSet []*SecurityGroupIPSet `json:"security_group_ipset_set" name:"security_group_ipset_set" location:"elements"`
	TotalCount            *int                  `json:"total_count" name:"total_count" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/sg/describe_security_group_rules.html
func (s *SecurityGroupService) DescribeSecurityGroupRules(i *DescribeSecurityGroupRulesInput) (*DescribeSecurityGroupRulesOutput, error) {
	if i == nil {
		i = &DescribeSecurityGroupRulesInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DescribeSecurityGroupRules",
		RequestMethod: "GET",
	}

	x := &DescribeSecurityGroupRulesOutput{}
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

type DescribeSecurityGroupRulesInput struct {

	// Direction's available values: 0, 1
	Direction          *int      `json:"direction" name:"direction" location:"params"`
	Limit              *int      `json:"limit" name:"limit" default:"20" location:"params"`
	Offset             *int      `json:"offset" name:"offset" default:"0" location:"params"`
	SecurityGroup      *string   `json:"security_group" name:"security_group" location:"params"`
	SecurityGroupRules []*string `json:"security_group_rules" name:"security_group_rules" location:"params"`
}

func (v *DescribeSecurityGroupRulesInput) Validate() error {

	if v.Direction != nil {
		directionValidValues := []string{"0", "1"}
		directionParameterValue := fmt.Sprint(*v.Direction)

		directionIsValid := false
		for _, value := range directionValidValues {
			if value == directionParameterValue {
				directionIsValid = true
			}
		}

		if !directionIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "Direction",
				ParameterValue: directionParameterValue,
				AllowedValues:  directionValidValues,
			}
		}
	}

	return nil
}

type DescribeSecurityGroupRulesOutput struct {
	Message              *string              `json:"message" name:"message"`
	Action               *string              `json:"action" name:"action" location:"elements"`
	RetCode              *int                 `json:"ret_code" name:"ret_code" location:"elements"`
	SecurityGroupRuleSet []*SecurityGroupRule `json:"security_group_rule_set" name:"security_group_rule_set" location:"elements"`
	TotalCount           *int                 `json:"total_count" name:"total_count" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/sg/describe_security_group_snapshots.html
func (s *SecurityGroupService) DescribeSecurityGroupSnapshots(i *DescribeSecurityGroupSnapshotsInput) (*DescribeSecurityGroupSnapshotsOutput, error) {
	if i == nil {
		i = &DescribeSecurityGroupSnapshotsInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DescribeSecurityGroupSnapshots",
		RequestMethod: "GET",
	}

	x := &DescribeSecurityGroupSnapshotsOutput{}
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

type DescribeSecurityGroupSnapshotsInput struct {
	Limit                  *int      `json:"limit" name:"limit" default:"20" location:"params"`
	Offset                 *int      `json:"offset" name:"offset" default:"0" location:"params"`
	Reverse                *int      `json:"reverse" name:"reverse" default:"1" location:"params"`
	SecurityGroup          *string   `json:"security_group" name:"security_group" location:"params"` // Required
	SecurityGroupSnapshots []*string `json:"security_group_snapshots" name:"security_group_snapshots" location:"params"`
}

func (v *DescribeSecurityGroupSnapshotsInput) Validate() error {

	if v.SecurityGroup == nil {
		return errors.ParameterRequiredError{
			ParameterName: "SecurityGroup",
			ParentName:    "DescribeSecurityGroupSnapshotsInput",
		}
	}

	return nil
}

type DescribeSecurityGroupSnapshotsOutput struct {
	Message                  *string                  `json:"message" name:"message"`
	Action                   *string                  `json:"action" name:"action" location:"elements"`
	RetCode                  *int                     `json:"ret_code" name:"ret_code" location:"elements"`
	SecurityGroupSnapshotSet []*SecurityGroupSnapshot `json:"security_group_snapshot_set" name:"security_group_snapshot_set" location:"elements"`
	TotalCount               *int                     `json:"total_count" name:"total_count" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/sg/describe_security_groups.html
func (s *SecurityGroupService) DescribeSecurityGroups(i *DescribeSecurityGroupsInput) (*DescribeSecurityGroupsOutput, error) {
	if i == nil {
		i = &DescribeSecurityGroupsInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DescribeSecurityGroups",
		RequestMethod: "GET",
	}

	x := &DescribeSecurityGroupsOutput{}
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

type DescribeSecurityGroupsInput struct {
	Limit          *int      `json:"limit" name:"limit" default:"20" location:"params"`
	Offset         *int      `json:"offset" name:"offset" default:"0" location:"params"`
	SearchWord     *string   `json:"search_word" name:"search_word" location:"params"`
	SecurityGroups []*string `json:"security_groups" name:"security_groups" location:"params"`
	Tags           []*string `json:"tags" name:"tags" location:"params"`
	Verbose        *int      `json:"verbose" name:"verbose" default:"0" location:"params"`
}

func (v *DescribeSecurityGroupsInput) Validate() error {

	return nil
}

type DescribeSecurityGroupsOutput struct {
	Message          *string          `json:"message" name:"message"`
	Action           *string          `json:"action" name:"action" location:"elements"`
	RetCode          *int             `json:"ret_code" name:"ret_code" location:"elements"`
	SecurityGroupSet []*SecurityGroup `json:"security_group_set" name:"security_group_set" location:"elements"`
	TotalCount       *int             `json:"total_count" name:"total_count" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/sg/modify_security_group_attributes.html
func (s *SecurityGroupService) ModifySecurityGroupAttributes(i *ModifySecurityGroupAttributesInput) (*ModifySecurityGroupAttributesOutput, error) {
	if i == nil {
		i = &ModifySecurityGroupAttributesInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "ModifySecurityGroupAttributes",
		RequestMethod: "GET",
	}

	x := &ModifySecurityGroupAttributesOutput{}
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

type ModifySecurityGroupAttributesInput struct {
	Description       *string `json:"description" name:"description" location:"params"`
	SecurityGroup     *string `json:"security_group" name:"security_group" location:"params"` // Required
	SecurityGroupName *string `json:"security_group_name" name:"security_group_name" location:"params"`
}

func (v *ModifySecurityGroupAttributesInput) Validate() error {

	if v.SecurityGroup == nil {
		return errors.ParameterRequiredError{
			ParameterName: "SecurityGroup",
			ParentName:    "ModifySecurityGroupAttributesInput",
		}
	}

	return nil
}

type ModifySecurityGroupAttributesOutput struct {
	Message         *string `json:"message" name:"message"`
	Action          *string `json:"action" name:"action" location:"elements"`
	RetCode         *int    `json:"ret_code" name:"ret_code" location:"elements"`
	SecurityGroupID *string `json:"security_group_id" name:"security_group_id" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/sg/modify_security_group_ipset_attributes.html
func (s *SecurityGroupService) ModifySecurityGroupIPSetAttributes(i *ModifySecurityGroupIPSetAttributesInput) (*ModifySecurityGroupIPSetAttributesOutput, error) {
	if i == nil {
		i = &ModifySecurityGroupIPSetAttributesInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "ModifySecurityGroupIPSetAttributes",
		RequestMethod: "GET",
	}

	x := &ModifySecurityGroupIPSetAttributesOutput{}
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

type ModifySecurityGroupIPSetAttributesInput struct {
	Description            *string `json:"description" name:"description" location:"params"`
	SecurityGroupIPSet     *string `json:"security_group_ipset" name:"security_group_ipset" location:"params"` // Required
	SecurityGroupIPSetName *string `json:"security_group_ipset_name" name:"security_group_ipset_name" location:"params"`
	Val                    *string `json:"val" name:"val" location:"params"`
}

func (v *ModifySecurityGroupIPSetAttributesInput) Validate() error {

	if v.SecurityGroupIPSet == nil {
		return errors.ParameterRequiredError{
			ParameterName: "SecurityGroupIPSet",
			ParentName:    "ModifySecurityGroupIPSetAttributesInput",
		}
	}

	return nil
}

type ModifySecurityGroupIPSetAttributesOutput struct {
	Message              *string `json:"message" name:"message"`
	Action               *string `json:"action" name:"action" location:"elements"`
	RetCode              *int    `json:"ret_code" name:"ret_code" location:"elements"`
	SecurityGroupIPSetID *string `json:"security_group_ipset_id" name:"security_group_ipset_id" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/sg/modify_security_group_rule_attributes.html
func (s *SecurityGroupService) ModifySecurityGroupRuleAttributes(i *ModifySecurityGroupRuleAttributesInput) (*ModifySecurityGroupRuleAttributesOutput, error) {
	if i == nil {
		i = &ModifySecurityGroupRuleAttributesInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "ModifySecurityGroupRuleAttributes",
		RequestMethod: "GET",
	}

	x := &ModifySecurityGroupRuleAttributesOutput{}
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

type ModifySecurityGroupRuleAttributesInput struct {

	// Direction's available values: 0, 1
	Direction *int    `json:"direction" name:"direction" location:"params"`
	Priority  *int    `json:"priority" name:"priority" location:"params"`
	Protocol  *string `json:"protocol" name:"protocol" location:"params"`
	// RuleAction's available values: accept, drop
	RuleAction            *string `json:"rule_action" name:"rule_action" location:"params"`
	SecurityGroup         *string `json:"security_group" name:"security_group" location:"params"`
	SecurityGroupRule     *string `json:"security_group_rule" name:"security_group_rule" location:"params"` // Required
	SecurityGroupRuleName *string `json:"security_group_rule_name" name:"security_group_rule_name" location:"params"`
	Val1                  *string `json:"val1" name:"val1" location:"params"`
	Val2                  *string `json:"val2" name:"val2" location:"params"`
	Val3                  *string `json:"val3" name:"val3" location:"params"`
}

func (v *ModifySecurityGroupRuleAttributesInput) Validate() error {

	if v.Direction != nil {
		directionValidValues := []string{"0", "1"}
		directionParameterValue := fmt.Sprint(*v.Direction)

		directionIsValid := false
		for _, value := range directionValidValues {
			if value == directionParameterValue {
				directionIsValid = true
			}
		}

		if !directionIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "Direction",
				ParameterValue: directionParameterValue,
				AllowedValues:  directionValidValues,
			}
		}
	}

	if v.RuleAction != nil {
		ruleActionValidValues := []string{"accept", "drop"}
		ruleActionParameterValue := fmt.Sprint(*v.RuleAction)

		ruleActionIsValid := false
		for _, value := range ruleActionValidValues {
			if value == ruleActionParameterValue {
				ruleActionIsValid = true
			}
		}

		if !ruleActionIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "RuleAction",
				ParameterValue: ruleActionParameterValue,
				AllowedValues:  ruleActionValidValues,
			}
		}
	}

	if v.SecurityGroupRule == nil {
		return errors.ParameterRequiredError{
			ParameterName: "SecurityGroupRule",
			ParentName:    "ModifySecurityGroupRuleAttributesInput",
		}
	}

	return nil
}

type ModifySecurityGroupRuleAttributesOutput struct {
	Message             *string `json:"message" name:"message"`
	Action              *string `json:"action" name:"action" location:"elements"`
	RetCode             *int    `json:"ret_code" name:"ret_code" location:"elements"`
	SecurityGroupRuleID *string `json:"security_group_rule_id" name:"security_group_rule_id" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/sg/rollback_security_group.html
func (s *SecurityGroupService) RollbackSecurityGroup(i *RollbackSecurityGroupInput) (*RollbackSecurityGroupOutput, error) {
	if i == nil {
		i = &RollbackSecurityGroupInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "RollbackSecurityGroup",
		RequestMethod: "GET",
	}

	x := &RollbackSecurityGroupOutput{}
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

type RollbackSecurityGroupInput struct {
	SecurityGroup         *string `json:"security_group" name:"security_group" location:"params"`                   // Required
	SecurityGroupSnapshot *string `json:"security_group_snapshot" name:"security_group_snapshot" location:"params"` // Required
}

func (v *RollbackSecurityGroupInput) Validate() error {

	if v.SecurityGroup == nil {
		return errors.ParameterRequiredError{
			ParameterName: "SecurityGroup",
			ParentName:    "RollbackSecurityGroupInput",
		}
	}

	if v.SecurityGroupSnapshot == nil {
		return errors.ParameterRequiredError{
			ParameterName: "SecurityGroupSnapshot",
			ParentName:    "RollbackSecurityGroupInput",
		}
	}

	return nil
}

type RollbackSecurityGroupOutput struct {
	Message                 *string `json:"message" name:"message"`
	Action                  *string `json:"action" name:"action" location:"elements"`
	RetCode                 *int    `json:"ret_code" name:"ret_code" location:"elements"`
	SecurityGroupID         *string `json:"security_group_id" name:"security_group_id" location:"elements"`
	SecurityGroupSnapshotID *string `json:"security_group_snapshot_id" name:"security_group_snapshot_id" location:"elements"`
}
