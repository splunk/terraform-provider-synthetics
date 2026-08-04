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

// ValidateNewHttpCheckV2 checks whether an HTTP check payload would be
// accepted when creating a new test, without saving the test or triggering a
// run.
func (c Client) ValidateNewHttpCheckV2(HttpCheckV2Details *HttpCheckV2Input) (*ValidateResponse, *RequestDetails, error) {
	if HttpCheckV2Details.Test.Validations == nil {
		HttpCheckV2Details.Test.Validations = make([]Validations, 0)
	}

	body, err := json.Marshal(HttpCheckV2Details)
	if err != nil {
		return nil, nil, err
	}

	details, err := c.makePublicAPICall("POST", "/tests/http/validate", bytes.NewBuffer(body), nil)
	if err != nil {
		return nil, details, err
	}

	validateHttpCheckV2, err := parseValidateResponse(details.ResponseBody)
	if err != nil {
		return nil, details, err
	}

	return validateHttpCheckV2, details, nil
}

// ValidateHttpCheckV2 checks whether an HTTP check payload would be accepted
// when updating the existing test identified by id, without saving the
// change or triggering a run.
func (c Client) ValidateHttpCheckV2(id int, HttpCheckV2Details *HttpCheckV2Input) (*ValidateResponse, *RequestDetails, error) {
	if HttpCheckV2Details.Test.Validations == nil {
		HttpCheckV2Details.Test.Validations = make([]Validations, 0)
	}

	body, err := json.Marshal(HttpCheckV2Details)
	if err != nil {
		return nil, nil, err
	}

	details, err := c.makePublicAPICall("PUT", fmt.Sprintf("/tests/http/%d/validate", id), bytes.NewBuffer(body), nil)
	if err != nil {
		return nil, details, err
	}

	validateHttpCheckV2, err := parseValidateResponse(details.ResponseBody)
	if err != nil {
		return nil, details, err
	}

	return validateHttpCheckV2, details, nil
}

// ValidateNewHttpCheckV2WithNullablePort checks whether an HTTP check payload
// with a nullable port would be accepted when creating a new test, without
// saving the test or triggering a run.
func (c Client) ValidateNewHttpCheckV2WithNullablePort(HttpCheckV2Details *HttpCheckV2InputWithNullablePort) (*ValidateResponse, *RequestDetails, error) {
	if HttpCheckV2Details.Test.Validations == nil {
		HttpCheckV2Details.Test.Validations = make([]Validations, 0)
	}

	body, err := json.Marshal(HttpCheckV2Details)
	if err != nil {
		return nil, nil, err
	}

	details, err := c.makePublicAPICall("POST", "/tests/http/validate", bytes.NewBuffer(body), nil)
	if err != nil {
		return nil, details, err
	}

	validateHttpCheckV2, err := parseValidateResponse(details.ResponseBody)
	if err != nil {
		return nil, details, err
	}

	return validateHttpCheckV2, details, nil
}

// ValidateHttpCheckV2WithNullablePort checks whether an HTTP check payload
// with a nullable port would be accepted when updating the existing test
// identified by id, without saving the change or triggering a run.
func (c Client) ValidateHttpCheckV2WithNullablePort(id int, HttpCheckV2Details *HttpCheckV2InputWithNullablePort) (*ValidateResponse, *RequestDetails, error) {
	if HttpCheckV2Details.Test.Validations == nil {
		HttpCheckV2Details.Test.Validations = make([]Validations, 0)
	}

	body, err := json.Marshal(HttpCheckV2Details)
	if err != nil {
		return nil, nil, err
	}

	details, err := c.makePublicAPICall("PUT", fmt.Sprintf("/tests/http/%d/validate", id), bytes.NewBuffer(body), nil)
	if err != nil {
		return nil, details, err
	}

	validateHttpCheckV2, err := parseValidateResponse(details.ResponseBody)
	if err != nil {
		return nil, details, err
	}

	return validateHttpCheckV2, details, nil
}
