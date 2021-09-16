package syntheticsclient

import (
	"bytes"
	"encoding/json"
)

func parseCreateBrowserCheckResponse(response string) (*BrowserCheckResponse, error) {

	var createBrowserCheck BrowserCheckResponse
	JSONResponse := []byte(response)
	err := json.Unmarshal(JSONResponse, &createBrowserCheck)
	if err != nil {
		return nil, err
	}

	return &createBrowserCheck, err
}

func (c Client) CreateBrowserCheck(browserCheckDetails *BrowserCheckInput) (*BrowserCheckResponse, *RequestDetails, error) {

	body, err := json.Marshal(browserCheckDetails)
	if err != nil {
		return nil, nil, err
	}

	details, err := c.makePublicAPICall("POST", "/v2/checks/real_browsers", bytes.NewBuffer(body), nil)
	if err != nil {
		return nil, details, err
	}

	newBrowserCheck, err := parseCreateBrowserCheckResponse(details.ResponseBody)
	if err != nil {
		return newBrowserCheck, details, err
	}

	return newBrowserCheck, details, nil
}
