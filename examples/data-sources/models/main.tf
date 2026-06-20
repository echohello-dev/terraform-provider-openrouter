data "openrouter_models" "text" {
  output_modality = "text"
}

output "text_model_ids" {
  value = data.openrouter_models.text.models[*].id
}
