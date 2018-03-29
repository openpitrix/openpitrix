// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package app

type Metadata struct {
	ApiVersion  string        `json:"api_version,omitempty"`
	Name        string        `json:"name"`
	Version     string        `json:"version"`
	AppVersion  string        `json:"app_version,omitempty"`
	Description string        `json:"description,omitempty"`
	Home        string        `json:"home,omitempty"`
	Icon        string        `json:"icon,omitempty"`
	Maintainers []*Maintainer `json:"maintainers,omitempty"`
	Screenshots []string      `json:"screenshots,omitempty"`
	Keywords    []string      `json:"keywords,omitempty"`
	Sources     []string      `json:"sources,omitempty"`
}

type Maintainer struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
}
