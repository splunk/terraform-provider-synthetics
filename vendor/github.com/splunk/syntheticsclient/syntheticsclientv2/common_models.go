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
	Downloadbandwidth int    `json:"downloadBandwidth,omitempty"`
	Latency           int    `json:"latency,omitempty"`
	Packetloss        int    `json:"packetLoss,omitempty"`
	Uploadbandwidth   int    `json:"uploadBandwidth,omitempty"`
}

type Advancedsettings struct {
	Authentication     `json:"authentication"`
	Cookiesv2          []Cookiesv2      `json:"cookies"`
	BrowserHeaders     []BrowserHeaders `json:"headers,omitempty"`
	HostOverrides      []HostOverrides  `json:"hostOverrides,omitempty"`
	UserAgent          string           `json:"userAgent,omitempty"`
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
	KeepHostHeader bool   `json:"keepHostHeader"`
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
	WaitForNav   bool   `json:"waitForNav"`
	SelectorType string `json:"selectorType,omitempty"`
	Selector     string `json:"selector,omitempty"`
	Options      `json:"options,omitempty"`
}

type Options struct {
	URL string `json:"url,omitempty"`
}

type Device struct {
	ID                int    `json:"id,omitempty"`
	Label             string `json:"label,omitempty"`
	UserAgent         string `json:"userAgent,omitempty"`
	Networkconnection `json:"networkConnection,omitempty"`
	Viewportheight    int `json:"viewportHeight"`
	Viewportwidth     int `json:"viewportWidth"`
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
	RequestMethod string `json:"requestMethod,omitempty"`
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
	Createdat          time.Time `json:"createdAt"`
	Frequency          int       `json:"frequency"`
	ID                 int       `json:"id"`
	Locationids        []string  `json:"locationIds"`
	Name               string    `json:"name"`
	Schedulingstrategy string    `json:"schedulingStrategy"`
	Type               string    `json:"type"`
	Updatedat          time.Time `json:"updatedAt"`
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
	Createdat   time.Time `json:"createdAt,omitempty"`
	Description string    `json:"description,omitempty"`
	ID          int       `json:"id,omitempty"`
	Name        string    `json:"name"`
	Secret      bool      `json:"secret"`
	Updatedat   time.Time `json:"updatedAt,omitempty"`
	Value       string    `json:"value"`
}

type DeleteCheck struct {
	Result  string `json:"result"`
	Message string `json:"message"`
	Errors  Errors `json:"errors"`
}

type Location struct {
	ID      string `json:"id"`
	Label   string `json:"label"`
	Default bool   `json:"default"`
	Type    string `json:"type,omitempty"`
	Country string `json:"country,omitempty"`
}

type Meta struct {
	ActiveTestIds []int `json:"activeTestIds"`
	PausedTestIds []int `json:"pausedTestIds"`
}

type DevicesV2Response struct {
	Devices []Device `json:"devices"`
}

type VariableV2Response struct {
	Variable `json:"variable"`
}

type VariableV2Input struct {
	Variable `json:"variable"`
}

type VariablesV2Response struct {
	Variable []Variable `json:"variables"`
}

type LocationsV2Response struct {
	Location           []Location `json:"locations"`
	DefaultLocationIds []string   `json:"defaultLocationIds"`
}

type LocationV2Response struct {
	Location `json:"location"`
	Meta     `json:"meta"`
}

type LocationV2Input struct {
	Location `json:"location"`
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
		SchedulingStrategy string    `json:"schedulingStrategy"`
		CreatedAt          time.Time `json:"createdAt"`
		UpdatedAt          time.Time `json:"updatedAt"`
		LocationIds        []string  `json:"locationIds"`
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
		LocationIds        []string `json:"locationIds"`
		Frequency          int      `json:"frequency"`
		SchedulingStrategy string   `json:"schedulingStrategy"`
		Active             bool     `json:"active"`
	} `json:"test"`
}

type HttpCheckV2Response struct {
	Test struct {
		ID                 int           `json:"id"`
		Name               string        `json:"name"`
		Active             bool          `json:"active"`
		Frequency          int           `json:"frequency"`
		SchedulingStrategy string        `json:"schedulingStrategy"`
		CreatedAt          time.Time     `json:"createdAt,omitempty"`
		UpdatedAt          time.Time     `json:"updatedAt,omitempty"`
		LocationIds        []string      `json:"locationIds"`
		Type               string        `json:"type"`
		URL                string        `json:"url"`
		RequestMethod      string        `json:"requestMethod"`
		Body               string        `json:"body,omitempty"`
		HttpHeaders        []HttpHeaders `json:"headers,omitempty"`
	} `json:"test"`
}

type HttpCheckV2Input struct {
	Test struct {
		Name               string        `json:"name"`
		Type               string        `json:"type"`
		URL                string        `json:"url"`
		LocationIds        []string      `json:"locationIds"`
		Frequency          int           `json:"frequency"`
		SchedulingStrategy string        `json:"schedulingStrategy"`
		Active             bool          `json:"active"`
		RequestMethod      string        `json:"requestMethod"`
		Body               string        `json:"body,omitempty"`
		HttpHeaders        []HttpHeaders `json:"headers,omitempty"`
	} `json:"test"`
}

type ApiCheckV2Input struct {
	Test struct {
		Active             bool       `json:"active"`
		Deviceid           int        `json:"deviceId"`
		Frequency          int        `json:"frequency"`
		Locationids        []string   `json:"locationIds"`
		Name               string     `json:"name"`
		Requests           []Requests `json:"requests"`
		Schedulingstrategy string     `json:"schedulingStrategy"`
	} `json:"test"`
}

type ApiCheckV2Response struct {
	Test struct {
		Active             bool      `json:"active,omitempty"`
		Createdat          time.Time `json:"createdAt"`
		Device             `json:"device,omitempty"`
		Frequency          int        `json:"frequency,omitempty"`
		ID                 int        `json:"id,omitempty"`
		Locationids        []string   `json:"locationIds,omitempty"`
		Name               string     `json:"name,omitempty"`
		Requests           []Requests `json:"requests,omitempty"`
		Schedulingstrategy string     `json:"schedulingStrategy,omitempty"`
		Type               string     `json:"type,omitempty"`
		Updatedat          time.Time  `json:"updatedAt,omitempty"`
	}
}

type BrowserCheckV2Input struct {
	Test struct {
		Name               string         `json:"name"`
		Transactions       []Transactions `json:"transactions"`
		Urlprotocol        string         `json:"urlProtocol"`
		Starturl           string         `json:"startUrl"`
		LocationIds        []string       `json:"locationIds"`
		DeviceID           int            `json:"deviceId"`
		Frequency          int            `json:"frequency"`
		Schedulingstrategy string         `json:"schedulingStrategy"`
		Active             bool           `json:"active"`
		Advancedsettings   `json:"advancedSettings"`
	} `json:"test"`
}

type BrowserCheckV2Response struct {
	Test struct {
		Active             bool `json:"active"`
		Advancedsettings   `json:"advancedSettings"`
		Createdat          time.Time `json:"createdAt"`
		Device             `json:"device"`
		Frequency          int            `json:"frequency"`
		ID                 int            `json:"id"`
		Locationids        []string       `json:"locationIds"`
		Name               string         `json:"name"`
		Schedulingstrategy string         `json:"schedulingStrategy"`
		Transactions       []Transactions `json:"transactions"`
		Type               string         `json:"type"`
		Updatedat          time.Time      `json:"updatedAt"`
	} `json:"test"`
}
