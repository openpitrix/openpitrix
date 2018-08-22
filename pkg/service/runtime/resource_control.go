// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package runtime

import (
	"context"
	"time"

	"github.com/gocraft/dbr"

	"openpitrix.io/openpitrix/pkg/pi"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
)

func (p *Server) getLabelsMap(ctx context.Context, runtimeIds []string) (labelsMap map[string][]*models.RuntimeLabel, err error) {
	var runtimeLabels []*models.RuntimeLabel
	_, err = pi.Global().DB(ctx).
		Select(models.RuntimeLabelColumns...).
		From(models.RuntimeLabelTableName).
		Where(db.Eq(models.ColumnRuntimeId, runtimeIds)).
		Load(&runtimeLabels)
	if err != nil {
		return
	}
	labelsMap = models.RuntimeLabelsMap(runtimeLabels)
	return
}

func (p *Server) getRuntimeCredential(ctx context.Context, runtimeCredentialId string) (*models.RuntimeCredential, error) {
	runtimeCredential := &models.RuntimeCredential{}
	err := pi.Global().DB(ctx).
		Select(models.RuntimeCredentialColumns...).
		From(models.RuntimeCredentialTableName).
		Where(db.Eq(RuntimeCredentialIdColumn, runtimeCredentialId)).
		LoadOne(&runtimeCredential)
	if err != nil {
		return nil, err
	}
	return runtimeCredential, nil
}

func (p *Server) getRuntime(ctx context.Context, runtimeId string) (*models.Runtime, error) {
	runtime := &models.Runtime{}
	err := pi.Global().DB(ctx).
		Select(models.RuntimeColumns...).
		From(models.RuntimeTableName).
		Where(db.Eq(models.ColumnRuntimeId, runtimeId)).
		LoadOne(runtime)
	if err != nil {
		return nil, err
	}
	return runtime, nil
}

func (p *Server) getRuntimeLabelsById(ctx context.Context, runtimeIds ...string) ([]*models.RuntimeLabel, error) {
	if len(runtimeIds) == 0 {
		return nil, nil
	}
	var runtimeLabels []*models.RuntimeLabel
	query := pi.Global().DB(ctx).
		Select(models.RuntimeLabelColumns...).
		From(models.RuntimeLabelTableName).
		Where(db.Eq(models.ColumnRuntimeId, runtimeIds))

	_, err := query.Load(&runtimeLabels)
	if err != nil {
		return nil, err
	}
	return runtimeLabels, nil
}

func (p *Server) getCredentialMap(ctx context.Context, credentialIds ...string) (map[string]*models.RuntimeCredential, error) {
	if len(credentialIds) == 0 {
		return nil, nil
	}
	var runtimeCredneitals []*models.RuntimeCredential
	query := pi.Global().DB(ctx).
		Select(models.RuntimeCredentialColumns...).
		From(models.RuntimeCredentialTableName).
		Where(db.Eq(RuntimeCredentialIdColumn, credentialIds))

	_, err := query.Load(&runtimeCredneitals)
	if err != nil {
		return nil, err
	}
	credentialMap := models.RuntimeCredentialMap(runtimeCredneitals)
	return credentialMap, nil
}

func (p *Server) getRuntimesByFilterCondition(ctx context.Context, filterCondition dbr.Builder, limit, offset uint64) (
	runtimes []*models.Runtime, count uint32, err error) {
	query := pi.Global().DB(ctx).
		Select(models.RuntimeColumns...).
		From(models.RuntimeTableName).
		Offset(offset).
		Limit(limit).Where(filterCondition)
	_, err = query.Load(&runtimes)
	if err != nil {
		return nil, 0, err
	}
	count, err = query.Count()
	if err != nil {
		return nil, 0, err
	}
	return runtimes, count, nil
}

func (p *Server) insertRuntimeLabels(ctx context.Context, runtimeId string, labelMap map[string]string) error {
	if len(labelMap) == 0 {
		return nil
	}
	insertQuery := pi.Global().DB(ctx).
		InsertInto(models.RuntimeLabelTableName).
		Columns(models.RuntimeLabelColumns...)
	for labelKey, labelValue := range labelMap {
		insertQuery = insertQuery.Record(
			models.NewRuntimeLabel(runtimeId, labelKey, labelValue),
		)
	}
	_, err := insertQuery.Exec()
	if err != nil {
		return err
	}
	return nil
}

// deleteRuntimeLabels by runtimeId and labelMap.if len(labelMap) =0 , record are not deleted.
func (p *Server) deleteRuntimeLabels(ctx context.Context, runtimeId string, labelMap map[string]string) error {
	if len(labelMap) == 0 {
		return nil
	}
	var conditions []dbr.Builder
	for labelKey, labelValue := range labelMap {
		conditions = append(conditions, BuildDeleteLabelFilterCondition(runtimeId, labelKey, labelValue))
	}
	_, err := pi.Global().DB(ctx).
		DeleteFrom(models.RuntimeLabelTableName).
		Where(db.Or(conditions...)).
		Exec()
	if err != nil {
		return err
	}
	return nil
}

func (p *Server) deleteRuntimes(ctx context.Context, runtimeIds []string) error {
	_, err := pi.Global().DB(ctx).
		Update(models.RuntimeTableName).
		Set(StatusColumn, constants.StatusDeleted).
		Set(StatusTimeColumn, time.Now()).
		Where(db.Eq(models.ColumnRuntimeId, runtimeIds)).
		Exec()
	return err
}

func (p *Server) insertRuntime(ctx context.Context, runtime models.Runtime) error {
	_, err := pi.Global().DB(ctx).
		InsertInto(models.RuntimeTableName).
		Columns(models.RuntimeColumns...).
		Record(runtime).
		Exec()
	return err
}

