// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package jsontool

import (
	"encoding/json"

	"openpitrix.io/openpitrix/pkg/logger"
)

func Encode(o interface{}) ([]byte, error) {
	return json.Marshal(o)
}

func Decode(y []byte, o interface{}) error {
	return json.Unmarshal(y, o)
}

func ToString(o interface{}) string {
	b, err := Encode(o)
	if err != nil {
		logger.Errorf("Failed to encode [%+v], error: %+v", o, err)
		return ""
	}
	return string(b)
}
