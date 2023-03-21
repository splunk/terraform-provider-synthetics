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
	"bytes"
	"encoding/json"
	"fmt"
)

type GetChecks struct {
	CurrentPage  int    `json:"current_page"`
	PerPage      int    `json:"per_page"`
	NextPage     int    `json:"next_page"`
	PreviousPage int    `json:"previous_page"`
	TotalPages   int    `json:"total_pages"`
	TotalCount   int    `json:"total_count"`
	Checks       Checks `json:"checks"`
}

// Leaving off "Enabled" filter setting. Can be added later if required.
type GetChecksOptions struct {
	Type    string `json:"type"`
	PerPage int    `json:"per_page"`
	Page    int    `json:"page"`
	Muted   bool   `json:"muted"`
}

func parseChecksResponse(response string) (*GetChecks, error) {
	// Parse the response and return the check object
	var checks GetChecks
	err := json.Unmarshal([]byte(response), &checks)
	if err != nil {
		return nil, err
	}

	return &checks, err
}

// GetChecks returns all checks
func (c Client) GetChecks(params *GetChecksOptions) (*GetChecks, *RequestDetails, error) {
	// Check for default params
	if params.Type == "" {
		params.Type = "all"
	}
	if params.Page == 0 {
		params.Page = int(1)
	}
	if params.PerPage == 0 {
		params.PerPage = int(50)
	}

	// Make the request
	details, err := c.makePublicAPICall(
		"GET",
		fmt.Sprintf("/v2/checks?type=%s&page=%d&per_page=%d&muted=%t", params.Type, params.Page, params.PerPage, params.Muted),
		bytes.NewBufferString("{}"),
		nil)

	// Check for errors
	if err != nil {
		return nil, details, err
	}

	check, err := parseChecksResponse(details.ResponseBody)
	if err != nil {
		return check, details, err
	}

	return check, details, nil
}
