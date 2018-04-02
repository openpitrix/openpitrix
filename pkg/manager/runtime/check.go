// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package runtime

import (
	"openpitrix.io/openpitrix/pkg/constants"
)

func (p *Server) checkRuntimeDeleted(runtimeEnvId string) (bool, error) {
	runtimeEnv, err := p.getRuntime(runtimeEnvId)
	if err != nil {
		return true, err
	}
	if runtimeEnv.Status == constants.StatusDeleted {
		return true, nil
	}
	return false, nil
}

func (p *Server) checkRuntimeCredentialDeleted(runtimeEnvCredentialId string) (bool, error) {
	runtimeEnvCredential, err := p.getRuntimeCredential(runtimeEnvCredentialId)
	if err != nil {
		return true, err
	}
	if runtimeEnvCredential.Status == constants.StatusDeleted {
		return true, nil
	}
	return false, nil
}
