// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package opapp

import "fmt"

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

func (m *Metadata) GetAppVersion() string {
	if m != nil {
		return m.AppVersion
	}
	return ""
}

func (m *Metadata) GetApiVersion() string {
	if m != nil {
		return m.ApiVersion
	}
	return ""
}

func (m *Metadata) GetIcon() string {
	if m != nil {
		return m.Icon
	}
	return ""
}
func (m *Metadata) GetMaintainers() []*Maintainer {
	if m != nil {
		return m.Maintainers
	}
	return nil
}

func (m *Metadata) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Metadata) GetHome() string {
	if m != nil {
		return m.Home
	}
	return ""
}

func (m *Metadata) GetSources() []string {
	if m != nil {
		return m.Sources
	}
	return nil
}

func (m *Metadata) GetKeywords() []string {
	if m != nil {
		return m.Keywords
	}
	return nil
}

func (m *Metadata) GetVersion() string {
	if m != nil {
		return m.Version
	}
	return ""
}

func (m *Metadata) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *Metadata) GetScreenshots() []string {
	if m != nil {
		return m.Screenshots
	}
	return []string{}
}

func (m *Metadata) GetPackageName() string {
	return fmt.Sprintf("%s-%s.tgz", m.Name, m.Version)
}
