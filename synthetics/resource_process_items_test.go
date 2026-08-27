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
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Each processXItems function is a thin pass-through to an already-tested buildXData
// function (verified by reading the source before writing these tests). These tests
// confirm the pass-through itself: same ResourceData in, structurally identical output.

func TestProcessApiCheckV2Items(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceApiCheckV2().Schema, map[string]interface{}{
		"test": []interface{}{map[string]interface{}{"name": "api-test"}},
	})
	if got, want := processApiCheckV2Items(d), buildApiV2Data(d); !reflect.DeepEqual(got, want) {
		t.Fatalf("processApiCheckV2Items() = %#v, want %#v", got, want)
	}
}

func TestProcessVariableV2Items(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceVariableV2().Schema, map[string]interface{}{
		"variable": []interface{}{map[string]interface{}{"name": "var-test", "value": "v"}},
	})
	if got, want := processVariableV2Items(d), buildVariableV2Data(d); !reflect.DeepEqual(got, want) {
		t.Fatalf("processVariableV2Items() = %#v, want %#v", got, want)
	}
}

func TestProcessPortCheckV2Items(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourcePortCheckV2().Schema, map[string]interface{}{
		"test": []interface{}{map[string]interface{}{"name": "port-test"}},
	})
	if got, want := processPortCheckV2Items(d), buildPortCheckV2Data(d); !reflect.DeepEqual(got, want) {
		t.Fatalf("processPortCheckV2Items() = %#v, want %#v", got, want)
	}
}

func TestProcessHttpCheckV2Items(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceHttpCheckV2().Schema, map[string]interface{}{
		"test": []interface{}{map[string]interface{}{"name": "http-test"}},
	})
	if got, want := processHttpCheckV2Items(d), buildHttpV2Data(d); !reflect.DeepEqual(got, want) {
		t.Fatalf("processHttpCheckV2Items() = %#v, want %#v", got, want)
	}
}

func TestProcessDowntimeConfigurationV2Items(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceDowntimeConfigurationV2().Schema, map[string]interface{}{
		"downtime_configuration": []interface{}{map[string]interface{}{"name": "downtime-test"}},
	})
	if got, want := processDowntimeConfigurationV2Items(d), buildDowntimeConfigurationV2Data(d); !reflect.DeepEqual(got, want) {
		t.Fatalf("processDowntimeConfigurationV2Items() = %#v, want %#v", got, want)
	}
}

func TestProcessLocationV2Items(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceLocationV2().Schema, map[string]interface{}{
		"location": []interface{}{map[string]interface{}{"label": "location-test"}},
	})
	if got, want := processLocationV2Items(d), buildLocationV2Data(d); !reflect.DeepEqual(got, want) {
		t.Fatalf("processLocationV2Items() = %#v, want %#v", got, want)
	}
}

func TestProcessBrowserCheckV2Items(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceBrowserCheckV2().Schema, map[string]interface{}{
		"test": []interface{}{map[string]interface{}{"name": "browser-test"}},
	})
	got, gotErr := processBrowserCheckV2Items(d)
	want, wantErr := buildBrowserV2Data(d)
	if (gotErr == nil) != (wantErr == nil) {
		t.Fatalf("processBrowserCheckV2Items() error = %v, buildBrowserV2Data() error = %v", gotErr, wantErr)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("processBrowserCheckV2Items() = %#v, want %#v", got, want)
	}
}

func TestClientCertificateKeySchema(t *testing.T) {
	t.Run("required", func(t *testing.T) {
		s := clientCertificateKeySchema(true)
		if !s.Required || s.Computed {
			t.Fatalf("clientCertificateKeySchema(true) = %#v, want Required=true, Computed=false", s)
		}
	})

	t.Run("computed", func(t *testing.T) {
		s := clientCertificateKeySchema(false)
		if s.Required || !s.Computed {
			t.Fatalf("clientCertificateKeySchema(false) = %#v, want Required=false, Computed=true", s)
		}
	})
}
