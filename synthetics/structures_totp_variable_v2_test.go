package synthetics

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sc2 "github.com/splunk/syntheticsclient/v2/syntheticsclientv2"
)

func TestTotpVariableV2SecretSchemaIsRequiredAndSensitive(t *testing.T) {
	totpSchema := resourceTotpVariableV2().Schema["totp_variable"].Elem.(*schema.Resource).Schema
	secretSchema := totpSchema["secret"]
	if !secretSchema.Required {
		t.Fatal("secret schema must be required")
	}
	if !secretSchema.Sensitive {
		t.Fatal("secret schema must be sensitive")
	}
}

func TestTotpVariableV2Defaults(t *testing.T) {
	totpSchema := resourceTotpVariableV2().Schema["totp_variable"].Elem.(*schema.Resource).Schema
	if got := totpSchema["digits"].Default; got != 6 {
		t.Fatalf("digits default = %#v, want 6", got)
	}
	if got := totpSchema["interval"].Default; got != 30 {
		t.Fatalf("interval default = %#v, want 30", got)
	}
	if got := totpSchema["hmac_digest"].Default; got != "sha1" {
		t.Fatalf("hmac_digest default = %#v, want sha1", got)
	}
}

func TestTotpVariableDataSourcesDoNotExposeSecret(t *testing.T) {
	singleSchema := dataSourceTotpVariableV2().Schema["totp_variable"].Elem.(*schema.Resource).Schema
	if _, ok := singleSchema["secret"]; ok {
		t.Fatal("single TOTP data source must not expose secret")
	}
	listSchema := dataSourceTotpVariablesV2().Schema["totp_variables"].Elem.(*schema.Resource).Schema
	if _, ok := listSchema["secret"]; ok {
		t.Fatal("TOTP list data source must not expose secret")
	}
}

func TestBuildTotpVariableV2DataRequiresSecret(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceTotpVariableV2().Schema, map[string]interface{}{
		"totp_variable": []interface{}{
			map[string]interface{}{
				"name":        "login_mfa",
				"description": "login MFA",
				"digits":      6,
				"interval":    30,
				"hmac_digest": "sha1",
			},
		},
	})

	_, err := buildTotpVariableV2Data(d)
	if err == nil {
		t.Fatal("expected error when TOTP secret is missing")
	}
}

func TestBuildTotpVariableV2Data(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceTotpVariableV2().Schema, map[string]interface{}{
		"totp_variable": []interface{}{
			map[string]interface{}{
				"name":        "login_mfa",
				"description": "login MFA",
				"secret":      "JBSWY3DPEHPK3PXP",
				"digits":      6,
				"interval":    30,
				"hmac_digest": "sha1",
			},
		},
	})

	got, err := buildTotpVariableV2Data(d)
	if err != nil {
		t.Fatalf("buildTotpVariableV2Data returned error: %v", err)
	}
	if got.Totp.Name != "login_mfa" || got.Totp.Secret != "JBSWY3DPEHPK3PXP" || got.Totp.HmacDigest != "sha1" {
		t.Fatalf("buildTotpVariableV2Data = %#v", got.Totp)
	}
}

func TestBuildTotpVariableV2UpdateDataDoesNotSendRedactedSecret(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceTotpVariableV2().Schema, map[string]interface{}{
		"totp_variable": []interface{}{
			map[string]interface{}{
				"name":        "login_mfa",
				"description": "updated",
				"secret":      totpVariableRedactedSecret,
				"digits":      8,
				"interval":    45,
				"hmac_digest": "sha1",
			},
		},
	})

	got := buildTotpVariableV2UpdateData(d)
	if got.Totp.Secret != nil {
		t.Fatalf("Secret = %#v, want nil when secret is redacted", *got.Totp.Secret)
	}
	if got.Totp.Description == nil || *got.Totp.Description != "updated" {
		t.Fatalf("Description = %#v, want updated", got.Totp.Description)
	}
}

func TestFlattenTotpVariableV2ReadPreservesExistingStateSecret(t *testing.T) {
	resp := &sc2.TotpVariableV2Response{}
	resp.Totp.ID = 123
	resp.Totp.Name = "login_mfa"
	resp.Totp.Description = "login MFA"
	resp.Totp.Secret = totpVariableRedactedSecret
	resp.Totp.Digits = 6
	resp.Totp.Interval = 30
	resp.Totp.HmacDigest = "sha1"

	got := flattenTotpVariableV2Read(resp, "existing-sensitive-secret")
	if len(got) != 1 {
		t.Fatalf("len(flattened) = %d, want 1", len(got))
	}
	totp := got[0].(map[string]interface{})
	if totp["secret"] != "existing-sensitive-secret" {
		t.Fatalf("secret = %#v, want existing-sensitive-secret", totp["secret"])
	}
}

func TestFlattenTotpVariableV2MetadataDoesNotSetSecret(t *testing.T) {
	totp := sc2.TotpVariable{
		ID:          123,
		Name:        "login_mfa",
		Description: "login MFA",
		Secret:      totpVariableRedactedSecret,
		Digits:      6,
		Interval:    30,
		HmacDigest:  "sha1",
	}

	got := flattenTotpVariableV2Metadata(totp)
	if _, ok := got["secret"]; ok {
		t.Fatalf("metadata flattened secret unexpectedly: %#v", got)
	}
}

func TestBuildBrowserV2DataAcceptsTotpStepReference(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceBrowserCheckV2().Schema, map[string]interface{}{
		"test": []interface{}{
			map[string]interface{}{
				"name":                "browser-with-totp",
				"active":              true,
				"frequency":           5,
				"device_id":           1,
				"location_ids":        []interface{}{"aws-us-east-1"},
				"scheduling_strategy": "round_robin",
				"advanced_settings": []interface{}{
					map[string]interface{}{
						"verify_certificates":         true,
						"collect_interactive_metrics": false,
					},
				},
				"transactions": []interface{}{
					map[string]interface{}{
						"name": "Login",
						"steps": []interface{}{
							map[string]interface{}{
								"name":          "Enter MFA code",
								"type":          "enter_value",
								"selector":      "mfa-code",
								"selector_type": "id",
								"value":         "{{totp.login_mfa}}",
							},
						},
					},
				},
			},
		},
	})

	got, err := buildBrowserV2Data(d)
	if err != nil {
		t.Fatalf("buildBrowserV2Data returned error: %v", err)
	}
	if got.Test.Transactions[0].StepsV2[0].Value != "{{totp.login_mfa}}" {
		t.Fatalf("step value = %#v, want {{totp.login_mfa}}", got.Test.Transactions[0].StepsV2[0].Value)
	}
}
