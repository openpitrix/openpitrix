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

type LoadBalancerService struct {
	Config     *config.Config
	Properties *LoadBalancerServiceProperties
}

type LoadBalancerServiceProperties struct {
	// QingCloud Zone ID
	Zone *string `json:"zone" name:"zone"` // Required
}

func (s *QingCloudService) LoadBalancer(zone string) (*LoadBalancerService, error) {
	properties := &LoadBalancerServiceProperties{
		Zone: &zone,
	}

	return &LoadBalancerService{Config: s.Config, Properties: properties}, nil
}

// Documentation URL: https://docs.qingcloud.com/api/lb/add_loadbalancer_backends.html
func (s *LoadBalancerService) AddLoadBalancerBackends(i *AddLoadBalancerBackendsInput) (*AddLoadBalancerBackendsOutput, error) {
	if i == nil {
		i = &AddLoadBalancerBackendsInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "AddLoadBalancerBackends",
		RequestMethod: "GET",
	}

	x := &AddLoadBalancerBackendsOutput{}
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

type AddLoadBalancerBackendsInput struct {
	Backends             []*LoadBalancerBackend `json:"backends" name:"backends" location:"params"`                           // Required
	LoadBalancerListener *string                `json:"loadbalancer_listener" name:"loadbalancer_listener" location:"params"` // Required
}

func (v *AddLoadBalancerBackendsInput) Validate() error {

	if len(v.Backends) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "Backends",
			ParentName:    "AddLoadBalancerBackendsInput",
		}
	}

	if len(v.Backends) > 0 {
		for _, property := range v.Backends {
			if err := property.Validate(); err != nil {
				return err
			}
		}
	}

	if v.LoadBalancerListener == nil {
		return errors.ParameterRequiredError{
			ParameterName: "LoadBalancerListener",
			ParentName:    "AddLoadBalancerBackendsInput",
		}
	}

	return nil
}

