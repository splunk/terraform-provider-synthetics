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
)

func parseCreateBrowserCheckResponse(response string) (*BrowserCheckResponse, error) {

	var createBrowserCheck BrowserCheckResponse
	JSONResponse := []byte(response)
	err := json.Unmarshal(JSONResponse, &createBrowserCheck)
	if err != nil {
		return nil, err
	}

	return &createBrowserCheck, err
}

func (c Client) CreateBrowserCheck(browserCheckDetails *BrowserCheckInput) (*BrowserCheckResponse, *RequestDetails, error) {

	body, err := json.Marshal(browserCheckDetails)
	if err != nil {
		return nil, nil, err
	}

	details, err := c.makePublicAPICall("POST", "/v2/checks/real_browsers", bytes.NewBuffer(body), nil)
	if err != nil {
		return nil, details, err
	}

	newBrowserCheck, err := parseCreateBrowserCheckResponse(details.ResponseBody)
	if err != nil {
		return newBrowserCheck, details, err
	}

	return newBrowserCheck, details, nil
}
