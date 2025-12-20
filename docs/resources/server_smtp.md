---
page_title: "Globalscape EFT: server_smtp Resource"
description: |-
  Configures the Globalscape EFT server SMTP settings.
---

# Resource `globalscapeeft_server_smtp`

Controls the singleton SMTP configuration exposed by `PATCH /admin/v2/server`. Use this resource to configure the mail server that EFT uses for alerts and notifications.

Deleting this resource removes it from Terraform state only. It does not reset the remote SMTP settings because EFT exposes a single global configuration.

## Example Usage

```hcl
resource "globalscapeeft_server_smtp" "default" {
  server              = "smtp.example.net"
  port                = 587
  sender_address      = "eft@example.net"
  sender_name         = "Globalscape EFT"
  login               = "smtp-user"
  password            = var.smtp_password
  use_authentication  = true
  use_implicit_tls    = false
}
```

## Schema

### Required

- `server` (String) SMTP host name or IP address.
- `port` (Number) SMTP port number.
- `sender_address` (String) Email address used in the `From` header.
- `sender_name` (String) Display name portion of the `From` header.

### Optional

- `login` (String) SMTP account username.
- `password` (String, Sensitive) SMTP account password.
- `use_authentication` (Boolean) Whether SMTP AUTH is required.
- `use_implicit_tls` (Boolean) Enable implicit TLS (SMTPS).

### Read-only

- `id` (String) Static identifier for the singleton server settings.
