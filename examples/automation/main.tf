# Automated Documentation Generator
# Generate documentation from project description

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

variable "project_description" {
  type = object({
    name        = string
    description = string
    language    = string
    features    = list(string)
    endpoints   = list(string)
  })

  default = {
    name        = "user-management-api"
    description = "RESTful API for user authentication and profile management"
    language    = "Go"
    features = [
      "JWT token authentication",
      "User registration and login",
      "Password reset via email",
      "Profile CRUD operations",
      "Role-based access control"
    ]
    endpoints = [
      "POST /auth/register",
      "POST /auth/login",
      "POST /auth/refresh",
      "GET /users/:id",
      "PUT /users/:id",
      "DELETE /users/:id"
    ]
  }
}

resource "openrouter_chat_completion" "readme" {
  model = "openai/gpt-4"

  messages {
    role    = "system"
    content = <<-EOF
      You are an expert technical writer. Generate comprehensive README documentation.
      Include: Project title, description, features, prerequisites, installation,
      configuration, API endpoints, environment variables, and contributing guidelines.
      Format output in Markdown.
    EOF
  }

  messages {
    role    = "user"
    content = <<-EOF
      Generate README for: ${var.project_description.name}
      Description: ${var.project_description.description}
      Language: ${var.project_description.language}
      Features: ${join(", ", var.project_description.features)}
      Endpoints: ${join(", ", var.project_description.endpoints)}
    EOF
  }

  max_tokens   = 2000
  temperature   = 0.5
}

resource "local_file" "readme" {
  content  = openrouter_chat_completion.readme.content
  filename = "${path.module}/generated_README.md"
}

output "generated_readme" {
  value = openrouter_chat_completion.readme.content
}
