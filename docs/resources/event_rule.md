---
page_title: "Globalscape EFT: event_rule Resource"
description: |-
  Manage Globalscape EFT event rules via raw JSON payloads.
---

# Resource `globalscapeeft_event_rule`

This resource creates, updates, and deletes event rules for a site using the `/admin/v2/sites/{siteId}/event-rules` REST endpoints. It exposes the REST `attributes` and `relationships` bodies as JSON strings so you can paste definitions straight from the EFT API reference or an existing rule.

**Important Notes:**
- Sensitive fields (passwords, passphrases) in the JSON are automatically sanitized and not stored in Terraform state for security.
- The JSON is normalized when stored in state, so formatting differences are expected.

## Example Usage

```hcl
resource "globalscapeeft_event_rule" "timer" {
  site_id = var.site_id

  attributes_json = jsonencode({
    info = {
      Name        = "Terraform Timer"
      Description = "Runs a timer"
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

## Import

Import an existing rule using `<site_id>/<rule_id>`:

```bash
terraform import globalscapeeft_event_rule.example 5ceae6e3-11b1-40c6-b4e4-3078a8e88a35/8d11ec4f-3c6c-4ab9-8045-fb5c94b38ca0
```

## Schema

### Required

- `site_id` (String) ID of the site where the rule lives.
- `attributes_json` (String) JSON document for the event rule `attributes` block.

### Optional

- `relationships_json` (String) JSON document for the `relationships` block, when needed.

### Timeouts

This resource supports customizable timeouts for operations:
- `create` - Default: 5 minutes
- `read` - Default: 5 minutes
- `update` - Default: 5 minutes
- `delete` - Default: 5 minutes

Example:
```hcl
resource "globalscapeeft_event_rule" "timer" {
  # ... other configuration ...

  timeouts {
    create = "10m"
    update = "10m"
  }
}
```

### Read-only

- `id` (String) Event rule identifier assigned by EFT.
