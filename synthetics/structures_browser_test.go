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

package synthetics

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sc2 "github.com/splunk/syntheticsclient/v3/syntheticsclientv2"
)

// Test stepSelectorRepresentationDiffers
func TestStepSelectorRepresentationDiffers(t *testing.T) {
	tests := []struct {
		name string
		a    stepSelectorInput
		b    stepSelectorInput
		want bool
	}{
		{
			name: "both use legacy fields",
			a:    stepSelectorInput{selectorType: "id", selector: "btn"},
			b:    stepSelectorInput{selectorType: "id", selector: "btn"},
			want: false,
		},
		{
			name: "both use selectors list",
			a:    stepSelectorInput{selectors: []sc2.Selector{{Type: "id", Value: "btn"}}},
			b:    stepSelectorInput{selectors: []sc2.Selector{{Type: "id", Value: "btn"}}},
			want: false,
		},
		{
			name: "state legacy, config selectors list",
			a:    stepSelectorInput{selectorType: "id", selector: "btn"},
			b:    stepSelectorInput{selectors: []sc2.Selector{{Type: "id", Value: "btn"}}},
			want: true,
		},
		{
			name: "state selectors, config legacy",
			a:    stepSelectorInput{selectors: []sc2.Selector{{Type: "id", Value: "btn"}}},
			b:    stepSelectorInput{selectorType: "id", selector: "btn"},
			want: true,
		},
		{
			name: "both empty",
			a:    stepSelectorInput{},
			b:    stepSelectorInput{},
			want: false,
		},
		{
			name: "one has empty selectors array",
			a:    stepSelectorInput{},
			b:    stepSelectorInput{selectors: []sc2.Selector{}},
			want: false,
		},
		{
			name: "state has legacy, config has neither",
			a:    stepSelectorInput{selectorType: "id", selector: "btn"},
			b:    stepSelectorInput{},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stepSelectorRepresentationDiffers(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("stepSelectorRepresentationDiffers(%+v, %+v) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

// Test flattenBrowserV2Read
func TestFlattenBrowserV2Read(t *testing.T) {
	tests := []struct {
		name     string
		resp     *sc2.BrowserCheckV2Response
		checkKey string
		checkVal interface{}
	}{
		{
			name:     "nil response",
			resp:     &sc2.BrowserCheckV2Response{},
			checkKey: "active",
			checkVal: false,
		},
		{
			name: "populated response",
			resp: &sc2.BrowserCheckV2Response{
				Test: sc2.BrowserCheckV2ResponseTest{
					Active:             true,
					Automaticretries:   2,
					Deviceid:           5,
					Frequency:          60,
					Name:               "Test Browser",
					Schedulingstrategy: "round_robin",
					Locationids:        []string{"loc1", "loc2"},
					Transactions: []sc2.Transactions{
						{Name: "txn1", StepsV2: []sc2.StepsV2{}},
					},
					Customproperties: []sc2.CustomProperties{
						{Key: "env", Value: "prod"},
					},
					Advancedsettings: sc2.Advancedsettings{},
				},
			},
			checkKey: "name",
			checkVal: "Test Browser",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := flattenBrowserV2Read(tt.resp)
			if len(got) == 0 {
				t.Fatal("expected non-empty result")
			}
			result := got[0].(map[string]interface{})
			if result[tt.checkKey] != tt.checkVal {
				t.Errorf("flattenBrowserV2Read()[%s] = %v, want %v", tt.checkKey, result[tt.checkKey], tt.checkVal)
			}
		})
	}
}

// Test flattenBrowserV2Data
func TestFlattenBrowserV2Data(t *testing.T) {
	tests := []struct {
		name   string
		resp   *sc2.BrowserCheckV2Response
		device *sc2.Device
	}{
		{
			name: "empty response",
			resp: &sc2.BrowserCheckV2Response{},
		},
		{
			name: "populated response",
			resp: &sc2.BrowserCheckV2Response{
				Test: sc2.BrowserCheckV2ResponseTest{
					Active:    true,
					ID:        123,
					Name:      "BrowserTest",
					Frequency: 30,
				},
			},
		},
		{
			name: "response with device",
			resp: &sc2.BrowserCheckV2Response{
				Test: sc2.BrowserCheckV2ResponseTest{
					Deviceid: 10,
					Name:     "TestWithDevice",
				},
			},
			device: &sc2.Device{ID: 10, Label: "iPhone"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			devices := []sc2.Device{}
			if tt.device != nil {
				devices = append(devices, *tt.device)
			}
			got := flattenBrowserV2Data(tt.resp, devices)
			if len(got) == 0 {
				t.Fatal("expected non-empty result")
			}
			if _, ok := got[0].(map[string]interface{}); !ok {
				t.Fatal("expected map result")
			}
		})
	}
}

// Test flattenSetupData
func TestFlattenSetupData(t *testing.T) {
	tests := []struct {
		name  string
		setup *[]sc2.Setup
		want  int
	}{
		{
			name:  "nil setup",
			setup: nil,
			want:  0,
		},
		{
			name:  "empty setup",
			setup: &[]sc2.Setup{},
			want:  0,
		},
		{
			name: "single setup",
			setup: &[]sc2.Setup{
				{
					Name:      "setup1",
					Type:      "variable",
					Extractor: "json",
					Source:    "response",
					Variable:  "var1",
					Code:      "code1",
					Value:     "value1",
				},
			},
			want: 1,
		},
		{
			name: "multiple setups",
			setup: &[]sc2.Setup{
				{Name: "setup1", Type: "variable"},
				{Name: "setup2", Type: "variable"},
				{Name: "setup3", Type: "javascript"},
			},
			want: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := flattenSetupData(tt.setup)
			if len(got) != tt.want {
				t.Errorf("flattenSetupData() = %d items, want %d", len(got), tt.want)
			}
			if tt.want > 0 {
				for i, item := range got {
					if _, ok := item.(map[string]interface{}); !ok {
						t.Errorf("item %d is not a map", i)
					}
				}
			}
		})
	}
}

// Test flattenTransactionsData
func TestFlattenTransactionsData(t *testing.T) {
	tests := []struct {
		name string
		txns *[]sc2.Transactions
		want int
	}{
		{
			name: "nil transactions",
			txns: nil,
			want: 0,
		},
		{
			name: "empty transactions",
			txns: &[]sc2.Transactions{},
			want: 0,
		},
		{
			name: "single transaction",
			txns: &[]sc2.Transactions{
				{
					Name: "txn1",
					StepsV2: []sc2.StepsV2{
						{Name: "step1", Type: "navigate"},
					},
				},
			},
			want: 1,
		},
		{
			name: "multiple transactions",
			txns: &[]sc2.Transactions{
				{Name: "txn1"},
				{Name: "txn2"},
				{Name: "txn3"},
			},
			want: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := flattenTransactionsData(tt.txns)
			if len(got) != tt.want {
				t.Errorf("flattenTransactionsData() = %d items, want %d", len(got), tt.want)
			}
			for _, item := range got {
				if m, ok := item.(map[string]interface{}); ok {
					if _, hasName := m["name"]; !hasName && tt.want > 0 {
						t.Error("expected name field in transaction")
					}
				}
			}
		})
	}
}

// Test flattenBusinessTransactionsData
func TestFlattenBusinessTransactionsData(t *testing.T) {
	tests := []struct {
		name string
		txns *[]sc2.Transactions
		want int
	}{
		{
			name: "nil",
			txns: nil,
			want: 0,
		},
		{
			name: "empty",
			txns: &[]sc2.Transactions{},
			want: 0,
		},
		{
			name: "single business transaction",
			txns: &[]sc2.Transactions{
				{
					Name: "BizTxn1",
					StepsV2: []sc2.StepsV2{
						{Name: "step1", Type: "click_element"},
					},
				},
			},
			want: 1,
		},
		{
			name: "multiple business transactions",
			txns: &[]sc2.Transactions{
				{Name: "BizTxn1"},
				{Name: "BizTxn2"},
			},
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := flattenBusinessTransactionsData(tt.txns)
			if len(got) != tt.want {
				t.Errorf("flattenBusinessTransactionsData() = %d items, want %d", len(got), tt.want)
			}
		})
	}
}

// Test buildSetupData
func TestBuildSetupData(t *testing.T) {
	tests := []struct {
		name  string
		setup []interface{}
		want  int
	}{
		{
			name:  "empty input",
			setup: []interface{}{},
			want:  0,
		},
		{
			name: "single setup",
			setup: []interface{}{
				map[string]interface{}{
					"name":      "setup1",
					"type":      "variable",
					"extractor": "json",
					"source":    "response",
					"variable":  "var1",
					"code":      "",
					"value":     "",
				},
			},
			want: 1,
		},
		{
			name: "multiple setups",
			setup: []interface{}{
				map[string]interface{}{
					"name": "setup1", "type": "variable", "extractor": "", "source": "", "variable": "", "code": "", "value": "",
				},
				map[string]interface{}{
					"name": "setup2", "type": "javascript", "extractor": "", "source": "", "variable": "", "code": "code", "value": "",
				},
			},
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildSetupData(tt.setup)
			if len(got) != tt.want {
				t.Errorf("buildSetupData() = %d items, want %d", len(got), tt.want)
			}
		})
	}
}

