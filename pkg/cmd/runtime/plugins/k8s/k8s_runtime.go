// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package k8s_runtime

import (
	"errors"

	"openpitrix.io/openpitrix/pkg/cmd/runtime"
)

func init() {
	runtime.RegisterRuntime(new(K8sRuntime))
}

type K8sRuntime struct{}

func (p *K8sRuntime) Name() string { return "k8s" }

func (p *K8sRuntime) Run(app string, args ...string) error {
	return errors.New("TODO")
}
