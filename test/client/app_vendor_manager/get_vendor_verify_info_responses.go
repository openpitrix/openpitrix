// Code generated by go-swagger; DO NOT EDIT.

package app_vendor_manager

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	"openpitrix.io/openpitrix/test/models"
)

// GetVendorVerifyInfoReader is a Reader for the GetVendorVerifyInfo structure.
type GetVendorVerifyInfoReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetVendorVerifyInfoReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewGetVendorVerifyInfoOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewGetVendorVerifyInfoOK creates a GetVendorVerifyInfoOK with default headers values
func NewGetVendorVerifyInfoOK() *GetVendorVerifyInfoOK {
	return &GetVendorVerifyInfoOK{}
}

/*GetVendorVerifyInfoOK handles this case with default header values.

A successful response.
*/
type GetVendorVerifyInfoOK struct {
	Payload *models.OpenpitrixVendorVerifyInfo
}

func (o *GetVendorVerifyInfoOK) Error() string {
	return fmt.Sprintf("[GET /v1/vendor_verify_infos/{user_id.value}][%d] getVendorVerifyInfoOK  %+v", 200, o.Payload)
}

func (o *GetVendorVerifyInfoOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.OpenpitrixVendorVerifyInfo)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
