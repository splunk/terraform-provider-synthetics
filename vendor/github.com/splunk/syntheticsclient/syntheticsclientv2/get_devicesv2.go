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

func parseDevicesV2Response(response string) (*DevicesV2Response, error) {
	// Parse the response and return the device object
	var device DevicesV2Response
	err := json.Unmarshal([]byte(response), &device)
	if err != nil {
		return nil, err
	}

	return &device, err
}

func (c Client) GetDevicesV2() (*DevicesV2Response, *RequestDetails, error) {

	details, err := c.makePublicAPICall("GET",
		"/devices",
		bytes.NewBufferString("{}"),
		nil)

	if err != nil {
		return nil, details, err
	}

	check, err := parseDevicesV2Response(details.ResponseBody)
	if err != nil {
		return check, details, err
	}

	return check, details, nil
}
