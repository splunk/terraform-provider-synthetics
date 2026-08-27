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
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sc2 "github.com/splunk/syntheticsclient/v3/syntheticsclientv2"
)

// Tests for flattenIdData
func TestFlattenIdData(t *testing.T) {
	t.Run("single id in set", func(t *testing.T) {
		testSet := schema.NewSet(
			func(i interface{}) int {
				return i.(map[string]interface{})["id"].(int)
			},
			[]interface{}{
				map[string]interface{}{
					"id": 42,
				},
			},
		)
		got := flattenIdData(testSet)
		if got != 42 {
			t.Fatalf("flattenIdData() = %d, want 42", got)
		}
	})
}

// Tests for flattenStringIdData
func TestFlattenStringIdData(t *testing.T) {
	t.Run("string id in set", func(t *testing.T) {
		testSet := schema.NewSet(
			func(i interface{}) int { return 0 },
			[]interface{}{
				map[string]interface{}{
					"id": "location-123",
				},
			},
		)
		got := flattenStringIdData(testSet)
		if got != "location-123" {
			t.Fatalf("flattenStringIdData() = %s, want location-123", got)
		}
	})
}

// Tests for flattenApiV2Read
func TestFlattenApiV2Read(t *testing.T) {
	t.Run("minimal api response", func(t *testing.T) {
		response := &sc2.ApiCheckV2Response{}
		response.Test.Active = true
		response.Test.Automaticretries = 0
		response.Test.Frequency = 0
		response.Test.Name = ""
		response.Test.Schedulingstrategy = ""
		response.Test.Deviceid = 0
		response.Test.Locationids = []string{}
		response.Test.Requests = []sc2.Requests{}
		response.Test.Customproperties = []sc2.CustomProperties{}

		result := flattenApiV2Read(response)
		if len(result) != 1 {
			t.Fatalf("flattenApiV2Read() len = %d, want 1", len(result))
		}
		m := result[0].(map[string]interface{})
		if m["active"] != true || m["automatic_retries"] != 0 {
			t.Fatalf("flattenApiV2Read() = %v", m)
		}
	})

	t.Run("full api response", func(t *testing.T) {
		response := &sc2.ApiCheckV2Response{}
		response.Test.Active = true
		response.Test.Automaticretries = 2
		response.Test.Frequency = 60
		response.Test.Name = "test-api-check"
		response.Test.Schedulingstrategy = "auto"
		response.Test.Deviceid = 5
		response.Test.Locationids = []string{"loc1", "loc2"}
		response.Test.Requests = []sc2.Requests{}
		response.Test.Customproperties = []sc2.CustomProperties{}

		result := flattenApiV2Read(response)
		if len(result) != 1 {
			t.Fatalf("flattenApiV2Read() len = %d, want 1", len(result))
		}
		m := result[0].(map[string]interface{})
		if m["name"] != "test-api-check" || m["frequency"] != 60 || m["device_id"] != 5 {
			t.Fatalf("flattenApiV2Read() = %v", m)
		}
	})
}

