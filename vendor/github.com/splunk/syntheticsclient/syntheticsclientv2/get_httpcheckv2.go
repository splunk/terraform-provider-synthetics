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

func parseGetHttpCheckV2Response(response string) (*HttpCheckV2Response, error) {
	// Parse the response and return the user object
	var HttpCheckV2 HttpCheckV2Response
	err := json.Unmarshal([]byte(response), &HttpCheckV2)
	if err != nil {
		return nil, err
	}

	return &HttpCheckV2, err
}

func (c Client) GetHttpCheckV2(id int) (*HttpCheckV2Response, *RequestDetails, error) {
	details, err := c.makePublicAPICall("GET", fmt.Sprintf("/tests/http/%d", id), bytes.NewBufferString("{}"), nil)

	// Check for errors
	if err != nil {
		return nil, details, err
	}

	HttpCheckV2, err := parseGetHttpCheckV2Response(details.ResponseBody)
	if err != nil {
		return HttpCheckV2, details, err
	}

	return HttpCheckV2, details, nil
}