// Test selectorsFromFields
func TestSelectorsFromFields(t *testing.T) {
	tests := []struct {
		name  string
		stype string
		sval  string
		want  int
	}{
		{
			name:  "both empty",
			stype: "",
			sval:  "",
			want:  0,
		},
		{
			name:  "type only",
			stype: "id",
			sval:  "",
			want:  0,
		},
		{
			name:  "value only",
			stype: "",
			sval:  "selector",
			want:  0,
		},
		{
			name:  "both populated",
			stype: "id",
			sval:  "submit",
			want:  1,
		},
		{
			name:  "css selector",
			stype: "css",
			sval:  ".button",
			want:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := selectorsFromFields(tt.stype, tt.sval)
			if len(got) != tt.want {
				t.Errorf("selectorsFromFields() = %d selectors, want %d", len(got), tt.want)
			}
			if tt.want > 0 && len(got) > 0 {
				if got[0].Type != tt.stype || got[0].Value != tt.sval {
					t.Errorf("selector mismatch: got %+v, want type=%q value=%q", got[0], tt.stype, tt.sval)
				}
			}
		})
	}
}

// Test parseSelectorsList
func TestParseSelectorsList(t *testing.T) {
	tests := []struct {
		name string
		raw  interface{}
		want int
	}{
		{
			name: "nil input",
			raw:  nil,
			want: 0,
		},
		{
			name: "not a list",
			raw:  "string",
			want: 0,
		},
		{
			name: "empty list",
			raw:  []interface{}{},
			want: 0,
		},
		{
			name: "single selector",
			raw: []interface{}{
				map[string]interface{}{"type": "id", "value": "btn"},
			},
			want: 1,
		},
		{
			name: "multiple selectors",
			raw: []interface{}{
				map[string]interface{}{"type": "css", "value": ".btn"},
				map[string]interface{}{"type": "id", "value": "submit"},
			},
			want: 2,
		},
		{
			name: "selector missing type",
			raw: []interface{}{
				map[string]interface{}{"type": "", "value": "btn"},
			},
			want: 0,
		},
		{
			name: "selector missing value",
			raw: []interface{}{
				map[string]interface{}{"type": "id", "value": ""},
			},
			want: 0,
		},
		{
			name: "invalid item type",
			raw: []interface{}{
				"not a map",
			},
			want: 0,
		},
		{
			name: "mixed valid and invalid",
			raw: []interface{}{
				map[string]interface{}{"type": "id", "value": "btn"},
				"invalid",
				map[string]interface{}{"type": "css", "value": ".primary"},
			},
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseSelectorsList(tt.raw)
			if len(got) != tt.want {
				t.Errorf("parseSelectorsList() = %d selectors, want %d", len(got), tt.want)
			}
		})
	}
}

