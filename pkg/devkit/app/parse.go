package app

import (
	"encoding/json"
)

func UnmarshalMetadata(data []byte) (*Metadata, error) {
	y := &Metadata{}
	err := json.Unmarshal(data, y)
	if err != nil {
		return nil, err
	}
	return y, nil
}

func UnmarshalConfigTemplate(data []byte) (*ConfigTemplate, error) {
	y := &ConfigTemplate{}
	err := json.Unmarshal(data, y)
	if err != nil {
		return nil, err
	}
	y.Raw = string(data)
	return y, nil
}
