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

func parseVariableV2Response(response string) (*VariableV2Response, error) {
	// Parse the response and return the check object
	var check VariableV2Response
	err := json.Unmarshal([]byte(response), &check)
	if err != nil {
		return nil, err
	}

	return &check, err
}

func (c Client) GetVariableV2(id int) (*VariableV2Response, *RequestDetails, error) {

	details, err := c.makePublicAPICall("GET",
		fmt.Sprintf("/variables/%d", id),
		bytes.NewBufferString("{}"),
		nil)

	if err != nil {
		return nil, details, err
	}

	check, err := parseVariableV2Response(details.ResponseBody)
	if err != nil {
		return check, details, err
	}

	return check, details, nil
}

func parseVariablesV2Response(response string) (*VariablesV2Response, error) {
	// Parse the response and return the check object
	var check VariablesV2Response
	err := json.Unmarshal([]byte(response), &check)
	if err != nil {
		return nil, err
	}

	return &check, err
}

func (c Client) GetVariablesV2() (*VariablesV2Response, *RequestDetails, error) {

	details, err := c.makePublicAPICall("GET",
		"/variables",
		bytes.NewBufferString("{}"),
		nil)

	if err != nil {
		return nil, details, err
	}

	check, err := parseVariablesV2Response(details.ResponseBody)
	if err != nil {
		return check, details, err
	}

	return check, details, nil
}
