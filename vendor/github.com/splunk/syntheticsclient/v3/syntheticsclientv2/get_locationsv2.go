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

func parseLocationsV2Response(response string) (*LocationsV2Response, error) {
	// Parse the response and return the locations object
	var locations LocationsV2Response
	err := json.Unmarshal([]byte(response), &locations)
	if err != nil {
		return nil, err
	}

	return &locations, err
}

func parseLocationV2Response(response string) (*LocationV2Response, error) {
	// Parse the response and return the locations object
	var location LocationV2Response
	err := json.Unmarshal([]byte(response), &location)
	if err != nil {
		return nil, err
	}

	return &location, err
}

func (c Client) GetLocationsV2() (*LocationsV2Response, *RequestDetails, error) {

	details, err := c.makePublicAPICall("GET",
		"/locations",
		bytes.NewBufferString("{}"),
		nil)

	if err != nil {
		return nil, details, err
	}

	locations, err := parseLocationsV2Response(details.ResponseBody)
	if err != nil {
		return locations, details, err
	}

	return locations, details, nil
}

func (c Client) GetLocationV2(id string) (*LocationV2Response, *RequestDetails, error) {

	details, err := c.makePublicAPICall("GET",
		fmt.Sprintf("/locations/%s", id),
		bytes.NewBufferString("{}"),
		nil)

	if err != nil {
		return nil, details, err
	}

	location, err := parseLocationV2Response(details.ResponseBody)
	if err != nil {
		return location, details, err
	}

	return location, details, nil
}