// Tests for findDeviceByID
func TestFindDeviceByID(t *testing.T) {
	testCases := []struct {
		name       string
		devices    []sc2.Device
		deviceID   int
		shouldFind bool
	}{
		{
			name:       "device found",
			deviceID:   2,
			shouldFind: true,
			devices: []sc2.Device{
				{ID: 1, Label: "device1"},
				{ID: 2, Label: "device2"},
				{ID: 3, Label: "device3"},
			},
		},
		{
			name:       "device not found",
			deviceID:   99,
			shouldFind: false,
			devices: []sc2.Device{
				{ID: 1, Label: "device1"},
				{ID: 2, Label: "device2"},
			},
		},
		{
			name:       "empty device list",
			deviceID:   1,
			shouldFind: false,
			devices:    []sc2.Device{},
		},
		{
			name:       "find first device",
			deviceID:   1,
			shouldFind: true,
			devices: []sc2.Device{
				{ID: 1, Label: "device1"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := findDeviceByID(tc.devices, tc.deviceID)
			if tc.shouldFind {
				if got == nil {
					t.Fatalf("findDeviceByID() = nil, want device with ID %d", tc.deviceID)
				}
				if got.ID != tc.deviceID {
					t.Fatalf("findDeviceByID() = %d, want %d", got.ID, tc.deviceID)
				}
			} else {
				if got != nil {
					t.Fatalf("findDeviceByID() = %v, want nil", got)
				}
			}
		})
	}
}

// Tests for flattenDeviceFromID
func TestFlattenDeviceFromID(t *testing.T) {
	testCases := []struct {
		name     string
		deviceID int
		devices  []sc2.Device
		expLen   int
	}{
		{
			name:     "device found in list",
			deviceID: 2,
			devices: []sc2.Device{
				{ID: 1, Label: "device1"},
				{ID: 2, Label: "device2"},
			},
			expLen: 1,
		},
		{
			name:     "device not found, non-zero id",
			deviceID: 99,
			devices: []sc2.Device{
				{ID: 1, Label: "device1"},
			},
			expLen: 1,
		},
		{
			name:     "zero device id",
			deviceID: 0,
			devices: []sc2.Device{
				{ID: 1, Label: "device1"},
			},
			expLen: 0,
		},
		{
			name:     "empty device list with zero id",
			deviceID: 0,
			devices:  []sc2.Device{},
			expLen:   0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := flattenDeviceFromID(tc.deviceID, tc.devices)
			if len(result) != tc.expLen {
				t.Fatalf("flattenDeviceFromID() len = %d, want %d", len(result), tc.expLen)
			}
		})
	}
}

// Tests for flattenDeviceData
func TestFlattenDeviceData(t *testing.T) {
	t.Run("zero device", func(t *testing.T) {
		device := &sc2.Device{
			ID:                0,
			Label:             "",
			UserAgent:         "",
			Networkconnection: sc2.Networkconnection{},
			Viewportheight:    0,
			Viewportwidth:     0,
		}
		result := flattenDeviceData(device)
		if len(result) != 1 {
			t.Fatalf("flattenDeviceData() len = %d, want 1", len(result))
		}
	})

	t.Run("populated device", func(t *testing.T) {
		device := &sc2.Device{
			ID:             5,
			Label:          "Chrome Desktop",
			UserAgent:      "Mozilla/5.0...",
			Viewportheight: 768,
			Viewportwidth:  1024,
			Networkconnection: sc2.Networkconnection{
				Description:       "Test Network",
				Downloadbandwidth: 100,
				Uploadbandwidth:   50,
				Latency:           20,
				Packetloss:        1,
			},
		}
		result := flattenDeviceData(device)
		if len(result) != 1 {
			t.Fatalf("flattenDeviceData() len = %d, want 1", len(result))
		}
		m := result[0].(map[string]interface{})
		if m["id"] != 5 || m["label"] != "Chrome Desktop" || m["viewport_height"] != 768 {
			t.Fatalf("flattenDeviceData() = %v", m)
		}
	})
}

// Tests for flattenDevicesV2Data
func TestFlattenDevicesV2Data(t *testing.T) {
	testCases := []struct {
		name    string
		devices *[]sc2.Device
		expLen  int
	}{
		{
			name:    "nil devices",
			devices: nil,
			expLen:  0,
		},
		{
			name:    "empty device list",
			devices: &[]sc2.Device{},
			expLen:  0,
		},
		{
			name: "multiple devices",
			devices: &[]sc2.Device{
				{ID: 1, Label: "device1"},
				{ID: 2, Label: "device2"},
				{ID: 3, Label: "device3"},
			},
			expLen: 3,
		},
		{
			name: "single device",
			devices: &[]sc2.Device{
				{ID: 42, Label: "single-device"},
			},
			expLen: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := flattenDevicesV2Data(tc.devices)
			if len(result) != tc.expLen {
				t.Fatalf("flattenDevicesV2Data() len = %d, want %d", len(result), tc.expLen)
			}
		})
	}
}

