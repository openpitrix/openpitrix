// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package yaml

import (
	"github.com/ghodss/yaml"
)

// Marshals the object into JSON then converts JSON to YAML and returns the
// YAML.
func Marshal(o interface{}) ([]byte, error) {
	return yaml.Marshal(o)
}

// Converts YAML to JSON then uses JSON to unmarshal into an object.
func Unmarshal(y []byte, o interface{}) error {
	return yaml.Unmarshal(y, o)
}

// Convert JSON to YAML.
func JSONToYAML(j []byte) ([]byte, error) {
	return yaml.JSONToYAML(j)
}

// Convert YAML to JSON. Since JSON is a subset of YAML, passing JSON through
// this method should be a no-op.
//
// Things YAML can do that are not supported by JSON:
// * In YAML you can have binary and null keys in your maps. These are invalid
//   in JSON. (int and float keys are converted to strings.)
// * Binary data in YAML with the !!binary tag is not supported. If you want to
//   use binary data with this library, encode the data as base64 as usual but do
//   not use the !!binary tag in your YAML. This will ensure the original base64
//   encoded data makes it all the way through to the JSON.
func YAMLToJSON(y []byte) ([]byte, error) {
	return yaml.YAMLToJSON(y)
}
