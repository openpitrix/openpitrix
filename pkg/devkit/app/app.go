// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package app

// BufferedFile represents an archive file buffered for later processing.
type BufferedFile struct {
	Name string
	Data []byte
}

type App struct {
	Metadata *Metadata `json:"metadata,omitempty"`

	ConfigTemplate *ConfigTemplate `json:"config,omitempty"`

	ClusterTemplate *ClusterTemplate `json:"cluster_template,omitempty"`

	Files []BufferedFile `json:"files,omitempty"`
}
