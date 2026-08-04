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

// ValidateNewSslCheckV2 checks whether an SSL check payload would be accepted
// when creating a new test, without saving the test or triggering a run.
func (c Client) ValidateNewSslCheckV2(SslCheckV2Details *SslCheckV2Input) (*ValidateResponse, *RequestDetails, error) {
	if SslCheckV2Details.Test.Validations == nil {
		SslCheckV2Details.Test.Validations = make([]Validations, 0)
	}

	body, err := json.Marshal(SslCheckV2Details)
	if err != nil {
		return nil, nil, err
	}

	details, err := c.makePublicAPICall("POST", "/tests/ssl/validate", bytes.NewBuffer(body), nil)
	if err != nil {
		return nil, details, err
	}

	validateSslCheckV2, err := parseValidateResponse(details.ResponseBody)
	if err != nil {
		return nil, details, err
	}

	return validateSslCheckV2, details, nil
}

// ValidateSslCheckV2 checks whether an SSL check payload would be accepted
// when updating the existing test identified by id, without saving the
// change or triggering a run.
func (c Client) ValidateSslCheckV2(id int, SslCheckV2Details *SslCheckV2UpdateInput) (*ValidateResponse, *RequestDetails, error) {
	body, err := json.Marshal(SslCheckV2Details)
	if err != nil {
		return nil, nil, err
	}

	details, err := c.makePublicAPICall("PUT", fmt.Sprintf("/tests/ssl/%d/validate", id), bytes.NewBuffer(body), nil)
	if err != nil {
		return nil, details, err
	}

	validateSslCheckV2, err := parseValidateResponse(details.ResponseBody)
	if err != nil {
		return nil, details, err
	}

	return validateSslCheckV2, details, nil
}
