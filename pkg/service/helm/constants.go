// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package helm

const (
	Provider       = "kubernetes"
	ProviderConfig = `
provider_type: helm
enable: true
`
)

const (
	RuntimeAnnotationKey = "openpitrix_runtime"

	DeploymentFlag  = "-Deployment"
	StatefulSetFlag = "-StatefulSet"
	DaemonSetFlag   = "-DaemonSet"
)
