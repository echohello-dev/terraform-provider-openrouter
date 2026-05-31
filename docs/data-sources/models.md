---
subcategory: "Models"
layout: ""
page_title: "OpenRouter Models Data Source - terraform-provider-openrouter"
description: |-
  List all available AI models from OpenRouter.
---

# openrouter_models

Use this data source to get a list of all AI models available through OpenRouter.

## Example Usage

```hcl
data "openrouter_models" "all" {}

data "openrouter_models" "text_only" {
  output_modality = "text"
}
```

## Argument Reference

- `output_modality` - (Optional) Filter models by output capability. Valid values: `text`, `image`, `audio`, `embeddings`, `all`. Defaults to `text`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `models` - A list of model objects with the following attributes:
  - `id` - Unique identifier for the model (e.g., `openai/gpt-4`)
  - `name` - Display name of the model
  - `description` - Description of the model
  - `context_length` - Maximum context length in tokens
  - `created` - Unix timestamp of when the model was created
  - `pricing` - Map of pricing tiers to input prices
  - `top_provider` - Top provider for this model
  - `architecture` - Model architecture class
