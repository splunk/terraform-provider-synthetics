data "synthetics_totp_variable_v2_check" "login_mfa" {
  totp_variable {
    id = 123
  }
}
