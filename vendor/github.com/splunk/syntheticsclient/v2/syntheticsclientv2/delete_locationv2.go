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
	"errors"
	"fmt"
	"strconv"
)

func (c Client) DeleteLocationV2(id string) (int, error) {
	requestDetails, err := c.makePublicAPICall("DELETE", fmt.Sprintf("/locations/%s", id), bytes.NewBufferString("{}"), nil)
	if err != nil {
		return 1, err
	}
	var status = requestDetails.StatusCode

	fmt.Println(status)

	if status >= 300 || status < 200 {
		errorMsg := fmt.Sprintf("error: Response code %v. Expecting 2XX.", strconv.Itoa(status))
		return status, errors.New(errorMsg)
	}

	return status, err
}
