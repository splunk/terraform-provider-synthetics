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

// ValidateNewApiCheckV2 checks whether an API check payload would be accepted
// when creating a new test, without saving the test or triggering a run.
func (c Client) ValidateNewApiCheckV2(ApiCheckV2Details *ApiCheckV2Input) (*ValidateResponse, *RequestDetails, error) {
	if ApiCheckV2Details.Test.Requests[0].Setup == nil {
		ApiCheckV2Details.Test.Requests[0].Setup = make([]Setup, 0)
	}

	if ApiCheckV2Details.Test.Requests[0].Validations == nil {
		ApiCheckV2Details.Test.Requests[0].Validations = make([]Validations, 0)
	}

	body, err := json.Marshal(ApiCheckV2Details)
	if err != nil {
		return nil, nil, err
	}

	details, err := c.makePublicAPICall("POST", "/v2/tests/api/validate", bytes.NewBuffer(body), nil)
	if err != nil {
		return nil, details, err
	}

	validateApiCheckV2, err := parseValidateResponse(details.ResponseBody)
	if err != nil {
		return nil, details, err
	}

	return validateApiCheckV2, details, nil
}

// ValidateApiCheckV2 checks whether an API check payload would be accepted
// when updating the existing test identified by id, without saving the
// change or triggering a run.
func (c Client) ValidateApiCheckV2(id int, ApiCheckV2Details *ApiCheckV2Input) (*ValidateResponse, *RequestDetails, error) {
	if ApiCheckV2Details.Test.Requests[0].Setup == nil {
		ApiCheckV2Details.Test.Requests[0].Setup = make([]Setup, 0)
	}

	if ApiCheckV2Details.Test.Requests[0].Validations == nil {
		ApiCheckV2Details.Test.Requests[0].Validations = make([]Validations, 0)
	}

	body, err := json.Marshal(ApiCheckV2Details)
	if err != nil {
		return nil, nil, err
	}

	details, err := c.makePublicAPICall("PUT", fmt.Sprintf("/v2/tests/api/%d/validate", id), bytes.NewBuffer(body), nil)
	if err != nil {
		return nil, details, err
	}

	validateApiCheckV2, err := parseValidateResponse(details.ResponseBody)
	if err != nil {
		return nil, details, err
	}

	return validateApiCheckV2, details, nil
}
