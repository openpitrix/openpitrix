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
	JobTableName: {
		"job_id", "cluster_id", "app_id", "app_version", "status",
	},
	TaskTableName: {
		"job_id", "task_id", "status",
	},
	RepoTableName: {
		"repo_id", "name", "visibility", "status",
	},
}

var SearchWordColumnTable = []string{
	AppTableName,
}

// columns that can be search through sql 'like' operator
var SearchColumns = map[string][]string{
	AppTableName: {
		"app_id", "name", "repo_id", "owner", "chart_name",
	},
	JobTableName: {
		"job_id", "cluster_id", "app_id", "app_version", "status",
	},
	TaskTableName: {
		"job_id", "task_id", "status",
	},
}

func GetColumnsFromStruct(s interface{}) []string {
	names := structs.Names(s)
	for i, name := range names {
		names[i] = utils.CamelCaseToUnderscore(name)
	}
	return names
}
