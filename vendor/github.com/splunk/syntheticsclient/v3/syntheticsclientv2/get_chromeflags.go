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

func parseChromeFlagsResponse(response string) (*ChromeFlagsResponse, error) {
	var chromeFlags ChromeFlagsResponse
	err := json.Unmarshal([]byte(response), &chromeFlags)
	if err != nil {
		return nil, err
	}

	return &chromeFlags, err
}

func (c Client) GetChromeFlags() (*ChromeFlagsResponse, *RequestDetails, error) {
	details, err := c.makePublicAPICall("GET",
		"/chrome_flags",
		bytes.NewBufferString("{}"),
		nil)

	if err != nil {
		return nil, details, err
	}

	chromeFlags, err := parseChromeFlagsResponse(details.ResponseBody)
	if err != nil {
		return chromeFlags, details, err
	}

	return chromeFlags, details, nil
}
