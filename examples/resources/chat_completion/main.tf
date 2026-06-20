terraform {
  required_providers {
    openrouter = {
      source  = "echohello-dev/openrouter"
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

resource "openrouter_chat_completion" "example" {
  model = "openai/gpt-4"

  messages {
    role    = "system"
    content = "You are a helpful assistant."
  }

  messages {
    role    = "user"
    content = "What is Terraform?"
  }

  max_tokens   = 500
  temperature  = 0.7
  seed         = 42
  user         = "terraform-user"
  session_id   = "session-123"
}

output "response" {
  value = openrouter_chat_completion.example.content
}

# Audit the cost of the generation after the fact
data "openrouter_generation" "cost_audit" {
  id = openrouter_chat_completion.example.response_id
}

output "generation_cost" {
  value = data.openrouter_generation.cost_audit.total_cost
}

output "generation_provider" {
  value = data.openrouter_generation.cost_audit.provider_name
}

output "generation_latency_ms" {
  value = data.openrouter_generation.cost_audit.latency
}
