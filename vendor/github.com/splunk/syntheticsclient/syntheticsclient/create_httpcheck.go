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

func parseCreateHttpCheckResponse(response string) (*HttpCheckResponse, error) {

	var createHttpCheck HttpCheckResponse
	JSONResponse := []byte(response)
	err := json.Unmarshal(JSONResponse, &createHttpCheck)
	if err != nil {
		return nil, err
	}

	return &createHttpCheck, err
}

func (c Client) CreateHttpCheck(httpCheckDetails *HttpCheckInput) (*HttpCheckResponse, *RequestDetails, error) {

	body, err := json.Marshal(httpCheckDetails)
	if err != nil {
		return nil, nil, err
	}

	details, err := c.makePublicAPICall("POST", "/v2/checks/http", bytes.NewBuffer(body), nil)
	if err != nil {
		return nil, details, err
	}

	newHttpCheck, err := parseCreateHttpCheckResponse(details.ResponseBody)
	if err != nil {
		return newHttpCheck, details, err
	}

	return newHttpCheck, details, nil
}
