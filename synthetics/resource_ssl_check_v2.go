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
	"strconv"

	sc2 "github.com/splunk/syntheticsclient/v2/syntheticsclientv2"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceSslCheckV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSslCheckV2Create,
		ReadContext:   resourceSslCheckV2Read,
		UpdateContext: resourceSslCheckV2Update,
		DeleteContext: resourceSslCheckV2Delete,

		Schema: map[string]*schema.Schema{
			"test": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: sslCheckV2ResourceTestSchema(),
				},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func sslCheckV2ResourceTestSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeInt,
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
		"last_run_location_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"last_run_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"last_run_core_metrics_published_at": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"type": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "ssl",
		},
		"active": {
			Type:     schema.TypeBool,
			Required: true,
		},
		"frequency": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"scheduling_strategy": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      "round_robin",
			ValidateFunc: validation.StringMatch(regexp.MustCompile(`(^concurrent$|^round_robin$)`), "Setting must match concurrent or round_robin"),
		},
		"location_ids": {
			Type:     schema.TypeList,
			Required: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"automatic_retries": {
			Type:     schema.TypeInt,
			Computed: true,
			Optional: true,
		},
		"host": {
			Type:     schema.TypeString,
			Required: true,
		},
		"port": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"server_name": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"allow_self_signed": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"allow_untrusted_root": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"ca_certificate_id": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"validations": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: sslCheckV2ValidationSchema(),
			},
		},
		"custom_properties": {
			Type:     schema.TypeSet,
			Computed: true,
			Optional: true,
			Elem: &schema.Resource{
				Schema: sslCheckV2CustomPropertiesSchema(true),
			},
		},
	}
}

func sslCheckV2DataSourceTestSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeInt,
			Required: true,
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
		"last_run_location_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"last_run_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"last_run_core_metrics_published_at": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"type": {
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
		"location_ids": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"automatic_retries": {
			Type:     schema.TypeInt,
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
		"server_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"allow_self_signed": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"allow_untrusted_root": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"ca_certificate_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"validations": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: sslCheckV2ValidationSchema(),
			},
		},
		"custom_properties": {
			Type:     schema.TypeSet,
			Computed: true,
			Elem: &schema.Resource{
				Schema: sslCheckV2CustomPropertiesSchema(false),
			},
		},
	}
}

func sslCheckV2ValidationSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"actual": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringMatch(regexp.MustCompile(`^\{\{(response|headers)\.[^}]+\}\}$`), "actual must follow the format {{response.<VARIABLE_NAME>}} or {{headers.<VARIABLE_NAME>}}"),
		},
		"comparator": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"expected": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"name": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"type": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice([]string{"assert_numeric", "assert_string"}, false),
		},
	}
}

func sslCheckV2CustomPropertiesSchema(validate bool) map[string]*schema.Schema {
	keySchema := &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}
	valueSchema := &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}
	if validate {
		keySchema.ValidateFunc = validation.StringMatch(regexp.MustCompile(`^[a-zA-Z][\w.-]{0,127}$`), "custom_properties key must start with a letter and may contain letters, numbers, underscore, dot, and hyphen, up to 128 characters total with no whitespace")
		valueSchema.ValidateFunc = validation.StringMatch(regexp.MustCompile(`^.{0,256}$`), "custom_properties value must be at most 256 characters")
	}
	return map[string]*schema.Schema{
		"key":   keySchema,
		"value": valueSchema,
	}
}

func resourceSslCheckV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*sc2.Client)

	var diags diag.Diagnostics

	checkData := buildSslCheckV2Data(d)

	o, _, err := c.CreateSslCheckV2(&checkData)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.Itoa(o.Test.ID))

	return append(diags, resourceSslCheckV2Read(ctx, d, meta)...)
}

func resourceSslCheckV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*sc2.Client)

	var diags diag.Diagnostics

	checkID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	o, r, err := c.GetSslCheckV2(checkID)
	if r != nil && r.StatusCode == http.StatusNotFound {
		d.SetId("")
		log.Println("[WARN] SSL check exists in state but not in API. Removing resource from state.")
		return diags
	}
	if err != nil {
		statusCode := 0
		if r != nil {
			statusCode = r.StatusCode
		}
		log.Println("[WARN] Synthetics API error for SSL check.", checkID, err.Error(), statusCode)
		return diag.FromErr(err)
	}
	if err := d.Set("test", flattenSslCheckV2Read(o)); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceSslCheckV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*sc2.Client)

	checkID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	checkData := buildSslCheckV2UpdateData(d)
	_, _, err = c.UpdateSslCheckV2(checkID, &checkData)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceSslCheckV2Read(ctx, d, meta)
}

func resourceSslCheckV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*sc2.Client)

	var diags diag.Diagnostics

	checkID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = c.DeleteSslCheckV2(checkID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return diags
}
