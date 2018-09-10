// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package opapp

import (
	"fmt"
	"strings"

	"github.com/xeipuuv/gojsonschema"

	"openpitrix.io/openpitrix/pkg/util/jsonutil"
)

var schemaLoader = gojsonschema.NewStringLoader(ClusterSchema)

func getError(result *gojsonschema.Result) error {
	var errs []string
	for _, desc := range result.Errors() {
		errs = append(errs, desc.String())
	}
	return fmt.Errorf(strings.Join(errs, "\n\t"))
}

func (c ClusterConf) Validate() error {
	documentLoader := gojsonschema.NewStringLoader(c.RenderJson)
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return err
	}
	if !result.Valid() {
		return getError(result)
	}
	return nil
}

func ValidateClusterConfTmpl(clusterTmpl *ClusterConfTemplate, input jsonutil.Json) error {
	cluster, err := clusterTmpl.Render(input)
	if err != nil {
		return err
	}
	return cluster.Validate()
}
