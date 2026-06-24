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
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sc2 "github.com/splunk/syntheticsclient/v2/syntheticsclientv2"
)

func dataSourceExcludedFileTypesV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceExcludedFileTypesV2Read,
		Schema: map[string]*schema.Schema{
			"excluded_file_types": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceExcludedFileTypesV2Read(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*sc2.Client)
	var diags diag.Diagnostics

	excludedFileTypes, _, err := c.GetExcludedFileTypesV2()
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("excluded_file_types", excludedFileTypes.ExcludedFileTypes); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("global_excluded_file_types_synthetics")
	return diags
}
