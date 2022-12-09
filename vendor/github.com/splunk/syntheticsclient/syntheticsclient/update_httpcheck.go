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

func parseUpdateHttpCheckResponse(response string) (*HttpCheckResponse, error) {
	var updateHttpCheck HttpCheckResponse
	err := json.Unmarshal([]byte(response), &updateHttpCheck)
	if err != nil {
		return nil, err
	}

	return &updateHttpCheck, err
}

// CreateContact creates a new contact for a user
func (c Client) UpdateHttpCheck(id int, httpCheckDetails *HttpCheckInput) (*HttpCheckResponse, *RequestDetails, error) {

	body, err := json.Marshal(httpCheckDetails)
	if err != nil {
		return nil, nil, err
	}

	requestDetails, err := c.makePublicAPICall("PUT", fmt.Sprintf("/v2/checks/http/%d", id), bytes.NewBuffer(body), nil)
	if err != nil {
		return nil, requestDetails, err
	}

	updateHttpCheck, err := parseUpdateHttpCheckResponse(requestDetails.ResponseBody)
	if err != nil {
		return updateHttpCheck, requestDetails, err
	}

	return updateHttpCheck, requestDetails, nil
}
