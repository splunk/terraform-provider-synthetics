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

func parseUpdateClientCertificateV2Response(response string) (*ClientCertificateV2Response, error) {
	var updateClientCertificateV2 ClientCertificateV2Response
	if response != "" {
		err := json.Unmarshal([]byte(response), &updateClientCertificateV2)
		if err != nil {
			return nil, err
		}
		return &updateClientCertificateV2, nil
	}
	return &updateClientCertificateV2, nil
}

func (c Client) UpdateClientCertificateV2(id int, input *ClientCertificateV2UpdateInput) (*ClientCertificateV2Response, *RequestDetails, error) {
	body, err := json.Marshal(input)
	if err != nil {
		return nil, nil, err
	}

	details, err := c.makePublicAPICall("PUT", fmt.Sprintf("/certificates/%d", id), bytes.NewBuffer(body), nil)
	if err != nil {
		return nil, details, err
	}

	updateClientCertificateV2, err := parseUpdateClientCertificateV2Response(details.ResponseBody)
	if err != nil {
		return updateClientCertificateV2, details, err
	}

	return updateClientCertificateV2, details, nil
}
