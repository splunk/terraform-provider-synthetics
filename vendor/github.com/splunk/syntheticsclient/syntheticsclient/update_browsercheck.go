package syntheticsclient

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func parseUpdateBrowserCheckResponse(response string) (*BrowserCheckResponse, error) {
	var updateBrowserCheck BrowserCheckResponse
	err := json.Unmarshal([]byte(response), &updateBrowserCheck)
	if err != nil {
		return nil, err
	}

	return &updateBrowserCheck, err
}

func (c Client) UpdateBrowserCheck(id int, browserCheckDetails *BrowserCheckInput) (*BrowserCheckResponse, *RequestDetails, error) {

	body, err := json.Marshal(browserCheckDetails)
	if err != nil {
		return nil, nil, err
	}

	requestDetails, err := c.makePublicAPICall("PUT", fmt.Sprintf("/v2/checks/real_browsers/%d", id), bytes.NewBuffer(body), nil)
	if err != nil {
		return nil, requestDetails, err
	}

	updateBrowserCheck, err := parseUpdateBrowserCheckResponse(requestDetails.ResponseBody)
	if err != nil {
		return updateBrowserCheck, requestDetails, err
	}

	return updateBrowserCheck, requestDetails, nil
}
