// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package runtime

import (
	"time"

	"github.com/gocraft/dbr"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
)

func (p *Server) getLabelsMap(runtimeIds []string) (labelsMap map[string][]*models.RuntimeLabel, err error) {
	var runtimeLabels []*models.RuntimeLabel
	_, err = p.Db.
		Select(models.RuntimeLabelColumns...).
		From(models.RuntimeLabelTableName).
		Where(db.Eq(RuntimeIdColumn, runtimeIds)).
		Load(&runtimeLabels)
	if err != nil {
		return
	}
	labelsMap = models.RuntimeLabelsMap(runtimeLabels)
	return
}

func (p *Server) getRuntimeCredential(runtimeCredentialId string) (*models.RuntimeCredential, error) {
	runtimeCredential := &models.RuntimeCredential{}
	err := p.Db.
		Select(models.RuntimeCredentialColumns...).
		From(models.RuntimeCredentialTableName).
		Where(db.Eq(RuntimeCredentialIdColumn, runtimeCredentialId)).
		LoadOne(&runtimeCredential)
	if err != nil {
		return nil, err
	}
	return runtimeCredential, nil
}

func (p *Server) getRuntime(runtimeId string) (*models.Runtime, error) {
	runtime := &models.Runtime{}
	err := p.Db.
		Select(models.RuntimeColumns...).
		From(models.RuntimeTableName).
		Where(db.Eq(RuntimeIdColumn, runtimeId)).
		LoadOne(runtime)
	if err != nil {
		return nil, err
	}
	return runtime, nil
}

func (p *Server) getRuntimeLabelsById(runtimeIds ...string) ([]*models.RuntimeLabel, error) {
	if len(runtimeIds) == 0 {
		return nil, nil
	}
	var runtimeLabels []*models.RuntimeLabel
	query := p.Db.
		Select(models.RuntimeLabelColumns...).
		From(models.RuntimeLabelTableName).
		Where(db.Eq(RuntimeIdColumn, runtimeIds))

	_, err := query.Load(&runtimeLabels)
	if err != nil {
		return nil, err
	}
	return runtimeLabels, nil
}

