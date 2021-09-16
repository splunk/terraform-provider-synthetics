package syntheticsclient

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func parseUpdateHttpCheckResponse(response string) (*HttpCheckResponse, error) {
	var updateHttpCheck HttpCheckResponse
	err := json.Unmarshal([]byte(response), &updateHttpCheck)
	if err != nil {
		return nil, err
	}

	return &updateHttpCheck, err
}

// CreateContact creates a new contact for a user
func (c Client) UpdateHttpCheck(id int, httpCheckDetails *HttpCheckInput) (*HttpCheckResponse, *RequestDetails, error) {

	body, err := json.Marshal(httpCheckDetails)
	if err != nil {
		return nil, nil, err
	}

	requestDetails, err := c.makePublicAPICall("PUT", fmt.Sprintf("/v2/checks/http/%d", id), bytes.NewBuffer(body), nil)
	if err != nil {
		return nil, requestDetails, err
	}

	updateHttpCheck, err := parseUpdateHttpCheckResponse(requestDetails.ResponseBody)
	if err != nil {
		return updateHttpCheck, requestDetails, err
	}

	return updateHttpCheck, requestDetails, nil
}
