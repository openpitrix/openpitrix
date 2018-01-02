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
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/yunify/qingcloud-sdk-go/logger"
	"github.com/yunify/qingcloud-sdk-go/request/data"
	"github.com/yunify/qingcloud-sdk-go/utils"
)

// Builder is the request builder for QingCloud service.
type Builder struct {
	parsedURL        string
	parsedProperties *map[string]string
	parsedParams     *map[string]string

	operation *data.Operation
	input     *reflect.Value
}

// BuildHTTPRequest builds http request with an operation and an input.
func (b *Builder) BuildHTTPRequest(o *data.Operation, i *reflect.Value) (*http.Request, error) {
	b.operation = o
	b.input = i

	err := b.parse()
	if err != nil {
		return nil, err
	}

	return b.build()
}

func (b *Builder) build() (*http.Request, error) {
	httpRequest, err := http.NewRequest(b.operation.RequestMethod, b.parsedURL, nil)
	if err != nil {
		return nil, err
	}

	logger.Info(fmt.Sprintf(
		"Built QingCloud request: [%d] %s",
		utils.StringToUnixInt(httpRequest.Header.Get("Date"), "RFC 822"),
		httpRequest.URL.String()))

	return httpRequest, nil
}

func (b *Builder) parse() error {

	err := b.parseRequestProperties()
	if err != nil {
		return err
	}
	err = b.parseRequestParams()
	if err != nil {
		return err
	}
	err = b.parseRequestURL()
	if err != nil {
		return err
	}

	return nil
}

func (b *Builder) parseRequestProperties() error {
	propertiesMap := map[string]string{}

	fields := reflect.ValueOf(b.operation.Properties).Elem()
	for i := 0; i < fields.NumField(); i++ {
		switch value := fields.Field(i).Interface().(type) {
		case *string:
			if value != nil {
				propertiesMap[fields.Type().Field(i).Tag.Get("name")] = *value
			}
		case *int:
			if value != nil {
				propertiesMap[fields.Type().Field(i).Tag.Get("name")] = strconv.Itoa(int(*value))
			}
		}
	}

	b.parsedProperties = &propertiesMap

	return nil
}

func (b *Builder) parseRequestParams() error {
	var requestParams map[string]string

	if b.parsedParams != nil {
		requestParams = *b.parsedParams
	} else {
		requestParams = map[string]string{}
	}

	b.parsedParams = &requestParams

	requestParams["action"] = b.operation.APIName

	if !b.input.Elem().IsValid() {
		return nil
	}

	for i := 0; i < b.input.Elem().NumField(); i++ {
		tagName := b.input.Elem().Type().Field(i).Tag.Get("name")
		tagLocation := b.input.Elem().Type().Field(i).Tag.Get("location")
		tagDefault := b.input.Elem().Type().Field(i).Tag.Get("default")
		if tagName != "" && tagLocation != "" && requestParams != nil {
			switch value := b.input.Elem().Field(i).Interface().(type) {
			case *string:
				if tagDefault != "" {
					requestParams[tagName] = tagDefault
				}
				if value != nil {
					requestParams[tagName] = *value
				}
			case *int:
				if tagDefault != "" {
					requestParams[tagName] = tagDefault
				}
				if value != nil {
					requestParams[tagName] = strconv.Itoa(int(*value))
				}
			case *bool:
			case *time.Time:
				if tagDefault != "" {
					requestParams[tagName] = tagDefault
				}
				if value != nil {
					format := b.input.Elem().Type().Field(i).Tag.Get("format")
					requestParams[tagName] = utils.TimeToString(*value, format)
				}
			case []*string:
				for index, item := range value {
					key := tagName + "." + strconv.Itoa(index+1)
					if tagDefault != "" {
						requestParams[tagName] = tagDefault
					}
					if item != nil {
						requestParams[key] = *item
					}
				}
			case []*int:
				for index, item := range value {
					key := tagName + "." + strconv.Itoa(index+1)
					if tagDefault != "" {
						requestParams[tagName] = tagDefault
					}
					if item != nil {
						requestParams[key] = strconv.Itoa(int(*item))
					}
				}
			default:
				if value != nil {
					value = value.(interface{})
					typeName := reflect.TypeOf(value.(interface{})).String()

					if strings.HasPrefix(typeName, "[]*") {
						valueArray := reflect.ValueOf(value)

						for i := 0; i < valueArray.Len(); i++ {
							item := valueArray.Index(i).Elem()

							for j := 0; j < item.NumField(); j++ {
								fieldTagName := item.Type().Field(j).Tag.Get("name")
								tagKey := tagName + "." + strconv.Itoa(i+1) + "." + fieldTagName

								switch fieldValue := item.Field(j).Interface().(type) {
								case *int:
									if fieldValue != nil {
										requestParams[tagKey] = strconv.Itoa(int(*fieldValue))
									}
								case *string:
									if fieldValue != nil {
										requestParams[tagKey] = *fieldValue
									}
								}
							}
						}
					}
				}
			}
		}
	}

	return nil
}

func (b *Builder) parseRequestURL() error {
	conf := b.operation.Config

	endpoint := conf.Protocol + "://" + conf.Host + ":" + strconv.Itoa(conf.Port)
	requestURI := regexp.MustCompile(`/+`).ReplaceAllString(conf.URI, "/")

	b.parsedURL = endpoint + requestURI

	if b.parsedParams != nil {
		zone := (*b.parsedProperties)["zone"]
		if zone != "" {
			(*b.parsedParams)["zone"] = zone
		}
		paramsParts := []string{}
		for key, value := range *b.parsedParams {
			paramsParts = append(paramsParts, fmt.Sprintf("%s=%s", key, url.QueryEscape(value)))
		}

		joined := strings.Join(paramsParts, "&")
		if joined != "" {
			b.parsedURL += "?" + joined
		}
	}

	return nil
}
