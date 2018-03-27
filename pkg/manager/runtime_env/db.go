// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package runtime_env

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gocraft/dbr"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/models"
)

func (p *Server) getRuntimeEnvCredential(runtimeEnvCredentialId string) (*models.RuntimeEnvCredential, error) {
	runtimeEnvCredential := &models.RuntimeEnvCredential{}
	err := p.Db.
		Select(models.RuntimeEnvCredentialColumns...).
		From(models.RuntimeEnvCredentialTableName).
		Where(db.Eq(RuntimeEnvCredentialIdColumn, runtimeEnvCredentialId)).
		LoadOne(&runtimeEnvCredential)
	if err != nil {
		return nil, err
	}
	return runtimeEnvCredential, nil
}

func (p *Server) getAttachmentsByRuntimeEnvIds(runtimeEnvIds []string) ([]*models.RuntimeEnvAttachedCredential, error) {
	var runtimeEnvAttachedCredentials []*models.RuntimeEnvAttachedCredential
	_, err := p.Db.
		Select(models.RuntimeEnvAttachedCredentialColumns...).
		From(models.RuntimeEnvAttachedCredentialTableName).
		Where(db.Eq(RuntimeEnvIdColumn, runtimeEnvIds)).
		Load(&runtimeEnvAttachedCredentials)
	if err != nil {
		return nil, err
	}
	return runtimeEnvAttachedCredentials, nil
}

func (p *Server) getAttachmentsByRuntimeEnvCredentialIds(runtimeEnvCredentialIds []string) ([]*models.RuntimeEnvAttachedCredential, error) {
	var runtimeEnvAttachedCredentials []*models.RuntimeEnvAttachedCredential
	_, err := p.Db.
		Select(models.RuntimeEnvAttachedCredentialColumns...).
		From(models.RuntimeEnvAttachedCredentialTableName).
		Where(db.Eq(RuntimeEnvCredentialIdColumn, runtimeEnvCredentialIds)).
		Load(&runtimeEnvAttachedCredentials)
	if err != nil {
		return nil, err
	}
	return runtimeEnvAttachedCredentials, nil
}

func (p *Server) getRuntimeEnv(runtimeEnvId string) (*models.RuntimeEnv, error) {
	runtimeEnv := &models.RuntimeEnv{}
	err := p.Db.
		Select(models.RuntimeEnvColumns...).
		From(models.RuntimeEnvTableName).
		Where(db.Eq(RuntimeEnvIdColumn, runtimeEnvId)).
		LoadOne(runtimeEnv)
	if err != nil {
		return nil, err
	}
	return runtimeEnv, nil
}

func (p *Server) getRuntimeEnvLabelsByEnvId(runtimeEnvId ...string) ([]*models.RuntimeEnvLabel, error) {
	var runtimeEnvLabels []*models.RuntimeEnvLabel
	query := p.Db.
		Select(models.RuntimeEnvLabelColumns...).
		From(models.RuntimeEnvLabelTableName).
		Where(db.Eq(RuntimeEnvIdColumn, runtimeEnvId))

	_, err := query.Load(&runtimeEnvLabels)
	if err != nil {
		return nil, err
	}
	return runtimeEnvLabels, nil
}

func (p *Server) getRuntimeEnvsByFilterCondition(filterCondition dbr.Builder, limit, offset uint64) (
	runtimeEnvs []*models.RuntimeEnv, count uint32, err error) {
	query := p.Db.
		Select(models.RuntimeEnvColumns...).
		From(models.RuntimeEnvTableName).
		Offset(offset).
		Limit(limit).Where(filterCondition)
	_, err = query.Load(&runtimeEnvs)
	if err != nil {
		return nil, 0, err
	}
	count, err = query.Count()
	if err != nil {
		return nil, 0, err
	}
	return runtimeEnvs, count, nil
}

func (p *Server) getRuntimeEnvCredentialsByFilterCondition(filterCondition dbr.Builder, limit, offset uint64) (
	runtimeEnvCredentials []*models.RuntimeEnvCredential, count uint32, err error) {
	query := p.Db.
		Select(models.RuntimeEnvCredentialColumns...).
		From(models.RuntimeEnvCredentialTableName).
		Offset(offset).
		Limit(limit).
		Where(filterCondition)
	_, err = query.Load(&runtimeEnvCredentials)
	if err != nil {
		return nil, 0, err
	}
	count, err = query.Count()
	if err != nil {
		return nil, 0, err
	}
	return runtimeEnvCredentials, count, err
}

