// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package gerr

import "fmt"

type ErrorMessage struct {
	Name string
	En   string
}

func (em ErrorMessage) Message(locale string, a ...interface{}) string {
	format := ""
	switch locale {
	case EN:
		format = em.En
	}
	return fmt.Sprintf(format, a...)
}

var (
	ErrorCreateResourceFailed = ErrorMessage{
		Name: "create_resource_failed",
		En:   "create resource failed",
	}
	ErrorDeleteResourceFailed = ErrorMessage{
		Name: "delete_resource_failed",
		En:   "delete resource failed",
	}
	ErrorDescribeResourcesFailed = ErrorMessage{
		Name: "describe_resources_failed",
		En:   "describe resource failed",
	}
	ErrorModifyResourcesFailed = ErrorMessage{
		Name: "modify_resources_failed",
		En:   "modify resource failed",
	}
	ErrorResourcesNotFound = ErrorMessage{
		Name: "resource_not_found",
		En:   "resource [%s] not found",
	}
	ErrorInternalError = ErrorMessage{
		Name: "internal_error",
		En:   "internal error",
	}
	ErrorMissingParameter = ErrorMessage{
		Name: "missing_parameter",
		En:   "missing parameter [%s]",
	}
	ErrorValidateFailed = ErrorMessage{
		Name: "validate_failed",
		En:   "validate failed",
	}
	ErrorParameterParseFailed = ErrorMessage{
		Name: "parameter_parse_failed",
		En:   "parameter [%s] parse failed",
	}
	ErrorResourceAlreadyDeleted = ErrorMessage{
		Name: "resource_already_deleted",
		En:   "resource [%s] has already been deleted",
	}
	ErrorIllegalParameterLength = ErrorMessage{
		Name: "illegal_parameter_length",
		En:   "illegal parameter [%s] length",
	}
	ErrorUnsupportedParameterValue = ErrorMessage{
		Name: "unsupported_parameter_value",
		En:   "unsupported parameter [%s] value [%s]",
	}
	ErrorIllegalUrlFormat = ErrorMessage{
		Name: "illegal_url_format",
		En:   "illegal URL format [%s]",
	}
	ErrorIllegalLabelFormat = ErrorMessage{
		Name: "illegal_label_format",
		En:   "illegal label format",
	}
)
