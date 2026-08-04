// Copyright 2026 Splunk, Inc.
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
	"context"

	sc2 "github.com/splunk/syntheticsclient/v2/syntheticsclientv2"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// dataSourceValidateBrowserCheckV2 validates a browser check test payload
// against the Synthetics API without creating, updating, deleting, or
// running a test. Set test_id to validate as if updating that existing
// test; omit it to validate as if creating a new test.
func dataSourceValidateBrowserCheckV2() *schema.Resource {
	schemaMap := map[string]*schema.Schema{
		"test_id": validateTestIDSchema(),
		"test": {
			Type:     schema.TypeList,
			Required: true,
			Elem: &schema.Resource{
				Schema: browserCheckV2ResourceTestSchema(),
			},
		},
	}
	for key, s := range validateResultSchema() {
		schemaMap[key] = s
	}

	return &schema.Resource{
		ReadContext: dataSourceValidateBrowserCheckV2Read,
		Schema:      schemaMap,
	}
}

func dataSourceValidateBrowserCheckV2Read(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*sc2.Client)

	checkData, err := buildBrowserV2Data(d)
	if err != nil {
		return diag.FromErr(err)
	}

	var resp *sc2.ValidateResponse

	if testID, ok := d.GetOk("test_id"); ok {
		resp, _, err = c.ValidateBrowserCheckV2(testID.(int), &checkData)
	} else {
		resp, _, err = c.ValidateNewBrowserCheckV2(&checkData)
	}
	if err != nil {
		return diag.FromErr(err)
	}

	if diags := setValidateResultData(d, resp); diags != nil {
		return diags
	}

	d.SetId("validate_browser_check_v2")
	return nil
}
