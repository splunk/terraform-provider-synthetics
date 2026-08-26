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

func parseCreateVariableV2Response(response string) (*VariableV2Response, error) {

	var createVariableV2 VariableV2Response
	JSONResponse := []byte(response)
	err := json.Unmarshal(JSONResponse, &createVariableV2)
	if err != nil {
		return nil, err
	}

	return &createVariableV2, err
}

func (c Client) CreateVariableV2(VariableV2Details *VariableV2Input) (*VariableV2Response, *RequestDetails, error) {

	body, err := json.Marshal(VariableV2Details)
	if err != nil {
		return nil, nil, err
	}

	details, err := c.makePublicAPICall("POST", "/variables", bytes.NewBuffer(body), nil)
	if err != nil {
		return nil, details, err
	}

	newVariableV2, err := parseCreateVariableV2Response(details.ResponseBody)
	if err != nil {
		return newVariableV2, details, err
	}

	return newVariableV2, details, nil
}
