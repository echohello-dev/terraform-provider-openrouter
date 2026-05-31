# Basic Usage Example
# Lists available models and creates a simple chat completion

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

# Get all available text models
data "openrouter_models" "available" {
  output_modality = "text"
}

# Get specific model details
data "openrouter_model" "claude" {
  id = "anthropic/claude-3-haiku"
}

# Check account balance
data "openrouter_balance" "my_balance" {}

output "available_models_count" {
  value = length(data.openrouter_models.available.models)
}

output "claude_context_length" {
  value = data.openrouter_model.claude.context_length
}

output "credits_remaining" {
  value     = data.openrouter_balance.my_balance.limit_remaining
  sensitive = true
}

# Create a simple chat completion
resource "openrouter_chat_completion" "hello" {
  model = data.openrouter_model.claude.id

  messages {
    role    = "user"
    content = "Say hello in one sentence."
  }

  max_tokens = 50
}

output "hello_response" {
  value = openrouter_chat_completion.hello.content
}
