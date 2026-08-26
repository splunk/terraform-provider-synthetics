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

func TestTotpVariableIDFromList(t *testing.T) {
	tests := []struct {
		name  string
		value interface{}
		want  int
	}{
		{"not a slice", "nope", 0},
		{"nil value", nil, 0},
		{"empty slice", []interface{}{}, 0},
		{"nil first element", []interface{}{nil}, 0},
		{"first element not a map", []interface{}{"nope"}, 0},
		{"missing id", []interface{}{map[string]interface{}{}}, 0},
		{"populated", []interface{}{map[string]interface{}{"id": 42}}, 42},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := totpVariableIDFromList(tc.value); got != tc.want {
				t.Fatalf("totpVariableIDFromList(%#v) = %d, want %d", tc.value, got, tc.want)
			}
		})
	}
}

func TestTotpVariableSecretFromState(t *testing.T) {
	t.Run("no totp_variable block", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, resourceTotpVariableV2().Schema, map[string]interface{}{})
		if got := totpVariableSecretFromState(d); got != "" {
			t.Fatalf("totpVariableSecretFromState() = %q, want empty", got)
		}
	})

	t.Run("populated block", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, resourceTotpVariableV2().Schema, map[string]interface{}{
			"totp_variable": []interface{}{
				map[string]interface{}{
					"name":   "login_mfa",
					"secret": "JBSWY3DPEHPK3PXP",
				},
			},
		})
		if got := totpVariableSecretFromState(d); got != "JBSWY3DPEHPK3PXP" {
			t.Fatalf("totpVariableSecretFromState() = %q, want JBSWY3DPEHPK3PXP", got)
		}
	})
}

func TestFlattenTotpVariableV2Data(t *testing.T) {
	t.Run("nil response", func(t *testing.T) {
		got := flattenTotpVariableV2Data(nil)
		if len(got) != 0 {
			t.Fatalf("len(flattenTotpVariableV2Data(nil)) = %d, want 0", len(got))
		}
	})

	t.Run("populated response", func(t *testing.T) {
		resp := &sc2.TotpVariableV2Response{}
		resp.Totp.ID = 7
		resp.Totp.Name = "login_mfa"

		got := flattenTotpVariableV2Data(resp)
		if len(got) != 1 {
			t.Fatalf("len(flattenTotpVariableV2Data()) = %d, want 1", len(got))
		}
		m := got[0].(map[string]interface{})
		if m["id"] != 7 || m["name"] != "login_mfa" {
			t.Fatalf("flattenTotpVariableV2Data() = %#v", m)
		}
	})
}

func TestFlattenTotpVariablesV2Data(t *testing.T) {
	t.Run("empty slice", func(t *testing.T) {
		got := flattenTotpVariablesV2Data([]sc2.TotpVariable{})
		if len(got) != 0 {
			t.Fatalf("len(flattenTotpVariablesV2Data()) = %d, want 0", len(got))
		}
	})

	t.Run("multiple elements", func(t *testing.T) {
		got := flattenTotpVariablesV2Data([]sc2.TotpVariable{
			{ID: 1, Name: "one"},
			{ID: 2, Name: "two"},
		})
		if len(got) != 2 {
			t.Fatalf("len(flattenTotpVariablesV2Data()) = %d, want 2", len(got))
		}
		if got[0].(map[string]interface{})["name"] != "one" || got[1].(map[string]interface{})["name"] != "two" {
			t.Fatalf("flattenTotpVariablesV2Data() = %#v", got)
		}
	})
}

func TestBuildTotpVariableV2DataRejectsRedactedSecret(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceTotpVariableV2().Schema, map[string]interface{}{
		"totp_variable": []interface{}{
			map[string]interface{}{
				"name":   "login_mfa",
				"secret": totpVariableRedactedSecret,
			},
		},
	})

	_, err := buildTotpVariableV2Data(d)
	if err == nil {
		t.Fatal("expected error when TOTP secret is the redacted placeholder")
	}
}

func TestBuildTotpVariableV2DataWithoutTotpVariableBlock(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceTotpVariableV2().Schema, map[string]interface{}{})

	_, err := buildTotpVariableV2Data(d)
	if err == nil {
		t.Fatal("expected error when totp_variable block is entirely absent")
	}
}

