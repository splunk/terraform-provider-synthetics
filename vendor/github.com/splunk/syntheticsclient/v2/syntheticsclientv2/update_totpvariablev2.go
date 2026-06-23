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

func parseUpdateTotpVariableV2Response(response string) (*TotpVariableV2Response, error) {
	var updateTotpVariableV2 TotpVariableV2Response
	if response != "" {
		err := json.Unmarshal([]byte(response), &updateTotpVariableV2)
		if err != nil {
			return nil, err
		}
		return &updateTotpVariableV2, err
	}
	return &updateTotpVariableV2, nil
}

func (c Client) UpdateTotpVariableV2(id int, TotpVariableV2Details *TotpVariableV2UpdateInput) (*TotpVariableV2Response, *RequestDetails, error) {
	body, err := json.Marshal(TotpVariableV2Details)
	if err != nil {
		return nil, nil, err
	}

	requestDetails, err := c.makePublicAPICall("PUT", fmt.Sprintf("/totps/%d", id), bytes.NewBuffer(body), nil)
	scrubTotpVariableRequestDetails(requestDetails, totpVariableUpdateSecret(TotpVariableV2Details))
	if err != nil {
		return nil, requestDetails, err
	}

	updateTotpVariableV2, err := parseUpdateTotpVariableV2Response(requestDetails.ResponseBody)
	if err != nil {
		return updateTotpVariableV2, requestDetails, err
	}

	return updateTotpVariableV2, requestDetails, nil
}

func totpVariableUpdateSecret(details *TotpVariableV2UpdateInput) string {
	if details == nil || details.Totp.Secret == nil {
		return ""
	}

	return *details.Totp.Secret
}
