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
)

func parseCreateBrowserCheckV2Response(response string) (*BrowserCheckV2Response, error) {

	var createBrowserCheckV2 BrowserCheckV2Response
	JSONResponse := []byte(response)
	err := json.Unmarshal(JSONResponse, &createBrowserCheckV2)
	if err != nil {
		return nil, err
	}

	return &createBrowserCheckV2, err
}

func (c Client) CreateBrowserCheckV2(browserCheckV2Details *BrowserCheckV2Input) (*BrowserCheckV2Response, *RequestDetails, error) {

	body, err := json.Marshal(browserCheckV2Details)
	if err != nil {
		return nil, nil, err
	}

	details, err := c.makePublicAPICall("POST", "/tests/browser", bytes.NewBuffer(body), nil)
	if err != nil {
		return nil, details, err
	}

	newBrowserCheckV2, err := parseCreateBrowserCheckV2Response(details.ResponseBody)
	if err != nil {
		return newBrowserCheckV2, details, err
	}

	return newBrowserCheckV2, details, nil
}