func (p *Server) insertRuntimeCredential(ctx context.Context, runtimeCredential models.RuntimeCredential) error {
	_, err := pi.Global().DB(ctx).
		InsertInto(models.RuntimeCredentialTableName).
		Columns(models.RuntimeCredentialColumns...).
		Record(runtimeCredential).
		Exec()
	return err
}

func (p *Server) updateRuntimeByMap(ctx context.Context, runtimeId string, attributes map[string]interface{}) error {
	if attributes == nil {
		return nil
	}
	_, err := pi.Global().DB(ctx).
		Update(models.RuntimeTableName).
		SetMap(attributes).
		Where(db.Eq(models.ColumnRuntimeId, runtimeId)).
		Exec()
	return err
}

func BuildDeleteLabelFilterCondition(runtimeId, labelKey, labelValue string) dbr.Builder {
	var conditions []dbr.Builder
	conditions = append(conditions, db.Eq(models.ColumnRuntimeId, runtimeId))
	conditions = append(conditions, db.Eq(LabelKeyColumn, labelKey))
	conditions = append(conditions, db.Eq(LabelValueColumn, labelValue))
	return db.And(conditions...)
}

func (p *Server) createRuntime(ctx context.Context, name, description, provider, url, runtimeCredentialId, zone, userId string) (runtimeId string, err error) {
	newRuntime := models.NewRuntime(name, description, provider, url, runtimeCredentialId, zone, userId)
	err = p.insertRuntime(ctx, *newRuntime)
	if err != nil {
		return "", err
	}
	return newRuntime.RuntimeId, err
}

func (p *Server) createRuntimeLabels(ctx context.Context, runtimeId, labelString string) error {
	labelMap, err := LabelStringToMap(labelString)
	if err != nil {
		return err
	}
	err = p.insertRuntimeLabels(ctx, runtimeId, labelMap)
	if err != nil {
		return err
	}
	return nil
}

func (p *Server) formatRuntimeSet(ctx context.Context, runtimes []*models.Runtime) (pbRuntimes []*pb.Runtime, err error) {
	pbRuntimes = models.RuntimeToPbs(runtimes)
	var runtimeIds []string
	for _, runtime := range runtimes {
		runtimeIds = append(runtimeIds, runtime.RuntimeId)
	}

	labelsMap, err := p.getLabelsMap(ctx, runtimeIds)
	if err != nil {
		return
	}
	for _, pbRuntime := range pbRuntimes {
		runtimeId := pbRuntime.GetRuntimeId().GetValue()
		pbRuntime.Labels = models.RuntimeLabelsToPbs(labelsMap[runtimeId])
	}

	return
}

func (p *Server) formatRuntimeDetailSet(ctx context.Context, runtimes []*models.Runtime) (pbRuntimeDetails []*pb.RuntimeDetail, err error) {
	pbRuntimes := models.RuntimeToPbs(runtimes)
	var runtimeIds []string
	var credentialIds []string
	runtimeCredentialMap := map[string]string{}
	for _, runtime := range runtimes {
		runtimeIds = append(runtimeIds, runtime.RuntimeId)
		credentialIds = append(credentialIds, runtime.RuntimeCredentialId)
		runtimeCredentialMap[runtime.RuntimeId] = runtime.RuntimeCredentialId
	}

	labelsMap, err := p.getLabelsMap(ctx, runtimeIds)
	if err != nil {
		return
	}
	runtimeCredentials, err := p.getCredentialMap(ctx, credentialIds...)
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

func (p *Server) updateRuntime(ctx context.Context, req *pb.ModifyRuntimeRequest) error {
	attributes := manager.BuildUpdateAttributes(req, NameColumn, DescriptionColumn)
	err := p.updateRuntimeByMap(ctx, req.RuntimeId.GetValue(), attributes)
	if err != nil {
		return err
	}
	return nil
}

func (p *Server) updateRuntimeLabels(ctx context.Context, runtimeId string, labelString string) error {
	newLabelMap, err := LabelStringToMap(labelString)
	if err != nil {
		return err
	}
	oldRuntimeLabels, err := p.getRuntimeLabelsById(ctx, runtimeId)
	if err != nil {
		return err
	}
	oldLabelMap := LabelStructToMap(oldRuntimeLabels)
	additionLabelMap, deletionLabelMap := LabelMapDiff(oldLabelMap, newLabelMap)
	err = p.deleteRuntimeLabels(ctx, runtimeId, deletionLabelMap)
	if err != nil {
		return err
	}
	err = p.insertRuntimeLabels(ctx, runtimeId, additionLabelMap)
	if err != nil {
		return err
	}
	return nil
}

func (p *Server) updateRuntimeCredential(ctx context.Context, credentialId, provider, credential string) error {
	content := CredentialStringToJsonString(provider, credential)
	attributes := map[string]interface{}{
		RuntimeCredentialContentColumn: content,
	}
	_, err := pi.Global().DB(ctx).
		Update(models.RuntimeCredentialTableName).
		SetMap(attributes).
		Where(db.Eq(RuntimeCredentialIdColumn, credentialId)).
		Exec()
	return err
}

func (p *Server) createRuntimeCredential(ctx context.Context, provider, content string) (
	runtimeCredentialId string, err error) {

	newRunTimeCredential := models.NewRuntimeCredential(CredentialStringToJsonString(provider, content))
	err = p.insertRuntimeCredential(ctx, *newRunTimeCredential)
	if err != nil {
		return "", err
	}
	return newRunTimeCredential.RuntimeCredentialId, nil
}
