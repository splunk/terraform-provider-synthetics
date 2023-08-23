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

func parseUpdateBrowserCheckV2Response(response string) (*BrowserCheckV2Response, error) {
	var updateBrowserCheckV2 BrowserCheckV2Response
	if response != "" {
		err := json.Unmarshal([]byte(response), &updateBrowserCheckV2)
		if err != nil {
			return nil, err
		}
		return &updateBrowserCheckV2, err
	}
	return &updateBrowserCheckV2, nil
}

func (c Client) UpdateBrowserCheckV2(id int, BrowserCheckV2Details *BrowserCheckV2Input) (*BrowserCheckV2Response, *RequestDetails, error) {

	body, err := json.Marshal(BrowserCheckV2Details)
	if err != nil {
		return nil, nil, err
	}

	requestDetails, err := c.makePublicAPICall("PUT", fmt.Sprintf("/tests/browser/%d", id), bytes.NewBuffer(body), nil)
	if err != nil {
		return nil, requestDetails, err
	}

	updateBrowserCheckV2, err := parseUpdateBrowserCheckV2Response(requestDetails.ResponseBody)
	if err != nil {
		return updateBrowserCheckV2, requestDetails, err
	}

	return updateBrowserCheckV2, requestDetails, nil
}
