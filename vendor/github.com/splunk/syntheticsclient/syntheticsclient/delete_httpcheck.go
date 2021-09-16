package syntheticsclient

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func parseDeleteHttpCheckResponse(response string) (*DeleteCheck, error) {
	var deleteHttpCheck DeleteCheck
	err := json.Unmarshal([]byte(response), &deleteHttpCheck)
	if err != nil {
		return nil, err
	}

	return &deleteHttpCheck, err
}

func (c Client) DeleteHttpCheck(id int) (*DeleteCheck, error) {
	requestDetails, err := c.makePublicAPICall("DELETE", fmt.Sprintf("/v2/checks/http/%d", id), bytes.NewBufferString("{}"), nil)
	if err != nil {
		return nil, err
	}

	deleteHttpCheck, err := parseDeleteHttpCheckResponse(requestDetails.ResponseBody)
	if err != nil {
		return deleteHttpCheck, err
	}

	return deleteHttpCheck, nil
}
