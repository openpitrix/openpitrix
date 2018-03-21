package runtime_env

import (
	"fmt"

	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/models"
)

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

func (p *Server) getRuntimeEnvLabelsByEnvId(runtimeEnvId string) ([]*models.RuntimeEnvLabel, error) {
	var runtimeEnvLabels []*models.RuntimeEnvLabel
	query := p.Db.
		Select(models.RuntimeEnvLabelColumns...).
		From(models.RuntimeEnvLabelTableName).
		Where(db.Eq(RuntimeEnvIdColumn, runtimeEnvId))

	_, err := query.Load(&runtimeEnvLabels)
	if err != nil {
		return nil, fmt.Errorf("get runtime_env_labels error %+v", err)
	}
	return runtimeEnvLabels, nil
}

func (p *Server) insertRuntimeEnv(runtimeEnv models.RuntimeEnv) error {
	_, err := p.Db.
		InsertInto(models.RuntimeEnvTableName).
		Columns(models.RuntimeEnvColumns...).
		Record(runtimeEnv).
		Exec()
	return err
}

func (p *Server) insertRuntimeEnvLabels(runtimeEnvId string, labelMap map[string]string) error {
	for labelKey, labelValue := range labelMap {
		newRuntimeEnvLabels := models.NewRuntimeEnvLabel(runtimeEnvId, labelKey, labelValue)
		_, err := p.Db.
			InsertInto(models.RuntimeEnvLabelTableName).
			Columns(models.RuntimeEnvLabelColumns...).
			Record(newRuntimeEnvLabels).
			Exec()
		if err != nil {
			return err
		}
	}
	return nil
}
