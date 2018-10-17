// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package helm

const (
	RuntimeAnnotationKey = "openpitrix_runtime"

	DeploymentFlag  = "-Deployment"
	StatefulSetFlag = "-StatefulSet"
	DaemonSetFlag   = "-DaemonSet"
	ServiceFlag     = "-Service"
	ConfigMapFlag   = "-ConfigMap"
	SecretFlag      = "-Secret"
	PVCFlag         = "-PVC"
	IngressFlag     = "-Ingress"
)
