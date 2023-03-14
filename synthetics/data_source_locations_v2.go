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
	"log"

	sc2 "syntheticsclientv2"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceLocationsV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLocationsV2Read,
		Schema: map[string]*schema.Schema{
			"locations": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"label": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"default": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"country": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"default_location_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceLocationsV2Read(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	c := m.(*sc2.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	check, _, err := c.GetLocationsV2()
	println(check)
	if err != nil {
		return diag.FromErr(err)
	}
	
	
	locations := flattenLocationsV2Data(&check.Location)
	if err := d.Set("locations", locations); err != nil {
		return diag.FromErr(err)
	}

	defaulty := flattenDefaultLocationData(check.DefaultLocationIds)
	if err := d.Set("default_location_ids", defaulty); err != nil {
		return diag.FromErr(err)
	}


	log.Println("[DEBUG] *******************************************************************", check)



	id := fmt.Sprint(check.Location[0].ID)
	d.SetId(id)
	return diags
}
