// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
)

// OpenpitrixDeleteRepoLabelResponse openpitrix delete repo label response
// swagger:model openpitrixDeleteRepoLabelResponse
type OpenpitrixDeleteRepoLabelResponse struct {

	// repo label
	RepoLabel *OpenpitrixRepoLabel `json:"repo_label,omitempty"`
}

// Validate validates this openpitrix delete repo label response
func (m *OpenpitrixDeleteRepoLabelResponse) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateRepoLabel(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *OpenpitrixDeleteRepoLabelResponse) validateRepoLabel(formats strfmt.Registry) error {

	if swag.IsZero(m.RepoLabel) { // not required
		return nil
	}

	if m.RepoLabel != nil {

		if err := m.RepoLabel.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("repo_label")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *OpenpitrixDeleteRepoLabelResponse) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *OpenpitrixDeleteRepoLabelResponse) UnmarshalBinary(b []byte) error {
	var res OpenpitrixDeleteRepoLabelResponse
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
