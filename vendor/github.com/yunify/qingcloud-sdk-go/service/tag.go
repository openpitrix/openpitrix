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

type TagService struct {
	Config     *config.Config
	Properties *TagServiceProperties
}

type TagServiceProperties struct {
	// QingCloud Zone ID
	Zone *string `json:"zone" name:"zone"` // Required
}

func (s *QingCloudService) Tag(zone string) (*TagService, error) {
	properties := &TagServiceProperties{
		Zone: &zone,
	}

	return &TagService{Config: s.Config, Properties: properties}, nil
}

// Documentation URL: https://docs.qingcloud.com/api/tag/attach_tags.html
func (s *TagService) AttachTags(i *AttachTagsInput) (*AttachTagsOutput, error) {
	if i == nil {
		i = &AttachTagsInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "AttachTags",
		RequestMethod: "GET",
	}

	x := &AttachTagsOutput{}
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

type AttachTagsInput struct {
	ResourceTagPairs []*ResourceTagPair `json:"resource_tag_pairs" name:"resource_tag_pairs" location:"params"` // Required
}

func (v *AttachTagsInput) Validate() error {

	if len(v.ResourceTagPairs) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "ResourceTagPairs",
			ParentName:    "AttachTagsInput",
		}
	}

	if len(v.ResourceTagPairs) > 0 {
		for _, property := range v.ResourceTagPairs {
			if err := property.Validate(); err != nil {
				return err
			}
		}
	}

	return nil
}

type AttachTagsOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/tag/create_tag.html
func (s *TagService) CreateTag(i *CreateTagInput) (*CreateTagOutput, error) {
	if i == nil {
		i = &CreateTagInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "CreateTag",
		RequestMethod: "GET",
	}

	x := &CreateTagOutput{}
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

type CreateTagInput struct {
	Color   *string `json:"color" name:"color" location:"params"`
	TagName *string `json:"tag_name" name:"tag_name" location:"params"`
}

func (v *CreateTagInput) Validate() error {

	return nil
}

type CreateTagOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
	TagID   *string `json:"tag_id" name:"tag_id" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/tag/delete_tags.html
func (s *TagService) DeleteTags(i *DeleteTagsInput) (*DeleteTagsOutput, error) {
	if i == nil {
		i = &DeleteTagsInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DeleteTags",
		RequestMethod: "GET",
	}

	x := &DeleteTagsOutput{}
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

type DeleteTagsInput struct {
	Tags []*string `json:"tags" name:"tags" location:"params"` // Required
}

func (v *DeleteTagsInput) Validate() error {

	if len(v.Tags) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "Tags",
			ParentName:    "DeleteTagsInput",
		}
	}

	return nil
}

type DeleteTagsOutput struct {
	Message *string   `json:"message" name:"message"`
	Action  *string   `json:"action" name:"action" location:"elements"`
	RetCode *int      `json:"ret_code" name:"ret_code" location:"elements"`
	Tags    []*string `json:"tags" name:"tags" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/tag/describe_tags.html
func (s *TagService) DescribeTags(i *DescribeTagsInput) (*DescribeTagsOutput, error) {
	if i == nil {
		i = &DescribeTagsInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DescribeTags",
		RequestMethod: "GET",
	}

	x := &DescribeTagsOutput{}
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

type DescribeTagsInput struct {
	Limit      *int      `json:"limit" name:"limit" default:"0" location:"params"`
	Offset     *int      `json:"offset" name:"offset" default:"0" location:"params"`
	SearchWord *string   `json:"search_word" name:"search_word" location:"params"`
	Tags       []*string `json:"tags" name:"tags" location:"params"`
	// Verbose's available values: 0, 1
	Verbose *int `json:"verbose" name:"verbose" default:"0" location:"params"`
}

func (v *DescribeTagsInput) Validate() error {

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

type DescribeTagsOutput struct {
	Message    *string `json:"message" name:"message"`
	Action     *string `json:"action" name:"action" location:"elements"`
	RetCode    *int    `json:"ret_code" name:"ret_code" location:"elements"`
	TagSet     []*Tag  `json:"tag_set" name:"tag_set" location:"elements"`
	TotalCount *int    `json:"total_count" name:"total_count" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/tag/detach_tags.html
func (s *TagService) DetachTags(i *DetachTagsInput) (*DetachTagsOutput, error) {
	if i == nil {
		i = &DetachTagsInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DetachTags",
		RequestMethod: "GET",
	}

	x := &DetachTagsOutput{}
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

type DetachTagsInput struct {
	ResourceTagPairs []*ResourceTagPair `json:"resource_tag_pairs" name:"resource_tag_pairs" location:"params"` // Required
}

func (v *DetachTagsInput) Validate() error {

	if len(v.ResourceTagPairs) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "ResourceTagPairs",
			ParentName:    "DetachTagsInput",
		}
	}

	if len(v.ResourceTagPairs) > 0 {
		for _, property := range v.ResourceTagPairs {
			if err := property.Validate(); err != nil {
				return err
			}
		}
	}

	return nil
}

type DetachTagsOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/tag/modify_tag_attributes.html
func (s *TagService) ModifyTagAttributes(i *ModifyTagAttributesInput) (*ModifyTagAttributesOutput, error) {
	if i == nil {
		i = &ModifyTagAttributesInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "ModifyTagAttributes",
		RequestMethod: "GET",
	}

	x := &ModifyTagAttributesOutput{}
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

type ModifyTagAttributesInput struct {
	Color       *string `json:"color" name:"color" location:"params"`
	Description *string `json:"description" name:"description" location:"params"`
	Tag         *string `json:"tag" name:"tag" location:"params"` // Required
	TagName     *string `json:"tag_name" name:"tag_name" location:"params"`
}

func (v *ModifyTagAttributesInput) Validate() error {

	if v.Tag == nil {
		return errors.ParameterRequiredError{
			ParameterName: "Tag",
			ParentName:    "ModifyTagAttributesInput",
		}
	}

	return nil
}

type ModifyTagAttributesOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}
