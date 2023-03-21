// Copyright 2021 Splunk, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package syntheticsclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"time"
)

type Client struct {
	publicBaseURL string
	apiKey        string
	httpClient    http.Client
}

type ClientArgs struct {
	timeoutSeconds int
	publicBaseUrl  string
}

type RequestDetails struct {
	StatusCode   int
	ResponseBody string
	RequestBody  string
	RawResponse  *http.Response
	RawRequest   *http.Request
}

type errorResponse struct {
	Status  string `json:"status,omitempty"`
	Error   string `json:"error,omitempty"`
	Result  string `json:"result,omitempty"`
	Message string `json:"message,omitempty"`
	Errors  Errors `json:"errors,omitempty"`
}

func (c Client) String() string {
	return fmt.Sprintf("Splunk Synthetics Client: URL: %s ", c.publicBaseURL)
}

func (c Client) makePublicAPICall(method string, endpoint string, requestBody io.Reader, queryParams map[string]string) (*RequestDetails, error) {
	details := RequestDetails{}
	// Create the request
	req, err := http.NewRequest(method, c.publicBaseURL+endpoint, requestBody)
	if err != nil {
		return &details, err
	}

	// Set the auth headers needed for the public api
	req.Header.Set("api-key", c.apiKey)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Set the query params
	q := req.URL.Query()
	for key, value := range queryParams {
		q.Add(key, value)
	}
	req.URL.RawQuery = q.Encode()

	// Add the request to the details
	details.RawRequest = req
	requestDump, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		return &details, err
	}
	details.RequestBody = string(requestDump)

	// Make the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return &details, err
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
		var errRes errorResponse
		if err = json.NewDecoder(resp.Body).Decode(&errRes); err == nil {
			errorField, err2 := json.Marshal(errRes)
			if err2 != nil {
				return &details, fmt.Errorf("unknown issue while parsing API error response, status code: %d", resp.StatusCode)
			}
			return &details, errors.New("Status Code: " + resp.Status + "\n" + string(errorField))
		}
		return &details, fmt.Errorf("unknown error, status code: %d", resp.StatusCode)
	}

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &details, err
	}

	details.StatusCode = resp.StatusCode
	details.ResponseBody = string(responseBody)
	details.RawResponse = resp

	return &details, nil
}

func NewClient(apiKey string) *Client {
	args := ClientArgs{timeoutSeconds: 30}
	return NewConfigurableClient(apiKey, args)
}

func NewConfigurableClient(apiKey string, args ClientArgs) *Client {
	client := Client{
		apiKey:     apiKey,
		httpClient: http.Client{Timeout: time.Duration(args.timeoutSeconds) * time.Second},
	}
	if args.publicBaseUrl == "" {
		client.publicBaseURL = "https://monitoring-api.rigor.com"
	} else {
		client.publicBaseURL = args.publicBaseUrl
	}

	return &client
}

// GetHTTPClient returns http client for the purpose of test
func (c Client) GetHTTPClient() *http.Client {
	return &c.httpClient
}

// Helper for tests and output
func JsonPrint(data interface{}) {
	var p []byte
	p, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%s \n", p)
}
