// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package runtime

import (
	"context"
	"fmt"
	"time"

	"github.com/ghodss/yaml"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/pi"
	"openpitrix.io/openpitrix/pkg/plugins"
)

func getRuntimeCredentials(ctx context.Context, credentialIds ...string) ([]*models.RuntimeCredential, error) {
	if len(credentialIds) == 0 {
		return nil, nil
	}
	var runtimeCredentials []*models.RuntimeCredential
	query := pi.Global().DB(ctx).
		Select(models.RuntimeCredentialColumns...).
		From(constants.TableRuntimeCredential).
		Where(db.Eq(constants.ColumnRuntimeCredentialId, credentialIds))

	_, err := query.Load(&runtimeCredentials)
	if err != nil {
		return nil, err
	}
	return runtimeCredentials, nil
}

func formatRuntimeDetailSet(ctx context.Context, runtimes []*models.Runtime) (pbRuntimeDetails []*pb.RuntimeDetail, err error) {
	pbRuntimes := models.RuntimeToPbs(runtimes)
	var credentialIds []string
	for _, runtime := range runtimes {
		credentialIds = append(credentialIds, runtime.RuntimeCredentialId)
	}

	runtimeCredentials, err := getRuntimeCredentials(ctx, credentialIds...)
	if err != nil {
		return
	}
	runtimeCredentialMap := models.RuntimeCredentialMap(runtimeCredentials)

	for _, pbRuntime := range pbRuntimes {
		pbRuntimeDetail := new(pb.RuntimeDetail)
		pbRuntimeDetail.Runtime = pbRuntime
		credentialId := pbRuntime.GetRuntimeCredentialId().GetValue()
		runtimeCredential := runtimeCredentialMap[credentialId]
		runtimeCredential.RuntimeCredentialContent, err = encodeRuntimeCredentialContent(runtimeCredential.Provider, runtimeCredential.RuntimeCredentialContent)
		if err != nil {
			return
		}
		pbRuntimeDetail.RuntimeCredential = models.RuntimeCredentialToPb(runtimeCredential)
		pbRuntimeDetails = append(pbRuntimeDetails, pbRuntimeDetail)
	}

	return
}

// Decode to json string,then can be inserted into db
func decodeRuntimeCredentialContent(provider, runtimeCredentialContent string) (string, error) {
	if plugins.IsVmbasedProviders(provider) {
		return runtimeCredentialContent, nil
	} else if constants.ProviderKubernetes == provider {
		content, err := yaml.YAMLToJSON([]byte(runtimeCredentialContent))
		if err != nil {
			return "", err
		}
		return string(content), nil
	}
	return "", fmt.Errorf("unsupport provider [%s]", provider)
}

// Encode to origin string,then can be used in the provider interface
func encodeRuntimeCredentialContent(provider, content string) (string, error) {
	if plugins.IsVmbasedProviders(provider) {
		return content, nil
	} else if constants.ProviderKubernetes == provider {
		content, err := yaml.JSONToYAML([]byte(content))
		if err != nil {
			return "", err
		}
		return string(content), nil
	}
	return "", fmt.Errorf("unsupport provider [%s]", provider)
}

func deleteRuntime(ctx context.Context, runtimeIds []string) error {
	if len(runtimeIds) == 0 {
		return nil
	}
	_, err := pi.Global().DB(ctx).
		Update(constants.TableRuntime).
		Set(constants.ColumnStatus, constants.StatusDeleted).
		Set(constants.ColumnStatusTime, time.Now()).
		Where(db.Eq(constants.ColumnRuntimeId, runtimeIds)).
		Exec()
	return err
}