// Test resolveSingleSelector
func TestResolveSingleSelector(t *testing.T) {
	tests := []struct {
		name   string
		in     stepSelectorInput
		wantOK bool
		wantT  string
		wantV  string
	}{
		{
			name:   "empty input",
			in:     stepSelectorInput{},
			wantOK: false,
		},
		{
			name: "legacy fields",
			in: stepSelectorInput{
				selectorType: "id",
				selector:     "btn",
			},
			wantOK: true,
			wantT:  "id",
			wantV:  "btn",
		},
		{
			name: "single selector in list",
			in: stepSelectorInput{
				selectors: []sc2.Selector{
					{Type: "css", Value: ".primary"},
				},
			},
			wantOK: true,
			wantT:  "css",
			wantV:  ".primary",
		},
		{
			name: "multiple selectors",
			in: stepSelectorInput{
				selectors: []sc2.Selector{
					{Type: "id", Value: "btn1"},
					{Type: "css", Value: ".btn2"},
				},
			},
			wantOK: false,
		},
		{
			name: "legacy type but no value",
			in: stepSelectorInput{
				selectorType: "id",
			},
			wantOK: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotT, gotV, gotOK := tt.in.resolveSingleSelector()
			if gotOK != tt.wantOK {
				t.Errorf("resolveSingleSelector() ok = %v, want %v", gotOK, tt.wantOK)
			}
			if gotOK {
				if gotT != tt.wantT {
					t.Errorf("resolveSingleSelector() type = %q, want %q", gotT, tt.wantT)
				}
				if gotV != tt.wantV {
					t.Errorf("resolveSingleSelector() value = %q, want %q", gotV, tt.wantV)
				}
			}
		})
	}
}

// Test flattenSelectorsData
func TestFlattenSelectorsData(t *testing.T) {
	tests := []struct {
		name      string
		selectors []sc2.Selector
		want      int
	}{
		{
			name:      "empty selectors",
			selectors: []sc2.Selector{},
			want:      0,
		},
		{
			name: "single selector",
			selectors: []sc2.Selector{
				{Type: "id", Value: "btn"},
			},
			want: 1,
		},
		{
			name: "multiple selectors",
			selectors: []sc2.Selector{
				{Type: "css", Value: ".primary"},
				{Type: "id", Value: "submit"},
				{Type: "xpath", Value: "//button"},
			},
			want: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := flattenSelectorsData(tt.selectors)
			if len(got) != tt.want {
				t.Errorf("flattenSelectorsData() = %d items, want %d", len(got), tt.want)
			}
			for i, item := range got {
				m, ok := item.(map[string]interface{})
				if !ok {
					t.Errorf("item %d is not a map", i)
					continue
				}
				if _, hasType := m["type"]; !hasType {
					t.Errorf("item %d missing type field", i)
				}
				if _, hasValue := m["value"]; !hasValue {
					t.Errorf("item %d missing value field", i)
				}
			}
		})
	}
}

