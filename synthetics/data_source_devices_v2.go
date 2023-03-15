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
	"log"

	sc2 "syntheticsclientv2"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDevicesV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDevicesV2Read,
		Schema: map[string]*schema.Schema{
			"devices": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"label": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_agent": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"viewport_height": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"viewport_width": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"network_connection": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"description": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"download_bandwidth": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"latency": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"packet_loss": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"upload_bandwidth": {
										Type:     schema.TypeInt,
										Optional: true,
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

func dataSourceDevicesV2Read(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	c := m.(*sc2.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	check, _, err := c.GetDevicesV2()
	println(check)
	if err != nil {
		return diag.FromErr(err)
	}
	
	
	devices := flattenDevicesV2Data(&check.Devices)
	if err := d.Set("devices", devices); err != nil {
		return diag.FromErr(err)
	}


	log.Println("[DEBUG] *******************************************************************", check)



	id := "global_devices_synthetics"
	d.SetId(id)
	return diags
}
