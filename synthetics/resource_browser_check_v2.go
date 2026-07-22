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

func resourceBrowserCheckV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBrowserCheckV2Create,
		ReadContext:   resourceBrowserCheckV2Read,
		UpdateContext: resourceBrowserCheckV2Update,
		DeleteContext: resourceBrowserCheckV2Delete,

		Schema: map[string]*schema.Schema{
			"test": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"active": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"frequency": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"scheduling_strategy": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "round_robin",
							ValidateFunc: validation.StringMatch(regexp.MustCompile(`(^concurrent$|^round_robin$)`), "Setting must match concurrent or round_robin"),
						},
						"url_protocol": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"start_url": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"location_ids": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"device_id": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"advanced_settings": {
							Type:     schema.TypeSet,
							Required: true,
							Elem:     browserCheckV2AdvancedSettingsResource(false),
						},
						"transactions": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"steps": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "Unique steps for the transaction. See official [API documentation](https://dev.splunk.com/observability/reference/api/synthetics_browser/latest#endpoint-createbrowsertest) as the source of truth for descriptions and options for these values.",
										Elem: &schema.Resource{
											Schema: browserCheckV2StepSchema(false),
										},
									},
								},
							},
						},
						"custom_properties": {
							Type:     schema.TypeSet,
							Computed: true,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[a-zA-Z][\w.-]{0,127}$`), "custom_properties key must start with a letter and may contain letters, numbers, underscore, dot, and hyphen, up to 128 characters total with no whitespace"),
									},
									"value": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringMatch(regexp.MustCompile(`^.{0,256}$`), "custom_properties value must be at most 256 characters"),
									},
								},
							},
						},
						"automatic_retries": {
							Type:     schema.TypeInt,
							Computed: true,
							Optional: true,
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

func browserCheckV2AdvancedSettingsResource(computed bool) *schema.Resource {
	domainValidation := validation.StringMatch(regexp.MustCompile(`^([a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,6}$`), "Setting must be a valid domain")
	pathValidation := validation.StringMatch(regexp.MustCompile(`^\/`), "Setting must be a valid path starting with /")

	stringSchema := func(sensitive bool, validateFunc func(interface{}, string) ([]string, []error)) *schema.Schema {
		s := &schema.Schema{
			Type:      schema.TypeString,
			Sensitive: sensitive,
		}
		if computed {
			s.Computed = true
		} else {
			s.Optional = true
			if validateFunc != nil {
				s.ValidateDiagFunc = validation.ToDiagFunc(validateFunc)
			}
		}
		return s
	}

	boolSchema := func() *schema.Schema {
		s := &schema.Schema{
			Type: schema.TypeBool,
		}
		if computed {
			s.Computed = true
		} else {
			s.Optional = true
		}
		return s
	}

	setSchema := func(resource *schema.Resource) *schema.Schema {
		s := &schema.Schema{
			Type: schema.TypeSet,
			Elem: resource,
		}
		if computed {
			s.Computed = true
		} else {
			s.Optional = true
		}
		return s
	}

	verifyCertificatesSchema := &schema.Schema{
		Type: schema.TypeBool,
	}
	if computed {
		verifyCertificatesSchema.Computed = true
	} else {
		verifyCertificatesSchema.Required = true
	}

	collectInteractiveMetricsSchema := &schema.Schema{
		Type: schema.TypeBool,
	}
	if computed {
		collectInteractiveMetricsSchema.Computed = true
	} else {
		collectInteractiveMetricsSchema.Optional = true
		collectInteractiveMetricsSchema.Default = false
	}

	certificateIDsSchema := &schema.Schema{
		Type: schema.TypeList,
		Elem: &schema.Schema{
			Type:         schema.TypeInt,
			ValidateFunc: validation.IntAtLeast(1),
		},
	}
	if computed {
		certificateIDsSchema.Computed = true
	} else {
		certificateIDsSchema.Optional = true
	}

	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"user_agent":                  stringSchema(false, nil),
			"verify_certificates":         verifyCertificatesSchema,
			"collect_interactive_metrics": collectInteractiveMetricsSchema,
			"certificate_ids":             certificateIDsSchema,
			"authentication": setSchema(&schema.Resource{
				Schema: map[string]*schema.Schema{
					"username": stringSchema(false, nil),
					"password": stringSchema(true, nil),
				},
			}),
			"chrome_flags": setSchema(&schema.Resource{
				Schema: map[string]*schema.Schema{
					"name":  stringSchema(false, nil),
					"value": stringSchema(false, nil),
				},
			}),
			"cookies": setSchema(&schema.Resource{
				Schema: map[string]*schema.Schema{
					"key":    stringSchema(false, nil),
					"value":  stringSchema(true, nil),
					"domain": stringSchema(false, domainValidation),
					"path":   stringSchema(false, pathValidation),
				},
			}),
			"headers": setSchema(&schema.Resource{
				Schema: map[string]*schema.Schema{
					"name":   stringSchema(false, nil),
					"value":  stringSchema(true, nil),
					"domain": stringSchema(false, domainValidation),
				},
			}),
			"host_overrides": setSchema(&schema.Resource{
				Schema: map[string]*schema.Schema{
					"source":           stringSchema(false, nil),
					"target":           stringSchema(false, nil),
					"keep_host_header": boolSchema(),
				},
			}),
			"excluded_files": browserCheckV2ExcludedFilesSchema(computed),
		},
	}
}

func resourceBrowserCheckV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*sc2.Client)
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	checkData, err := processBrowserCheckV2Items(d)
	if err != nil {
		return diag.FromErr(err)
	}
	o, _, err := c.CreateBrowserCheckV2(&checkData)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.Itoa(o.Test.ID))

	resourceBrowserCheckV2Read(ctx, d, meta)

	return diags
}

func resourceBrowserCheckV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*sc2.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	checkID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	o, r, err := c.GetBrowserCheckV2(checkID)

	if r.StatusCode == http.StatusNotFound {
		d.SetId("")
		log.Println("[WARN] Resource exists in state but not in API. Removing resource from state.")
		return diags
	}

	if err != nil {
		log.Println("[WARN] Synthetics API error.", checkID, err.Error(), r.StatusCode)
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] read browser v2 check id=%d", checkID)
	if err := d.Set("test", flattenBrowserV2Read(o)); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceBrowserCheckV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*sc2.Client)

	var diags diag.Diagnostics

	checkID := d.Id()

	checkIdString, err := strconv.Atoi(checkID)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = c.DeleteBrowserCheckV2(checkIdString)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}

func resourceBrowserCheckV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*sc2.Client)

	checkID := d.Id()

	log.Println("[DEBUG] UPDATE BROWSER CHECK ID: ", checkID)

	checkData, err := processBrowserCheckV2Items(d)
	if err != nil {
		return diag.FromErr(err)
	}

	checkIdString, err := strconv.Atoi(checkID)
	if err != nil {
		return diag.FromErr(err)
	}

	o, _, err := c.UpdateBrowserCheckV2(checkIdString, &checkData)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] updated browser v2 check id=%s api_id=%d", checkID, o.Test.ID)

	return resourceBrowserCheckV2Read(ctx, d, meta)
}

func processBrowserCheckV2Items(d *schema.ResourceData) (sc2.BrowserCheckV2Input, error) {
	check, err := buildBrowserV2Data(d)
	if err != nil {
		return check, err
	}
	log.Printf("[DEBUG] built browser v2 check transactions=%d locations=%d", len(check.Test.Transactions), len(check.Test.LocationIds))
	return check, nil
}

var browserCheckV2SelectorTypes = []string{
	"id",
	"name",
	"xpath",
	"css",
	"link",
	"jspath",
}

func browserCheckV2SelectorSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(browserCheckV2SelectorTypes, false),
		},
		"value": {
			Type:     schema.TypeString,
			Required: true,
		},
	}
}

// browserCheckV2StepSchema returns the schema for a browser test transaction step.
// When computed is true, fields are marked computed (for data sources).
func browserCheckV2StepSchema(computed bool) map[string]*schema.Schema {
	optional := !computed

	selectorsSchema := &schema.Schema{
		Type:     schema.TypeList,
		Elem:     &schema.Resource{Schema: browserCheckV2SelectorSchema()},
		Optional: optional,
		Computed: computed,
		Description: "Element locators for this step (1-10). When set, this is sent to the API as the selectors array. " +
			"selector and selector_type are still supported as a shorthand for a single locator.",
	}

	selectorTypeSchema := &schema.Schema{
		Type:        schema.TypeString,
		Optional:    optional,
		Computed:    computed,
		Description: "Shorthand for the first selector when selectors is not used.",
	}
	selectorSchema := &schema.Schema{
		Type:        schema.TypeString,
		Optional:    optional,
		Computed:    computed,
		Description: "Shorthand for the first selector when selectors is not used.",
	}
	if !computed {
		suppress := browserCheckV2SelectorRepresentationDiffSuppress
		selectorsSchema.DiffSuppressFunc = suppress
		selectorTypeSchema.DiffSuppressFunc = suppress
		selectorSchema.DiffSuppressFunc = suppress
	}

	return map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Optional: optional,
			Computed: computed,
		},
		"type": {
			Type:     schema.TypeString,
			Optional: optional,
			Computed: computed,
		},
		"url": {
			Type:     schema.TypeString,
			Optional: optional,
			Computed: computed,
		},
		"action": {
			Type:     schema.TypeString,
			Optional: optional,
			Computed: computed,
		},
		"selectors":     selectorsSchema,
		"selector_type": selectorTypeSchema,
		"selector":      selectorSchema,
		"option_selector_type": {
			Type:     schema.TypeString,
			Optional: optional,
			Computed: computed,
		},
		"option_selector": {
			Type:     schema.TypeString,
			Optional: optional,
			Computed: computed,
		},
		"variable_name": {
			Type:     schema.TypeString,
			Optional: optional,
			Computed: computed,
		},
		"value": {
			Type:     schema.TypeString,
			Optional: optional,
			Computed: computed,
		},
		"duration": {
			Type:     schema.TypeInt,
			Optional: optional,
			Computed: computed,
		},
		"wait_for_nav": func() *schema.Schema {
			s := &schema.Schema{
				Type:     schema.TypeBool,
				Optional: optional,
				Computed: computed,
			}
			if !computed {
				s.Default = false
			}
			return s
		}(),
		"wait_for_nav_timeout": func() *schema.Schema {
			s := &schema.Schema{
				Type:     schema.TypeInt,
				Optional: optional,
				Computed: computed,
			}
			if !computed {
				s.ValidateFunc = validation.All(validation.IntAtLeast(1), validation.IntAtMost(20000))
			}
			return s
		}(),
		"wait_for_nav_timeout_default": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"max_wait_time": func() *schema.Schema {
			s := &schema.Schema{
				Type:     schema.TypeInt,
				Optional: optional,
				Computed: computed,
			}
			if !computed {
				s.ValidateFunc = validation.All(validation.IntAtLeast(1), validation.IntAtMost(90000))
			}
			return s
		}(),
		"max_wait_time_default": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"options": {
			Type:     schema.TypeSet,
			Optional: optional,
			Computed: computed,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"url": {
						Type:     schema.TypeString,
						Optional: optional,
						Computed: computed,
					},
				},
			},
		},
	}
}
