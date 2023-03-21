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
)

func parseCreatePortCheckV2Response(response string) (*PortCheckV2Response, error) {

	var createPortCheckV2 PortCheckV2Response
	JSONResponse := []byte(response)
	err := json.Unmarshal(JSONResponse, &createPortCheckV2)
	if err != nil {
		return nil, err
	}

	return &createPortCheckV2, err
}

func (c Client) CreatePortCheckV2(PortCheckV2Details *PortCheckV2Input) (*PortCheckV2Response, *RequestDetails, error) {

	body, err := json.Marshal(PortCheckV2Details)
	if err != nil {
		return nil, nil, err
	}

	details, err := c.makePublicAPICall("POST", "/tests/port", bytes.NewBuffer(body), nil)
	if err != nil {
		return nil, details, err
	}

	newPortCheckV2, err := parseCreatePortCheckV2Response(details.ResponseBody)
	if err != nil {
		return newPortCheckV2, details, err
	}

	return newPortCheckV2, details, nil
}
