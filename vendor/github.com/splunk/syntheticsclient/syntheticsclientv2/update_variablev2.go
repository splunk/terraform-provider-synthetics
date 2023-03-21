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

func parseUpdateVariableV2Response(response string) (*VariableV2Response, error) {
	var updateVariableV2 VariableV2Response
	err := json.Unmarshal([]byte(response), &updateVariableV2)
	if err != nil {
		return nil, err
	}

	return &updateVariableV2, err
}

func (c Client) UpdateVariableV2(id int, VariableV2Details *VariableV2Input) (*VariableV2Response, *RequestDetails, error) {

	body, err := json.Marshal(VariableV2Details)
	if err != nil {
		return nil, nil, err
	}

	requestDetails, err := c.makePublicAPICall("PUT", fmt.Sprintf("/variables/%d", id), bytes.NewBuffer(body), nil)
	if err != nil {
		return nil, requestDetails, err
	}

	updateVariableV2, err := parseUpdateVariableV2Response(requestDetails.ResponseBody)
	if err != nil {
		return updateVariableV2, requestDetails, err
	}

	return updateVariableV2, requestDetails, nil
}
