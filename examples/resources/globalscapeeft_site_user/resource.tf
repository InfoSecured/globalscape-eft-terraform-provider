variable "site_id" {}

resource "globalscapeeft_site_user" "example" {
  site_id         = var.site_id
  login_name      = "tf-example"
  password        = "Secur3P@ss!"
  password_type   = "Default"
  display_name    = "Terraform Example"
  email           = "tf@example.com"
  account_enabled = "yes"
}
