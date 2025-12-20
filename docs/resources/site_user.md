---
page_title: "Globalscape EFT: site_user Resource"
description: |-
  Manage Globalscape EFT site users via the REST API.
---

# Resource `globalscapeeft_site_user`

Creates, updates, and deletes a Globalscape EFT user scoped to a particular site using the `/admin/v2/sites/{siteId}/users` endpoints.

Only a subset of the available account attributes are modeled today (login name, password, account enablement, and basic personal information). Additional fields can be added to the resource as needed.

## Example Usage

```hcl
resource "globalscapeeft_site_user" "example" {
  site_id            = "892b16dc-24a8-473f-a74e-c597b824c879"
  login_name         = "terraform-user"
  password           = var.user_password
  password_type      = "Default"
  display_name       = "Terraform Automation"
  email              = "automation@example.com"
  account_enabled    = "yes"
  home_folder_path   = "/Automation"
  home_folder_enabled = "yes"
  home_folder_root   = "yes"
}
```

## Schema

### Required

- `site_id` (String) ID of the site that owns the user.
- `login_name` (String) Unique login name for the user.

### Optional

- `password` (String, Sensitive) Password for local EFT accounts.
- `password_type` (String) Password type value expected by EFT (defaults to `Default`).
- `display_name` (String) Friendly display name.
- `email` (String) Email address.
- `account_enabled` (String) Whether the account is enabled (`yes`, `no`, or `inherit`). Defaults to `inherit`.
- `home_folder_path` (String) Physical path for the user's home directory.
- `home_folder_enabled` (String) Enables or disables the home folder entry (`yes`, `no`, or `inherit`). Defaults to `inherit`.
- `home_folder_root` (String) Whether the home folder is treated as the user's root (`yes`, `no`, or `inherit`). Defaults to `inherit`.

### Read-only

- `id` (String) User GUID assigned by EFT.
