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

func parseCreateCaCertificateV2Response(response string) (*CaCertificateV2Response, error) {
	var createCaCertificateV2 CaCertificateV2Response
	JSONResponse := []byte(response)
	err := json.Unmarshal(JSONResponse, &createCaCertificateV2)
	if err != nil {
		return nil, err
	}

	return &createCaCertificateV2, err
}

func (c Client) CreateCaCertificateV2(CaCertificateV2Details *CaCertificateV2Input) (*CaCertificateV2Response, *RequestDetails, error) {
	body, err := json.Marshal(CaCertificateV2Details)
	if err != nil {
		return nil, nil, err
	}

	details, err := c.makePublicAPICall("POST", "/cacerts", bytes.NewBuffer(body), nil)
	if err != nil {
		return nil, details, err
	}

	newCaCertificateV2, err := parseCreateCaCertificateV2Response(details.ResponseBody)
	if err != nil {
		return newCaCertificateV2, details, err
	}

	return newCaCertificateV2, details, nil
}
