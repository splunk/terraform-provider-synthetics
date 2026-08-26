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

func parseUpdateSslCheckV2Response(response string) (*SslCheckV2Response, error) {
	var updateSslCheckV2 SslCheckV2Response
	if response != "" {
		err := json.Unmarshal([]byte(response), &updateSslCheckV2)
		if err != nil {
			return nil, err
		}
		return &updateSslCheckV2, err
	}
	return &updateSslCheckV2, nil
}

func (c Client) UpdateSslCheckV2(id int, SslCheckV2Details *SslCheckV2UpdateInput) (*SslCheckV2Response, *RequestDetails, error) {
	body, err := json.Marshal(SslCheckV2Details)
	if err != nil {
		return nil, nil, err
	}

	requestDetails, err := c.makePublicAPICall("PUT", fmt.Sprintf("/tests/ssl/%d", id), bytes.NewBuffer(body), nil)
	if err != nil {
		return nil, requestDetails, err
	}

	updateSslCheckV2, err := parseUpdateSslCheckV2Response(requestDetails.ResponseBody)
	if err != nil {
		return updateSslCheckV2, requestDetails, err
	}

	return updateSslCheckV2, requestDetails, nil
}
