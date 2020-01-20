// Code generated by go-swagger; DO NOT EDIT.

package direct_debits

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

	"github.com/form3tech/go-form3-web/httpclient/support/swagger-client/paymentapi/models"
)

// NewPostDirectdebitsIDAdmissionsParams creates a new PostDirectdebitsIDAdmissionsParams object
// with the default values initialized.
func NewPostDirectdebitsIDAdmissionsParams() *PostDirectdebitsIDAdmissionsParams {
	var ()
	return &PostDirectdebitsIDAdmissionsParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewPostDirectdebitsIDAdmissionsParamsWithTimeout creates a new PostDirectdebitsIDAdmissionsParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewPostDirectdebitsIDAdmissionsParamsWithTimeout(timeout time.Duration) *PostDirectdebitsIDAdmissionsParams {
	var ()
	return &PostDirectdebitsIDAdmissionsParams{

		timeout: timeout,
	}
}

// NewPostDirectdebitsIDAdmissionsParamsWithContext creates a new PostDirectdebitsIDAdmissionsParams object
// with the default values initialized, and the ability to set a context for a request
func NewPostDirectdebitsIDAdmissionsParamsWithContext(ctx context.Context) *PostDirectdebitsIDAdmissionsParams {
	var ()
	return &PostDirectdebitsIDAdmissionsParams{

		Context: ctx,
	}
}

// NewPostDirectdebitsIDAdmissionsParamsWithHTTPClient creates a new PostDirectdebitsIDAdmissionsParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewPostDirectdebitsIDAdmissionsParamsWithHTTPClient(client *http.Client) *PostDirectdebitsIDAdmissionsParams {
	var ()
	return &PostDirectdebitsIDAdmissionsParams{
		HTTPClient: client,
	}
}

/*PostDirectdebitsIDAdmissionsParams contains all the parameters to send to the API endpoint
for the post directdebits ID admissions operation typically these are written to a http.Request
*/
type PostDirectdebitsIDAdmissionsParams struct {

	/*DirectDebitAdmissionCreationRequest*/
	DirectDebitAdmissionCreationRequest *models.DirectDebitAdmissionCreation
	/*ID
	  Direct debit Id

	*/
	ID strfmt.UUID

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the post directdebits ID admissions params
func (o *PostDirectdebitsIDAdmissionsParams) WithTimeout(timeout time.Duration) *PostDirectdebitsIDAdmissionsParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the post directdebits ID admissions params
func (o *PostDirectdebitsIDAdmissionsParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the post directdebits ID admissions params
func (o *PostDirectdebitsIDAdmissionsParams) WithContext(ctx context.Context) *PostDirectdebitsIDAdmissionsParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the post directdebits ID admissions params
func (o *PostDirectdebitsIDAdmissionsParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the post directdebits ID admissions params
func (o *PostDirectdebitsIDAdmissionsParams) WithHTTPClient(client *http.Client) *PostDirectdebitsIDAdmissionsParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the post directdebits ID admissions params
func (o *PostDirectdebitsIDAdmissionsParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithDirectDebitAdmissionCreationRequest adds the directDebitAdmissionCreationRequest to the post directdebits ID admissions params
func (o *PostDirectdebitsIDAdmissionsParams) WithDirectDebitAdmissionCreationRequest(directDebitAdmissionCreationRequest *models.DirectDebitAdmissionCreation) *PostDirectdebitsIDAdmissionsParams {
	o.SetDirectDebitAdmissionCreationRequest(directDebitAdmissionCreationRequest)
	return o
}

// SetDirectDebitAdmissionCreationRequest adds the directDebitAdmissionCreationRequest to the post directdebits ID admissions params
func (o *PostDirectdebitsIDAdmissionsParams) SetDirectDebitAdmissionCreationRequest(directDebitAdmissionCreationRequest *models.DirectDebitAdmissionCreation) {
	o.DirectDebitAdmissionCreationRequest = directDebitAdmissionCreationRequest
}

// WithID adds the id to the post directdebits ID admissions params
func (o *PostDirectdebitsIDAdmissionsParams) WithID(id strfmt.UUID) *PostDirectdebitsIDAdmissionsParams {
	o.SetID(id)
	return o
}

// SetID adds the id to the post directdebits ID admissions params
func (o *PostDirectdebitsIDAdmissionsParams) SetID(id strfmt.UUID) {
	o.ID = id
}

// WriteToRequest writes these params to a swagger request
func (o *PostDirectdebitsIDAdmissionsParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.DirectDebitAdmissionCreationRequest != nil {
		if err := r.SetBodyParam(o.DirectDebitAdmissionCreationRequest); err != nil {
			return err
		}
	}

	// path param id
	if err := r.SetPathParam("id", o.ID.String()); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}