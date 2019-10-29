// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
)

type Out struct {
	action string
	out    io.Writer
}

type httpBody interface {
	MarshalBinary() ([]byte, error)
}
type httpRequest interface {
	WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error
}

func (o Out) GetBodyWithIndent(body httpBody) ([]byte, error) {
	var buf bytes.Buffer
	b, err := body.MarshalBinary()
	if err != nil {
		return buf.Bytes(), err
	}
	err = json.Indent(&buf, b, "", "  ")
	return buf.Bytes(), err
}

func (o Out) GetParamsWithIndent(request httpRequest) ([]byte, error) {
	var ret []byte
	reg := strfmt.NewFormats()
	r := &clientRequest{}
	err := request.WriteToRequest(r, reg)
	if err != nil {
		return ret, err
	}
	if r.GetQueryParams() != nil {
		return []byte(r.GetQueryParams().Encode()), nil
		//return toJson(r.GetQueryParams())
	}
	return r.GetBody(), nil
}

func (o Out) WriteRequest(request httpRequest) error {
	b, err := o.GetParamsWithIndent(request)
	if err != nil {
		return err
	}

	o.out.Write([]byte(fmt.Sprintf("------ sending request ------\n%s\n", o.action)))
	o.out.Write([]byte("------ params ------\n"))
	o.out.Write(b)
	o.out.Write([]byte("\n"))
	return nil
}

func (o Out) WriteResponse(response httpBody) error {
	b, err := o.GetBodyWithIndent(response)
	if err != nil {
		return err
	}
	o.out.Write([]byte("------ response ------\n"))
	o.out.Write(b)
	o.out.Write([]byte("\n"))
	return nil
}

func toJson(i interface{}) ([]byte, error) {
	return json.MarshalIndent(i, "", "  ")
}

type clientRequest struct {
	query   url.Values
	payload interface{}
}

func (c *clientRequest) GetHeaderParams() http.Header {
	return make(http.Header)
}

func (c *clientRequest) SetHeaderParam(string, ...string) error {
	return nil
}

func (c *clientRequest) SetQueryParam(name string, values ...string) error {
	if c.query == nil {
		c.query = make(url.Values)
	}
	c.query[name] = values
	return nil
}

func (c *clientRequest) SetFormParam(string, ...string) error {
	return nil
}

func (c *clientRequest) SetPathParam(string, string) error {
	return nil
}

func (c *clientRequest) GetQueryParams() url.Values {
	return c.query
}

func (c *clientRequest) SetFileParam(string, ...runtime.NamedReadCloser) error {
	return nil
}

func (c *clientRequest) SetBodyParam(payload interface{}) error {
	c.payload = payload
	return nil
}

func (c *clientRequest) SetTimeout(time.Duration) error {
	return nil
}

func (c *clientRequest) GetMethod() string {
	return ""
}

func (c *clientRequest) GetPath() string {
	return ""
}

func (c *clientRequest) GetBody() []byte {
	if c.payload == nil {
		return nil
	}
	ret, err := toJson(c.payload)
	if err != nil {
		panic(err)
	}
	return ret
}

func (c *clientRequest) GetBodyParam() interface{} {
	return c.payload
}

func (c *clientRequest) GetFileParam() map[string][]runtime.NamedReadCloser {
	return nil
}
