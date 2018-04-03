// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package runtime

import (
	"openpitrix.io/openpitrix/pkg/constants"
)

func (p *Server) checkRuntimeDeleted(runtimeId string) (bool, error) {
	runtime, err := p.getRuntime(runtimeId)
	if err != nil {
		return true, err
	}
	if runtime.Status == constants.StatusDeleted {
		return true, nil
	}
	return false, nil
}

func (p *Server) checkRuntimeCredentialDeleted(runtimeCredentialId string) (bool, error) {
	runtimeCredential, err := p.getRuntimeCredential(runtimeCredentialId)
	if err != nil {
		return true, err
	}
	if runtimeCredential.Status == constants.StatusDeleted {
		return true, nil
	}
	return false, nil
}
