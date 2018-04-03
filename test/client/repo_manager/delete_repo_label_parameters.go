// Code generated by go-swagger; DO NOT EDIT.

package repo_manager

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

	"openpitrix.io/openpitrix/test/models"
)

// NewDeleteRepoLabelParams creates a new DeleteRepoLabelParams object
// with the default values initialized.
func NewDeleteRepoLabelParams() *DeleteRepoLabelParams {
	var ()
	return &DeleteRepoLabelParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewDeleteRepoLabelParamsWithTimeout creates a new DeleteRepoLabelParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewDeleteRepoLabelParamsWithTimeout(timeout time.Duration) *DeleteRepoLabelParams {
	var ()
	return &DeleteRepoLabelParams{

		timeout: timeout,
	}
}

// NewDeleteRepoLabelParamsWithContext creates a new DeleteRepoLabelParams object
// with the default values initialized, and the ability to set a context for a request
func NewDeleteRepoLabelParamsWithContext(ctx context.Context) *DeleteRepoLabelParams {
	var ()
	return &DeleteRepoLabelParams{

		Context: ctx,
	}
}

// NewDeleteRepoLabelParamsWithHTTPClient creates a new DeleteRepoLabelParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewDeleteRepoLabelParamsWithHTTPClient(client *http.Client) *DeleteRepoLabelParams {
	var ()
	return &DeleteRepoLabelParams{
		HTTPClient: client,
	}
}

/*DeleteRepoLabelParams contains all the parameters to send to the API endpoint
for the delete repo label operation typically these are written to a http.Request
*/
type DeleteRepoLabelParams struct {

	/*Body*/
	Body *models.OpenpitrixDeleteRepoLabelRequest

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the delete repo label params
func (o *DeleteRepoLabelParams) WithTimeout(timeout time.Duration) *DeleteRepoLabelParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the delete repo label params
func (o *DeleteRepoLabelParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the delete repo label params
func (o *DeleteRepoLabelParams) WithContext(ctx context.Context) *DeleteRepoLabelParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the delete repo label params
func (o *DeleteRepoLabelParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the delete repo label params
func (o *DeleteRepoLabelParams) WithHTTPClient(client *http.Client) *DeleteRepoLabelParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the delete repo label params
func (o *DeleteRepoLabelParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithBody adds the body to the delete repo label params
func (o *DeleteRepoLabelParams) WithBody(body *models.OpenpitrixDeleteRepoLabelRequest) *DeleteRepoLabelParams {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the delete repo label params
func (o *DeleteRepoLabelParams) SetBody(body *models.OpenpitrixDeleteRepoLabelRequest) {
	o.Body = body
}

// WriteToRequest writes these params to a swagger request
func (o *DeleteRepoLabelParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.Body != nil {
		if err := r.SetBodyParam(o.Body); err != nil {
			return err
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
