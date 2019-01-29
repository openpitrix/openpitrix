// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
)

// OpenpitrixCanDoRequest openpitrix can do request
// swagger:model openpitrixCanDoRequest
type OpenpitrixCanDoRequest struct {

	// url
	URL string `json:"url,omitempty"`

	// url method
	URLMethod string `json:"url_method,omitempty"`

	// user id
	UserID string `json:"user_id,omitempty"`
}

// Validate validates this openpitrix can do request
func (m *OpenpitrixCanDoRequest) Validate(formats strfmt.Registry) error {
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// MarshalBinary interface implementation
func (m *OpenpitrixCanDoRequest) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *OpenpitrixCanDoRequest) UnmarshalBinary(b []byte) error {
	var res OpenpitrixCanDoRequest
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}