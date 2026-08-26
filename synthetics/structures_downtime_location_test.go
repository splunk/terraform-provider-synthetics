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

// TestBuildDowntimeConfigurationV2Data tests buildDowntimeConfigurationV2Data with various inputs
func TestBuildDowntimeConfigurationV2Data(t *testing.T) {
	tests := []struct {
		name     string
		config   map[string]interface{}
		validate func(t *testing.T, result sc2.DowntimeConfigurationV2Input)
	}{
		{
			name: "populated downtime configuration",
			config: map[string]interface{}{
				"downtime_configuration": []interface{}{
					map[string]interface{}{
						"name":        "test downtime",
						"description": "test description",
						"rule":        "ALWAYS",
						"start_time":  "2024-01-01T00:00:00.000Z",
						"end_time":    "2024-01-02T00:00:00.000Z",
						"timezone":    "UTC",
						"test_ids":    []interface{}{1, 2, 3},
						"recurrence":  []interface{}{},
					},
				},
			},
			validate: func(t *testing.T, result sc2.DowntimeConfigurationV2Input) {
				if result.Name != "test downtime" {
					t.Fatalf("expected name 'test downtime', got %q", result.Name)
				}
				if result.Description != "test description" {
					t.Fatalf("expected description 'test description', got %q", result.Description)
				}
				if result.Rule != "ALWAYS" {
					t.Fatalf("expected rule 'ALWAYS', got %q", result.Rule)
				}
				if result.Timezone == nil || *result.Timezone != "UTC" {
					t.Fatalf("expected timezone 'UTC', got %v", result.Timezone)
				}
				if len(result.Testids) != 3 {
					t.Fatalf("expected 3 test ids, got %d", len(result.Testids))
				}
			},
		},
		{
			name: "minimal downtime configuration",
			config: map[string]interface{}{
				"downtime_configuration": []interface{}{
					map[string]interface{}{
						"name":        "minimal downtime",
						"description": "",
						"rule":        "ALWAYS",
						"start_time":  "2024-01-01T00:00:00.000Z",
						"end_time":    "2024-01-02T00:00:00.000Z",
						"test_ids":    []interface{}{},
						"recurrence":  []interface{}{},
					},
				},
			},
			validate: func(t *testing.T, result sc2.DowntimeConfigurationV2Input) {
				if result.Name != "minimal downtime" {
					t.Fatalf("expected name 'minimal downtime', got %q", result.Name)
				}
				if result.Description != "" {
					t.Fatalf("expected empty description, got %q", result.Description)
				}
				if len(result.Testids) != 0 {
					t.Fatalf("expected 0 test ids, got %d", len(result.Testids))
				}
			},
		},
		{
			name: "downtime configuration with single test id",
			config: map[string]interface{}{
				"downtime_configuration": []interface{}{
					map[string]interface{}{
						"name":        "single test",
						"description": "test",
						"rule":        "ALWAYS",
						"start_time":  "2024-01-01T00:00:00.000Z",
						"end_time":    "2024-01-02T00:00:00.000Z",
						"test_ids":    []interface{}{42},
						"recurrence":  []interface{}{},
					},
				},
			},
			validate: func(t *testing.T, result sc2.DowntimeConfigurationV2Input) {
				if len(result.Testids) != 1 || result.Testids[0] != 42 {
					t.Fatalf("expected test id 42, got %v", result.Testids)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, resourceDowntimeConfigurationV2().Schema, tt.config)
			result := buildDowntimeConfigurationV2Data(d)
			tt.validate(t, result)
		})
	}
}