// Test buildSelectorsFromStep
func TestBuildSelectorsFromStep(t *testing.T) {
	tests := []struct {
		name    string
		step    map[string]interface{}
		want    int
		wantErr bool
	}{
		{
			name:    "empty step",
			step:    map[string]interface{}{},
			want:    0,
			wantErr: false,
		},
		{
			name: "legacy fields only",
			step: map[string]interface{}{
				"selector_type": "id",
				"selector":      "btn",
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "selectors list",
			step: map[string]interface{}{
				"selectors": []interface{}{
					map[string]interface{}{"type": "css", "value": ".btn"},
				},
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "selectors preferred over legacy",
			step: map[string]interface{}{
				"selector_type": "id",
				"selector":      "ignored",
				"selectors": []interface{}{
					map[string]interface{}{"type": "css", "value": ".btn"},
				},
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "multiple selectors in list",
			step: map[string]interface{}{
				"selectors": []interface{}{
					map[string]interface{}{"type": "css", "value": ".btn1"},
					map[string]interface{}{"type": "id", "value": "btn2"},
				},
			},
			want:    2,
			wantErr: false,
		},
		{
			name: "selector missing type",
			step: map[string]interface{}{
				"selectors": []interface{}{
					map[string]interface{}{"type": "", "value": "btn"},
				},
				"name": "Test Step",
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "invalid selector item",
			step: map[string]interface{}{
				"selectors": []interface{}{
					"not a map",
				},
				"name": "Test Step",
			},
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := buildSelectorsFromStep(tt.step)
			if (err != nil) != tt.wantErr {
				t.Errorf("buildSelectorsFromStep() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.want {
				t.Errorf("buildSelectorsFromStep() = %d selectors, want %d", len(got), tt.want)
			}
		})
	}
}

// Test flattenStepsData
func TestFlattenStepsData(t *testing.T) {
	tests := []struct {
		name  string
		steps *[]sc2.StepsV2
		want  int
	}{
		{
			name:  "nil steps",
			steps: nil,
			want:  0,
		},
		{
			name:  "empty steps",
			steps: &[]sc2.StepsV2{},
			want:  0,
		},
		{
			name: "single step",
			steps: &[]sc2.StepsV2{
				{
					Name: "step1",
					Type: "navigate",
					URL:  "https://example.com",
				},
			},
			want: 1,
		},
		{
			name: "multiple steps with selectors",
			steps: &[]sc2.StepsV2{
				{
					Name:      "Click",
					Type:      "click_element",
					Selectors: []sc2.Selector{{Type: "id", Value: "btn"}},
				},
				{
					Name: "Navigate",
					Type: "navigate",
					URL:  "https://example.com",
				},
			},
			want: 2,
		},
		{
			name: "step with all fields",
			steps: &[]sc2.StepsV2{
				{
					Name:                     "ComplexStep",
					Type:                     "click_element",
					URL:                      "https://test.com",
					WaitForNav:               true,
					WaitForNavTimeout:        5000,
					WaitForNavTimeoutDefault: false,
					MaxWaitTime:              10000,
					MaxWaitTimeDefault:       false,
					Selectors:                []sc2.Selector{{Type: "css", Value: ".btn"}},
					OptionSelectorType:       "css",
					OptionSelector:           ".option",
					VariableName:             "var1",
					Value:                    "testval",
					Duration:                 1000,
				},
			},
			want: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := flattenStepsData(tt.steps)
			if len(got) != tt.want {
				t.Errorf("flattenStepsData() = %d items, want %d", len(got), tt.want)
			}
		})
	}
}

// Test flattenCookiesData
func TestFlattenCookiesData(t *testing.T) {
	tests := []struct {
		name    string
		cookies *[]sc2.Cookiesv2
		want    int
	}{
		{
			name:    "nil cookies",
			cookies: nil,
			want:    0,
		},
		{
			name:    "empty cookies",
			cookies: &[]sc2.Cookiesv2{},
			want:    0,
		},
		{
			name: "single cookie",
			cookies: &[]sc2.Cookiesv2{
				{Key: "session", Value: "abc123", Domain: "example.com", Path: "/"},
			},
			want: 1,
		},
		{
			name: "multiple cookies",
			cookies: &[]sc2.Cookiesv2{
				{Key: "session", Value: "abc123"},
				{Key: "token", Value: "xyz789"},
				{Key: "prefs", Value: "dark_mode"},
			},
			want: 3,
		},
		{
			name: "cookie with partial fields",
			cookies: &[]sc2.Cookiesv2{
				{Key: "id", Value: ""},
				{Key: "", Value: "val"},
			},
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := flattenCookiesData(tt.cookies)
			if len(got) != tt.want {
				t.Errorf("flattenCookiesData() = %d items, want %d", len(got), tt.want)
			}
		})
	}
}

// Test buildCookiesData
func TestBuildCookiesData(t *testing.T) {
	tests := []struct {
		name    string
		cookies *schema.Set
		want    int
	}{
		{
			name: "single cookie",
			cookies: testCookiesSet(
				map[string]interface{}{
					"key":    "session",
					"value":  "abc123",
					"domain": "example.com",
					"path":   "/",
				},
			),
			want: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildCookiesData(tt.cookies)
			if len(got) != tt.want {
				t.Errorf("buildCookiesData() = %d items, want %d", len(got), tt.want)
			}
		})
	}
}

// Test flattenBrowserHeadersData
func TestFlattenBrowserHeadersData(t *testing.T) {
	tests := []struct {
		name    string
		headers *[]sc2.BrowserHeaders
		want    int
	}{
		{
			name:    "nil headers",
			headers: nil,
			want:    0,
		},
		{
			name:    "empty headers",
			headers: &[]sc2.BrowserHeaders{},
			want:    0,
		},
		{
			name: "single header",
			headers: &[]sc2.BrowserHeaders{
				{Name: "X-Custom", Value: "value1", Domain: "example.com"},
			},
			want: 1,
		},
		{
			name: "multiple headers",
			headers: &[]sc2.BrowserHeaders{
				{Name: "X-Custom-1", Value: "val1"},
				{Name: "X-Custom-2", Value: "val2"},
				{Name: "X-Auth", Value: "Bearer token"},
			},
			want: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := flattenBrowserHeadersData(tt.headers)
			if len(got) != tt.want {
				t.Errorf("flattenBrowserHeadersData() = %d items, want %d", len(got), tt.want)
			}
		})
	}
}

