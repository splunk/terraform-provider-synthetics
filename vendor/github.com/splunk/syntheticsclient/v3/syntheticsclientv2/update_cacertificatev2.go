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

func parseUpdateCaCertificateV2Response(response string) (*CaCertificateV2Response, error) {
	var updateCaCertificateV2 CaCertificateV2Response
	if response != "" {
		err := json.Unmarshal([]byte(response), &updateCaCertificateV2)
		if err != nil {
			return nil, err
		}
		return &updateCaCertificateV2, err
	}
	return &updateCaCertificateV2, nil
}

func (c Client) UpdateCaCertificateV2(id int, CaCertificateV2Details *CaCertificateV2UpdateInput) (*CaCertificateV2Response, *RequestDetails, error) {
	body, err := json.Marshal(CaCertificateV2Details)
	if err != nil {
		return nil, nil, err
	}

	requestDetails, err := c.makePublicAPICall("PUT", fmt.Sprintf("/cacerts/%d", id), bytes.NewBuffer(body), nil)
	if err != nil {
		return nil, requestDetails, err
	}

	updateCaCertificateV2, err := parseUpdateCaCertificateV2Response(requestDetails.ResponseBody)
	if err != nil {
		return updateCaCertificateV2, requestDetails, err
	}

	return updateCaCertificateV2, requestDetails, nil
}
