data "openrouter_credits" "account" {}

output "credits_purchased" {
  value = data.openrouter_credits.account.total_credits
}

output "credits_used" {
  value = data.openrouter_credits.account.total_usage
}

output "credits_remaining" {
  value = data.openrouter_credits.account.total_credits - data.openrouter_credits.account.total_usage
}