func (p *Server) insertRuntimeEnvLabels(runtimeEnvId string, labelMap map[string]string) error {
	insertQuery := p.Db.
		InsertInto(models.RuntimeEnvLabelTableName).
		Columns(models.RuntimeEnvLabelColumns...)
	for labelKey, labelValue := range labelMap {
		insertQuery = insertQuery.Record(
			models.NewRuntimeEnvLabel(runtimeEnvId, labelKey, labelValue),
		)
	}
	_, err := insertQuery.Exec()
	if err != nil {
		return err
	}
	return nil
}

func (p *Server) deleteRuntimeEnvLabels(runtimeEnvId string, labelMap map[string]string) error {
	var conditions []dbr.Builder
	for labelKey, labelValue := range labelMap {
		conditions = append(conditions, BuildDeleteLabelFilterCondition(runtimeEnvId, labelKey, labelValue))
	}
	_, err := p.Db.
		DeleteFrom(models.RuntimeEnvLabelTableName).
		Where(db.Or(conditions...)).
		Exec()
	if err != nil {
		return err
	}
	return nil
}

func (p *Server) getRuntimeEnvIdsBySelectorMap(selectorMap map[string][]string) ([]string, error) {
	var runtimeEnvIds []string
	sql := BuildSelectRuntimeEnvIdBySelectorSql(selectorMap)
	_, err := p.Db.
		SelectBySql(sql).
		Load(&runtimeEnvIds)
	if err != nil {
		return nil, err
	}
	return runtimeEnvIds, nil
}

func (p *Server) getAttachmentByRuntimeEnvId(runtimeEnvId string) (*models.RuntimeEnvAttachedCredential, error) {
	runtimeEnvAttachedCredential := &models.RuntimeEnvAttachedCredential{}
	err := p.Db.
		Select(models.RuntimeEnvAttachedCredentialColumns...).
		From(models.RuntimeEnvAttachedCredentialTableName).
		Where(db.Eq(RuntimeEnvIdColumn, runtimeEnvId)).
		LoadOne(runtimeEnvAttachedCredential)
	if err != nil {
		return nil, err
	}
	return runtimeEnvAttachedCredential, nil
}

func (p *Server) getAttachmentCountByRuntimeEnvCredentialId(runtimeEnvCredentialId string) (count uint32, err error) {
	count, err = p.Db.
		Select("*").
		From(models.RuntimeEnvAttachedCredentialTableName).
		Where(db.Eq(RuntimeEnvCredentialIdColumn, runtimeEnvCredentialId)).
		Count()
	return count, err
}

func (p *Server) getAttachmentCountByRuntimeEnvId(runtimeEnvId string) (count uint32, err error) {
	count, err = p.Db.
		Select("*").
		From(models.RuntimeEnvAttachedCredentialTableName).
		Where(db.Eq(RuntimeEnvIdColumn, runtimeEnvId)).
		Count()
	return count, err
}

func (p *Server) detachCredentialFromRuntimeEnv(runtimeEnvId, runtimeEnvCredentialId string) error {
	_, err := p.Db.
		DeleteFrom(models.RuntimeEnvAttachedCredentialTableName).
		Where(db.And(db.Eq(RuntimeEnvIdColumn, runtimeEnvId), db.Eq(RuntimeEnvCredentialIdColumn, runtimeEnvCredentialId))).
		Exec()
	return err
}

func (p *Server) deleteRuntimeEnvCredential(runtimeEnvCredentialId string) error {
	_, err := p.Db.
		Update(models.RuntimeEnvCredentialTableName).
		Set(StatusColumn, constants.StatusDeleted).
		Set(StatusTimeColumn, time.Now()).
		Where(db.Eq(RuntimeEnvCredentialIdColumn, runtimeEnvCredentialId)).
		Exec()
	return err
}

func (p *Server) deleteRuntimeEnv(runtimeEnvId string) error {
	_, err := p.Db.
		Update(models.RuntimeEnvTableName).
		Set(StatusColumn, constants.StatusDeleted).
		Set(StatusTimeColumn, time.Now()).
		Where(db.Eq(RuntimeEnvIdColumn, runtimeEnvId)).
		Exec()
	return err
}

