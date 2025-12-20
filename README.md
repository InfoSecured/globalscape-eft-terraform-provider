# Globalscape EFT Terraform Provider

This repository contains a Terraform provider built with the [HashiCorp Plugin Framework](https://developer.hashicorp.com/terraform/plugin/framework) that manages parts of the Globalscape EFT Server REST API. The implementation currently focuses on:

- Authenticating to the EFT Admin REST API using local or AD credentials.
- Reading high-level server metadata (version, general settings, listener configuration, SMTP) via the `globalscapeeft_server` data source.
- Enumerating configured sites with the `globalscapeeft_sites` data source.
- Managing the server-wide SMTP configuration with the `globalscapeeft_server_smtp` resource.
- Managing site users via the `globalscapeeft_site_user` resource.
- Creating, updating, and deleting event rules with the `globalscapeeft_event_rule` resource by manipulating EFT's JSON payloads directly.

## Building the provider

```bash
# Run from the repository root
GOMODCACHE=$(pwd)/.gomodcache GOCACHE=$(pwd)/.gocache go build ./...
```

Use Go 1.22 or newer. The custom cache paths ensure that module downloads stay inside the workspace when sandboxed.

## Provider configuration

```hcl
provider "globalscapeeft" {
  host                 = "https://eft.example.com:4450/admin"
  username             = var.eft_username
  password             = var.eft_password
  auth_type            = "EFT"         # or "AD"
  insecure_skip_verify = true           # when using lab/self-signed certs
}
```

- `host` must include the `/admin` base path. All REST calls append `/v1` or `/v2` to this path.
- TLS verification can be disabled for appliances with self-signed certificates.

## Resources and data sources

### Data source `globalscapeeft_server`

Returns read-only information about the EFT server, including SMTP settings. Example:

```hcl
data "globalscapeeft_server" "current" {}

output "eft_version" {
  value = data.globalscapeeft_server.current.version
}
```

### Data source `globalscapeeft_sites`

Returns every site configured on the EFT server, exposing each site's ID for use with user resources.

```hcl
data "globalscapeeft_sites" "all" {}

output "first_site" {
  value = data.globalscapeeft_sites.all.sites[0].name
}
```

### Resource `globalscapeeft_server_smtp`

Manages the singleton SMTP configuration returned by `PATCH /admin/v2/server`.

```hcl
resource "globalscapeeft_server_smtp" "default" {
  server         = "smtp.example.net"
  port           = 587
  sender_address = "eft@example.net"
  sender_name    = "Globalscape EFT"
  login          = "smtp-user"
  password       = var.smtp_password
  use_authentication = true
  use_implicit_tls   = false
}
```

Deleting the resource only removes it from state because EFT exposes a single set of SMTP settings per server instance.

### Resource `globalscapeeft_site_user`

Creates and manages a user for a given site. Only the most common account fields are currently exposed; additional attributes can be added as needed.

```hcl
resource "globalscapeeft_site_user" "example" {
  site_id         = "892b16dc-24a8-473f-a74e-c597b824c879"
  login_name      = "terraform-user"
  password        = var.user_password
  password_type   = "Default"
  display_name    = "Terraform Automation"
  email           = "automation@example.com"
  account_enabled = "yes"
}
```

### Resource `globalscapeeft_event_rule`

Allows you to manage an event rule using the raw JSON `attributes`/`relationships` payloads from the EFT REST API. This is useful when importing an existing rule, tweaking it, and removing it once no longer needed.

```hcl
resource "globalscapeeft_event_rule" "auto_cleanup" {
  site_id = data.globalscapeeft_sites.all.sites[0].id

  attributes_json = jsonencode({
    info = {
      Name        = "Terraform Rule"
      Description = "Cleanup automation"
      Enabled     = true
      Type        = "Timer"
      Folder      = "00000000-0000-0000-0000-000000000000"
      Remote      = false
    }
    statements = {
      StatementsList = []
    }
  })
}
```

## Examples

See the `examples/` directory for copy/paste ready snippets covering provider configuration, data sources, and resources.
