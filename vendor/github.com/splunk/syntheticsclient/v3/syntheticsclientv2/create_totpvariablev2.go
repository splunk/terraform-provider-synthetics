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
)

func parseCreateTotpVariableV2Response(response string) (*TotpVariableV2Response, error) {
	var createTotpVariableV2 TotpVariableV2Response
	err := json.Unmarshal([]byte(response), &createTotpVariableV2)
	if err != nil {
		return nil, err
	}

	return &createTotpVariableV2, err
}

func (c Client) CreateTotpVariableV2(TotpVariableV2Details *TotpVariableV2Input) (*TotpVariableV2Response, *RequestDetails, error) {
	body, err := json.Marshal(TotpVariableV2Details)
	if err != nil {
		return nil, nil, err
	}

	details, err := c.makePublicAPICall("POST", "/totps", bytes.NewBuffer(body), nil)
	scrubTotpVariableRequestDetails(details, totpVariableCreateSecret(TotpVariableV2Details))
	if err != nil {
		return nil, details, err
	}

	newTotpVariableV2, err := parseCreateTotpVariableV2Response(details.ResponseBody)
	if err != nil {
		return newTotpVariableV2, details, err
	}

	return newTotpVariableV2, details, nil
}

func scrubTotpVariableRequestDetails(details *RequestDetails, sensitiveValues ...string) {
	if details == nil {
		return
	}
	for _, sensitiveValue := range sensitiveValues {
		details.RequestBody = redactSensitiveValue(details.RequestBody, sensitiveValue)
	}
	details.RawRequest = nil
}

func totpVariableCreateSecret(details *TotpVariableV2Input) string {
	if details == nil {
		return ""
	}

	return details.Totp.Secret
}
