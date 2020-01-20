// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
)

// ReversalSubmissionAttributes reversal submission attributes
// swagger:model reversalSubmissionAttributes
type ReversalSubmissionAttributes struct {

	// scheme status code
	SchemeStatusCode string `json:"scheme_status_code,omitempty"`

	// status
	Status ReversalSubmissionStatus `json:"status,omitempty"`

	// status reason
	StatusReason string `json:"status_reason,omitempty"`

	// submission datetime
	// Read Only: true
	SubmissionDatetime strfmt.DateTime `json:"submission_datetime,omitempty"`

	// transaction start datetime
	// Read Only: true
	TransactionStartDatetime strfmt.DateTime `json:"transaction_start_datetime,omitempty"`
}

// Validate validates this reversal submission attributes
func (m *ReversalSubmissionAttributes) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateStatus(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ReversalSubmissionAttributes) validateStatus(formats strfmt.Registry) error {

	if swag.IsZero(m.Status) { // not required
		return nil
	}

	if err := m.Status.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("status")
		}
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *ReversalSubmissionAttributes) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ReversalSubmissionAttributes) UnmarshalBinary(b []byte) error {
	var res ReversalSubmissionAttributes
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}