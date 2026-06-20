data "openrouter_model" "claude" {
  id = "anthropic/claude-3-haiku"
}

output "claude_context_length" {
  value = data.openrouter_model.claude.context_length
}

output "claude_pricing" {
  value = data.openrouter_model.claude.pricing
}