// Tests for flattenApiV2Data
func TestFlattenApiV2Data(t *testing.T) {
	t.Run("minimal api response with no devices", func(t *testing.T) {
		response := &sc2.ApiCheckV2Response{}
		response.Test.Active = true
		response.Test.Automaticretries = 0
		response.Test.Frequency = 0
		response.Test.Locationids = []string{}
		response.Test.Requests = []sc2.Requests{}
		response.Test.Customproperties = []sc2.CustomProperties{}

		result := flattenApiV2Data(response, []sc2.Device{})
		if len(result) != 1 {
			t.Fatalf("flattenApiV2Data() len = %d, want 1", len(result))
		}
		m := result[0].(map[string]interface{})
		if m["active"] != true {
			t.Fatalf("flattenApiV2Data() = %v", m)
		}
	})

	t.Run("full api response with timestamps", func(t *testing.T) {
		response := &sc2.ApiCheckV2Response{}
		response.Test.Active = true
		response.Test.Createdat = time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
		response.Test.Updatedat = time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC)
		response.Test.Lastrunat = time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC)
		response.Test.ID = 123
		response.Test.Name = "test-api"
		response.Test.Frequency = 60
		response.Test.Type = "api"
		response.Test.Lastrunstatus = "success"
		response.Test.Createdby = "user1"
		response.Test.Updatedby = "user2"
		response.Test.Automaticretries = 2
		response.Test.Schedulingstrategy = "round_robin"
		response.Test.Locationids = []string{"loc1"}
		response.Test.Requests = []sc2.Requests{}
		response.Test.Customproperties = []sc2.CustomProperties{}

		result := flattenApiV2Data(response, []sc2.Device{})
		if len(result) != 1 {
			t.Fatalf("flattenApiV2Data() len = %d, want 1", len(result))
		}
		m := result[0].(map[string]interface{})
		if m["id"] != 123 || m["name"] != "test-api" || m["last_run_status"] != "success" || m["created_by"] != "user1" {
			t.Fatalf("flattenApiV2Data() = %v", m)
		}
	})
}

// Tests for flattenVariableV2Read
func TestFlattenVariableV2Read(t *testing.T) {
	t.Run("basic variable", func(t *testing.T) {
		input := &sc2.VariableV2Response{
			Variable: sc2.Variable{
				Name:        "test_var",
				Description: "A test variable",
				Value:       "test_value",
				Secret:      false,
			},
		}
		result := flattenVariableV2Read(input)
		if len(result) != 1 {
			t.Fatalf("flattenVariableV2Read() len = %d, want 1", len(result))
		}
		m := result[0].(map[string]interface{})
		if m["name"] != "test_var" || m["value"] != "test_value" || m["secret"] != false {
			t.Fatalf("flattenVariableV2Read() = %v", m)
		}
	})

	t.Run("secret variable", func(t *testing.T) {
		input := &sc2.VariableV2Response{
			Variable: sc2.Variable{
				Name:        "secret_var",
				Description: "A secret",
				Value:       "***",
				Secret:      true,
			},
		}
		result := flattenVariableV2Read(input)
		if len(result) != 1 {
			t.Fatalf("flattenVariableV2Read() len = %d, want 1", len(result))
		}
		m := result[0].(map[string]interface{})
		if m["secret"] != true {
			t.Fatalf("flattenVariableV2Read() secret = %v, want true", m["secret"])
		}
	})
}

