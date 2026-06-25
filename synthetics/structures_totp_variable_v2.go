package synthetics

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sc2 "github.com/splunk/syntheticsclient/v2/syntheticsclientv2"
)

const totpVariableRedactedSecret = "<REDACTED>"

func buildTotpVariableV2Data(d *schema.ResourceData) (sc2.TotpVariableV2Input, error) {
	var input sc2.TotpVariableV2Input
	totp, ok := firstTotpVariableBlock(d)
	if !ok {
		return input, fmt.Errorf("totp_variable block is required")
	}

	secret := totpStringField(totp, "secret")
	if secret == "" || secret == totpVariableRedactedSecret {
		return input, fmt.Errorf("totp_variable secret is required")
	}

	input.Totp.Name = totpStringField(totp, "name")
	input.Totp.Description = totpStringField(totp, "description")
	input.Totp.Secret = secret
	input.Totp.Digits = totpIntField(totp, "digits")
	input.Totp.Interval = totpIntField(totp, "interval")
	input.Totp.HmacDigest = totpStringField(totp, "hmac_digest")
	return input, nil
}

func buildTotpVariableV2UpdateData(d *schema.ResourceData) sc2.TotpVariableV2UpdateInput {
	var input sc2.TotpVariableV2UpdateInput
	totp, ok := firstTotpVariableBlock(d)
	if !ok {
		return input
	}

	description := totpStringField(totp, "description")
	digits := totpIntField(totp, "digits")
	interval := totpIntField(totp, "interval")
	hmacDigest := totpStringField(totp, "hmac_digest")

	input.Totp.Description = &description
	input.Totp.Digits = &digits
	input.Totp.Interval = &interval
	input.Totp.HmacDigest = &hmacDigest

	secret := totpStringField(totp, "secret")
	if d.HasChange("totp_variable.0.secret") && secret != "" && secret != totpVariableRedactedSecret {
		input.Totp.Secret = &secret
	}

	return input
}

func firstTotpVariableBlock(d *schema.ResourceData) (map[string]interface{}, bool) {
	raw, ok := d.Get("totp_variable").([]interface{})
	if !ok || len(raw) == 0 || raw[0] == nil {
		return nil, false
	}
	totp, ok := raw[0].(map[string]interface{})
	return totp, ok
}

func totpVariableIDFromList(value interface{}) int {
	raw, ok := value.([]interface{})
	if !ok || len(raw) == 0 || raw[0] == nil {
		return 0
	}
	totp, ok := raw[0].(map[string]interface{})
	if !ok {
		return 0
	}
	return totpIntField(totp, "id")
}

func totpVariableSecretFromState(d *schema.ResourceData) string {
	totp, ok := firstTotpVariableBlock(d)
	if !ok {
		return ""
	}
	return totpStringField(totp, "secret")
}

func flattenTotpVariableV2Read(resp *sc2.TotpVariableV2Response, existingSecret string) []interface{} {
	if resp == nil {
		return []interface{}{}
	}
	totp := flattenTotpVariableV2Metadata(resp.Totp)
	if secret := totpVariableSecretForState(resp.Totp.Secret, existingSecret); secret != "" {
		totp["secret"] = secret
	}
	return []interface{}{totp}
}

func flattenTotpVariableV2Data(resp *sc2.TotpVariableV2Response) []interface{} {
	if resp == nil {
		return []interface{}{}
	}
	return []interface{}{flattenTotpVariableV2Metadata(resp.Totp)}
}

func flattenTotpVariablesV2Data(totps []sc2.TotpVariable) []interface{} {
	result := make([]interface{}, len(totps))
	for i, totp := range totps {
		result[i] = flattenTotpVariableV2Metadata(totp)
	}
	return result
}

func flattenTotpVariableV2Metadata(totp sc2.TotpVariable) map[string]interface{} {
	result := make(map[string]interface{})
	if totp.ID != 0 {
		result["id"] = totp.ID
	}
	if totp.Name != "" {
		result["name"] = totp.Name
	}
	if totp.Description != "" {
		result["description"] = totp.Description
	}
	if totp.Digits != 0 {
		result["digits"] = totp.Digits
	}
	if totp.Interval != 0 {
		result["interval"] = totp.Interval
	}
	if totp.HmacDigest != "" {
		result["hmac_digest"] = totp.HmacDigest
	}
	if !totp.CreatedAt.IsZero() {
		result["created_at"] = totp.CreatedAt.Format(time.RFC3339)
	}
	if totp.CreatedBy != "" {
		result["created_by"] = totp.CreatedBy
	}
	if !totp.UpdatedAt.IsZero() {
		result["updated_at"] = totp.UpdatedAt.Format(time.RFC3339)
	}
	if totp.UpdatedBy != "" {
		result["updated_by"] = totp.UpdatedBy
	}
	return result
}

func totpVariableSecretForState(apiSecret, existingSecret string) string {
	if existingSecret != "" && (apiSecret == "" || apiSecret == totpVariableRedactedSecret) {
		return existingSecret
	}
	if apiSecret == totpVariableRedactedSecret {
		return ""
	}
	return apiSecret
}

func totpStringField(totp map[string]interface{}, key string) string {
	if value, ok := totp[key].(string); ok {
		return value
	}
	return ""
}

func totpIntField(totp map[string]interface{}, key string) int {
	if value, ok := totp[key].(int); ok {
		return value
	}
	return 0
}
