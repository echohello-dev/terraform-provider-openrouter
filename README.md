# Terraform Provider for OpenRouter

A Terraform provider for managing [OpenRouter](https://openrouter.ai) resources. OpenRouter provides unified API access to 300+ AI models from providers like Anthropic, OpenAI, Meta, Google, and more.

## Features

- **Data Sources**
  - `openrouter_models` - List all available AI models
  - `openrouter_model` - Get details for a specific model
  - `openrouter_balance` - View account credit balance and usage
  - `openrouter_generation` - Query generation stats and cost by ID
  - `openrouter_credits` - View total credits purchased and used

- **Resources**
  - `openrouter_chat_completion` - Create chat completions and store results in state

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- Go >= 1.25 (for building from source)

## Installation

### Pre-built Binaries

Download the latest release for your platform from the [GitHub Releases](https://github.com/echohello-dev/terraform-provider-openrouter/releases) page.

### From Source

```bash
git clone https://github.com/echohello-dev/terraform-provider-openrouter.git
cd terraform-provider-openrouter
go build -o terraform-provider-openrouter
```

### Terraform Registry

The provider is available on the Terraform Registry. Add to your Terraform configuration:

```hcl
terraform {
  required_providers {
    openrouter = {
      source  = "echohello-dev/openrouter"
      version = "~> 1.0"
    }
  }
}
```

## Configuration

### Provider Configuration Options

| Option | Type | Required | Default | Description |
|--------|------|----------|---------|-------------|
| `api_key` | string | Yes | - | OpenRouter API key |
| `referer` | string | No | - | Site URL for OpenRouter rankings |
| `app_title` | string | No | - | App display name for OpenRouter dashboard |

### Environment Variables

```bash
export OPENROUTER_API_KEY="sk-..."
export OPENROUTER_REFERER="https://your-app.com"
export OPENROUTER_APP_TITLE="Your App Name"
```

## Usage Examples

### List Available Models

```hcl
data "openrouter_models" "available" {
  output_modality = "text"  # text, image, audio, embeddings, all
}

output "model_list" {
  value = data.openrouter_models.available.models[*].id
}
```

### Get Specific Model Details

```hcl
data "openrouter_model" "gpt4" {
  id = "openai/gpt-4"
}

output "model_info" {
  value = data.openrouter_model.gpt4
}
```

### Check Account Balance

```hcl
data "openrouter_balance" "my_balance" {}

output "credits_remaining" {
  value = data.openrouter_balance.my_balance.limit_remaining
}
```

### Query Generation Stats and Cost

Use the generation ID returned by `openrouter_chat_completion` to audit cost, latency, and provider details after the fact.

```hcl
resource "openrouter_chat_completion" "example" {
  model = "openai/gpt-4"
  # ... messages
}

data "openrouter_generation" "cost_audit" {
  id = openrouter_chat_completion.example.response_id
}

output "generation_cost" {
  value = data.openrouter_generation.cost_audit.total_cost
}

output "generation_provider" {
  value = data.openrouter_generation.cost_audit.provider_name
}
```

### View Total Credits

```hcl
data "openrouter_credits" "account" {}

output "credits_purchased" {
  value = data.openrouter_credits.account.total_credits
}

output "credits_used" {
  value = data.openrouter_credits.account.total_usage
}
```

### Create a Chat Completion

```hcl
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
}

output "response" {
  value = resource.openrouter_chat_completion.example.content
}
```

### Complete Example: AI-Powered README Generator

```hcl
provider "openrouter" {
  api_key = var.openrouter_api_key
}

variable "openrouter_api_key" {
  type      = string
  sensitive = true
}

data "openrouter_model" "claude" {
  id = "anthropic/claude-3-opus"
}

resource "openrouter_chat_completion" "readme_generator" {
  model = data.openrouter_model.claude.id

  messages {
    role    = "system"
    content = "You are an expert technical writer. Generate clear, concise README documentation."
  }

  messages {
    role    = "user"
    content = <<-EOF
      Generate a README for a Go project with the following structure:
      - Project name: my-api-service
      - Description: RESTful API for user management
      - Features: JWT auth, PostgreSQL, Redis caching
      - Endpoints: /users, /auth/login, /auth/register
    EOF
  }

  max_tokens = 1000
  temperature = 0.5
}

output "generated_readme" {
  value = openrouter_chat_completion.readme_generator.content
}
```

## Import Support

The `openrouter_chat_completion` resource is create-only. There is no OpenRouter
endpoint to fetch a previously-created completion by ID, so importing existing
completions into Terraform state is not supported. Recreate the resource in
config and the API call will be made again on the next `terraform apply`.

## Edge Cases and Error Handling

### Rate Limiting

The provider handles 429 Too Many Requests errors and retries with exponential backoff.

### Insufficient Credits

If your account has insufficient credits (402 errors), the provider will return a descriptive error:

```
Error: API error (402): Insufficient credits. Please add credits to your OpenRouter account.
```

### Model Not Found

When querying a specific model that doesn't exist:

```
Error: Model not found: invalid/model-id
```

### Invalid Messages

Message roles must be one of: `system`, `user`, `assistant`, or `tool`.

## Automation Patterns

### Automated Model Selection Based on Cost

```hcl
data "openrouter_models" "all" {}

locals {
  # Find cheapest text model under $0.01/1M tokens
  cheapest_model = merge([
    for m in data.openrouter_models.all.models :
    m if m.pricing.input < 0.01
  ]...)[0]
}

resource "openrouter_chat_completion" "budget_completion" {
  model = local.cheapest_model.id
  # ... messages
}
```

### Batch Completion Processing

```hcl
variable "prompts" {
  type    = list(string)
  default = ["What is AI?", "What is ML?", "What is DL?"]
}

resource "openrouter_chat_completion" "batch" {
  count = length(var.prompts)

  model = "openai/gpt-3.5-turbo"

  messages {
    role    = "user"
    content = var.prompts[count.index]
  }

  max_tokens = 200
}

output "responses" {
  value = openrouter_chat_completion.batch[*].content
}
```

### Dynamic Model Selection with Fallback

```hcl
resource "openrouter_chat_completion" "primary" {
  model = "anthropic/claude-3-opus"

  messages {
    role    = "user"
    content = "Explain quantum computing"
  }
}

# If primary fails, fallback to cheaper model
resource "null_resource" "fallback" {
  count = openrouter_chat_completion.primary.id == "" ? 1 : 0

  # Trigger fallback logic via local-exec
  provisioner "local-exec" {
    command = "echo 'Primary model failed, consider using gpt-3.5-turbo'"
  }
}
```

## Migration Guide

### From v0.x to v1.0

If you were using an older version of this provider:

1. Update provider source in your Terraform configuration:
   ```hcl
   terraform {
     required_providers {
       openrouter = {
         source  = "echohello-dev/openrouter"
         version = "~> 1.0"
       }
     }
   }
   ```

2. Rename any existing resources:
   - `openrouter_completion` → `openrouter_chat_completion`

3. Update message format if using completions:
   ```hcl
   # Old format
   messages = ["prompt text"]

   # New format
   messages {
     role    = "user"
     content = "prompt text"
   }
   ```

## Debugging

Enable debug logging:

```bash
TF_LOG=DEBUG terraform plan
```

## Contributing

Contributions are welcome! Please read our contributing guidelines before submitting PRs.

## License

MIT License - see LICENSE file for details.
