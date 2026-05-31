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
				Description: "Pricing information for the model (input price per million tokens).",
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
		pricing[tier] = price.Input
	}

	d.SetId(foundModel.ID)
	d.Set("name", foundModel.Name)
	d.Set("description", foundModel.Description)
	d.Set("context_length", foundModel.ContextLength)
	d.Set("created", foundModel.Created)
	d.Set("pricing", pricing)
	d.Set("top_provider", foundModel.TopProvider.SelectedMode)
	d.Set("architecture", foundModel.Architecture.ModelClass)
	d.Set("supported_parameters", foundModel.SupportedParams)
	d.Set("knowledge_cutoff", foundModel.KnowledgeCutoff)
	d.Set("features", foundModel.Features)

	return nil
}