// Test buildBrowserHeadersData
func TestBuildBrowserHeadersData(t *testing.T) {
	tests := []struct {
		name    string
		headers *schema.Set
		want    int
	}{
		{
			name: "with headers data populated directly",
			headers: testBrowserHeadersSet(
				map[string]interface{}{
					"name":   "X-Custom",
					"value":  "value1",
					"domain": "example.com",
				},
			),
			want: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildBrowserHeadersData(tt.headers)
			if len(got) != tt.want {
				t.Errorf("buildBrowserHeadersData() = %d items, want %d", len(got), tt.want)
			}
		})
	}
}

// Test flattenHostOverridesData
func TestFlattenHostOverridesData(t *testing.T) {
	tests := []struct {
		name         string
		hostOverride *[]sc2.HostOverrides
		want         int
	}{
		{
			name:         "nil host overrides",
			hostOverride: nil,
			want:         0,
		},
		{
			name:         "empty host overrides",
			hostOverride: &[]sc2.HostOverrides{},
			want:         0,
		},
		{
			name: "single host override",
			hostOverride: &[]sc2.HostOverrides{
				{Source: "example.com", Target: "192.168.1.1", KeepHostHeader: true},
			},
			want: 1,
		},
		{
			name: "multiple host overrides",
			hostOverride: &[]sc2.HostOverrides{
				{Source: "api.example.com", Target: "10.0.0.1"},
				{Source: "cdn.example.com", Target: "10.0.0.2", KeepHostHeader: true},
			},
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := flattenHostOverridesData(tt.hostOverride)
			if len(got) != tt.want {
				t.Errorf("flattenHostOverridesData() = %d items, want %d", len(got), tt.want)
			}
		})
	}
}

// Test buildHostOverridesData
func TestBuildHostOverridesData(t *testing.T) {
	tests := []struct {
		name         string
		hostOverride *schema.Set
		want         int
	}{
		{
			name: "single host override",
			hostOverride: testHostOverridesSet(
				map[string]interface{}{
					"source":           "example.com",
					"target":           "192.168.1.1",
					"keep_host_header": true,
				},
			),
			want: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildHostOverridesData(tt.hostOverride)
			if len(got) != tt.want {
				t.Errorf("buildHostOverridesData() = %d items, want %d", len(got), tt.want)
			}
		})
	}
}

// Test flattenValidationsData
func TestFlattenValidationsData(t *testing.T) {
	tests := []struct {
		name        string
		validations *[]sc2.Validations
		want        int
	}{
		{
			name:        "nil validations",
			validations: nil,
			want:        0,
		},
		{
			name:        "empty validations",
			validations: &[]sc2.Validations{},
			want:        0,
		},
		{
			name: "single validation",
			validations: &[]sc2.Validations{
				{
					Name:       "check_status",
					Type:       "response_code",
					Actual:     "200",
					Expected:   "200",
					Comparator: "equals",
				},
			},
			want: 1,
		},
		{
			name: "multiple validations",
			validations: &[]sc2.Validations{
				{
					Name: "check_status",
					Type: "response_code",
				},
				{
					Name:   "check_body",
					Type:   "response_body",
					Source: "body",
				},
				{
					Name:   "check_xpath",
					Type:   "xpath",
					Source: "response",
				},
			},
			want: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := flattenValidationsData(tt.validations)
			if len(got) != tt.want {
				t.Errorf("flattenValidationsData() = %d items, want %d", len(got), tt.want)
			}
		})
	}
}

// Test flattenSslValidationsData
func TestFlattenSslValidationsData(t *testing.T) {
	tests := []struct {
		name        string
		validations *[]sc2.Validations
		want        int
	}{
		{
			name:        "nil ssl validations",
			validations: nil,
			want:        0,
		},
		{
			name:        "empty ssl validations",
			validations: &[]sc2.Validations{},
			want:        0,
		},
		{
			name: "single ssl validation",
			validations: &[]sc2.Validations{
				{
					Name:       "check_cert",
					Type:       "certificate",
					Actual:     "valid",
					Expected:   "valid",
					Comparator: "equals",
				},
			},
			want: 1,
		},
		{
			name: "multiple ssl validations",
			validations: &[]sc2.Validations{
				{Name: "cert_check", Type: "certificate"},
				{Name: "expiry_check", Type: "certificate_expiry"},
			},
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := flattenSslValidationsData(tt.validations)
			if len(got) != tt.want {
				t.Errorf("flattenSslValidationsData() = %d items, want %d", len(got), tt.want)
			}
		})
	}
}

// Test flattenHeaderData
func TestFlattenHeaderData(t *testing.T) {
	tests := []struct {
		name    string
		headers *sc2.Headers
		want    int
	}{
		{
			name:    "nil headers",
			headers: nil,
			want:    0,
		},
		{
			name:    "empty headers",
			headers: &sc2.Headers{},
			want:    0,
		},
		{
			name: "single header",
			headers: &sc2.Headers{
				"X-Custom": "value1",
			},
			want: 1,
		},
		{
			name: "multiple headers",
			headers: &sc2.Headers{
				"X-Custom-1":    "val1",
				"X-Custom-2":    "val2",
				"Authorization": "Bearer token",
			},
			want: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := flattenHeaderData(tt.headers)
			if len(got) != tt.want {
				t.Errorf("flattenHeaderData() = %d items, want %d", len(got), tt.want)
			}
		})
	}
}

