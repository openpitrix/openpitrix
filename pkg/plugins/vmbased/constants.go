// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package vmbased

const (
	ActionRunInstances       = "RunInstances"
	ActionStartInstances     = "StartInstances"
	ActionStopInstances      = "StopInstances"
	ActionTerminateInstances = "TerminateInstances"

	ActionCreateVolumes = "CreateVolumes"
	ActionAttachVolumes = "AttachVolumes"
	ActionDetachVolumes = "DetachVolumes"
	ActionDeleteVolumes = "DeleteVolumes"
	ActionResizeVolumes = "ResizeVolumes"

	ActionFormatAndMountVolume = "FormatAndMountVolume"

	ActionWaitFrontgateAvailable = "WaitFrontgateAvailable"
	ActionRegisterMetadata       = "RegisterMetadata"
	ActionDeregisterMetadata     = "DeregisterMetadata"
	ActionRegisterCmd            = "RegisterCmd"
	ActionDeregesterCmd          = "DeregisterCmd"
	ActionStartConfd             = "StartConfd"
	ActionStopConfd              = "StopConfd"
	ActionGetTaskStatus          = "GetTaskStatus"
)

const (
	RegisterClustersRootPath         = "clusters"
	RegisterNodeHosts                = "hosts"
	RegisterNodeHost                 = "host"
	RegisterNodeCluster              = "cluster"
	RegisterNodeEnv                  = "env"
	RegisterNodeLoadbalancer         = "loadbalancer"
	RegisterNodeCmd                  = "cmd"
	RegisterNodeEndpoint             = "endpoints"
	RegisterNodeAdding               = "adding-hosts"
	RegisterNodeDeleting             = "deleting-hosts"
	RegisterNodeVerticalScalingRoles = "vertical-scaling-roles"
)

const (
	// second
	TimeoutStartConfd           = 600
	TimeoutStopConfd            = 600
	TimeoutDeregister           = 60
	TimeoutRegister             = 60
	TimeoutFormatAndMountVolume = 600
	TimeoutUmountVolume         = 120
	TimeoutSshKeygen            = 120
)
