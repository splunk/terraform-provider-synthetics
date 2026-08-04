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
	"encoding/json"
	"strings"
)

// ValidateResponse is the response returned by the Synthetics API's test
// validate endpoints (POST .../validate for new test definitions and
// PUT/PATCH .../{id}/validate for existing tests). Validate checks whether a
// test payload would be accepted by the API without saving any changes or
// triggering a test run.
//
// The API always responds HTTP 200 for validate calls, using Valid to signal
// whether the payload passed. Details is a JSON array when Valid is true and
// a JSON object of per-attribute error messages when Valid is false; use
// FieldErrors to normalize both shapes.
type ValidateResponse struct {
	Valid   bool            `json:"valid"`
	Message string          `json:"message"`
	Details json.RawMessage `json:"details"`
}

func parseValidateResponse(response string) (*ValidateResponse, error) {
	var validateResponse ValidateResponse
	if response == "" {
		return &validateResponse, nil
	}

	if err := json.Unmarshal([]byte(response), &validateResponse); err != nil {
		return nil, err
	}

	return &validateResponse, nil
}

// FieldErrors returns field-level validation error messages keyed by
// attribute name. It returns an empty map when validation succeeded, since
// the API responds with an empty array rather than an object in that case.
func (v *ValidateResponse) FieldErrors() (map[string][]string, error) {
	fieldErrors := map[string][]string{}

	trimmed := strings.TrimSpace(string(v.Details))
	if trimmed == "" || trimmed == "null" || strings.HasPrefix(trimmed, "[") {
		return fieldErrors, nil
	}

	if err := json.Unmarshal(v.Details, &fieldErrors); err != nil {
		return nil, err
	}

	return fieldErrors, nil
}
