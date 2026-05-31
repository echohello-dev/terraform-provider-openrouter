---
subcategory: "AI"
layout: ""
page_title: "OpenRouter Chat Completion Resource - terraform-provider-openrouter"
description: |-
  Create and manage chat completions with OpenRouter AI models.
---

# openrouter_chat_completion

Use this resource to create chat completions and store their results in Terraform state.

## Example Usage

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
  temperature   = 0.7
}

output "response" {
  value = openrouter_chat_completion.example.content
}
```

## Argument Reference

- `model` - (Required) The model to use for the completion. Example: `openai/gpt-4`.
- `messages` - (Required) List of message objects. Each message must have `role` and `content`.
- `max_tokens` - (Optional) Maximum number of tokens to generate.
- `temperature` - (Optional) Sampling temperature (0-2). Higher values make output more random. Default: `0.7`.
- `top_p` - (Optional) Nucleus sampling parameter. Higher values make output more focused.

### Messages

Each message object supports the following:

- `role` - (Required) Role of the message. Must be one of: `system`, `user`, `assistant`.
- `content` - (Required) Content of the message.
- `name` - (Optional) Name of the sender.

## Attributes Reference

- `response_id` - Unique identifier for this completion.
- `response_object` - Object type (e.g., `chat.completion`).
- `response_created` - Unix timestamp of when the response was created.
- `response_model` - Actual model used for the completion.
- `finish_reason` - Reason the completion finished (stop, length, etc.).
- `content` - The generated completion content.
- `prompt_tokens` - Number of tokens in the prompt.
- `completion_tokens` - Number of tokens in the completion.
- `total_tokens` - Total number of tokens used.

## Import

Chat completions can be imported by their response ID:

```bash
terraform import openrouter_chat_completion.example <response_id>
```

## Edge Cases

- **Empty response**: If the model returns an empty completion, `content` will be an empty string.
- **Long responses**: If the completion exceeds `max_tokens`, `finish_reason` will be `"length"`.
- **Rate limiting**: The provider will retry on 429 errors with exponential backoff.
- **Insufficient credits**: Returns a 402 error with a descriptive message.
