---
page_title: "Globalscape EFT: server_smtp Resource"
description: |-
  Configures the Globalscape EFT server SMTP settings.
---

# Resource `globalscapeeft_server_smtp`

Controls the singleton SMTP configuration exposed by `PATCH /admin/v2/server`. Use this resource to configure the mail server that EFT uses for alerts and notifications.

**Important Notes:**
- Deleting this resource removes it from Terraform state only. The SMTP configuration remains on the EFT server. To clear SMTP settings, update the resource with empty/default values before destroying.
- The `password` attribute is write-only and will not be stored in Terraform state after creation or refresh for security reasons.

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
- `password` (String, Sensitive) SMTP account password. **Note:** This value is write-only and will not be stored in state after creation or updates for security.
- `use_authentication` (Boolean) Whether SMTP AUTH is required.
- `use_implicit_tls` (Boolean) Enable implicit TLS (SMTPS).

### Timeouts

This resource supports customizable timeouts for operations:
- `create` - Default: 5 minutes
- `read` - Default: 5 minutes
- `update` - Default: 5 minutes
- `delete` - Default: 5 minutes

Example:
```hcl
resource "globalscapeeft_server_smtp" "default" {
  # ... other configuration ...

  timeouts {
    create = "10m"
    update = "10m"
  }
}
```

### Read-only

- `id` (String) Static identifier for the singleton server settings.

## Import

Import the server SMTP settings using any identifier (typically "server" or "1"):

```bash
terraform import globalscapeeft_server_smtp.default server
```
