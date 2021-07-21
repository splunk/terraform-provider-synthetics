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
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sc "github.com/splunk/syntheticsclient/syntheticsclient"
)

func resourceBrowserCheck() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBrowserCheckCreate,
		ReadContext:   resourceBrowserCheckRead,
		UpdateContext: resourceBrowserCheckUpdate,
		DeleteContext: resourceBrowserCheckDelete,

		Schema: map[string]*schema.Schema{
			"last_updated": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "real_browser",
			},
			"frequency": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"round_robin": {
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
			},
			"auto_retry": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"http_request_body": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"http_method": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"user_agent": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Mozilla/5.0 (X11; Linux x86_64; Rigor) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36",
			},
			"auto_update_user_agent": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"enforce_ssl_validation": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"wait_for_full_metrics": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"integrations": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"viewport": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"width": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"height": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"dns_overrides": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"original_domain": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"original_host": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"browser": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "chrome",
						},
					},
				},
			},
			"steps": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"item_method": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"value": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"how": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"what": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"variable_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"position": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"cookies": {
				Type:     schema.TypeSet,
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
						"domain": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"path": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"javascript_files": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"threshold_monitors": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"matcher": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"metric_name": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
								v := val.(string)
								switch v {
								case
									"first_byte_time_ms",
									"dom_interactive_time_ms",
									"dom_load_time_ms",
									"dom_complete_time_ms",
									"start_render_ms",
									"onload_time_ms",
									"visually_complete_ms",
									"fully_loaded_time_ms",
									"first_paint_time_ms",
									"first_contentful_paint_time_ms",
									"first_meaningful_paint_time_ms",
									"first_interactive_time_ms",
									"first_cpu_idle_time_ms",
									"first_request_dns_time_ms",
									"first_request_connect_time_ms",
									"first_request_ssl_time_ms",
									"first_request_send_time_ms",
									"first_request_wait_time_ms",
									"first_request_receive_time_ms",
									"speed_index",
									"requests",
									"content_bytes",
									"html_files",
									"html_bytes",
									"image_files",
									"image_bytes",
									"javascript_files",
									"javascript_bytes",
									"css_files",
									"css_bytes",
									"video_files",
									"video_bytes",
									"font_files",
									"font_bytes",
									"other_files",
									"other_bytes",
									"client_errors",
									"connection_errors",
									"server_errors",
									"errors",
									"run_count",
									"success_count",
									"failure_count",
									"lighthouse_performance_score",
									"availability",
									"downtime",
									"total_blocking_time_ms",
									"largest_contentful_paint_time_ms",
									"cumulative_layout_shift":
									return
								}
								errs = append(errs, fmt.Errorf("%s is not a valid metric_name. Please check your input against API docs at https://monitoring-api.rigor.com", v))
								return
							},
						},
						"comparison_type": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
								v := val.(string)
								switch v {
								case
									"less_than",
									"equals",
									"greater_than":
									return
								}
								errs = append(errs, fmt.Errorf("%s is not a valid comparison_type. Please check your input against API docs at https://monitoring-api.rigor.com", v))
								return
							},
						},
						"value": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"excluded_files": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"pattern": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"preset_name": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
								v := val.(string)
								switch v {
								case
									"chartbeat",
									"clicktale",
									"comscore",
									"coremetrics",
									"crazy-egg",
									"eloqua",
									"gomez",
									"google-analytics",
									"hubspot",
									"liveperson",
									"mixpanel",
									"omniture",
									"optimizely",
									"pardot",
									"quantcast",
									"spectate",
									"tealium",
									"white-ops":
									return
								}
								errs = append(errs, fmt.Errorf("%s is not a valid preset_name. Please check your input against API docs at https://monitoring-api.rigor.com", v))
								return
							},
						},
						"exclusion_type": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
								v := val.(string)
								switch v {
								case
									"preset",
									"custom",
									"all_except":
									return
								}
								errs = append(errs, fmt.Errorf("%s is not a valid exclusion_type. Please check your input against API docs at https://monitoring-api.rigor.com", v))
								return
							},
						},
					},
				},
			},
			"notifications": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"sms": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"call": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"email": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"notify_after_failure_count": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"notify_on_location_failure": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"muted": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"notify_who": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"sms": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"call": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"email": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"custom_user_email": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"type": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"links": {
										Type:     schema.TypeSet,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"self_html": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"id": {
										Type:     schema.TypeInt,
										Computed: true,
										Optional: true,
									},
								},
							},
						},
						"notification_windows": {
							Type:     schema.TypeSet,
							Computed: true,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"start_time": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"end_time": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"duration_in_minutes": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"time_zone": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"escalations": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"sms": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"call": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"email": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"after_minutes": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"is_repeat": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"notify_who": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": {
													Type:     schema.TypeInt,
													Optional: true,
												},
												"sms": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"call": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"email": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"custom_user_email": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"type": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"links": {
													Type:     schema.TypeSet,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"self_html": {
																Type:     schema.TypeString,
																Computed: true,
															},
														},
													},
												},
											},
										},
									},
									"notification_window": {
										Type:     schema.TypeSet,
										Computed: true,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"start_time": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"end_time": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"duration_in_minutes": {
													Type:     schema.TypeInt,
													Optional: true,
												},
												"time_zone": {
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
			"http_request_headers": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_agent": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"blackout_periods": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"end_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"timezone": {
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
						"repeat_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"duration_in_minutes": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"is_repeat": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"monthly_repeat_type": {
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
					},
				},
			},
			"locations": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"success_criteria": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"created_at": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"updated_at": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"comparison_string": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"check_connection": {
				// `connection` is a reserved field name in Terraform
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"download_bandwidth": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"upload_bandwidth": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"latency": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"packet_loss": {
							Type:     schema.TypeFloat,
							Optional: true,
						},
					},
				},
			},
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceBrowserCheckCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*sc.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	check := sc.BrowserCheckInput{}

	checkData := processBrowserCheckItems(d, check)

	o, _, err := c.CreateBrowserCheck(&checkData)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.Itoa(o.ID))

	resourceBrowserCheckRead(ctx, d, meta)

	return diags
	// return nil
}

func resourceBrowserCheckRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*sc.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	checkID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	check, _, err := c.GetCheck(checkID)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Println("[DEBUG] GET check response data: ", check)

	return diags
}

func resourceBrowserCheckDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*sc.Client)

	var diags diag.Diagnostics

	checkID := d.Id()

	checkIdString, err := strconv.Atoi(checkID)
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := c.DeleteBrowserCheck(checkIdString)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Println("[DEBUG] Delete check response data: ", resp)
	d.SetId("")

	return diags
}

func resourceBrowserCheckUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*sc.Client)

	checkID := d.Id()

	check := sc.BrowserCheckInput{}

	checkData := processBrowserCheckItems(d, check)

	checkIdString, err := strconv.Atoi(checkID)
	if err != nil {
		return diag.FromErr(err)
	}

	o, _, err := c.UpdateBrowserCheck(checkIdString, &checkData)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Println("[DEBUG] Update check response data: ", o)
	d.Set("last_updated", time.Now().Format(time.RFC850))

	return resourceBrowserCheckRead(ctx, d, meta)
}

func processBrowserCheckItems(d *schema.ResourceData, check sc.BrowserCheckInput) sc.BrowserCheckInput {

	//These MUST exist or the request will fail
	check.Name = d.Get("name").(string)
	check.Frequency = d.Get("frequency").(int)
	check.Type = d.Get("type").(string)
	check.URL = d.Get("url").(string)

	check.RoundRobin = d.Get("round_robin").(bool)
	check.AutoRetry = d.Get("auto_retry").(bool)
	check.Enabled = d.Get("enabled").(bool)
	check.HTTPRequestBody = d.Get("http_request_body").(string)
	check.HTTPMethod = d.Get("http_method").(string)
	check.AutoUpdateUserAgent = d.Get("auto_update_user_agent").(bool)
	check.EnforceSslValidation = d.Get("enforce_ssl_validation").(bool)
	check.WaitForFullMetrics = d.Get("wait_for_full_metrics").(bool)

	check.Viewport = buildViewportData(d)
	check.Connection = buildConnectionData(d)
	check.Tags = buildTagsData(d)
	check.Locations = buildLocationsData(d)
	check.Integrations = buildIntegrationsData(d)
	check.Notifications = buildNotificationsData(check.Notifications, d)
	check.ExcludedFiles = buildExcludedFilesData(d)
	check.Cookies = buildCookieData(d)
	check.DNSOverrides = buildDnsOverridesData(d)
	check.ThresholdMonitors = buildThresholdMonitorsData(d)
	// Currently javascript_files and steps settings are not available via public API endpoints
	check.JavascriptFiles = buildJavascriptFilesData(d)
	check.Steps = buildStepData(d)

	return check
}
