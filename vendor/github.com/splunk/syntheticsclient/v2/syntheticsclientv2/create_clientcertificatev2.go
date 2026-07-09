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

func parseCreateClientCertificateV2Response(response string) (*ClientCertificateV2Response, error) {
	var createClientCertificateV2 ClientCertificateV2Response
	err := json.Unmarshal([]byte(response), &createClientCertificateV2)
	if err != nil {
		return nil, err
	}

	return &createClientCertificateV2, nil
}

func (c Client) CreateClientCertificateV2(input *ClientCertificateV2Input) (*ClientCertificateV2Response, *RequestDetails, error) {
	body, err := json.Marshal(input)
	if err != nil {
		return nil, nil, err
	}

	details, err := c.makePublicAPICall("POST", "/certificates", bytes.NewBuffer(body), nil)
	if err != nil {
		return nil, details, err
	}

	newClientCertificateV2, err := parseCreateClientCertificateV2Response(details.ResponseBody)
	if err != nil {
		return newClientCertificateV2, details, err
	}

	return newClientCertificateV2, details, nil
}
