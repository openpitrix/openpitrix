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
	ActionSetFrontgateConfig     = "SetFrontgateConfig"
	ActionSetDroneConfig         = "SetDroneConfig"
)

const (
	RegisterClustersRootPath         = "clusters"
	RegisterNodeHosts                = "hosts"
	RegisterNodeHost                 = "host"
	RegisterNodeCluster              = "cluster"
	RegisterNodeEnv                  = "env"
	RegisterNodeLoadbalancer         = "loadbalancer"
	RegisterNodeCmd                  = "cmd"
	RegisterNodeCmdId                = "id"
	RegisterNodeCmdTimeout           = "timeout"
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

const (
	MetadataConfPath   = "/opt/openpitrix/conf/"
	OpenPitrixConfFile = "openpitrix.conf"
	DroneConfFile      = "drone.conf"
	FrontgateConfFile  = "frontgate.conf"
	ConfdPath          = "/etc/confd/"
	MetadataLogLevel   = "debug"
	ConfdBackendType   = "libconfd-backend-etcdv3"
	ConfdCmdLogPath    = "/opt/openpitrix/logs/cmd.log"
)

const (
	DefaultLoginPasswd = "Password"
)
