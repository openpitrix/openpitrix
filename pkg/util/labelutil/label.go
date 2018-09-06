// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

//go:generate go run gen_helper.go
//go:generate go fmt

package labelutil

import (
	"fmt"
	"net/url"
)

func Parse(labelString string) (map[string][]string, error) {
	m, err := url.ParseQuery(labelString)
	if err != nil {
		return nil, err
	}
	for mKey, mValue := range m {
		for _, v := range mValue {
			if v == "" {
				return nil, fmt.Errorf("bad key [%s] had empty value", mKey)
			}
		}
	}
	return m, nil
}