// Test flattenLocationData
func TestFlattenLocationData(t *testing.T) {
	tests := []struct {
		name      string
		locations *[]string
		want      int
	}{
		{
			name:      "nil locations",
			locations: nil,
			want:      0,
		},
		{
			name:      "empty locations",
			locations: &[]string{},
			want:      0,
		},
		{
			name: "single location",
			locations: &[]string{
				"us-west-1",
			},
			want: 1,
		},
		{
			name: "multiple locations",
			locations: &[]string{
				"us-west-1",
				"us-east-1",
				"eu-west-1",
			},
			want: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := flattenLocationData(tt.locations)
			if len(got) != tt.want {
				t.Errorf("flattenLocationData() = %d items, want %d", len(got), tt.want)
			}
			for i, loc := range got {
				if _, ok := loc.(string); !ok {
					t.Errorf("item %d is not a string, got %T", i, loc)
				}
			}
		})
	}
}

// Test flattenAdvancedSettingsData
func TestFlattenAdvancedSettingsData(t *testing.T) {
	tests := []struct {
		name     string
		settings *sc2.Advancedsettings
	}{
		{
			name:     "empty settings",
			settings: &sc2.Advancedsettings{},
		},
		{
			name: "settings with verify certificates",
			settings: &sc2.Advancedsettings{
				Verifycertificates: true,
			},
		},
		{
			name: "settings with user agent",
			settings: &sc2.Advancedsettings{
				UserAgent: strPtr("Mozilla/5.0"),
			},
		},
		{
			name: "settings with cookies and headers",
			settings: &sc2.Advancedsettings{
				Cookiesv2: []sc2.Cookiesv2{
					{Key: "session", Value: "abc"},
				},
				BrowserHeaders: []sc2.BrowserHeaders{
					{Name: "X-Custom", Value: "val"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := flattenAdvancedSettingsData(tt.settings)
			if len(got) == 0 {
				t.Fatal("expected at least one item")
			}
			if _, ok := got[0].(map[string]interface{}); !ok {
				t.Fatal("expected map result")
			}
		})
	}
}

// Test buildAdvancedSettingsData
func TestBuildAdvancedSettingsData(t *testing.T) {
	tests := []struct {
		name     string
		settings *schema.Set
		wantErr  bool
	}{
		{
			name: "settings with basic fields",
			settings: testAdvancedSettingsSet(
				map[string]interface{}{
					"verify_certificates":         true,
					"collect_interactive_metrics": false,
					"user_agent":                  "Custom UA",
				},
			),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := buildAdvancedSettingsData(tt.settings)
			if (err != nil) != tt.wantErr {
				t.Errorf("buildAdvancedSettingsData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Test buildBusinessTransactionsData
func TestBuildBusinessTransactionsData(t *testing.T) {
	tests := []struct {
		name string
		txns []interface{}
		want int
		err  bool
	}{
		{
			name: "empty transactions",
			txns: []interface{}{},
			want: 0,
			err:  false,
		},
		{
			name: "single transaction",
			txns: []interface{}{
				map[string]interface{}{
					"name":  "txn1",
					"steps": []interface{}{},
				},
			},
			want: 1,
			err:  false,
		},
		{
			name: "transaction with steps",
			txns: []interface{}{
				map[string]interface{}{
					"name": "txn1",
					"steps": []interface{}{
						map[string]interface{}{
							"name": "step1", "type": "navigate", "url": "https://example.com",
							"wait_for_nav": false, "wait_for_nav_timeout": 0,
							"max_wait_time": 0, "option_selector_type": "", "option_selector": "",
							"variable_name": "", "value": "", "duration": 0,
						},
					},
				},
			},
			want: 1,
			err:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := buildBusinessTransactionsData(tt.txns)
			if (err != nil) != tt.err {
				t.Errorf("buildBusinessTransactionsData() error = %v, wantErr %v", err, tt.err)
				return
			}
			if len(got) != tt.want {
				t.Errorf("buildBusinessTransactionsData() = %d items, want %d", len(got), tt.want)
			}
		})
	}
}

// Test flattenConfigurationData
func TestFlattenConfigurationData(t *testing.T) {
	tests := []struct {
		name   string
		config *sc2.Configuration
	}{
		{
			name:   "nil config",
			config: &sc2.Configuration{},
		},
		{
			name: "populated config",
			config: &sc2.Configuration{
				Name:          "config1",
				URL:           "https://example.com",
				Body:          "request body",
				RequestMethod: "POST",
				Headers: sc2.Headers{
					"X-Custom": "value",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := flattenConfigurationData(tt.config)
			if len(got) == 0 {
				t.Fatal("expected at least one item")
			}
			if _, ok := got[0].(map[string]interface{}); !ok {
				t.Fatal("expected map result")
			}
		})
	}
}

// Test flattenHttpHeadersData
func TestFlattenHttpHeadersData(t *testing.T) {
	tests := []struct {
		name    string
		headers *[]sc2.HttpHeaders
		want    int
	}{
		{
			name:    "nil headers",
			headers: nil,
			want:    0,
		},
		{
			name:    "empty headers",
			headers: &[]sc2.HttpHeaders{},
			want:    0,
		},
		{
			name: "single header",
			headers: &[]sc2.HttpHeaders{
				{Name: "X-Custom", Value: "value1"},
			},
			want: 1,
		},
		{
			name: "multiple headers",
			headers: &[]sc2.HttpHeaders{
				{Name: "X-Custom-1", Value: "val1"},
				{Name: "Authorization", Value: "Bearer token"},
			},
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := flattenHttpHeadersData(tt.headers)
			if len(got) != tt.want {
				t.Errorf("flattenHttpHeadersData() = %d items, want %d", len(got), tt.want)
			}
		})
	}
}

// Test buildHttpHeadersData
func TestBuildHttpHeadersData(t *testing.T) {
	tests := []struct {
		name    string
		headers *schema.Set
		want    int
	}{
		{
			name: "single header",
			headers: testHttpHeadersSet(
				map[string]interface{}{
					"name":  "X-Custom",
					"value": "value1",
				},
			),
			want: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildHttpHeadersData(tt.headers)
			if len(*got) != tt.want {
				t.Errorf("buildHttpHeadersData() = %d items, want %d", len(*got), tt.want)
			}
		})
	}
}

// Test flattenCustomProperties
func TestFlattenCustomProperties(t *testing.T) {
	tests := []struct {
		name  string
		props *[]sc2.CustomProperties
		want  int
	}{
		{
			name:  "nil properties",
			props: nil,
			want:  0,
		},
		{
			name:  "empty properties",
			props: &[]sc2.CustomProperties{},
			want:  0,
		},
		{
			name: "single property",
			props: &[]sc2.CustomProperties{
				{Key: "env", Value: "prod"},
			},
			want: 1,
		},
		{
			name: "multiple properties",
			props: &[]sc2.CustomProperties{
				{Key: "env", Value: "prod"},
				{Key: "team", Value: "platform"},
				{Key: "version", Value: "1.0"},
			},
			want: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := flattenCustomProperties(tt.props)
			if len(got) != tt.want {
				t.Errorf("flattenCustomProperties() = %d items, want %d", len(got), tt.want)
			}
		})
	}
}

// Test dropStaleStepSelectors
func TestDropStaleStepSelectors(t *testing.T) {
	tests := []struct {
		name        string
		step        map[string]interface{}
		wantSelDrop bool
	}{
		{
			name:        "no legacy fields, no stale selectors",
			step:        map[string]interface{}{},
			wantSelDrop: false,
		},
		{
			name: "legacy fields but no selector field",
			step: map[string]interface{}{
				"selector_type": "id",
				"selector":      "btn",
			},
			wantSelDrop: false,
		},
		{
			name: "stale selector matching legacy",
			step: map[string]interface{}{
				"selector_type": "id",
				"selector":      "btn",
				"selectors": []interface{}{
					map[string]interface{}{"type": "id", "value": "btn"},
				},
			},
			wantSelDrop: false,
		},
		{
			name: "stale selector not matching legacy",
			step: map[string]interface{}{
				"selector_type": "id",
				"selector":      "btn",
				"selectors": []interface{}{
					map[string]interface{}{"type": "css", "value": ".btn"},
				},
			},
			wantSelDrop: true,
		},
		{
			name: "multiple selectors not dropped",
			step: map[string]interface{}{
				"selector_type": "id",
				"selector":      "btn",
				"selectors": []interface{}{
					map[string]interface{}{"type": "id", "value": "btn1"},
					map[string]interface{}{"type": "id", "value": "btn2"},
				},
			},
			wantSelDrop: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dropStaleStepSelectors(tt.step)
			_, hasSelectors := tt.step["selectors"]
			if tt.wantSelDrop && hasSelectors {
				t.Errorf("expected selectors to be dropped but they remain")
			}
			if !tt.wantSelDrop && !hasSelectors && len(tt.step) > 2 {
				// if originally had selectors but dropped when shouldn't
				t.Errorf("expected selectors to remain but were dropped")
			}
		})
	}
}

// Test buildStepV2Data
func TestBuildStepV2Data(t *testing.T) {
	tests := []struct {
		name  string
		steps []interface{}
		want  int
		err   bool
	}{
		{
			name:  "empty steps",
			steps: []interface{}{},
			want:  0,
			err:   false,
		},
		{
			name: "single step",
			steps: []interface{}{
				map[string]interface{}{
					"name": "step1", "type": "navigate", "url": "https://example.com",
					"wait_for_nav": false, "wait_for_nav_timeout": 0,
					"max_wait_time": 0, "option_selector_type": "", "option_selector": "",
					"variable_name": "", "value": "", "duration": 0,
				},
			},
			want: 1,
			err:  false,
		},
		{
			name: "multiple steps",
			steps: []interface{}{
				map[string]interface{}{
					"name": "step1", "type": "navigate", "url": "https://example.com",
					"wait_for_nav": false, "wait_for_nav_timeout": 0,
					"max_wait_time": 0, "option_selector_type": "", "option_selector": "",
					"variable_name": "", "value": "", "duration": 0,
				},
				map[string]interface{}{
					"name": "step2", "type": "click_element", "url": "",
					"wait_for_nav": true, "wait_for_nav_timeout": 5000,
					"max_wait_time": 10000, "option_selector_type": "", "option_selector": "",
					"variable_name": "", "value": "", "duration": 0,
				},
			},
			want: 2,
			err:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := buildStepV2Data(tt.steps)
			if (err != nil) != tt.err {
				t.Errorf("buildStepV2Data() error = %v, wantErr %v", err, tt.err)
				return
			}
			if len(got) != tt.want {
				t.Errorf("buildStepV2Data() = %d items, want %d", len(got), tt.want)
			}
		})
	}
}

// Test buildBrowserV2Data - integration test done in resource tests
func TestBuildBrowserV2Data(t *testing.T) {
	// This function requires complex ResourceData setup which is tested
	// through resource integration tests. Pure function logic is covered
	// by other unit tests.
	t.Run("buildBrowserV2Data exists", func(t *testing.T) {
		// Placeholder to indicate function exists and is tested elsewhere
	})
}

// Helper function
func strPtr(s string) *string {
	return &s
}

// testBrowserHeadersSet creates a Set with proper BrowserHeaders schema
func testBrowserHeadersSet(items ...map[string]interface{}) *schema.Set {
	itemList := make([]interface{}, len(items))
	for i, item := range items {
		itemList[i] = item
	}
	// Create a resource with the expected schema for browser headers
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name":   {Type: schema.TypeString},
			"value":  {Type: schema.TypeString},
			"domain": {Type: schema.TypeString},
		},
	}
	return schema.NewSet(schema.HashResource(resource), itemList)
}

