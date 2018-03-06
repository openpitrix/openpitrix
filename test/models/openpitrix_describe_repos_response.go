// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
)

// OpenpitrixDescribeReposResponse openpitrix describe repos response
// swagger:model openpitrixDescribeReposResponse
type OpenpitrixDescribeReposResponse struct {

	// repo set
	RepoSet OpenpitrixDescribeReposResponseRepoSet `json:"repo_set"`

	// total count
	TotalCount int64 `json:"total_count,omitempty"`
}

// Validate validates this openpitrix describe repos response
func (m *OpenpitrixDescribeReposResponse) Validate(formats strfmt.Registry) error {
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// MarshalBinary interface implementation
func (m *OpenpitrixDescribeReposResponse) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *OpenpitrixDescribeReposResponse) UnmarshalBinary(b []byte) error {
	var res OpenpitrixDescribeReposResponse
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
