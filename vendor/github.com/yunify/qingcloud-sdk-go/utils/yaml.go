// +-------------------------------------------------------------------------
// | Copyright (C) 2016 Yunify, Inc.
// +-------------------------------------------------------------------------
// | Licensed under the Apache License, Version 2.0 (the "License");
// | you may not use this work except in compliance with the License.
// | You may obtain a copy of the License in the LICENSE file, or at:
// |
// | http://www.apache.org/licenses/LICENSE-2.0
// |
// | Unless required by applicable law or agreed to in writing, software
// | distributed under the License is distributed on an "AS IS" BASIS,
// | WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// | See the License for the specific language governing permissions and
// | limitations under the License.
// +-------------------------------------------------------------------------

package utils

import (
	"gopkg.in/yaml.v2"
)

// YAMLEncode encode given interface to yaml byte slice.
func YAMLEncode(source interface{}) ([]byte, error) {
	bytesResult, err := yaml.Marshal(source)
	if err != nil {
		return []byte{}, err
	}

	return bytesResult, nil
}

// YAMLDecode decode given yaml byte slice to corresponding struct.
func YAMLDecode(content []byte, destinations ...interface{}) (interface{}, error) {
	var destination interface{}
	var err error
	if len(destinations) == 1 {
		destination = destinations[0]
		err = yaml.Unmarshal(content, destination)
	} else {
		err = yaml.Unmarshal(content, &destination)
	}

	if err != nil {
		return nil, err
	}
	return destination, err
}
