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

// Common and shared struct models used for more complex requests
type Links struct {
	Self     string `json:"self,omitempty"`
	SelfHTML string `json:"self_html,omitempty"`
	Metrics  string `json:"metrics,omitempty"`
	LastRun  string `json:"last_run,omitempty"`
}

type Tags []struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type Status struct {
	LastCode           int    `json:"last_code,omitempty"`
	LastMessage        string `json:"last_message,omitempty"`
	LastResponseTime   int    `json:"last_response_time,omitempty"`
	LastRunAt          string `json:"last_run_at,omitempty"`
	LastFailureAt      string `json:"last_failure_at,omitempty"`
	LastAlertAt        string `json:"last_alert_at,omitempty"`
	HasFailure         bool   `json:"has_failure,omitempty"`
	HasLocationFailure bool   `json:"has_location_failure,omitempty"`
}

type NotifyWho struct {
	Sms             bool   `json:"sms,omitempty"`
	Call            bool   `json:"call,omitempty"`
	Email           bool   `json:"email,omitempty"`
	CustomUserEmail string `json:"custom_email"`
	Type            string `json:"type,omitempty"`
	Links           Links  `json:"links,omitempty"`
	ID              int    `json:"id,omitempty"`
}

type NotificationWindows []struct {
	StartTime         string `json:"start_time,omitempty"`
	EndTime           string `json:"end_time,omitempty"`
	DurationInMinutes int    `json:"duration_in_minutes,omitempty"`
	TimeZone          string `json:"time_zone,omitempty"`
}

type NotificationWindow struct {
	StartTime         string `json:"start_time,omitempty"`
	EndTime           string `json:"end_time,omitempty"`
	DurationInMinutes int    `json:"duration_in_minutes,omitempty"`
	TimeZone          string `json:"time_zone,omitempty"`
}

type Escalations struct {
	Sms                bool               `json:"sms,omitempty"`
	Email              bool               `json:"email,omitempty"`
	Call               bool               `json:"call,omitempty"`
	AfterMinutes       int                `json:"after_minutes,omitempty"`
	NotifyWho          []NotifyWho        `json:"notify_who,omitempty"`
	IsRepeat           bool               `json:"is_repeat,omitempty"`
	NotificationWindow NotificationWindow `json:"notification_window,omitempty"`
}

type Notifications struct {
	Sms                     bool                `json:"sms,omitempty"`
	Email                   bool                `json:"email,omitempty"`
	Call                    bool                `json:"call,omitempty"`
	NotifyWho               []NotifyWho         `json:"notify_who,omitempty"`
	NotifyAfterFailureCount int                 `json:"notify_after_failure_count,omitempty"`
	NotifyOnLocationFailure bool                `json:"notify_on_location_failure,omitempty"`
	NotificationWindows     NotificationWindows `json:"notification_windows,omitempty"`
	Escalations             []Escalations       `json:"escalations,omitempty"`
	Muted                   bool                `json:"muted,omitempty"`
}

type SuccessCriteria struct {
	ActionType       string `json:"action_type,omitempty"`
	ComparisonString string `json:"comparison_string,omitempty"`
	CreatedAt        string `json:"created_at,omitempty"`
	UpdatedAt        string `json:"updated_at,omitempty"`
}

type Connection struct {
	UploadBandwidth   int     `json:"upload_bandwidth,omitempty"`
	DownloadBandwidth int     `json:"download_bandwidth,omitempty"`
	Latency           int     `json:"latency,omitempty"`
	PacketLoss        float64 `json:"packet_loss,omitempty"`
}

type Locations []struct {
	ID          int    `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	WorldRegion string `json:"world_region,omitempty"`
	RegionCode  string `json:"region_code,omitempty"`
}

type Integrations []struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type HTTPRequestHeaders struct {
	UserAgent string `json:"User-Agent,omitempty"`
}

type Browser struct {
	Label string `json:"label,omitempty"`
	Code  string `json:"code,omitempty"`
}

type Steps struct {
	ItemMethod   string `json:"item_method,omitempty"`
	Value        string `json:"value,omitempty"`
	How          string `json:"how,omitempty"`
	What         string `json:"what,omitempty"`
	UpdatedAt    string `json:"updated_at,omitempty"`
	CreatedAt    string `json:"created_at,omitempty"`
	VariableName string `json:"variable_name,omitempty"`
	Name         string `json:"name,omitempty"`
	Position     int    `json:"position,omitempty"`
}

type Cookies struct {
	Key    string `json:"key,omitempty"`
	Value  string `json:"value,omitempty"`
	Domain string `json:"domain,omitempty"`
	Path   string `json:"path,omitempty"`
}

type JavascriptFiles struct {
	ID        int    `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
	Links     Links  `json:"links,omitempty"`
}

type ExcludedFiles struct {
	ExclusionType string `json:"exclusion_type,omitempty"`
	PresetName    string `json:"preset_name,omitempty"`
	URL           string `json:"url,omitempty"`
	CreatedAt     string `json:"created_at,omitempty"`
	UpdatedAt     string `json:"updated_at,omitempty"`
}