// Tests for flattenVariableV2Data
func TestFlattenVariableV2Data(t *testing.T) {
	t.Run("variable with timestamps", func(t *testing.T) {
		input := &sc2.VariableV2Response{
			Variable: sc2.Variable{
				ID:          42,
				Name:        "test_var",
				Description: "A test variable",
				Value:       "test_value",
				Secret:      false,
				Createdat:   time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				Updatedat:   time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC),
			},
		}
		result := flattenVariableV2Data(input)
		if len(result) != 1 {
			t.Fatalf("flattenVariableV2Data() len = %d, want 1", len(result))
		}
		m := result[0].(map[string]interface{})
		if m["id"] != 42 || m["name"] != "test_var" {
			t.Fatalf("flattenVariableV2Data() = %v", m)
		}
		// Check timestamps are present
		if m["created_at"] == "" || m["updated_at"] == "" {
			t.Fatalf("flattenVariableV2Data() timestamps empty")
		}
	})

	t.Run("variable with zero timestamps", func(t *testing.T) {
		input := &sc2.VariableV2Response{
			Variable: sc2.Variable{
				ID:        10,
				Name:      "var2",
				Value:     "value2",
				Secret:    false,
				Createdat: time.Time{},
				Updatedat: time.Time{},
			},
		}
		result := flattenVariableV2Data(input)
		if len(result) != 1 {
			t.Fatalf("flattenVariableV2Data() len = %d, want 1", len(result))
		}
		m := result[0].(map[string]interface{})
		if m["id"] != 10 {
			t.Fatalf("flattenVariableV2Data() id = %v, want 10", m["id"])
		}
	})
}

// Tests for flattenVariablesV2Data
func TestFlattenVariablesV2Data(t *testing.T) {
	testCases := []struct {
		name   string
		input  *[]sc2.Variable
		expLen int
	}{
		{
			name:   "nil variables",
			input:  nil,
			expLen: 0,
		},
		{
			name:   "empty variables",
			input:  &[]sc2.Variable{},
			expLen: 0,
		},
		{
			name: "multiple variables",
			input: &[]sc2.Variable{
				{
					ID:        1,
					Name:      "var1",
					Value:     "value1",
					Secret:    false,
					Createdat: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
					Updatedat: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				},
				{
					ID:        2,
					Name:      "var2",
					Value:     "value2",
					Secret:    true,
					Createdat: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
					Updatedat: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			expLen: 2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := flattenVariablesV2Data(tc.input)
			if len(result) != tc.expLen {
				t.Fatalf("flattenVariablesV2Data() len = %d, want %d", len(result), tc.expLen)
			}
		})
	}
}

// Tests for flattenRequestData
func TestFlattenRequestData(t *testing.T) {
	testCases := []struct {
		name   string
		input  *[]sc2.Requests
		expLen int
	}{
		{
			name:   "nil requests",
			input:  nil,
			expLen: 0,
		},
		{
			name:   "empty requests",
			input:  &[]sc2.Requests{},
			expLen: 0,
		},
		{
			name: "single request",
			input: &[]sc2.Requests{
				{
					Configuration: sc2.Configuration{
						Name:          "test",
						RequestMethod: "GET",
						URL:           "https://example.com",
						Body:          "",
					},
					Setup:       []sc2.Setup{},
					Validations: []sc2.Validations{},
				},
			},
			expLen: 1,
		},
		{
			name: "multiple requests",
			input: &[]sc2.Requests{
				{
					Configuration: sc2.Configuration{
						Name:          "req1",
						RequestMethod: "GET",
						URL:           "https://example.com/1",
					},
					Setup:       []sc2.Setup{},
					Validations: []sc2.Validations{},
				},
				{
					Configuration: sc2.Configuration{
						Name:          "req2",
						RequestMethod: "POST",
						URL:           "https://example.com/2",
					},
					Setup:       []sc2.Setup{},
					Validations: []sc2.Validations{},
				},
			},
			expLen: 2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := flattenRequestData(tc.input)
			if len(result) != tc.expLen {
				t.Fatalf("flattenRequestData() len = %d, want %d", len(result), tc.expLen)
			}
		})
	}
}

