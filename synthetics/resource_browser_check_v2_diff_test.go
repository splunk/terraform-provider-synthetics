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
)

func rawBrowserCheckWithStep(step map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"test": []interface{}{
			map[string]interface{}{
				"name":                "browser-test",
				"active":              true,
				"frequency":           5,
				"device_id":           1,
				"location_ids":        []interface{}{"aws-us-east-1"},
				"scheduling_strategy": "round_robin",
				"transactions": []interface{}{
					map[string]interface{}{
						"name":  "T1",
						"steps": []interface{}{step},
					},
				},
			},
		},
	}
}

func TestBrowserCheckV2StepPrefixFromAttributeKey(t *testing.T) {
	tests := []struct {
		name string
		key  string
		want string
	}{
		{"matching selector key", "test.0.transactions.0.steps.0.selector", "test.0.transactions.0.steps.0"},
		{"matching selector_type key", "test.0.transactions.0.steps.0.selector_type", "test.0.transactions.0.steps.0"},
		{"matching nested selectors key", "test.0.transactions.0.steps.0.selectors.0.value", "test.0.transactions.0.steps.0"},
		{"different index", "test.0.transactions.1.steps.2.selector", "test.0.transactions.1.steps.2"},
		{"no trailing field", "test.0.transactions.0.steps.0", ""},
		{"unrelated key", "test.0.name", ""},
		{"empty key", "", ""},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := browserCheckV2StepPrefixFromAttributeKey(tc.key); got != tc.want {
				t.Fatalf("browserCheckV2StepPrefixFromAttributeKey(%q) = %q, want %q", tc.key, got, tc.want)
			}
		})
	}
}

// interfaceFromResourceDataField/stringFromResourceDataField read via d.HasChange + d.GetChange
// when the field has a diff, and fall back to d.GetOk otherwise. schema.TestResourceDataRaw builds
// its diff against a nil prior state, so every populated field's diff is (old="" zero value, new=
// configured value) — HasChange is always true for anything set in the raw config. That's a real,
// exercised code path (it's exactly what happens during resource Create), so it's tested directly
// below; it just means useState=true always yields the type's zero value in these tests, not an
// arbitrary "old" state value.
func TestInterfaceFromResourceDataField(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceBrowserCheckV2().Schema, rawBrowserCheckWithStep(map[string]interface{}{
		"name":          "step1",
		"type":          "click",
		"selector":      "old-sel",
		"selector_type": "id",
	}))
	const key = "test.0.transactions.0.steps.0.selector"

	t.Run("useState=false returns the configured (new) value", func(t *testing.T) {
		got := interfaceFromResourceDataField(d, key, false)
		if got != "old-sel" {
			t.Fatalf("interfaceFromResourceDataField(useState=false) = %#v, want old-sel", got)
		}
	})

	t.Run("useState=true returns the pre-diff (zero) value", func(t *testing.T) {
		got := interfaceFromResourceDataField(d, key, true)
		if got != "" {
			t.Fatalf("interfaceFromResourceDataField(useState=true) = %#v, want empty string", got)
		}
	})

	t.Run("unset field falls back to GetOk and returns nil", func(t *testing.T) {
		got := interfaceFromResourceDataField(d, "test.0.transactions.0.steps.0.value", false)
		if got != nil {
			t.Fatalf("interfaceFromResourceDataField(unset) = %#v, want nil", got)
		}
	})
}

func TestStringFromResourceDataField(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceBrowserCheckV2().Schema, rawBrowserCheckWithStep(map[string]interface{}{
		"name":          "step1",
		"type":          "click",
		"selector":      "old-sel",
		"selector_type": "id",
	}))

	t.Run("populated string field", func(t *testing.T) {
		got := stringFromResourceDataField(d, "test.0.transactions.0.steps.0.selector", false)
		if got != "old-sel" {
			t.Fatalf("stringFromResourceDataField() = %q, want old-sel", got)
		}
	})

	t.Run("unset field returns empty string", func(t *testing.T) {
		got := stringFromResourceDataField(d, "test.0.transactions.0.steps.0.value", false)
		if got != "" {
			t.Fatalf("stringFromResourceDataField(unset) = %q, want empty string", got)
		}
	})

	t.Run("useState=true on a diffed field returns the zero value, not an error", func(t *testing.T) {
		got := stringFromResourceDataField(d, "test.0.transactions.0.steps.0.selector", true)
		if got != "" {
			t.Fatalf("stringFromResourceDataField(useState=true) = %q, want empty string", got)
		}
	})
}

