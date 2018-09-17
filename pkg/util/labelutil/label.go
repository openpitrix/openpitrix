// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

//go:generate go run gen_helper.go
//go:generate go fmt

package labelutil

import (
	"bytes"
	"fmt"
	"net/url"
	"strings"
)

func Parse(labelString string) (Query, error) {
	m, err := parseQuery(labelString)
	if err != nil {
		return nil, err
	}
	for _, val := range m {
		if val.V == "" {
			return nil, fmt.Errorf("bad key [%s] had empty value", val.K)
		}
	}
	return m, nil
}

type Q struct {
	K string
	V string
}

type Query []Q

func (query Query) Append(k string, v string) Query {
	q := Q{
		K: k,
		V: v,
	}
	return append(query, q)
}

// Ref: net/url/url.go.Values.Encode
func (query Query) ToString() string {
	if len(query) == 0 {
		return ""
	}
	var buf bytes.Buffer
	for _, q := range query {
		prefix := url.QueryEscape(q.K) + "="
		if buf.Len() > 0 {
			buf.WriteByte('&')
		}
		buf.WriteString(prefix)
		buf.WriteString(url.QueryEscape(q.V))
	}
	return buf.String()
}

// Ref: net/url/url.go.ParseQuery
func parseQuery(query string) (Query, error) {
	var m Query
	var err error
	for query != "" {
		key := query
		if i := strings.IndexAny(key, "&;"); i >= 0 {
			key, query = key[:i], key[i+1:]
		} else {
			query = ""
		}
		if key == "" {
			continue
		}
		value := ""
		if i := strings.Index(key, "="); i >= 0 {
			key, value = key[:i], key[i+1:]
		}
		key, err1 := url.QueryUnescape(key)
		if err1 != nil {
			if err == nil {
				err = err1
			}
			continue
		}
		value, err1 = url.QueryUnescape(value)
		if err1 != nil {
			if err == nil {
				err = err1
			}
			continue
		}
		m = m.Append(key, value)
	}
	return m, err
}
