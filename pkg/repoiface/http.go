// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package repoiface

import (
	"context"
	"io/ioutil"
	"net/http"
	neturl "net/url"
	"strings"
)

type HttpInterface struct {
	url *neturl.URL
}

func NewHttpInterface(ctx context.Context, url *neturl.URL) (*HttpInterface, error) {
	return &HttpInterface{
		url: url,
	}, nil
}

func (i *HttpInterface) ReadFile(ctx context.Context, filename string) ([]byte, error) {
	u := strings.TrimSuffix(i.url.String(), "/") + "/" + filename

	resp, err := http.Get(u)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (i *HttpInterface) WriteFile(ctx context.Context, filename string, data []byte) error {
	return nil
}
