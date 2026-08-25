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

func parseCreateSslCheckV2Response(response string) (*SslCheckV2Response, error) {
	var createSslCheckV2 SslCheckV2Response
	JSONResponse := []byte(response)
	err := json.Unmarshal(JSONResponse, &createSslCheckV2)
	if err != nil {
		return nil, err
	}

	return &createSslCheckV2, err
}

func (c Client) CreateSslCheckV2(SslCheckV2Details *SslCheckV2Input) (*SslCheckV2Response, *RequestDetails, error) {
	if SslCheckV2Details.Test.Validations == nil {
		validation := make([]Validations, 0)
		SslCheckV2Details.Test.Validations = validation
	}

	body, err := json.Marshal(SslCheckV2Details)
	if err != nil {
		return nil, nil, err
	}

	details, err := c.makePublicAPICall("POST", "/tests/ssl", bytes.NewBuffer(body), nil)
	if err != nil {
		return nil, details, err
	}

	newSslCheckV2, err := parseCreateSslCheckV2Response(details.ResponseBody)
	if err != nil {
		return newSslCheckV2, details, err
	}

	return newSslCheckV2, details, nil
}
