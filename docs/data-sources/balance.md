---
subcategory: "Account"
layout: ""
page_title: "OpenRouter Balance Data Source - terraform-provider-openrouter"
description: |-
  Get account balance and usage statistics from OpenRouter.
---

# openrouter_balance

Use this data source to view your OpenRouter account balance and usage statistics.

## Example Usage

```hcl
data "openrouter_balance" "my_balance" {}

output "credits_remaining" {
  value = data.openrouter_balance.my_balance.limit_remaining
}

output "monthly_usage" {
  value = data.openrouter_balance.my_balance.usage_monthly
}
```

## Attributes Reference

- `label` - Label for the API key
- `limit` - Credit limit for the key, or null if unlimited
- `limit_remaining` - Remaining credits for the key, or null if unlimited
- `usage` - Number of credits used (all time)
- `usage_daily` - Number of credits used (current UTC day)
- `usage_weekly` - Number of credits used (current UTC week, starting Monday)
- `usage_monthly` - Number of credits used (current UTC month)
- `is_free_tier` - Whether the user has paid for credits before
