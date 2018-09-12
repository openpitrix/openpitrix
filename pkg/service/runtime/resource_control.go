// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package runtime

import (
	"context"
	"time"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/pi"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
)

func getLabelsMap(ctx context.Context, runtimeIds []string) (labelsMap map[string][]*models.RuntimeLabel, err error) {
	var runtimeLabels []*models.RuntimeLabel
	_, err = pi.Global().DB(ctx).
		Select(models.RuntimeLabelColumns...).
		From(constants.TableRuntimeLabel).
		Where(db.Eq(constants.ColumnRuntimeId, runtimeIds)).
		OrderDir(constants.ColumnCreateTime, true).
		Load(&runtimeLabels)
	if err != nil {
		return
	}
	labelsMap = models.RuntimeLabelsMap(runtimeLabels)
	return
}

func getRuntime(ctx context.Context, runtimeId string) (*models.Runtime, error) {
	runtime := &models.Runtime{}
	err := pi.Global().DB(ctx).
		Select(models.RuntimeColumns...).
		From(constants.TableRuntime).
		Where(db.Eq(constants.ColumnRuntimeId, runtimeId)).
		LoadOne(runtime)
	if err != nil {
		return nil, err
	}
	return runtime, nil
}

func getCredentialMap(ctx context.Context, credentialIds ...string) (map[string]*models.RuntimeCredential, error) {
	if len(credentialIds) == 0 {
		return nil, nil
	}
	var runtimeCredneitals []*models.RuntimeCredential
	query := pi.Global().DB(ctx).
		Select(models.RuntimeCredentialColumns...).
		From(constants.TableRuntimeCredential).
		Where(db.Eq(RuntimeCredentialIdColumn, credentialIds))

	_, err := query.Load(&runtimeCredneitals)
	if err != nil {
		return nil, err
	}
	credentialMap := models.RuntimeCredentialMap(runtimeCredneitals)
	return credentialMap, nil
}

func deleteRuntimes(ctx context.Context, runtimeIds []string) error {
	_, err := pi.Global().DB(ctx).
		Update(constants.TableRuntime).
		Set(StatusColumn, constants.StatusDeleted).
		Set(StatusTimeColumn, time.Now()).
		Where(db.Eq(constants.ColumnRuntimeId, runtimeIds)).
		Exec()
	return err
}

func insertRuntime(ctx context.Context, runtime models.Runtime) error {
	_, err := pi.Global().DB(ctx).
		InsertInto(constants.TableRuntime).
		Columns(models.RuntimeColumns...).
		Record(runtime).
		Exec()
	return err
}

func insertRuntimeCredential(ctx context.Context, runtimeCredential models.RuntimeCredential) error {
	_, err := pi.Global().DB(ctx).
		InsertInto(constants.TableRuntimeCredential).
		Columns(models.RuntimeCredentialColumns...).
		Record(runtimeCredential).
		Exec()
	return err
}

func updateRuntimeByMap(ctx context.Context, runtimeId string, attributes map[string]interface{}) error {
	if attributes == nil {
		return nil
	}
	_, err := pi.Global().DB(ctx).
		Update(constants.TableRuntime).
		SetMap(attributes).
		Where(db.Eq(constants.ColumnRuntimeId, runtimeId)).
		Exec()
	return err
}

func createRuntime(ctx context.Context, name, description, provider, url, runtimeCredentialId, zone, userId string) (runtimeId string, err error) {
	newRuntime := models.NewRuntime(name, description, provider, url, runtimeCredentialId, zone, userId)
	err = insertRuntime(ctx, *newRuntime)
	if err != nil {
		return "", err
	}
	return newRuntime.RuntimeId, err
}

func formatRuntimeSet(ctx context.Context, runtimes []*models.Runtime) (pbRuntimes []*pb.Runtime, err error) {
	pbRuntimes = models.RuntimeToPbs(runtimes)
	var runtimeIds []string
	for _, runtime := range runtimes {
		runtimeIds = append(runtimeIds, runtime.RuntimeId)
	}

	labelsMap, err := getLabelsMap(ctx, runtimeIds)
	if err != nil {
		return
	}
	for _, pbRuntime := range pbRuntimes {
		runtimeId := pbRuntime.GetRuntimeId().GetValue()
		pbRuntime.Labels = models.RuntimeLabelsToPbs(labelsMap[runtimeId])
	}

	return
}

func formatRuntimeDetailSet(ctx context.Context, runtimes []*models.Runtime) (pbRuntimeDetails []*pb.RuntimeDetail, err error) {
	pbRuntimes := models.RuntimeToPbs(runtimes)
	var runtimeIds []string
	var credentialIds []string
	runtimeCredentialMap := map[string]string{}
	for _, runtime := range runtimes {
		runtimeIds = append(runtimeIds, runtime.RuntimeId)
		credentialIds = append(credentialIds, runtime.RuntimeCredentialId)
		runtimeCredentialMap[runtime.RuntimeId] = runtime.RuntimeCredentialId
	}

	labelsMap, err := getLabelsMap(ctx, runtimeIds)
	if err != nil {
		return
	}
	runtimeCredentials, err := getCredentialMap(ctx, credentialIds...)
	if err != nil {
		return
	}
	for _, pbRuntime := range pbRuntimes {
		pbRuntimeDetail := new(pb.RuntimeDetail)
		pbRuntimeDetail.Runtime = pbRuntime
		runtimeId := pbRuntime.GetRuntimeId().GetValue()
		credentialId := runtimeCredentialMap[runtimeId]
		pbRuntime.Labels = models.RuntimeLabelsToPbs(labelsMap[runtimeId])
		pbRuntimeDetail.RuntimeCredential = pbutil.ToProtoString(
			CredentialJsonStringToString(
				pbRuntime.Provider.GetValue(), runtimeCredentials[credentialId].Content))
		pbRuntimeDetails = append(pbRuntimeDetails, pbRuntimeDetail)
	}

	return
}

func updateRuntime(ctx context.Context, req *pb.ModifyRuntimeRequest) error {
	attributes := manager.BuildUpdateAttributes(req, NameColumn, DescriptionColumn)
	err := updateRuntimeByMap(ctx, req.RuntimeId.GetValue(), attributes)
	if err != nil {
		return err
	}
	return nil
}

func updateRuntimeCredential(ctx context.Context, credentialId, provider, credential string) error {
	content := CredentialStringToJsonString(provider, credential)
	attributes := map[string]interface{}{
		RuntimeCredentialContentColumn: content,
	}
	_, err := pi.Global().DB(ctx).
		Update(constants.TableRuntimeCredential).
		SetMap(attributes).
		Where(db.Eq(RuntimeCredentialIdColumn, credentialId)).
		Exec()
	return err
}

func createRuntimeCredential(ctx context.Context, provider, content string) (
	runtimeCredentialId string, err error) {

	newRunTimeCredential := models.NewRuntimeCredential(CredentialStringToJsonString(provider, content))
	err = insertRuntimeCredential(ctx, *newRunTimeCredential)
	if err != nil {
		return "", err
	}
	return newRunTimeCredential.RuntimeCredentialId, nil
}
