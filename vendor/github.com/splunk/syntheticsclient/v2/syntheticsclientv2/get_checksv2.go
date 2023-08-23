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
	"bytes"
	"encoding/json"
	"fmt"
)

func parseChecksV2Response(response string) (*ChecksV2Response, error) {
	// Parse the response and return the check object
	var checks ChecksV2Response
	err := json.Unmarshal([]byte(response), &checks)
	if err != nil {
		return nil, err
	}

	return &checks, err
}

// GetChecks returns all checks
func (c Client) GetChecksV2(params *GetChecksV2Options) (*ChecksV2Response, *RequestDetails, error) {
	// Check for default params
	if params.TestType == "" {
		params.TestType = ""
	}
	if params.Search == "" {
		params.Search = ""
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
		fmt.Sprintf("/tests?testType=%s&page=%d&perPage=%d&orderBy=%s&search=%s", params.TestType, params.Page, params.PerPage, params.OrderBy, params.Search),
		bytes.NewBufferString("{}"),
		nil)

	// Check for errors
	if err != nil {
		return nil, details, err
	}

	check, err := parseChecksV2Response(details.ResponseBody)
	if err != nil {
		return check, details, err
	}

	return check, details, nil
}
