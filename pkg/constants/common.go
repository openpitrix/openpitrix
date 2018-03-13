// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package constants

import "time"

const (
	prefix                = "openpitrix-"
	ApiGatewayHost        = prefix + "api-gateway"
	RepoManagerHost       = prefix + "repo-manager"
	AppManagerHost        = prefix + "app-manager"
	RuntimeEnvManagerHost = prefix + "runtime-env-manager"
	ClusterManagerHost    = prefix + "cluster-manager"
	JobManagerHost        = prefix + "job-manager"
	TaskManagerHost       = prefix + "task-manager"
	PilotManagerHost      = prefix + "pilot-manager"
	DbCtrlHost            = prefix + "db-ctrl"
	RepoIndexerHost       = prefix + "repo-indexer"
)

const (
	ApiGatewayPort        = 9100 // 91 is similar as Pi, Open[Pi]trix
	RepoManagerPort       = 9101
	AppManagerPort        = 9102
	RuntimeEnvManagerPort = 9103
	ClusterManagerPort    = 9104
	DbCtrlPort            = 9105
	JobManagerPort        = 9106
	TaskManagerPort       = 9107
	RepoIndexerPort       = 9108
	PilotManagerPort      = 9110
)

const (
	StatusActive     = "active"
	StatusDeleted    = "deleted"
	StatusWorking    = "working"
	StatusPending    = "pending"
	StatusSuccessful = "successful"
	StatusFailed     = "failed"
)

const (
	JobLength      = 20
	TaskLength     = 20
	RepoTaskLength = 20
)

const (
	WaitTaskTimeout  = 600 * time.Second
	WaitTaskInterval = 10 * time.Second
)
