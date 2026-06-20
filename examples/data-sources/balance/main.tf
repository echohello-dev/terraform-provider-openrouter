data "openrouter_balance" "current" {}

output "label" {
  value = data.openrouter_balance.current.label
}

output "limit_remaining" {
  value     = data.openrouter_balance.current.limit_remaining
  sensitive = true
}
