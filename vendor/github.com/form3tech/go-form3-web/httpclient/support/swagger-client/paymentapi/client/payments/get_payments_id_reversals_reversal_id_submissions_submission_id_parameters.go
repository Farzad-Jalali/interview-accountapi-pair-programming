// Code generated by go-swagger; DO NOT EDIT.

package payments

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"
	"time"

	"golang.org/x/net/context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"

	strfmt "github.com/go-openapi/strfmt"
)

// NewGetPaymentsIDReversalsReversalIDSubmissionsSubmissionIDParams creates a new GetPaymentsIDReversalsReversalIDSubmissionsSubmissionIDParams object
// with the default values initialized.
func NewGetPaymentsIDReversalsReversalIDSubmissionsSubmissionIDParams() *GetPaymentsIDReversalsReversalIDSubmissionsSubmissionIDParams {
	var ()
	return &GetPaymentsIDReversalsReversalIDSubmissionsSubmissionIDParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewGetPaymentsIDReversalsReversalIDSubmissionsSubmissionIDParamsWithTimeout creates a new GetPaymentsIDReversalsReversalIDSubmissionsSubmissionIDParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewGetPaymentsIDReversalsReversalIDSubmissionsSubmissionIDParamsWithTimeout(timeout time.Duration) *GetPaymentsIDReversalsReversalIDSubmissionsSubmissionIDParams {
	var ()
	return &GetPaymentsIDReversalsReversalIDSubmissionsSubmissionIDParams{

		timeout: timeout,
	}
}

// NewGetPaymentsIDReversalsReversalIDSubmissionsSubmissionIDParamsWithContext creates a new GetPaymentsIDReversalsReversalIDSubmissionsSubmissionIDParams object
// with the default values initialized, and the ability to set a context for a request
func NewGetPaymentsIDReversalsReversalIDSubmissionsSubmissionIDParamsWithContext(ctx context.Context) *GetPaymentsIDReversalsReversalIDSubmissionsSubmissionIDParams {
	var ()
	return &GetPaymentsIDReversalsReversalIDSubmissionsSubmissionIDParams{

		Context: ctx,
	}
}

// NewGetPaymentsIDReversalsReversalIDSubmissionsSubmissionIDParamsWithHTTPClient creates a new GetPaymentsIDReversalsReversalIDSubmissionsSubmissionIDParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewGetPaymentsIDReversalsReversalIDSubmissionsSubmissionIDParamsWithHTTPClient(client *http.Client) *GetPaymentsIDReversalsReversalIDSubmissionsSubmissionIDParams {
	var ()
	return &GetPaymentsIDReversalsReversalIDSubmissionsSubmissionIDParams{
		HTTPClient: client,
	}
}

/*GetPaymentsIDReversalsReversalIDSubmissionsSubmissionIDParams contains all the parameters to send to the API endpoint
for the get payments ID reversals reversal ID submissions submission ID operation typically these are written to a http.Request
*/
type GetPaymentsIDReversalsReversalIDSubmissionsSubmissionIDParams struct {

	/*ID
	  Payment Id

	*/
	ID strfmt.UUID
	/*ReversalID
	  Reversal Id

	*/
	ReversalID strfmt.UUID
	/*SubmissionID
	  Submission Id

	*/
	SubmissionID strfmt.UUID

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the get payments ID reversals reversal ID submissions submission ID params
func (o *GetPaymentsIDReversalsReversalIDSubmissionsSubmissionIDParams) WithTimeout(timeout time.Duration) *GetPaymentsIDReversalsReversalIDSubmissionsSubmissionIDParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get payments ID reversals reversal ID submissions submission ID params
func (o *GetPaymentsIDReversalsReversalIDSubmissionsSubmissionIDParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get payments ID reversals reversal ID submissions submission ID params
func (o *GetPaymentsIDReversalsReversalIDSubmissionsSubmissionIDParams) WithContext(ctx context.Context) *GetPaymentsIDReversalsReversalIDSubmissionsSubmissionIDParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get payments ID reversals reversal ID submissions submission ID params
func (o *GetPaymentsIDReversalsReversalIDSubmissionsSubmissionIDParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get payments ID reversals reversal ID submissions submission ID params
func (o *GetPaymentsIDReversalsReversalIDSubmissionsSubmissionIDParams) WithHTTPClient(client *http.Client) *GetPaymentsIDReversalsReversalIDSubmissionsSubmissionIDParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get payments ID reversals reversal ID submissions submission ID params
func (o *GetPaymentsIDReversalsReversalIDSubmissionsSubmissionIDParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithID adds the id to the get payments ID reversals reversal ID submissions submission ID params
func (o *GetPaymentsIDReversalsReversalIDSubmissionsSubmissionIDParams) WithID(id strfmt.UUID) *GetPaymentsIDReversalsReversalIDSubmissionsSubmissionIDParams {
	o.SetID(id)
	return o
}

// SetID adds the id to the get payments ID reversals reversal ID submissions submission ID params
func (o *GetPaymentsIDReversalsReversalIDSubmissionsSubmissionIDParams) SetID(id strfmt.UUID) {
	o.ID = id
}

// WithReversalID adds the reversalID to the get payments ID reversals reversal ID submissions submission ID params
func (o *GetPaymentsIDReversalsReversalIDSubmissionsSubmissionIDParams) WithReversalID(reversalID strfmt.UUID) *GetPaymentsIDReversalsReversalIDSubmissionsSubmissionIDParams {
	o.SetReversalID(reversalID)
	return o
}

// SetReversalID adds the reversalId to the get payments ID reversals reversal ID submissions submission ID params
func (o *GetPaymentsIDReversalsReversalIDSubmissionsSubmissionIDParams) SetReversalID(reversalID strfmt.UUID) {
	o.ReversalID = reversalID
}

// WithSubmissionID adds the submissionID to the get payments ID reversals reversal ID submissions submission ID params
func (o *GetPaymentsIDReversalsReversalIDSubmissionsSubmissionIDParams) WithSubmissionID(submissionID strfmt.UUID) *GetPaymentsIDReversalsReversalIDSubmissionsSubmissionIDParams {
	o.SetSubmissionID(submissionID)
	return o
}

// SetSubmissionID adds the submissionId to the get payments ID reversals reversal ID submissions submission ID params
func (o *GetPaymentsIDReversalsReversalIDSubmissionsSubmissionIDParams) SetSubmissionID(submissionID strfmt.UUID) {
	o.SubmissionID = submissionID
}

// WriteToRequest writes these params to a swagger request
func (o *GetPaymentsIDReversalsReversalIDSubmissionsSubmissionIDParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param id
	if err := r.SetPathParam("id", o.ID.String()); err != nil {
		return err
	}

	// path param reversalId
	if err := r.SetPathParam("reversalId", o.ReversalID.String()); err != nil {
		return err
	}

	// path param submissionId
	if err := r.SetPathParam("submissionId", o.SubmissionID.String()); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}