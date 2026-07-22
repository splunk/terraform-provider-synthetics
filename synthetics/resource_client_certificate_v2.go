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
	"net/http"
	"regexp"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sc2 "github.com/splunk/syntheticsclient/v2/syntheticsclientv2"
)

var clientCertificateNamePattern = regexp.MustCompile(`^[A-Za-z0-9_-]+$`)

const clientCertificateAuthMaterialDescription = "Client certificates are authentication material for synthetic tests. They are separate from CA certificates used as SSL trust material."

const clientCertificateStateWarning = "Client certificate public key content, private key content, and private key password are stored in Terraform state when managed by this resource. Use encrypted, access-controlled remote state and do not commit Terraform state files or certificate material. API reads return certificate material redacted; imports can recover metadata but cannot recover key content or private-key password."

func validateClientCertificateName(value interface{}, key string) ([]string, []error) {
	name, ok := value.(string)
	if !ok || name == "" {
		return nil, []error{fmt.Errorf("%s must be a non-empty string", key)}
	}
	if !clientCertificateNamePattern.MatchString(name) {
		return nil, []error{fmt.Errorf("%s must contain only letters, numbers, underscores, and hyphens", key)}
	}
	return nil, nil
}

func resourceClientCertificateV2() *schema.Resource {
	return &schema.Resource{
		Description:   "Manages a Synthetics client certificate for mTLS authentication. " + clientCertificateAuthMaterialDescription + " " + clientCertificateStateWarning,
		CreateContext: resourceClientCertificateV2Create,
		ReadContext:   resourceClientCertificateV2Read,
		UpdateContext: resourceClientCertificateV2Update,
		DeleteContext: resourceClientCertificateV2Delete,
		Schema: map[string]*schema.Schema{
			"client_certificate": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: clientCertificateV2ResourceSchema(),
				},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func clientCertificateV2ResourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateClientCertificateName,
		},
		"description": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"domain": {
			Type:     schema.TypeString,
			Required: true,
		},
		"expires_at": {
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
		"public_key":  clientCertificateKeySchema(true),
		"private_key": clientCertificatePrivateKeySchema(),
	}
}

func clientCertificateKeySchema(required bool) *schema.Schema {
	s := &schema.Schema{
		Type:     schema.TypeList,
		MaxItems: 1,
		Elem: &schema.Resource{Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"content": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"filename": {
				Type:     schema.TypeString,
				Required: true,
			},
			"file_extension": {
				Type:     schema.TypeString,
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
		}},
	}
	if required {
		s.Required = true
	} else {
		s.Computed = true
	}
	return s
}

func clientCertificatePrivateKeySchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		MaxItems: 1,
		Elem: &schema.Resource{Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"content": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"filename": {
				Type:     schema.TypeString,
				Required: true,
			},
			"file_extension": {
				Type:     schema.TypeString,
				Required: true,
			},
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
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
		}},
	}
}

func resourceClientCertificateV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*sc2.Client)
	var diags diag.Diagnostics

	response, _, err := c.CreateClientCertificateV2(buildClientCertificateV2Data(d))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.Itoa(response.Certificate.ID))

	return append(diags, resourceClientCertificateV2Read(ctx, d, meta)...)
}

func resourceClientCertificateV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*sc2.Client)
	var diags diag.Diagnostics

	clientCertificateID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	response, details, err := c.GetClientCertificateV2(clientCertificateID)
	if details != nil && details.StatusCode == http.StatusNotFound {
		d.SetId("")
		log.Println("[WARN] Client certificate exists in state but not in API. Removing resource from state.")
		return diags
	}
	if err != nil {
		statusCode := 0
		if details != nil {
			statusCode = details.StatusCode
		}
		log.Println("[WARN] Synthetics API error for client certificate.", clientCertificateID, err.Error(), statusCode)
		return diag.FromErr(err)
	}

	existing := firstMapFromList(d.Get("client_certificate"))
	if err := d.Set("client_certificate", []interface{}{flattenClientCertificateV2Read(response.Certificate, existing)}); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceClientCertificateV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*sc2.Client)

	clientCertificateID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	clientCertificateData, err := buildClientCertificateV2UpdateData(d)
	if err != nil {
		return diag.FromErr(err)
	}

	_, _, err = c.UpdateClientCertificateV2(clientCertificateID, clientCertificateData)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceClientCertificateV2Read(ctx, d, meta)
}

func resourceClientCertificateV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*sc2.Client)
	var diags diag.Diagnostics

	clientCertificateID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	statusCode, err := c.DeleteClientCertificateV2(clientCertificateID)
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