// Tests for timeString
func TestTimeString(t *testing.T) {
	t.Run("zero time", func(t *testing.T) {
		result := timeString(time.Time{})
		if result != "" {
			t.Fatalf("timeString() = %q, want empty string", result)
		}
	})

	t.Run("non-zero time", func(t *testing.T) {
		result := timeString(time.Date(2021, 1, 1, 12, 0, 0, 0, time.UTC))
		if result == "" {
			t.Fatalf("timeString() = empty string, want non-empty")
		}
	})
}

// Tests for flattenNetworkConnectionData
func TestFlattenNetworkConnectionData(t *testing.T) {
	t.Run("empty network connection", func(t *testing.T) {
		input := &sc2.Networkconnection{
			Description:       "",
			Downloadbandwidth: 0,
			Uploadbandwidth:   0,
			Latency:           0,
			Packetloss:        0,
		}
		result := flattenNetworkConnectionData(input)
		if len(result) != 1 {
			t.Fatalf("flattenNetworkConnectionData() len = %d, want 1", len(result))
		}
		m := result[0].(map[string]interface{})
		if m["description"] != "" || m["download_bandwidth"] != 0 {
			t.Fatalf("flattenNetworkConnectionData() = %v", m)
		}
	})

	t.Run("populated network connection", func(t *testing.T) {
		input := &sc2.Networkconnection{
			Description:       "LTE",
			Downloadbandwidth: 10000,
			Uploadbandwidth:   5000,
			Latency:           50,
			Packetloss:        1,
		}
		result := flattenNetworkConnectionData(input)
		if len(result) != 1 {
			t.Fatalf("flattenNetworkConnectionData() len = %d, want 1", len(result))
		}
		m := result[0].(map[string]interface{})
		if m["description"] != "LTE" || m["download_bandwidth"] != 10000 || m["upload_bandwidth"] != 5000 {
			t.Fatalf("flattenNetworkConnectionData() = %v", m)
		}
	})
}

// Tests for flattenAuthenticationData
func TestFlattenAuthenticationData(t *testing.T) {
	t.Run("empty authentication", func(t *testing.T) {
		input := &sc2.Authentication{
			Username: "",
			Password: "",
		}
		result := flattenAuthenticationData(input)
		if len(result) != 1 {
			t.Fatalf("flattenAuthenticationData() len = %d, want 1", len(result))
		}
	})

	t.Run("username only", func(t *testing.T) {
		input := &sc2.Authentication{
			Username: "testuser",
			Password: "",
		}
		result := flattenAuthenticationData(input)
		if len(result) != 1 {
			t.Fatalf("flattenAuthenticationData() len = %d, want 1", len(result))
		}
		m := result[0].(map[string]interface{})
		if m["username"] != "testuser" {
			t.Fatalf("flattenAuthenticationData() = %v", m)
		}
	})

	t.Run("full authentication", func(t *testing.T) {
		input := &sc2.Authentication{
			Username: "testuser",
			Password: "testpass",
		}
		result := flattenAuthenticationData(input)
		if len(result) != 1 {
			t.Fatalf("flattenAuthenticationData() len = %d, want 1", len(result))
		}
		m := result[0].(map[string]interface{})
		if m["username"] != "testuser" || m["password"] != "testpass" {
			t.Fatalf("flattenAuthenticationData() = %v", m)
		}
	})
}