func TestBuildTotpVariableV2UpdateDataWithoutTotpVariableBlock(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceTotpVariableV2().Schema, map[string]interface{}{})

	got := buildTotpVariableV2UpdateData(d)
	if got.Totp.Description != nil || got.Totp.Secret != nil {
		t.Fatalf("buildTotpVariableV2UpdateData() = %#v, want zero value when no totp_variable block is set", got.Totp)
	}
}

func TestFlattenTotpVariableV2ReadNilResponse(t *testing.T) {
	got := flattenTotpVariableV2Read(nil, "existing-secret")
	if len(got) != 0 {
		t.Fatalf("len(flattenTotpVariableV2Read(nil)) = %d, want 0", len(got))
	}
}

func TestFlattenTotpVariableV2MetadataFullyPopulated(t *testing.T) {
	createdAt := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2021, 2, 2, 0, 0, 0, 0, time.UTC)
	totp := sc2.TotpVariable{
		ID:          9,
		Name:        "login_mfa",
		Description: "login MFA",
		Digits:      6,
		Interval:    30,
		HmacDigest:  "sha1",
		CreatedAt:   createdAt,
		CreatedBy:   "alice",
		UpdatedAt:   updatedAt,
		UpdatedBy:   "bob",
	}

	got := flattenTotpVariableV2Metadata(totp)
	if got["created_at"] != createdAt.Format(time.RFC3339) || got["created_by"] != "alice" {
		t.Fatalf("created fields = %#v", got)
	}
	if got["updated_at"] != updatedAt.Format(time.RFC3339) || got["updated_by"] != "bob" {
		t.Fatalf("updated fields = %#v", got)
	}
}

func TestFlattenTotpVariableV2MetadataZeroValue(t *testing.T) {
	got := flattenTotpVariableV2Metadata(sc2.TotpVariable{})
	if len(got) != 0 {
		t.Fatalf("flattenTotpVariableV2Metadata(zero value) = %#v, want empty map", got)
	}
}

func TestTotpVariableSecretForState(t *testing.T) {
	tests := []struct {
		name           string
		apiSecret      string
		existingSecret string
		want           string
	}{
		{"redacted with existing secret falls back to state", totpVariableRedactedSecret, "existing", "existing"},
		{"empty api secret with existing secret falls back to state", "", "existing", "existing"},
		{"redacted with no existing secret returns empty", totpVariableRedactedSecret, "", ""},
		{"real api secret passes through", "JBSWY3DPEHPK3PXP", "", "JBSWY3DPEHPK3PXP"},
		{"real api secret overrides stale existing", "JBSWY3DPEHPK3PXP", "stale", "JBSWY3DPEHPK3PXP"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := totpVariableSecretForState(tc.apiSecret, tc.existingSecret); got != tc.want {
				t.Fatalf("totpVariableSecretForState(%q, %q) = %q, want %q", tc.apiSecret, tc.existingSecret, got, tc.want)
			}
		})
	}
}

func TestTotpStringField(t *testing.T) {
	tests := []struct {
		name string
		totp map[string]interface{}
		key  string
		want string
	}{
		{"present string", map[string]interface{}{"name": "login_mfa"}, "name", "login_mfa"},
		{"missing key", map[string]interface{}{}, "name", ""},
		{"wrong type", map[string]interface{}{"name": 42}, "name", ""},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := totpStringField(tc.totp, tc.key); got != tc.want {
				t.Fatalf("totpStringField() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestTotpIntField(t *testing.T) {
	tests := []struct {
		name string
		totp map[string]interface{}
		key  string
		want int
	}{
		{"present int", map[string]interface{}{"digits": 6}, "digits", 6},
		{"missing key", map[string]interface{}{}, "digits", 0},
		{"wrong type", map[string]interface{}{"digits": "six"}, "digits", 0},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := totpIntField(tc.totp, tc.key); got != tc.want {
				t.Fatalf("totpIntField() = %d, want %d", got, tc.want)
			}
		})
	}
}
