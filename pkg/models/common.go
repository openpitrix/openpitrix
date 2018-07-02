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
		ColumnAppId, ColumnName, ColumnRepoId, ColumnDescription, ColumnStatus,
		ColumnHome, ColumnIcon, ColumnScreenshots, ColumnMaintainers, ColumnSources,
		ColumnReadme, ColumnOwner, ColumnChartName,
	},
	AppVersionTableName: {
		ColumnVersionId, ColumnAppId, ColumnName, ColumnOwner, ColumnDescription,
		ColumnPackageName, ColumnStatus,
	},
	JobTableName: {
		ColumnJobId, ColumnClusterId, ColumnAppId, ColumnVersionId,
		ColumnExecutor, ColumnProvider, ColumnStatus, ColumnOwner,
	},
	TaskTableName: {
		ColumnJobId, ColumnTaskId, ColumnExecutor, ColumnStatus, ColumnOwner,
	},
	RepoTableName: {
		ColumnRepoId, ColumnName, ColumnType, ColumnVisibility, ColumnStatus,
	},
	RuntimeTableName: {
		ColumnRuntimeId, ColumnProvider, ColumnZone, ColumnStatus, ColumnOwner,
	},
	RepoLabelTableName: {
		ColumnRepoId, ColumnRepoLabelId, ColumnStatus,
	},
	RepoSelectorTableName: {
		ColumnRepoId, ColumnRepoSelectorId, ColumnStatus,
	},
	RepoEventTableName: {
		ColumnRepoEventId, ColumnRepoId, ColumnStatus,
	},
	ClusterTableName: {
		ColumnClusterId, ColumnAppId, ColumnVersionId, ColumnStatus,
		ColumnRuntimeId, ColumnFrontgateId, ColumnOwner,
	},
	ClusterNodeTableName: {
		ColumnClusterId, ColumnNodeId, ColumnStatus, ColumnOwner,
	},
	CategoryTableName: {
		ColumnCategoryId, ColumnStatus, ColumnLocale, ColumnOwner, ColumnName,
	},
}

var SearchWordColumnTable = []string{
	RuntimeTableName,
	AppTableName,
	AppVersionTableName,
	RepoTableName,
	JobTableName,
	TaskTableName,
	ClusterTableName,
	ClusterNodeTableName,
}

// columns that can be search through sql 'like' operator
var SearchColumns = map[string][]string{
	AppTableName: {
		ColumnAppId, ColumnName, ColumnRepoId, ColumnOwner, ColumnChartName, ColumnKeywords,
	},
	AppVersionTableName: {
		ColumnVersionId, ColumnAppId, ColumnName, ColumnDescription, ColumnOwner, ColumnPackageName,
	},
	JobTableName: {
		ColumnJobId, ColumnClusterId, ColumnOwner, ColumnJobAction, ColumnExecutor, ColumnProvider, ColumnExecutor, ColumnProvider,
	},
	TaskTableName: {
		ColumnJobId, ColumnTaskId, ColumnTaskAction, ColumnOwner, ColumnNodeId, ColumnTarget,
	},
	RuntimeTableName: {
		ColumnRuntimeId, ColumnName, ColumnOwner, ColumnProvider, ColumnZone,
	},
	ClusterTableName: {
		ColumnClusterId, ColumnName, ColumnOwner, ColumnAppId, ColumnVersionId, ColumnRuntimeId,
	},
	ClusterNodeTableName: {
		ColumnNodeId, ColumnClusterId, ColumnName, ColumnInstanceId, ColumnVolumeId, ColumnPrivateIp, ColumnRole, ColumnOwner,
	},
	RepoTableName: {
		ColumnName, ColumnDescription,
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
