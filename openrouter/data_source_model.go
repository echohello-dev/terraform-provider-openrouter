package openrouter

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceModel() *schema.Resource {
	return &schema.Resource{
		ReadContext: modelRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The model identifier (e.g., openai/gpt-4).",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Display name of the model.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Description of the model.",
			},
			"context_length": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Maximum context length in tokens.",
			},
			"created": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Unix timestamp of when the model was created.",
			},
		"pricing": {
			Type:        schema.TypeMap,
			Computed:    true,
			Elem:        &schema.Schema{Type: schema.TypeFloat},
			Description: "Pricing per million tokens, keyed by `<tier>_input` and `<tier>_output` (e.g. `text_input`, `text_output`).",
		},
			"top_provider": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Top provider for this model.",
			},
			"architecture": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Model architecture.",
			},
			"supported_parameters": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of supported parameters for this model.",
			},
			"knowledge_cutoff": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date up to which the model was trained on data.",
			},
			"features": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Feature flags for the model.",
			},
		},
	}
}

func modelRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client)

	modelID := d.Get("id").(string)

	models, diags := client.ListModels(ctx, "")
	if diags.HasError() {
		return diags
	}

	var foundModel *Model
	for i := range models {
		if models[i].ID == modelID {
			foundModel = &models[i]
			break
		}
	}

	if foundModel == nil {
		return diag.Errorf("Model not found: %s. Available models can be queried using the openrouter_models data source.", modelID)
	}

	pricing := make(map[string]interface{})
	for tier, price := range foundModel.Pricing {
		pricing[tier+"_input"] = price.Input
		pricing[tier+"_output"] = price.Output
	}

	d.SetId(foundModel.ID)
	if err := d.Set("name", foundModel.Name); err != nil {
		return diag.Errorf("Failed to set name: %s", err)
	}
	if err := d.Set("description", foundModel.Description); err != nil {
		return diag.Errorf("Failed to set description: %s", err)
	}
	if err := d.Set("context_length", foundModel.ContextLength); err != nil {
		return diag.Errorf("Failed to set context_length: %s", err)
	}
	if err := d.Set("created", foundModel.Created); err != nil {
		return diag.Errorf("Failed to set created: %s", err)
	}
	if err := d.Set("pricing", pricing); err != nil {
		return diag.Errorf("Failed to set pricing: %s", err)
	}
	if err := d.Set("top_provider", foundModel.TopProvider.SelectedMode); err != nil {
		return diag.Errorf("Failed to set top_provider: %s", err)
	}
	if err := d.Set("architecture", foundModel.Architecture.ModelClass); err != nil {
		return diag.Errorf("Failed to set architecture: %s", err)
	}
	if err := d.Set("supported_parameters", foundModel.SupportedParams); err != nil {
		return diag.Errorf("Failed to set supported_parameters: %s", err)
	}
	if err := d.Set("knowledge_cutoff", foundModel.KnowledgeCutoff); err != nil {
		return diag.Errorf("Failed to set knowledge_cutoff: %s", err)
	}
	if err := d.Set("features", foundModel.Features); err != nil {
		return diag.Errorf("Failed to set features: %s", err)
	}

	return nil
}
