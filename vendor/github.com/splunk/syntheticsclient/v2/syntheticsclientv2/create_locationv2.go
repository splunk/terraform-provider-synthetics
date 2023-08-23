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

func parseCreateLocationV2Response(response string) (*LocationV2Response, error) {

	var createLocationV2 LocationV2Response
	JSONResponse := []byte(response)
	err := json.Unmarshal(JSONResponse, &createLocationV2)
	if err != nil {
		return nil, err
	}

	return &createLocationV2, err
}

func (c Client) CreateLocationV2(LocationV2Details *LocationV2Input) (*LocationV2Response, *RequestDetails, error) {

	body, err := json.Marshal(LocationV2Details)
	if err != nil {
		return nil, nil, err
	}

	details, err := c.makePublicAPICall("POST", "/locations", bytes.NewBuffer(body), nil)
	if err != nil {
		return nil, details, err
	}

	newLocationV2, err := parseCreateLocationV2Response(details.ResponseBody)
	if err != nil {
		return newLocationV2, details, err
	}

	return newLocationV2, details, nil
}
