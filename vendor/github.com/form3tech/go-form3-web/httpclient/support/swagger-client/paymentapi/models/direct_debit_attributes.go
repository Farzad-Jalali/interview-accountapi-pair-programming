// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// DirectDebitAttributes direct debit attributes
// swagger:model directDebitAttributes
type DirectDebitAttributes struct {

	// amount
	// Pattern: ^[0-9.]{0,20}$
	Amount string `json:"amount,omitempty"`

	// beneficiary party
	BeneficiaryParty *DirectDebitAttributesBeneficiaryParty `json:"beneficiary_party,omitempty"`

	// clearing id
	ClearingID string `json:"clearing_id,omitempty"`

	// currency
	Currency string `json:"currency,omitempty"`

	// debtor party
	DebtorParty *DirectDebitAttributesDebtorParty `json:"debtor_party,omitempty"`

	// numeric reference
	NumericReference string `json:"numeric_reference,omitempty"`

	// payment scheme
	PaymentScheme string `json:"payment_scheme,omitempty"`

	// processing date
	ProcessingDate strfmt.Date `json:"processing_date,omitempty"`

	// reference
	Reference string `json:"reference,omitempty"`

	// scheme payment sub type
	SchemePaymentSubType string `json:"scheme_payment_sub_type,omitempty"`
}

// Validate validates this direct debit attributes
func (m *DirectDebitAttributes) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateAmount(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateBeneficiaryParty(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateDebtorParty(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *DirectDebitAttributes) validateAmount(formats strfmt.Registry) error {

	if swag.IsZero(m.Amount) { // not required
		return nil
	}

	if err := validate.Pattern("amount", "body", string(m.Amount), `^[0-9.]{0,20}$`); err != nil {
		return err
	}

	return nil
}

func (m *DirectDebitAttributes) validateBeneficiaryParty(formats strfmt.Registry) error {

	if swag.IsZero(m.BeneficiaryParty) { // not required
		return nil
	}

	if m.BeneficiaryParty != nil {

		if err := m.BeneficiaryParty.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("beneficiary_party")
			}
			return err
		}
	}

	return nil
}

func (m *DirectDebitAttributes) validateDebtorParty(formats strfmt.Registry) error {

	if swag.IsZero(m.DebtorParty) { // not required
		return nil
	}

	if m.DebtorParty != nil {

		if err := m.DebtorParty.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("debtor_party")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *DirectDebitAttributes) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *DirectDebitAttributes) UnmarshalBinary(b []byte) error {
	var res DirectDebitAttributes
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}