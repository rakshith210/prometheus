// Copyright 2023 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package azuread

import (
	"context"
	"net/http"
)

// azureADRoundTripper is round tripper adding Azure AD authorization to calls.
type azureADRoundTripper struct {
	next          http.RoundTripper
	tokenProvider tokenProvider
}

// NewAzureADRoundTripper creates round tripper adding Azure AD authorization to calls.
func NewAzureADRoundTripper(cfg *AzureADConfig, next http.RoundTripper) (http.RoundTripper, error) {
	if next == nil {
		next = http.DefaultTransport
	}

	cred, err := newTokenCredential(cfg)
	if err != nil {
		return nil, err
	}

	tokenProvider, err := newTokenProvider(context.Background(), cfg, cred)
	if err != nil {
		return nil, err
	}

	rt := &azureADRoundTripper{
		next:          next,
		tokenProvider: TokenProvider,
	}
	return rt, nil
}

// RoundTrip sets Authorization header for requests.
func (rt *azureADRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	accessToken, err := rt.tokenProvider.getAccessToken()
	if err != nil {
		return nil, err
	}
	bearerAccessToken := "Bearer " + accessToken
	req.Header.Set("Authorization", bearerAccessToken)

	return rt.next.RoundTrip(req)
}
