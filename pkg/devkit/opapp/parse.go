package opapp

import (
	"encoding/json"
	"fmt"
)

func DecodePackageJson(data []byte) (*Metadata, error) {
	y := &Metadata{}
	err := json.Unmarshal(data, y)
	if err != nil {
		return nil, fmt.Errorf("failed to decode package.json: %+v", err)
	}
	return y, nil
}

func DecodeConfigJson(data []byte) (*ConfigTemplate, error) {
	y := &ConfigTemplate{}
	err := json.Unmarshal(data, y)
	if err != nil {
		return nil, fmt.Errorf("failed to decode config.json: %+v", err)
	}
	y.Raw = string(data)
	return y, nil
}
