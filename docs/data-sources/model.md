---
subcategory: "Models"
layout: ""
page_title: "OpenRouter Model Data Source - terraform-provider-openrouter"
description: |-
  Get details for a specific AI model from OpenRouter.
---

# openrouter_model

Use this data source to get details for a specific AI model available through OpenRouter.

## Example Usage

```hcl
data "openrouter_model" "gpt4" {
  id = "openai/gpt-4"
}

output "model_context_length" {
  value = data.openrouter_model.gpt4.context_length
}
```

## Argument Reference

- `id` - (Required) The model identifier (e.g., `openai/gpt-4`).

## Attributes Reference

- `id` - Model identifier
- `name` - Display name of the model
- `description` - Description of the model
- `context_length` - Maximum context length in tokens
- `created` - Unix timestamp of when the model was created
- `pricing` - Map of pricing tiers to input prices
- `top_provider` - Top provider for this model
- `architecture` - Model architecture class
- `supported_parameters` - List of supported API parameters
