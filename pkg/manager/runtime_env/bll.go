package runtime_env

import (
	"fmt"

	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/logger"
)

func (p *Server) getRuntimeEnvPbWithLabel(runtimeEnvId string) (*pb.RuntimeEnv, error) {
	runtimeEnv, err := p.getRuntimeEnv(runtimeEnvId)
	if err != nil {
		logger.Errorf("failed to get runtime_env [%v] ", err)
		return nil, fmt.Errorf("failed to get runtime_env [%v] ", err)
	}
	pbRuntimeEnv := models.RuntimeEnvToPb(runtimeEnv)
	runtimeEnvLabels, err := p.getRuntimeEnvLabelsByEnvId(runtimeEnvId)
	if err != nil {
		logger.Errorf("failed to get runtime_env_label [%v] ", err)
		return nil, fmt.Errorf("failed to get runtime_env label [%v] ", err)
	}
	pbRuntimeEnv.Labels = models.RuntimeEnvLabelsToPbs(runtimeEnvLabels)

	return pbRuntimeEnv, nil
}

func (p *Server) createRuntimeEnv(name, description, url, userId string) (runtimeEnvId string, err error) {
	newRuntimeEnv := models.NewRuntimeEnv(name, description, url, userId)
	err = p.insertRuntimeEnv(*newRuntimeEnv)
	if err != nil {
		logger.Errorf("failed to insert runtime_env [%v] ", err)
		return "", fmt.Errorf("failed to insert runtime_env [%v] ", err)
	}
	return newRuntimeEnv.RuntimeEnvId, err
}

func (p *Server) createRuntimeEnvLabels(runtimeEnvId, labelString string) error {
	labelMap, err := LabelStringToMap(labelString)
	if err != nil {
		logger.Error(err)
		return err
	}
	err = p.insertRuntimeEnvLabels(runtimeEnvId, labelMap)
	return err
}
