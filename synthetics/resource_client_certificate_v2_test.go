package synthetics

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"math/big"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	sc2 "github.com/splunk/syntheticsclient/v2/syntheticsclientv2"
)

func TestClientCertificateV2SchemaMarksSecretFieldsSensitive(t *testing.T) {
	resource := resourceClientCertificateV2()
	certificate := resource.Schema["client_certificate"].Elem.(*schema.Resource)
	publicKey := certificate.Schema["public_key"].Elem.(*schema.Resource)
	privateKey := certificate.Schema["private_key"].Elem.(*schema.Resource)

	if !publicKey.Schema["content"].Sensitive {
		t.Fatal("public_key.content should be sensitive")
	}
	if !privateKey.Schema["content"].Sensitive {
		t.Fatal("private_key.content should be sensitive")
	}
	if !privateKey.Schema["password"].Sensitive {
		t.Fatal("private_key.password should be sensitive")
	}
	if !certificate.Schema["name"].ForceNew {
		t.Fatal("client_certificate.name should be ForceNew")
	}
}

func TestValidateClientCertificateName(t *testing.T) {
	for _, name := range []string{"mtls_api_example", "mtls-api-example", "mtlsAPI123"} {
		_, errors := validateClientCertificateName(name, "client_certificate.0.name")
		if len(errors) != 0 {
			t.Fatalf("name %q should be valid: %#v", name, errors)
		}
	}

	for _, name := range []string{"", "contains space", "contains.dot", "contains/slash"} {
		_, errors := validateClientCertificateName(name, "client_certificate.0.name")
		if len(errors) == 0 {
			t.Fatalf("name %q should be invalid", name)
		}
	}
}

func TestFlattenClientCertificateV2PreservesStateSecretsWhenAPIIsRedacted(t *testing.T) {
	existing := map[string]interface{}{
		"name":        "mtls_api_example",
		"description": "existing",
		"domain":      "api.example.com",
		"public_key": []interface{}{map[string]interface{}{
			"content":        "base64-public",
			"filename":       "client.crt",
			"file_extension": "pem",
		}},
		"private_key": []interface{}{map[string]interface{}{
			"content":        "base64-private",
			"filename":       "client.key",
			"file_extension": "pem",
			"password":       "key-password",
		}},
	}

	actual := flattenClientCertificateV2Read(sc2.ClientCertificate{
		ID:          123,
		Name:        "mtls_api_example",
		Description: "from-api",
		Domain:      "api.example.com",
		PublicKey: sc2.ClientCertificateKey{
			Content:       "<REDACTED>",
			Filename:      "client.crt",
			FileExtension: "pem",
		},
		PrivateKey: sc2.ClientCertificatePrivateKey{
			Content:       "<REDACTED>",
			Filename:      "client.key",
			FileExtension: "pem",
			Password:      "<REDACTED>",
		},
	}, existing)

	publicKey := actual["public_key"].([]interface{})[0].(map[string]interface{})
	privateKey := actual["private_key"].([]interface{})[0].(map[string]interface{})
	if publicKey["content"] != "base64-public" {
		t.Fatalf("public key content = %#v, want base64-public", publicKey["content"])
	}
	if privateKey["content"] != "base64-private" {
		t.Fatalf("private key content = %#v, want base64-private", privateKey["content"])
	}
	if privateKey["password"] != "key-password" {
		t.Fatalf("private key password = %#v, want key-password", privateKey["password"])
	}
}

func TestClientCertificatePasswordRemovalWithUnchangedPrivateKeyIsRejected(t *testing.T) {
	oldBlock := map[string]interface{}{
		"private_key": []interface{}{map[string]interface{}{
			"content":  "base64-private",
			"password": "key-password",
		}},
	}
	newBlock := map[string]interface{}{
		"private_key": []interface{}{map[string]interface{}{
			"content": "base64-private",
		}},
	}

	err := validateClientCertificatePasswordChange(oldBlock, newBlock)
	if err == nil {
		t.Fatal("expected private key password removal to be rejected")
	}
	if got, want := err.Error(), "private_key.password cannot be removed while private_key.content is unchanged"; got != want {
		t.Fatalf("error = %q, want %q", got, want)
	}
}

