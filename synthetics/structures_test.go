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

	sc2 "github.com/splunk/syntheticsclient/v2/syntheticsclientv2"
)

func TestBuildSelectorsFromStep_multipleSelectors(t *testing.T) {
	step := map[string]interface{}{
		"name": "Click submit",
		"selectors": []interface{}{
			map[string]interface{}{"type": "css", "value": ".primary"},
			map[string]interface{}{"type": "id", "value": "submit-btn"},
		},
	}

	got, err := buildSelectorsFromStep(step)
	if err != nil {
		t.Fatalf("buildSelectorsFromStep() error = %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("len(selectors) = %d, want 2", len(got))
	}
	if got[0].Type != "css" || got[0].Value != ".primary" {
		t.Fatalf("first selector = %+v, want css/.primary", got[0])
	}
	if got[1].Type != "id" || got[1].Value != "submit-btn" {
		t.Fatalf("second selector = %+v, want id/submit-btn", got[1])
	}
}

func TestBuildSelectorsFromStep_legacyFields(t *testing.T) {
	step := map[string]interface{}{
		"selector_type": "id",
		"selector":      "checkout-btn",
	}

	got, err := buildSelectorsFromStep(step)
	if err != nil {
		t.Fatalf("buildSelectorsFromStep() error = %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len(selectors) = %d, want 1", len(got))
	}
	if got[0].Type != "id" || got[0].Value != "checkout-btn" {
		t.Fatalf("selector = %+v", got[0])
	}
}

func TestDropStaleStepSelectors_legacyUpdateWins(t *testing.T) {
	step := map[string]interface{}{
		"name":          "06 Select-Val-Val",
		"selector_type": "id",
		"selector":      "valz",
		"selectors": []interface{}{
			map[string]interface{}{"type": "id", "value": "beep"},
		},
	}
	dropStaleStepSelectors(step)
	got, err := buildSelectorsFromStep(step)
	if err != nil {
		t.Fatalf("buildSelectorsFromStep() error = %v", err)
	}
	if len(got) != 1 || got[0].Value != "valz" {
		t.Fatalf("selectors = %+v, want id/valz after dropping stale state", got)
	}
}

func TestBuildSelectorsFromStep_prefersSelectorsList(t *testing.T) {
	step := map[string]interface{}{
		"selector_type": "id",
		"selector":      "ignored",
		"selectors": []interface{}{
			map[string]interface{}{"type": "css", "value": ".btn"},
		},
	}

	got, err := buildSelectorsFromStep(step)
	if err != nil {
		t.Fatalf("buildSelectorsFromStep() error = %v", err)
	}
	if len(got) != 1 || got[0].Value != ".btn" {
		t.Fatalf("selectors = %+v, want only .btn from selectors list", got)
	}
}

func TestFlattenStepsData_singleSelectorUsesSelectorsBlock(t *testing.T) {
	steps := []sc2.StepsV2{{
		Name:      "click",
		Type:      "click_element",
		Selectors: []sc2.Selector{{Type: "id", Value: "#order"}},
	}}

	got := flattenStepsData(&steps)
	if len(got) != 1 {
		t.Fatalf("len(steps) = %d, want 1", len(got))
	}
	step := got[0].(map[string]interface{})
	selectors, ok := step["selectors"].([]interface{})
	if !ok || len(selectors) != 1 {
		t.Fatalf("selectors = %#v, want one block", step["selectors"])
	}
	sel := selectors[0].(map[string]interface{})
	if sel["type"] != "id" || sel["value"] != "#order" {
		t.Fatalf("selector block = %#v, want id/#order", sel)
	}
	if _, ok := step["selector"]; ok {
		t.Fatal("expected no legacy selector field for single API selector")
	}
	if _, ok := step["selector_type"]; ok {
		t.Fatal("expected no legacy selector_type for single API selector")
	}
}

func TestStepSelectorInputsEquivalent_legacyAndSelectorsBlock(t *testing.T) {
	oldIn := stepSelectorInput{
		selectorType: "id",
		selector:     "#order",
	}
	newIn := stepSelectorInput{
		selectors: []sc2.Selector{{Type: "id", Value: "#order"}},
	}
	if !stepSelectorInputsEquivalent(oldIn, newIn) {
		t.Fatal("expected legacy and single selectors block to be equivalent")
	}
	if !migratingFromLegacyToSelectors(oldIn, newIn) {
		t.Fatal("expected legacy state to selectors config to be a migration")
	}
}

func TestMigratingFromLegacyToSelectors_afterApplyNotMigration(t *testing.T) {
	in := stepSelectorInput{
		selectors: []sc2.Selector{{Type: "id", Value: "#order"}},
	}
	if migratingFromLegacyToSelectors(in, in) {
		t.Fatal("selectors in state and config should not be treated as migration")
	}
}

func TestFlattenStepsData_multipleSelectorsExposed(t *testing.T) {
	steps := []sc2.StepsV2{{
		Name: "click",
		Type: "click_element",
		Selectors: []sc2.Selector{
			{Type: "css", Value: ".primary"},
			{Type: "id", Value: "submit-btn"},
		},
	}}

	got := flattenStepsData(&steps)
	step := got[0].(map[string]interface{})
	selectors, ok := step["selectors"].([]interface{})
	if !ok || len(selectors) != 2 {
		t.Fatalf("selectors = %#v, want 2 blocks", step["selectors"])
	}
	if _, ok := step["selector"]; ok {
		t.Fatal("expected no legacy selector fields for multiple selectors")
	}
	if _, ok := step["selector_type"]; ok {
		t.Fatal("expected no legacy selector_type for multiple selectors")
	}
}

func TestBuildSelectorsFromStep_tooManySelectors(t *testing.T) {
	selectors := make([]interface{}, browserCheckV2MaxSelectors+1)
	for i := range selectors {
		selectors[i] = map[string]interface{}{"type": "id", "value": "x"}
	}
	step := map[string]interface{}{"selectors": selectors}

	_, err := buildSelectorsFromStep(step)
	if err == nil {
		t.Fatal("expected error for too many selectors")
	}
}
