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
