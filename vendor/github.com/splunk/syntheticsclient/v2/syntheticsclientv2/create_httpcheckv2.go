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

func parseCreateHttpCheckV2Response(response string) (*HttpCheckV2Response, error) {

	var createHttpCheckV2 HttpCheckV2Response
	JSONResponse := []byte(response)
	err := json.Unmarshal(JSONResponse, &createHttpCheckV2)
	if err != nil {
		return nil, err
	}

	return &createHttpCheckV2, err
}

func (c Client) CreateHttpCheckV2(HttpCheckV2Details *HttpCheckV2Input) (*HttpCheckV2Response, *RequestDetails, error) {

	if HttpCheckV2Details.Test.Validations == nil {
		validation := make([]Validations, 0)
		HttpCheckV2Details.Test.Validations = validation
	}

	body, err := json.Marshal(HttpCheckV2Details)
	if err != nil {
		return nil, nil, err
	}

	details, err := c.makePublicAPICall("POST", "/tests/http", bytes.NewBuffer(body), nil)
	if err != nil {
		return nil, details, err
	}

	newHttpCheckV2, err := parseCreateHttpCheckV2Response(details.ResponseBody)
	if err != nil {
		return newHttpCheckV2, details, err
	}

	return newHttpCheckV2, details, nil
}
