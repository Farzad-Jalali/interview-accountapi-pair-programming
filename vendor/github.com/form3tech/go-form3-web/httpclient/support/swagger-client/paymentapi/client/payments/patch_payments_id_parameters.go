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

	"github.com/form3tech/go-form3-web/httpclient/support/swagger-client/paymentapi/models"
)

// NewPatchPaymentsIDParams creates a new PatchPaymentsIDParams object
// with the default values initialized.
func NewPatchPaymentsIDParams() *PatchPaymentsIDParams {
	var ()
	return &PatchPaymentsIDParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewPatchPaymentsIDParamsWithTimeout creates a new PatchPaymentsIDParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewPatchPaymentsIDParamsWithTimeout(timeout time.Duration) *PatchPaymentsIDParams {
	var ()
	return &PatchPaymentsIDParams{

		timeout: timeout,
	}
}

// NewPatchPaymentsIDParamsWithContext creates a new PatchPaymentsIDParams object
// with the default values initialized, and the ability to set a context for a request
func NewPatchPaymentsIDParamsWithContext(ctx context.Context) *PatchPaymentsIDParams {
	var ()
	return &PatchPaymentsIDParams{

		Context: ctx,
	}
}

// NewPatchPaymentsIDParamsWithHTTPClient creates a new PatchPaymentsIDParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewPatchPaymentsIDParamsWithHTTPClient(client *http.Client) *PatchPaymentsIDParams {
	var ()
	return &PatchPaymentsIDParams{
		HTTPClient: client,
	}
}

/*PatchPaymentsIDParams contains all the parameters to send to the API endpoint
for the patch payments ID operation typically these are written to a http.Request
*/
type PatchPaymentsIDParams struct {

	/*PaymentAmendOrDeleteRequest*/
	PaymentAmendOrDeleteRequest *models.PaymentAmendment
	/*ID
	  Payment Id

	*/
	ID strfmt.UUID

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the patch payments ID params
func (o *PatchPaymentsIDParams) WithTimeout(timeout time.Duration) *PatchPaymentsIDParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the patch payments ID params
func (o *PatchPaymentsIDParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the patch payments ID params
func (o *PatchPaymentsIDParams) WithContext(ctx context.Context) *PatchPaymentsIDParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the patch payments ID params
func (o *PatchPaymentsIDParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the patch payments ID params
func (o *PatchPaymentsIDParams) WithHTTPClient(client *http.Client) *PatchPaymentsIDParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the patch payments ID params
func (o *PatchPaymentsIDParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithPaymentAmendOrDeleteRequest adds the paymentAmendOrDeleteRequest to the patch payments ID params
func (o *PatchPaymentsIDParams) WithPaymentAmendOrDeleteRequest(paymentAmendOrDeleteRequest *models.PaymentAmendment) *PatchPaymentsIDParams {
	o.SetPaymentAmendOrDeleteRequest(paymentAmendOrDeleteRequest)
	return o
}

// SetPaymentAmendOrDeleteRequest adds the paymentAmendOrDeleteRequest to the patch payments ID params
func (o *PatchPaymentsIDParams) SetPaymentAmendOrDeleteRequest(paymentAmendOrDeleteRequest *models.PaymentAmendment) {
	o.PaymentAmendOrDeleteRequest = paymentAmendOrDeleteRequest
}

// WithID adds the id to the patch payments ID params
func (o *PatchPaymentsIDParams) WithID(id strfmt.UUID) *PatchPaymentsIDParams {
	o.SetID(id)
	return o
}

// SetID adds the id to the patch payments ID params
func (o *PatchPaymentsIDParams) SetID(id strfmt.UUID) {
	o.ID = id
}

// WriteToRequest writes these params to a swagger request
func (o *PatchPaymentsIDParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.PaymentAmendOrDeleteRequest != nil {
		if err := r.SetBodyParam(o.PaymentAmendOrDeleteRequest); err != nil {
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