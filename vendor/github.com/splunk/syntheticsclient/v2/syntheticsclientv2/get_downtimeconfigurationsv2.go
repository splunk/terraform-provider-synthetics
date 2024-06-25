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
	"fmt"
	"strings"
)

func parseDowntimeConfigurationV2Response(response string) (*DowntimeConfigurationV2Response, error) {
	// Parse the response and return the check object
	var check DowntimeConfigurationV2Response
	err := json.Unmarshal([]byte(response), &check)
	if err != nil {
		return nil, err
	}

	return &check, err
}

func (c Client) GetDowntimeConfigurationV2(id int) (*DowntimeConfigurationV2Response, *RequestDetails, error) {

	details, err := c.makePublicAPICall("GET",
		fmt.Sprintf("/downtime_configurations/%d", id),
		bytes.NewBufferString("{}"),
		nil)

	if err != nil {
		return nil, details, err
	}

	check, err := parseDowntimeConfigurationV2Response(details.ResponseBody)
	if err != nil {
		return check, details, err
	}

	return check, details, nil
}

func parseDowntimeConfigurationsV2Response(response string) (*DowntimeConfigurationsV2Response, error) {
	// Parse the response and return the check object
	var check DowntimeConfigurationsV2Response
	err := json.Unmarshal([]byte(response), &check)
	if err != nil {
		return nil, err
	}

	return &check, err
}

func (c Client) GetDowntimeConfigurationsV2(params *GetDowntimeConfigurationsV2Options) (*DowntimeConfigurationsV2Response, *RequestDetails, error) {
	// Check for default params
	if params.Search == "" {
		params.Search = ""
	}
	if params.Page == 0 {
		params.Page = int(1)
	}
	if params.PerPage == 0 {
		params.PerPage = int(50)
	}
	details, err := c.makePublicAPICall(
		"GET",
		fmt.Sprintf("/downtime_configurations?page=%d&perPage=%d&orderBy=%s&search=%s%s%s",
			params.Page,
			params.PerPage,
			params.OrderBy,
			params.Search,
			dcStringsQueryParam(params.Rule, "&rule[]="),
			dcStringsQueryParam(params.Status, "&status[]="),
		),
		bytes.NewBufferString("{}"),
		nil)

	if err != nil {
		return nil, details, err
	}

	check, err := parseDowntimeConfigurationsV2Response(details.ResponseBody)
	if err != nil {
		return check, details, err
	}

	return check, details, nil
}

func dcStringsQueryParam(params []string, queryParamName string) string {
	if len(params) == 0 {
		return ""
	}
	return queryParamName + strings.Join(params, queryParamName)
}
