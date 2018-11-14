// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package opapp

import (
	"github.com/golang/protobuf/ptypes/any"
	"github.com/pkg/errors"

	"openpitrix.io/openpitrix/pkg/util/jsonutil"
)

// BufferedFile represents an archive file buffered for later processing.
type BufferedFile struct {
	Name string
	Data []byte
}

type OpApp struct {
	Metadata *Metadata `json:"metadata,omitempty"`

	ConfigTemplate *ConfigTemplate `json:"config,omitempty"`

	ClusterConfTemplate *ClusterConfTemplate `json:"cluster_template,omitempty"`

	Files []*any.Any `json:"files,omitempty"`
}

func (a *OpApp) Validate(config jsonutil.Json) error {
	conf := a.ConfigTemplate.GetRenderedConfig(config)
	err := a.ConfigTemplate.Validate(conf)
	if err != nil {
		return errors.Wrap(err, "validate config.json failed")
	}
	err = ValidateClusterConfTmpl(a.ClusterConfTemplate, conf)
	if err != nil {
		return errors.Wrap(err, "validate cluster.json.tmpl failed")
	}
	return nil
}
