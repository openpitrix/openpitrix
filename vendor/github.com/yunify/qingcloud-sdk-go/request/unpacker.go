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

package request

import (
	"bytes"
	"fmt"
	"net/http"
	"reflect"

	"github.com/yunify/qingcloud-sdk-go/logger"
	"github.com/yunify/qingcloud-sdk-go/request/data"
	"github.com/yunify/qingcloud-sdk-go/request/errors"
	"github.com/yunify/qingcloud-sdk-go/utils"
)

// Unpacker is the response unpacker.
type Unpacker struct {
	operation *data.Operation

	httpResponse *http.Response
	output       *reflect.Value
}

// UnpackHTTPRequest unpack the http response with an operation, http response and an output.
func (u *Unpacker) UnpackHTTPRequest(o *data.Operation, r *http.Response, x *reflect.Value) error {
	u.operation = o
	u.httpResponse = r
	u.output = x

	err := u.parseResponse()
	if err != nil {
		return err
	}

	err = u.parseError()
	if err != nil {
		return err
	}

	return nil
}

func (u *Unpacker) parseResponse() error {
	if u.httpResponse.StatusCode == 200 {
		if u.httpResponse.Header.Get("Content-Type") == "application/json" {
			buffer := &bytes.Buffer{}
			buffer.ReadFrom(u.httpResponse.Body)
			u.httpResponse.Body.Close()

			logger.Info(fmt.Sprintf(
				"Response json string: [%d] %s",
				utils.StringToUnixInt(u.httpResponse.Header.Get("Date"), "RFC 822"),
				string(buffer.Bytes())))

			_, err := utils.JSONDecode(buffer.Bytes(), u.output.Interface())
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (u *Unpacker) parseError() error {
	retCodeValue := u.output.Elem().FieldByName("RetCode")
	messageValue := u.output.Elem().FieldByName("Message")

	if retCodeValue.IsValid() && retCodeValue.Type().String() == "*int" &&
		messageValue.IsValid() && messageValue.Type().String() == "*string" &&
		retCodeValue.Elem().Int() != 0 {

		return &errors.QingCloudError{
			RetCode: int(retCodeValue.Elem().Int()),
			Message: messageValue.Elem().String(),
		}
	}

	return nil
}
