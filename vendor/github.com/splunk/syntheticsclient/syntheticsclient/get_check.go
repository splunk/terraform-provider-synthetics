package syntheticsclient

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type GetCheck struct {
	ID                              int
	Name                            string              `json:"name"`
	Type                            string              `json:"type"`
	Frequency                       int                 `json:"frequency,omitempty"`
	Paused                          bool                `json:"paused,omitempty"`
	Muted                           bool                `json:"muted,omitempty"`
	CreatedAt                       string              `json:"created_at,omitempty"`
	UpdatedAt                       string              `json:"updated_at,omitempty"`
	Links                           Links               `json:"links,omitempty"`
	Status                          Status              `json:"status,omitempty"`
	Notifications                   Notifications       `json:"notifications,omitempty"`
	ResponseTimeMonitorMilliseconds int                 `json:"response_time_monitor_milliseconds,omitempty"`
	HTTPRequestHeaders              HTTPRequestHeaders  `json:"http_request_headers,omitempty"`
	HTTPRequestBody                 string              `json:"http_request_body,omitempty"`
	HTTPMethod                      string              `json:"http_method,omitempty"`
	RoundRobin                      bool                `json:"round_robin,omitempty"`
	AutoRetry                       bool                `json:"auto_retry,omitempty"`
	Enabled                         bool                `json:"enabled,omitempty"`
	Integrations                    Integrations        `json:"integrations,omitempty"`
	URL                             string              `json:"url,omitempty"`
	UserAgent                       string              `json:"user_agent,omitempty"`
	AutoUpdateUserAgent             bool                `json:"auto_update_user_agent,omitempty"`
	Viewport                        Viewport            `json:"viewport,omitempty"`
	EnforceSslValidation            bool                `json:"enforce_ssl_validation,omitempty"`
	Browser                         Browser             `json:"browser,omitempty"`
	DNSOverrides                    DNSOverrides        `json:"dns_overrides,omitempty"`
	WaitForFullMetrics              bool                `json:"wait_for_full_metrics,omitempty"`
	Tags                            Tags                `json:"tags,omitempty"`
	BlackoutPeriods                 BlackoutPeriods     `json:"blackout_periods,omitempty"`
	Locations                       Locations           `json:"locations,omitempty"`
	Steps                           []Steps             `json:"steps,omitempty"`
	JavascriptFiles                 []JavascriptFiles   `json:"javascript_files,omitempty"`
	ThresholdMonitors               []ThresholdMonitors `json:"threshold_monitors,omitempty"`
	ExcludedFiles                   []ExcludedFiles     `json:"excluded_files,omitempty"`
	Cookies                         []Cookies           `json:"cookies,omitempty"`
	Connection                      Connection          `json:"connection,omitempty"`
	SuccessCriteria                 []SuccessCriteria   `json:"success_criteria,omitempty"`
}

func parseCheckResponse(response string) (*GetCheck, error) {
	// Parse the response and return the check object
	var check GetCheck
	err := json.Unmarshal([]byte(response), &check)
	if err != nil {
		return nil, err
	}

	return &check, err
}

func (c Client) GetCheck(id int) (*GetCheck, *RequestDetails, error) {

	details, err := c.makePublicAPICall("GET",
		fmt.Sprintf("/v2/checks/%d", id),
		bytes.NewBufferString("{}"),
		nil)

	if err != nil {
		return nil, details, err
	}

	check, err := parseCheckResponse(details.ResponseBody)
	if err != nil {
		return check, details, err
	}

	return check, details, nil
}