// testCookiesSet creates a Set with proper Cookies schema
func testCookiesSet(items ...map[string]interface{}) *schema.Set {
	itemList := make([]interface{}, len(items))
	for i, item := range items {
		itemList[i] = item
	}
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"key":    {Type: schema.TypeString},
			"value":  {Type: schema.TypeString},
			"domain": {Type: schema.TypeString},
			"path":   {Type: schema.TypeString},
		},
	}
	return schema.NewSet(schema.HashResource(resource), itemList)
}

// testHostOverridesSet creates a Set with proper HostOverrides schema
func testHostOverridesSet(items ...map[string]interface{}) *schema.Set {
	itemList := make([]interface{}, len(items))
	for i, item := range items {
		itemList[i] = item
	}
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"source":           {Type: schema.TypeString},
			"target":           {Type: schema.TypeString},
			"keep_host_header": {Type: schema.TypeBool},
		},
	}
	return schema.NewSet(schema.HashResource(resource), itemList)
}

// testHttpHeadersSet creates a Set with proper HttpHeaders schema
func testHttpHeadersSet(items ...map[string]interface{}) *schema.Set {
	itemList := make([]interface{}, len(items))
	for i, item := range items {
		itemList[i] = item
	}
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name":  {Type: schema.TypeString},
			"value": {Type: schema.TypeString},
		},
	}
	return schema.NewSet(schema.HashResource(resource), itemList)
}

