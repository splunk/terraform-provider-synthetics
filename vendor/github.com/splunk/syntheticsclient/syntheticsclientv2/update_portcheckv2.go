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

func parseUpdatePortCheckV2Response(response string) (*PortCheckV2Response, error) {
	var updatePortCheckV2 PortCheckV2Response
	if response != "" {
		err := json.Unmarshal([]byte(response), &updatePortCheckV2)
		if err != nil {
			return nil, err
		}
		return &updatePortCheckV2, err
	}
	return &updatePortCheckV2, nil
}

func (c Client) UpdatePortCheckV2(id int, PortCheckV2Details *PortCheckV2Input) (*PortCheckV2Response, *RequestDetails, error) {

	body, err := json.Marshal(PortCheckV2Details)
	if err != nil {
		return nil, nil, err
	}

	requestDetails, err := c.makePublicAPICall("PUT", fmt.Sprintf("/tests/port/%d", id), bytes.NewBuffer(body), nil)
	if err != nil {
		return nil, requestDetails, err
	}

	updatePortCheckV2, err := parseUpdatePortCheckV2Response(requestDetails.ResponseBody)
	if err != nil {
		return updatePortCheckV2, requestDetails, err
	}

	return updatePortCheckV2, requestDetails, nil
}
