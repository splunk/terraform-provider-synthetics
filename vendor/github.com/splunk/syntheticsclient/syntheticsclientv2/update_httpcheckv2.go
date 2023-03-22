// Copyright 2021 Splunk, Inc.
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

func parseUpdateHttpCheckV2Response(response string) (*HttpCheckV2Response, error) {
	var updateHttpCheckV2 HttpCheckV2Response
	if response != "" {
		err := json.Unmarshal([]byte(response), &updateHttpCheckV2)
		if err != nil {
			return nil, err
		}
		return &updateHttpCheckV2, err
	}
	return &updateHttpCheckV2, nil
}

func (c Client) UpdateHttpCheckV2(id int, HttpCheckV2Details *HttpCheckV2Input) (*HttpCheckV2Response, *RequestDetails, error) {

	body, err := json.Marshal(HttpCheckV2Details)
	if err != nil {
		return nil, nil, err
	}

	requestDetails, err := c.makePublicAPICall("PUT", fmt.Sprintf("/tests/http/%d", id), bytes.NewBuffer(body), nil)
	if err != nil {
		return nil, requestDetails, err
	}

	updateHttpCheckV2, err := parseUpdateHttpCheckV2Response(requestDetails.ResponseBody)
	if err != nil {
		return updateHttpCheckV2, requestDetails, err
	}

	return updateHttpCheckV2, requestDetails, nil
}