func (p *Server) getCredentialMap(credentialIds ...string) (map[string]*models.RuntimeCredential, error) {
	if len(credentialIds) == 0 {
		return nil, nil
	}
	var runtimeCredneitals []*models.RuntimeCredential
	query := p.Db.
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

func (p *Server) getRuntimesByFilterCondition(filterCondition dbr.Builder, limit, offset uint64) (
	runtimes []*models.Runtime, count uint32, err error) {
	query := p.Db.
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

func (p *Server) insertRuntimeLabels(runtimeId string, labelMap map[string]string) error {
	if len(labelMap) == 0 {
		return nil
	}
	insertQuery := p.Db.
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
func (p *Server) deleteRuntimeLabels(runtimeId string, labelMap map[string]string) error {
	if len(labelMap) == 0 {
		return nil
	}
	var conditions []dbr.Builder
	for labelKey, labelValue := range labelMap {
		conditions = append(conditions, BuildDeleteLabelFilterCondition(runtimeId, labelKey, labelValue))
	}
	_, err := p.Db.
		DeleteFrom(models.RuntimeLabelTableName).
		Where(db.Or(conditions...)).
		Exec()
	if err != nil {
		return err
	}
	return nil
}

func (p *Server) deleteRuntimes(runtimeIds []string) error {
	_, err := p.Db.
		Update(models.RuntimeTableName).
		Set(StatusColumn, constants.StatusDeleted).
		Set(StatusTimeColumn, time.Now()).
		Where(db.Eq(RuntimeIdColumn, runtimeIds)).
		Exec()
	return err
}

func (p *Server) insertRuntime(runtime models.Runtime) error {
	_, err := p.Db.
		InsertInto(models.RuntimeTableName).
		Columns(models.RuntimeColumns...).
		Record(runtime).
		Exec()
	return err
}

func (p *Server) insertRuntimeCredential(runtimeCredential models.RuntimeCredential) error {
	_, err := p.Db.
		InsertInto(models.RuntimeCredentialTableName).
		Columns(models.RuntimeCredentialColumns...).
		Record(runtimeCredential).
		Exec()
	return err
}

func (p *Server) updateRuntimeByMap(runtimeId string, attributes map[string]interface{}) error {
	if attributes == nil {
		return nil
	}
	_, err := p.Db.
		Update(models.RuntimeTableName).
		SetMap(attributes).
		Where(db.Eq(RuntimeIdColumn, runtimeId)).
		Exec()
	return err
}

func BuildDeleteLabelFilterCondition(runtimeId, labelKey, labelValue string) dbr.Builder {
	var conditions []dbr.Builder
	conditions = append(conditions, db.Eq(RuntimeIdColumn, runtimeId))
	conditions = append(conditions, db.Eq(LabelKeyColumn, labelKey))
	conditions = append(conditions, db.Eq(LabelValueColumn, labelValue))
	return db.And(conditions...)
}

func (p *Server) createRuntime(name, description, provider, url, runtimeCredentialId, zone, userId string) (runtimeId string, err error) {
	newRuntime := models.NewRuntime(name, description, provider, url, runtimeCredentialId, zone, userId)
	err = p.insertRuntime(*newRuntime)
	if err != nil {
		return "", err
	}
	return newRuntime.RuntimeId, err
}

func (p *Server) createRuntimeLabels(runtimeId, labelString string) error {
	labelMap, err := LabelStringToMap(labelString)
	if err != nil {
		return err
	}
	err = p.insertRuntimeLabels(runtimeId, labelMap)
	if err != nil {
		return err
	}
	return nil
}

func (p *Server) formatRuntimeSet(runtimes []*models.Runtime) (pbRuntimes []*pb.Runtime, err error) {
	pbRuntimes = models.RuntimeToPbs(runtimes)
	var runtimeIds []string
	for _, runtime := range runtimes {
		runtimeIds = append(runtimeIds, runtime.RuntimeId)
	}

	labelsMap, err := p.getLabelsMap(runtimeIds)
	if err != nil {
		return
	}
	for _, pbRuntime := range pbRuntimes {
		runtimeId := pbRuntime.GetRuntimeId().GetValue()
		pbRuntime.Labels = models.RuntimeLabelsToPbs(labelsMap[runtimeId])
	}

	return
}

func (p *Server) formatRuntimeDetailSet(runtimes []*models.Runtime) (pbRuntimeDetails []*pb.RuntimeDetail, err error) {
	pbRuntimes := models.RuntimeToPbs(runtimes)
	var runtimeIds []string
	var credentialIds []string
	runtimeCredentialMap := map[string]string{}
	for _, runtime := range runtimes {
		runtimeIds = append(runtimeIds, runtime.RuntimeId)
		credentialIds = append(credentialIds, runtime.RuntimeCredentialId)
		runtimeCredentialMap[runtime.RuntimeId] = runtime.RuntimeCredentialId
	}

	labelsMap, err := p.getLabelsMap(runtimeIds)
	if err != nil {
		return
	}
	runtimeCredentials, err := p.getCredentialMap(credentialIds...)
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
			RuntimeCredentialJsonStringToString(
				pbRuntime.Provider.GetValue(), runtimeCredentials[credentialId].Content))
		pbRuntimeDetails = append(pbRuntimeDetails, pbRuntimeDetail)
	}

	return
}

func (p *Server) updateRuntime(req *pb.ModifyRuntimeRequest) error {
	attributes := manager.BuildUpdateAttributes(req, NameColumn, DescriptionColumn)
	err := p.updateRuntimeByMap(req.RuntimeId.GetValue(), attributes)
	if err != nil {
		return err
	}
	return nil
}

func (p *Server) updateRuntimeLabels(runtimeId string, labelString string) error {
	newLabelMap, err := LabelStringToMap(labelString)
	if err != nil {
		return err
	}
	oldRuntimeLabels, err := p.getRuntimeLabelsById(runtimeId)
	if err != nil {
		return err
	}
	oldLabelMap := LabelStructToMap(oldRuntimeLabels)
	additionLabelMap, deletionLabelMap := LabelMapDiff(oldLabelMap, newLabelMap)
	err = p.deleteRuntimeLabels(runtimeId, deletionLabelMap)
	if err != nil {
		return err
	}
	err = p.insertRuntimeLabels(runtimeId, additionLabelMap)
	if err != nil {
		return err
	}
	return nil
}

func (p *Server) createRuntimeCredential(provider, content string) (
	runtimeCredentialId string, err error) {

	newRunTimeCredential := models.NewRuntimeCredential(RuntimeCredentialStringToJsonString(provider, content))
	err = p.insertRuntimeCredential(*newRunTimeCredential)
	if err != nil {
		return "", err
	}
	return newRunTimeCredential.RuntimeCredentialId, nil
}
