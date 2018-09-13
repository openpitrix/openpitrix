// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package runtime

import (
	"net/url"

	"github.com/ghodss/yaml"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/plugins"
)

func SelectorStringToMap(selectorString string) (map[string][]string, error) {
	selectorMap, err := url.ParseQuery(selectorString)
	if err != nil {
		return nil, err
	}
	return selectorMap, nil
}

func CredentialStringToJsonString(provider, content string) string {
	if plugins.IsVmbasedProviders(provider) {
		return content
	} else if constants.ProviderKubernetes == provider {
		content, err := yaml.YAMLToJSON([]byte(content))
		if err != nil {
			panic(err)
		}
		return string(content)
	}
	panic("unsupport provider")
}

func CredentialJsonStringToString(provider, content string) string {
	if plugins.IsVmbasedProviders(provider) {
		return content
	} else if constants.ProviderKubernetes == provider {
		content, err := yaml.JSONToYAML([]byte(content))
		if err != nil {
			panic(err)
		}
		return string(content)
	}
	panic("unsupport provider")
}
