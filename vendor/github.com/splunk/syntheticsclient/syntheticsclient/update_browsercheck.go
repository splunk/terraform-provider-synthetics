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

func parseUpdateBrowserCheckResponse(response string) (*BrowserCheckResponse, error) {
	var updateBrowserCheck BrowserCheckResponse
	err := json.Unmarshal([]byte(response), &updateBrowserCheck)
	if err != nil {
		return nil, err
	}

	return &updateBrowserCheck, err
}

func (c Client) UpdateBrowserCheck(id int, browserCheckDetails *BrowserCheckInput) (*BrowserCheckResponse, *RequestDetails, error) {

	body, err := json.Marshal(browserCheckDetails)
	if err != nil {
		return nil, nil, err
	}

	requestDetails, err := c.makePublicAPICall("PUT", fmt.Sprintf("/v2/checks/real_browsers/%d", id), bytes.NewBuffer(body), nil)
	if err != nil {
		return nil, requestDetails, err
	}

	updateBrowserCheck, err := parseUpdateBrowserCheckResponse(requestDetails.ResponseBody)
	if err != nil {
		return updateBrowserCheck, requestDetails, err
	}

	return updateBrowserCheck, requestDetails, nil
}
