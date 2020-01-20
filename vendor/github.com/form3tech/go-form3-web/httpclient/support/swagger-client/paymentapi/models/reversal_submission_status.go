// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/validate"
)

// ReversalSubmissionStatus reversal submission status
// swagger:model ReversalSubmissionStatus
type ReversalSubmissionStatus string

const (
	// ReversalSubmissionStatusAccepted captures enum value "accepted"
	ReversalSubmissionStatusAccepted ReversalSubmissionStatus = "accepted"
	// ReversalSubmissionStatusValidationPassed captures enum value "validation_passed"
	ReversalSubmissionStatusValidationPassed ReversalSubmissionStatus = "validation_passed"
	// ReversalSubmissionStatusReleasedToGateway captures enum value "released_to_gateway"
	ReversalSubmissionStatusReleasedToGateway ReversalSubmissionStatus = "released_to_gateway"
	// ReversalSubmissionStatusQueuedForDelivery captures enum value "queued_for_delivery"
	ReversalSubmissionStatusQueuedForDelivery ReversalSubmissionStatus = "queued_for_delivery"
	// ReversalSubmissionStatusDeliveryConfirmed captures enum value "delivery_confirmed"
	ReversalSubmissionStatusDeliveryConfirmed ReversalSubmissionStatus = "delivery_confirmed"
	// ReversalSubmissionStatusDeliveryFailed captures enum value "delivery_failed"
	ReversalSubmissionStatusDeliveryFailed ReversalSubmissionStatus = "delivery_failed"
	// ReversalSubmissionStatusSubmitted captures enum value "submitted"
	ReversalSubmissionStatusSubmitted ReversalSubmissionStatus = "submitted"
	// ReversalSubmissionStatusValidationPending captures enum value "validation_pending"
	ReversalSubmissionStatusValidationPending ReversalSubmissionStatus = "validation_pending"
)

// for schema
var reversalSubmissionStatusEnum []interface{}

func init() {
	var res []ReversalSubmissionStatus
	if err := json.Unmarshal([]byte(`["accepted","validation_passed","released_to_gateway","queued_for_delivery","delivery_confirmed","delivery_failed","submitted","validation_pending"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		reversalSubmissionStatusEnum = append(reversalSubmissionStatusEnum, v)
	}
}

func (m ReversalSubmissionStatus) validateReversalSubmissionStatusEnum(path, location string, value ReversalSubmissionStatus) error {
	if err := validate.Enum(path, location, value, reversalSubmissionStatusEnum); err != nil {
		return err
	}
	return nil
}

// Validate validates this reversal submission status
func (m ReversalSubmissionStatus) Validate(formats strfmt.Registry) error {
	var res []error

	// value enum
	if err := m.validateReversalSubmissionStatusEnum("", "body", m); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}