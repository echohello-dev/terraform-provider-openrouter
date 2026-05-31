# Multi-Model Comparison Example
# Compare responses from multiple AI models

terraform {
  required_providers {
    openrouter = {
      source  = "echohello-dev/terraform-provider-openrouter"
      version = "~> 1.0"
    }
  }
}

provider "openrouter" {
  api_key = var.openrouter_api_key
}

variable "openrouter_api_key" {
  type      = string
  sensitive = true
}

variable "prompt" {
  type    = string
  default = "Explain what Terraform is in one paragraph."
}

variable "models" {
  type    = list(string)
  default = ["openai/gpt-3.5-turbo", "anthropic/claude-3-haiku", "google/gemini-pro"]
}

locals {
  completion_for_model = { for i, model_id in var.models : model_id => openrouter_chat_completion.comparisons[i].content }
}

resource "openrouter_chat_completion" "comparisons" {
  count = length(var.models)

  model = var.models[count.index]

  messages {
    role    = "user"
    content = var.prompt
  }

  max_tokens = 300
  temperature = 0.7
}

output "model_responses" {
  value = local.completion_for_model
}
