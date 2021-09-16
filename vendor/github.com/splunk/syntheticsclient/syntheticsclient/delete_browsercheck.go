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