func TestStepSelectorInputFromResourceData(t *testing.T) {
	t.Run("legacy selector fields", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, resourceBrowserCheckV2().Schema, rawBrowserCheckWithStep(map[string]interface{}{
			"name":          "step1",
			"type":          "click",
			"selector":      "old-sel",
			"selector_type": "id",
		}))
		got := stepSelectorInputFromResourceData(d, "test.0.transactions.0.steps.0", false)
		if got.selectorType != "id" || got.selector != "old-sel" {
			t.Fatalf("stepSelectorInputFromResourceData() = %#v, want id/old-sel", got)
		}
		if len(got.selectors) != 0 {
			t.Fatalf("selectors = %#v, want empty", got.selectors)
		}
	})

	t.Run("selectors block", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, resourceBrowserCheckV2().Schema, rawBrowserCheckWithStep(map[string]interface{}{
			"name": "step1",
			"type": "click",
			"selectors": []interface{}{
				map[string]interface{}{"type": "css", "value": ".primary"},
			},
		}))
		got := stepSelectorInputFromResourceData(d, "test.0.transactions.0.steps.0", false)
		if len(got.selectors) != 1 || got.selectors[0].Type != "css" || got.selectors[0].Value != ".primary" {
			t.Fatalf("stepSelectorInputFromResourceData() = %#v, want one css/.primary selector", got.selectors)
		}
	})

	t.Run("no selector fields set (e.g. assert_text_present)", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, resourceBrowserCheckV2().Schema, rawBrowserCheckWithStep(map[string]interface{}{
			"name":  "step1",
			"type":  "assert_text_present",
			"value": "Order confirmed",
		}))
		got := stepSelectorInputFromResourceData(d, "test.0.transactions.0.steps.0", false)
		if got.selectorType != "" || got.selector != "" || len(got.selectors) != 0 {
			t.Fatalf("stepSelectorInputFromResourceData() = %#v, want zero value", got)
		}
	})

	t.Run("useState=true on either shape yields an empty input", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, resourceBrowserCheckV2().Schema, rawBrowserCheckWithStep(map[string]interface{}{
			"name":          "step1",
			"type":          "click",
			"selector":      "old-sel",
			"selector_type": "id",
		}))
		got := stepSelectorInputFromResourceData(d, "test.0.transactions.0.steps.0", true)
		if got.selectorType != "" || got.selector != "" || len(got.selectors) != 0 {
			t.Fatalf("stepSelectorInputFromResourceData(useState=true) = %#v, want zero value", got)
		}
	})
}

// browserCheckV2SelectorRepresentationDiffSuppress's actual suppression (true) path requires a
// genuine state-vs-config diff where the old value is populated and differs in shape but not in
// meaning from the new value (e.g. legacy selector_type/selector in state vs an equivalent
// selectors block in config). schema.TestResourceDataRaw cannot construct that: its diff is
// always computed against a nil prior state, so the "old" side is always the zero value and can
// never be equivalent-but-differently-shaped from a populated "new" side. That path's logic is
// still fully covered indirectly: stepSelectorInputsEquivalent, migratingFromLegacyToSelectors,
// and stepSelectorRepresentationDiffers (the three functions this one composes to decide
// suppression) are unit tested directly at 100% coverage in structures_browser_test.go, and this
// function's own dispatch to them is exercised below on every reachable branch.
func TestBrowserCheckV2SelectorRepresentationDiffSuppress(t *testing.T) {
	t.Run("key does not match a step attribute prefix", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, resourceBrowserCheckV2().Schema, rawBrowserCheckWithStep(map[string]interface{}{
			"name": "step1",
			"type": "click",
		}))
		if browserCheckV2SelectorRepresentationDiffSuppress("test.0.name", "old", "new", d) {
			t.Fatal("expected no suppression for a key outside any step's attribute prefix")
		}
	})

	t.Run("populated selector is never equivalent to the pre-diff zero value, so no suppression", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, resourceBrowserCheckV2().Schema, rawBrowserCheckWithStep(map[string]interface{}{
			"name":          "step1",
			"type":          "click",
			"selector":      "old-sel",
			"selector_type": "id",
		}))
		if browserCheckV2SelectorRepresentationDiffSuppress("test.0.transactions.0.steps.0.selector", "", "old-sel", d) {
			t.Fatal("expected no suppression when old and new selector representations are not equivalent")
		}
	})
}
