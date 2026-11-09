variable "client_key_password" {
  type      = string
  sensitive = true
}

resource "synthetics_create_client_certificate_v2" "mtls" {
  client_certificate {
    name        = "api-example-client-cert"
    description = "mTLS certificate for api.example.com"
    domain      = "api.example.com"

    public_key {
      # The file may also hold the client certificate followed by its intermediate
      # certificates (up to 256 KB in total).
      content        = filebase64("${path.module}/client.crt")
      filename       = "client.crt"
      file_extension = "pem"
    }

    private_key {
      content        = filebase64("${path.module}/client.key")
      filename       = "client.key"
      file_extension = "pem"
      password       = var.client_key_password
    }
  }
}
