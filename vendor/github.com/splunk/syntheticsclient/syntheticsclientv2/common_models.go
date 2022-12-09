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
	"time"
)

// Common and shared struct models used for more complex requests
type Networkconnection struct {
	Description       string `json:"description,omitempty"`
	Downloadbandwidth int    `json:"download_bandwidth,omitempty"`
	Latency           int    `json:"latency,omitempty"`
	Packetloss        int    `json:"packet_loss,omitempty"`
	Uploadbandwidth   int    `json:"upload_bandwidth,omitempty"`
}

type Advancedsettings struct {
	Authentication     `json:"authentication"`
	Cookiesv2          []Cookiesv2      `json:"cookies"`
	BrowserHeaders     []BrowserHeaders `json:"headers,omitempty"`
	HostOverrides      []HostOverrides  `json:"host_overrides,omitempty"`
	UserAgent          string           `json:"user_agent,omitempty"`
	Verifycertificates bool             `json:"verifyCertificates,omitempty"`
}

type Authentication struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

type Cookiesv2 struct {
	Key    string `json:"key"`
	Value  string `json:"value"`
	Domain string `json:"domain"`
	Path   string `json:"path"`
}

type BrowserHeaders struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Domain string `json:"domain"`
}

type HostOverrides struct {
	Source         string `json:"source"`
	Target         string `json:"target"`
	KeepHostHeader bool   `json:"keep_host_header"`
}

type Transactions struct {
	Name    string    `json:"name"`
	StepsV2 []StepsV2 `json:"steps"`
}

type BusinessTransactions struct {
	Name    string    `json:"name"`
	StepsV2 []StepsV2 `json:"steps"`
}

type StepsV2 struct {
	Name         string `json:"name"`
	Type         string `json:"type"`
	URL          string `json:"url,omitempty"`
	Action       string `json:"action,omitempty"`
	WaitForNav   bool   `json:"wait_for_nav"`
	SelectorType string `json:"selector_type,omitempty"`
	Selector     string `json:"selector,omitempty"`
	Options      `json:"options,omitempty"`
}

type BusinessTransactionStepsV2 struct {
	Name         string `json:"name"`
	Type         string `json:"type"`
	URL          string `json:"url,omitempty"`
	Action       string `json:"action,omitempty"`
	WaitForNav   bool   `json:"wait_for_nav"`
	SelectorType string `json:"selector_type,omitempty"`
	Selector     string `json:"selector,omitempty"`
	Options      `json:"options,omitempty"`
}

type Options struct {
	URL string `json:"url,omitempty"`
}

type Device struct {
	ID                int    `json:"id,omitempty"`
	Label             string `json:"label,omitempty"`
	UserAgent         string `json:"user_agent,omitempty"`
	Networkconnection `json:"network_connection,omitempty"`
	Viewportheight    int `json:"viewport_height"`
	Viewportwidth     int `json:"viewport_width"`
}

type Requests struct {
	Configuration `json:"configuration,omitempty"`
	Setup         []Setup       `json:"setup,omitempty"`
	Validations   []Validations `json:"validations,omitempty"`
}

type Configuration struct {
	Body          string `json:"body"`
	Headers       `json:"headers,omitempty"`
	Name          string `json:"name,omitempty"`
	Requestmethod string `json:"requestMethod,omitempty"`
	URL           string `json:"url,omitempty"`
}

type Headers map[string]interface{}

type Setup struct {
	Extractor string `json:"extractor,omitempty"`
	Name      string `json:"name,omitempty"`
	Source    string `json:"source,omitempty"`
	Type      string `json:"type,omitempty"`
	Variable  string `json:"variable,omitempty"`
}

type Validations struct {
	Actual     string `json:"actual,omitempty"`
	Comparator string `json:"comparator,omitempty"`
	Expected   string `json:"expected,omitempty"`
	Name       string `json:"name,omitempty"`
	Type       string `json:"type,omitempty"`
}

type Tests []struct {
	Active             bool      `json:"active"`
	Createdat          time.Time `json:"created_at"`
	Frequency          int       `json:"frequency"`
	ID                 int       `json:"id"`
	Locationids        []string  `json:"locationIds"`
	Name               string    `json:"name"`
	Schedulingstrategy string    `json:"scheduling_strategy"`
	Type               string    `json:"type"`
	Updatedat          time.Time `json:"updated_at"`
}

type GetChecksV2Options struct {
	TestType string `json:"testType"`
	PerPage  int    `json:"perPage"`
	Page     int    `json:"page"`
	Search   string `json:"search"`
	OrderBy  string `json:"orderBy"`
}

type Errors []struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
}

type HttpHeaders struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

type Variable struct {
	Createdat   time.Time `json:"created_at,omitempty"`
	Description string    `json:"description,omitempty"`
	ID          int       `json:"id,omitempty"`
	Name        string    `json:"name"`
	Secret      bool      `json:"secret"`
	Updatedat   time.Time `json:"updated_at,omitempty"`
	Value       string    `json:"value"`
}

type DeleteCheck struct {
	Result  string `json:"result"`
	Message string `json:"message"`
	Errors  Errors `json:"errors"`
}

type VariableV2Response struct {
	Variable `json:"variable"`
}

