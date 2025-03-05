// Copyright 2024 Splunk, Inc.
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

func dataSourceDowntimeConfigurationV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDowntimeConfigurationV2Read,
		Schema: map[string]*schema.Schema{
			"downtime_configuration": {
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
						"description": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
						},
						"rule": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"start_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"end_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"updated_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tests_updated_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"test_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"timezone": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
						},
						"recurrence": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"repeats": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"custom_value": {
													Type:     schema.TypeInt,
													Optional: true,
												},
												"custom_frequency": {
													Type:     schema.TypeString,
													Optional: true,
												},
											},
										},
									},
									"end": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": {
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
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceDowntimeConfigurationV2Read(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	c := m.(*sc2.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	downtimeConfigID := flattenIdData(d.Get("downtime_configuration"))

	downtimeConfig, _, err := c.GetDowntimeConfigurationV2(downtimeConfigID)
	println(downtimeConfig)
	if err != nil {
		return diag.FromErr(err)
	}

	DowntimeConfiguration := flattenDowntimeConfigurationV2Data(downtimeConfig)
	if err := d.Set("downtime_configuration", DowntimeConfiguration); err != nil {
		return diag.FromErr(err)
	}

	id := fmt.Sprint(downtimeConfig.ID)
	d.SetId(id)
	return diags
}
