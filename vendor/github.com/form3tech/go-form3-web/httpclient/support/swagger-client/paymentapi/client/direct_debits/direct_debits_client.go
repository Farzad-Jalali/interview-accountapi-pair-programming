// Code generated by go-swagger; DO NOT EDIT.

package direct_debits

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"
)

// New creates a new direct debits API client.
func New(transport runtime.ClientTransport, formats strfmt.Registry) *Client {
	return &Client{transport: transport, formats: formats}
}

/*
Client for direct debits API
*/
type Client struct {
	transport runtime.ClientTransport
	formats   strfmt.Registry
}

/*
PostDirectdebits creates direct debit
*/
func (a *Client) PostDirectdebits(params *PostDirectdebitsParams) (*PostDirectdebitsCreated, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewPostDirectdebitsParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "PostDirectdebits",
		Method:             "POST",
		PathPattern:        "/directdebits",
		ProducesMediaTypes: []string{"application/json; charset=utf-8", "application/vnd.api+json; charset=utf-8"},
		ConsumesMediaTypes: []string{"application/json", "application/vnd.api+json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &PostDirectdebitsReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*PostDirectdebitsCreated), nil

}

/*
PostDirectdebitsIDAdmissions creates direct debit admission
*/
func (a *Client) PostDirectdebitsIDAdmissions(params *PostDirectdebitsIDAdmissionsParams) (*PostDirectdebitsIDAdmissionsCreated, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewPostDirectdebitsIDAdmissionsParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "PostDirectdebitsIDAdmissions",
		Method:             "POST",
		PathPattern:        "/directdebits/{id}/admissions",
		ProducesMediaTypes: []string{"application/json; charset=utf-8", "application/vnd.api+json; charset=utf-8"},
		ConsumesMediaTypes: []string{"application/json", "application/vnd.api+json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &PostDirectdebitsIDAdmissionsReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*PostDirectdebitsIDAdmissionsCreated), nil

}

// SetTransport changes the transport on the client
func (a *Client) SetTransport(transport runtime.ClientTransport) {
	a.transport = transport
}