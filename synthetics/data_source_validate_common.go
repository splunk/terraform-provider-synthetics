// Copyright 2026 Splunk, Inc.
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
	"sort"

	sc2 "github.com/splunk/syntheticsclient/v2/syntheticsclientv2"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// validateTestIDSchema returns the optional identifier used to select
// update-style validation (validate as if updating the test with this ID)
// over create-style validation (validate as if creating a new test).
func validateTestIDSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeInt,
		Optional:    true,
		Description: "ID of an existing test to validate an update-style payload against. Omit to validate the payload as a new test.",
	}
}

// validateResultSchema returns the computed attributes populated from the
// Synthetics API's ValidateResponse: whether the payload is valid, the
// API's message, and any field-level validation errors.
func validateResultSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"valid": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"message": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"field_errors": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"field": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"messages": {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
				},
			},
		},
	}
}

// setValidateResultData sets the valid/message/field_errors attributes on d
// from a ValidateResponse returned by a Synthetics API validate call.
func setValidateResultData(d *schema.ResourceData, resp *sc2.ValidateResponse) diag.Diagnostics {
	fieldErrors, err := resp.FieldErrors()
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("valid", resp.Valid); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("message", resp.Message); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("field_errors", flattenValidateFieldErrors(fieldErrors)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func flattenValidateFieldErrors(fieldErrors map[string][]string) []interface{} {
	fields := make([]string, 0, len(fieldErrors))
	for field := range fieldErrors {
		fields = append(fields, field)
	}
	sort.Strings(fields)

	result := make([]interface{}, 0, len(fields))
	for _, field := range fields {
		messages := make([]interface{}, len(fieldErrors[field]))
		for i, message := range fieldErrors[field] {
			messages[i] = message
		}
		result = append(result, map[string]interface{}{
			"field":    field,
			"messages": messages,
		})
	}
	return result
}
