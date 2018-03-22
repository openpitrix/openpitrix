// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package runtime_env

import (
	"openpitrix.io/openpitrix/pkg/constants"
)

func (p *Server) checkRuntimeEnvDeleted(runtimeEnvId string) (bool, error) {
	runtimeEnv, err := p.getRuntimeEnv(runtimeEnvId)
	if err != nil {
		return true, err
	}
	if runtimeEnv.Status == constants.StatusDeleted {
		return true, nil
	}
	return false, nil
}

func (p *Server) checkRuntimeEnvCredentialDeleted(runtimeEnvCredentialId string) (bool, error) {
	runtimeEnvCredential, err := p.getRuntimeEnvCredential(runtimeEnvCredentialId)
	if err != nil {
		return true, err
	}
	if runtimeEnvCredential.Status == constants.StatusDeleted {
		return true, nil
	}
	return false, nil
}

func (p *Server) checkRuntimeEnvCredentialAttached(runtimeEnvCredentialId string) (bool, error) {
	count, err := p.getAttachmentCountByRuntimeEnvCredentialId(runtimeEnvCredentialId)
	if err != nil {
		return true, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

func (p *Server) checkRuntimeEnvAttached(runtimeEnvId string) (bool, error) {
	count, err := p.getAttachmentCountByRuntimeEnvId(runtimeEnvId)
	if err != nil {
		return true, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}
