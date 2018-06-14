package plugins

import (
	"testing"

	"openpitrix.io/openpitrix/pkg/constants"
)

func TestProviderInterface(t *testing.T) {
	_, err := GetProviderPlugin(constants.ProviderQingCloud, nil)
	if err != nil {
		t.Errorf("Error: %+v", err)
	}
	_, err = GetProviderPlugin(constants.ProviderKubernetes, nil)
	if err != nil {
		t.Errorf("Error: %+v", err)
	}
}
