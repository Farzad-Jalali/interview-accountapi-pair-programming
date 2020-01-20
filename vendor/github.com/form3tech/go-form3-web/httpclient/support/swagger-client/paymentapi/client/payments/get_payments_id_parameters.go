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

// NewGetPaymentsIDParams creates a new GetPaymentsIDParams object
// with the default values initialized.
func NewGetPaymentsIDParams() *GetPaymentsIDParams {
	var ()
	return &GetPaymentsIDParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewGetPaymentsIDParamsWithTimeout creates a new GetPaymentsIDParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewGetPaymentsIDParamsWithTimeout(timeout time.Duration) *GetPaymentsIDParams {
	var ()
	return &GetPaymentsIDParams{

		timeout: timeout,
	}
}

// NewGetPaymentsIDParamsWithContext creates a new GetPaymentsIDParams object
// with the default values initialized, and the ability to set a context for a request
func NewGetPaymentsIDParamsWithContext(ctx context.Context) *GetPaymentsIDParams {
	var ()
	return &GetPaymentsIDParams{

		Context: ctx,
	}
}

// NewGetPaymentsIDParamsWithHTTPClient creates a new GetPaymentsIDParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewGetPaymentsIDParamsWithHTTPClient(client *http.Client) *GetPaymentsIDParams {
	var ()
	return &GetPaymentsIDParams{
		HTTPClient: client,
	}
}

/*GetPaymentsIDParams contains all the parameters to send to the API endpoint
for the get payments ID operation typically these are written to a http.Request
*/
type GetPaymentsIDParams struct {

	/*ID
	  Payment Id

	*/
	ID strfmt.UUID

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the get payments ID params
func (o *GetPaymentsIDParams) WithTimeout(timeout time.Duration) *GetPaymentsIDParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get payments ID params
func (o *GetPaymentsIDParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get payments ID params
func (o *GetPaymentsIDParams) WithContext(ctx context.Context) *GetPaymentsIDParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get payments ID params
func (o *GetPaymentsIDParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get payments ID params
func (o *GetPaymentsIDParams) WithHTTPClient(client *http.Client) *GetPaymentsIDParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get payments ID params
func (o *GetPaymentsIDParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithID adds the id to the get payments ID params
func (o *GetPaymentsIDParams) WithID(id strfmt.UUID) *GetPaymentsIDParams {
	o.SetID(id)
	return o
}

// SetID adds the id to the get payments ID params
func (o *GetPaymentsIDParams) SetID(id strfmt.UUID) {
	o.ID = id
}

// WriteToRequest writes these params to a swagger request
func (o *GetPaymentsIDParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param id
	if err := r.SetPathParam("id", o.ID.String()); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}