variable "site_id" {}

resource "globalscapeeft_event_rule" "example" {
  site_id = var.site_id

  attributes_json = jsonencode({
    info = {
      Name        = "Example Terraform Rule"
      Description = "Created via Terraform"
      Enabled     = true
      Folder      = "00000000-0000-0000-0000-000000000000"
      Remote      = false
      Type        = "Timer"
    }
    statements = {
      StatementsList = []
    }
  })
}
