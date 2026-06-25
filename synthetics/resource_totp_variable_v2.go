package synthetics

import (
	"context"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	sc2 "github.com/splunk/syntheticsclient/v2/syntheticsclientv2"
)

func resourceTotpVariableV2() *schema.Resource {
	return &schema.Resource{
		Description:   "Manages a Synthetics TOTP variable. The TOTP secret is stored in Terraform state even when marked sensitive. Use encrypted, access-controlled remote state and do not commit state files or real TOTP seeds.",
		CreateContext: resourceTotpVariableV2Create,
		ReadContext:   resourceTotpVariableV2Read,
		UpdateContext: resourceTotpVariableV2Update,
		DeleteContext: resourceTotpVariableV2Delete,
		Schema: map[string]*schema.Schema{
			"totp_variable": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: totpVariableV2ResourceSchema(),
				},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func totpVariableV2ResourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[A-Za-z0-9_-]+$`), "TOTP variable name may contain only letters, numbers, underscore, and hyphen"),
		},
		"description": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"secret": {
			Type:         schema.TypeString,
			Required:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"digits": {
			Type:         schema.TypeInt,
			Optional:     true,
			Default:      6,
			ValidateFunc: validation.IntBetween(1, 10),
		},
		"interval": {
			Type:         schema.TypeInt,
			Optional:     true,
			Default:      30,
			ValidateFunc: validation.IntBetween(10, 120),
		},
		"hmac_digest": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      "sha1",
			ValidateFunc: validation.StringInSlice([]string{"sha1"}, false),
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
	}
}

func resourceTotpVariableV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*sc2.Client)
	var diags diag.Diagnostics

	totpVariableData, err := buildTotpVariableV2Data(d)
	if err != nil {
		return diag.FromErr(err)
	}

	o, _, err := c.CreateTotpVariableV2(&totpVariableData)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.Itoa(o.Totp.ID))

	return append(diags, resourceTotpVariableV2Read(ctx, d, meta)...)
}

func resourceTotpVariableV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*sc2.Client)
	var diags diag.Diagnostics

	totpVariableID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	existingSecret := totpVariableSecretFromState(d)
	o, r, err := c.GetTotpVariableV2(totpVariableID)
	if r != nil && r.StatusCode == http.StatusNotFound {
		d.SetId("")
		log.Println("[WARN] TOTP variable exists in state but not in API. Removing resource from state.")
		return diags
	}
	if err != nil {
		statusCode := 0
		if r != nil {
			statusCode = r.StatusCode
		}
		log.Println("[WARN] Synthetics API error for TOTP variable.", totpVariableID, err.Error(), statusCode)
		return diag.FromErr(err)
	}

	if err := d.Set("totp_variable", flattenTotpVariableV2Read(o, existingSecret)); err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func resourceTotpVariableV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*sc2.Client)

	totpVariableID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	totpVariableData := buildTotpVariableV2UpdateData(d)
	_, _, err = c.UpdateTotpVariableV2(totpVariableID, &totpVariableData)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceTotpVariableV2Read(ctx, d, meta)
}

func resourceTotpVariableV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*sc2.Client)
	var diags diag.Diagnostics

	totpVariableID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	statusCode, err := c.DeleteTotpVariableV2(totpVariableID)
	if err != nil {
		if statusCode == http.StatusNotFound {
			d.SetId("")
			return diags
		}
		return diag.FromErr(err)
	}

	d.SetId("")
	return diags
}
