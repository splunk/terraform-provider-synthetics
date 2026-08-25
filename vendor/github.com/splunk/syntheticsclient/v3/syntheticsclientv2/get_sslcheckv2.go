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

func parseGetSslCheckV2Response(response string) (*SslCheckV2Response, error) {
	var sslCheckV2 SslCheckV2Response
	err := json.Unmarshal([]byte(response), &sslCheckV2)
	if err != nil {
		return nil, err
	}

	return &sslCheckV2, err
}

func (c Client) GetSslCheckV2(id int) (*SslCheckV2Response, *RequestDetails, error) {
	details, err := c.makePublicAPICall("GET", fmt.Sprintf("/tests/ssl/%d", id), bytes.NewBufferString("{}"), nil)
	if err != nil {
		return nil, details, err
	}

	sslCheckV2, err := parseGetSslCheckV2Response(details.ResponseBody)
	if err != nil {
		return sslCheckV2, details, err
	}

	return sslCheckV2, details, nil
}
