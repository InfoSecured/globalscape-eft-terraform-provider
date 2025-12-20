page_title: "Globalscape EFT Provider"
description: |-
  Provider plugin for managing Globalscape EFT via REST.
---

# Globalscape EFT Provider

Use this provider to interact with Globalscape EFT Admin REST endpoints exposed at `https://<host>:<port>/admin`.

## Authentication

Provide credentials for an EFT local admin or an AD account that is authorized for the Admin API. The provider exchanges the credentials for an `EFTAdminAuthToken` via `POST /admin/v1/authentication` and attaches this token to all follow-up requests.

```hcl
provider "globalscapeeft" {
  host                 = "https://eft.example.com:4450/admin"
  username             = var.eft_username
  password             = var.eft_password
  auth_type            = "EFT" # or "AD"
  insecure_skip_verify = true   # optional, useful for lab systems
}
```

## Schema

- `host` (String, Required) Admin API base URL including the `/admin` suffix.
- `username` (String, Required) Admin account username.
- `password` (String, Required, Sensitive) Admin account password.
- `auth_type` (String, Optional) Authentication realm. Defaults to `EFT`.
- `insecure_skip_verify` (Boolean, Optional) Skip TLS verification when connecting to EFT.

## Supported Resources

- [`globalscapeeft_server_smtp`](resources/server_smtp.md)
- [`globalscapeeft_site_user`](resources/site_user.md)
- [`globalscapeeft_event_rule`](resources/event_rule.md)

## Supported Data Sources

- [`globalscapeeft_server`](data-sources/server.md)
- [`globalscapeeft_sites`](data-sources/sites.md)
