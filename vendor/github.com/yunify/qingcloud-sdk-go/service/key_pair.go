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

type KeyPairService struct {
	Config     *config.Config
	Properties *KeyPairServiceProperties
}

type KeyPairServiceProperties struct {
	// QingCloud Zone ID
	Zone *string `json:"zone" name:"zone"` // Required
}

func (s *QingCloudService) KeyPair(zone string) (*KeyPairService, error) {
	properties := &KeyPairServiceProperties{
		Zone: &zone,
	}

	return &KeyPairService{Config: s.Config, Properties: properties}, nil
}

// Documentation URL: https://docs.qingcloud.com/api/keypair/attach_key_pairs.html
func (s *KeyPairService) AttachKeyPairs(i *AttachKeyPairsInput) (*AttachKeyPairsOutput, error) {
	if i == nil {
		i = &AttachKeyPairsInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "AttachKeyPairs",
		RequestMethod: "GET",
	}

	x := &AttachKeyPairsOutput{}
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

type AttachKeyPairsInput struct {
	Instances []*string `json:"instances" name:"instances" location:"params"` // Required
	KeyPairs  []*string `json:"keypairs" name:"keypairs" location:"params"`   // Required
}

func (v *AttachKeyPairsInput) Validate() error {

	if len(v.Instances) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "Instances",
			ParentName:    "AttachKeyPairsInput",
		}
	}

	if len(v.KeyPairs) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "KeyPairs",
			ParentName:    "AttachKeyPairsInput",
		}
	}

	return nil
}

type AttachKeyPairsOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	JobID   *string `json:"job_id" name:"job_id" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/keypair/create_key_pairs.html
func (s *KeyPairService) CreateKeyPair(i *CreateKeyPairInput) (*CreateKeyPairOutput, error) {
	if i == nil {
		i = &CreateKeyPairInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "CreateKeyPair",
		RequestMethod: "GET",
	}

	x := &CreateKeyPairOutput{}
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

type CreateKeyPairInput struct {

	// EncryptMethod's available values: ssh-rsa, ssh-dss
	EncryptMethod *string `json:"encrypt_method" name:"encrypt_method" default:"ssh-rsa" location:"params"`
	KeyPairName   *string `json:"keypair_name" name:"keypair_name" location:"params"`
	// Mode's available values: system, user
	Mode      *string `json:"mode" name:"mode" default:"system" location:"params"`
	PublicKey *string `json:"public_key" name:"public_key" location:"params"`
}

func (v *CreateKeyPairInput) Validate() error {

	if v.EncryptMethod != nil {
		encryptMethodValidValues := []string{"ssh-rsa", "ssh-dss"}
		encryptMethodParameterValue := fmt.Sprint(*v.EncryptMethod)

		encryptMethodIsValid := false
		for _, value := range encryptMethodValidValues {
			if value == encryptMethodParameterValue {
				encryptMethodIsValid = true
			}
		}

		if !encryptMethodIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "EncryptMethod",
				ParameterValue: encryptMethodParameterValue,
				AllowedValues:  encryptMethodValidValues,
			}
		}
	}

	if v.Mode != nil {
		modeValidValues := []string{"system", "user"}
		modeParameterValue := fmt.Sprint(*v.Mode)

		modeIsValid := false
		for _, value := range modeValidValues {
			if value == modeParameterValue {
				modeIsValid = true
			}
		}

		if !modeIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "Mode",
				ParameterValue: modeParameterValue,
				AllowedValues:  modeValidValues,
			}
		}
	}

	return nil
}

type CreateKeyPairOutput struct {
	Message    *string `json:"message" name:"message"`
	Action     *string `json:"action" name:"action" location:"elements"`
	KeyPairID  *string `json:"keypair_id" name:"keypair_id" location:"elements"`
	PrivateKey *string `json:"private_key" name:"private_key" location:"elements"`
	RetCode    *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/keypair/delete_key_pairs.html
func (s *KeyPairService) DeleteKeyPairs(i *DeleteKeyPairsInput) (*DeleteKeyPairsOutput, error) {
	if i == nil {
		i = &DeleteKeyPairsInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DeleteKeyPairs",
		RequestMethod: "GET",
	}

	x := &DeleteKeyPairsOutput{}
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

type DeleteKeyPairsInput struct {
	KeyPairs []*string `json:"keypairs" name:"keypairs" location:"params"` // Required
}

func (v *DeleteKeyPairsInput) Validate() error {

	if len(v.KeyPairs) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "KeyPairs",
			ParentName:    "DeleteKeyPairsInput",
		}
	}

	return nil
}

type DeleteKeyPairsOutput struct {
	Message  *string   `json:"message" name:"message"`
	Action   *string   `json:"action" name:"action" location:"elements"`
	KeyPairs []*string `json:"keypairs" name:"keypairs" location:"elements"`
	RetCode  *int      `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/keypair/describe_key_pairs.html