type VariableV2Input struct {
	Variable `json:"variable"`
}

type ChecksV2Response struct {
	Nextpagelink int `json:"nextPageLink"`
	Perpage      int `json:"perPage"`
	Tests        `json:"tests"`
	Totalcount   int `json:"totalCount"`
}

type PortCheckV2Response struct {
	Test struct {
		ID                 int       `json:"id"`
		Name               string    `json:"name"`
		Active             bool      `json:"active"`
		Frequency          int       `json:"frequency"`
		SchedulingStrategy string    `json:"scheduling_strategy"`
		CreatedAt          time.Time `json:"created_at"`
		UpdatedAt          time.Time `json:"updated_at"`
		LocationIds        []string  `json:"location_ids"`
		Type               string    `json:"type"`
		Protocol           string    `json:"protocol"`
		Host               string    `json:"host"`
		Port               int       `json:"port"`
	} `json:"test"`
}

type PortCheckV2Input struct {
	Test struct {
		Name               string   `json:"name"`
		Type               string   `json:"type"`
		URL                string   `json:"url"`
		Port               int      `json:"port"`
		Protocol           string   `json:"protocol"`
		Host               string   `json:"host"`
		LocationIds        []string `json:"location_ids"`
		Frequency          int      `json:"frequency"`
		SchedulingStrategy string   `json:"scheduling_strategy"`
		Active             bool     `json:"active"`
	} `json:"test"`
}

type HttpCheckV2Response struct {
	Test struct {
		ID                 int           `json:"id"`
		Name               string        `json:"name"`
		Active             bool          `json:"active"`
		Frequency          int           `json:"frequency"`
		SchedulingStrategy string        `json:"scheduling_strategy"`
		CreatedAt          time.Time     `json:"created_at,omitempty"`
		UpdatedAt          time.Time     `json:"updated_at,omitempty"`
		LocationIds        []string      `json:"location_ids"`
		Type               string        `json:"type"`
		URL                string        `json:"url"`
		RequestMethod      string        `json:"request_method"`
		Body               string        `json:"body,omitempty"`
		HttpHeaders        []HttpHeaders `json:"headers,omitempty"`
	} `json:"test"`
}

type HttpCheckV2Input struct {
	Test struct {
		Name               string        `json:"name"`
		Type               string        `json:"type"`
		URL                string        `json:"url"`
		LocationIds        []string      `json:"location_ids"`
		Frequency          int           `json:"frequency"`
		SchedulingStrategy string        `json:"scheduling_strategy"`
		Active             bool          `json:"active"`
		RequestMethod      string        `json:"request_method"`
		Body               string        `json:"body,omitempty"`
		HttpHeaders        []HttpHeaders `json:"headers,omitempty"`
	} `json:"test"`
}

type ApiCheckV2Input struct {
	Test struct {
		Active             bool       `json:"active"`
		Deviceid           int        `json:"device_id"`
		Frequency          int        `json:"frequency"`
		Locationids        []string   `json:"location_ids"`
		Name               string     `json:"name"`
		Requests           []Requests `json:"requests"`
		Schedulingstrategy string     `json:"scheduling_strategy"`
	} `json:"test"`
}

type ApiCheckV2Response struct {
	Test struct {
		Active             bool      `json:"active,omitempty"`
		Createdat          time.Time `json:"created_at"`
		Device             `json:"device,omitempty"`
		Frequency          int        `json:"frequency,omitempty"`
		ID                 int        `json:"id,omitempty"`
		Locationids        []string   `json:"location_ids,omitempty"`
		Name               string     `json:"name,omitempty"`
		Requests           []Requests `json:"requests,omitempty"`
		Schedulingstrategy string     `json:"scheduling_strategy,omitempty"`
		Type               string     `json:"type,omitempty"`
		Updatedat          time.Time  `json:"updated_at,omitempty"`
	}
}

type BrowserCheckV2Input struct {
	Test struct {
		Name                 string                 `json:"name"`
		BusinessTransactions []BusinessTransactions `json:"business_transactions"`
		Urlprotocol          string                 `json:"urlProtocol"`
		Starturl             string                 `json:"startUrl"`
		LocationIds          []string               `json:"location_ids"`
		DeviceID             int                    `json:"device_id"`
		Frequency            int                    `json:"frequency"`
		Schedulingstrategy   string                 `json:"scheduling_strategy"`
		Active               bool                   `json:"active"`
		Advancedsettings     `json:"advanced_settings"`
	} `json:"test"`
}

type BrowserCheckV2Response struct {
	Test struct {
		Active               bool `json:"active"`
		Advancedsettings     `json:"advanced_settings"`
		BusinessTransactions []BusinessTransactions `json:"business_transactions"`
		Createdat            time.Time              `json:"created_at"`
		Device               `json:"device"`
		Frequency            int            `json:"frequency"`
		ID                   int            `json:"id"`
		Locationids          []string       `json:"location_ids"`
		Name                 string         `json:"name"`
		Schedulingstrategy   string         `json:"scheduling_strategy"`
		Transactions         []Transactions `json:"transactions"`
		Type                 string         `json:"type"`
		Updatedat            time.Time      `json:"updated_at"`
	} `json:"test"`
}