type BlackoutPeriods []struct {
	StartDate         string `json:"start_date,omitempty"`
	EndDate           string `json:"end_date,omitempty"`
	Timezone          string `json:"timezone,omitempty"`
	StartTime         string `json:"start_time,omitempty"`
	EndTime           string `json:"end_time,omitempty"`
	RepeatType        string `json:"repeat_type,omitempty"`
	DurationInMinutes int    `json:"duration_in_minutes,omitempty"`
	IsRepeat          bool   `json:"is_repeat,omitempty"`
	MonthlyRepeatType string `json:"monthly_repeat_type,omitempty"`
	CreatedAt         string `json:"created_at,omitempty"`
	UpdatedAt         string `json:"updated_at,omitempty"`
}

type Viewport struct {
	Height int `json:"height,omitempty"`
	Width  int `json:"width,omitempty"`
}

type ThresholdMonitors struct {
	Matcher        string `json:"matcher,omitempty"`
	MetricName     string `json:"metric_name,omitempty"`
	ComparisonType string `json:"comparison_type,omitempty"`
	Value          int    `json:"value,omitempty"`
	CreatedAt      string `json:"created_at,omitempty"`
	UpdatedAt      string `json:"updated_at,omitempty"`
}

type DNSOverrides struct {
	OriginalDomainCom string `json:"original.domain.com,omitempty"`
	OriginalHostCom   string `json:"original.host.com,omitempty"`
}