// Tests for buildAuthenticationData
func TestBuildAuthenticationData(t *testing.T) {
	t.Run("empty authentication set", func(t *testing.T) {
		input := schema.NewSet(
			func(i interface{}) int { return 0 },
			[]interface{}{},
		)
		result := buildAuthenticationData(input)
		if result != nil {
			t.Fatalf("buildAuthenticationData() = %v, want nil", result)
		}
	})

	t.Run("full authentication", func(t *testing.T) {
		input := schema.NewSet(
			func(i interface{}) int { return 0 },
			[]interface{}{
				map[string]interface{}{
					"username": "testuser",
					"password": "testpass",
				},
			},
		)
		result := buildAuthenticationData(input)
		if result == nil {
			t.Fatalf("buildAuthenticationData() = nil, want auth")
		}
		if result.Username != "testuser" || result.Password != "testpass" {
			t.Fatalf("buildAuthenticationData() = %v", result)
		}
	})

	t.Run("username only", func(t *testing.T) {
		input := schema.NewSet(
			func(i interface{}) int { return 0 },
			[]interface{}{
				map[string]interface{}{
					"username": "user",
					"password": "",
				},
			},
		)
		result := buildAuthenticationData(input)
		if result == nil {
			t.Fatalf("buildAuthenticationData() = nil, want auth")
		}
		if result.Username != "user" {
			t.Fatalf("buildAuthenticationData() = %v", result)
		}
	})
}

func TestBuildApiV2DataWithRequestsAndCustomProperties(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceApiCheckV2().Schema, map[string]interface{}{
		"test": []interface{}{
			map[string]interface{}{
				"name":                "api-create",
				"active":              true,
				"device_id":           1,
				"frequency":           5,
				"automatic_retries":   1,
				"location_ids":        []interface{}{"aws-us-east-1"},
				"scheduling_strategy": "round_robin",
				"requests": []interface{}{
					map[string]interface{}{
						"configuration": []interface{}{
							map[string]interface{}{
								"name":           "get-config",
								"url":            "https://api.example.com/endpoint",
								"request_method": "GET",
								"body":           "",
								"headers": map[string]interface{}{
									"Authorization": "Bearer token",
								},
								"certificate_id": 0,
							},
						},
						"setup": []interface{}{
							map[string]interface{}{
								"name":      "setup-step",
								"type":      "variable_from_response",
								"extractor": "json",
								"source":    "body",
								"variable":  "response_var",
								"code":      "",
								"value":     "",
							},
						},
						"validations": []interface{}{
							map[string]interface{}{
								"name":       "assert-status",
								"type":       "assert_numeric",
								"actual":     "{{response.code}}",
								"comparator": "equals",
								"expected":   "200",
								"extractor":  "",
								"source":     "",
								"variable":   "",
								"value":      "",
								"code":       "",
							},
						},
					},
				},
				"custom_properties": []interface{}{
					map[string]interface{}{
						"key":   "owner",
						"value": "api-team",
					},
				},
			},
		},
	})

	got := buildApiV2Data(d)

	if got.Test.Name != "api-create" {
		t.Fatalf("Name = %v, want api-create", got.Test.Name)
	}
	if got.Test.Active != true {
		t.Fatalf("Active = %v, want true", got.Test.Active)
	}
	if got.Test.Deviceid != 1 {
		t.Fatalf("Deviceid = %v, want 1", got.Test.Deviceid)
	}
	if got.Test.Frequency != 5 {
		t.Fatalf("Frequency = %v, want 5", got.Test.Frequency)
	}
	if got.Test.Automaticretries != 1 {
		t.Fatalf("Automaticretries = %v, want 1", got.Test.Automaticretries)
	}
	if len(got.Test.Locationids) != 1 || got.Test.Locationids[0] != "aws-us-east-1" {
		t.Fatalf("Locationids = %v, want [aws-us-east-1]", got.Test.Locationids)
	}
	if got.Test.Schedulingstrategy != "round_robin" {
		t.Fatalf("Schedulingstrategy = %v, want round_robin", got.Test.Schedulingstrategy)
	}
	if len(got.Test.Requests) != 1 {
		t.Fatalf("Requests len = %d, want 1", len(got.Test.Requests))
	}
	req := got.Test.Requests[0]
	if req.Name != "get-config" || req.URL != "https://api.example.com/endpoint" || req.RequestMethod != "GET" {
		t.Fatalf("Configuration = %#v, want get-config/https/GET", req.Configuration)
	}
	if len(req.Setup) != 1 {
		t.Fatalf("Setup len = %d, want 1", len(req.Setup))
	}
	if req.Setup[0].Name != "setup-step" || req.Setup[0].Type != "variable_from_response" {
		t.Fatalf("Setup = %#v, want setup-step", req.Setup[0])
	}
	if len(req.Validations) != 1 {
		t.Fatalf("Validations len = %d, want 1", len(req.Validations))
	}
	if req.Validations[0].Name != "assert-status" || req.Validations[0].Type != "assert_numeric" {
		t.Fatalf("Validation = %#v, want assert-status", req.Validations[0])
	}
	if len(got.Test.Customproperties) != 1 || got.Test.Customproperties[0].Key != "owner" || got.Test.Customproperties[0].Value != "api-team" {
		t.Fatalf("Customproperties = %#v, want owner=api-team", got.Test.Customproperties)
	}
}

