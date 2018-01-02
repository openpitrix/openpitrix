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

type RouterService struct {
	Config     *config.Config
	Properties *RouterServiceProperties
}

type RouterServiceProperties struct {
	// QingCloud Zone ID
	Zone *string `json:"zone" name:"zone"` // Required
}

func (s *QingCloudService) Router(zone string) (*RouterService, error) {
	properties := &RouterServiceProperties{
		Zone: &zone,
	}

	return &RouterService{Config: s.Config, Properties: properties}, nil
}

// Documentation URL: https://docs.qingcloud.com/api/router/add_router_static_entries.html
func (s *RouterService) AddRouterStaticEntries(i *AddRouterStaticEntriesInput) (*AddRouterStaticEntriesOutput, error) {
	if i == nil {
		i = &AddRouterStaticEntriesInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "AddRouterStaticEntries",
		RequestMethod: "GET",
	}

	x := &AddRouterStaticEntriesOutput{}
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

type AddRouterStaticEntriesInput struct {
	Entries      []*RouterStaticEntry `json:"entries" name:"entries" location:"params"`
	RouterStatic *string              `json:"router_static" name:"router_static" location:"params"` // Required
}

func (v *AddRouterStaticEntriesInput) Validate() error {

	if len(v.Entries) > 0 {
		for _, property := range v.Entries {
			if err := property.Validate(); err != nil {
				return err
			}
		}
	}

	if v.RouterStatic == nil {
		return errors.ParameterRequiredError{
			ParameterName: "RouterStatic",
			ParentName:    "AddRouterStaticEntriesInput",
		}
	}

	return nil
}

type AddRouterStaticEntriesOutput struct {
	Message             *string   `json:"message" name:"message"`
	Action              *string   `json:"action" name:"action" location:"elements"`
	RetCode             *int      `json:"ret_code" name:"ret_code" location:"elements"`
	RouterStaticEntries []*string `json:"router_static_entries" name:"router_static_entries" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/router/add_router_statics.html
func (s *RouterService) AddRouterStatics(i *AddRouterStaticsInput) (*AddRouterStaticsOutput, error) {
	if i == nil {
		i = &AddRouterStaticsInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "AddRouterStatics",
		RequestMethod: "GET",
	}

	x := &AddRouterStaticsOutput{}
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

type AddRouterStaticsInput struct {
	Router  *string         `json:"router" name:"router" location:"params"`   // Required
	Statics []*RouterStatic `json:"statics" name:"statics" location:"params"` // Required
	VxNet   *string         `json:"vxnet" name:"vxnet" location:"params"`
}

func (v *AddRouterStaticsInput) Validate() error {

	if v.Router == nil {
		return errors.ParameterRequiredError{
			ParameterName: "Router",
			ParentName:    "AddRouterStaticsInput",
		}
	}

	if len(v.Statics) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "Statics",
			ParentName:    "AddRouterStaticsInput",
		}
	}

	if len(v.Statics) > 0 {
		for _, property := range v.Statics {
			if err := property.Validate(); err != nil {
				return err
			}
		}
	}

	return nil
}

type AddRouterStaticsOutput struct {
	Message       *string   `json:"message" name:"message"`
	Action        *string   `json:"action" name:"action" location:"elements"`
	RetCode       *int      `json:"ret_code" name:"ret_code" location:"elements"`
	RouterStatics []*string `json:"router_statics" name:"router_statics" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/router/create_routers.html
func (s *RouterService) CreateRouters(i *CreateRoutersInput) (*CreateRoutersOutput, error) {
	if i == nil {
		i = &CreateRoutersInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "CreateRouters",
		RequestMethod: "GET",
	}

	x := &CreateRoutersOutput{}
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

type CreateRoutersInput struct {
	Count      *int    `json:"count" name:"count" default:"1" location:"params"`
	RouterName *string `json:"router_name" name:"router_name" location:"params"`
	// RouterType's available values: 0, 1, 2, 3
	RouterType    *int    `json:"router_type" name:"router_type" default:"1" location:"params"`
	SecurityGroup *string `json:"security_group" name:"security_group" location:"params"`
	VpcNetwork    *string `json:"vpc_network" name:"vpc_network" location:"params"`
}

func (v *CreateRoutersInput) Validate() error {

	if v.RouterType != nil {
		routerTypeValidValues := []string{"0", "1", "2", "3"}
		routerTypeParameterValue := fmt.Sprint(*v.RouterType)

		routerTypeIsValid := false
		for _, value := range routerTypeValidValues {
			if value == routerTypeParameterValue {
				routerTypeIsValid = true
			}
		}

		if !routerTypeIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "RouterType",
				ParameterValue: routerTypeParameterValue,
				AllowedValues:  routerTypeValidValues,
			}
		}
	}

	return nil
}

type CreateRoutersOutput struct {
	Message *string   `json:"message" name:"message"`
	Action  *string   `json:"action" name:"action" location:"elements"`
	JobID   *string   `json:"job_id" name:"job_id" location:"elements"`
	RetCode *int      `json:"ret_code" name:"ret_code" location:"elements"`
	Routers []*string `json:"routers" name:"routers" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/router/delete_router_static_entries.html
func (s *RouterService) DeleteRouterStaticEntries(i *DeleteRouterStaticEntriesInput) (*DeleteRouterStaticEntriesOutput, error) {
	if i == nil {
		i = &DeleteRouterStaticEntriesInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DeleteRouterStaticEntries",
		RequestMethod: "GET",
	}

	x := &DeleteRouterStaticEntriesOutput{}
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

type DeleteRouterStaticEntriesInput struct {
	RouterStaticEntries []*string `json:"router_static_entries" name:"router_static_entries" location:"params"` // Required
}

func (v *DeleteRouterStaticEntriesInput) Validate() error {

	if len(v.RouterStaticEntries) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "RouterStaticEntries",
			ParentName:    "DeleteRouterStaticEntriesInput",
		}
	}

	return nil
}

type DeleteRouterStaticEntriesOutput struct {
	Message             *string   `json:"message" name:"message"`
	Action              *string   `json:"action" name:"action" location:"elements"`
	RetCode             *int      `json:"ret_code" name:"ret_code" location:"elements"`
	RouterStaticEntries []*string `json:"router_static_entries" name:"router_static_entries" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/router/delete_router_statics.html
func (s *RouterService) DeleteRouterStatics(i *DeleteRouterStaticsInput) (*DeleteRouterStaticsOutput, error) {
	if i == nil {
		i = &DeleteRouterStaticsInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DeleteRouterStatics",
		RequestMethod: "GET",
	}

	x := &DeleteRouterStaticsOutput{}
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

type DeleteRouterStaticsInput struct {
	RouterStatics []*string `json:"router_statics" name:"router_statics" location:"params"` // Required
}

func (v *DeleteRouterStaticsInput) Validate() error {

	if len(v.RouterStatics) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "RouterStatics",
			ParentName:    "DeleteRouterStaticsInput",
		}
	}

	return nil
}

type DeleteRouterStaticsOutput struct {
	Message       *string   `json:"message" name:"message"`
	Action        *string   `json:"action" name:"action" location:"elements"`
	RetCode       *int      `json:"ret_code" name:"ret_code" location:"elements"`
	RouterStatics []*string `json:"router_statics" name:"router_statics" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/router/delete_routers.html
func (s *RouterService) DeleteRouters(i *DeleteRoutersInput) (*DeleteRoutersOutput, error) {
	if i == nil {
		i = &DeleteRoutersInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DeleteRouters",
		RequestMethod: "GET",
	}

	x := &DeleteRoutersOutput{}
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

type DeleteRoutersInput struct {
	Routers []*string `json:"routers" name:"routers" location:"params"` // Required
}

func (v *DeleteRoutersInput) Validate() error {

	if len(v.Routers) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "Routers",
			ParentName:    "DeleteRoutersInput",
		}
	}

	return nil
}

type DeleteRoutersOutput struct {
	Message *string   `json:"message" name:"message"`
	Action  *string   `json:"action" name:"action" location:"elements"`
	JobID   *string   `json:"job_id" name:"job_id" location:"elements"`
	RetCode *int      `json:"ret_code" name:"ret_code" location:"elements"`
	Routers []*string `json:"routers" name:"routers" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/router/describe_router_static_entries.html
func (s *RouterService) DescribeRouterStaticEntries(i *DescribeRouterStaticEntriesInput) (*DescribeRouterStaticEntriesOutput, error) {
	if i == nil {
		i = &DescribeRouterStaticEntriesInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DescribeRouterStaticEntries",
		RequestMethod: "GET",
	}

	x := &DescribeRouterStaticEntriesOutput{}
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

type DescribeRouterStaticEntriesInput struct {
	Limit               *int    `json:"limit" name:"limit" location:"params"`
	Offset              *int    `json:"offset" name:"offset" location:"params"`
	RouterStatic        *string `json:"router_static" name:"router_static" location:"params"`
	RouterStaticEntryID *string `json:"router_static_entry_id" name:"router_static_entry_id" location:"params"`
}

func (v *DescribeRouterStaticEntriesInput) Validate() error {

	return nil
}

type DescribeRouterStaticEntriesOutput struct {
	Message              *string              `json:"message" name:"message"`
	Action               *string              `json:"action" name:"action" location:"elements"`
	RetCode              *int                 `json:"ret_code" name:"ret_code" location:"elements"`
	RouterStaticEntrySet []*RouterStaticEntry `json:"router_static_entry_set" name:"router_static_entry_set" location:"elements"`
	TotalCount           *int                 `json:"total_count" name:"total_count" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/router/describe_router_statics.html
func (s *RouterService) DescribeRouterStatics(i *DescribeRouterStaticsInput) (*DescribeRouterStaticsOutput, error) {
	if i == nil {
		i = &DescribeRouterStaticsInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DescribeRouterStatics",
		RequestMethod: "GET",
	}

	x := &DescribeRouterStaticsOutput{}
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

type DescribeRouterStaticsInput struct {
	Limit         *int      `json:"limit" name:"limit" default:"20" location:"params"`
	Offset        *int      `json:"offset" name:"offset" default:"0" location:"params"`
	Router        *string   `json:"router" name:"router" location:"params"` // Required
	RouterStatics []*string `json:"router_statics" name:"router_statics" location:"params"`
	// StaticType's available values: 1, 2, 3, 4, 5, 6, 7, 8
	StaticType *int `json:"static_type" name:"static_type" location:"params"`
	// Verbose's available values: 0, 1
	Verbose *int    `json:"verbose" name:"verbose" location:"params"`
	VxNet   *string `json:"vxnet" name:"vxnet" location:"params"`
}

func (v *DescribeRouterStaticsInput) Validate() error {

	if v.Router == nil {
		return errors.ParameterRequiredError{
			ParameterName: "Router",
			ParentName:    "DescribeRouterStaticsInput",
		}
	}

	if v.StaticType != nil {
		staticTypeValidValues := []string{"1", "2", "3", "4", "5", "6", "7", "8"}
		staticTypeParameterValue := fmt.Sprint(*v.StaticType)

		staticTypeIsValid := false
		for _, value := range staticTypeValidValues {
			if value == staticTypeParameterValue {
				staticTypeIsValid = true
			}
		}

		if !staticTypeIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "StaticType",
				ParameterValue: staticTypeParameterValue,
				AllowedValues:  staticTypeValidValues,
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

type DescribeRouterStaticsOutput struct {
	Message         *string         `json:"message" name:"message"`
	Action          *string         `json:"action" name:"action" location:"elements"`
	RetCode         *int            `json:"ret_code" name:"ret_code" location:"elements"`
	RouterStaticSet []*RouterStatic `json:"router_static_set" name:"router_static_set" location:"elements"`
	TotalCount      *int            `json:"total_count" name:"total_count" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/router/describe_router_vxnets.html
func (s *RouterService) DescribeRouterVxNets(i *DescribeRouterVxNetsInput) (*DescribeRouterVxNetsOutput, error) {
	if i == nil {
		i = &DescribeRouterVxNetsInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DescribeRouterVxnets",
		RequestMethod: "GET",
	}

	x := &DescribeRouterVxNetsOutput{}
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

type DescribeRouterVxNetsInput struct {
	Limit  *int    `json:"limit" name:"limit" default:"20" location:"params"`
	Offset *int    `json:"offset" name:"offset" default:"0" location:"params"`
	Router *string `json:"router" name:"router" location:"params"` // Required
	// Verbose's available values: 0, 1
	Verbose *int    `json:"verbose" name:"verbose" location:"params"`
	VxNet   *string `json:"vxnet" name:"vxnet" location:"params"`
}

func (v *DescribeRouterVxNetsInput) Validate() error {

	if v.Router == nil {
		return errors.ParameterRequiredError{
			ParameterName: "Router",
			ParentName:    "DescribeRouterVxNetsInput",
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

type DescribeRouterVxNetsOutput struct {
	Message        *string        `json:"message" name:"message"`
	Action         *string        `json:"action" name:"action" location:"elements"`
	RetCode        *int           `json:"ret_code" name:"ret_code" location:"elements"`
	RouterVxNetSet []*RouterVxNet `json:"router_vxnet_set" name:"router_vxnet_set" location:"elements"`
	TotalCount     *int           `json:"total_count" name:"total_count" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/router/describe_routers.html
func (s *RouterService) DescribeRouters(i *DescribeRoutersInput) (*DescribeRoutersOutput, error) {
	if i == nil {
		i = &DescribeRoutersInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "DescribeRouters",
		RequestMethod: "GET",
	}

	x := &DescribeRoutersOutput{}
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

type DescribeRoutersInput struct {
	Limit      *int      `json:"limit" name:"limit" location:"params"`
	Offset     *int      `json:"offset" name:"offset" location:"params"`
	Routers    []*string `json:"routers" name:"routers" location:"params"`
	SearchWord *string   `json:"search_word" name:"search_word" location:"params"`
	Status     []*string `json:"status" name:"status" location:"params"`
	Tags       []*string `json:"tags" name:"tags" location:"params"`
	// Verbose's available values: 0, 1
	Verbose *int    `json:"verbose" name:"verbose" location:"params"`
	VxNet   *string `json:"vxnet" name:"vxnet" location:"params"`
}

func (v *DescribeRoutersInput) Validate() error {

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

type DescribeRoutersOutput struct {
	Message    *string   `json:"message" name:"message"`
	Action     *string   `json:"action" name:"action" location:"elements"`
	RetCode    *int      `json:"ret_code" name:"ret_code" location:"elements"`
	RouterSet  []*Router `json:"router_set" name:"router_set" location:"elements"`
	TotalCount *int      `json:"total_count" name:"total_count" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/monitor/get_monitor.html
func (s *RouterService) GetRouterMonitor(i *GetRouterMonitorInput) (*GetRouterMonitorOutput, error) {
	if i == nil {
		i = &GetRouterMonitorInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "GetMonitor",
		RequestMethod: "GET",
	}

	x := &GetRouterMonitorOutput{}
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

type GetRouterMonitorInput struct {
	EndTime   *time.Time `json:"end_time" name:"end_time" format:"ISO 8601" location:"params"`     // Required
	Meters    []*string  `json:"meters" name:"meters" location:"params"`                           // Required
	Resource  *string    `json:"resource" name:"resource" location:"params"`                       // Required
	StartTime *time.Time `json:"start_time" name:"start_time" format:"ISO 8601" location:"params"` // Required
	// Step's available values: 5m, 15m, 2h, 1d
	Step *string `json:"step" name:"step" location:"params"` // Required
}

func (v *GetRouterMonitorInput) Validate() error {

	if len(v.Meters) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "Meters",
			ParentName:    "GetRouterMonitorInput",
		}
	}

	if v.Resource == nil {
		return errors.ParameterRequiredError{
			ParameterName: "Resource",
			ParentName:    "GetRouterMonitorInput",
		}
	}

	if v.Step == nil {
		return errors.ParameterRequiredError{
			ParameterName: "Step",
			ParentName:    "GetRouterMonitorInput",
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

type GetRouterMonitorOutput struct {
	Message    *string  `json:"message" name:"message"`
	Action     *string  `json:"action" name:"action" location:"elements"`
	MeterSet   []*Meter `json:"meter_set" name:"meter_set" location:"elements"`
	ResourceID *string  `json:"resource_id" name:"resource_id" location:"elements"`
	RetCode    *int     `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/router/get_vpn_certs.html
func (s *RouterService) GetVPNCerts(i *GetVPNCertsInput) (*GetVPNCertsOutput, error) {
	if i == nil {
		i = &GetVPNCertsInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "GetVPNCerts",
		RequestMethod: "GET",
	}

	x := &GetVPNCertsOutput{}
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

type GetVPNCertsInput struct {

	// Platform's available values: windows, linux, mac
	Platform *string `json:"platform" name:"platform" location:"params"`
	Router   *string `json:"router" name:"router" location:"params"` // Required
}

func (v *GetVPNCertsInput) Validate() error {

	if v.Platform != nil {
		platformValidValues := []string{"windows", "linux", "mac"}
		platformParameterValue := fmt.Sprint(*v.Platform)

		platformIsValid := false
		for _, value := range platformValidValues {
			if value == platformParameterValue {
				platformIsValid = true
			}
		}

		if !platformIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "Platform",
				ParameterValue: platformParameterValue,
				AllowedValues:  platformValidValues,
			}
		}
	}

	if v.Router == nil {
		return errors.ParameterRequiredError{
			ParameterName: "Router",
			ParentName:    "GetVPNCertsInput",
		}
	}

	return nil
}

type GetVPNCertsOutput struct {
	Message           *string `json:"message" name:"message"`
	Action            *string `json:"action" name:"action" location:"elements"`
	CaCert            *string `json:"ca_cert" name:"ca_cert" location:"elements"`
	ClientCrt         *string `json:"client_crt" name:"client_crt" location:"elements"`
	ClientKey         *string `json:"client_key" name:"client_key" location:"elements"`
	LinuxConfSample   *string `json:"linux_conf_sample" name:"linux_conf_sample" location:"elements"`
	MacConfSample     *string `json:"mac_conf_sample" name:"mac_conf_sample" location:"elements"`
	RetCode           *int    `json:"ret_code" name:"ret_code" location:"elements"`
	RouterID          *string `json:"router_id" name:"router_id" location:"elements"`
	StaticKey         *string `json:"static_key" name:"static_key" location:"elements"`
	WindowsConfSample *string `json:"windows_conf_sample" name:"windows_conf_sample" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/router/join_router.html
func (s *RouterService) JoinRouter(i *JoinRouterInput) (*JoinRouterOutput, error) {
	if i == nil {
		i = &JoinRouterInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "JoinRouter",
		RequestMethod: "GET",
	}

	x := &JoinRouterOutput{}
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

type JoinRouterInput struct {
	DYNIPEnd   *string `json:"dyn_ip_end" name:"dyn_ip_end" location:"params"`
	DYNIPStart *string `json:"dyn_ip_start" name:"dyn_ip_start" location:"params"`
	// Features's available values: 1
	Features  *int    `json:"features" name:"features" default:"1" location:"params"`
	IPNetwork *string `json:"ip_network" name:"ip_network" location:"params"` // Required
	ManagerIP *string `json:"manager_ip" name:"manager_ip" location:"params"`
	Router    *string `json:"router" name:"router" location:"params"` // Required
	VxNet     *string `json:"vxnet" name:"vxnet" location:"params"`   // Required
}

func (v *JoinRouterInput) Validate() error {

	if v.Features != nil {
		featuresValidValues := []string{"1"}
		featuresParameterValue := fmt.Sprint(*v.Features)

		featuresIsValid := false
		for _, value := range featuresValidValues {
			if value == featuresParameterValue {
				featuresIsValid = true
			}
		}

		if !featuresIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "Features",
				ParameterValue: featuresParameterValue,
				AllowedValues:  featuresValidValues,
			}
		}
	}

	if v.IPNetwork == nil {
		return errors.ParameterRequiredError{
			ParameterName: "IPNetwork",
			ParentName:    "JoinRouterInput",
		}
	}

	if v.Router == nil {
		return errors.ParameterRequiredError{
			ParameterName: "Router",
			ParentName:    "JoinRouterInput",
		}
	}

	if v.VxNet == nil {
		return errors.ParameterRequiredError{
			ParameterName: "VxNet",
			ParentName:    "JoinRouterInput",
		}
	}

	return nil
}

type JoinRouterOutput struct {
	Message  *string `json:"message" name:"message"`
	Action   *string `json:"action" name:"action" location:"elements"`
	JobID    *string `json:"job_id" name:"job_id" location:"elements"`
	RetCode  *int    `json:"ret_code" name:"ret_code" location:"elements"`
	RouterID *string `json:"router_id" name:"router_id" location:"elements"`
	VxNetID  *string `json:"vxnet_id" name:"vxnet_id" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/router/leave_router.html
func (s *RouterService) LeaveRouter(i *LeaveRouterInput) (*LeaveRouterOutput, error) {
	if i == nil {
		i = &LeaveRouterInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "LeaveRouter",
		RequestMethod: "GET",
	}

	x := &LeaveRouterOutput{}
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

type LeaveRouterInput struct {
	Router *string   `json:"router" name:"router" location:"params"` // Required
	VxNets []*string `json:"vxnets" name:"vxnets" location:"params"` // Required
}

func (v *LeaveRouterInput) Validate() error {

	if v.Router == nil {
		return errors.ParameterRequiredError{
			ParameterName: "Router",
			ParentName:    "LeaveRouterInput",
		}
	}

	if len(v.VxNets) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "VxNets",
			ParentName:    "LeaveRouterInput",
		}
	}

	return nil
}

type LeaveRouterOutput struct {
	Message  *string   `json:"message" name:"message"`
	Action   *string   `json:"action" name:"action" location:"elements"`
	JobID    *string   `json:"job_id" name:"job_id" location:"elements"`
	RetCode  *int      `json:"ret_code" name:"ret_code" location:"elements"`
	RouterID *string   `json:"router_id" name:"router_id" location:"elements"`
	VxNets   []*string `json:"vxnets" name:"vxnets" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/router/modify_router_attributes.html
func (s *RouterService) ModifyRouterAttributes(i *ModifyRouterAttributesInput) (*ModifyRouterAttributesOutput, error) {
	if i == nil {
		i = &ModifyRouterAttributesInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "ModifyRouterAttributes",
		RequestMethod: "GET",
	}

	x := &ModifyRouterAttributesOutput{}
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

type ModifyRouterAttributesInput struct {
	Description *string `json:"description" name:"description" location:"params"`
	DYNIPEnd    *string `json:"dyn_ip_end" name:"dyn_ip_end" location:"params"`
	DYNIPStart  *string `json:"dyn_ip_start" name:"dyn_ip_start" location:"params"`
	EIP         *string `json:"eip" name:"eip" location:"params"`
	// Features's available values: 1, 2
	Features      *int    `json:"features" name:"features" location:"params"`
	Router        *string `json:"router" name:"router" location:"params"` // Required
	RouterName    *string `json:"router_name" name:"router_name" location:"params"`
	SecurityGroup *string `json:"security_group" name:"security_group" location:"params"`
	VxNet         *string `json:"vxnet" name:"vxnet" location:"params"`
}

func (v *ModifyRouterAttributesInput) Validate() error {

	if v.Features != nil {
		featuresValidValues := []string{"1", "2"}
		featuresParameterValue := fmt.Sprint(*v.Features)

		featuresIsValid := false
		for _, value := range featuresValidValues {
			if value == featuresParameterValue {
				featuresIsValid = true
			}
		}

		if !featuresIsValid {
			return errors.ParameterValueNotAllowedError{
				ParameterName:  "Features",
				ParameterValue: featuresParameterValue,
				AllowedValues:  featuresValidValues,
			}
		}
	}

	if v.Router == nil {
		return errors.ParameterRequiredError{
			ParameterName: "Router",
			ParentName:    "ModifyRouterAttributesInput",
		}
	}

	return nil
}

type ModifyRouterAttributesOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/router/modify_router_static_attributes.html
func (s *RouterService) ModifyRouterStaticAttributes(i *ModifyRouterStaticAttributesInput) (*ModifyRouterStaticAttributesOutput, error) {
	if i == nil {
		i = &ModifyRouterStaticAttributesInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "ModifyRouterStaticAttributes",
		RequestMethod: "GET",
	}

	x := &ModifyRouterStaticAttributesOutput{}
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

type ModifyRouterStaticAttributesInput struct {
	RouterStatic     *string `json:"router_static" name:"router_static" location:"params"` // Required
	RouterStaticName *string `json:"router_static_name" name:"router_static_name" location:"params"`
	Val1             *string `json:"val1" name:"val1" location:"params"`
	Val2             *string `json:"val2" name:"val2" location:"params"`
	Val3             *string `json:"val3" name:"val3" location:"params"`
	Val4             *string `json:"val4" name:"val4" location:"params"`
	Val5             *string `json:"val5" name:"val5" location:"params"`
	Val6             *string `json:"val6" name:"val6" location:"params"`
}

func (v *ModifyRouterStaticAttributesInput) Validate() error {

	if v.RouterStatic == nil {
		return errors.ParameterRequiredError{
			ParameterName: "RouterStatic",
			ParentName:    "ModifyRouterStaticAttributesInput",
		}
	}

	return nil
}

type ModifyRouterStaticAttributesOutput struct {
	Message        *string `json:"message" name:"message"`
	Action         *string `json:"action" name:"action" location:"elements"`
	RetCode        *int    `json:"ret_code" name:"ret_code" location:"elements"`
	RouterStaticID *string `json:"router_static_id" name:"router_static_id" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/router/modify_router_static_entry_attributes.html
func (s *RouterService) ModifyRouterStaticEntryAttributes(i *ModifyRouterStaticEntryAttributesInput) (*ModifyRouterStaticEntryAttributesOutput, error) {
	if i == nil {
		i = &ModifyRouterStaticEntryAttributesInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "ModifyRouterStaticEntryAttributes",
		RequestMethod: "GET",
	}

	x := &ModifyRouterStaticEntryAttributesOutput{}
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

type ModifyRouterStaticEntryAttributesInput struct {
	RouterStaticEntry     *string `json:"router_static_entry" name:"router_static_entry" location:"params"` // Required
	RouterStaticEntryName *string `json:"router_static_entry_name" name:"router_static_entry_name" location:"params"`
	Val1                  *string `json:"val1" name:"val1" location:"params"`
	Val2                  *string `json:"val2" name:"val2" location:"params"`
}

func (v *ModifyRouterStaticEntryAttributesInput) Validate() error {

	if v.RouterStaticEntry == nil {
		return errors.ParameterRequiredError{
			ParameterName: "RouterStaticEntry",
			ParentName:    "ModifyRouterStaticEntryAttributesInput",
		}
	}

	return nil
}

type ModifyRouterStaticEntryAttributesOutput struct {
	Message           *string `json:"message" name:"message"`
	Action            *string `json:"action" name:"action" location:"elements"`
	RetCode           *int    `json:"ret_code" name:"ret_code" location:"elements"`
	RouterStaticEntry *string `json:"router_static_entry" name:"router_static_entry" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/router/poweroff_routers.html
func (s *RouterService) PowerOffRouters(i *PowerOffRoutersInput) (*PowerOffRoutersOutput, error) {
	if i == nil {
		i = &PowerOffRoutersInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "PowerOffRouters",
		RequestMethod: "GET",
	}

	x := &PowerOffRoutersOutput{}
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

type PowerOffRoutersInput struct {
	Routers []*string `json:"routers" name:"routers" location:"params"` // Required
}

func (v *PowerOffRoutersInput) Validate() error {

	if len(v.Routers) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "Routers",
			ParentName:    "PowerOffRoutersInput",
		}
	}

	return nil
}

type PowerOffRoutersOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	JobID   *string `json:"job_id" name:"job_id" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/router/poweron_routers.html
func (s *RouterService) PowerOnRouters(i *PowerOnRoutersInput) (*PowerOnRoutersOutput, error) {
	if i == nil {
		i = &PowerOnRoutersInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "PowerOnRouters",
		RequestMethod: "GET",
	}

	x := &PowerOnRoutersOutput{}
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

type PowerOnRoutersInput struct {
	Routers []*string `json:"routers" name:"routers" location:"params"` // Required
}

func (v *PowerOnRoutersInput) Validate() error {

	if len(v.Routers) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "Routers",
			ParentName:    "PowerOnRoutersInput",
		}
	}

	return nil
}

type PowerOnRoutersOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	JobID   *string `json:"job_id" name:"job_id" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}

// Documentation URL: https://docs.qingcloud.com/api/router/update_routers.html
func (s *RouterService) UpdateRouters(i *UpdateRoutersInput) (*UpdateRoutersOutput, error) {
	if i == nil {
		i = &UpdateRoutersInput{}
	}
	o := &data.Operation{
		Config:        s.Config,
		Properties:    s.Properties,
		APIName:       "UpdateRouters",
		RequestMethod: "GET",
	}

	x := &UpdateRoutersOutput{}
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

type UpdateRoutersInput struct {
	Routers []*string `json:"routers" name:"routers" location:"params"` // Required
}

func (v *UpdateRoutersInput) Validate() error {

	if len(v.Routers) == 0 {
		return errors.ParameterRequiredError{
			ParameterName: "Routers",
			ParentName:    "UpdateRoutersInput",
		}
	}

	return nil
}

type UpdateRoutersOutput struct {
	Message *string `json:"message" name:"message"`
	Action  *string `json:"action" name:"action" location:"elements"`
	JobID   *string `json:"job_id" name:"job_id" location:"elements"`
	RetCode *int    `json:"ret_code" name:"ret_code" location:"elements"`
}
