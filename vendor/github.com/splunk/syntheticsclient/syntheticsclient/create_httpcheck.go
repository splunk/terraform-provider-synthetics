package syntheticsclient

import (
	"bytes"
	"encoding/json"
)

func parseCreateHttpCheckResponse(response string) (*HttpCheckResponse, error) {

	var createHttpCheck HttpCheckResponse
	JSONResponse := []byte(response)
	err := json.Unmarshal(JSONResponse, &createHttpCheck)
	if err != nil {
		return nil, err
	}

	return &createHttpCheck, err
}

func (c Client) CreateHttpCheck(httpCheckDetails *HttpCheckInput) (*HttpCheckResponse, *RequestDetails, error) {

	body, err := json.Marshal(httpCheckDetails)
	if err != nil {
		return nil, nil, err
	}

	details, err := c.makePublicAPICall("POST", "/v2/checks/http", bytes.NewBuffer(body), nil)
	if err != nil {
		return nil, details, err
	}

	newHttpCheck, err := parseCreateHttpCheckResponse(details.ResponseBody)
	if err != nil {
		return newHttpCheck, details, err
	}

	return newHttpCheck, details, nil
}
