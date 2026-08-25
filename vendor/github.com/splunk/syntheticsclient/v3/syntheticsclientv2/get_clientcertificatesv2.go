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

func parseGetClientCertificatesV2Response(response string) (*ClientCertificatesV2Response, error) {
	var clientCertificatesV2 ClientCertificatesV2Response
	err := json.Unmarshal([]byte(response), &clientCertificatesV2)
	if err != nil {
		return nil, err
	}

	return &clientCertificatesV2, nil
}

func (c Client) GetClientCertificatesV2() (*ClientCertificatesV2Response, *RequestDetails, error) {
	details, err := c.makePublicAPICall("GET", "/certificates", bytes.NewBufferString("{}"), nil)
	if err != nil {
		return nil, details, err
	}

	clientCertificatesV2, err := parseGetClientCertificatesV2Response(details.ResponseBody)
	if err != nil {
		return clientCertificatesV2, details, err
	}

	return clientCertificatesV2, details, nil
}
