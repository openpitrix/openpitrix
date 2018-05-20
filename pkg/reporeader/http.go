// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package reporeader

import (
	"io/ioutil"
	"net/http"
	neturl "net/url"
	"strings"
)

type HttpReader struct {
	url *neturl.URL
}

func NewHttpReader(url *neturl.URL) *HttpReader {
	return &HttpReader{
		url: url,
	}
}

func (h *HttpReader) GetIndexYaml() ([]byte, error) {
	u := strings.TrimSuffix(h.url.String(), "/") + "/index.yaml"

	resp, err := http.Get(u)
	if err != nil {
		return nil, ErrGetIndexYamlFailed
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, ErrGetIndexYamlFailed
	}
	return body, nil
}
