// Copyright 2026 Splunk, Inc.
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

// ValidateNewPortCheckV2 checks whether a port check payload would be
// accepted when creating a new test, without saving the test or triggering a
// run.
func (c Client) ValidateNewPortCheckV2(PortCheckV2Details *PortCheckV2Input) (*ValidateResponse, *RequestDetails, error) {
	body, err := json.Marshal(PortCheckV2Details)
	if err != nil {
		return nil, nil, err
	}

	details, err := c.makePublicAPICall("POST", "/tests/port/validate", bytes.NewBuffer(body), nil)
	if err != nil {
		return nil, details, err
	}

	validatePortCheckV2, err := parseValidateResponse(details.ResponseBody)
	if err != nil {
		return nil, details, err
	}

	return validatePortCheckV2, details, nil
}

// ValidatePortCheckV2 checks whether a port check payload would be accepted
// when updating the existing test identified by id, without saving the
// change or triggering a run.
func (c Client) ValidatePortCheckV2(id int, PortCheckV2Details *PortCheckV2Input) (*ValidateResponse, *RequestDetails, error) {
	body, err := json.Marshal(PortCheckV2Details)
	if err != nil {
		return nil, nil, err
	}

	details, err := c.makePublicAPICall("PUT", fmt.Sprintf("/tests/port/%d/validate", id), bytes.NewBuffer(body), nil)
	if err != nil {
		return nil, details, err
	}

	validatePortCheckV2, err := parseValidateResponse(details.ResponseBody)
	if err != nil {
		return nil, details, err
	}

	return validatePortCheckV2, details, nil
}
