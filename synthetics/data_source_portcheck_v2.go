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
	"fmt"

	sc2 "github.com/splunk/syntheticsclient/v2/syntheticsclientv2"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourcePortCheckV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePortCheckV2Read,
		Schema: map[string]*schema.Schema{
			"test": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"active": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"frequency": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"scheduling_strategy": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_by": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"updated_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"updated_by": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_run_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_run_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"location_ids": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"custom_properties": {
							Type:     schema.TypeSet,
							Computed: true,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"value": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"automatic_retries": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourcePortCheckV2Read(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	c := m.(*sc2.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	checkID := flattenIdData(d.Get("test"))

	check, _, err := c.GetPortCheckV2(checkID)
	println(check)
	if err != nil {
		return diag.FromErr(err)
	}

	checkTest := flattenPortCheckV2Data(check)
	if err := d.Set("test", checkTest); err != nil {
		return diag.FromErr(err)
	}

	id := fmt.Sprint(check.Test.ID)
	d.SetId(id)
	return diags
}