func TestBuildVariableV2Data(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceVariableV2().Schema, map[string]interface{}{
		"variable": []interface{}{
			map[string]interface{}{
				"description": "my api key",
				"name":        "api_key",
				"secret":      true,
				"value":       "super-secret-value",
			},
		},
	})

	got := buildVariableV2Data(d)

	if got.Description != "my api key" {
		t.Fatalf("Description = %q, want my api key", got.Description)
	}
	if got.Name != "api_key" {
		t.Fatalf("Name = %q, want api_key", got.Name)
	}
	if got.Secret != true {
		t.Fatalf("Secret = %v, want true", got.Secret)
	}
	if got.Value != "super-secret-value" {
		t.Fatalf("Value = %q, want super-secret-value", got.Value)
	}
}

func TestBuildRequestsDataWithContent(t *testing.T) {
	requests := []interface{}{
		map[string]interface{}{
			"configuration": []interface{}{
				map[string]interface{}{
					"name":           "config1",
					"url":            "https://example.com",
					"request_method": "POST",
					"body":           `{"key": "value"}`,
					"headers": map[string]interface{}{
						"Content-Type": "application/json",
					},
					"certificate_id": 0,
				},
			},
			"setup": []interface{}{
				map[string]interface{}{
					"name":      "setup1",
					"type":      "variable_from_response",
					"extractor": "json_path",
					"source":    "body",
					"variable":  "extracted_var",
					"code":      "",
					"value":     "",
				},
			},
			"validations": []interface{}{
				map[string]interface{}{
					"name":       "validation1",
					"type":       "assert_numeric",
					"actual":     "{{response.code}}",
					"comparator": "equals",
					"expected":   "200",
					"extractor":  "",
					"source":     "",
					"variable":   "",
					"value":      "",
					"code":       "",
				},
			},
		},
	}

	got := buildRequestsData(requests)

	if len(got) != 1 {
		t.Fatalf("Requests len = %d, want 1", len(got))
	}
	req := got[0]
	if req.Name != "config1" || req.URL != "https://example.com" || req.RequestMethod != "POST" {
		t.Fatalf("Configuration = %#v, want config1/https/POST", req.Configuration)
	}
	if req.Body != `{"key": "value"}` {
		t.Fatalf("Configuration body = %q, want json", req.Body)
	}
	if len(req.Setup) != 1 {
		t.Fatalf("Setup len = %d, want 1", len(req.Setup))
	}
	if req.Setup[0].Name != "setup1" || req.Setup[0].Variable != "extracted_var" {
		t.Fatalf("Setup = %#v, want setup1", req.Setup[0])
	}
	if len(req.Validations) != 1 {
		t.Fatalf("Validations len = %d, want 1", len(req.Validations))
	}
	if req.Validations[0].Name != "validation1" || req.Validations[0].Comparator != "equals" {
		t.Fatalf("Validation = %#v, want validation1", req.Validations[0])
	}
}

func TestBuildRequestsDataEmptySlice(t *testing.T) {
	requests := []interface{}{}
	got := buildRequestsData(requests)
	if len(got) != 0 {
		t.Fatalf("Requests len = %d, want 0", len(got))
	}
}
