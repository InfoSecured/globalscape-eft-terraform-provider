---
page_title: "Globalscape EFT: server Data Source"
description: |-
  Reads high-level Globalscape EFT server settings via GET /admin/v2/server.
---

# Data Source `globalscapeeft_server`

Fetches server metadata, listener details, and SMTP configuration using the REST API documented in the bundled Globalscape EFT reference PDF.

## Example Usage

```hcl
data "globalscapeeft_server" "current" {}

output "version" {
  value = data.globalscapeeft_server.current.version
}
```

## Schema

### Read-only

- `id` (String) Internal server identifier returned by EFT.
- `version` (String) Software version string.
- `general` (Block)
  - `config_file_path` (String) File path to the EFT configuration.
  - `enable_utc_in_listings` (Boolean) Whether UTC timestamps are shown in listings.
  - `last_modified_by` (String) Most recent admin to change the config.
  - `last_modified_time` (Number) Unix timestamp for the last modification.
- `listener_settings` (Block)
  - `admin_port` (Number) Administrative listener port.
  - `enable_remote_administration` (Boolean) Whether the remote admin service is enabled.
  - `listen_ips` (List of String) IP addresses where the admin service listens.
- `smtp` (Block)
  - `login` (String)
  - `password` (String, Sensitive)
  - `port` (Number)
  - `sender_address` (String)
  - `sender_name` (String)
  - `server` (String)
  - `use_authentication` (Boolean)
  - `use_implicit_tls` (Boolean)
