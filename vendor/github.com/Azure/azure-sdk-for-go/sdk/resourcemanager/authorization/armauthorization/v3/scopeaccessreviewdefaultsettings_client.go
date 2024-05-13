//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armauthorization

import (
	"context"
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
	"net/url"
	"strings"
)

// ScopeAccessReviewDefaultSettingsClient contains the methods for the ScopeAccessReviewDefaultSettings group.
// Don't use this type directly, use NewScopeAccessReviewDefaultSettingsClient() instead.
type ScopeAccessReviewDefaultSettingsClient struct {
	internal *arm.Client
}

// NewScopeAccessReviewDefaultSettingsClient creates a new instance of ScopeAccessReviewDefaultSettingsClient with the specified values.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewScopeAccessReviewDefaultSettingsClient(credential azcore.TokenCredential, options *arm.ClientOptions) (*ScopeAccessReviewDefaultSettingsClient, error) {
	cl, err := arm.NewClient(moduleName, moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &ScopeAccessReviewDefaultSettingsClient{
		internal: cl,
	}
	return client, nil
}

// Get - Get access review default settings for the subscription
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-12-01-preview
//   - scope - The scope of the resource.
//   - options - ScopeAccessReviewDefaultSettingsClientGetOptions contains the optional parameters for the ScopeAccessReviewDefaultSettingsClient.Get
//     method.
func (client *ScopeAccessReviewDefaultSettingsClient) Get(ctx context.Context, scope string, options *ScopeAccessReviewDefaultSettingsClientGetOptions) (ScopeAccessReviewDefaultSettingsClientGetResponse, error) {
	var err error
	const operationName = "ScopeAccessReviewDefaultSettingsClient.Get"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.getCreateRequest(ctx, scope, options)
	if err != nil {
		return ScopeAccessReviewDefaultSettingsClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return ScopeAccessReviewDefaultSettingsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return ScopeAccessReviewDefaultSettingsClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *ScopeAccessReviewDefaultSettingsClient) getCreateRequest(ctx context.Context, scope string, options *ScopeAccessReviewDefaultSettingsClientGetOptions) (*policy.Request, error) {
	urlPath := "/{scope}/providers/Microsoft.Authorization/accessReviewScheduleSettings/default"
	if scope == "" {
		return nil, errors.New("parameter scope cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{scope}", url.PathEscape(scope))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-12-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *ScopeAccessReviewDefaultSettingsClient) getHandleResponse(resp *http.Response) (ScopeAccessReviewDefaultSettingsClientGetResponse, error) {
	result := ScopeAccessReviewDefaultSettingsClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.AccessReviewDefaultSettings); err != nil {
		return ScopeAccessReviewDefaultSettingsClientGetResponse{}, err
	}
	return result, nil
}

// Put - Get access review default settings for the subscription
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-12-01-preview
//   - scope - The scope of the resource.
//   - properties - Access review schedule settings.
//   - options - ScopeAccessReviewDefaultSettingsClientPutOptions contains the optional parameters for the ScopeAccessReviewDefaultSettingsClient.Put
//     method.
func (client *ScopeAccessReviewDefaultSettingsClient) Put(ctx context.Context, scope string, properties AccessReviewScheduleSettings, options *ScopeAccessReviewDefaultSettingsClientPutOptions) (ScopeAccessReviewDefaultSettingsClientPutResponse, error) {
	var err error
	const operationName = "ScopeAccessReviewDefaultSettingsClient.Put"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.putCreateRequest(ctx, scope, properties, options)
	if err != nil {
		return ScopeAccessReviewDefaultSettingsClientPutResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return ScopeAccessReviewDefaultSettingsClientPutResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return ScopeAccessReviewDefaultSettingsClientPutResponse{}, err
	}
	resp, err := client.putHandleResponse(httpResp)
	return resp, err
}

// putCreateRequest creates the Put request.
func (client *ScopeAccessReviewDefaultSettingsClient) putCreateRequest(ctx context.Context, scope string, properties AccessReviewScheduleSettings, options *ScopeAccessReviewDefaultSettingsClientPutOptions) (*policy.Request, error) {
	urlPath := "/{scope}/providers/Microsoft.Authorization/accessReviewScheduleSettings/default"
	if scope == "" {
		return nil, errors.New("parameter scope cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{scope}", url.PathEscape(scope))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-12-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, properties); err != nil {
		return nil, err
	}
	return req, nil
}

// putHandleResponse handles the Put response.
func (client *ScopeAccessReviewDefaultSettingsClient) putHandleResponse(resp *http.Response) (ScopeAccessReviewDefaultSettingsClientPutResponse, error) {
	result := ScopeAccessReviewDefaultSettingsClientPutResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.AccessReviewDefaultSettings); err != nil {
		return ScopeAccessReviewDefaultSettingsClientPutResponse{}, err
	}
	return result, nil
}
