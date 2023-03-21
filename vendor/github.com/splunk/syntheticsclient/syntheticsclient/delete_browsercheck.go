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

package syntheticsclient

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func parseDeleteBrowserCheckResponse(response string) (*DeleteCheck, error) {
	var deleteBrowserCheck DeleteCheck
	err := json.Unmarshal([]byte(response), &deleteBrowserCheck)
	if err != nil {
		return nil, err
	}

	return &deleteBrowserCheck, err
}

func (c Client) DeleteBrowserCheck(id int) (*DeleteCheck, error) {
	requestDetails, err := c.makePublicAPICall("DELETE", fmt.Sprintf("/v2/checks/real_browsers/%d", id), bytes.NewBufferString("{}"), nil)
	if err != nil {
		return nil, err
	}

	deleteBrowserCheck, err := parseDeleteBrowserCheckResponse(requestDetails.ResponseBody)
	if err != nil {
		return deleteBrowserCheck, err
	}

	return deleteBrowserCheck, nil
}
