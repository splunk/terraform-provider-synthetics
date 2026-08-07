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
	"net/http"
	"regexp"

	sc2 "github.com/splunk/syntheticsclient/v2/syntheticsclientv2"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var privateLocationIDPattern = regexp.MustCompile(`\Aprivate-[a-z\-]*[a-z]\z`)

func resourceLocationV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLocationV2Create,
		ReadContext:   resourceLocationV2Read,
		// Due to forcing new on every change (locations are immutable) this code shouldn't ever execute
		// Leaving code here in case this changes in the future
		// UpdateContext: resourceLocationV2Update,
		DeleteContext: resourceLocationV2Delete,

		Schema: map[string]*schema.Schema{
			"location": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringMatch(privateLocationIDPattern, "name must start with 'private-'"),
							ForceNew:     true,
						},
						"label": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"country": {
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
					},
				},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceLocationV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*sc2.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	checkData := processLocationV2Items(d)

	o, _, err := c.CreateLocationV2(&checkData)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Println("[DEBUG] created location id:", o.ID)
	d.SetId(o.ID)

	resourceLocationV2Read(ctx, d, meta)

	return diags
}

func resourceLocationV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*sc2.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	var locationID = d.Id()

	location, r, err := c.GetLocationV2(locationID)

	if r.StatusCode == http.StatusNotFound {
		d.SetId("")
		log.Println("[WARN] Resource exists in state but not in API. Removing resource from state.")
		return diags
	}
	if err != nil {
		log.Println("[WARN] Synthetics API error.", locationID, err.Error(), r.StatusCode)
		return diag.FromErr(err)
	}
	log.Println("[DEBUG] read location id:", locationID)
	if err := d.Set("location", flattenLocationV2Data(location.Location)); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceLocationV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*sc2.Client)

	var diags diag.Diagnostics

	locationID := d.Id()

	var locationIdString = locationID

	_, err := c.DeleteLocationV2(locationIdString)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Println("[DEBUG] deleted location id:", locationIdString)
	d.SetId("")

	return diags
}

func processLocationV2Items(d *schema.ResourceData) sc2.LocationV2Input {
	var check = buildLocationV2Data(d)
	return check
}
