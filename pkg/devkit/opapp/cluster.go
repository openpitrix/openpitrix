// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package opapp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"regexp"

	"openpitrix.io/openpitrix/pkg/util/jsonutil"
)

type ClusterConfTemplate struct {
	Raw string
}

func replaceTemplateExpression(s string) string {
	var re = regexp.MustCompile(`"{{`)
	s = re.ReplaceAllString(s, `{{ getv `)
	re = regexp.MustCompile(`}}"`)
	s = re.ReplaceAllString(s, `}}`)
	return s
}

func getv(input interface{}) interface{} {
	switch i := input.(type) {
	case string:
		return template.HTML(jsonutil.ToString(i))
	case nil:
		return template.HTML("null")
	default:
		return i
	}
}

func (c *ClusterConfTemplate) Render(input jsonutil.Json) (ClusterConf, error) {
	var cluster ClusterConf
	var tmpl *template.Template
	raw := replaceTemplateExpression(c.Raw)
	tmpl, err := template.New("render_cluster_config").Funcs(template.FuncMap{
		"getv": getv,
	}).Parse(raw)
	if err != nil {
		return cluster, fmt.Errorf("failed to parse cluster.json.tmpl: %+v", err)
	}
	b := bytes.NewBuffer([]byte{})
	err = tmpl.Execute(b, input.Interface())
	if err != nil {
		return cluster, fmt.Errorf("failed to render cluster.json: %+v", err)
	}
	cluster.RenderJson = b.String()
	err = json.Unmarshal(b.Bytes(), &cluster)
	if err != nil {
		return cluster, fmt.Errorf("failed to decode cluster.json: %+v", err)
	}
	return cluster, nil
}

type ClusterConf struct {
	RenderJson string `json:"-"`

	AppId                      string                 `json:"app_id"`
	VersionId                  string                 `json:"version_id"`
	GlobalUuid                 string                 `json:"global_uuid"`
	Name                       string                 `json:"name"`
	Description                string                 `json:"description"`
	Subnet                     string                 `json:"subnet"`
	Links                      map[string]string      `json:"links"`
	BackupPolicy               string                 `json:"backup_policy"`
	IncrementalBackupSupported *bool                  `json:"incremental_backup_supported"`
	UpgradePolicy              []string               `json:"upgrade_policy"`
	Nodes                      []Node                 `json:"nodes"`
	Env                        map[string]interface{} `json:"env"`
	AdvancedActions            []string               `json:"advanced_actions"`
	Endpoints                  map[string]struct {
		Port     uint32 `json:"port"`
		Protocol string `json:"protocol"`
	} `json:"endpoints"`
	MetadataRootAccess *bool        `json:"metadata_root_access"`
	HealthCheck        *HealthCheck `json:"health_check"`
	Monitor            *Monitor     `json:"monitor"`
	DisplayTabs        struct {
		DisplayTabsItems map[string]struct {
			Cmd              string   `json:"cmd"`
			RolesToExecuteOn []string `json:"roles_to_execute_on"`
			Description      string   `json:"description"`
			Timeout          uint32   `json:"timeout"`
		}
	} `json:"display_tabs"`
}

type ServiceParams struct {
	Params map[string]interface{}
}

type HealthCheck struct {
	Enable             *bool  `json:"enable"`
	IntervalSec        uint32 `json:"interval_sec"`
	TimeoutSec         uint32 `json:"timeout_sec"`
	ActionTimeoutSec   uint32 `json:"action_timeout_sec"`
	HealthyThreshold   uint32 `json:"healthy_threshold"`
	UnhealthyThreshold uint32 `json:"unhealthy_threshold"`
	CheckCmd           string `json:"check_cmd"`
	ActionCmd          string `json:"action_cmd"`
}

type Monitor struct {
	Enable *bool  `json:"enable"`
	Cmd    string `json:"cmd"`
	Items  map[string]struct {
		Unit                   string   `json:"unit"`
		ValueType              string   `json:"value_type"`
		StatisticsType         string   `json:"statistics_type"`
		ScaleFactorWhenDisplay uint32   `json:"scale_factor_when_display"`
		Enums                  []string `json:"enums"`
	} `json:"items"`
	Groups  map[string][]string `json:"groups"`
	Display []string            `json:"display"`
	Alarm   []string            `json:"alarm"`
}

type Node struct {
	Role            string   `json:"role"`
	AdvancedActions []string `json:"advanced_actions"`
	Loadbalancer    []struct {
		Listener string `json:"listener"`
		Port     uint32 `json:"port"`
		Policy   string `json:"policy"`
	} `json:"loadbalancer"`
	Container struct {
		Type  string `json:"type"`
		Image string `json:"image"`
	} `json:"container"`
	Count  uint32 `json:"count"`
	CPU    uint32 `json:"cpu"`
	Memory uint32 `json:"memory"`
	GPU    uint32 `json:"gpu"`
	Volume struct {
		Size         uint32      `json:"size"`
		InstanceSize uint32      `json:"instance_size"`
		MountPoint   interface{} `json:"mount_point"`
		MountOptions string      `json:"mount_options"`
		Filesystem   string      `json:"filesystem"`
	} `json:"volume"`
	Replica               uint32                 `json:"replica"`
	Passphraseless        string                 `json:"passphraseless"`
	VerticalScalingPolicy string                 `json:"vertical_scaling_policy"`
	UserAccess            *bool                  `json:"user_access"`
	Services              map[string]interface{} `json:"services"`
	ServerIDUpperBound    uint32                 `json:"server_id_upper_bound"`
	Env                   map[string]interface{} `json:"env"`
	AgentInstalled        *bool                  `json:"agent_installed"`
	CustomMetadata        map[string]interface{} `json:"custom_metadata"`
	HealthCheck           *HealthCheck           `json:"health_check"`
	Monitor               *Monitor               `json:"monitor"`
}

type Service struct {
	NodesToExecuteOn *uint32                `json:"nodes_to_execute_on"`
	PostStartService *bool                  `json:"post_start_service"`
	PostStopService  *bool                  `json:"post_stop_service"`
	Timeout          *uint32                `json:"timeout"`
	ServiceParams    map[string]interface{} `json:"service_params"`
	PreCheck         string                 `json:"pre_check"`
	Cmd              string                 `json:"cmd"`
	Order            *uint32                `json:"order"`
}
