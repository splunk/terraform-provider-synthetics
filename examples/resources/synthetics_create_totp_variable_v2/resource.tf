variable "login_mfa_totp_secret" {
  type      = string
  sensitive = true
}

resource "synthetics_create_totp_variable_v2" "login_mfa" {
  totp_variable {
    name        = "login_mfa"
    description = "TOTP seed for login browser test"
    secret      = var.login_mfa_totp_secret
    digits      = 6
    interval    = 30
    hmac_digest = "sha1"
  }
}
