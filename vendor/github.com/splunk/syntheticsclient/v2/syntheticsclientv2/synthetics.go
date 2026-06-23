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

package syntheticsclientv2

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"
)

type Client struct {
	publicBaseURL string
	apiKey        string
	realm         string
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
	Status  string                 `json:"status,omitempty"`
	Error   string                 `json:"error,omitempty"`
	Result  string                 `json:"result,omitempty"`
	Message string                 `json:"message,omitempty"`
	Errors  Errors                 `json:"errors,omitempty"`
	Details map[string]interface{} `json:"details,omitempty"`
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
	req.Header.Set("X-SF-TOKEN", c.apiKey)

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
	details.RequestBody = sanitizeRequestDump(string(requestDump), c.apiKey, endpoint)

	// Make the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return &details, err
	}

	details.StatusCode = resp.StatusCode
	details.RawResponse = resp

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
		var errRes errorResponse
		if err = json.NewDecoder(resp.Body).Decode(&errRes); err == nil {
			errorField, err2 := json.Marshal(errRes)
			if err2 != nil {
				return &details, fmt.Errorf("unknown issue while parsing API error response, status code: %d", resp.StatusCode)
			}
			return &details, errors.New("Status Code: " + resp.Status + "\n" + "Response: " + string(errorField))
		}
		return &details, fmt.Errorf("unknown error, status code: %d", resp.StatusCode)
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return &details, err
	}

	details.ResponseBody = string(responseBody)

	return &details, nil
}

func redactSensitiveValue(value string, sensitiveValue string) string {
	if sensitiveValue == "" {
		return value
	}

	return strings.ReplaceAll(value, sensitiveValue, "[REDACTED]")
}

func sanitizeRequestDump(requestDump string, apiKey string, endpoint string) string {
	sanitizedRequestDump := redactSensitiveValue(requestDump, apiKey)
	if !strings.Contains(endpoint, "/cacerts") {
		return sanitizedRequestDump
	}

	return redactCaCertificateContent(sanitizedRequestDump)
}

func redactCaCertificateContent(requestDump string) string {
	headers, body, separator, ok := splitRequestDump(requestDump)
	if !ok || body == "" {
		return requestDump
	}

	var requestBody interface{}
	if err := json.Unmarshal([]byte(body), &requestBody); err != nil {
		return replaceRequestDumpBody(headers, separator)
	}

	redactContentFields(requestBody)

	redactedRequestBody, err := json.Marshal(requestBody)
	if err != nil {
		return replaceRequestDumpBody(headers, separator)
	}

	return headers + separator + string(redactedRequestBody)
}

func splitRequestDump(requestDump string) (string, string, string, bool) {
	for _, separator := range []string{"\r\n\r\n", "\n\n"} {
		requestParts := strings.SplitN(requestDump, separator, 2)
		if len(requestParts) == 2 {
			return requestParts[0], requestParts[1], separator, true
		}
	}

	return "", "", "", false
}

func replaceRequestDumpBody(headers string, separator string) string {
	return headers + separator + "[REDACTED]"
}

func redactContentFields(value interface{}) {
	switch typedValue := value.(type) {
	case map[string]interface{}:
		for key, nestedValue := range typedValue {
			if strings.EqualFold(key, "content") {
				typedValue[key] = "[REDACTED]"
				continue
			}
			redactContentFields(nestedValue)
		}
	case []interface{}:
		for _, nestedValue := range typedValue {
			redactContentFields(nestedValue)
		}
	}
}

func NewClientArgs(timeout int, baseUrl string) ClientArgs {
	return ClientArgs{
		timeoutSeconds: timeout,
		publicBaseUrl:  baseUrl,
	}
}

func NewClient(apiKey string, realm string) *Client {
	args := ClientArgs{timeoutSeconds: 30}
	return NewConfigurableClient(apiKey, realm, args)
}

func NewConfigurableClient(apiKey string, realm string, args ClientArgs) *Client {
	client := Client{
		apiKey:     apiKey,
		realm:      realm,
		httpClient: http.Client{Timeout: time.Duration(args.timeoutSeconds) * time.Second},
	}
	if args.publicBaseUrl == "" {
		client.publicBaseURL = "https://api." + realm + ".signalfx.com/v2/synthetics"
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