func TestClientCertificatePasswordOmittedWithReplacedPrivateKeyIsAllowed(t *testing.T) {
	oldBlock := map[string]interface{}{
		"private_key": []interface{}{map[string]interface{}{
			"content":  "base64-private-old",
			"password": "key-password",
		}},
	}
	newBlock := map[string]interface{}{
		"private_key": []interface{}{map[string]interface{}{
			"content": "base64-private-new",
		}},
	}

	if err := validateClientCertificatePasswordChange(oldBlock, newBlock); err != nil {
		t.Fatalf("expected replacement private key without password to be allowed: %s", err)
	}
}

func TestAccCreateUpdateClientCertificateV2(t *testing.T) {
	resourceName := "synthetics_create_client_certificate_v2.mtls"
	name := fmt.Sprintf("terraform-client-cert-%d", time.Now().UnixNano())
	publicKeyContent, privateKeyContent := testAccClientCertificateV2Material(t)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckClientCertificateV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + testAccClientCertificateV2Config(name, "Terraform acceptance client certificate", "client.crt", "client.key", publicKeyContent, privateKeyContent),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "client_certificate.0.name", name),
					resource.TestCheckResourceAttr(resourceName, "client_certificate.0.domain", "api.example.com"),
					resource.TestCheckResourceAttr(resourceName, "client_certificate.0.description", "Terraform acceptance client certificate"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccStateIdFunc(resourceName),
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"client_certificate.0.public_key.0.content",
					"client_certificate.0.private_key.0.content",
				},
			},
			{
				Config: providerConfig + testAccClientCertificateV2Config(name, "Terraform acceptance client certificate updated", "client-updated.crt", "client-updated.key", publicKeyContent, privateKeyContent),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "client_certificate.0.description", "Terraform acceptance client certificate updated"),
					resource.TestCheckResourceAttr(resourceName, "client_certificate.0.public_key.0.filename", "client-updated.crt"),
					resource.TestCheckResourceAttr(resourceName, "client_certificate.0.private_key.0.filename", "client-updated.key"),
				),
			},
		},
	})
}

func TestAccClientCertificateV2Attachments(t *testing.T) {
	name := fmt.Sprintf("terraform-client-cert-attach-%d", time.Now().UnixNano())
	publicKeyContent, privateKeyContent := testAccClientCertificateV2Material(t)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckClientCertificateV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + testAccClientCertificateV2AttachmentConfig(name, publicKeyContent, privateKeyContent, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("synthetics_create_client_certificate_v2.mtls", "id"),
					resource.TestCheckResourceAttrSet("synthetics_create_http_check_v2.http_mtls", "test.0.certificate_id"),
					resource.TestCheckResourceAttrSet("synthetics_create_browser_check_v2.browser_mtls", "test.0.advanced_settings.0.certificate_ids.0"),
					resource.TestCheckResourceAttrSet("synthetics_create_api_check_v2.api_mtls", "test.0.requests.0.configuration.0.certificate_id"),
				),
			},
			{
				Config: providerConfig + testAccClientCertificateV2AttachmentConfig(name, publicKeyContent, privateKeyContent, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("synthetics_create_http_check_v2.http_mtls", "test.0.certificate_id", "0"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_mtls", "test.0.advanced_settings.0.certificate_ids.#", "0"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_mtls", "test.0.requests.0.configuration.0.certificate_id", "0"),
				),
			},
		},
	})
}

func testAccClientCertificateV2Material(t *testing.T) (string, string) {
	t.Helper()

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate test private key: %s", err)
	}

	now := time.Now()
	template := x509.Certificate{
		SerialNumber: big.NewInt(now.UnixNano()),
		Subject: pkix.Name{
			CommonName: "api.example.com",
		},
		NotBefore:             now.Add(-time.Hour),
		NotAfter:              now.Add(24 * time.Hour),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
		DNSNames:              []string{"api.example.com"},
	}

	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		t.Fatalf("generate test certificate: %s", err)
	}

	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})

	return base64.StdEncoding.EncodeToString(certPEM), base64.StdEncoding.EncodeToString(keyPEM)
}

