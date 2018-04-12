// Code generated by go-swagger; DO NOT EDIT.

package runtime_manager

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

// NewDescribeRuntimeProviderZonesParams creates a new DescribeRuntimeProviderZonesParams object
// with the default values initialized.
func NewDescribeRuntimeProviderZonesParams() *DescribeRuntimeProviderZonesParams {
	var ()
	return &DescribeRuntimeProviderZonesParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewDescribeRuntimeProviderZonesParamsWithTimeout creates a new DescribeRuntimeProviderZonesParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewDescribeRuntimeProviderZonesParamsWithTimeout(timeout time.Duration) *DescribeRuntimeProviderZonesParams {
	var ()
	return &DescribeRuntimeProviderZonesParams{

		timeout: timeout,
	}
}

// NewDescribeRuntimeProviderZonesParamsWithContext creates a new DescribeRuntimeProviderZonesParams object
// with the default values initialized, and the ability to set a context for a request
func NewDescribeRuntimeProviderZonesParamsWithContext(ctx context.Context) *DescribeRuntimeProviderZonesParams {
	var ()
	return &DescribeRuntimeProviderZonesParams{

		Context: ctx,
	}
}

// NewDescribeRuntimeProviderZonesParamsWithHTTPClient creates a new DescribeRuntimeProviderZonesParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewDescribeRuntimeProviderZonesParamsWithHTTPClient(client *http.Client) *DescribeRuntimeProviderZonesParams {
	var ()
	return &DescribeRuntimeProviderZonesParams{
		HTTPClient: client,
	}
}

/*DescribeRuntimeProviderZonesParams contains all the parameters to send to the API endpoint
for the describe runtime provider zones operation typically these are written to a http.Request
*/
type DescribeRuntimeProviderZonesParams struct {

	/*Provider*/
	Provider *string
	/*RuntimeCredential*/
	RuntimeCredential *string
	/*RuntimeURL*/
	RuntimeURL *string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the describe runtime provider zones params
func (o *DescribeRuntimeProviderZonesParams) WithTimeout(timeout time.Duration) *DescribeRuntimeProviderZonesParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the describe runtime provider zones params
func (o *DescribeRuntimeProviderZonesParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the describe runtime provider zones params
func (o *DescribeRuntimeProviderZonesParams) WithContext(ctx context.Context) *DescribeRuntimeProviderZonesParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the describe runtime provider zones params
func (o *DescribeRuntimeProviderZonesParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the describe runtime provider zones params
func (o *DescribeRuntimeProviderZonesParams) WithHTTPClient(client *http.Client) *DescribeRuntimeProviderZonesParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the describe runtime provider zones params
func (o *DescribeRuntimeProviderZonesParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithProvider adds the provider to the describe runtime provider zones params
func (o *DescribeRuntimeProviderZonesParams) WithProvider(provider *string) *DescribeRuntimeProviderZonesParams {
	o.SetProvider(provider)
	return o
}

// SetProvider adds the provider to the describe runtime provider zones params
func (o *DescribeRuntimeProviderZonesParams) SetProvider(provider *string) {
	o.Provider = provider
}

// WithRuntimeCredential adds the runtimeCredential to the describe runtime provider zones params
func (o *DescribeRuntimeProviderZonesParams) WithRuntimeCredential(runtimeCredential *string) *DescribeRuntimeProviderZonesParams {
	o.SetRuntimeCredential(runtimeCredential)
	return o
}

// SetRuntimeCredential adds the runtimeCredential to the describe runtime provider zones params
func (o *DescribeRuntimeProviderZonesParams) SetRuntimeCredential(runtimeCredential *string) {
	o.RuntimeCredential = runtimeCredential
}

// WithRuntimeURL adds the runtimeURL to the describe runtime provider zones params
func (o *DescribeRuntimeProviderZonesParams) WithRuntimeURL(runtimeURL *string) *DescribeRuntimeProviderZonesParams {
	o.SetRuntimeURL(runtimeURL)
	return o
}

// SetRuntimeURL adds the runtimeUrl to the describe runtime provider zones params
func (o *DescribeRuntimeProviderZonesParams) SetRuntimeURL(runtimeURL *string) {
	o.RuntimeURL = runtimeURL
}

// WriteToRequest writes these params to a swagger request
func (o *DescribeRuntimeProviderZonesParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.Provider != nil {

		// query param provider
		var qrProvider string
		if o.Provider != nil {
			qrProvider = *o.Provider
		}
		qProvider := qrProvider
		if qProvider != "" {
			if err := r.SetQueryParam("provider", qProvider); err != nil {
				return err
			}
		}

	}

	if o.RuntimeCredential != nil {

		// query param runtime_credential
		var qrRuntimeCredential string
		if o.RuntimeCredential != nil {
			qrRuntimeCredential = *o.RuntimeCredential
		}
		qRuntimeCredential := qrRuntimeCredential
		if qRuntimeCredential != "" {
			if err := r.SetQueryParam("runtime_credential", qRuntimeCredential); err != nil {
				return err
			}
		}

	}

	if o.RuntimeURL != nil {

		// query param runtime_url
		var qrRuntimeURL string
		if o.RuntimeURL != nil {
			qrRuntimeURL = *o.RuntimeURL
		}
		qRuntimeURL := qrRuntimeURL
		if qRuntimeURL != "" {
			if err := r.SetQueryParam("runtime_url", qRuntimeURL); err != nil {
				return err
			}
		}

	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}