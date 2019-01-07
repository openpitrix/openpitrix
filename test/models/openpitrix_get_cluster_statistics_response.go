// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
)

// OpenpitrixGetClusterStatisticsResponse openpitrix get cluster statistics response
// swagger:model openpitrixGetClusterStatisticsResponse
type OpenpitrixGetClusterStatisticsResponse struct {

	// cluster count
	ClusterCount int64 `json:"cluster_count,omitempty"`

	// cluster create time range -> cluster count, max length is 14
	LastTwoWeekCreated map[string]int64 `json:"last_two_week_created,omitempty"`

	// runtime count
	RuntimeCount int64 `json:"runtime_count,omitempty"`

	// app id -> cluster count, max length is 10
	TopTenApps map[string]int64 `json:"top_ten_apps,omitempty"`

	// runtime id -> cluster count, max length is 10
	TopTenRuntimes map[string]int64 `json:"top_ten_runtimes,omitempty"`
}

// Validate validates this openpitrix get cluster statistics response
func (m *OpenpitrixGetClusterStatisticsResponse) Validate(formats strfmt.Registry) error {
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// MarshalBinary interface implementation
func (m *OpenpitrixGetClusterStatisticsResponse) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *OpenpitrixGetClusterStatisticsResponse) UnmarshalBinary(b []byte) error {
	var res OpenpitrixGetClusterStatisticsResponse
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
