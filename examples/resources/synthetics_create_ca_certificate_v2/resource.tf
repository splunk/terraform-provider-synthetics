# CA certificate content is stored in Terraform state even when marked sensitive.
# Use encrypted, access-controlled remote state and do not commit private CA
# material to source control.
resource "synthetics_create_ca_certificate_v2" "ca_certificate_v2" {
  ca_certificate {
    name           = "Terraform - CA Certificate V2"
    description    = "Example private CA certificate for SSL tests"
    content        = filebase64("${path.module}/ca.pem")
    file_extension = "pem"
    filename       = "ca.pem"
  }
}
