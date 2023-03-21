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

func parseUpdateApiCheckV2Response(response string) (*ApiCheckV2Response, error) {
	var updateApiCheckV2 ApiCheckV2Response
	err := json.Unmarshal([]byte(response), &updateApiCheckV2)
	if err != nil {
		return nil, err
	}

	return &updateApiCheckV2, err
}

func (c Client) UpdateApiCheckV2(id int, ApiCheckV2Details *ApiCheckV2Input) (*ApiCheckV2Response, *RequestDetails, error) {

	body, err := json.Marshal(ApiCheckV2Details)
	if err != nil {
		return nil, nil, err
	}

	requestDetails, err := c.makePublicAPICall("PUT", fmt.Sprintf("/tests/api/%d", id), bytes.NewBuffer(body), nil)
	if err != nil {
		return nil, requestDetails, err
	}

	updateApiCheckV2, err := parseUpdateApiCheckV2Response(requestDetails.ResponseBody)
	if err != nil {
		return updateApiCheckV2, requestDetails, err
	}

	return updateApiCheckV2, requestDetails, nil
}
