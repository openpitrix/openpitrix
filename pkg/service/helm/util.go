// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package helm

import (
	"bytes"
	"fmt"
	"strings"

	"openpitrix.io/openpitrix/pkg/util/jsonutil"
	"openpitrix.io/openpitrix/pkg/util/yamlutil"
)

func ConvertJsonToYaml(data []byte) ([]byte, error) {
	var v map[string]interface{}
	err := jsonutil.Decode(data, &v)
	if err != nil {
		return nil, err
	}

	rawVals, err := yamlutil.Encode(v)
	if err != nil {
		return nil, err
	}
	return rawVals, nil
}

func GetLabelString(m map[string]string) string {
	b := new(bytes.Buffer)
	for k, v := range m {
		fmt.Fprintf(b, "%s=%s ", k, v)
	}
	return b.String()
}

func GetStringFromValues(vals map[string]interface{}, key string) (string, bool) {
	v, ok := vals[key]
	if !ok {
		return "", false
	}
	s, ok := v.(string)
	if !ok {
		return "", false
	}
	return s, true
}

func isConnectionError(err error) bool {
	if err == nil {
		return false
	}
	return strings.HasPrefix(err.Error(), "connection error:")
}

func mergeValues(dest map[string]interface{}, src map[string]interface{}) map[string]interface{} {
	for k, v := range src {
		// If the key doesn't exist already, then just set the key to that value
		if _, exists := dest[k]; !exists {
			dest[k] = v
			continue
		}
		nextMap, ok := v.(map[string]interface{})
		// If it isn't another map, overwrite the value
		if !ok {
			dest[k] = v
			continue
		}
		// Edge case: If the key exists in the destination, but isn't a map
		destMap, isMap := dest[k].(map[string]interface{})
		// If the source map has a map for this key, prefer it
		if !isMap {
			dest[k] = v
			continue
		}
		// If we got to this point, it is a map in both, so merge them
		dest[k] = mergeValues(destMap, nextMap)
	}
	return dest
}
