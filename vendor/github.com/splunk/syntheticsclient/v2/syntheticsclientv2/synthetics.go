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
	details.RequestBody = sanitizeRequestDump(endpoint, requestDump)
	fmt.Println("************")
	fmt.Println(details.RequestBody)
	fmt.Println("************")

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

func sanitizeRequestDump(endpoint string, requestDump []byte) string {
	dump := redactHeader(string(requestDump), "X-Sf-Token")
	if strings.Contains(endpoint, "/cacerts") {
		return redactRequestJSONField(dump, "content")
	}
	return dump
}

func redactHeader(dump string, header string) string {
	lines := strings.Split(dump, "\n")
	prefix := strings.ToLower(header) + ":"
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(strings.ToLower(trimmed), prefix) {
			suffix := ""
			if strings.HasSuffix(line, "\r") {
				suffix = "\r"
			}
			lines[i] = header + ": <REDACTED>" + suffix
		}
	}
	return strings.Join(lines, "\n")
}

func redactRequestJSONField(dump string, field string) string {
	separator := "\r\n\r\n"
	parts := strings.SplitN(dump, separator, 2)
	if len(parts) != 2 {
		separator = "\n\n"
		parts = strings.SplitN(dump, separator, 2)
	}
	if len(parts) != 2 {
		return dump
	}

	var body interface{}
	if err := json.Unmarshal([]byte(parts[1]), &body); err != nil {
		return dump
	}
	redactJSONField(body, field)
	redactedBody, err := json.Marshal(body)
	if err != nil {
		return dump
	}
	return parts[0] + separator + string(redactedBody)
}

func redactJSONField(value interface{}, field string) {
	switch typed := value.(type) {
	case map[string]interface{}:
		for key, nested := range typed {
			if strings.EqualFold(key, field) {
				typed[key] = "<REDACTED>"
				continue
			}
			redactJSONField(nested, field)
		}
	case []interface{}:
		for _, nested := range typed {
			redactJSONField(nested, field)
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
