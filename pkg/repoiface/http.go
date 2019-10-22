// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package repoiface

import (
	"context"
	"fmt"
	"io/ioutil"
	neturl "net/url"

	"openpitrix.io/openpitrix/pkg/util/httputil"
	"openpitrix.io/openpitrix/pkg/util/yamlutil"
)

type HttpInterface struct {
	url *neturl.URL
}

func NewHttpInterface(ctx context.Context, url *neturl.URL) (*HttpInterface, error) {
	return &HttpInterface{
		url: url,
	}, nil
}

func (i *HttpInterface) CheckFile(ctx context.Context, filename string) (bool, error) {
	u := URLJoin(i.url.String(), filename)

	resp, err := httputil.HttpGet(u)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || 299 < resp.StatusCode {
		return false, fmt.Errorf("http status code is %d", resp.StatusCode)
	}

	return true, nil
}

func (i *HttpInterface) ReadFile(ctx context.Context, filename string) ([]byte, error) {
	u := URLJoin(i.url.String(), filename)

	resp, err := httputil.HttpGet(u)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf(`looks like "%s" is not a valid chart repository or cannot be reached: Failed to fetch %s : %s`, i.url.String(), u, resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (i *HttpInterface) DeleteFile(ctx context.Context, filename string) error {
	return ErrWriteIsUnsupported
}

func (i *HttpInterface) WriteFile(ctx context.Context, filename string, data []byte) error {
	return ErrWriteIsUnsupported
}

type indexYaml struct {
	ApiVersion string                 `yaml:"apiVersion"`
	Entries    map[string]interface{} `yaml:"entries"`
	Generated  string                 `yaml:"generated"`
}

func (i *HttpInterface) CheckRead(ctx context.Context) error {
	data, err := i.ReadFile(ctx, "index.yaml")
	if err != nil {
		return err
	}
	var y indexYaml
	err = yamlutil.Decode(data, &y)
	return err
}

func (i *HttpInterface) CheckWrite(ctx context.Context) error {
	return ErrWriteIsUnsupported
}