// testAdvancedSettingsSet creates a Set with proper AdvancedSettings schema
func testAdvancedSettingsSet(items ...map[string]interface{}) *schema.Set {
	itemList := make([]interface{}, len(items))
	for i, item := range items {
		// Set default empty sets for required nested fields
		if _, ok := item["authentication"]; !ok {
			item["authentication"] = schema.NewSet(schema.HashResource(
				&schema.Resource{
					Schema: map[string]*schema.Schema{
						"username": {Type: schema.TypeString},
						"password": {Type: schema.TypeString},
					},
				},
			), nil)
		}
		if _, ok := item["headers"]; !ok {
			item["headers"] = schema.NewSet(schema.HashResource(
				&schema.Resource{
					Schema: map[string]*schema.Schema{
						"name":   {Type: schema.TypeString},
						"value":  {Type: schema.TypeString},
						"domain": {Type: schema.TypeString},
					},
				},
			), nil)
		}
		if _, ok := item["cookies"]; !ok {
			item["cookies"] = schema.NewSet(schema.HashResource(
				&schema.Resource{
					Schema: map[string]*schema.Schema{
						"key":    {Type: schema.TypeString},
						"value":  {Type: schema.TypeString},
						"domain": {Type: schema.TypeString},
						"path":   {Type: schema.TypeString},
					},
				},
			), nil)
		}
		if _, ok := item["host_overrides"]; !ok {
			item["host_overrides"] = schema.NewSet(schema.HashResource(
				&schema.Resource{
					Schema: map[string]*schema.Schema{
						"source":           {Type: schema.TypeString},
						"target":           {Type: schema.TypeString},
						"keep_host_header": {Type: schema.TypeBool},
					},
				},
			), nil)
		}
		if _, ok := item["chrome_flags"]; !ok {
			item["chrome_flags"] = schema.NewSet(schema.HashResource(
				&schema.Resource{
					Schema: map[string]*schema.Schema{
						"name":  {Type: schema.TypeString},
						"value": {Type: schema.TypeString},
					},
				},
			), nil)
		}
		if _, ok := item["certificate_ids"]; !ok {
			item["certificate_ids"] = []interface{}{}
		}
		if _, ok := item["excluded_files"]; !ok {
			item["excluded_files"] = schema.NewSet(schema.HashResource(
				&schema.Resource{
					Schema: map[string]*schema.Schema{
						"type":  {Type: schema.TypeString},
						"regex": {Type: schema.TypeString},
					},
				},
			), nil)
		}
		itemList[i] = item
	}
	resource := browserCheckV2AdvancedSettingsResource(false)
	return schema.NewSet(schema.HashResource(resource), itemList)
}
