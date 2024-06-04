// Copyright 2024 Splunk, Inc.
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

func parseCreateDowntimeConfigurationV2Response(response string) (*DowntimeConfigurationV2Response, error) {

	var createDowntimeConfigurationV2 DowntimeConfigurationV2Response
	JSONResponse := []byte(response)
	err := json.Unmarshal(JSONResponse, &createDowntimeConfigurationV2)
	if err != nil {
		return nil, err
	}

	return &createDowntimeConfigurationV2, err
}

func (c Client) CreateDowntimeConfigurationV2(DowntimeConfigurationV2Details *DowntimeConfigurationV2Input) (*DowntimeConfigurationV2Response, *RequestDetails, error) {

	body, err := json.Marshal(DowntimeConfigurationV2Details)
	if err != nil {
		return nil, nil, err
	}

	details, err := c.makePublicAPICall("POST", "/downtime_configurations", bytes.NewBuffer(body), nil)
	if err != nil {
		return nil, details, err
	}

	newDowntimeConfigurationV2, err := parseCreateDowntimeConfigurationV2Response(details.ResponseBody)
	if err != nil {
		return newDowntimeConfigurationV2, details, err
	}

	return newDowntimeConfigurationV2, details, nil
}