func (p *Server) insertRuntimeEnv(runtimeEnv models.RuntimeEnv) error {
	_, err := p.Db.
		InsertInto(models.RuntimeEnvTableName).
		Columns(models.RuntimeEnvColumns...).
		Record(runtimeEnv).
		Exec()
	return err
}

func (p *Server) insertRuntimeEnvAttachedCredential(runtimeEnvAttachedCredential models.RuntimeEnvAttachedCredential) error {
	_, err := p.Db.
		InsertInto(models.RuntimeEnvAttachedCredentialTableName).
		Columns(models.RuntimeEnvAttachedCredentialColumns...).
		Record(runtimeEnvAttachedCredential).
		Exec()
	return err
}

func (p *Server) insertRuntimeEnvCredential(runtimeEnvCredential models.RuntimeEnvCredential) error {
	_, err := p.Db.
		InsertInto(models.RuntimeEnvCredentialTableName).
		Columns(models.RuntimeEnvCredentialColumns...).
		Record(runtimeEnvCredential).
		Exec()
	return err
}

func (p *Server) updateRuntimeEnvByMap(runtimeEnvId string, attributes map[string]interface{}) error {
	_, err := p.Db.
		Update(models.RuntimeEnvTableName).
		SetMap(attributes).
		Where(db.Eq(RuntimeEnvIdColumn, runtimeEnvId)).
		Exec()
	return err
}

func (p *Server) updateRuntimeEnvCredentialByMap(runtimeEnvCredentialId string, attributes map[string]interface{}) error {
	_, err := p.Db.
		Update(models.RuntimeEnvCredentialTableName).
		SetMap(attributes).
		Where(db.Eq(RuntimeEnvCredentialIdColumn, runtimeEnvCredentialId)).
		Exec()
	return err
}

func BuildDeleteLabelFilterCondition(runtimeEnvId, labelKey, labelValue string) dbr.Builder {
	var conditions []dbr.Builder
	conditions = append(conditions, db.Eq(RuntimeEnvIdColumn, runtimeEnvId))
	conditions = append(conditions, db.Eq(LabelKeyColumn, labelKey))
	conditions = append(conditions, db.Eq(LabelValueColumn, labelValue))
	return db.And(conditions...)
}

func BuildSelectRuntimeEnvIdBySelectorSql(selectorMap map[string][]string) string {
	i := 0
	var selectConditionArray []string
	sqlBody := ""
	const baseTableName = "t"
	for labelKey, labelValues := range selectorMap {
		tableAliasName := baseTableName + strconv.Itoa(i+1)
		if i == 0 {
			sqlBody = BuildSelectOneColumnString(models.RuntimeEnvLabelTableName, tableAliasName, RuntimeEnvIdColumn)
		} else {
			sqlBody += BuildInnerJoinTableString("t1", models.RuntimeEnvLabelTableName,
				tableAliasName, RuntimeEnvIdColumn)
		}

		labelKeyConditionString := "(" + BuildSqlEqString(tableAliasName, RuntimeEnvLabelKeyColumn, labelKey) + ")"

		var labelValueConditionArray []string
		for _, labelValue := range labelValues {
			labelValueConditionArray = append(labelValueConditionArray,
				BuildSqlEqString(tableAliasName, RuntimeEnvLabelValueColumn, labelValue))
		}
		labelValueConditionString := "(" + strings.Join(labelValueConditionArray, " or ") + ")"

		selectConditionArray = append(selectConditionArray, "("+labelKeyConditionString+" and "+labelValueConditionString+")\n")
		i++
	}
	sqlBody += "where \n" + strings.Join(selectConditionArray, "and")
	return sqlBody
}

func BuildSelectOneColumnString(tableName, tableAliasName, column string) string {
	return fmt.Sprintf(
		"select %v.%v \n from %v %v \n",
		tableAliasName, column, models.RuntimeEnvLabelTableName, tableAliasName)
}

func BuildInnerJoinTableString(baseTableName, tabelName, tableAliasName, column string) string {
	return fmt.Sprintf("inner join %v %v \n on %v.%v=%v.%v \n",
		tabelName, tableAliasName, baseTableName, column,
		tableAliasName, column)
}

func BuildSqlEqString(tableName, column, value string) string {
	return tableName + "." + column + "=" + "'" + value + "'"
}
