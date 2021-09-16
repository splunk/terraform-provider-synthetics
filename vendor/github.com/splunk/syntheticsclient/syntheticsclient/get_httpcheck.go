package syntheticsclient

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type GetHttpCheck struct {
	ID                              int
	Name                            string             `json:"name"`
	Type                            string             `json:"type"`
	Frequency                       int                `json:"frequency,omitempty"`
	Paused                          bool               `json:"paused,omitempty"`
	Muted                           bool               `json:"muted,omitempty"`
	CreatedAt                       string             `json:"created_at,omitempty"`
	UpdatedAt                       string             `json:"updated_at,omitempty"`
	Links                           Links              `json:"links,omitempty"`
	Status                          Status             `json:"status,omitempty"`
	Notifications                   Notifications      `json:"notifications,omitempty"`
	ResponseTimeMonitorMilliseconds int                `json:"response_time_monitor_milliseconds,omitempty"`
	HTTPRequestHeaders              HTTPRequestHeaders `json:"http_request_headers,omitempty"`
	HTTPRequestBody                 string             `json:"http_request_body,omitempty"`
	HTTPMethod                      string             `json:"http_method,omitempty"`
	RoundRobin                      bool               `json:"round_robin,omitempty"`
	AutoRetry                       bool               `json:"auto_retry,omitempty"`
	Enabled                         bool               `json:"enabled,omitempty"`
	Integrations                    Integrations       `json:"integrations,omitempty"`
	URL                             string             `json:"url,omitempty"`
	UserAgent                       string             `json:"user_agent,omitempty"`
	Tags                            Tags               `json:"tags,omitempty"`
	BlackoutPeriods                 BlackoutPeriods    `json:"blackout_periods,omitempty"`
	Locations                       Locations          `json:"locations,omitempty"`
	Connection                      Connection         `json:"connection"`
	SuccessCriteria                 []SuccessCriteria  `json:"success_criteria,omitempty"`
}

func parseGetHttpCheckResponse(response string) (*GetHttpCheck, error) {
	// Parse the response and return the user object
	var httpcheck GetHttpCheck
	err := json.Unmarshal([]byte(response), &httpcheck)
	if err != nil {
		return nil, err
	}

	return &httpcheck, err
}

func (c Client) GetHttpCheck(id int) (*GetHttpCheck, *RequestDetails, error) {
	details, err := c.makePublicAPICall("GET", fmt.Sprintf("/v2/checks/http/%d", id), bytes.NewBufferString("{}"), nil)

	// Check for errors
	if err != nil {
		return nil, details, err
	}

	httpcheck, err := parseGetHttpCheckResponse(details.ResponseBody)
	if err != nil {
		return httpcheck, details, err
	}

	return httpcheck, details, nil
}