func testAccClientCertificateV2Config(name, description, publicFilename, privateFilename, publicKeyContent, privateKeyContent string) string {
	return fmt.Sprintf(`
resource "synthetics_create_client_certificate_v2" "mtls" {
  provider = synthetics.synthetics

  client_certificate {
    name        = %[1]q
    description = %[2]q
    domain      = "api.example.com"

    public_key {
      content        = %[5]q
      filename       = %[3]q
      file_extension = "pem"
    }

    private_key {
      content        = %[6]q
      filename       = %[4]q
      file_extension = "pem"
    }
  }
}
`, name, description, publicFilename, privateFilename, publicKeyContent, privateKeyContent)
}

func testAccClientCertificateV2AttachmentConfig(name, publicKeyContent, privateKeyContent string, attach bool) string {
	certificateIDLine := ""
	browserCertificateIDsLine := ""
	apiCertificateIDLine := ""
	if attach {
		certificateIDLine = "    certificate_id = tonumber(synthetics_create_client_certificate_v2.mtls.id)\n"
		browserCertificateIDsLine = "      certificate_ids = [tonumber(synthetics_create_client_certificate_v2.mtls.id)]\n"
		apiCertificateIDLine = "        certificate_id = tonumber(synthetics_create_client_certificate_v2.mtls.id)\n"
	}

	return testAccClientCertificateV2Config(name, "Terraform acceptance client certificate attachment", "client.crt", "client.key", publicKeyContent, privateKeyContent) + fmt.Sprintf(`
resource "synthetics_create_http_check_v2" "http_mtls" {
  provider = synthetics.synthetics

  test {
    active              = true
    frequency           = 5
    location_ids        = ["aws-us-west-2"]
    name                = %[1]q
    type                = "http"
    url                 = "https://api.example.com/health"
    automatic_retries   = 0
    scheduling_strategy = "round_robin"
    request_method      = "GET"
%[2]s    verify_certificates = true
  }
}

resource "synthetics_create_browser_check_v2" "browser_mtls" {
  provider = synthetics.synthetics

  test {
    active              = true
    device_id           = 1
    frequency           = 5
    location_ids        = ["aws-us-west-2"]
    automatic_retries   = 0
    name                = %[3]q
    scheduling_strategy = "round_robin"

    advanced_settings {
      verify_certificates         = true
%[4]s      collect_interactive_metrics = false
    }

    transactions {
      name = "Load application"

      steps {
        name = "Go to application"
        type = "go_to_url"
        url  = "https://app.example.com"
      }
    }
  }
}

resource "synthetics_create_api_check_v2" "api_mtls" {
  provider = synthetics.synthetics

  test {
    active              = true
    device_id           = 1
    frequency           = 5
    location_ids        = ["aws-us-west-2"]
    name                = %[5]q
    scheduling_strategy = "round_robin"
    automatic_retries   = 0

    requests {
      configuration {
        body           = ""
        headers        = {}
        name           = "Health"
        request_method = "GET"
        url            = "https://api.example.com/health"
%[6]s      }
    }
  }
}
`, name+"-http", certificateIDLine, name+"-browser", browserCertificateIDsLine, name+"-api", apiCertificateIDLine)
}

func testAccCheckClientCertificateV2Destroy(s *terraform.State) error {
	c := testAccProvider.Meta().(*sc2.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "synthetics_create_client_certificate_v2" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return err
		}

		_, details, err := c.GetClientCertificateV2(id)
		if details != nil && details.StatusCode == http.StatusNotFound {
			continue
		}
		if err != nil {
			return err
		}

		return fmt.Errorf("client certificate %s still exists", rs.Primary.ID)
	}

	return nil
}
