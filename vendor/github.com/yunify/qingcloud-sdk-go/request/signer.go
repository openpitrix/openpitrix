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
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/yunify/qingcloud-sdk-go/logger"
	"github.com/yunify/qingcloud-sdk-go/utils"
)

// Signer is the http request signer for IaaS service.
type Signer struct {
	AccessKeyID     string
	SecretAccessKey string

	BuiltURL string
}

// WriteSignature calculates signature and write it to http request.
func (is *Signer) WriteSignature(request *http.Request) error {
	_, err := is.BuildSignature(request)
	if err != nil {
		return err
	}

	newRequest, err := http.NewRequest(request.Method,
		request.URL.Scheme+"://"+request.URL.Host+is.BuiltURL, nil)
	if err != nil {
		return err
	}
	request.URL = newRequest.URL

	logger.Info(fmt.Sprintf(
		"Signed QingCloud request: [%d] %s",
		utils.StringToUnixInt(request.Header.Get("Date"), "RFC 822"),
		request.URL.String()))

	return nil
}

// BuildSignature calculates the signature string.
func (is *Signer) BuildSignature(request *http.Request) (string, error) {
	stringToSign, err := is.BuildStringToSign(request)
	if err != nil {
		return "", err
	}

	h := hmac.New(sha256.New, []byte(is.SecretAccessKey))
	h.Write([]byte(stringToSign))

	signature := strings.TrimSpace(base64.StdEncoding.EncodeToString(h.Sum(nil)))
	signature = strings.Replace(signature, " ", "+", -1)
	signature = url.QueryEscape(signature)

	logger.Debug(fmt.Sprintf(
		"QingCloud signature: [%d] %s",
		utils.StringToUnixInt(request.Header.Get("Date"), "RFC 822"),
		signature))

	is.BuiltURL += "&signature=" + signature

	return signature, nil
}

// BuildStringToSign build the string to sign.
func (is *Signer) BuildStringToSign(request *http.Request) (string, error) {
	query := request.URL.Query()

	query.Set("access_key_id", is.AccessKeyID)
	query.Set("signature_method", "HmacSHA256")
	query.Set("signature_version", "1")

	var timeValue time.Time
	if request.Header.Get("Date") != "" {
		var err error
		timeValue, err = utils.StringToTime(request.Header.Get("Date"), "RFC 822")
		if err != nil {
			return "", err
		}
	} else {
		timeValue = time.Now()
	}
	query.Set("time_stamp", utils.TimeToString(timeValue, "ISO 8601"))

	keys := []string{}
	for key := range query {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	parts := []string{}
	for _, key := range keys {
		values := query[key]
		if len(values) > 0 {
			if values[0] != "" {
				value := strings.TrimSpace(strings.Join(values, ""))
				value = url.QueryEscape(value)
				value = strings.Replace(value, "+", "%20", -1)
				parts = append(parts, key+"="+value)
			} else {
				parts = append(parts, key)
			}
		} else {
			parts = append(parts, key)
		}
	}

	urlParams := strings.Join(parts, "&")

	stringToSign := request.Method + "\n" + request.URL.Path + "\n" + urlParams

	logger.Debug(fmt.Sprintf(
		"QingCloud string to sign: [%d] %s",
		utils.StringToUnixInt(request.Header.Get("Date"), "RFC 822"),
		stringToSign))

	is.BuiltURL = request.URL.Path + "?" + urlParams

	return stringToSign, nil
}