type Checks []struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Frequency int    `json:"frequency"`
	Paused    bool   `json:"paused"`
	Muted     bool   `json:"muted"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Links     Links  `json:"links"`
	Status    Status `json:"status"`
	Tags      Tags   `json:"tags"`
}

type Errors []struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
}

type DeleteCheck struct {
	Result  string `json:"result"`
	Message string `json:"message"`
	Errors  Errors `json:"errors"`
}

type HttpCheckInput struct {
	ID                 int                `json:"id,omitempty"`
	Name               string             `json:"name,omitempty"`
	Type               string             `json:"type,omitempty"`
	Frequency          int                `json:"frequency,omitempty"`
	Paused             bool               `json:"paused,omitempty"`
	Muted              bool               `json:"muted,omitempty"`
	CreatedAt          string             `json:"created_at,omitempty"`
	UpdatedAt          string             `json:"updated_at,omitempty"`
	Links              Links              `json:"links,omitempty"`
	Tags               []string           `json:"tags"`
	Status             Status             `json:"status,omitempty"`
	RoundRobin         bool               `json:"round_robin,omitempty"`
	AutoRetry          bool               `json:"auto_retry,omitempty"`
	Enabled            bool               `json:"enabled,omitempty"`
	BlackoutPeriods    BlackoutPeriods    `json:"blackout_periods,omitempty"`
	Locations          []int              `json:"locations,omitempty"`
	Integrations       []int              `json:"integrations,omitempty"`
	HTTPRequestHeaders HTTPRequestHeaders `json:"http_request_headers,omitempty"`
	HTTPRequestBody    string             `json:"http_request_body,omitempty"`
	Notifications      Notifications      `json:"notifications,omitempty"`
	URL                string             `json:"url,omitempty"`
	HTTPMethod         string             `json:"http_method,omitempty"`
	SuccessCriteria    []SuccessCriteria  `json:"success_criteria,omitempty"`
	Connection         Connection         `json:"connection,omitempty"`
}

type HttpCheckResponse struct {
	ID                 int                `json:"id,omitempty"`
	Name               string             `json:"name,omitempty"`
	Type               string             `json:"type,omitempty"`
	Frequency          int                `json:"frequency,omitempty"`
	Paused             bool               `json:"paused,omitempty"`
	Muted              bool               `json:"muted,omitempty"`
	CreatedAt          string             `json:"created_at,omitempty"`
	UpdatedAt          string             `json:"updated_at,omitempty"`
	Links              Links              `json:"links,omitempty"`
	Tags               Tags               `json:"tags,omitempty"`
	Status             Status             `json:"status,omitempty"`
	RoundRobin         bool               `json:"round_robin,omitempty"`
	AutoRetry          bool               `json:"auto_retry,omitempty"`
	Enabled            bool               `json:"enabled,omitempty"`
	BlackoutPeriods    BlackoutPeriods    `json:"blackout_periods,omitempty"`
	Locations          Locations          `json:"locations,omitempty"`
	Integrations       Integrations       `json:"integrations,omitempty"`
	HTTPRequestHeaders HTTPRequestHeaders `json:"http_request_headers,omitempty"`
	HTTPRequestBody    string             `json:"http_request_body,omitempty"`
	Notifications      Notifications      `json:"notifications,omitempty"`
	URL                string             `json:"url,omitempty"`
	HTTPMethod         string             `json:"http_method,omitempty"`
	SuccessCriteria    []SuccessCriteria  `json:"success_criteria,omitempty"`
	Connection         Connection         `json:"connection,omitempty"`
}

type BrowserCheckInput struct {
	ID                   int                 `json:"id,omitempty"`
	Name                 string              `json:"name,omitempty"`
	Type                 string              `json:"type,omitempty"`
	Frequency            int                 `json:"frequency,omitempty"`
	Paused               bool                `json:"paused,omitempty"`
	Muted                bool                `json:"muted,omitempty"`
	CreatedAt            string              `json:"created_at,omitempty"`
	UpdatedAt            string              `json:"updated_at,omitempty"`
	Links                Links               `json:"links,omitempty"`
	Tags                 []string            `json:"tags"`
	Status               Status              `json:"status,omitempty"`
	RoundRobin           bool                `json:"round_robin,omitempty"`
	AutoRetry            bool                `json:"auto_retry,omitempty"`
	Enabled              bool                `json:"enabled,omitempty"`
	BlackoutPeriods      BlackoutPeriods     `json:"blackout_periods,omitempty"`
	Locations            []int               `json:"locations,omitempty"`
	Integrations         []int               `json:"integrations,omitempty"`
	HTTPRequestHeaders   HTTPRequestHeaders  `json:"http_request_headers,omitempty"`
	HTTPRequestBody      string              `json:"http_request_body,omitempty"`
	HTTPMethod           string              `json:"http_method,omitempty"`
	Notifications        Notifications       `json:"notifications,omitempty"`
	URL                  string              `json:"url,omitempty"`
	UserAgent            string              `json:"user_agent,omitempty"`
	AutoUpdateUserAgent  bool                `json:"auto_update_user_agent,omitempty"`
	Browser              Browser             `json:"browser,omitempty"`
	Steps                []Steps             `json:"steps,omitempty"`
	Cookies              []Cookies           `json:"cookies,omitempty"`
	JavascriptFiles      []JavascriptFiles   `json:"javascript_files,omitempty"`
	ExcludedFiles        []ExcludedFiles     `json:"excluded_files,omitempty"`
	Viewport             Viewport            `json:"viewport,omitempty"`
	EnforceSslValidation bool                `json:"enforce_ssl_validation,omitempty"`
	ThresholdMonitors    []ThresholdMonitors `json:"threshold_monitors,omitempty"`
	DNSOverrides         DNSOverrides        `json:"dns_overrides,omitempty"`
	Connection           Connection          `json:"connection,omitempty"`
	WaitForFullMetrics   bool                `json:"wait_for_full_metrics,omitempty"`
}

type BrowserCheckResponse struct {
	ID                   int                 `json:"id,omitempty"`
	Name                 string              `json:"name,omitempty"`
	Type                 string              `json:"type,omitempty"`
	Frequency            int                 `json:"frequency,omitempty"`
	Paused               bool                `json:"paused,omitempty"`
	Muted                bool                `json:"muted,omitempty"`
	CreatedAt            string              `json:"created_at,omitempty"`
	UpdatedAt            string              `json:"updated_at,omitempty"`
	Links                Links               `json:"links,omitempty"`
	Tags                 Tags                `json:"tags,omitempty"`
	Status               Status              `json:"status,omitempty"`
	RoundRobin           bool                `json:"round_robin,omitempty"`
	AutoRetry            bool                `json:"auto_retry,omitempty"`
	Enabled              bool                `json:"enabled,omitempty"`
	BlackoutPeriods      BlackoutPeriods     `json:"blackout_periods,omitempty"`
	Locations            Locations           `json:"locations,omitempty"`
	Integrations         Integrations        `json:"integrations,omitempty"`
	HTTPRequestHeaders   HTTPRequestHeaders  `json:"http_request_headers,omitempty"`
	HTTPRequestBody      string              `json:"http_request_body,omitempty"`
	HTTPMethod           string              `json:"http_method,omitempty"`
	Notifications        Notifications       `json:"notifications,omitempty"`
	URL                  string              `json:"url,omitempty"`
	UserAgent            string              `json:"user_agent,omitempty"`
	AutoUpdateUserAgent  bool                `json:"auto_update_user_agent,omitempty"`
	Browser              Browser             `json:"browser,omitempty"`
	Steps                []Steps             `json:"steps,omitempty"`
	Cookies              []Cookies           `json:"cookies,omitempty"`
	JavascriptFiles      []JavascriptFiles   `json:"javascript_files,omitempty"`
	ExcludedFiles        []ExcludedFiles     `json:"excluded_files,omitempty"`
	Viewport             Viewport            `json:"viewport,omitempty"`
	EnforceSslValidation bool                `json:"enforce_ssl_validation,omitempty"`
	ThresholdMonitors    []ThresholdMonitors `json:"threshold_monitors,omitempty"`
	DNSOverrides         DNSOverrides        `json:"dns_overrides,omitempty"`
	Connection           Connection          `json:"connection,omitempty"`
	WaitForFullMetrics   bool                `json:"wait_for_full_metrics,omitempty"`
}
