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

type ImageService struct {
	Config     *config.Config
	Properties *ImageServiceProperties
}

type ImageServiceProperties struct {
	// QingCloud Zone ID
	Zone *string `json:"zone" name:"zone"` // Required
}

func (s *QingCloudService) Image(zone string) (*ImageService, error) {
	properties := &ImageServiceProperties{
		Zone: &zone,
	}

	return &ImageService{Config: s.Config, Properties: properties}, nil
}

// Documentation URL: https://docs.qingcloud.com/api/image/capture_instance.html
func (s *ImageService) CaptureInstance(i *CaptureInstanceInput) (*CaptureInstanceOutput, error) {
	if i == nil {
		i = &CaptureInstanceInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "CaptureInstance",
		RequestMethod: "GET",
	}

	x := &CaptureInstanceOutput{}
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

type CaptureInstanceInput struct {
	ImageName *string `json:"image_name" name:"image_name" location:"params"`
	Instance  *string `json:"instance" name:"instance" location:"params"` // Required
}

func (v *CaptureInstanceInput) Validate() error {

	if v.Instance == nil {
		return errors.ParameterRequiredError{
			ParameterName: "Instance",
			ParentName:    "CaptureInstanceInput",
		}
	}

	return nil
}

type CaptureInstanceOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	ImageID *string `json:"image_id" name:"image_id" location:"elements"`
	JobID   *string `json:"job_id" name:"job_id" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/image/delete_images.html
func (s *ImageService) DeleteImages(i *DeleteImagesInput) (*DeleteImagesOutput, error) {
	if i == nil {
		i = &DeleteImagesInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DeleteImages",
		RequestMethod: "GET",
	}

	x := &DeleteImagesOutput{}
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

type DeleteImagesInput struct {
	Images []*string `json:"images" name:"images" location:"params"` // Required
}

func (v *DeleteImagesInput) Validate() error {

	if len(v.Images) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "Images",
			ParentName:    "DeleteImagesInput",
		}
	}

	return nil
}

type DeleteImagesOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	JobID   *string `json:"job_id" name:"job_id" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/image/describe-image-users.html
func (s *ImageService) DescribeImageUsers(i *DescribeImageUsersInput) (*DescribeImageUsersOutput, error) {
	if i == nil {
		i = &DescribeImageUsersInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DescribeImageUsers",
		RequestMethod: "GET",
	}

	x := &DescribeImageUsersOutput{}
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

type DescribeImageUsersInput struct {
	ImageID *string `json:"image_id" name:"image_id" location:"params"` // Required
	Limit   *int    `json:"limit" name:"limit" default:"20" location:"params"`
	Offset  *int    `json:"offset" name:"offset" default:"0" location:"params"`
}

func (v *DescribeImageUsersInput) Validate() error {

	if v.ImageID == nil {
		return errors.ParameterRequiredError{
			ParameterName: "ImageID",
			ParentName:    "DescribeImageUsersInput",
		}
	}

	return nil
}

type DescribeImageUsersOutput struct {
	Message      *string      `json:"message" name:"message"`
	Action       *string      `json:"action" name:"action" location:"elements"`
	ImageUserSet []*ImageUser `json:"image_user_set" name:"image_user_set" location:"elements"`
	RetCode      *int         `json:"ret_code" name:"ret_code" location:"elements"`
	TotalCount   *int         `json:"total_count" name:"total_count" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/image/describe_images.html
func (s *ImageService) DescribeImages(i *DescribeImagesInput) (*DescribeImagesOutput, error) {
	if i == nil {
		i = &DescribeImagesInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DescribeImages",
		RequestMethod: "GET",
	}

	x := &DescribeImagesOutput{}
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

type DescribeImagesInput struct {
	Images   []*string `json:"images" name:"images" location:"params"`
	Limit    *int      `json:"limit" name:"limit" default:"20" location:"params"`
	Offset   *int      `json:"offset" name:"offset" default:"0" location:"params"`
	OSFamily *string   `json:"os_family" name:"os_family" location:"params"`
	// ProcessorType's available values: 64bit, 32bit
	ProcessorType *string `json:"processor_type" name:"processor_type" location:"params"`
	// Provider's available values: system, self
	Provider   *string   `json:"provider" name:"provider" location:"params"`
	SearchWord *string   `json:"search_word" name:"search_word" location:"params"`
	Status     []*string `json:"status" name:"status" location:"params"`
	// Verbose's available values: 0
	Verbose *int `json:"verbose" name:"verbose" default:"0" location:"params"`
	// Visibility's available values: public, private
	Visibility *string `json:"visibility" name:"visibility" location:"params"`
}

func (v *DescribeImagesInput) Validate() error {

	if v.ProcessorType != nil {
		processorTypeValidValues := []string{"64bit", "32bit"}
		processorTypeParameterValue := fmt.Sprint(*v.ProcessorType)

		processorTypeIsValid := false
		for _, value := range processorTypeValidValues {
			if value == processorTypeParameterValue {
				processorTypeIsValid = true
			}
		}

		if !processorTypeIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "ProcessorType",
				ParameterValue: processorTypeParameterValue,
				AllowedValues:  processorTypeValidValues,
			}
		}
	}

	if v.Provider != nil {
		providerValidValues := []string{"system", "self"}
		providerParameterValue := fmt.Sprint(*v.Provider)

		providerIsValid := false
		for _, value := range providerValidValues {
			if value == providerParameterValue {
				providerIsValid = true
			}
		}

		if !providerIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "Provider",
				ParameterValue: providerParameterValue,
				AllowedValues:  providerValidValues,
			}
		}
	}

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

	if v.Visibility != nil {
		visibilityValidValues := []string{"public", "private"}
		visibilityParameterValue := fmt.Sprint(*v.Visibility)

		visibilityIsValid := false
		for _, value := range visibilityValidValues {
			if value == visibilityParameterValue {
				visibilityIsValid = true
			}
		}

		if !visibilityIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "Visibility",
				ParameterValue: visibilityParameterValue,
				AllowedValues:  visibilityValidValues,
			}
		}
	}

	return nil
}

type DescribeImagesOutput struct {
	Message    *string  `json:"message" name:"message"`
	Action     *string  `json:"action" name:"action" location:"elements"`
	ImageSet   []*Image `json:"image_set" name:"image_set" location:"elements"`
	RetCode    *int     `json:"ret_code" name:"ret_code" location:"elements"`
	TotalCount *int     `json:"total_count" name:"total_count" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/image/grant-image-to-users.html
func (s *ImageService) GrantImageToUsers(i *GrantImageToUsersInput) (*GrantImageToUsersOutput, error) {
	if i == nil {
		i = &GrantImageToUsersInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "GrantImageToUsers",
		RequestMethod: "GET",
	}

	x := &GrantImageToUsersOutput{}
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

type GrantImageToUsersInput struct {
	Image *string   `json:"image" name:"image" location:"params"` // Required
	Users []*string `json:"users" name:"users" location:"params"` // Required
}

func (v *GrantImageToUsersInput) Validate() error {

	if v.Image == nil {
		return errors.ParameterRequiredError{
			ParameterName: "Image",
			ParentName:    "GrantImageToUsersInput",
		}
	}

	if len(v.Users) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "Users",
			ParentName:    "GrantImageToUsersInput",
		}
	}

	return nil
}

type GrantImageToUsersOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/image/modify_image_attributes.html
func (s *ImageService) ModifyImageAttributes(i *ModifyImageAttributesInput) (*ModifyImageAttributesOutput, error) {
	if i == nil {
		i = &ModifyImageAttributesInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "ModifyImageAttributes",
		RequestMethod: "GET",
	}

	x := &ModifyImageAttributesOutput{}
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

type ModifyImageAttributesInput struct {
	Description *string `json:"description" name:"description" location:"params"`
	Image       *string `json:"image" name:"image" location:"params"` // Required
	ImageName   *string `json:"image_name" name:"image_name" location:"params"`
}

func (v *ModifyImageAttributesInput) Validate() error {

	if v.Image == nil {
		return errors.ParameterRequiredError{
			ParameterName: "Image",
			ParentName:    "ModifyImageAttributesInput",
		}
	}

	return nil
}

type ModifyImageAttributesOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	ImageID *string `json:"image_id" name:"image_id" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/image/revoke-image-from-users.html
func (s *ImageService) RevokeImageFromUsers(i *RevokeImageFromUsersInput) (*RevokeImageFromUsersOutput, error) {
	if i == nil {
		i = &RevokeImageFromUsersInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "RevokeImageFromUsers",
		RequestMethod: "GET",
	}

	x := &RevokeImageFromUsersOutput{}
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

type RevokeImageFromUsersInput struct {
	Image *string   `json:"image" name:"image" location:"params"` // Required
	Users []*string `json:"users" name:"users" location:"params"` // Required
}

func (v *RevokeImageFromUsersInput) Validate() error {

	if v.Image == nil {
		return errors.ParameterRequiredError{
			ParameterName: "Image",
			ParentName:    "RevokeImageFromUsersInput",
		}
	}

	if len(v.Users) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "Users",
			ParentName:    "RevokeImageFromUsersInput",
		}
	}

	return nil
}

type RevokeImageFromUsersOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}
