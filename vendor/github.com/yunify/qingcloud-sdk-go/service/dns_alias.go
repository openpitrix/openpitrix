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

type DNSAliasService struct {
	Config     *config.Config
	Properties *DNSAliasServiceProperties
}

type DNSAliasServiceProperties struct {
	// QingCloud Zone ID
	Zone *string `json:"zone" name:"zone"` // Required
}

func (s *QingCloudService) DNSAlias(zone string) (*DNSAliasService, error) {
	properties := &DNSAliasServiceProperties{
		Zone: &zone,
	}

	return &DNSAliasService{Config: s.Config, Properties: properties}, nil
}

// Documentation URL: https://docs.qingcloud.com/api/dns_alias/associate_dns_alias.html
func (s *DNSAliasService) AssociateDNSAlias(i *AssociateDNSAliasInput) (*AssociateDNSAliasOutput, error) {
	if i == nil {
		i = &AssociateDNSAliasInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "AssociateDNSAlias",
		RequestMethod: "GET",
	}

	x := &AssociateDNSAliasOutput{}
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

type AssociateDNSAliasInput struct {
	Prefix   *string `json:"prefix" name:"prefix" location:"params"`     // Required
	Resource *string `json:"resource" name:"resource" location:"params"` // Required
}

func (v *AssociateDNSAliasInput) Validate() error {

	if v.Prefix == nil {
		return errors.ParameterRequiredError{
			ParameterName: "Prefix",
			ParentName:    "AssociateDNSAliasInput",
		}
	}

	if v.Resource == nil {
		return errors.ParameterRequiredError{
			ParameterName: "Resource",
			ParentName:    "AssociateDNSAliasInput",
		}
	}

	return nil
}

type AssociateDNSAliasOutput struct {
	Message    *string `json:"message" name:"message"`
	Action     *string `json:"action" name:"action" location:"elements"`
	DNSAliasID *string `json:"dns_alias_id" name:"dns_alias_id" location:"elements"`
	DomainName *string `json:"domain_name" name:"domain_name" location:"elements"`
	JobID      *string `json:"job_id" name:"job_id" location:"elements"`
	RetCode    *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/dns_alias/describe_dns_aliases.html
func (s *DNSAliasService) DescribeDNSAliases(i *DescribeDNSAliasesInput) (*DescribeDNSAliasesOutput, error) {
	if i == nil {
		i = &DescribeDNSAliasesInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DescribeDNSAliases",
		RequestMethod: "GET",
	}

	x := &DescribeDNSAliasesOutput{}
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

type DescribeDNSAliasesInput struct {
	DNSAliases []*string `json:"dns_aliases" name:"dns_aliases" location:"params"`
	Limit      *int      `json:"limit" name:"limit" default:"20" location:"params"`
	Offset     *int      `json:"offset" name:"offset" default:"0" location:"params"`
	ResourceID *string   `json:"resource_id" name:"resource_id" location:"params"`
	SearchWord *string   `json:"search_word" name:"search_word" location:"params"`
}

func (v *DescribeDNSAliasesInput) Validate() error {

	return nil
}

type DescribeDNSAliasesOutput struct {
	Message     *string     `json:"message" name:"message"`
	Action      *string     `json:"action" name:"action" location:"elements"`
	DNSAliasSet []*DNSAlias `json:"dns_alias_set" name:"dns_alias_set" location:"elements"`
	RetCode     *int        `json:"ret_code" name:"ret_code" location:"elements"`
	TotalCount  *int        `json:"total_count" name:"total_count" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/dns_alias/dissociate_dns_aliases.html
func (s *DNSAliasService) DissociateDNSAliases(i *DissociateDNSAliasesInput) (*DissociateDNSAliasesOutput, error) {
	if i == nil {
		i = &DissociateDNSAliasesInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DissociateDNSAliases",
		RequestMethod: "GET",
	}

	x := &DissociateDNSAliasesOutput{}
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

type DissociateDNSAliasesInput struct {
	DNSAliases []*string `json:"dns_aliases" name:"dns_aliases" location:"params"` // Required
}

func (v *DissociateDNSAliasesInput) Validate() error {

	if len(v.DNSAliases) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "DNSAliases",
			ParentName:    "DissociateDNSAliasesInput",
		}
	}

	return nil
}

type DissociateDNSAliasesOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	JobID   *string `json:"job_id" name:"job_id" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/dns_alias/get_dns_label.html
func (s *DNSAliasService) GetDNSLabel(i *GetDNSLabelInput) (*GetDNSLabelOutput, error) {
	if i == nil {
		i = &GetDNSLabelInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "GetDNSLabel",
		RequestMethod: "GET",
	}

	x := &GetDNSLabelOutput{}
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

type GetDNSLabelInput struct {
}

func (v *GetDNSLabelInput) Validate() error {

	return nil
}

type GetDNSLabelOutput struct {
	Message    *string `json:"message" name:"message"`
	Action     *string `json:"action" name:"action" location:"elements"`
	DNSLabel   *string `json:"dns_label" name:"dns_label" location:"elements"`
	DomainName *string `json:"domain_name" name:"domain_name" location:"elements"`
	RetCode    *int    `json:"ret_code" name:"ret_code" location:"elements"`
}
