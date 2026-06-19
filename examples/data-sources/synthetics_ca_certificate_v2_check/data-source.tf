# CA certificate content is read into Terraform state even when marked sensitive.
# Use encrypted, access-controlled remote state and do not commit private CA
# material to source control.
data "synthetics_ca_certificate_v2_check" "ca_certificate_v2" {
  ca_certificate {
    id = 42
  }
}
