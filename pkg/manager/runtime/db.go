// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package runtime

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

func (p *Server) getRuntimeCredentialsByFilterCondition(filterCondition dbr.Builder, limit, offset uint64) (
	runtimeCredentials []*models.RuntimeCredential, count uint32, err error) {
	query := p.Db.
		Select(models.RuntimeCredentialColumns...).
		From(models.RuntimeCredentialTableName).
		Offset(offset).
		Limit(limit).
		Where(filterCondition)
	_, err = query.Load(&runtimeCredentials)
	if err != nil {
		return nil, 0, err
	}
	count, err = query.Count()
	if err != nil {
		return nil, 0, err
	}
	return runtimeCredentials, count, err
}

func (p *Server) insertRuntimeLabels(runtimeId string, labelMap map[string]string) error {
	if labelMap == nil {
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

// deleteRuntimeLabels by runtimeId and labelMap.if labelMap = nil, record are not deleted.
func (p *Server) deleteRuntimeLabels(runtimeId string, labelMap map[string]string) error {
	if labelMap == nil {
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

func (p *Server) getRuntimeIdsBySelectorMap(selectorMap map[string][]string) ([]string, error) {
	var runtimeIds []string
	sql := BuildSelectRuntimeIdBySelectorSql(selectorMap)
	_, err := p.Db.
		SelectBySql(sql).
		Load(&runtimeIds)
	if err != nil {
		return nil, err
	}
	return runtimeIds, nil
}

func (p *Server) deleteRuntimeCredential(runtimeCredentialId string) error {
	_, err := p.Db.
		Update(models.RuntimeCredentialTableName).
		Set(StatusColumn, constants.StatusDeleted).
		Set(StatusTimeColumn, time.Now()).
		Where(db.Eq(RuntimeCredentialIdColumn, runtimeCredentialId)).
		Exec()
	return err
}

func (p *Server) deleteRuntime(runtimeId string) error {
	_, err := p.Db.
		Update(models.RuntimeTableName).
		Set(StatusColumn, constants.StatusDeleted).
		Set(StatusTimeColumn, time.Now()).
		Where(db.Eq(RuntimeIdColumn, runtimeId)).
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

func (p *Server) updateRuntimeCredentialByMap(runtimeCredentialId string, attributes map[string]interface{}) error {
	if attributes == nil {
		return nil
	}
	_, err := p.Db.
		Update(models.RuntimeCredentialTableName).
		SetMap(attributes).
		Where(db.Eq(RuntimeCredentialIdColumn, runtimeCredentialId)).
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

func BuildSelectRuntimeIdBySelectorSql(selectorMap map[string][]string) string {
	i := 0
	var selectConditionArray []string
	sqlBody := ""
	const baseTableName = "t"
	for labelKey, labelValues := range selectorMap {
		tableAliasName := baseTableName + strconv.Itoa(i+1)
		if i == 0 {
			sqlBody = BuildSelectOneColumnString(models.RuntimeLabelTableName, tableAliasName, RuntimeIdColumn)
		} else {
			sqlBody += BuildInnerJoinTableString("t1", models.RuntimeLabelTableName,
				tableAliasName, RuntimeIdColumn)
		}

		labelKeyConditionString := "(" + BuildSqlEqString(tableAliasName, RuntimeLabelKeyColumn, labelKey) + ")"

		var labelValueConditionArray []string
		for _, labelValue := range labelValues {
			labelValueConditionArray = append(labelValueConditionArray,
				BuildSqlEqString(tableAliasName, RuntimeLabelValueColumn, labelValue))
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
		tableAliasName, column, models.RuntimeLabelTableName, tableAliasName)
}

func BuildInnerJoinTableString(baseTableName, tabelName, tableAliasName, column string) string {
	return fmt.Sprintf("inner join %v %v \n on %v.%v=%v.%v \n",
		tabelName, tableAliasName, baseTableName, column,
		tableAliasName, column)
}

func BuildSqlEqString(tableName, column, value string) string {
	return tableName + "." + column + "=" + "'" + value + "'"
}
