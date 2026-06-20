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
}

output "response" {
  value = openrouter_chat_completion.example.content
}
