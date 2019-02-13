// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
)

// OpenpitrixAppVersionAudit openpitrix app version audit
// swagger:model openpitrixAppVersionAudit
type OpenpitrixAppVersionAudit struct {

	// app id
	AppID string `json:"app_id,omitempty"`

	// app name
	AppName string `json:"app_name,omitempty"`

	// message
	Message string `json:"message,omitempty"`

	// operator
	Operator string `json:"operator,omitempty"`

	// review id
	ReviewID string `json:"review_id,omitempty"`

	// role
	Role string `json:"role,omitempty"`

	// status
	Status string `json:"status,omitempty"`

	// status time
	StatusTime strfmt.DateTime `json:"status_time,omitempty"`

	// version id
	VersionID string `json:"version_id,omitempty"`

	// version name
	VersionName string `json:"version_name,omitempty"`
}

// Validate validates this openpitrix app version audit
func (m *OpenpitrixAppVersionAudit) Validate(formats strfmt.Registry) error {
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// MarshalBinary interface implementation
func (m *OpenpitrixAppVersionAudit) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *OpenpitrixAppVersionAudit) UnmarshalBinary(b []byte) error {
	var res OpenpitrixAppVersionAudit
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
