// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"github.com/fatih/structs"

	"openpitrix.io/openpitrix/pkg/utils"
)

// columns that can be search through sql '=' operator
var IndexedColumns = map[string][]string{
	AppTableName: {
		"app_id", "name", "repo_id", "description", "status",
		"home", "icon", "screenshots", "maintainers", "sources",
		"readme", "owner", "chart_name",
	},
	AppVersionTableName: {
		"version_id", "app_id", "name", "owner", "description",
		"package_name", "status",
	},
	JobTableName: {
		"job_id", "cluster_id", "app_id", "app_version", "status",
	},
	TaskTableName: {
		"job_id", "task_id", "status",
	},
	RepoTableName: {
		"repo_id", "name", "type", "visibility", "status",
	},
	RuntimeEnvTableName: {
		"runtime_env_id", "status", "owner",
	},
	RuntimeEnvCredentialTableName: {
		"runtime_env_credential_id", "status", "owner",
	},
	RepoLabelTableName: {
		"repo_id", "repo_label_id", "status",
	},
	RepoSelectorTableName: {
		"repo_id", "repo_selector_id", "status",
	},
	RepoTaskTableName: {
		"repo_task_id", "repo_id", "status",
	},
}

var SearchWordColumnTable = []string{
	RuntimeEnvTableName,
	RuntimeEnvCredentialTableName,
	AppTableName, AppVersionTableName,
}

// columns that can be search through sql 'like' operator
var SearchColumns = map[string][]string{
	AppTableName: {
		"app_id", "name", "repo_id", "owner", "chart_name",
	},
	AppVersionTableName: {
		"version_id", "app_id", "name", "description", "owner", "package_name",
	},
	JobTableName: {
		"job_id", "cluster_id", "app_id", "app_version", "status",
	},
	TaskTableName: {
		"job_id", "task_id", "status",
	},
	RuntimeEnvTableName: {
		"runtime_env_id", "name",
	},
	RuntimeEnvCredentialTableName: {
		"runtime_env_credential_id", "name",
	},
}

func GetColumnsFromStruct(s interface{}) []string {
	names := structs.Names(s)
	for i, name := range names {
		names[i] = utils.CamelCaseToUnderscore(name)
	}
	return names
}

func GetColumnsFromStructWithPrefix(prefix string, s interface{}) []string {
	names := structs.Names(s)
	for i, name := range names {
		names[i] = prefix + "." + utils.CamelCaseToUnderscore(name)
	}
	return names
}
