// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"github.com/fatih/structs"

	"openpitrix.io/openpitrix/pkg/util/stringutil"
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
		"job_id", "cluster_id", "app_id", "version_id", "executor", "provider", "status", "owner",
	},
	TaskTableName: {
		"job_id", "task_id", "executor", "status", "owner",
	},
	RepoTableName: {
		"repo_id", "name", "type", "visibility", "status",
	},
	RuntimeTableName: {
		"runtime_id", "provider", "zone", "status", "owner",
	},
	RepoLabelTableName: {
		"repo_id", "repo_label_id", "status",
	},
	RepoSelectorTableName: {
		"repo_id", "repo_selector_id", "status",
	},
	RepoEventTableName: {
		"repo_event_id", "repo_id", "status",
	},
	ClusterTableName: {
		"cluster_id", "app_id", "version_id", "status", "runtime_id", "frontgate_id", "owner",
	},
	ClusterNodeTableName: {
		"cluster_id", "node_id", "status", "owner",
	},
	CategoryTableName: {
		"category_id", "status", "locale", "owner", "name",
	},
}

var SearchWordColumnTable = []string{
	RuntimeTableName,
	AppTableName, AppVersionTableName,
}

// columns that can be search through sql 'like' operator
var SearchColumns = map[string][]string{
	AppTableName: {
		"app_id", "name", "repo_id", "owner", "chart_name", "keywords",
	},
	AppVersionTableName: {
		"version_id", "app_id", "name", "description", "owner", "package_name",
	},
	JobTableName: {
		"executor", "provider",
	},
	TaskTableName: {
		"job_id", "task_id", "status",
	},
	RuntimeTableName: {
		"runtime_id", "name",
	},
	ClusterTableName: {
		"name", "description",
	},
	ClusterNodeTableName: {
		"name",
	},
}

func GetColumnsFromStruct(s interface{}) []string {
	names := structs.Names(s)
	for i, name := range names {
		names[i] = stringutil.CamelCaseToUnderscore(name)
	}
	return names
}

func GetColumnsFromStructWithPrefix(prefix string, s interface{}) []string {
	names := structs.Names(s)
	for i, name := range names {
		names[i] = prefix + "." + stringutil.CamelCaseToUnderscore(name)
	}
	return names
}