func (s *KeyPairService) DescribeKeyPairs(i *DescribeKeyPairsInput) (*DescribeKeyPairsOutput, error) {
	if i == nil {
		i = &DescribeKeyPairsInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DescribeKeyPairs",
		RequestMethod: "GET",
	}

	x := &DescribeKeyPairsOutput{}
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

type DescribeKeyPairsInput struct {

	// EncryptMethod's available values: ssh-rsa, ssh-dss
	EncryptMethod *string   `json:"encrypt_method" name:"encrypt_method" location:"params"`
	InstanceID    *string   `json:"instance_id" name:"instance_id" location:"params"`
	KeyPairs      []*string `json:"keypairs" name:"keypairs" location:"params"`
	Limit         *int      `json:"limit" name:"limit" default:"20" location:"params"`
	Offset        *int      `json:"offset" name:"offset" default:"0" location:"params"`
	SearchWord    *string   `json:"search_word" name:"search_word" location:"params"`
	Tags          []*string `json:"tags" name:"tags" location:"params"`
	Verbose       *int      `json:"verbose" name:"verbose" default:"0" location:"params"`
}

func (v *DescribeKeyPairsInput) Validate() error {

	if v.EncryptMethod != nil {
		encryptMethodValidValues := []string{"ssh-rsa", "ssh-dss"}
		encryptMethodParameterValue := fmt.Sprint(*v.EncryptMethod)

		encryptMethodIsValid := false
		for _, value := range encryptMethodValidValues {
			if value == encryptMethodParameterValue {
				encryptMethodIsValid = true
			}
		}

		if !encryptMethodIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "EncryptMethod",
				ParameterValue: encryptMethodParameterValue,
				AllowedValues:  encryptMethodValidValues,
			}
		}
	}

	return nil
}

type DescribeKeyPairsOutput struct {
	Message    *string    `json:"message" name:"message"`
	Action     *string    `json:"action" name:"action" location:"elements"`
	KeyPairSet []*KeyPair `json:"keypair_set" name:"keypair_set" location:"elements"`
	RetCode    *int       `json:"ret_code" name:"ret_code" location:"elements"`
	TotalCount *int       `json:"total_count" name:"total_count" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/keypair/detach_key_pairs.html
func (s *KeyPairService) DetachKeyPairs(i *DetachKeyPairsInput) (*DetachKeyPairsOutput, error) {
	if i == nil {
		i = &DetachKeyPairsInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DetachKeyPairs",
		RequestMethod: "GET",
	}

	x := &DetachKeyPairsOutput{}
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

type DetachKeyPairsInput struct {
	Instances []*string `json:"instances" name:"instances" location:"params"` // Required
	KeyPairs  []*string `json:"keypairs" name:"keypairs" location:"params"`   // Required
}

func (v *DetachKeyPairsInput) Validate() error {

	if len(v.Instances) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "Instances",
			ParentName:    "DetachKeyPairsInput",
		}
	}

	if len(v.KeyPairs) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "KeyPairs",
			ParentName:    "DetachKeyPairsInput",
		}
	}

	return nil
}

type DetachKeyPairsOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	JobID   *string `json:"job_id" name:"job_id" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/keypair/modify_key_pair_attributes.html
func (s *KeyPairService) ModifyKeyPairAttributes(i *ModifyKeyPairAttributesInput) (*ModifyKeyPairAttributesOutput, error) {
	if i == nil {
		i = &ModifyKeyPairAttributesInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "ModifyKeyPairAttributes",
		RequestMethod: "GET",
	}

	x := &ModifyKeyPairAttributesOutput{}
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

type ModifyKeyPairAttributesInput struct {
	Description *string `json:"description" name:"description" location:"params"`
	KeyPair     *string `json:"keypair" name:"keypair" location:"params"` // Required
	KeyPairName *string `json:"keypair_name" name:"keypair_name" location:"params"`
}

func (v *ModifyKeyPairAttributesInput) Validate() error {

	if v.KeyPair == nil {
		return errors.ParameterRequiredError{
			ParameterName: "KeyPair",
			ParentName:    "ModifyKeyPairAttributesInput",
		}
	}

	return nil
}

type ModifyKeyPairAttributesOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}