type AddLoadBalancerBackendsOutput struct {
	Message              *string   `json:"message" name:"message"`
	Action               *string   `json:"action" name:"action" location:"elements"`
	LoadBalancerBackends []*string `json:"loadbalancer_backends" name:"loadbalancer_backends" location:"elements"`
	RetCode              *int      `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/lb/add_loadbalancer_listeners.html
func (s *LoadBalancerService) AddLoadBalancerListeners(i *AddLoadBalancerListenersInput) (*AddLoadBalancerListenersOutput, error) {
	if i == nil {
		i = &AddLoadBalancerListenersInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "AddLoadBalancerListeners",
		RequestMethod: "GET",
	}

	x := &AddLoadBalancerListenersOutput{}
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

type AddLoadBalancerListenersInput struct {
	Listeners    []*LoadBalancerListener `json:"listeners" name:"listeners" location:"params"`
	LoadBalancer *string                 `json:"loadbalancer" name:"loadbalancer" location:"params"`
}

func (v *AddLoadBalancerListenersInput) Validate() error {

	if len(v.Listeners) > 0 {
		for _, property := range v.Listeners {
			if err := property.Validate(); err != nil {
				return err
			}
		}
	}

	return nil
}

type AddLoadBalancerListenersOutput struct {
	Message               *string   `json:"message" name:"message"`
	Action                *string   `json:"action" name:"action" location:"elements"`
	LoadBalancerListeners []*string `json:"loadbalancer_listeners" name:"loadbalancer_listeners" location:"elements"`
	RetCode               *int      `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/lb/add_loadbalancer_policy_rules.html
func (s *LoadBalancerService) AddLoadBalancerPolicyRules(i *AddLoadBalancerPolicyRulesInput) (*AddLoadBalancerPolicyRulesOutput, error) {
	if i == nil {
		i = &AddLoadBalancerPolicyRulesInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "AddLoadBalancerPolicyRules",
		RequestMethod: "GET",
	}

	x := &AddLoadBalancerPolicyRulesOutput{}
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

type AddLoadBalancerPolicyRulesInput struct {
	LoadBalancerPolicy *string                   `json:"loadbalancer_policy" name:"loadbalancer_policy" location:"params"`
	Rules              []*LoadBalancerPolicyRule `json:"rules" name:"rules" location:"params"`
}

func (v *AddLoadBalancerPolicyRulesInput) Validate() error {

	if len(v.Rules) > 0 {
		for _, property := range v.Rules {
			if err := property.Validate(); err != nil {
				return err
			}
		}
	}

	return nil
}

type AddLoadBalancerPolicyRulesOutput struct {
	Message                 *string   `json:"message" name:"message"`
	Action                  *string   `json:"action" name:"action" location:"elements"`
	LoadBalancerPolicyRules []*string `json:"loadbalancer_policy_rules" name:"loadbalancer_policy_rules" location:"elements"`
	RetCode                 *int      `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/lb/apply_loadbalancer_policy.html
func (s *LoadBalancerService) ApplyLoadBalancerPolicy(i *ApplyLoadBalancerPolicyInput) (*ApplyLoadBalancerPolicyOutput, error) {
	if i == nil {
		i = &ApplyLoadBalancerPolicyInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "ApplyLoadBalancerPolicy",
		RequestMethod: "GET",
	}

	x := &ApplyLoadBalancerPolicyOutput{}
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

type ApplyLoadBalancerPolicyInput struct {
	LoadBalancerPolicy *string `json:"loadbalancer_policy" name:"loadbalancer_policy" location:"params"` // Required
}

func (v *ApplyLoadBalancerPolicyInput) Validate() error {

	if v.LoadBalancerPolicy == nil {
		return errors.ParameterRequiredError{
			ParameterName: "LoadBalancerPolicy",
			ParentName:    "ApplyLoadBalancerPolicyInput",
		}
	}

	return nil
}

type ApplyLoadBalancerPolicyOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	JobID   *string `json:"job_id" name:"job_id" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/lb/associate_eips_to_loadbalancer.html
func (s *LoadBalancerService) AssociateEIPsToLoadBalancer(i *AssociateEIPsToLoadBalancerInput) (*AssociateEIPsToLoadBalancerOutput, error) {
	if i == nil {
		i = &AssociateEIPsToLoadBalancerInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "AssociateEipsToLoadBalancer",
		RequestMethod: "GET",
	}

	x := &AssociateEIPsToLoadBalancerOutput{}
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

type AssociateEIPsToLoadBalancerInput struct {
	EIPs         []*string `json:"eips" name:"eips" location:"params"`                 // Required
	LoadBalancer *string   `json:"loadbalancer" name:"loadbalancer" location:"params"` // Required
}

func (v *AssociateEIPsToLoadBalancerInput) Validate() error {

	if len(v.EIPs) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "EIPs",
			ParentName:    "AssociateEIPsToLoadBalancerInput",
		}
	}

	if v.LoadBalancer == nil {
		return errors.ParameterRequiredError{
			ParameterName: "LoadBalancer",
			ParentName:    "AssociateEIPsToLoadBalancerInput",
		}
	}

	return nil
}

type AssociateEIPsToLoadBalancerOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	JobID   *string `json:"job_id" name:"job_id" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/lb/create_loadbalancer.html
func (s *LoadBalancerService) CreateLoadBalancer(i *CreateLoadBalancerInput) (*CreateLoadBalancerOutput, error) {
	if i == nil {
		i = &CreateLoadBalancerInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "CreateLoadBalancer",
		RequestMethod: "GET",
	}

	x := &CreateLoadBalancerOutput{}
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

type CreateLoadBalancerInput struct {
	EIPs             []*string `json:"eips" name:"eips" location:"params"`
	HTTPHeaderSize   *int      `json:"http_header_size" name:"http_header_size" location:"params"`
	LoadBalancerName *string   `json:"loadbalancer_name" name:"loadbalancer_name" location:"params"`
	// LoadBalancerType's available values: 0, 1, 2, 3, 4, 5
	LoadBalancerType *int    `json:"loadbalancer_type" name:"loadbalancer_type" default:"0" location:"params"`
	NodeCount        *int    `json:"node_count" name:"node_count" location:"params"`
	PrivateIP        *string `json:"private_ip" name:"private_ip" location:"params"`
	SecurityGroup    *string `json:"security_group" name:"security_group" location:"params"`
	VxNet            *string `json:"vxnet" name:"vxnet" location:"params"`
}

func (v *CreateLoadBalancerInput) Validate() error {

	if v.LoadBalancerType != nil {
		loadBalancerTypeValidValues := []string{"0", "1", "2", "3", "4", "5"}
		loadBalancerTypeParameterValue := fmt.Sprint(*v.LoadBalancerType)

		loadBalancerTypeIsValid := false
		for _, value := range loadBalancerTypeValidValues {
			if value == loadBalancerTypeParameterValue {
				loadBalancerTypeIsValid = true
			}
		}

		if !loadBalancerTypeIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "LoadBalancerType",
				ParameterValue: loadBalancerTypeParameterValue,
				AllowedValues:  loadBalancerTypeValidValues,
			}
		}
	}

	return nil
}

type CreateLoadBalancerOutput struct {
	Message        *string `json:"message" name:"message"`
	Action         *string `json:"action" name:"action" location:"elements"`
	JobID          *string `json:"job_id" name:"job_id" location:"elements"`
	LoadBalancerID *string `json:"loadbalancer_id" name:"loadbalancer_id" location:"elements"`
	RetCode        *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/lb/create_loadbalancer_policy.html
func (s *LoadBalancerService) CreateLoadBalancerPolicy(i *CreateLoadBalancerPolicyInput) (*CreateLoadBalancerPolicyOutput, error) {
	if i == nil {
		i = &CreateLoadBalancerPolicyInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "CreateLoadBalancerPolicy",
		RequestMethod: "GET",
	}

	x := &CreateLoadBalancerPolicyOutput{}
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

type CreateLoadBalancerPolicyInput struct {
	LoadBalancerPolicyName *string `json:"loadbalancer_policy_name" name:"loadbalancer_policy_name" location:"params"` // Required
	// Operator's available values: or, and
	Operator *string `json:"operator" name:"operator" default:"or" location:"params"`
}

func (v *CreateLoadBalancerPolicyInput) Validate() error {

	if v.LoadBalancerPolicyName == nil {
		return errors.ParameterRequiredError{
			ParameterName: "LoadBalancerPolicyName",
			ParentName:    "CreateLoadBalancerPolicyInput",
		}
	}

	if v.Operator != nil {
		operatorValidValues := []string{"or", "and"}
		operatorParameterValue := fmt.Sprint(*v.Operator)

		operatorIsValid := false
		for _, value := range operatorValidValues {
			if value == operatorParameterValue {
				operatorIsValid = true
			}
		}

		if !operatorIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "Operator",
				ParameterValue: operatorParameterValue,
				AllowedValues:  operatorValidValues,
			}
		}
	}

	return nil
}

type CreateLoadBalancerPolicyOutput struct {
	Message              *string `json:"message" name:"message"`
	Action               *string `json:"action" name:"action" location:"elements"`
	LoadBalancerPolicyID *string `json:"loadbalancer_policy_id" name:"loadbalancer_policy_id" location:"elements"`
	RetCode              *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/lb/create_server_certificate.html
func (s *LoadBalancerService) CreateServerCertificate(i *CreateServerCertificateInput) (*CreateServerCertificateOutput, error) {
	if i == nil {
		i = &CreateServerCertificateInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "CreateServerCertificate",
		RequestMethod: "POST",
	}

	x := &CreateServerCertificateOutput{}
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

type CreateServerCertificateInput struct {
	CertificateContent    *string `json:"certificate_content" name:"certificate_content" location:"params"` // Required
	PrivateKey            *string `json:"private_key" name:"private_key" location:"params"`                 // Required
	ServerCertificateName *string `json:"server_certificate_name" name:"server_certificate_name" location:"params"`
}

func (v *CreateServerCertificateInput) Validate() error {

	if v.CertificateContent == nil {
		return errors.ParameterRequiredError{
			ParameterName: "CertificateContent",
			ParentName:    "CreateServerCertificateInput",
		}
	}

	if v.PrivateKey == nil {
		return errors.ParameterRequiredError{
			ParameterName: "PrivateKey",
			ParentName:    "CreateServerCertificateInput",
		}
	}

	return nil
}

type CreateServerCertificateOutput struct {
	Message             *string `json:"message" name:"message"`
	Action              *string `json:"action" name:"action" location:"elements"`
	RetCode             *int    `json:"ret_code" name:"ret_code" location:"elements"`
	ServerCertificateID *string `json:"server_certificate_id" name:"server_certificate_id" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/lb/delete_loadbalancer_backends.html
func (s *LoadBalancerService) DeleteLoadBalancerBackends(i *DeleteLoadBalancerBackendsInput) (*DeleteLoadBalancerBackendsOutput, error) {
	if i == nil {
		i = &DeleteLoadBalancerBackendsInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DeleteLoadBalancerBackends",
		RequestMethod: "GET",
	}

	x := &DeleteLoadBalancerBackendsOutput{}
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

type DeleteLoadBalancerBackendsInput struct {
	LoadBalancerBackends []*string `json:"loadbalancer_backends" name:"loadbalancer_backends" location:"params"` // Required
}

func (v *DeleteLoadBalancerBackendsInput) Validate() error {

	if len(v.LoadBalancerBackends) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "LoadBalancerBackends",
			ParentName:    "DeleteLoadBalancerBackendsInput",
		}
	}

	return nil
}

type DeleteLoadBalancerBackendsOutput struct {
	Message              *string   `json:"message" name:"message"`
	Action               *string   `json:"action" name:"action" location:"elements"`
	LoadBalancerBackends []*string `json:"loadbalancer_backends" name:"loadbalancer_backends" location:"elements"`
	RetCode              *int      `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/lb/delete_loadbalancer_listeners.html
func (s *LoadBalancerService) DeleteLoadBalancerListeners(i *DeleteLoadBalancerListenersInput) (*DeleteLoadBalancerListenersOutput, error) {
	if i == nil {
		i = &DeleteLoadBalancerListenersInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DeleteLoadBalancerListeners",
		RequestMethod: "GET",
	}

	x := &DeleteLoadBalancerListenersOutput{}
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

type DeleteLoadBalancerListenersInput struct {
	LoadBalancerListeners []*string `json:"loadbalancer_listeners" name:"loadbalancer_listeners" location:"params"` // Required
}

func (v *DeleteLoadBalancerListenersInput) Validate() error {

	if len(v.LoadBalancerListeners) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "LoadBalancerListeners",
			ParentName:    "DeleteLoadBalancerListenersInput",
		}
	}

	return nil
}

type DeleteLoadBalancerListenersOutput struct {
	Message               *string   `json:"message" name:"message"`
	Action                *string   `json:"action" name:"action" location:"elements"`
	LoadBalancerListeners []*string `json:"loadbalancer_listeners" name:"loadbalancer_listeners" location:"elements"`
	RetCode               *int      `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/lb/delete_loadbalancer_policies.html
func (s *LoadBalancerService) DeleteLoadBalancerPolicies(i *DeleteLoadBalancerPoliciesInput) (*DeleteLoadBalancerPoliciesOutput, error) {
	if i == nil {
		i = &DeleteLoadBalancerPoliciesInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DeleteLoadBalancerPolicies",
		RequestMethod: "GET",
	}

	x := &DeleteLoadBalancerPoliciesOutput{}
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

type DeleteLoadBalancerPoliciesInput struct {
	LoadBalancerPolicies []*string `json:"loadbalancer_policies" name:"loadbalancer_policies" location:"params"` // Required
}

func (v *DeleteLoadBalancerPoliciesInput) Validate() error {

	if len(v.LoadBalancerPolicies) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "LoadBalancerPolicies",
			ParentName:    "DeleteLoadBalancerPoliciesInput",
		}
	}

	return nil
}

type DeleteLoadBalancerPoliciesOutput struct {
	Message              *string   `json:"message" name:"message"`
	Action               *string   `json:"action" name:"action" location:"elements"`
	LoadBalancerPolicies []*string `json:"loadbalancer_policies" name:"loadbalancer_policies" location:"elements"`
	RetCode              *int      `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/lb/delete_loadbalancer_policy_rules.html
func (s *LoadBalancerService) DeleteLoadBalancerPolicyRules(i *DeleteLoadBalancerPolicyRulesInput) (*DeleteLoadBalancerPolicyRulesOutput, error) {
	if i == nil {
		i = &DeleteLoadBalancerPolicyRulesInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DeleteLoadBalancerPolicyRules",
		RequestMethod: "GET",
	}

	x := &DeleteLoadBalancerPolicyRulesOutput{}
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

type DeleteLoadBalancerPolicyRulesInput struct {
	LoadBalancerPolicyRules []*string `json:"loadbalancer_policy_rules" name:"loadbalancer_policy_rules" location:"params"` // Required
}

func (v *DeleteLoadBalancerPolicyRulesInput) Validate() error {

	if len(v.LoadBalancerPolicyRules) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "LoadBalancerPolicyRules",
			ParentName:    "DeleteLoadBalancerPolicyRulesInput",
		}
	}

	return nil
}

type DeleteLoadBalancerPolicyRulesOutput struct {
	Message                 *string   `json:"message" name:"message"`
	Action                  *string   `json:"action" name:"action" location:"elements"`
	LoadBalancerPolicyRules []*string `json:"loadbalancer_policy_rules" name:"loadbalancer_policy_rules" location:"elements"`
	RetCode                 *int      `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/lb/delete_loadbalancers.html
func (s *LoadBalancerService) DeleteLoadBalancers(i *DeleteLoadBalancersInput) (*DeleteLoadBalancersOutput, error) {
	if i == nil {
		i = &DeleteLoadBalancersInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DeleteLoadBalancers",
		RequestMethod: "GET",
	}

	x := &DeleteLoadBalancersOutput{}
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

type DeleteLoadBalancersInput struct {
	LoadBalancers []*string `json:"loadbalancers" name:"loadbalancers" location:"params"` // Required
}

func (v *DeleteLoadBalancersInput) Validate() error {

	if len(v.LoadBalancers) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "LoadBalancers",
			ParentName:    "DeleteLoadBalancersInput",
		}
	}

	return nil
}

type DeleteLoadBalancersOutput struct {
	Message       *string   `json:"message" name:"message"`
	Action        *string   `json:"action" name:"action" location:"elements"`
	JobID         *string   `json:"job_id" name:"job_id" location:"elements"`
	LoadBalancers []*string `json:"loadbalancers" name:"loadbalancers" location:"elements"`
	RetCode       *int      `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/lb/delete_server_certificates.html
func (s *LoadBalancerService) DeleteServerCertificates(i *DeleteServerCertificatesInput) (*DeleteServerCertificatesOutput, error) {
	if i == nil {
		i = &DeleteServerCertificatesInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DeleteServerCertificates",
		RequestMethod: "GET",
	}

	x := &DeleteServerCertificatesOutput{}
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

type DeleteServerCertificatesInput struct {
	ServerCertificates []*string `json:"server_certificates" name:"server_certificates" location:"params"` // Required
}

func (v *DeleteServerCertificatesInput) Validate() error {

	if len(v.ServerCertificates) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "ServerCertificates",
			ParentName:    "DeleteServerCertificatesInput",
		}
	}

	return nil
}

type DeleteServerCertificatesOutput struct {
	Message            *string   `json:"message" name:"message"`
	Action             *string   `json:"action" name:"action" location:"elements"`
	RetCode            *int      `json:"ret_code" name:"ret_code" location:"elements"`
	ServerCertificates []*string `json:"server_certificates" name:"server_certificates" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/lb/describe_loadbalancer_backends.html
func (s *LoadBalancerService) DescribeLoadBalancerBackends(i *DescribeLoadBalancerBackendsInput) (*DescribeLoadBalancerBackendsOutput, error) {
	if i == nil {
		i = &DescribeLoadBalancerBackendsInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DescribeLoadBalancerBackends",
		RequestMethod: "GET",
	}

	x := &DescribeLoadBalancerBackendsOutput{}
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

type DescribeLoadBalancerBackendsInput struct {
	Limit                *int    `json:"limit" name:"limit" default:"20" location:"params"`
	LoadBalancer         *string `json:"loadbalancer" name:"loadbalancer" location:"params"`
	LoadBalancerBackends *string `json:"loadbalancer_backends" name:"loadbalancer_backends" location:"params"`
	LoadBalancerListener *string `json:"loadbalancer_listener" name:"loadbalancer_listener" location:"params"`
	Offset               *int    `json:"offset" name:"offset" default:"0" location:"params"`
	Verbose              *int    `json:"verbose" name:"verbose" location:"params"`
}

func (v *DescribeLoadBalancerBackendsInput) Validate() error {

	return nil
}

type DescribeLoadBalancerBackendsOutput struct {
	Message                *string                `json:"message" name:"message"`
	Action                 *string                `json:"action" name:"action" location:"elements"`
	LoadBalancerBackendSet []*LoadBalancerBackend `json:"loadbalancer_backend_set" name:"loadbalancer_backend_set" location:"elements"`
	RetCode                *int                   `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/lb/describe_loadbalancer_listeners.html
func (s *LoadBalancerService) DescribeLoadBalancerListeners(i *DescribeLoadBalancerListenersInput) (*DescribeLoadBalancerListenersOutput, error) {
	if i == nil {
		i = &DescribeLoadBalancerListenersInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DescribeLoadBalancerListeners",
		RequestMethod: "GET",
	}

	x := &DescribeLoadBalancerListenersOutput{}
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

type DescribeLoadBalancerListenersInput struct {
	Limit                 *int      `json:"limit" name:"limit" default:"20" location:"params"`
	LoadBalancer          *string   `json:"loadbalancer" name:"loadbalancer" location:"params"`
	LoadBalancerListeners []*string `json:"loadbalancer_listeners" name:"loadbalancer_listeners" location:"params"`
	Offset                *int      `json:"offset" name:"offset" default:"0" location:"params"`
	Verbose               *int      `json:"verbose" name:"verbose" location:"params"`
}

func (v *DescribeLoadBalancerListenersInput) Validate() error {

	return nil
}

type DescribeLoadBalancerListenersOutput struct {
	Message                 *string                 `json:"message" name:"message"`
	Action                  *string                 `json:"action" name:"action" location:"elements"`
	LoadBalancerListenerSet []*LoadBalancerListener `json:"loadbalancer_listener_set" name:"loadbalancer_listener_set" location:"elements"`
	RetCode                 *int                    `json:"ret_code" name:"ret_code" location:"elements"`
	TotalCount              *int                    `json:"total_count" name:"total_count" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/lb/describe_loadbalancer_policies.html
func (s *LoadBalancerService) DescribeLoadBalancerPolicies(i *DescribeLoadBalancerPoliciesInput) (*DescribeLoadBalancerPoliciesOutput, error) {
	if i == nil {
		i = &DescribeLoadBalancerPoliciesInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DescribeLoadBalancerPolicies",
		RequestMethod: "GET",
	}

	x := &DescribeLoadBalancerPoliciesOutput{}
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

type DescribeLoadBalancerPoliciesInput struct {
	Limit                *int      `json:"limit" name:"limit" default:"20" location:"params"`
	LoadBalancerPolicies []*string `json:"loadbalancer_policies" name:"loadbalancer_policies" location:"params"`
	Offset               *int      `json:"offset" name:"offset" default:"0" location:"params"`
	Verbose              *int      `json:"verbose" name:"verbose" location:"params"`
}

func (v *DescribeLoadBalancerPoliciesInput) Validate() error {

	return nil
}

type DescribeLoadBalancerPoliciesOutput struct {
	Message               *string               `json:"message" name:"message"`
	Action                *string               `json:"action" name:"action" location:"elements"`
	LoadBalancerPolicySet []*LoadBalancerPolicy `json:"loadbalancer_policy_set" name:"loadbalancer_policy_set" location:"elements"`
	RetCode               *int                  `json:"ret_code" name:"ret_code" location:"elements"`
	TotalCount            *int                  `json:"total_count" name:"total_count" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/lb/describe_loadbalancer_policy_rules.html
func (s *LoadBalancerService) DescribeLoadBalancerPolicyRules(i *DescribeLoadBalancerPolicyRulesInput) (*DescribeLoadBalancerPolicyRulesOutput, error) {
	if i == nil {
		i = &DescribeLoadBalancerPolicyRulesInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DescribeLoadBalancerPolicyRules",
		RequestMethod: "GET",
	}

	x := &DescribeLoadBalancerPolicyRulesOutput{}
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

type DescribeLoadBalancerPolicyRulesInput struct {
	Limit                   *int      `json:"limit" name:"limit" default:"20" location:"params"`
	LoadBalancerPolicy      *string   `json:"loadbalancer_policy" name:"loadbalancer_policy" location:"params"`
	LoadBalancerPolicyRules []*string `json:"loadbalancer_policy_rules" name:"loadbalancer_policy_rules" location:"params"`
	Offset                  *int      `json:"offset" name:"offset" default:"0" location:"params"`
}

func (v *DescribeLoadBalancerPolicyRulesInput) Validate() error {

	return nil
}

type DescribeLoadBalancerPolicyRulesOutput struct {
	Message                   *string                   `json:"message" name:"message"`
	Action                    *string                   `json:"action" name:"action" location:"elements"`
	LoadBalancerPolicyRuleSet []*LoadBalancerPolicyRule `json:"loadbalancer_policy_rule_set" name:"loadbalancer_policy_rule_set" location:"elements"`
	RetCode                   *int                      `json:"ret_code" name:"ret_code" location:"elements"`
	TotalCount                *int                      `json:"total_count" name:"total_count" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/lb/describe_loadbalancers.html
func (s *LoadBalancerService) DescribeLoadBalancers(i *DescribeLoadBalancersInput) (*DescribeLoadBalancersOutput, error) {
	if i == nil {
		i = &DescribeLoadBalancersInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DescribeLoadBalancers",
		RequestMethod: "GET",
	}

	x := &DescribeLoadBalancersOutput{}
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

type DescribeLoadBalancersInput struct {
	Limit         *int      `json:"limit" name:"limit" default:"20" location:"params"`
	LoadBalancers []*string `json:"loadbalancers" name:"loadbalancers" location:"params"`
	Offset        *int      `json:"offset" name:"offset" default:"0" location:"params"`
	SearchWord    *string   `json:"search_word" name:"search_word" location:"params"`
	Status        []*string `json:"status" name:"status" location:"params"`
	Tags          []*string `json:"tags" name:"tags" location:"params"`
	Verbose       *int      `json:"verbose" name:"verbose" default:"0" location:"params"`
}

func (v *DescribeLoadBalancersInput) Validate() error {

	return nil
}

type DescribeLoadBalancersOutput struct {
	Message         *string         `json:"message" name:"message"`
	Action          *string         `json:"action" name:"action" location:"elements"`
	LoadBalancerSet []*LoadBalancer `json:"loadbalancer_set" name:"loadbalancer_set" location:"elements"`
	RetCode         *int            `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/lb/describe_server_certificates.html
func (s *LoadBalancerService) DescribeServerCertificates(i *DescribeServerCertificatesInput) (*DescribeServerCertificatesOutput, error) {
	if i == nil {
		i = &DescribeServerCertificatesInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DescribeServerCertificates",
		RequestMethod: "GET",
	}

	x := &DescribeServerCertificatesOutput{}
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

type DescribeServerCertificatesInput struct {
	Limit              *int      `json:"limit" name:"limit" default:"20" location:"params"`
	Offset             *int      `json:"offset" name:"offset" default:"0" location:"params"`
	SearchWord         *string   `json:"search_word" name:"search_word" location:"params"`
	ServerCertificates []*string `json:"server_certificates" name:"server_certificates" location:"params"`
	Verbose            *int      `json:"verbose" name:"verbose" default:"0" location:"params"`
}

func (v *DescribeServerCertificatesInput) Validate() error {

	return nil
}

type DescribeServerCertificatesOutput struct {
	Message              *string              `json:"message" name:"message"`
	Action               *string              `json:"action" name:"action" location:"elements"`
	RetCode              *int                 `json:"ret_code" name:"ret_code" location:"elements"`
	ServerCertificateSet []*ServerCertificate `json:"server_certificate_set" name:"server_certificate_set" location:"elements"`
	TotalCount           *int                 `json:"total_count" name:"total_count" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/lb/dissociate_eips_from_loadbalancer.html
func (s *LoadBalancerService) DissociateEIPsFromLoadBalancer(i *DissociateEIPsFromLoadBalancerInput) (*DissociateEIPsFromLoadBalancerOutput, error) {
	if i == nil {
		i = &DissociateEIPsFromLoadBalancerInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DissociateEipsFromLoadBalancer",
		RequestMethod: "GET",
	}

	x := &DissociateEIPsFromLoadBalancerOutput{}
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

type DissociateEIPsFromLoadBalancerInput struct {
	EIPs         []*string `json:"eips" name:"eips" location:"params"`                 // Required
	LoadBalancer *string   `json:"loadbalancer" name:"loadbalancer" location:"params"` // Required
}

func (v *DissociateEIPsFromLoadBalancerInput) Validate() error {

	if len(v.EIPs) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "EIPs",
			ParentName:    "DissociateEIPsFromLoadBalancerInput",
		}
	}

	if v.LoadBalancer == nil {
		return errors.ParameterRequiredError{
			ParameterName: "LoadBalancer",
			ParentName:    "DissociateEIPsFromLoadBalancerInput",
		}
	}

	return nil
}

type DissociateEIPsFromLoadBalancerOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	JobID   *string `json:"job_id" name:"job_id" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/monitor/get_loadbalancer_monitor.html
func (s *LoadBalancerService) GetLoadBalancerMonitor(i *GetLoadBalancerMonitorInput) (*GetLoadBalancerMonitorOutput, error) {
	if i == nil {
		i = &GetLoadBalancerMonitorInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "GetLoadBalancerMonitor",
		RequestMethod: "GET",
	}

	x := &GetLoadBalancerMonitorOutput{}
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

type GetLoadBalancerMonitorInput struct {
	EndTime      *time.Time `json:"end_time" name:"end_time" format:"ISO 8601" location:"params"` // Required
	Meters       []*string  `json:"meters" name:"meters" location:"params"`                       // Required
	Resource     *string    `json:"resource" name:"resource" location:"params"`                   // Required
	ResourceType *string    `json:"resource_type" name:"resource_type" default:"loadbalancer" location:"params"`
	StartTime    *time.Time `json:"start_time" name:"start_time" format:"ISO 8601" location:"params"` // Required
	// Step's available values: 5m, 15m, 2h, 1d
	Step *string `json:"step" name:"step" location:"params"` // Required
}

func (v *GetLoadBalancerMonitorInput) Validate() error {

	if len(v.Meters) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "Meters",
			ParentName:    "GetLoadBalancerMonitorInput",
		}
	}

	if v.Resource == nil {
		return errors.ParameterRequiredError{
			ParameterName: "Resource",
			ParentName:    "GetLoadBalancerMonitorInput",
		}
	}

	if v.Step == nil {
		return errors.ParameterRequiredError{
			ParameterName: "Step",
			ParentName:    "GetLoadBalancerMonitorInput",
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

type GetLoadBalancerMonitorOutput struct {
	Message    *string  `json:"message" name:"message"`
	Action     *string  `json:"action" name:"action" location:"elements"`
	MeterSet   []*Meter `json:"meter_set" name:"meter_set" location:"elements"`
	ResourceID *string  `json:"resource_id" name:"resource_id" location:"elements"`
	RetCode    *int     `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/lb/modify_loadbalancer_attributes.html
func (s *LoadBalancerService) ModifyLoadBalancerAttributes(i *ModifyLoadBalancerAttributesInput) (*ModifyLoadBalancerAttributesOutput, error) {
	if i == nil {
		i = &ModifyLoadBalancerAttributesInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "ModifyLoadBalancerAttributes",
		RequestMethod: "GET",
	}

	x := &ModifyLoadBalancerAttributesOutput{}
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

type ModifyLoadBalancerAttributesInput struct {
	Description      *string `json:"description" name:"description" location:"params"`
	HTTPHeaderSize   *int    `json:"http_header_size" name:"http_header_size" location:"params"`
	LoadBalancer     *string `json:"loadbalancer" name:"loadbalancer" location:"params"` // Required
	LoadBalancerName *string `json:"loadbalancer_name" name:"loadbalancer_name" location:"params"`
	NodeCount        *int    `json:"node_count" name:"node_count" location:"params"`
	PrivateIP        *string `json:"private_ip" name:"private_ip" location:"params"`
	SecurityGroup    *string `json:"security_group" name:"security_group" location:"params"`
}

func (v *ModifyLoadBalancerAttributesInput) Validate() error {

	if v.LoadBalancer == nil {
		return errors.ParameterRequiredError{
			ParameterName: "LoadBalancer",
			ParentName:    "ModifyLoadBalancerAttributesInput",
		}
	}

	return nil
}

type ModifyLoadBalancerAttributesOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/lb/modify_loadbalancer_backend_attributes.html
func (s *LoadBalancerService) ModifyLoadBalancerBackendAttributes(i *ModifyLoadBalancerBackendAttributesInput) (*ModifyLoadBalancerBackendAttributesOutput, error) {
	if i == nil {
		i = &ModifyLoadBalancerBackendAttributesInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "ModifyLoadBalancerBackendAttributes",
		RequestMethod: "GET",
	}

	x := &ModifyLoadBalancerBackendAttributesOutput{}
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

type ModifyLoadBalancerBackendAttributesInput struct {

	// Disabled's available values: 0, 1
	Disabled                *int    `json:"disabled" name:"disabled" location:"params"`
	LoadBalancerBackend     *string `json:"loadbalancer_backend" name:"loadbalancer_backend" location:"params"`
	LoadBalancerBackendName *string `json:"loadbalancer_backend_name" name:"loadbalancer_backend_name" location:"params"`
	LoadBalancerPolicyID    *string `json:"loadbalancer_policy_id" name:"loadbalancer_policy_id" location:"params"`
	Port                    *string `json:"port" name:"port" location:"params"`
	Weight                  *string `json:"weight" name:"weight" location:"params"`
}

func (v *ModifyLoadBalancerBackendAttributesInput) Validate() error {

	if v.Disabled != nil {
		disabledValidValues := []string{"0", "1"}
		disabledParameterValue := fmt.Sprint(*v.Disabled)

		disabledIsValid := false
		for _, value := range disabledValidValues {
			if value == disabledParameterValue {
				disabledIsValid = true
			}
		}

		if !disabledIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "Disabled",
				ParameterValue: disabledParameterValue,
				AllowedValues:  disabledValidValues,
			}
		}
	}

	return nil
}

type ModifyLoadBalancerBackendAttributesOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/lb/modify_loadbalancer_listener_attributes.html
func (s *LoadBalancerService) ModifyLoadBalancerListenerAttributes(i *ModifyLoadBalancerListenerAttributesInput) (*ModifyLoadBalancerListenerAttributesOutput, error) {
	if i == nil {
		i = &ModifyLoadBalancerListenerAttributesInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "ModifyLoadBalancerListenerAttributes",
		RequestMethod: "GET",
	}

	x := &ModifyLoadBalancerListenerAttributesOutput{}
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

type ModifyLoadBalancerListenerAttributesInput struct {
	BalanceMode              *string `json:"balance_mode" name:"balance_mode" location:"params"`
	Forwardfor               *int    `json:"forwardfor" name:"forwardfor" location:"params"`
	HealthyCheckMethod       *string `json:"healthy_check_method" name:"healthy_check_method" location:"params"`
	HealthyCheckOption       *string `json:"healthy_check_option" name:"healthy_check_option" location:"params"`
	ListenerOption           *int    `json:"listener_option" name:"listener_option" location:"params"`
	LoadBalancerListener     *string `json:"loadbalancer_listener" name:"loadbalancer_listener" location:"params"` // Required
	LoadBalancerListenerName *string `json:"loadbalancer_listener_name" name:"loadbalancer_listener_name" location:"params"`
	ServerCertificateID      *string `json:"server_certificate_id" name:"server_certificate_id" location:"params"`
	SessionSticky            *string `json:"session_sticky" name:"session_sticky" location:"params"`
}

func (v *ModifyLoadBalancerListenerAttributesInput) Validate() error {

	if v.LoadBalancerListener == nil {
		return errors.ParameterRequiredError{
			ParameterName: "LoadBalancerListener",
			ParentName:    "ModifyLoadBalancerListenerAttributesInput",
		}
	}

	return nil
}

type ModifyLoadBalancerListenerAttributesOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/lb/modify_loadbalancer_policy_attributes.html
func (s *LoadBalancerService) ModifyLoadBalancerPolicyAttributes(i *ModifyLoadBalancerPolicyAttributesInput) (*ModifyLoadBalancerPolicyAttributesOutput, error) {
	if i == nil {
		i = &ModifyLoadBalancerPolicyAttributesInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "ModifyLoadBalancerPolicyAttributes",
		RequestMethod: "GET",
	}

	x := &ModifyLoadBalancerPolicyAttributesOutput{}
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

type ModifyLoadBalancerPolicyAttributesInput struct {
	LoadBalancerPolicy     *string `json:"loadbalancer_policy" name:"loadbalancer_policy" location:"params"` // Required
	LoadBalancerPolicyName *string `json:"loadbalancer_policy_name" name:"loadbalancer_policy_name" location:"params"`
	Operator               *string `json:"operator" name:"operator" location:"params"`
}

func (v *ModifyLoadBalancerPolicyAttributesInput) Validate() error {

	if v.LoadBalancerPolicy == nil {
		return errors.ParameterRequiredError{
			ParameterName: "LoadBalancerPolicy",
			ParentName:    "ModifyLoadBalancerPolicyAttributesInput",
		}
	}

	return nil
}

type ModifyLoadBalancerPolicyAttributesOutput struct {
	Message              *string `json:"message" name:"message"`
	Action               *string `json:"action" name:"action" location:"elements"`
	LoadBalancerPolicyID *string `json:"loadbalancer_policy_id" name:"loadbalancer_policy_id" location:"elements"`
	RetCode              *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/lb/modify_loadbalancer_policy_rule_attributes.html
func (s *LoadBalancerService) ModifyLoadBalancerPolicyRuleAttributes(i *ModifyLoadBalancerPolicyRuleAttributesInput) (*ModifyLoadBalancerPolicyRuleAttributesOutput, error) {
	if i == nil {
		i = &ModifyLoadBalancerPolicyRuleAttributesInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "ModifyLoadBalancerPolicyRuleAttributes",
		RequestMethod: "GET",
	}

	x := &ModifyLoadBalancerPolicyRuleAttributesOutput{}
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

type ModifyLoadBalancerPolicyRuleAttributesInput struct {
	LoadBalancerPolicyRule     *string `json:"loadbalancer_policy_rule" name:"loadbalancer_policy_rule" location:"params"` // Required
	LoadBalancerPolicyRuleName *string `json:"loadbalancer_policy_rule_name" name:"loadbalancer_policy_rule_name" location:"params"`
	Val                        *string `json:"val" name:"val" location:"params"`
}

func (v *ModifyLoadBalancerPolicyRuleAttributesInput) Validate() error {

	if v.LoadBalancerPolicyRule == nil {
		return errors.ParameterRequiredError{
			ParameterName: "LoadBalancerPolicyRule",
			ParentName:    "ModifyLoadBalancerPolicyRuleAttributesInput",
		}
	}

	return nil
}

type ModifyLoadBalancerPolicyRuleAttributesOutput struct {
	Message                  *string `json:"message" name:"message"`
	Action                   *string `json:"action" name:"action" location:"elements"`
	LoadBalancerPolicyRuleID *string `json:"loadbalancer_policy_rule_id" name:"loadbalancer_policy_rule_id" location:"elements"`
	RetCode                  *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/lb/modify_server_certificate_attributes.html
func (s *LoadBalancerService) ModifyServerCertificateAttributes(i *ModifyServerCertificateAttributesInput) (*ModifyServerCertificateAttributesOutput, error) {
	if i == nil {
		i = &ModifyServerCertificateAttributesInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "ModifyServerCertificateAttributes",
		RequestMethod: "GET",
	}

	x := &ModifyServerCertificateAttributesOutput{}
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

type ModifyServerCertificateAttributesInput struct {
	Description           *string `json:"description" name:"description" location:"params"`
	ServerCertificate     *string `json:"server_certificate" name:"server_certificate" location:"params"` // Required
	ServerCertificateName *string `json:"server_certificate_name" name:"server_certificate_name" location:"params"`
}

func (v *ModifyServerCertificateAttributesInput) Validate() error {

	if v.ServerCertificate == nil {
		return errors.ParameterRequiredError{
			ParameterName: "ServerCertificate",
			ParentName:    "ModifyServerCertificateAttributesInput",
		}
	}

	return nil
}

type ModifyServerCertificateAttributesOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/lb/resize_loadbalancers.html
func (s *LoadBalancerService) ResizeLoadBalancers(i *ResizeLoadBalancersInput) (*ResizeLoadBalancersOutput, error) {
	if i == nil {
		i = &ResizeLoadBalancersInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "ResizeLoadBalancers",
		RequestMethod: "GET",
	}

	x := &ResizeLoadBalancersOutput{}
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

type ResizeLoadBalancersInput struct {

	// LoadBalancerType's available values: 0, 1, 2, 3, 4, 5
	LoadBalancerType *int      `json:"loadbalancer_type" name:"loadbalancer_type" location:"params"`
	LoadBalancers    []*string `json:"loadbalancers" name:"loadbalancers" location:"params"`
}

func (v *ResizeLoadBalancersInput) Validate() error {

	if v.LoadBalancerType != nil {
		loadBalancerTypeValidValues := []string{"0", "1", "2", "3", "4", "5"}
		loadBalancerTypeParameterValue := fmt.Sprint(*v.LoadBalancerType)

		loadBalancerTypeIsValid := false
		for _, value := range loadBalancerTypeValidValues {
			if value == loadBalancerTypeParameterValue {
				loadBalancerTypeIsValid = true
			}
		}

		if !loadBalancerTypeIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "LoadBalancerType",
				ParameterValue: loadBalancerTypeParameterValue,
				AllowedValues:  loadBalancerTypeValidValues,
			}
		}
	}

	return nil
}

type ResizeLoadBalancersOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	JobID   *string `json:"job_id" name:"job_id" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/lb/start_loadbalancers.html
func (s *LoadBalancerService) StartLoadBalancers(i *StartLoadBalancersInput) (*StartLoadBalancersOutput, error) {
	if i == nil {
		i = &StartLoadBalancersInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "StartLoadBalancers",
		RequestMethod: "GET",
	}

	x := &StartLoadBalancersOutput{}
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

type StartLoadBalancersInput struct {
	LoadBalancers []*string `json:"loadbalancers" name:"loadbalancers" location:"params"` // Required
}

func (v *StartLoadBalancersInput) Validate() error {

	if len(v.LoadBalancers) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "LoadBalancers",
			ParentName:    "StartLoadBalancersInput",
		}
	}

	return nil
}

type StartLoadBalancersOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	JobID   *string `json:"job_id" name:"job_id" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/lb/stop_loadbalancers.html
func (s *LoadBalancerService) StopLoadBalancers(i *StopLoadBalancersInput) (*StopLoadBalancersOutput, error) {
	if i == nil {
		i = &StopLoadBalancersInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "StopLoadBalancers",
		RequestMethod: "GET",
	}

	x := &StopLoadBalancersOutput{}
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

type StopLoadBalancersInput struct {
	LoadBalancers []*string `json:"loadbalancers" name:"loadbalancers" location:"params"` // Required
}

func (v *StopLoadBalancersInput) Validate() error {

	if len(v.LoadBalancers) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "LoadBalancers",
			ParentName:    "StopLoadBalancersInput",
		}
	}

	return nil
}

type StopLoadBalancersOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	JobID   *string `json:"job_id" name:"job_id" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/lb/update_loadbalancers.html
func (s *LoadBalancerService) UpdateLoadBalancers(i *UpdateLoadBalancersInput) (*UpdateLoadBalancersOutput, error) {
	if i == nil {
		i = &UpdateLoadBalancersInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "UpdateLoadBalancers",
		RequestMethod: "GET",
	}

	x := &UpdateLoadBalancersOutput{}
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

type UpdateLoadBalancersInput struct {
	LoadBalancers []*string `json:"loadbalancers" name:"loadbalancers" location:"params"` // Required
}

func (v *UpdateLoadBalancersInput) Validate() error {

	if len(v.LoadBalancers) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "LoadBalancers",
			ParentName:    "UpdateLoadBalancersInput",
		}
	}

	return nil
}

type UpdateLoadBalancersOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	JobID   *string `json:"job_id" name:"job_id" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}
