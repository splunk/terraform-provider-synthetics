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
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var browserCheckV2StepAttrPrefixRe = regexp.MustCompile(`^(test\.\d+\.transactions\.\d+\.steps\.\d+)\.`)

// browserCheckV2SelectorRepresentationDiffSuppress suppresses cosmetic plan diffs when
// state and config describe the same single selector via different field shapes. It does
// not suppress the one-time migration from legacy fields in state to a selectors block
// in config.
func browserCheckV2SelectorRepresentationDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	prefix := browserCheckV2StepPrefixFromAttributeKey(k)
	if prefix == "" {
		return false
	}
	oldIn := stepSelectorInputFromResourceData(d, prefix, true)
	newIn := stepSelectorInputFromResourceData(d, prefix, false)
	if !stepSelectorInputsEquivalent(oldIn, newIn) {
		return false
	}
	if migratingFromLegacyToSelectors(oldIn, newIn) {
		return false
	}
	return stepSelectorRepresentationDiffers(oldIn, newIn)
}

func browserCheckV2StepPrefixFromAttributeKey(key string) string {
	m := browserCheckV2StepAttrPrefixRe.FindStringSubmatch(key)
	if len(m) < 2 {
		return ""
	}
	return m[1]
}

func stepSelectorInputFromResourceData(d *schema.ResourceData, prefix string, useState bool) stepSelectorInput {
	var in stepSelectorInput
	in.selectorType = stringFromResourceDataField(d, prefix+".selector_type", useState)
	in.selector = stringFromResourceDataField(d, prefix+".selector", useState)
	if raw := interfaceFromResourceDataField(d, prefix+".selectors", useState); raw != nil {
		in.selectors = parseSelectorsList(raw)
	}
	return in
}

func stringFromResourceDataField(d *schema.ResourceData, key string, useState bool) string {
	v := interfaceFromResourceDataField(d, key, useState)
	if v == nil {
		return ""
	}
	s, _ := v.(string)
	return s
}

func interfaceFromResourceDataField(d *schema.ResourceData, key string, useState bool) interface{} {
	if d.HasChange(key) {
		o, n := d.GetChange(key)
		if useState {
			return o
		}
		return n
	}
	v, ok := d.GetOk(key)
	if !ok {
		return nil
	}
	return v
}