// TestFlattenDowntimeConfigurationV2Read tests flattenDowntimeConfigurationV2Read
func TestFlattenDowntimeConfigurationV2Read(t *testing.T) {
	tests := []struct {
		name     string
		input    *sc2.DowntimeConfigurationV2Response
		validate func(t *testing.T, result []interface{})
	}{
		{
			name: "populated downtime configuration response",
			input: &sc2.DowntimeConfigurationV2Response{
				DowntimeConfiguration: sc2.DowntimeConfiguration{
					ID:          123,
					Name:        "test downtime",
					Description: "test description",
					Rule:        "ALWAYS",
					Starttime:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					Endtime:     time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
					Timezone:    stringPtr("UTC"),
					Testids:     []int{1, 2, 3},
				},
			},
			validate: func(t *testing.T, result []interface{}) {
				if len(result) != 1 {
					t.Fatalf("expected 1 result, got %d", len(result))
				}
				m := result[0].(map[string]interface{})
				if m["name"] != "test downtime" {
					t.Fatalf("expected name 'test downtime', got %v", m["name"])
				}
				if m["description"] != "test description" {
					t.Fatalf("expected description 'test description', got %v", m["description"])
				}
				tzPtr, ok := m["timezone"].(*string)
				if !ok || tzPtr == nil || *tzPtr != "UTC" {
					t.Fatalf("expected timezone 'UTC', got %v", m["timezone"])
				}
			},
		},
		{
			name: "downtime configuration without timezone",
			input: &sc2.DowntimeConfigurationV2Response{
				DowntimeConfiguration: sc2.DowntimeConfiguration{
					ID:        456,
					Name:      "no timezone",
					Rule:      "ALWAYS",
					Starttime: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					Endtime:   time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
					Timezone:  nil,
					Testids:   []int{},
				},
			},
			validate: func(t *testing.T, result []interface{}) {
				if len(result) != 1 {
					t.Fatalf("expected 1 result, got %d", len(result))
				}
				m := result[0].(map[string]interface{})
				if m["name"] != "no timezone" {
					t.Fatalf("expected name, got %v", m["name"])
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := flattenDowntimeConfigurationV2Read(tt.input)
			tt.validate(t, result)
		})
	}
}

// TestFlattenDowntimeConfigurationV2Data tests flattenDowntimeConfigurationV2Data
func TestFlattenDowntimeConfigurationV2Data(t *testing.T) {
	tests := []struct {
		name     string
		input    *sc2.DowntimeConfigurationV2Response
		validate func(t *testing.T, result []interface{})
	}{
		{
			name: "full downtime configuration response",
			input: &sc2.DowntimeConfigurationV2Response{
				DowntimeConfiguration: sc2.DowntimeConfiguration{
					ID:             123,
					Name:           "test downtime",
					Description:    "test description",
					Rule:           "ALWAYS",
					Status:         "ACTIVE",
					Starttime:      time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					Endtime:        time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
					Createdat:      time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					Updatedat:      time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
					Testsupdatedat: time.Date(2024, 1, 1, 6, 0, 0, 0, time.UTC),
					Testcount:      3,
					Testids:        []int{1, 2, 3},
					Timezone:       stringPtr("UTC"),
				},
			},
			validate: func(t *testing.T, result []interface{}) {
				if len(result) != 1 {
					t.Fatalf("expected 1 result, got %d", len(result))
				}
				m := result[0].(map[string]interface{})
				if m["id"] != 123 {
					t.Fatalf("expected id 123, got %v", m["id"])
				}
				if m["status"] != "ACTIVE" {
					t.Fatalf("expected status 'ACTIVE', got %v", m["status"])
				}
				if m["test_count"] != 3 {
					t.Fatalf("expected test_count 3, got %v", m["test_count"])
				}
			},
		},
		{
			name: "downtime configuration with no metadata",
			input: &sc2.DowntimeConfigurationV2Response{
				DowntimeConfiguration: sc2.DowntimeConfiguration{
					ID:             789,
					Name:           "minimal",
					Description:    "",
					Rule:           "ALWAYS",
					Status:         "INACTIVE",
					Starttime:      time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					Endtime:        time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
					Createdat:      time.Time{},
					Updatedat:      time.Time{},
					Testsupdatedat: time.Time{},
					Testcount:      0,
					Testids:        []int{},
					Timezone:       nil,
				},
			},
			validate: func(t *testing.T, result []interface{}) {
				if len(result) != 1 {
					t.Fatalf("expected 1 result, got %d", len(result))
				}
				m := result[0].(map[string]interface{})
				if m["id"] != 789 {
					t.Fatalf("expected id 789, got %v", m["id"])
				}
				if m["status"] != "INACTIVE" {
					t.Fatalf("expected status 'INACTIVE', got %v", m["status"])
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := flattenDowntimeConfigurationV2Data(tt.input)
			tt.validate(t, result)
		})
	}
}

// TestFlattenDowntimeConfigurationsV2Data tests flattenDowntimeConfigurationsV2Data
func TestFlattenDowntimeConfigurationsV2Data(t *testing.T) {
	tests := []struct {
		name     string
		input    *[]sc2.DowntimeConfiguration
		validate func(t *testing.T, result []interface{})
	}{
		{
			name: "multiple downtime configurations",
			input: &[]sc2.DowntimeConfiguration{
				{
					ID:        1,
					Name:      "first",
					Rule:      "ALWAYS",
					Status:    "ACTIVE",
					Starttime: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					Endtime:   time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
					Createdat: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					Updatedat: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					Testcount: 1,
				},
				{
					ID:        2,
					Name:      "second",
					Rule:      "ALWAYS",
					Status:    "INACTIVE",
					Starttime: time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
					Endtime:   time.Date(2024, 2, 2, 0, 0, 0, 0, time.UTC),
					Createdat: time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
					Updatedat: time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
					Testcount: 2,
				},
			},
			validate: func(t *testing.T, result []interface{}) {
				if len(result) != 2 {
					t.Fatalf("expected 2 results, got %d", len(result))
				}
				first := result[0].(map[string]interface{})
				if first["id"] != 1 || first["name"] != "first" {
					t.Fatalf("first config invalid: %v", first)
				}
				second := result[1].(map[string]interface{})
				if second["id"] != 2 || second["name"] != "second" {
					t.Fatalf("second config invalid: %v", second)
				}
			},
		},
		{
			name:  "empty list",
			input: &[]sc2.DowntimeConfiguration{},
			validate: func(t *testing.T, result []interface{}) {
				if len(result) != 0 {
					t.Fatalf("expected empty result, got %d items", len(result))
				}
			},
		},
		{
			name:  "nil input",
			input: nil,
			validate: func(t *testing.T, result []interface{}) {
				if len(result) != 0 {
					t.Fatalf("expected empty result for nil input, got %d items", len(result))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := flattenDowntimeConfigurationsV2Data(tt.input)
			tt.validate(t, result)
		})
	}
}

// TestFlattenRecurrenceData tests flattenRecurrenceData
func TestFlattenRecurrenceData(t *testing.T) {
	tests := []struct {
		name     string
		input    *sc2.Recurrence
		validate func(t *testing.T, result []interface{})
	}{
		{
			name: "recurrence with repeats and end",
			input: &sc2.Recurrence{
				Repeats: sc2.Repeats{
					Type: "WEEKLY",
				},
				End: &sc2.End{
					Type:  "AFTER_OCCURRENCES",
					Value: "10",
				},
			},
			validate: func(t *testing.T, result []interface{}) {
				if len(result) != 1 {
					t.Fatalf("expected 1 result, got %d", len(result))
				}
				m := result[0].(map[string]interface{})
				if repeats, ok := m["repeats"].([]interface{}); !ok || len(repeats) == 0 {
					t.Fatalf("expected repeats in recurrence")
				}
				if end, ok := m["end"].([]interface{}); !ok || len(end) == 0 {
					t.Fatalf("expected end in recurrence")
				}
			},
		},
		{
			name: "recurrence without end",
			input: &sc2.Recurrence{
				Repeats: sc2.Repeats{
					Type: "DAILY",
				},
				End: nil,
			},
			validate: func(t *testing.T, result []interface{}) {
				if len(result) != 1 {
					t.Fatalf("expected 1 result, got %d", len(result))
				}
				m := result[0].(map[string]interface{})
				if repeats, ok := m["repeats"].([]interface{}); !ok || len(repeats) == 0 {
					t.Fatalf("expected repeats in recurrence")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := flattenRecurrenceData(tt.input)
			tt.validate(t, result)
		})
	}
}

// TestFlattenRepeatsData tests flattenRepeatsData
func TestFlattenRepeatsData(t *testing.T) {
	tests := []struct {
		name     string
		input    sc2.Repeats
		validate func(t *testing.T, result []interface{})
	}{
		{
			name: "repeats with type only",
			input: sc2.Repeats{
				Type:            "WEEKLY",
				Customvalue:     nil,
				Customfrequency: nil,
			},
			validate: func(t *testing.T, result []interface{}) {
				if len(result) != 1 {
					t.Fatalf("expected 1 result, got %d", len(result))
				}
				m := result[0].(map[string]interface{})
				if m["type"] != "WEEKLY" {
					t.Fatalf("expected type 'WEEKLY', got %v", m["type"])
				}
			},
		},
		{
			name: "repeats with custom values",
			input: sc2.Repeats{
				Type:            "CUSTOM",
				Customvalue:     intPtr(5),
				Customfrequency: stringPtr("DAILY"),
			},
			validate: func(t *testing.T, result []interface{}) {
				if len(result) != 1 {
					t.Fatalf("expected 1 result, got %d", len(result))
				}
				m := result[0].(map[string]interface{})
				if m["type"] != "CUSTOM" {
					t.Fatalf("expected type 'CUSTOM', got %v", m["type"])
				}
				if m["custom_value"] != 5 {
					t.Fatalf("expected custom_value 5, got %v", m["custom_value"])
				}
				if m["custom_frequency"] != "DAILY" {
					t.Fatalf("expected custom_frequency 'DAILY', got %v", m["custom_frequency"])
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := flattenRepeatsData(tt.input)
			tt.validate(t, result)
		})
	}
}

// TestFlattenEndData tests flattenEndData
func TestFlattenEndData(t *testing.T) {
	tests := []struct {
		name     string
		input    *sc2.End
		validate func(t *testing.T, result []interface{})
	}{
		{
			name: "end with after occurrences",
			input: &sc2.End{
				Type:  "AFTER_OCCURRENCES",
				Value: "10",
			},
			validate: func(t *testing.T, result []interface{}) {
				if len(result) != 1 {
					t.Fatalf("expected 1 result, got %d", len(result))
				}
				m := result[0].(map[string]interface{})
				if m["type"] != "AFTER_OCCURRENCES" {
					t.Fatalf("expected type 'AFTER_OCCURRENCES', got %v", m["type"])
				}
				if m["value"] != "10" {
					t.Fatalf("expected value '10', got %v", m["value"])
				}
			},
		},
		{
			name: "end with on date",
			input: &sc2.End{
				Type:  "ON_DATE",
				Value: "2024-12-31",
			},
			validate: func(t *testing.T, result []interface{}) {
				if len(result) != 1 {
					t.Fatalf("expected 1 result, got %d", len(result))
				}
				m := result[0].(map[string]interface{})
				if m["type"] != "ON_DATE" {
					t.Fatalf("expected type 'ON_DATE', got %v", m["type"])
				}
				if m["value"] != "2024-12-31" {
					t.Fatalf("expected value '2024-12-31', got %v", m["value"])
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := flattenEndData(tt.input)
			tt.validate(t, result)
		})
	}
}

// TestBuildTestIdData tests buildTestIdData
func TestBuildTestIdData(t *testing.T) {
	tests := []struct {
		name     string
		input    []interface{}
		validate func(t *testing.T, result []int)
	}{
		{
			name:  "multiple test ids",
			input: []interface{}{1, 2, 3, 4, 5},
			validate: func(t *testing.T, result []int) {
				if len(result) != 5 {
					t.Fatalf("expected 5 test ids, got %d", len(result))
				}
				for i, id := range result {
					if id != i+1 {
						t.Fatalf("expected test id %d at index %d, got %d", i+1, i, id)
					}
				}
			},
		},
		{
			name:  "single test id",
			input: []interface{}{42},
			validate: func(t *testing.T, result []int) {
				if len(result) != 1 {
					t.Fatalf("expected 1 test id, got %d", len(result))
				}
				if result[0] != 42 {
					t.Fatalf("expected test id 42, got %d", result[0])
				}
			},
		},
		{
			name:  "empty test ids",
			input: []interface{}{},
			validate: func(t *testing.T, result []int) {
				if len(result) != 0 {
					t.Fatalf("expected 0 test ids, got %d", len(result))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildTestIdData(tt.input)
			tt.validate(t, result)
		})
	}
}

// TestFlattenLocationsV2Data tests flattenLocationsV2Data
func TestFlattenLocationsV2Data(t *testing.T) {
	tests := []struct {
		name     string
		input    *[]sc2.Location
		validate func(t *testing.T, result []interface{})
	}{
		{
			name: "multiple locations",
			input: &[]sc2.Location{
				{
					ID:      "us-east-1",
					Label:   "US East",
					Default: true,
					Type:    "PUBLIC",
					Country: "USA",
				},
				{
					ID:      "eu-west-1",
					Label:   "EU West",
					Default: false,
					Type:    "PUBLIC",
					Country: "Ireland",
				},
			},
			validate: func(t *testing.T, result []interface{}) {
				if len(result) != 2 {
					t.Fatalf("expected 2 locations, got %d", len(result))
				}
				first := result[0].(map[string]interface{})
				if first["id"] != "us-east-1" || first["default"] != true {
					t.Fatalf("first location invalid: %v", first)
				}
			},
		},
		{
			name:  "empty locations",
			input: &[]sc2.Location{},
			validate: func(t *testing.T, result []interface{}) {
				if len(result) != 0 {
					t.Fatalf("expected 0 locations, got %d", len(result))
				}
			},
		},
		{
			name:  "nil locations",
			input: nil,
			validate: func(t *testing.T, result []interface{}) {
				if len(result) != 0 {
					t.Fatalf("expected 0 locations for nil input, got %d", len(result))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := flattenLocationsV2Data(tt.input)
			tt.validate(t, result)
		})
	}
}

// TestFlattenDefaultLocationData tests flattenDefaultLocationData
func TestFlattenDefaultLocationData(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		validate func(t *testing.T, result []interface{})
	}{
		{
			name:  "multiple default locations",
			input: []string{"us-east-1", "eu-west-1", "ap-southeast-1"},
			validate: func(t *testing.T, result []interface{}) {
				if len(result) != 3 {
					t.Fatalf("expected 3 locations, got %d", len(result))
				}
				if result[0].(string) != "us-east-1" {
					t.Fatalf("expected first location 'us-east-1', got %v", result[0])
				}
			},
		},
		{
			name:  "single default location",
			input: []string{"us-east-1"},
			validate: func(t *testing.T, result []interface{}) {
				if len(result) != 1 {
					t.Fatalf("expected 1 location, got %d", len(result))
				}
			},
		},
		{
			name:  "empty locations",
			input: []string{},
			validate: func(t *testing.T, result []interface{}) {
				if len(result) != 0 {
					t.Fatalf("expected 0 locations, got %d", len(result))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := flattenDefaultLocationData(tt.input)
			tt.validate(t, result)
		})
	}
}

// TestBuildLocationV2Data tests buildLocationV2Data
func TestBuildLocationV2Data(t *testing.T) {
	tests := []struct {
		name     string
		config   map[string]interface{}
		validate func(t *testing.T, result sc2.LocationV2Input)
	}{
		{
			name: "public location",
			config: map[string]interface{}{
				"location": []interface{}{
					map[string]interface{}{
						"id":      "us-east-1",
						"label":   "US East",
						"default": true,
						"type":    "PUBLIC",
						"country": "USA",
					},
				},
			},
			validate: func(t *testing.T, result sc2.LocationV2Input) {
				if result.ID != "us-east-1" {
					t.Fatalf("expected id 'us-east-1', got %q", result.ID)
				}
				if result.Label != "US East" {
					t.Fatalf("expected label 'US East', got %q", result.Label)
				}
				if result.Default != true {
					t.Fatalf("expected default true, got %v", result.Default)
				}
			},
		},
		{
			name: "private location",
			config: map[string]interface{}{
				"location": []interface{}{
					map[string]interface{}{
						"id":      "private-custom-1",
						"label":   "Private Loc",
						"default": false,
						"type":    "PRIVATE",
						"country": "",
					},
				},
			},
			validate: func(t *testing.T, result sc2.LocationV2Input) {
				if result.ID != "private-custom-1" {
					t.Fatalf("expected id 'private-custom-1', got %q", result.ID)
				}
				if result.Type != "PRIVATE" {
					t.Fatalf("expected type 'PRIVATE', got %q", result.Type)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, resourceLocationV2().Schema, tt.config)
			result := buildLocationV2Data(d)
			tt.validate(t, result)
		})
	}
}

// TestFlattenLocationV2Data tests flattenLocationV2Data
func TestFlattenLocationV2Data(t *testing.T) {
	tests := []struct {
		name     string
		input    sc2.Location
		validate func(t *testing.T, result []interface{})
	}{
		{
			name: "public location",
			input: sc2.Location{
				ID:      "us-east-1",
				Label:   "US East",
				Default: true,
				Type:    "PUBLIC",
				Country: "USA",
			},
			validate: func(t *testing.T, result []interface{}) {
				if len(result) != 1 {
					t.Fatalf("expected 1 result, got %d", len(result))
				}
				m := result[0].(map[string]interface{})
				if m["id"] != "us-east-1" {
					t.Fatalf("expected id 'us-east-1', got %v", m["id"])
				}
				if m["default"] != true {
					t.Fatalf("expected default true, got %v", m["default"])
				}
			},
		},
		{
			name: "private location",
			input: sc2.Location{
				ID:      "private-custom-1",
				Label:   "Private Loc",
				Default: false,
				Type:    "PRIVATE",
				Country: "",
			},
			validate: func(t *testing.T, result []interface{}) {
				if len(result) != 1 {
					t.Fatalf("expected 1 result, got %d", len(result))
				}
				m := result[0].(map[string]interface{})
				if m["id"] != "private-custom-1" {
					t.Fatalf("expected id 'private-custom-1', got %v", m["id"])
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := flattenLocationV2Data(tt.input)
			tt.validate(t, result)
		})
	}
}

// TestFlattenLocationMetaV2Data tests flattenLocationMetaV2Data
func TestFlattenLocationMetaV2Data(t *testing.T) {
	tests := []struct {
		name     string
		input    sc2.Meta
		validate func(t *testing.T, result []interface{})
	}{
		{
			name: "meta with active and paused test ids",
			input: sc2.Meta{
				ActiveTestIds: []int{1, 2, 3},
				PausedTestIds: []int{4, 5},
			},
			validate: func(t *testing.T, result []interface{}) {
				if len(result) != 1 {
					t.Fatalf("expected 1 result, got %d", len(result))
				}
				m := result[0].(map[string]interface{})
				activeTests := m["active_test_ids"].([]int)
				if len(activeTests) != 3 {
					t.Fatalf("expected 3 active tests, got %d", len(activeTests))
				}
				pausedTests := m["paused_test_ids"].([]int)
				if len(pausedTests) != 2 {
					t.Fatalf("expected 2 paused tests, got %d", len(pausedTests))
				}
			},
		},
		{
			name: "meta with no active tests",
			input: sc2.Meta{
				ActiveTestIds: []int{},
				PausedTestIds: []int{1, 2},
			},
			validate: func(t *testing.T, result []interface{}) {
				if len(result) != 1 {
					t.Fatalf("expected 1 result, got %d", len(result))
				}
				m := result[0].(map[string]interface{})
				activeTests := m["active_test_ids"].([]int)
				if len(activeTests) != 0 {
					t.Fatalf("expected 0 active tests, got %d", len(activeTests))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := flattenLocationMetaV2Data(tt.input)
			tt.validate(t, result)
		})
	}
}

// TestBuildRepeatsData tests buildRepeatsData by directly testing its behavior
func TestBuildRepeatsData(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]interface{}
		validate func(t *testing.T, result *sc2.Repeats)
	}{
		{
			name: "repeats with type only",
			input: map[string]interface{}{
				"type":             "WEEKLY",
				"custom_value":     nil,
				"custom_frequency": nil,
			},
			validate: func(t *testing.T, result *sc2.Repeats) {
				if result == nil {
					t.Fatal("expected repeats, got nil")
				}
				if result.Type != "WEEKLY" {
					t.Fatalf("expected type 'WEEKLY', got %q", result.Type)
				}
				if result.Customvalue != nil {
					t.Fatalf("expected nil customvalue, got %v", result.Customvalue)
				}
			},
		},
		{
			name: "repeats with custom values",
			input: map[string]interface{}{
				"type":             "CUSTOM",
				"custom_value":     5,
				"custom_frequency": "DAILY",
			},
			validate: func(t *testing.T, result *sc2.Repeats) {
				if result == nil {
					t.Fatal("expected repeats, got nil")
				}
				if result.Type != "CUSTOM" {
					t.Fatalf("expected type 'CUSTOM', got %q", result.Type)
				}
				if result.Customvalue == nil || *result.Customvalue != 5 {
					t.Fatalf("expected customvalue 5, got %v", result.Customvalue)
				}
				if result.Customfrequency == nil || *result.Customfrequency != "DAILY" {
					t.Fatalf("expected customfrequency 'DAILY', got %v", result.Customfrequency)
				}
			},
		},
		{
			name:  "empty repeats",
			input: map[string]interface{}{},
			validate: func(t *testing.T, result *sc2.Repeats) {
				if result != nil {
					t.Fatalf("expected nil repeats for empty input, got %v", result)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var s *schema.Set
			if len(tt.input) > 0 {
				hash := schema.HashResource(&schema.Resource{
					Schema: map[string]*schema.Schema{
						"type":             {Type: schema.TypeString},
						"custom_value":     {Type: schema.TypeInt},
						"custom_frequency": {Type: schema.TypeString},
					},
				})
				s = schema.NewSet(hash, []interface{}{tt.input})
			} else {
				hash := schema.HashResource(&schema.Resource{
					Schema: map[string]*schema.Schema{
						"type":             {Type: schema.TypeString},
						"custom_value":     {Type: schema.TypeInt},
						"custom_frequency": {Type: schema.TypeString},
					},
				})
				s = schema.NewSet(hash, []interface{}{})
			}
			result := buildRepeatsData(s)
			tt.validate(t, result)
		})
	}
}

// TestBuildEndData tests buildEndData by directly testing its behavior
func TestBuildEndData(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]interface{}
		validate func(t *testing.T, result *sc2.End)
	}{
		{
			name: "end with after occurrences",
			input: map[string]interface{}{
				"type":  "AFTER_OCCURRENCES",
				"value": "10",
			},
			validate: func(t *testing.T, result *sc2.End) {
				if result == nil {
					t.Fatal("expected end, got nil")
				}
				if result.Type != "AFTER_OCCURRENCES" {
					t.Fatalf("expected type 'AFTER_OCCURRENCES', got %q", result.Type)
				}
				if result.Value != "10" {
					t.Fatalf("expected value '10', got %q", result.Value)
				}
			},
		},
		{
			name: "end with on date",
			input: map[string]interface{}{
				"type":  "ON_DATE",
				"value": "2024-12-31",
			},
			validate: func(t *testing.T, result *sc2.End) {
				if result == nil {
					t.Fatal("expected end, got nil")
				}
				if result.Type != "ON_DATE" {
					t.Fatalf("expected type 'ON_DATE', got %q", result.Type)
				}
				if result.Value != "2024-12-31" {
					t.Fatalf("expected value '2024-12-31', got %q", result.Value)
				}
			},
		},
		{
			name:  "empty end",
			input: map[string]interface{}{},
			validate: func(t *testing.T, result *sc2.End) {
				if result != nil {
					t.Fatalf("expected nil end for empty input, got %v", result)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var s *schema.Set
			if len(tt.input) > 0 {
				hash := schema.HashResource(&schema.Resource{
					Schema: map[string]*schema.Schema{
						"type":  {Type: schema.TypeString},
						"value": {Type: schema.TypeString},
					},
				})
				s = schema.NewSet(hash, []interface{}{tt.input})
			} else {
				hash := schema.HashResource(&schema.Resource{
					Schema: map[string]*schema.Schema{
						"type":  {Type: schema.TypeString},
						"value": {Type: schema.TypeString},
					},
				})
				s = schema.NewSet(hash, []interface{}{})
			}
			result := buildEndData(s)
			tt.validate(t, result)
		})
	}
}

// TestBuildRecurrenceData tests buildRecurrenceData through resource data integration
func TestBuildRecurrenceData(t *testing.T) {
	tests := []struct {
		name     string
		config   map[string]interface{}
		validate func(t *testing.T, result *sc2.Recurrence)
	}{
		{
			name: "populated recurrence in downtime configuration",
			config: map[string]interface{}{
				"downtime_configuration": []interface{}{
					map[string]interface{}{
						"name":        "with recurrence",
						"description": "test",
						"rule":        "ALWAYS",
						"start_time":  "2024-01-01T00:00:00.000Z",
						"end_time":    "2024-01-02T00:00:00.000Z",
						"test_ids":    []interface{}{},
						"recurrence": []interface{}{
							map[string]interface{}{
								"repeats": []interface{}{
									map[string]interface{}{
										"type":             "WEEKLY",
										"custom_value":     0,
										"custom_frequency": "",
									},
								},
								"end": []interface{}{
									map[string]interface{}{
										"type":  "AFTER_OCCURRENCES",
										"value": "10",
									},
								},
							},
						},
					},
				},
			},
			validate: func(t *testing.T, result *sc2.Recurrence) {
				if result == nil {
					t.Fatal("expected recurrence, got nil")
				}
				if result.Repeats.Type != "WEEKLY" {
					t.Fatalf("expected repeats type 'WEEKLY', got %q", result.Repeats.Type)
				}
				if result.End == nil {
					t.Fatal("expected end, got nil")
				}
				if result.End.Type != "AFTER_OCCURRENCES" {
					t.Fatalf("expected end type 'AFTER_OCCURRENCES', got %q", result.End.Type)
				}
			},
		},
		{
			name: "empty recurrence",
			config: map[string]interface{}{
				"downtime_configuration": []interface{}{
					map[string]interface{}{
						"name":        "no recurrence",
						"description": "",
						"rule":        "ALWAYS",
						"start_time":  "2024-01-01T00:00:00.000Z",
						"end_time":    "2024-01-02T00:00:00.000Z",
						"test_ids":    []interface{}{},
						"recurrence":  []interface{}{},
					},
				},
			},
			validate: func(t *testing.T, result *sc2.Recurrence) {
				if result == nil {
					t.Fatal("expected recurrence, got nil")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, resourceDowntimeConfigurationV2().Schema, tt.config)
			downtimeConfigData := d.Get("downtime_configuration").(*schema.Set).List()
			if len(downtimeConfigData) > 0 {
				downtimeConfiguration := downtimeConfigData[0].(map[string]interface{})
				result := buildRecurrenceData(downtimeConfiguration["recurrence"].(*schema.Set))
				tt.validate(t, result)
			} else {
				t.Fatal("expected downtime_configuration data")
			}
		})
	}
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}
