// Code generated by go-swagger; DO NOT EDIT.

package payments

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/form3tech/go-form3-web/httpclient/support/swagger-client/paymentapi/models"
)

// GetPaymentsIDReader is a Reader for the GetPaymentsID structure.
type GetPaymentsIDReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetPaymentsIDReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewGetPaymentsIDOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewGetPaymentsIDOK creates a GetPaymentsIDOK with default headers values
func NewGetPaymentsIDOK() *GetPaymentsIDOK {
	return &GetPaymentsIDOK{}
}

/*GetPaymentsIDOK handles this case with default header values.

Payment details
*/
type GetPaymentsIDOK struct {
	Payload *models.PaymentDetailsResponse
}

func (o *GetPaymentsIDOK) Error() string {
	return fmt.Sprintf("[GET /payments/{id}][%d] getPaymentsIdOK  %+v", 200, o.Payload)
}

func (o *GetPaymentsIDOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.PaymentDetailsResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}