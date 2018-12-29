// Code generated by go-swagger; DO NOT EDIT.

package app_vendor_manager

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"
)

// New creates a new app vendor manager API client.
func New(transport runtime.ClientTransport, formats strfmt.Registry) *Client {
	return &Client{transport: transport, formats: formats}
}

/*
Client for app vendor manager API
*/
type Client struct {
	transport runtime.ClientTransport
	formats   strfmt.Registry
}

/*
DescribeAppVendorStatistics describes app vendor statistics
*/
func (a *Client) DescribeAppVendorStatistics(params *DescribeAppVendorStatisticsParams, authInfo runtime.ClientAuthInfoWriter) (*DescribeAppVendorStatisticsOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewDescribeAppVendorStatisticsParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "DescribeAppVendorStatistics",
		Method:             "GET",
		PathPattern:        "/v1/DescribeAppVendorStatistics",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http", "https"},
		Params:             params,
		Reader:             &DescribeAppVendorStatisticsReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*DescribeAppVendorStatisticsOK), nil

}

/*
DescribeVendorVerifyInfos describes vendor verify infos
*/
func (a *Client) DescribeVendorVerifyInfos(params *DescribeVendorVerifyInfosParams, authInfo runtime.ClientAuthInfoWriter) (*DescribeVendorVerifyInfosOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewDescribeVendorVerifyInfosParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "DescribeVendorVerifyInfos",
		Method:             "GET",
		PathPattern:        "/v1/vendor_verify_infos",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http", "https"},
		Params:             params,
		Reader:             &DescribeVendorVerifyInfosReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*DescribeVendorVerifyInfosOK), nil

}

/*
GetVendorVerifyInfo gets vendor verify info
*/
func (a *Client) GetVendorVerifyInfo(params *GetVendorVerifyInfoParams, authInfo runtime.ClientAuthInfoWriter) (*GetVendorVerifyInfoOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetVendorVerifyInfoParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "GetVendorVerifyInfo",
		Method:             "GET",
		PathPattern:        "/v1/vendor_verify_infos/user_id=*",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http", "https"},
		Params:             params,
		Reader:             &GetVendorVerifyInfoReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*GetVendorVerifyInfoOK), nil

}

/*
PassVendorVerifyInfo passes vendor verify info
*/
func (a *Client) PassVendorVerifyInfo(params *PassVendorVerifyInfoParams, authInfo runtime.ClientAuthInfoWriter) (*PassVendorVerifyInfoOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewPassVendorVerifyInfoParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "PassVendorVerifyInfo",
		Method:             "POST",
		PathPattern:        "/v1/vendor_verify_infos/user_id=*/action:pass",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http", "https"},
		Params:             params,
		Reader:             &PassVendorVerifyInfoReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*PassVendorVerifyInfoOK), nil

}

/*
RejectVendorVerifyInfo rejects vendor verify info
*/
func (a *Client) RejectVendorVerifyInfo(params *RejectVendorVerifyInfoParams, authInfo runtime.ClientAuthInfoWriter) (*RejectVendorVerifyInfoOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewRejectVendorVerifyInfoParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "RejectVendorVerifyInfo",
		Method:             "POST",
		PathPattern:        "/v1/vendor_verify_infos/user_id=*/action:reject",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http", "https"},
		Params:             params,
		Reader:             &RejectVendorVerifyInfoReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*RejectVendorVerifyInfoOK), nil

}

/*
SubmitVendorVerifyInfo submits vendor verify info
*/
func (a *Client) SubmitVendorVerifyInfo(params *SubmitVendorVerifyInfoParams, authInfo runtime.ClientAuthInfoWriter) (*SubmitVendorVerifyInfoOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewSubmitVendorVerifyInfoParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "SubmitVendorVerifyInfo",
		Method:             "POST",
		PathPattern:        "/v1/vendor_verify_infos/{user_id}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http", "https"},
		Params:             params,
		Reader:             &SubmitVendorVerifyInfoReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*SubmitVendorVerifyInfoOK), nil

}

// SetTransport changes the transport on the client
func (a *Client) SetTransport(transport runtime.ClientTransport) {
	a.transport = transport
}
