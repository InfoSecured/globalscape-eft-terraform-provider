---
page_title: "Globalscape EFT: sites Data Source"
description: |-
  Lists the sites configured on the EFT server.
---

# Data Source `globalscapeeft_sites`

Use this data source to enumerate Globalscape EFT sites and discover the IDs required by other resources (such as `globalscapeeft_site_user`).

## Example Usage

```hcl
data "globalscapeeft_sites" "all" {}

output "first_site_id" {
  value = data.globalscapeeft_sites.all.sites[0].id
}
```

## Schema

### Read-only

- `sites` (List of Object) List of sites returned by the EFT API.
  - `id` (String) Site identifier used in REST endpoints.
  - `name` (String) Site label configured on the server.
